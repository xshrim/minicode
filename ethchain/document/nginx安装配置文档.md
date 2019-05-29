# Nginx安装配置

## 前置准备

### 系统环境

- 操作系统: CentOS7.4 (x86_64)
- Nginx: nginx-1.14.2.tar.gz

### 存储规划

| 文件系统 | 空间(GB) | 所属用户 | 所属组 | 目录权限 | 存储要求                   | 描述信息      |
| :------- | :------- | :------- | :----- | :------- | :------------------------- | :------------ |
| /nginx   |          | root     | root   |          | 本地/数据库HA/网络文件系统 | Nginx安装目录 |

### 依赖安装

```bash
yum install -y gcc-c++ pcre pcre-devel zlib zlib-devel openssl openssl-devel
# gcc-c++用于编译nginx源码
# pcre用于提供nginx的http模块的正则表达式支持
# zlib用于提供nginx对http数据包的压缩支持
# openssl用于提供nginx对https的加密传输支持
```

## 安装配置

### 软件安装

```bash
wget https://nginx.org/download/nginx-1.14.2.tar.gz
tar xvf nginx-1.14.2.tar.gz                                  # 解压nginx源码
cd nginx-1.14.2                                              # 编译安装nginx           
./configure \
--prefix=/nginx \
--conf-path=/nginx/conf/nginx.conf \
--pid-path=/nginx/conf/nginx.pid 
make && make install
echo "export PATH=/nginx/sbin:$PATH" >> /etc/profile          # 配置nginx环境变量
source /etc/profile
echo "/nginx/sbin/nginx" >> /etc/rc.local                    # 配置nginx开机自动启动
nginx                                                        # 启动nginx服务
firewall-cmd --zone=public --add-port=8080/tcp --permanent   # 永久开启nginx对外服务端口
firewall-cmd reload
firewall-cmd --query-port=8080/tcp --permanent
curl $(hostname -I | cut -d ' ' -f 1):80                     # 验证服务是否正常
nginx -V
```

### 配置与优化

```sh
########### nginx.conf文件#################
## 每个指令必须有分号结束。
#user administrator administrators;  #配置用户或者组，默认为nobody nobody。
worker_processes 2;  #允许生成的进程数，默认为1
#pid /nginx/conf/nginx.pid;   #指定nginx进程运行文件存放地址
error_log log/error.log debug;  #制定日志路径，级别。这个设置可以放入全局块，http块，server块，级别以此为：debug|info|notice|warn|error|crit|alert|emerg
events {
    accept_mutex on;   #设置网路连接序列化，防止惊群现象发生，默认为on
    multi_accept on;  #设置一个进程是否同时接受多个网络连接，默认为off
    use epoll;      #事件驱动模型，select|poll|kqueue|epoll|resig|/dev/poll|eventport, 默认根据操作系统自动选择
    worker_connections  1024;    #最大连接数，默认为512
}
http {
    include       mime.types;   #文件扩展名与文件类型映射表
    default_type  application/octet-stream; #默认文件类型，默认为text/plain
    charset utf-8;               #页面数据字符集
    #access_log off; #取消服务日志 
    log_format myFormat '$remote_addr–$remote_user [$time_local] $request $status $body_bytes_sent $http_referer $http_user_agent $http_x_forwarded_for'; #自定义格式
    access_log log/access.log myFormat;  #combined为日志格式的默认值
    client_header_buffer_size 4k;      #客户端请求头部缓冲区最大值
    gzip on;                           #开启数据包压缩, 默认关闭
    gzip_min_length 2k;                #压缩数据大小下限
    gzip_comp_level 2;                 #数据压缩级别
    gzip_types text/plain application/x-javascript text/css application/xml;    #压缩数据类型
    sendfile on;   #允许sendfile方式传输文件，默认为off，可以在http块，server块，location块。
    sendfile_max_chunk 100k;  #每个进程每次调用传输数量不能大于设定的值，默认为0，即不设上限。
    keepalive_timeout 65;  #连接超时时间，默认为75s，可以在http，server，location块。

    upstream mysvr {   
      server 127.0.0.1:7878;
      server 192.168.10.121:3333 backup;  #热备
    }
    error_page 404 https://www.baidu.com; #错误页
    server {
        keepalive_requests 120; #单连接请求上限次数。
        listen       4545;   #监听端口
        server_name  127.0.0.1;   #监听地址       
        location  ~*^.+$ {       #请求的url过滤，正则匹配，~为区分大小写，~*为不区分大小写。
           #root path;  #根目录
           #index vv.txt;  #设置默认页
           proxy_pass  http://mysvr;  #请求转向mysvr 定义的服务器列表
           deny 127.0.0.1;  #拒绝的ip
           allow 172.18.5.54; #允许的ip           
        } 
    }
}
```

## 内核调优

### sysctl

