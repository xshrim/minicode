FROM centos:7

LABEL xshrim "xshrim@yeah.net"

WORKDIR /root/

ADD alertsyslog /root/

RUN mkdir /data && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime  # 解决时区问题

EXPOSE 10901

CMD ["/root/alertsyslog"]
