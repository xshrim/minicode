# Apache安装配置

## 前置准备

### 系统环境

- 操作系统: CentOS7.4 (x86_64)
- Apache: httpd-2.4.25.tar.gz

### 存储规划

| 文件系统 | 空间(GB) | 所属用户 | 所属组 | 目录权限 | 存储要求                   | 描述信息       |
|:---------|:---------|:---------|:-------|:---------|:---------------------------|:---------------|
| /apache  |          | root     | root   |          | 本地/数据库HA/网络文件系统 | apache安装目录 |

### 依赖安装

#### apr

```bash
yum install -y gcc libtools
wget http://mirrors.cnnic.cn/apache//apr/apr-1.5.2.tar.gz
tar zxf apache/apr-1.5.2.tar.gz  
cd apr-1.5.2  
./configure --prefix=/usr/local/apr  
make && make instal
```

#### apr-util

```bash
wget http://mirrors.cnnic.cn/apache//apr/apr-util-1.5.4.tar.gz
tar zxf apr-util-1.5.4.tar.gz  
cd  apr-util-1.5.4  
./configure --prefix=/usr/local/apr-util --with-apr=/usr/local/apr  
make && make install
```

#### pcre

```bash
yum install -y gcc-c++ 
wget http://exim.mirror.fr/pcre/pcre-8.38.tar.gz
tar zxf pcre-8.38.tar.gz  
cd pcre-8.38  
./configure --prefix=/usr/local/pcre
make && make instal
```

[注] 也可直接从yum仓库中安装上述依赖

```bash
yum -y -q install gcc-c++ libtools apr apr-devel apr-util apr-util-devel openssl openssl-devel pcre pcre-devel mod_ssl expat expat-devel
```

## 安装配置

### 软件安装

```bash
yum install -y openssl-devel mod_ssl expat-devel
wget http://apache.fayea.com//httpd/httpd-2.4.25.tar.gz
tar zxf  httpd-2.4.25.tar.gz  
cd httpd-2.4.25 
./configure --prefix=/apache --with-apr=/usr/local/apr --with-apr-util=/usr/local/apr-util --with-pcre=/usr/local/pcre  
make && make install
echo "export PATH=$PATH:/apache/bin" >> /etc/profile
source /etc/profile
cp /apache/bin/apachectl /etc/init.d/httpd                      # 配置系统服务
sed -i '2i# chkconfig: 113 64 38' /etc/init.d/httpd
chmod 755 /etc/init.d/httpd
# service tomcat start
systemctl start httpd                                            # 启动apache服务
# chkconfig --add httpd
# chkconfig httpd on
systemctl enable httpd                                            # 设置apache开机启
# chkconfig --list  
systemctl list-units --type service --all|grep httpd
firewall-cmd --zone=public --add-port=80/tcp --permanent         # 永久开启apache对外服务端口
firewall-cmd reload
firewall-cmd --query-port=80/tcp --permanent
curl $(hostname -I | cut -d ' ' -f 1):80                         # 验证服务是否正常
apachectl fullstatus
```

### 配置与优化


#### httpd相关命令

- 查看当前安装模块mpm(多路处理器)

```bash
httpd -l
```

- 查看httpd进程数(即各个mpm模式下Apache能够处理的并发请求数)

```bash
ps -ef | grep httpd | wc -l
```

得到的结果数字就是表示可以同时并发的进程数据，一个父进程，5个子进程。apache默认是开启5个子进程

- 查看Apache的并发请求数及其TCP连接状态

```bash
netstat -n | awk '/^tcp/ { S[$NF]} END {for(a in S) print a, S[a]}'
```

[状态说明]

> ESTABLISHED: 连接已建立
> CLOSED：无连接是活动的或正在进行
> LISTEN：服务器在等待进入呼叫
> SYN_RECV：一个连接请求已经到达，等待确认
> SYN_SENT：应用已经开始，打开一个连接
> ESTABLISHED：正常数据传输状态
> FIN_WAIT1：应用说它已经完成
> FIN_WAIT2：另一边已同意释放
> ITMED_WAIT：等待所有分组死掉
> CLOSING：两边同时尝试关闭
> TIME_WAIT：另一边已初始化一个释放
> LAST_ACK：等待所有分组死掉

- 查看请求80服务的client ip按照连接数排序

```bash
netstat -nat|grep ":80"|awk '{print $5}' |awk -F: '{print $1}' | sort| uniq -c|sort -n
```

