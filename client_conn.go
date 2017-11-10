package main

import (
	"io"
	"strings"
	"context"
	"net/http"
	"sync"
	"log"
	"net"
	"bufio"
	"github.com/armon/go-socks5"
)

var socksServer *socks5.Server

func init() {
	conf := &socks5.Config{
		Dial: proxyDial,
	}
	var err error
	socksServer, err = socks5.New(conf)
	if err != nil {
	  panic(err)
	}
}

type clientConn struct {
	net.Conn
	r *bufio.Reader
}

func newClientConn(c net.Conn) *clientConn {
	return &clientConn {
		Conn: c,
		r: bufio.NewReader(c),
	}
}

func (cc *clientConn) Read(p []byte) (int, error) {
	return cc.r.Read(p)
}

func (cc *clientConn) serve() {
	defer func(){
		cc.Close()
	}()

	

	wg := &sync.WaitGroup{}

	buf,err:=cc.r.Peek(1)
	if err != nil {
		return
	}

	if buf[0] == uint8(5) {
		cc.serveSocks5(wg)
	} else {
		cc.serveHTTP(wg)
	}

	wg.Wait()
}

func (cc *clientConn) serveSocks5(wg *sync.WaitGroup) {
	wg.Add(1)
	socksServer.ServeConn(cc)
}

func (cc *clientConn) serveHTTP(wg *sync.WaitGroup) {
	req,err:=http.ReadRequest(cc.r)
	if err != nil {
		return
	}
	log.Println(req)

	host := req.Host
	if !strings.Contains(req.Host, ":") {
		host += ":80"
	}

	serverConn, err := proxyDial(context.TODO(), "tcp", host)
	if err != nil {
		return
	}

	if req.Method == "CONNECT" {
		rsp:="HTTP/1.1 200 Connection Established\r\nProxy-Agent: simpless\r\n\r\n"
		_, err = cc.Write([]byte(rsp))
		if err != nil {
			return
		}
	} else {
		err = req.Write(serverConn)
		if err != nil {
			return
		}
	}


	wg.Add(2)
	
	splice := func(r io.ReadCloser, w io.Writer, closeReader bool) {
		defer wg.Done()
		b := make([]byte, 1024)
		for {
			n,err:=r.Read(b)
			if n == 0 || err != nil {
				break
			}
			nn,err:=w.Write(b[:n])
			if n != nn || err != nil {
				break
			}
		}
	}

	go splice(cc, serverConn, false)
	go splice(serverConn, cc, true)
}