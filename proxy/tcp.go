package proxy

import (
	"net"
	"strconv"

	"github.com/net-byte/qsocks/common/constant"
	"github.com/net-byte/qsocks/config"
)

func TCPProxy(conn net.Conn, config config.Config, data []byte) {
	host, port := getAddr(data)
	if host == "" || port == "" {
		return
	}
	// bypass private ip
	if config.Bypass && net.ParseIP(host) != nil && net.ParseIP(host).IsPrivate() {
		DirectProxy(conn, host, port, config)
		return
	}

	stream := ConnectServer("tcp", host, port, config)
	if stream == nil {
		ResponseTCP(conn, constant.ConnectionRefused)
		return
	}

	ResponseTCP(conn, constant.SuccessReply)
	go Copy(stream, conn)
	go Copy(conn, stream)

}

func getAddr(b []byte) (host string, port string) {
	/**
	  +----+-----+-------+------+----------+----------+
	  |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	  +----+-----+-------+------+----------+----------+
	  | 1  |  1  | X'00' |  1   | Variable |    2     |
	  +----+-----+-------+------+----------+----------+
	*/
	len := len(b)
	switch b[3] {
	case constant.Ipv4Address:
		host = net.IPv4(b[4], b[5], b[6], b[7]).String()
	case constant.FqdnAddress:
		host = string(b[5 : len-2])
	case constant.Ipv6Address:
		host = net.IP(b[4:19]).String()
	default:
		return "", ""
	}
	port = strconv.Itoa(int(b[len-2])<<8 | int(b[len-1]))
	return host, port
}