- 查看apache详细链接情况

```bash
netstat -aptol
```

- 检测某一台服务器的端口是否开启状态

```bash
nc -v -w 10 -z 172.20.206.147 25801
```

- 验证apache2配置是否正确

```bash
httpd -t
```

#### apache模块

##### 模块介绍

Apache 各个模块功能: 

- 基本|(B)|模块默认包含，必须明确禁用；
- 扩展|(E)|/实验|(X)|模块默认不包含，必须明确启用。

| 模块名称            | 状态 | 简要描述                                                                     |
|:--------------------|:-----|:-----------------------------------------------------------------------------|
| mod_actions         | (B)  | 基于媒体类型或请求方法，为执行CGI脚本而提供                                  |
| mod_alias           | (B)  | 提供从文件系统的不同部分到文档树的映射和URL重定向                            |
| mod_asis            | (B)  | 发送自己包含HTTP头内容的文件                                                 |
| mod_auth_basic      | (B)  | 使用基本认证                                                                 |
| mod_authn_default   | (B)  | 在未正确配置认证模块的情况下简单拒绝一切认证信息                             |
| mod_authn_file      | (B)  | 使用纯文本文件为认证提供支持                                                 |
| mod_authz_default   | (B)  | 在未正确配置授权支持模块的情况下简单拒绝一切授权请求                         |
| mod_authz_groupfile | (B)  | 使用纯文本文件为组提供授权支持                                               |
| mod_authz_host      | (B)  | 供基于主机名、IP地址、请求特征的访问控制                                     |
| mod_authz_user      | (B)  | 基于每个用户提供授权支持                                                     |
| mod_autoindex       | (B)  | 自动对目录中的内容生成列表，类似于"ls"或"dir"命令                            |
| mod_cgi             | (B)  | 在非线程型MPM(prefork)上提供对CGI脚本执行的支持                              |
| mod_cgid            | (B)  | 在线程型MPM(worker)上用一个外部CGI守护进程执行CGI脚本                        |
| mod_dir             | (B)  | 指定目录索引文件以及为目录提供"尾斜杠"重定向                                 |
| mod_env             | (B)  | 允许Apache修改或清除传送到CGI脚本和SSI页面的环境变量                         |
| mod_filter          | (B)  | 根据上下文实际情况对输出过滤器进行动态配置                                   |
| mod_imagemap        | (B)  | 处理服务器端图像映射                                                         |
| mod_include         | (B)  | 实现服务端包含文档(SSI)处理                                                  |
| mod_isapi           | (B)  | 仅限于在Windows平台上实现ISAPI扩展                                           |
| mod_log_config      | (B)  | 允许记录日志和定制日志文件格式                                               |
| mod_mime            | (B)  | 根据文件扩展名决定应答的行为(处理器/过滤器)和内容(MIME类型/语言/字符集/编码) |
| mod_negotiation     | (B)  | 提供内容协商支持                                                             |
| mod_nw_ssl          | (B)  | 仅限于在NetWare平台上实现SSL加密支持                                         |
| mod_setenvif        | (B)  | 根据客户端请求头字段设置环境变量                                             |
| mod_status          | (B)  | 生成描述服务器状态的Web页面                                                  |
| mod_userdir         | (B)  | 允许用户从自己的主目录中提供页面(使用"/~username")                           |
| mod_auth_digest     | (X)  | 使用MD5摘要认证(更安全，但是只有最新的浏览器才支持)                          |
| mod_authn_alias     | (E)  | 基于实际认证支持者创建扩展的认证支持者，并为它起一个别名以便于引用           |
| mod_authn_anon      | (E)  | 提供匿名用户认证支持                                                         |
| mod_authn_dbd       | (E)  | 使用SQL数据库为认证提供支持                                                  |
| mod_authn_dbm       | (E)  | 使用DBM数据库为认证提供支持                                                  |
| mod_authnz_ldap     | (E)  | 允许使用一个LDAP目录存储用户名和密码数据库来执行基本认证和授权               |
| mod_authz_dbm       | (E)  | 使用DBM数据库文件为组提供授权支持                                            |
| mod_authz_owner     | (E)  | 基于文件的所有者进行授权                                                     |
| mod_cache           | (E)  | 基于URI键的内容动态缓冲(内存或磁盘)                                          |
| mod_cern_meta       | (E)  | 允许Apache使用CERN httpd元文件，从而可以在发送文件时对头进行修改             |
| mod_charset_lite    | (X)  | 允许对页面进行字符集转换                                                     |
| mod_dav             | (E)  | 允许Apache提供DAV协议支持                                                    |
| mod_dav_fs          | (E)  | 为mod_dav访问服务器上的文件系统提供支持                                      |
| mod_dav_lock        | (E)  | 为mod_dav锁定服务器上的文件提供支持                                          |
| mod_dbd             | (E)  | 管理SQL数据库连接，为需要数据库功能的模块提供支持                            |
| mod_deflate         | (E)  | 压缩发送给客户端的内容                                                       |
| mod_disk_cache      | (E)  | 基于磁盘的缓冲管理器                                                         |
| mod_dumpio          | (E)  | 将所有I/O操作转储到错误日志中                                                |
| mod_echo            | (X)  | 一个很简单的协议演示模块                                                     |
| mod_example         | (X)  | 一个很简单的Apache模块API演示模块                                            |
| mod_expires         | (E)  | 允许通过配置文件控制HTTP的"Expires:"和"Cache-Control:"头内容                 |
| mod_ext_filter      | (E)  | 使用外部程序作为过滤器                                                       |
| mod_file_cache      | (X)  | 提供文件描述符缓存支持，从而提高Apache性能                                   |
| mod_headers         | (E)  | 允许通过配置文件控制任意的HTTP请求和应答头信息                               |
| mod_ident           | (E)  | 实现RFC1413规定的ident查找                                                   |
| mod_info            | (E)  | 生成Apache配置情况的Web页面                                                  |
| mod_ldap            | (E)  | 为其它LDAP模块提供LDAP连接池和结果缓冲服务                                   |
| mod_log_forensic    | (E)  | 实现"对比日志"，即在请求被处理之前和处理完成之后进行两次记录                 |
| mod_logio           | (E)  | 对每个请求的输入/输出字节数以及HTTP头进行日志记录                            |
| mod_mem_cache       | (E)  | 基于内存的缓冲管理器                                                         |
| mod_mime_magic      | (E)  | 通过读取部分文件内容自动猜测文件的MIME类型                                   |
| mod_proxy           | (E)  | 提供HTTP/1.1的代理/网关功能支持                                              |
| mod_proxy_ajp       | (E)  | mod_proxy的扩展，提供Apache JServ Protocol支持                               |
| mod_proxy_balancer  | (E)  | mod_proxy的扩展，提供负载平衡支持                                            |
| mod_proxy_connect   | (E)  | mod_proxy的扩展，提供对处理HTTP CONNECT方法的支持                            |
| mod_proxy_ftp       | (E)  | mod_proxy的FTP支持模块                                                       |
| mod_proxy_http      | (E)  | mod_proxy的HTTP支持模块                                                      |
| mod_rewrite         | (E)  | 一个基于一定规则的实时重写URL请求的引擎                                      |
| mod_so              | (E)  | 允许运行时加载DSO模块                                                        |
| mod_speling         | (E)  | 自动纠正URL中的拼写错误                                                      |
| mod_ssl             | (E)  | 使用安全套接字层(SSL)和传输层安全(TLS)协议实现高强度加密传输                 |
| mod_suexec          | (E)  | 使用与调用web服务器的用户不同的用户身份来运行CGI和SSI程序                    |
| mod_unique_id       | (E)  | 为每个请求生成唯一的标识以便跟踪                                             |
| mod_usertrack       | (E)  | 使用Session跟踪用户(会发送很多Cookie)，以记录用户的点击流                    |
| mod_version         | (E)  | 提供基于版本的配置段支持                                                     |
| mod_vhost_alias     | (E)  | 提供大批量虚拟主机的动态配置支持                                             |

