package main

import (
	"flag"

	"github.com/net-byte/qsocks/client"
	"github.com/net-byte/qsocks/config"
	"github.com/net-byte/qsocks/server"
)

func main() {
	config := config.Config{}
	flag.StringVar(&config.LocalAddr, "l", "127.0.0.1:1083", "local address")
	flag.StringVar(&config.ServerAddr, "s", ":8443", "server address")
	flag.StringVar(&config.ClientKey, "ck", "../certs/client.key", "client key file path")
	flag.StringVar(&config.ClientPem, "cp", "../certs/client.pem", "client pem file path")
	flag.StringVar(&config.ServerKey, "sk", "../certs/server.key", "server key file path")
	flag.StringVar(&config.ServerPem, "sp", "../certs/server.pem", "server pem file path")
	flag.BoolVar(&config.ServerMode, "S", false, "server mode")
	flag.BoolVar(&config.Bypass, "bypass", false, "bypass private ip")
	flag.Parse()

	if config.ServerMode {
		server.Start(config)
	} else {
		client.Start(config)
	}
}
