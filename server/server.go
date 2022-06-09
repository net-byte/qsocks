package server

import (
	"bufio"
	"context"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/net-byte/qsocks/common/enum"
	"github.com/net-byte/qsocks/config"
	"github.com/net-byte/qsocks/proto"
	"github.com/net-byte/qsocks/proxy"
)

// Start starts the server
func Start(config config.Config) {
	log.Printf("qsocks server started on %s", config.ServerAddr)
	tlsConf, err := config.GetServerTLSConfig()
	if err != nil {
		log.Panic(err)
	}
	l, err := quic.ListenAddr(config.ServerAddr, tlsConf, nil)
	if err != nil {
		log.Panic(err)
	}
	for {
		session, err := l.Accept(context.Background())
		if err != nil {
			continue
		}
		go proxyConn(session, config)
	}

}

func proxyConn(session quic.Session, config config.Config) {
	stream, err := session.AcceptStream(context.Background())
	if err != nil {
		log.Println(err)
		return
	}
	defer stream.Close()
	reader := bufio.NewReader(stream)
	// handshake
	ok, req := handshake(config, reader)
	if !ok {
		return
	}
	// connect dst server
	conn, err := net.DialTimeout(req.Network, net.JoinHostPort(req.Host, req.Port), time.Duration(enum.Timeout)*time.Second)
	if err != nil {
		log.Printf("[server] failed to dial the dst server %v", err)
		return
	}
	defer conn.Close()
	// forward
	go toServer(config, reader, conn)
	toClient(config, stream, conn)
}

func handshake(config config.Config, reader *bufio.Reader) (bool, proxy.RequestAddr) {
	var req proxy.RequestAddr
	buf, _, err := proto.Decode(reader)
	if err != nil {
		return false, req
	}
	if req.UnmarshalBinary(buf) != nil {
		log.Printf("[server] failed to decode request addr %v", err)
		return false, req
	}
	reqTime, _ := strconv.ParseInt(req.Timestamp, 10, 64)
	if time.Now().Unix()-reqTime > int64(enum.Timeout) {
		log.Printf("[server] timestamp expired %v", reqTime)
		return false, req
	}
	return true, req
}

func toClient(config config.Config, stream quic.Stream, conn net.Conn) {
	buf := make([]byte, enum.BufferSize)
	for {
		conn.SetReadDeadline(time.Now().Add(time.Duration(enum.Timeout) * time.Second))
		n, err := conn.Read(buf)
		if err != nil || n == 0 {
			break
		}
		b := buf[:n]
		_, err = stream.Write(b)
		if err != nil {
			break
		}
	}
}

func toServer(config config.Config, reader *bufio.Reader, conn net.Conn) {
	buf := make([]byte, enum.BufferSize)
	for {
		n, err := reader.Read(buf)
		if err != nil || n == 0 {
			break
		}
		b := buf[:n]
		_, err = conn.Write(b)
		if err != nil {
			break
		}
	}
}
