#VERSION 1.0

#base image
FROM alpine:latest

MAINTAINER xshrim “xshrim@yeah.net”

COPY work/* /root/

WORKDIR /root

# RUN chmod a+x /root/fabrest
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2   # 解决go语言程序无法在alpine执行的问题
#RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && apk add -U util-linux && apk add -U tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime  # 解决go语言程序无法在alpine执行的问题和syslog不支持udp的问题和时区问题

EXPOSE 9999

#start fabric rest service
#ENTRYPOINT /root/fabrest -file /root/config/config.json
CMD ["/root/fabrest", "-file", "root/config/config.json"]
