# Tomcat安装配置

## 前置准备

### 系统环境

- 操作系统: CentOS7.4 (x86_64)
- JDK: jdk-8u191-linux-x64.tar.gz (可选)
- Tomcat: apache-tomcat-8.5.35.tar.gz

### 存储规划

| 文件系统                | 空间(GB) | 所属用户 | 所属组 | 目录权限 | 存储要求                   | 描述信息       |
|:------------------------|:---------|:---------|:-------|:---------|:---------------------------|:---------------|
| /tomcat                 |          | root     | root   |          | 本地/数据库HA/网络文件系统 | Tomcat安装目录 |
| /usr/local/jdk1.8.0_191 |          | root     | root   |          | 本地/数据库HA/网络文件系统 | JDK安装目录    |

## 安装配置

### JDK安装

```bash
tar xvf jdk-8u191-linux-x64.tar.gz
mv jdk1.8.0_191 /usr/local/
echo "JAVA_HOME=/usr/local/jdk1.8.0_191" >> /etc/profile           # 配置JDK环境变量
echo "PATH=$JAVA_HOME/bin:$PATH" >> /etc/profile
echo "export CLASSPATH=.:$JAVA_HOME/lib/dt.jar:$JAVA_HOME/lib/tools.jar"
source /etc/profile
java --version
```

### 软件安装

```bash
wget http://mirrors.shu.edu.cn/apache/tomcat/tomcat-8/v8.5.35/bin/apache-tomcat-8.5.35.tar.gz
tar xvf apache-tomcat-8.5.35.tar.gz
mv apache-tomcat-8.5.35 /tomcat
echo "export CATALINA_HOME=/tomcat" >> /etc/profile                # 配置Tomcat环境变量
echo 'export CATALINA_OPTS="-server -Xms2048m -Xmx4096m -XX:PermSize=1024m -XX:MaxPermSize=1024m"' >> /etc/profile  # Tomcat性能优化
source /etc/profile
cp /tomcat/bin/catalina.sh /etc/init.d/tomcat                      # 配置系统服务
sed -i '2i# chkconfig: 112 63 37' /etc/init.d/tomcat
sed -i '3i# . /etc/profile' /etc/init.d/tomcat
chmod 755 /etc/init.d/tomcat
# service tomcat start
systemctl start tomcat                                             # 启动tomcat服务
# chkconfig --add tomcat
# chkconfig tomcat on
systemctl enable tomcat                                            # 设置tomcat开机启
# chkconfig --list  
systemctl list-units --type service --all|grep tomcat
firewall-cmd --zone=public --add-port=8080/tcp --permanent         # 永久开启tomcat对外服务端口
firewall-cmd reload
firewall-cmd --query-port=8080/tcp --permanent
curl $(hostname -I | cut -d ' ' -f 1):8080                         # 验证服务是否正常
$CATALINA_HOME/bin/version.sh
```

### 安全配置与优化

#### 配置server.xml

##### 禁用8005端口

telnet localhost 8005 然后输入 SHUTDOWN 就可以关闭 Tomcat，为了安全我们要禁用该功能

```
默认值:
<Server port="8005" shutdown="SHUTDOWN">

修改为:
<Server port="-1" shutdown="SHUTDOWN">
```

##### 应用程序安全&关闭自动部署

```
默认值:
<Host name="localhost" appBase="webapps"
 unpackWARs="true" autoDeploy="true">

修改为:
<Host name="localhost" appBase="webapps"
 unpackWARs="false" autoDeploy="false" reloadable="false">
 ```

##### maxThreads 连接数限制修改配置

```
默认值:
<!--
 <Executor name="tomcatThreadPool" namePrefix="catalina-exec-"
 maxThreads="150" minSpareThreads="4"/>
 -->

修改为:
<Executor
 name="tomcatThreadPool"
 namePrefix="catalina-exec-"
 maxThreads="500"
 minSpareThreads="30"
 maxIdleTime="60000"
 prestartminSpareThreads = "true"
 maxQueueSize = "100"
/>
```

[参数解释]：
maxThreads：最大并发数，默认设置 200，一般建议在 500 ~ 800，根据硬件设施和业务来判断
minSpareThreads：Tomcat 初始化时创建的线程数，默认设置 25
maxIdleTime：如果当前线程大于初始化线程，那空闲线程存活的时间，单位毫秒，默认60000=60秒=1分钟。
prestartminSpareThreads：在 Tomcat 初始化的时候就初始化 minSpareThreads 的参数值，如果不等于 true，minSpareThreads 的值就没啥效果了
maxQueueSize：最大的等待队列数，超过则拒绝请求

