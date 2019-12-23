FROM alpine:latest

MAINTAINER xshrim "xshrim@yeah.net"

RUN set -ex \
    && echo "http://nl.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories \
    && echo "http://nl.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories \
    && echo "http://nl.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && sed -i 's/nl.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update \
    && apk upgrade \
    && apk add --no-cache \
    apache2-utils \
    bash \
    bat \
    bind-tools \
    bird \
    bridge-utils \
    busybox-extras \
    conntrack-tools \
    curl \
    docker \
    dhcping \
    drill \
    etcd \
    etcd-ctl \
    ethtool \
    exa \
    expect \
    file\
    figlet \
    fio \
    fping \
    git \
    httpie \
    htop \
    iftop \
    iperf \
    iproute2 \
    ipset \
    iptables \
    iptraf-ng \
    iputils \
    ipvsadm \
    libc6-compat \
    liboping \
    mtr \
    mysql-client \
    net-snmp-tools \
    netcat-openbsd \
    nftables \
    nfs-utils \
    nginx \
    ngrep \
    nmap \
    nmap-nping \
    openssl \
    openssh \
    openssh-client \
    redis \
    sshpass \
    py-crypto \
    python3 \
    postgresql-client \
    tree \
    samba \
    scapy \
    socat \
    sqlite \
    strace \
    tcpdump \
    tcptraceroute \
    util-linux \
    vim \
    zsh

#py2-virtualenv \
# apparmor issue #14140 
RUN mv /usr/sbin/tcpdump /usr/bin/tcpdump && echo 'alias python="python3"' >> /etc/profile && sed -i 's^/bin/ash^/bin/zsh^g' /etc/passwd

# Installing ctop - top-like container monitor
#RUN wget https://github.com/bcicen/ctop/releases/download/v0.7.1/ctop-0.7.1-linux-amd64 -O /usr/local/bin/ctop && chmod +x /usr/local/bin/ctop

# Installing calicoctl
#ARG CALICOCTL_VERSION=v3.10.2
#RUN wget https://github.com/projectcalico/calicoctl/releases/download/${CALICOCTL_VERSION}/calicoctl && chmod +x calicoctl && mv calicoctl /usr/bin/

# Installing kubectl
#curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl && chmod +x kubectl && mv kubectl /usr/bin/

# Settings
#ADD motd /etc/motd
#ADD profile  /etc/profile
ADD tools/* /usr/bin/
COPY oh-my-zsh /root/.oh-my-zsh
COPY zshrc /root/.zshrc
COPY vimrc /root/.vimrc

#RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2   # 解决go语言程序无法在alpine执行的问题
#RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && apk add -U util-linux && apk add -U tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime  # 解决go语言程序无法在alpine执行的问题和syslog不支持udp的问题和时区问题

CMD ["/bin/zsh","-l"]
