# Dockerfile References: https://docs.docker.com/engine/reference/builder/

######## 构建阶段 生成可执行文件 #######
FROM golang:1.15

LABEL maintainer="QianLu990613@foxmail.com"

WORKDIR /IntelligentTransfer

COPY . .

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o main.go

# 启动服务
CMD ["./main"]