#### 性能调优

##### 模块启用/关闭

1. 启用压缩

> LoadModule deflate_module modules/mod_deflate.so

2. 启用重写

> LoadModule rewrite_module modules/mod_rewrite.so  

3. 启用默认扩展，支持在这里进行修改httpd主要配置

> Include conf/extra/httpd-default.conf  

4. 提供文件描述符缓存支持，从而提高Apache性能

> LoadModule file_cache_module modules/mod_file_cache.so

5. 启用基于URI键的内容动态缓冲(内存或磁盘)  

> LoadModule cache_module modules/mod_cache.so  

6. 启用基于磁盘的缓冲管理器

> LoadModule cache_disk_module modules/mod_cache_disk.so  

7. 基于内存的缓冲管理器

> LoadModule socache_memcache_module modules/mod_socache_memcache.so 

8. 屏蔽所有不必要的模块

> #LoadModule authn_file_module modules/mod_authn_file.so  
> #LoadModule authn_dbm_module modules/mod_authn_dbm.so  
> #LoadModule authn_anon_module modules/mod_authn_anon.so  
> #LoadModule authn_dbd_module modules/mod_authn_dbd.so  
> #LoadModule authn_socache_module modules/mod_authn_socache.so  
> #LoadModule authn_core_module modules/mod_authn_core.so  
> #LoadModule authz_host_module modules/mod_authz_host.so  
> #LoadModule authz_groupfile_module modules/mod_authz_groupfile.so  
> #LoadModule authz_user_module modules/mod_authz_user.so  
> #LoadModule authz_dbm_module modules/mod_authz_dbm.so  
> #LoadModule authz_owner_module modules/mod_authz_owner.so  
> #LoadModule authz_dbd_module modules/mod_authz_dbd.so  
> LoadModule authz_core_module modules/mod_authz_core.so  
> LoadModule access_compat_module modules/mod_access_compat.so  
> #LoadModule auth_basic_module modules/mod_auth_basic.so  
> #LoadModule auth_form_module modules/mod_auth_form.so  
> #LoadModule auth_digest_module modules/mod_auth_digest.so 

