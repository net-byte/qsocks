# qsocks

A socks5 proxy over quic.

[![Travis](https://travis-ci.com/net-byte/qsocks.svg?branch=main)](https://github.com/net-byte/qsocks)
[![Go Report Card](https://goreportcard.com/badge/github.com/net-byte/qsocks)](https://goreportcard.com/report/github.com/net-byte/qsocks)
![image](https://img.shields.io/badge/License-MIT-orange)
![image](https://img.shields.io/badge/License-Anti--996-red)

# Usage
```
Usage of /main:
  -S    server mode
  -bypass
        bypass private ip
  -l string
        local address (default "127.0.0.1:1080")
  -s string
        server address (default ":8443")
  -ck string
        client key file path (default "../certs/client.key")
  -cp string
        client pem file path (default "../certs/client.pem")
  -sk string
        server key file path (default "../certs/server.key")
  -sp string
        server pem file path (default "../certs/server.pem")
```

# Docker

## Run client
```
docker run -d --restart=always --name qsocks-client -p 1083:1083 -p 1083:1083/udp netbyte/qsocks -l :1083 -s SERVER_IP:8443 -ck=/app/certs/client.key -cp=/app/certs/client.pem -sk=/app/certs/server.key -sp=/app/certs/server.pem

```

## Run server
```
docker run -d --restart=always --name qsocks-server -p 8443:8443/udp netbyte/qsocks -S -s :8443 -ck=/app/certs/client.key -cp=/app/certs/client.pem -sk=/app/certs/server.key -sp=/app/certs/server.pem
```


# License
[The MIT License (MIT)](https://raw.githubusercontent.com/net-byte/qsocks/main/LICENSE)


