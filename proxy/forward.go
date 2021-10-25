package proxy

import (
	"context"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/net-byte/qsocks/common/cipher"
	"github.com/net-byte/qsocks/config"
)

func ConnectServer(network string, host string, port string, config config.Config) quic.Stream {
	// handshake
	req := &RequestAddr{}
	req.Network = network
	req.Host = host
	req.Port = port
	req.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	req.Random = cipher.Random()
	data, err := req.MarshalBinary()
	if err != nil {
		log.Printf("[client] failed to marshal binary %v", err)
		return nil
	}
	tlsConf, err := config.GetClientTLSConfig()
	if err != nil {
		log.Println(err)
		return nil
	}
	session, err := quic.DialAddr(config.ServerAddr, tlsConf, nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		log.Println(err)
		return nil
	}
	stream.Write(data)
	return stream
}

func Copy(destination io.WriteCloser, source io.ReadCloser) {
	if destination == nil || source == nil {
		return
	}
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}
