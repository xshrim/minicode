# MySQL安装配置(单机)

## 前置准备

### 系统环境

- 操作系统: CentOS 7.4 (x86_x64)
- 数据库: MySQL 5.7.24 (linux-glibc1.2-x86_64)

### 存储规划

| 文件系统       | 空间(GB) | 所属用户 | 所属组 | 目录权限 | 存储要求                   | 描述信息          |
|:---------------|:---------|:---------|:-------|:---------|:---------------------------|:------------------|
| /dbdata/mysql  |          | mysql    | mysql  |          | 本地/数据库HA/网络文件系统 | MySQL安装目录     |
| /dbdata/dbdata |          | mysql    | mysql  |          | 本地/数据库HA/网络文件系统 | MySQL数据文件目录 |
| /dbdata/dblog  |          | mysql    | mysql  |          | 本地/数据库HA/网络文件系统 | MySQL日志文件目录 |

## 安装配置

### 用户和组添加

```bash
groupadd -g 801 mysql                      # 添加mysql
useradd -u 801 -g mysql -s /bin/bash mysql # 添加属组为mysql的mysql用户
```

[注]: 后续操作仍在root账号下执行

### 目录创建

```bash
mkdir -p /dbdata/{mysql,dbdata,dblog}      # 创建mysql相关目录
touch /dbdata/dblog/error.log              # 创建错误日志文件
chown mysql:mysql /dbdata/mysql            # 更改mysql相关目录属组和用户
chown mysql:mysql /dbdata/dbdata
chown mysql:mysql /dbdata/dblog
```

### 软件安装

```bash
wget https://dev.mysql.com/get/Downloads/MySQL-5.7/mysql-5.7.24-linux-glibc2.12-x86_64.tar.gz    # 下载linux通用版本
tar -xvf mysql-5.7.24-linux-glibc2.12-x86_64.tar.gz          # 解压安装文件
mv mysql-5.7.24-linux-glibc2.12-x86_64/* /dbdata/mysql/ 
chown -R mysql:mysql /dbdata/mysql/
echo "export PATH=/dbdata/mysql/bin:$PATH" >> /etc/profile    # 配置环境变量
source /etc/profile
mysql --version
```

### 初始化

```bash
mysqld --initialize-insecure --user=mysql --basedir=/dbdata/mysql/ --datadir=/dbdata/dbdata/     # --initialize-insecure选项在初始化数据库的同时为localhost上的root账号提供一个空密码
```

### 配置文件

```bash
vi /etc/my.cnf
```
[mysqld]

sql_mode=NO_ENGINE_SUBSTITUTION,STRICT_TRANS_TABLES 

basedir = /dbdata/mysql                                 # 数据库安装目录
datadir = /dbdata/dbdata                                # 数据文件目录
log-error = /dbdata/dblog/error.log                     # 记录错误日志
log-bin = /dbdata/dblog/mysql-bin                       # 开启log-bin
log-bin-index = /dbdata/dblog/mysql-bin.index           # log-bin index
max_binlog_size = 512M                                  # bin log最大值, 超出此值bin log将滚动覆盖
binlog_cache_size = 16M                                 # bin log缓存大小
long_query_time = 2                                     # 执行时间超过2s的sql记入慢查询日志
socket = /dbdata/dbdata/mysql.sock                      # 本地连接套接字路径, 需与[mysql]客户端部分的socket路径一致, 建议使用默认值
pid-file = /dbdata/dbdata/mysql.pid                     # pid文件
port = 3306                                             # 数据库连接端口
character-set-server=utf8                               # 数据库级别服务端字符编码
skip-grant-tables                                       # 忽略授权表, 免密登录数据库, 需注释掉

back_log = 1024                                         # 短时间内请求堆栈可以容纳的请求数量
max_connections = 3000                                  # 最大连接数
max_connect_errors = 50                                 # 连接重试次数
table_open_cache = 4096                                 # 表描述符缓存大小, 可减少文件打开/关闭次数
max_allowed_packet = 256M                               # 一次数据传输可接收的数据包大小