[参数说明]

- net.ipv4.tcp_max_tw_buckets = 6000

> timewait 的数量，默认是180000。

- net.ipv4.ip_local_port_range = 1024 65000

> 允许系统打开的端口范围。

- net.ipv4.tcp_tw_recycle = 1

> 启用timewait 快速回收。

- net.ipv4.tcp_tw_reuse = 1

> 开启重用。允许将TIME-WAIT sockets 重新用于新的TCP 连接。

- net.ipv4.tcp_syncookies = 1

> 开启SYN Cookies，当出现SYN 等待队列溢出时，启用cookies 来处理。

- net.core.somaxconn = 262144

> web 应用中listen 函数的backlog 默认会给我们内核参数的net.core.somaxconn 限制到128，而nginx 定义的NGX_LISTEN_BACKLOG 默认为511，所以有必要调整这个值。

- net.core.netdev_max_backlog = 262144

> 每个网络接口接收数据包的速率比内核处理这些包的速率快时，允许送到队列的数据包的最大数目。

- net.ipv4.tcp_max_orphans = 262144

> 系统中最多有多少个TCP 套接字不被关联到任何一个用户文件句柄上。如果超过这个数字，孤儿连接将即刻被复位并打印出警告信息。这个限制仅仅是为了防止简单的DoS 攻击，不能过分依靠它或者人为地减小这个值，更应该增加这个值(如果增加了内存之后)。

- net.ipv4.tcp_max_syn_backlog = 262144

> 记录的那些尚未收到客户端确认信息的连接请求的最大值。对于有128M 内存的系统而言，缺省值是1024，小内存的系统则是128。

- net.ipv4.tcp_timestamps = 0

> 时间戳可以避免序列号的卷绕。一个1Gbps 的链路肯定会遇到以前用过的序列号。时间戳能够让内核接受这种“异常”的数据包。这里需要将其关掉。

- net.ipv4.tcp_synack_retries = 1

> 为了打开对端的连接，内核需要发送一个SYN 并附带一个回应前面一个SYN 的ACK。也就是所谓三次握手中的第二次握手。这个设置决定了内核放弃连接之前发送SYN+ACK 包的数量。

- net.ipv4.tcp_syn_retries = 1

> 在内核放弃建立连接之前发送SYN 包的数量。

- net.ipv4.tcp_fin_timeout = 1

> 如果套接字由本端要求关闭，这个参数决定了它保持在FIN-WAIT-2 状态的时间。对端可以出错并永远不关闭连接，甚至意外当机。缺省值是60 秒。2.2 内核的通常值是180 秒，3你可以按这个设置，但要记住的是，即使你的机器是一个轻载的WEB 服务器，也有因为大量的死套接字而内存溢出的风险，FIN- WAIT-2 的危险性比FIN-WAIT-1 要小，因为它最多只能吃掉1.5K 内存，但是它们的生存期长些。

- net.ipv4.tcp_keepalive_time = 30

> 当keepalive 起用的时候，TCP 发送keepalive 消息的频度。缺省是2 小时。

[范例]

```sh
# cat /etc/sysctl.conf
# sysctl -p
net.ipv4.ip_forward = 0
net.ipv4.conf.default.rp_filter = 1
net.ipv4.conf.default.accept_source_route = 0
kernel.sysrq = 0
kernel.core_uses_pid = 1
net.ipv4.tcp_syncookies = 1
kernel.msgmnb = 65536
kernel.msgmax = 65536
kernel.shmmax = 68719476736
kernel.shmall = 4294967296
net.ipv4.tcp_max_tw_buckets = 6000
net.ipv4.tcp_sack = 1
net.ipv4.tcp_window_scaling = 1
net.ipv4.tcp_rmem = 4096 87380 4194304
net.ipv4.tcp_wmem = 4096 16384 4194304
net.core.wmem_default = 8388608
net.core.rmem_default = 8388608
net.core.rmem_max = 16777216
net.core.wmem_max = 16777216
net.core.netdev_max_backlog = 262144
net.core.somaxconn = 262144
net.ipv4.tcp_max_orphans = 3276800
net.ipv4.tcp_max_syn_backlog = 262144
net.ipv4.tcp_timestamps = 0
net.ipv4.tcp_synack_retries = 1
net.ipv4.tcp_syn_retries = 1
net.ipv4.tcp_tw_recycle = 1
net.ipv4.tcp_tw_reuse = 1
net.ipv4.tcp_mem = 94500000 915000000 927000000
net.ipv4.tcp_fin_timeout = 1
net.ipv4.tcp_keepalive_time = 30
net.ipv4.ip_local_port_range = 1024 65000
```

### limits

```sh
# /etc/security/limits.conf
# ulimit -a
* soft nofile 65535
* hard nofile 65535
* soft nproc 65535
* hard nproc 65535
```