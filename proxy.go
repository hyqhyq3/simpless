package main

import (
	"net/url"
	"context"
	"net"

	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
)

func proxyDial(ctx context.Context, network string, addr string) (c net.Conn, err error) {
	u,err := url.Parse(conf.parentProxy)
	if err != nil {
		return 
	}
	p,_ := u.User.Password()
	cipher,err:=ss.NewCipher(u.User.Username(), p)
	if err != nil {
		return
	}
	return ss.Dial(addr, u.Host, cipher)
}