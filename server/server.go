package server

import (
	"context"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/net-byte/qsocks/common/constant"
	"github.com/net-byte/qsocks/config"
	"github.com/net-byte/qsocks/proxy"
)

// Starts qsocks server
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

		go handleConn(session, config)
	}

}

func handshake(config config.Config, session quic.Session) (bool, proxy.RequestAddr) {
	var req proxy.RequestAddr
	stream, err := session.AcceptUniStream(context.Background())
	if err != nil {
		log.Println(err)
		return false, req
	}
	buf := make([]byte, constant.BufferSize)
	n, err := stream.Read(buf)
	if n == 0 || err != nil {
		return false, req
	}
	if req.UnmarshalBinary(buf[:n]) != nil {
		log.Printf("[server] failed to decode request addr %v", err)
		return false, req
	}
	reqTime, _ := strconv.ParseInt(req.Timestamp, 10, 64)
	if time.Now().Unix()-reqTime > int64(constant.Timeout) {
		log.Printf("[server] timestamp expired %v", reqTime)
		return false, req
	}

	return true, req
}

func handleConn(session quic.Session, config config.Config) {
	// handshake
	ok, req := handshake(config, session)
	if !ok {
		return
	}
	// connect real server
	// log.Printf("[server] dial the real server %v %v:%v", req.Network, req.Host, req.Port)
	conn, err := net.DialTimeout(req.Network, net.JoinHostPort(req.Host, req.Port), time.Duration(constant.Timeout)*time.Second)
	if err != nil {
		log.Printf("[server] failed to dial the real server %v", err)
		return
	}
	stream, err := session.AcceptStream(context.Background())
	if err != nil {
		log.Println(err)
		return
	}
	go proxy.Copy(conn, stream)
	proxy.Copy(stream, conn)
}
