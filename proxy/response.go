package proxy

import (
	"bytes"
	"encoding/binary"
	"net"

	"github.com/net-byte/qsocks/common/enum"
)

func Resp(conn net.Conn, rep byte) {
	/**
	  +----+-----+-------+------+----------+----------+
	  |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	  +----+-----+-------+------+----------+----------+
	  | 1  |  1  | X'00' |  1   | Variable |    2     |
	  +----+-----+-------+------+----------+----------+
	*/
	conn.Write([]byte{enum.Socks5Version, rep, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
}

func RespNoAuth(conn net.Conn) {
	/**
	  +----+--------+
	  |VER | METHOD |
	  +----+--------+
	  | 1  |   1    |
	  +----+--------+
	*/
	conn.Write([]byte{enum.Socks5Version, enum.NoAuth})
}

func RespSuccess(conn net.Conn, bindAddr *net.UDPAddr) {
	/**
	  +----+-----+-------+------+----------+----------+
	  |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	  +----+-----+-------+------+----------+----------+
	  | 1  |  1  | X'00' |  1   | Variable |    2     |
	  +----+-----+-------+------+----------+----------+
	*/
	response := []byte{enum.Socks5Version, enum.SuccessReply, 0x00, 0x01}
	buffer := bytes.NewBuffer(response)
	binary.Write(buffer, binary.BigEndian, bindAddr.IP.To4())
	binary.Write(buffer, binary.BigEndian, uint16(bindAddr.Port))
	conn.Write(buffer.Bytes())
}
