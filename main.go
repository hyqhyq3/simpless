package main

import (
	"flag"
	"net"
)

func init() {
	var configFile string
	flag.StringVar(&configFile, "config,c", "simpless.ini", "config file")
	flag.Parse()

	parseConfig(configFile)
}

func main() {
	
	l,err:=net.Listen("tcp", conf.listenAddr)
	if err != nil {
		panic(err)
	}
	for {
		c,err:=l.Accept()
		if err != nil {
			panic(err)
		}
		
		cc:=newClientConn(c)
		go cc.serve()
	}
}