package client

import (
	"io"
	"log"
	"net"

	"github.com/net-byte/qsocks/common/enum"
	"github.com/net-byte/qsocks/config"
	"github.com/net-byte/qsocks/proxy"
)

// Start starts the client
func Start(config config.Config) {
	udpConn := startUDPServer(config)
	startTCPServer(config, udpConn)
}

func startTCPServer(config config.Config, udpConn *net.UDPConn) {
	l, err := net.Listen("tcp", config.LocalAddr)
	if err != nil {
		log.Panicf("[tcp] failed to listen tcp %v", err)
	}
	log.Printf("qsocks [tcp] client started on %s", config.LocalAddr)
	for {
		tcpConn, err := l.Accept()
		if err != nil {
			continue
		}
		go tcpHandler(tcpConn, udpConn, config)
	}
}

func startUDPServer(config config.Config) *net.UDPConn {
	udpRelay := &proxy.UDPRelay{Config: config}
	return udpRelay.Start()
}

func tcpHandler(tcpConn net.Conn, udpConn *net.UDPConn, config config.Config) {
	buf := make([]byte, enum.BufferSize)
	//read version
	n, err := tcpConn.Read(buf[0:])
	if err != nil || err == io.EOF {
		return
	}
	b := buf[0:n]
	if b[0] != enum.Socks5Version {
		return
	}
	//no auth
	proxy.RespNoAuth(tcpConn)
	//read cmd
	n, err = tcpConn.Read(buf[0:])
	if err != nil || err == io.EOF {
		return
	}
	b = buf[0:n]
	switch b[1] {
	case enum.ConnectCommand:
		proxy.TCPProxy(tcpConn, config, b)
		return
	case enum.AssociateCommand:
		proxy.UDPProxy(tcpConn, udpConn, config)
		return
	case enum.BindCommand:
		proxy.Resp(tcpConn, enum.CommandNotSupported)
		return
	default:
		proxy.Resp(tcpConn, enum.CommandNotSupported)
		return
	}
}
