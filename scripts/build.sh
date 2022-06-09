#!bin/bash
#Linux amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/qsocks-linux-amd64 ./main.go
#Linux arm64
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./bin/qsocks-linux-arm64 ./main.go
#MacOS amd64
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/qsocks-darwin-amd64 ./main.go
#MacOS arm64
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./bin/qsocks-darwin-arm64 ./main.go
#Windows amd64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/qsocks-windows-amd64.exe ./main.go
#Windows arm64
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o ./bin/qsocks-windows-arm64.exe ./main.go
#OperWRT amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags="-s -w" -o ./bin/qsocks-openwrt-amd64 ./main.go

echo "DONE!!!"