9. 已经过时屏蔽

> #LoadModule autoindex_module modules/mod_autoindex.so

10. 用于定义缺省文档index.php、index.jsp等

> LoadModule dir_module modules/mod_dir.so

11. 用于定义记录文件格式

> LoadModule log_config_module modules/mod_log_config.so

12. 定义文件类型的关联

> LoadModule mime_module modules/mod_mime.so

13. 减少10%左右的重复请求

> LoadModule expires_module modules/mod_expires.so  

14. 允许apache修改或清除传递到cgi或ssi页面的环境变量

> LoadModule env_module modules/mod_env.so

15. 根据客户端请求头字段设置环境变量，如果不需要则屏蔽掉

> #LoadModule setenvif_module modules/mod_setenvif.so

16. 生成描述服务器状态的页面

> #LoadModule status_module modules/mod_status.so

17. 别名

> LoadModule alias_module modules/mod_alias.so

18. url地址重写模块

> LoadModule rewrite_module modules/mod_rewrite.so

19. jk_mod 负载均衡调度模块

> LoadModule    jk_module modules/mod_jk.so

20. 过滤模块，使用缓存必须启用过滤模块

> LoadModule filter_module modules/mod_filter.so

21. 关闭服务器版本信息

> LoadModule version_module modules/mod_version.so

22. 自动修正用户输入的url错误

> LoadModule speling_module modules/mod_speling.so

##### apache2扩展配置文件说明

| 配置文件                      | 描述说明                     |
|:------------------------------|:-----------------------------|
| httpd-autoindex.conf          | 自动索引配置                 |
| httpd-dav.conf                | WebDAV配置                   |
| httpd-default.conf            | Apache的默认配置             |
| httpd-info.conf               | mod_status, mod_info模块配置 |
| httpd-languages.conf          | Apache多语言配置支持         |
| httpd-manual.conf             | 在网站上提供Apache手册       |
| httpd-mpm.conf                | 多路处理模块配置文件         |
| httpd-multilang-errordoc.conf | 实现多语言的错误信息         |
| httpd-ssl.conf                | SSL配置                      |
| httpd-userdir.conf            | 配置用户目录                 |
| httpd-vhosts.conf             | 虚拟主机配置                 |

##### 性能指标计算方法

提供下面这个公式，以供大家在平时或者日常需要进行的性能测试中作为一个参考。

公式:

> 计算平均的并发用户数：C = nL/T

说明:

> C是平均的并发用户数；n 是 login session 的数量；L 是 login session 的平均长度；T指考察的时间段长度。
> 并发用户数峰值：C’ ≈ C+3根号C
> C’指并发用户数的峰值，C就是公式（1）中得到的平均的并发用户数。该公式的得出是假设用户的 loginsession 产生符合泊松分布而估算得到的。

##### apache2自带的压力测试工具ab

ab最常用的语法格式:

```bash
ab -n XXX -c YYY -k http://hostname.port/path/filename
```

