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
		stream, err := session.AcceptStream(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConn(stream, config)
	}

}

func handshake(config config.Config, stream quic.Stream) (bool, proxy.RequestAddr) {
	var req proxy.RequestAddr
	buf := make([]byte, constant.BufferSize)
	n, err := stream.Read(buf)
	if n == 0 || err != nil {
		return false, req
	}
	if req.UnmarshalBinary(buf[0:n]) != nil {
		log.Printf("[server] failed to unmarshal binary %v", err)
		return false, req
	}
	reqTime, _ := strconv.ParseInt(req.Timestamp, 10, 64)
	if time.Now().Unix()-reqTime > int64(constant.Timeout) {
		log.Printf("[server] timestamp expired %v", reqTime)
		return false, req
	}

	return true, req
}

func handleConn(stream quic.Stream, config config.Config) {
	// handshake
	ok, req := handshake(config, stream)
	if !ok {
		stream.Close()
		return
	}
	// connect real server
	conn, err := net.DialTimeout(req.Network, net.JoinHostPort(req.Host, req.Port), time.Duration(constant.Timeout)*time.Second)
	if err != nil {
		stream.Close()
		log.Printf("[server] failed to dial the real server%v", err)
		return
	}

	go proxy.Copy(conn, stream)
	go proxy.Copy(stream, conn)
}