##### Connector参数优化配置

```
默认值:
<Connector 
 port="8080" 
 protocol="HTTP/1.1" 
 connectionTimeout="20000" 
 redirectPort="8443" 
 />

修改为:
<Connector
 executor="tomcatThreadPool"
 port="8080"
 protocol="org.apache.coyote.http11.Http11Nio2Protocol"
 connectionTimeout="60000"
 maxConnections="10000"
 redirectPort="8443"
 enableLookups="false"
 acceptCount="100"
 maxPostSize="10485760"
 maxHttpHeaderSize="8192"
 compression="on"
 disableUploadTimeout="true"
 compressionMinSize="2048"
 acceptorThreadCount="2"
 compressableMimeType="text/html,text/plain,text/css,application/javascript,application/json,application/x-font-ttf,application/x-font-otf,image/svg+xml,image/jpeg,image/png,image/gif,audio/mpeg,video/mp4"
 URIEncoding="utf-8"
 processorCache="20000"
 tcpNoDelay="true"
 connectionLinger="5"
 server="Server Version 11.0"
 />
 ```

[参数解释]：
protocol：Tomcat 8 设置 nio2 更好：org.apache.coyote.http11.Http11Nio2Protocol
protocol：Tomcat 6 设置 nio 更好：org.apache.coyote.http11.Http11NioProtocol
protocol：Tomcat 8 设置 APR 性能飞快：org.apache.coyote.http11.Http11AprProtocol
connectionTimeout：Connector接受一个连接后等待的时间(milliseconds)，默认值是60000。
maxConnections：这个值表示最多可以有多少个socket连接到tomcat上
enableLookups：禁用DNS查询
acceptCount：当tomcat起动的线程数达到最大时，接受排队的请求个数，默认值为100。
maxPostSize：设置由容器解析的URL参数的最大长度，-1(小于0)为禁用这个属性，默认为2097152(2M) 请注意， FailedRequestFilter 过滤器可以用来拒绝达到了极限值的请求。
maxHttpHeaderSize：http请求头信息的最大程度，超过此长度的部分不予处理。一般8K。
compression：是否启用GZIP压缩 on为启用（文本数据压缩） off为不启用， force 压缩所有数据
disableUploadTimeout：这个标志允许servlet容器使用一个不同的,通常长在数据上传连接超时。 如果不指定,这个属性被设置为true,表示禁用该时间超时。
compressionMinSize：当超过最小数据大小才进行压缩
acceptorThreadCount：用于接受连接的线程数量。增加这个值在多CPU的机器上,尽管你永远不会真正需要超过2。默认值是1。
compressableMimeType：配置想压缩的数据类型
URIEncoding：网站一般采用UTF-8作为默认编码。
processorCache：协议处理器缓存的处理器对象来提高性能。 该设置决定多少这些对象的缓存。-1意味着无限的,默认是200。 如果不使用Servlet 3.0异步处理,默认是使用一样的maxThreads设置。 如果使用Servlet 3.0异步处理,默认是使用大maxThreads和预期的并发请求的最大数量(同步和异步)。
tcpNoDelay：如果设置为true,TCP_NO_DELAY选项将被设置在服务器套接字,而在大多数情况下提高性能。这是默认设置为true。
connectionLinger：秒数在这个连接器将持续使用的套接字时关闭。默认值是 -1,禁用socket 延迟时间。
server：隐藏Tomcat版本信息，首先隐藏HTTP头中的版本信息

#### 隐藏或修改Tomcat版本号

```bash
cd /tomcat/lib/
unzip catalina.jar
cd org/apache/catalina/util
cat ServerInfo.properties                  # 将去掉或修改以下参数
    server.info=Apache Tomcat/8.5.35
    server.number=8.5.35.0
    server.built=Jun 21 2017 17:01:09 UTC
```

#### 删除禁用默认管理页面以及相关配置文件

```bash
rm -rf /usr/local/apache-tomcat-8.5.35/webapps/*
rm -rf /usr/local/apache-tomcat-8.5.35/conf/tomcat-users.xml
```
