package config

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
)

type Config struct {
	LocalAddr  string
	ServerAddr string
	ClientKey  string
	ClientPem  string
	ServerKey  string
	ServerPem  string
	ServerMode bool
	Bypass     bool
}

func (config *Config) GetClientTLSConfig() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(config.ClientPem, config.ClientKey)
	if err != nil {
		return nil, err
	}

	certBytes, err := ioutil.ReadFile(config.ClientPem)
	if err != nil {
		return nil, err
	}

	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		panic("failed to parse client certificate")
	}

	return &tls.Config{
		RootCAs:            clientCertPool,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
		NextProtos:         []string{"qtun/1.0"},
	}, nil
}

func (config *Config) GetServerTLSConfig() (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(config.ServerPem, config.ServerKey)
	if err != nil {
		return nil, err
	}

	certBytes, err := ioutil.ReadFile(config.ClientPem)
	if err != nil {
		return nil, err
	}

	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		panic("failed to parse client certificate")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCertPool,
		NextProtos:   []string{"qtun/1.0"},
	}, nil
}
