package proxy

import (
	"context"
	"io"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/net-byte/qsocks/common/cipher"
	"github.com/net-byte/qsocks/config"
	"github.com/net-byte/qsocks/proto"
)

var _lock sync.Mutex
var _session quic.Session

func ConnectServer(config config.Config) quic.Session {
	_lock.Lock()
	if _session != nil {
		_lock.Unlock()
		return _session
	}
	_tlsConf, err := config.GetClientTLSConfig()
	if err != nil {
		log.Println(err)
		_lock.Unlock()
		return nil
	}
	quicConfig := &quic.Config{
		ConnectionIDLength:   12,
		HandshakeIdleTimeout: time.Second * 10,
		MaxIdleTimeout:       time.Second * 30,
		KeepAlive:            false,
	}
	_session, err = quic.DialAddr(config.ServerAddr, _tlsConf, quicConfig)
	if err != nil {
		log.Println(err)
		_lock.Unlock()
		return nil
	}
	_lock.Unlock()
	return _session
}

func Handshake(network string, host string, port string, session quic.Session) (bool, quic.Stream) {
	// handshake
	req := &RequestAddr{}
	req.Network = network
	req.Host = host
	req.Port = port
	req.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	req.Random = cipher.Random()
	data, err := req.MarshalBinary()
	if err != nil {
		log.Printf("[client] failed to encode request addr %v", err)
		return false, nil
	}
	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		log.Println(err)
		_session = nil
		return false, nil
	}
	edata, err := proto.Encode(data)
	if err != nil {
		log.Println(err)
		return false, nil
	}
	_, err = stream.Write(edata)
	if err != nil {
		log.Println(err)
		_session = nil
		return false, nil
	}
	return true, stream
}

func copy(destination io.WriteCloser, source io.ReadCloser) {
	if destination == nil || source == nil {
		return
	}
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}