说明:

> -n XXX:表示最多进行XXX次测试。也就是下载filename文件XXX次。
> -c YYY:客户端并发连接个数。
> -k:启用HTTP KeepAlive功能。默认不启用KeepAlive功能。

ab必须安装在客户端上，并且客户端机器配置性能要高些。
比如我们要对http://hostname:port/file.com下载10000次进行测试，并发访问为60个，启用HTTP KeepAlive功能，则访问指令为:

```bash
ab -n 10000 -c 60 -k http://hostname:port/file.htm
```

##### Java的压力测试工具Jmeter

Jmeter 是apache开发的基于Java的压力测试工具。

##### apache多路处理器MPM

目前apache2.4版本已经event MPM纳入正式版，不再是实验状态。安装时，apache已经自动将event MPM一起安装进去，通过apachectl -l可以查看到event.c模块。由此可以看到，event MPM已经成为apache默认的MPM工作模式。

1. 启用MPM

> Include conf/extra/httpd-mpm.conf 

2. 配置evnet MPM参数

```xml
<IfModule event.c>  
    #default 3  
    ServerLimit           15  

    #default 256 MaxRequestWorkers (2.3版本叫MaxClients) <= ServerLimit * ThreadsPerChild  
    MaxRequestWorkers            960  

    #default 3  
    StartServers          3  

    #default 64  
    ThreadsPerChild       64  

    #default 100000  ThreadLimit >= ThreadsPerChild  
    ThreadLimit           64  

    #default 75  
    MinSpareThreads       32  

    #default 400  MaxSpareThreads >= (MinSpareThreads + ThreadsPerChild)  
    MaxSpareThreads      112  

    #(2.3版本叫MaxRequestsPerChild)default 0, at 200 r/s, 20000 r results in a process lifetime of 2 minutes  
    MaxConnectionsPerChild 10000  
</IfModule>
```

说明:

- StartServers：初始数量的服务器进程开始，默认为3个
- MinSpareThreads：最小数量的工作线程,保存备用，默认是75个线程
- MaxSpareThreads：最大数量的工作线程,保存备用，默认是250线程
- ThreadsPerChild：固定数量的工作线程在每个服务器进程，默认是25个
- MaxRequestWorkers：最大数量的工作线程，默认是400
 -MaxConnectionsPerChild：最大连接数的一个服务器进程服务，默认为0，没有上限限制，但是为了避免内存异常，影响稳定性，需要设置一个数值进行限制在修改配置后,需要停止htppd,再启动httpd，不能通过apacherestart生效，而是先 apache stop 然后再 apache start才可以生效。
- ServerLimit：ServerLimit是活动子进程数量的硬限制，它必须大于或等于MaxClients除以ThreadsPerChild的值。serverLimit最大20000，默认是16。只有在你需要将MaxClients和ThreadsPerChild设置成需要超过默认值16个子进程的时候才需要使用这个指令。不要将该指令的值设置的比MaxClients和ThreadsPerChild需要的子进程数量高。使用这个指令时要特别当心。如果将ServerLimit设置成一个高出实际需要许多的值，将会有过多的共享内存被分配。如果将ServerLimit和MaxClients设置成超过系统的处理能力，Apache可能无法启动，或者系统将变得不稳定。RedHat Enterprise LinuxAS3.0Update2最大MaxClients只能设置到256。如果你需要设置其为更高，需要在MaxClients前面添加:ServerLimitxxx其中xxx不能少于MaxClients的数值。该设置方法适用于Apache 2.0系列。
- ThreadLimit：ThreadLimit是所有服务线程总数的硬限制，它必须大于或等于ThreadsPerChild指令。使用这个指令时要特别当心。如果将ThreadLimit设置成一个高出ThreadsPerChild实际需要很多的值，将会有过多的共享内存被分配。如果将ThreadLimit和ThreadsPerChild设置成超过系统的处理能力，Apache可能无法启动，或者系统将变得不稳定。该指令的值应当和ThreadsPerChild可能达到的最大值保持一致。

##### 计算event的相关参数

1. 计算服务器进程的平均内存

```bash
ps aux | grep 'httpd' | awk '{print $6/1024 " MB";}'
```

2. 计算正在通讯传输过程中的进程的平均内存，最好在一天之内不同的时间段内运行以下代码

```bash
ps aux | grep 'httpd' | awk '{print $6/1024;}' | awk '{avg += ($1 - avg) / NR;} END {print avg " MB";}'
```

