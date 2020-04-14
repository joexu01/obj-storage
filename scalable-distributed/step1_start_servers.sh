#!/bin/bash

# 用这个shell脚本为本机绑定多个IP地址
# 然后逐一启动服务

# 运行Docker的RabbitMQ镜像
sudo docker run --hostname localhost --name my-rabbitmq -p 15672:15672 -p 5672:5672 rabbitmq:3-management

# 绑定RabbitMQ服务器地址
sudo ifconfig lo:9 10.29.102.173/16

# 前六个是dataServers的IP，后两个是apiServers的IP
sudo ifconfig lo:1 10.29.1.1/16
sudo ifconfig lo:2 10.29.1.2/16
sudo ifconfig lo:3 10.29.1.3/16
sudo ifconfig lo:4 10.29.1.4/16
sudo ifconfig lo:5 10.29.1.5/16
sudo ifconfig lo:6 10.29.1.6/16
sudo ifconfig lo:7 10.29.2.1/16
sudo ifconfig lo:8 10.29.2.2/16

# shellcheck disable=SC2006
for i in `seq 1 6`; do mkdir -p ~/tmp/"$i"/objects; done

export RABBITMQ_SERVER=amqp://test:test@10.29.102.173:5672

# shellcheck disable=SC2034
LISTEN_ADDRESS=10.29.1.1:12345 STORAGE_ROOT=~/tmp/1 go run "$GOPATH"/src/github.com/joexu01/obj-storage/scalable-distributed/data_server/data_server.go &
LISTEN_ADDRESS=10.29.1.2:12345 STORAGE_ROOT=~/tmp/2 go run "$GOPATH"/src/github.com/joexu01/obj-storage/scalable-distributed/data_server/data_server.go &
LISTEN_ADDRESS=10.29.1.3:12345 STORAGE_ROOT=~/tmp/3 go run "$GOPATH"/src/github.com/joexu01/obj-storage/scalable-distributed/data_server/data_server.go &
LISTEN_ADDRESS=10.29.1.4:12345 STORAGE_ROOT=~/tmp/4 go run "$GOPATH"/src/github.com/joexu01/obj-storage/scalable-distributed/data_server/data_server.go &
LISTEN_ADDRESS=10.29.1.5:12345 STORAGE_ROOT=~/tmp/5 go run "$GOPATH"/src/github.com/joexu01/obj-storage/scalable-distributed/data_server/data_server.go &
LISTEN_ADDRESS=10.29.1.6:12345 STORAGE_ROOT=~/tmp/6 go run "$GOPATH"/src/github.com/joexu01/obj-storage/scalable-distributed/data_server/data_server.go &

# shellcheck disable=SC2034
LISTEN_ADDRESS=10.29.2.1:12345 go run "$GOPATH"/src/github.com/joexu01/obj-storage/scalable-distributed/api_server/api_server.go &
LISTEN_ADDRESS=10.29.2.2:12345 go run "$GOPATH"/src/github.com/joexu01/obj-storage/scalable-distributed/api_server/api_server.go &
