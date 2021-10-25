package proxy

import (
	"bytes"
	"encoding/binary"
	"net"

	"github.com/net-byte/qsocks/common/constant"
)

func ResponseTCP(conn net.Conn, rep byte) {
	/**
	  +----+-----+-------+------+----------+----------+
	  |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	  +----+-----+-------+------+----------+----------+
	  | 1  |  1  | X'00' |  1   | Variable |    2     |
	  +----+-----+-------+------+----------+----------+
	*/
	conn.Write([]byte{constant.Socks5Version, rep, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
}

func ResponseNoAuth(conn net.Conn) {
	/**
	  +----+--------+
	  |VER | METHOD |
	  +----+--------+
	  | 1  |   1    |
	  +----+--------+
	*/
	conn.Write([]byte{constant.Socks5Version, constant.NoAuth})
}

func ResponseUDP(conn net.Conn, bindAddr *net.UDPAddr) {
	/**
	  +----+-----+-------+------+----------+----------+
	  |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	  +----+-----+-------+------+----------+----------+
	  | 1  |  1  | X'00' |  1   | Variable |    2     |
	  +----+-----+-------+------+----------+----------+
	*/
	response := []byte{constant.Socks5Version, constant.SuccessReply, 0x00, 0x01}
	buffer := bytes.NewBuffer(response)
	binary.Write(buffer, binary.BigEndian, bindAddr.IP.To4())
	binary.Write(buffer, binary.BigEndian, uint16(bindAddr.Port))
	conn.Write(buffer.Bytes())
}
