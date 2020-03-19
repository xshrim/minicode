#VERSION 1.0

#base image
FROM alpine:latest

MAINTAINER xshrim “xshrim@yeah.net”

# 替换阿里云的源
RUN echo "http://mirrors.aliyun.com/alpine/latest-stable/main/" > /etc/apk/repositories
RUN echo "http://mirrors.aliyun.com/alpine/latest-stable/community/" >> /etc/apk/repositories

# 同步时间

# 更新源、安装openssh 并修改配置文件和生成key 并且同步时间
RUN apk update && \
  mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && \
  apk add --no-cache openssh-server tzdata && \
  cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
  sed -i "s/#PermitRootLogin.*/PermitRootLogin yes/g" /etc/ssh/sshd_config && \
  ssh-keygen -t rsa -P "" -f /etc/ssh/ssh_host_rsa_key && \
  ssh-keygen -t ecdsa -P "" -f /etc/ssh/ssh_host_ecdsa_key && \
  ssh-keygen -t ed25519 -P "" -f /etc/ssh/ssh_host_ed25519_key && \
  echo "root:admin" | chpasswd

ADD webssh /root/
ADD run.sh /root/
COPY static /root/static

WORKDIR /root

# 开放ssh端口
EXPOSE 22
# 开放webssh端口
EXPOSE 2222

CMD ["sh","/root/run.sh"]