package main

import (	
	ini "gopkg.in/go-ini/ini.v1"
)

var conf  struct {
	parentProxy string
	listenAddr string
}


func parseConfig(file string) {
	cfg,err:=ini.Load(file)
	if err != nil {
		panic(err)
	}
	conf.parentProxy = cfg.Section("").Key("parent").String()
	conf.listenAddr = cfg.Section("").Key("listen").String()
}