max_heap_table_size = 256M                              # 独立内存表大小上限
read_rnd_buffer_size = 32M                              # 随机读缓冲区大小, 提高查询速度, 每一个连接均会分配此缓存
sort_buffer_size = 256M                                 # 排序缓冲区大小, 值越大越能避免磁盘IO, 排序速度越快
join_buffer_size = 16M                                  # 表连接缓冲区大小,值越大表连接执行速度越快
thread_cache_size = 16                                  # 线程缓存数量, 提高重新连接的速度
query_cache_size = 128M                                 # 查询缓存大小, 提高重复查询效率(表未发生变化的情况下)
query_cache_limit = 4M                                  # 单个查询可使用的查询缓存的大小
ft_min_word_len = 2                                     # 开启全文索引, 设置索引关键字最小长度, 通常需注释

thread_stack = 512K                                     # 一个线程栈占用内存上限
transaction_isolation = READ-COMMITTED                  # 事务隔离级别, 不可重复读, 解决脏读问题, 如需要MVCC支持,需使用可重复读
tmp_table_size = 256M                                   # 内部内存临时表的最大值, 每个线程均分配
long_query_time = 6                                     # 查询超时秒数

server_id=1                                             # 服务端ID, 高可用需为每个服务端设置不同ID

innodb_buffer_pool_size = 2G                            # InnoDB高速数据和索引缓冲区大小上限, 值越大越能避免磁盘IO
innodb_thread_concurrency = 16
innodb_log_buffer_size = 64M                            # 内存中日志文件缓存写入磁盘前的临界值

innodb_log_file_size = 512M                             # 日志组中每个日志文件大小
innodb_log_files_in_group = 3                           # 日志组数量
innodb_max_dirty_pages_pct = 90                         # 内存中脏页占内存总数据页百分比上限, 超过此值脏数据将强制落盘
innodb_lock_wait_timeout = 120                          # 事务发生阻塞回滚前等待的超时时间
innodb_file_per_table = on                              # 每个表拥有独立的表空间, 分配单独的数据文件和表描述文件

[mysqldump]
quick                                                   # mysqldump从服务器取得数据后直接输出而不必先缓存到内存

max_allowed_packet = 128M                               # 服务端接受mysqldump数据包大小上限

[mysql]
auto-rehash                                             # 开启客户端自动补全
default-character-set=utf8                              # 客户端字符集
socket = /dbdata/dbdata/mysql.sock                      # 客户端与服务端建立本地socket连接的sock文件路径

[myisamchk]                                             # 数据库信息统计与优化修复相关参数
key_buffer = 16M
sort_buffer_size = 16M
read_buffer = 8M
write_buffer = 8M

[mysqld_safe]
open-files-limit = 10240                                # mysql最大打开文件数
```

```

[注]: 以上仅主要的配置项和调优项, 可根据需要调整更多调优内容, 日志记录部分可按需设置需要记录的日志类型.

### 服务启动

```bash
cp -a /dbdata/mysql/support-files/mysql.server /etc/init.d/mysqld            # 拷贝mysqld服务文件, 建议重命名为mysqld而不是mysql, 因为mysql服务默认指向mariadb
sed -i 's/^basedir=$/basedir=\/dbdata\/mysql/g' /etc/init.d/mysqld           # 根据my.cnf中的自定义设置, 修改服务脚本内容
sed -i 's/^datadir=$/datadir=\/dbdata\/dbdata/g' /etc/init.d/mysqld
sed -i 's/^mysqld_pid_file_path=$/mysqld_pid_file_path=\/dbdata\/dbdata\/mysql.pid/g' /etc/init.d/mysqld
service mysqld start
service mysqld status
```

### 修改密码
```bash
mysql -uroot -p'\n'                                      # 使用-insecures初始化的mysql的root账号密码为空
> use mysql;
> select * from user;
> UPDATE mysql.user SET authentication_string = PASSWORD('root') WHERE user = 'root' AND host = 'localhost';
> FLUSH PRIVILEGES;
> exit;
```

注释掉my.cnf文件中 skip-grant-tables 项, 重启mysql服务即可实现安全登录.

[注]: 如需开启远程访问, 向mysql库的user表中更新或插入可远程访问的主机ip或ip段即可.
