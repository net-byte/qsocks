package proxy

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/net-byte/qsocks/common/constant"
	"github.com/net-byte/qsocks/config"
)

func DirectProxy(conn net.Conn, host string, port string, config config.Config) {
	remoteConn := connectTCP(host, port, config)
	if remoteConn == nil {
		ResponseTCP(conn, constant.ConnectionRefused)
		return
	}

	ResponseTCP(conn, constant.SuccessReply)
	go forward(remoteConn, conn)
	go forward(conn, remoteConn)
}

func connectTCP(host string, port string, config config.Config) net.Conn {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), 60*time.Second)
	if err != nil {
		log.Printf("[tcp] failed to dial tcp %v", err)
		return nil
	}
	return conn
}

func forward(to io.WriteCloser, from io.ReadCloser) {
	defer to.Close()
	defer from.Close()
	io.Copy(to, from)
}
