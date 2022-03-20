#!bin/bash
#Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/qsocks-linux-amd64 ./main.go
#Linux arm
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./bin/qsocks-linux-arm64 ./main.go
#Mac OS
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/qsocks-darwin-amd64 ./main.go
#Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/qsocks-windows-amd64.exe ./main.go
#Operwrt
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags="-s -w" -o ./bin/qsocks-openwrt-amd64 ./main.go

echo "DONE!!!"
