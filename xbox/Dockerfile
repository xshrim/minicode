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
    dfc \
    dos2unix \
    drill \
    etcd \
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
    glances \
    highlight \
    httpie \
    htop \
    libc6-compat \
    liboping \
    lua5.1 \
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
    mtr \
    mysql-client \
    net-snmp-tools \
    netcat-openbsd \
    nftables \
    nfs-utils \
    nethogs \
    nginx \
    ngrep \
    nmap \
    nmap-nping \
    ntop \
    nuttcp \
    openssl \
    openssh \
    openssh-server \
    openssh-client \
    py-crypto \
    postgresql-client \
    qperf \
    tree \
    redis \
    rsync \
    samba \
    samba-client \
    samba-server \
    samba-common-tools \
    scapy \
    socat \
    sqlite \
    strace \
    sshpass \
    speedtest \
    speedtest-cli \
    tzdata \
    tcpdump \
    tcptraceroute \
    util-linux \
    vim \
    zsh

WORKDIR /root/
#py2-virtualenv \
# apparmor issue #14140 

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
ADD webssh /root/
ADD run.sh /root/
ADD deploy.yaml /root/
COPY static /root/static
COPY oh-my-zsh /root/.oh-my-zsh
COPY zshrc /root/.zshrc
COPY vimrc /root/.vimrc
COPY smb.conf /root/smb.conf

RUN mkdir /share && mv /usr/sbin/tcpdump /usr/bin/tcpdump && echo 'alias iperf="iperf3"' >> /etc/profile && sed -i 's^/bin/ash^/bin/zsh^g' /etc/passwd && sed -i "s/#PermitRootLogin.*/PermitRootLogin yes/g" /etc/ssh/sshd_config && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && ssh-keygen -t rsa -P "" -f /etc/ssh/ssh_host_rsa_key && ssh-keygen -t ecdsa -P "" -f /etc/ssh/ssh_host_ecdsa_key && ssh-keygen -t ed25519 -P "" -f /etc/ssh/ssh_host_ed25519_key && echo "root:root" | chpasswd && addgroup -g 1111 smb && adduser -D -H -G smb -s /bin/false -u 1111 smb && echo -e "smb\nsmb" | smbpasswd -a -s -c /root/smb.conf smb

#RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2   # 解决go语言程序无法在alpine执行的问题
#RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && apk add -U util-linux && apk add -U tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime  # 解决go语言程序无法在alpine执行的问题和syslog不支持udp的问题和时区问题

# 开放ssh端口
EXPOSE 22
# 开放smb端口
EXPOSE 135/tcp 137/udp 138/udp 139/tcp 445/tcp
# 开放webssh端口
EXPOSE 2222
# 开放gofs端口
EXPOSE 2333

# CMD ["/usr/sbin/sshd", "-D"]
CMD ["sh","/root/run.sh"]
# CMD ["/bin/zsh","-l"]