3. 通过上面两个指令计算出平均进程所使用的内存大小 ，再通过以下公式计算

> MaxRequestWorkers（MaxClients） =  (Total RAM - RAM used for Linux, MySQL, etc.) / Average httpd process size.

- StartServers: 30% of MaxRequestWorkers
- MinSpareThreads: 5% of MaxRequestWorkers
- MaxSpareThreads: 10% of MaxRequestWorkers
- ServerLimit: 100% of MaxRequestWorkers
- MaxConnectionsPerChild 10000 (as conservative alternative to address possible problem with memory leaky apps)

##### event MPM 与 worker MPM区别

可以支持比worker更高的并发数，主要安装在类unix/linux上的工作模式。event mpm是worker mpm的变种，但是具有比worker MPM更好的并发性能。在event mpm模式下，ssl是不被支持的，他会被切换到worker mpm下处理。event mpm在apache2.4版本时才被从实验状态转化成标准应用。

##### apache 缓存设置

apache涉及的缓存模块有mod_cache、mod_disk_cache、mod_file_cache、mod_mem_cache。如果要使用缓存必须启用这四个缓存模块。
同时修改缓存设置后，必须重启apache，刷新缓存，否则用户访问页面不是最新页面。

###### mod_cache、mod_disk_cache、mod_mem_cache、mod_file_cache关系

- apache缓存分为硬盘缓存和内存缓存
- mod_disk_cache mod_mem_cache 都依赖于mod_cache
- mod_file_cache是结合mod_cache使用，可以用于指定几个频繁访问，但是变化不大的文件

###### 配置硬盘缓存和内存缓存的缓存配置

```xml
<IfModule mod_cache.c>  
   #设置缓存过期时间，默认一小时  
   CacheDefaultExpire 3600  
   #设置缓存最大失效时间，默认最大一天  
   CacheMaxExpire 86400 CacheLastModifiedFactor 0.1 CacheIgnoreHeaders Set-Cookie CacheIgnoreCacheControl Off 
    <IfModule mod_disk_cache.c>  
       #启用缓存，并设定硬盘缓存目录（url路径）  
       CacheEnable disk /  
       #设定apache访问用户的缓存路径，需要进行授权配置，如linux设置为777  
       CacheRoot /home/apache/cache  
       #缓存目录深度  
       CacheDirLevels 5  
       #缓存目录名称字符串长度  
       CacheDirLength 5  
       #缓存文件最大值  
        CacheMaxFileSize 1048576  
       #缓存文件最小值
       CacheMinFileSize 10 
    </IfModule>

    <IfModule mod_mem_cache.c>  
       #缓存路径  
       CacheEnable mem /  

       #缓存对象最大个数  
       MCacheMaxObjectCount 20000  

       #单个缓存对象最大大小  
       MCacheMaxObjectSize 1048576  

       #单个缓存对象最小大小   
       MCacheMinObjectSize 10  

       #在缓冲区最多能够放置多少的将要被缓存对象的尺寸  
       MCacheMaxStreamingBuffer 65536  

       #清除缓存所使用的算法，默认是 GDSF，还有一个是LRU  
       MCacheRemovalAlgorithm GDSF  

       #缓存数据最多能使用的内存  
       MCacheSize 131072 
    </IfModule>
</IfModule>
```

###### 文件缓存的应用

1. 缓存文件

如果你的网站有几个文件的访问非常频繁而又不经常变动，则可以在 Apache 启动的时候就把它们的内容缓存到内存中（当然要启用内存缓存系统），使用的是 mod_file_cache 模块，有多个文件可以用空格格开，具体如下：

```xml
<IfModule mod_file_cache.c>  
     MMapFile /var/html/js/jquery.js  
</IfModule>  
```

2. 缓存句柄

```xml
<IfModule mod_file_cache.c>  
  CacheFile /usr/local/apache2/htdocs/index.html  
</IfModule>
```

##### apache压缩配置

apache通过mod_deflate模块实现页面压缩，要想进行页面压缩必须启用以下两个模块

- LoadModule deflate_module modules/mod_deflate.so  
- LoadModule filter_module modules/mod_filter.so

###### 页面压缩模块配置

