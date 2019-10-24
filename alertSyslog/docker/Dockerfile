#VERSION 1.0

#base image
FROM centos:latest

MAINTAINER xshrim “xshrim@yeah.net”

COPY syslogapi /root/

WORKDIR /root

# RUN chmod a+x /root/fabrest
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime  # 解决时区问题

EXPOSE 10901

#start fabric rest service
#ENTRYPOINT /root/fabrest -file /root/config/config.json
CMD ["/root/syslogapi"]
