# Dockerfile References: https://docs.docker.com/engine/reference/builder/

######## 构建阶段 生成可执行文件 #######
FROM golang:latest

LABEL maintainer="QianLu990613@foxmail.com"

RUN mkdir /IntelligentTransfer


WORKDIR /IntelligentTransfer

COPY . /IntelligentTransfer

EXPOSE 40000

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on GOARCH=amd64 go build -o main.go
RUN cd /IntelligentTransfer
# 启动服务
CMD ["/IntelligentTransfer/main"]
