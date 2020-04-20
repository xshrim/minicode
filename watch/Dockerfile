# VERSION 1.0

# build
# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags netgo -o filewatch main.go
# run
# docker run --rm -it -e WATCHFILE=/root/demo.json -e WATCHOPT=create_modify -e WATCHCMD=ls -p 2233:2233 -v /root/demo.json:/root/demo.json filewatch:1.0

# base image

FROM alpine:latest

LABEL maintainer="xshrim xshrim@yeah.net"

ADD filewatch /root/

WORKDIR /root

# RUN chmod a+x /root/fabrest
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
# RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && apk add -U util-linux && apk add -U tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime  # 解决go语言程序无法在alpine执行的问题和syslog不支持udp的问题和时区问题

EXPOSE 2233

# start fabric rest service
# ENTRYPOINT /root/fabrest -file /root/config/config.json
CMD ["/root/filewatch"]
