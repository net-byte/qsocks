package proxy

import (
	"net"
	"time"

	"github.com/net-byte/qsocks/common/enum"
	"github.com/net-byte/qsocks/config"
)

func DirectProxy(conn net.Conn, host string, port string, config config.Config) {
	rconn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), time.Duration(enum.Timeout)*time.Second)
	if err != nil {
		Resp(conn, enum.ConnectionRefused)
		return
	}
	Resp(conn, enum.SuccessReply)
	go copy(rconn, conn)
	copy(conn, rconn)
}
