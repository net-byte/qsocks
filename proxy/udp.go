package proxy

import (
	"log"
	"net"

	"github.com/net-byte/qsocks/common/constant"
	"github.com/net-byte/qsocks/config"
)

func UDPProxy(tcpConn net.Conn, udpConn *net.UDPConn, config config.Config) {
	defer tcpConn.Close()
	if udpConn == nil {
		log.Printf("[udp] failed to start udp server on %v", config.LocalAddr)
		return
	}
	bindAddr, _ := net.ResolveUDPAddr("udp", udpConn.LocalAddr().String())
	//response to client
	ResponseUDP(tcpConn, bindAddr)
	//keep tcp alive
	done := make(chan bool)
	go keepTCPAlive(tcpConn.(*net.TCPConn), done)
	<-done
}

func keepTCPAlive(tcpConn *net.TCPConn, done chan<- bool) {
	tcpConn.SetKeepAlive(true)
	buf := make([]byte, constant.BufferSize)
	for {
		_, err := tcpConn.Read(buf[0:])
		if err != nil {
			break
		}
	}
	done <- true
}