```xml
<ifmodule mod_deflate.c>  
#设定压缩率，压缩率1 -9, 6是建议值，不能太高，消耗过多的内存，影响服务器性能  
DeflateCompressionLevel 6  

AddOutputFilterByType DEFLATE text/plain  
AddOutputFilterByType DEFLATE text/html  
AddOutputFilterByType DEFLATE text/php  
AddOutputFilterByType DEFLATE text/xml  
AddOutputFilterByType DEFLATE text/css  
AddOutputFilterByType DEFLATE text/javascript  
AddOutputFilterByType DEFLATE application/xhtml+xml  
AddOutputFilterByType DEFLATE application/xml  
AddOutputFilterByType DEFLATE application/rss+xml  
AddOutputFilterByType DEFLATE application/atom_xml  
AddOutputFilterByType DEFLATE application/x-javascript  
AddOutputFilterByType DEFLATE application/x-httpd-php  
AddOutputFilterByType DEFLATE application/x-httpd-fastphp  
AddOutputFilterByType DEFLATE application/x-httpd-eruby  
AddOutputFilterByType DEFLATE image/svg+xml  
AddOutputFilterByType DEFLATE application/javascript  

#插入过滤器  
SetOutputFilter DEFLATE  

#排除不需要压缩的文件  
SetEnvIfNoCase Request_URI \.(?:exe|t?gz|zip|bz2|sit|rar)$ no-gzip dont-vary  
SetEnvIfNoCase Request_URI \.(?:gif|jpe?g|png)$ no-gzip don’t-vary  
SetEnvIfNoCase Request_URI \.pdf$ no-gzip dont-vary  
SetEnvIfNoCase Request_URI \.avi$ no-gzip dont-vary  
SetEnvIfNoCase Request_URI \.mov$ no-gzip dont-vary  
SetEnvIfNoCase Request_URI \.mp3$ no-gzip dont-vary  
SetEnvIfNoCase Request_URI \.mp4$ no-gzip dont-vary  
SetEnvIfNoCase Request_URI \.rm$ no-gzip dont-vary  
</ifmodule>
```

##### keepAlive

在HTTP 1.0中和Apache服务器的一次连接只能发出一次HTTP请求，而KeepAlive参数支持HTTP 1.1版本的一次连接，多次传输功能，这样就可以在一次连接中发出多个HTTP请求。从而避免对于同一个客户端需要打开不同的连接。很多请求通过同一个 TCP连接来发送，可以节约网络和系统资源。

1. keepAlive启用场景

- 如果有较多的js,css,图片访问，则需要开启长链接
- 如果内存较少，大量的动态页面请求，文件访问，则关闭长链接，节省内存，提高apache访问的稳定性
- 如果内存充足，cpu较好，服务器性能优越，则是否开启长链接对访问性能都不会产生影响

2. keepAlive配置

  在Apache的配置文件httpd.conf中，设置：

- Timeout  60 默认为60s修改为30s
- KeepAlive on  设置为on状态
- KeepAliveTimeout 默认为5s,如果值设置过高，由于每个进程都要保持一定时间对应该用户，而无法应付其他用户请求访问，从而导致服务器性能下降。
- MaxKeepAliveRequests 100  如果设置为0表示无限制，建议最好设置一个值, 把MaxKeepAliveRequests设置的尽量大，可以在一次连接中进行更多的HTTP请求。但在我们的测试中还发现，把 MaxKeepAliveRequests设置成1000，则评测的客户端容易出现“Send requesttimed out”的错误，所以具体数值还要根据自己的情形来设置。

#### 问题集锦

1. 加载

> LoadModule authz_core_module modules/mod_authz_core.so  
> Invalid command 'Require', perhaps misspelled or defined by a module not included in the server configuration`

2. 配置信息后面不能跟随注释，注释必须另起一行

> CacheDefaultExpire takes one argument, The default time in seconds to cache a document

3. 关键字错误 AddOutputFileByType 应该是

> AddOutputFitlerByType  
> Invalid command 'AddOutputFileByType', perhaps misspelled or defined by a module not included in the server configuration

4. 启用

> LoadModule setenvif_module modules/mod_setenvif.so  
> Invalid command 'SetEnvIfNoCase', perhaps misspelled or defined by a module not included in the server configuration

5. ifModule注释不能跟在配置参数后面，否则会导致配置解析失败

> AH00526: Syntax error on line 558 of /usr/local/cp-httpd-2.4.18/conf/httpd.conf:
> CacheDefaultExpire takes one argument, The default time in seconds to cache a document