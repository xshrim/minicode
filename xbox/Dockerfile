FROM alpine:latest

LABEL maintainer="xshrim@yeah.net"

WORKDIR /root/
# py2-virtualenv \
# apparmor issue #14140 

# Settings
COPY motd /etc/motd
ADD artifacts /root/

# mirrors.ustc.edu.cn
# mirrors.tuna.tsinghua.edu.cn
RUN set -ex \
    && echo "http://nl.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories \
    && echo "http://nl.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories \
    && echo "http://nl.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && sed -i 's/nl.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
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
    cmus \
    ctop \
    curl \
    docker \
    docker-compose \
    dhcping \
    dfc \
    dos2unix \
    drill \
    etcd-ctl \
    ethtool \
    exa \
    expect \
    fabric \
    file\
    figlet \
    fio \
    fping \
    git \
    gotop \
    gron \
    hexdump \
    highlight \
    httpie \
    htop \
    hyperfine \
    libc6-compat \
    liboping \
    lua5.3 \
    iftop \
    ipcalc \
    iperf3 \
    iproute2 \
    ipset \
    iptables \
    iptraf-ng \
    iptstate \
    iputils \
    ipvsadm \
    ioping \
    iozone \
    jq \
    kubectl \
    minio-client \
    mtr \
    mysql-client \
    mycli \
    net-snmp-tools \
    netcat-openbsd \
    nftables \
    nfs-utils \
    neofetch \
    nethogs \
    nginx \
    ngrep \
    nmap \
    nmap-nping \
    nuttcp \
    openssl \
    openssh \
    openssh-server \
    openssh-client \
    postgresql-client \
    pgcli \
    qperf \
    tree \
    redis \
    ranger \
    rsync \
    samba \
    samba-client \
    samba-server \
    samba-common-tools \
    socat \
    sqlite \
    strace \
    sshpass \
    speedtest-cli \
    websocat \
    websocketd \
    tzdata \
    tcpdump \
    tcptraceroute \
    tshark \
    tmux \
    util-linux \
    vim \
    zsh \
    && mkdir /share \
    && mv /usr/sbin/tcpdump /usr/bin/tcpdump \
    && echo 'alias iperf="iperf3"' >> /etc/profile \
    && echo 'export PATH=/root/tools/:$PATH' >> /etc/profile \
    && export TERM='xterm-256color'; ranger --copy-config=all \
    && chgrp root /usr/bin/dumpcap \
    && chmod 750 /usr/bin/dumpcap \
    && setcap cap_net_raw,cap_net_admin+eip /usr/bin/dumpcap \
    && sed -i 's^/bin/ash^/bin/zsh^g' /etc/passwd \
    && sed -i "s/#PermitRootLogin.*/PermitRootLogin yes/g" /etc/ssh/sshd_config \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && ssh-keygen -t rsa -P "" -f /etc/ssh/ssh_host_rsa_key \
    && ssh-keygen -t ecdsa -P "" -f /etc/ssh/ssh_host_ecdsa_key \
    && ssh-keygen -t ed25519 -P "" -f /etc/ssh/ssh_host_ed25519_key \
    && echo "root:root" | chpasswd \
    && addgroup -g 1111 smb \
    && adduser -D -H -G smb -s /bin/false -u 1111 smb \
    && echo -e "smb\nsmb" | smbpasswd -a -s -c /root/smb.conf smb

#RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2   # 解决go语言程序无法在alpine执行的问题
#RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && apk add -U util-linux && apk add -U tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime  # 解决go语言程序无法在alpine执行的问题和syslog不支持udp的问题和时区问题

# 开放ssh webssh gofs gohttpserver smb端口
EXPOSE 22 2222 2333 2444 135/tcp 137/udp 138/udp 139/tcp 445/tcp

# CMD ["/usr/sbin/sshd", "-D"]
CMD ["sh","/root/run.sh"]
# CMD ["/bin/zsh","-l"]
