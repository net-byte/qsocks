package proxy

import (
	"bytes"
	"context"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/net-byte/qsocks/common/constant"
	"github.com/net-byte/qsocks/config"
)

func UDPProxy(tcpConn net.Conn, udpConn *net.UDPConn, config config.Config) {
	defer tcpConn.Close()
	if udpConn == nil {
		log.Printf("[udp] failed to start udp server on %v", config.LocalAddr)
		return
	}
	bindAddr, _ := net.ResolveUDPAddr("udp", udpConn.LocalAddr().String())
	//response to client
	ResponseUDP(tcpConn, bindAddr)
	//keep tcp alive
	done := make(chan bool)
	go keepTCPAlive(tcpConn.(*net.TCPConn), done)
	<-done
}

func keepTCPAlive(tcpConn *net.TCPConn, done chan<- bool) {
	tcpConn.SetKeepAlive(true)
	buf := make([]byte, constant.BufferSize)
	for {
		_, err := tcpConn.Read(buf[0:])
		if err != nil {
			break
		}
	}
	done <- true
}

type UDPRelay struct {
	UDPConn    *net.UDPConn
	Config     config.Config
	headerMap  sync.Map
	sessionMap sync.Map
}

func (relay *UDPRelay) Start() *net.UDPConn {
	udpAddr, _ := net.ResolveUDPAddr("udp", relay.Config.LocalAddr)
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Printf("[udp] failed to listen udp %v", err)
		return nil
	}
	relay.UDPConn = udpConn
	go relay.toRemote()
	log.Printf("qsocks [udp] client started on %v", relay.Config.LocalAddr)
	return relay.UDPConn
}

func (relay *UDPRelay) toRemote() {
	defer relay.UDPConn.Close()
	buf := make([]byte, constant.BufferSize)
	for {
		relay.UDPConn.SetReadDeadline(time.Now().Add(time.Duration(constant.Timeout) * time.Second))
		n, cliAddr, err := relay.UDPConn.ReadFromUDP(buf)
		if err != nil || err == io.EOF || n == 0 {
			continue
		}
		b := buf[:n]
		dstAddr, header, data := relay.getAddr(b)
		if dstAddr == nil || header == nil || data == nil {
			continue
		}
		key := cliAddr.String()
		var session quic.Session
		if value, ok := relay.sessionMap.Load(key); ok {
			session = value.(quic.Session)
			stream, err := session.OpenStreamSync(context.Background())
			if err != nil {
				log.Println(err)
				continue
			}
			stream.Write(data)
		} else {
			session = ConnectServer(relay.Config)
			if session == nil {
				continue
			}
			ok := Handshake("udp", dstAddr.IP.String(), strconv.Itoa(dstAddr.Port), session)
			if !ok {
				continue
			}
			stream, err := session.OpenStreamSync(context.Background())
			if err != nil {
				log.Println(err)
				continue
			}
			go relay.toLocal(session, stream, cliAddr)
			stream.Write(data)
			relay.sessionMap.Store(key, session)
			relay.headerMap.Store(key, header)
		}
	}
}

func (relay *UDPRelay) toLocal(session quic.Session, stream quic.Stream, cliAddr *net.UDPAddr) {
	defer stream.Close()
	defer session.CloseWithError(0, "bye")
	key := cliAddr.String()
	buf := make([]byte, constant.BufferSize)
	for {
		n, err := stream.Read(buf)
		if n == 0 || err != nil {
			break
		}
		if header, ok := relay.headerMap.Load(key); ok {
			var data bytes.Buffer
			data.Write(header.([]byte))
			data.Write(buf[:n])
			relay.UDPConn.WriteToUDP(data.Bytes(), cliAddr)
		}
	}
	relay.headerMap.Delete(key)
	relay.sessionMap.Delete(key)
}

func (relay *UDPRelay) getAddr(b []byte) (dstAddr *net.UDPAddr, header []byte, data []byte) {
	/*
	   +----+------+------+----------+----------+----------+
	   |RSV | FRAG | ATYP | DST.ADDR | DST.PORT |   DATA   |
	   +----+------+------+----------+----------+----------+
	   |  2 |   1  |   1  | Variable |     2    | Variable |
	   +----+------+------+----------+----------+----------+
	*/
	if b[2] != 0x00 {
		log.Printf("[udp] not support frag %v", b[2])
		return nil, nil, nil
	}
	switch b[3] {
	case constant.Ipv4Address:
		dstAddr = &net.UDPAddr{
			IP:   net.IPv4(b[4], b[5], b[6], b[7]),
			Port: int(b[8])<<8 | int(b[9]),
		}
		header = b[0:10]
		data = b[10:]
	case constant.FqdnAddress:
		domainLength := int(b[4])
		domain := string(b[5 : 5+domainLength])
		ipAddr, err := net.ResolveIPAddr("ip", domain)
		if err != nil {
			log.Printf("[udp] failed to resolve dns %s:%v", domain, err)
			return nil, nil, nil
		}
		dstAddr = &net.UDPAddr{
			IP:   ipAddr.IP,
			Port: int(b[5+domainLength])<<8 | int(b[6+domainLength]),
		}
		header = b[0 : 7+domainLength]
		data = b[7+domainLength:]
	case constant.Ipv6Address:
		{
			dstAddr = &net.UDPAddr{
				IP:   net.IP(b[4:19]),
				Port: int(b[20])<<8 | int(b[21]),
			}
			header = b[0:22]
			data = b[22:]
		}
	default:
		return nil, nil, nil
	}
	return dstAddr, header, data
}
