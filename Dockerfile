FROM golang:alpine

WORKDIR /app
COPY . /app
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn
RUN go build -o ./bin/qsocks ./main.go

ENTRYPOINT ["./bin/qsocks"]

