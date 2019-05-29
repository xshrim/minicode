#!/bin/bash

echo "============================== Start to install MySQL-5.7.24 =============================="
# wget https://dev.mysql.com/get/Downloads/MySQL-5.7/mysql-5.7.24-linux-glibc2.12-x86_64.tar.gz
if [[ ! -f mysql-5.7.24-linux-glibc2.12-x86_64.tar.gz ]];
then
    echo 'mysql tar file not exist'
    exit
fi

tar -xvf mysql-5.7.24-linux-glibc2.12-x86_64.tar.gz

groupadd -g 801 mysql    
useradd -u 801 -g mysql -s /bin/bash mysql

mkdir -p /dbdata/{dbdata,dblog}

mv mysql-5.7.24-linux-glibc2.12-x86_64 /dbdata/mysql

touch /dbdata/dblog/error.log

chown -R mysql:mysql /dbdata/mysql    
chown -R mysql:mysql /dbdata/dbdata
chown -R mysql:mysql /dbdata/dblog

echo "export PATH=/dbdata/mysql/bin:$PATH" >> /etc/profile
source /etc/profile

mysqld --initialize-insecure --user=mysql --basedir=/dbdata/mysql/ --datadir=/dbdata/dbdata/ 

if [[ -f /etc/my.cnf ]];
then
    mv /etc/my.cnf /etc/my.cnf.old
fi

cat > /etc/my.cnf <<EOF
[mysqld]
 
sql_mode=NO_ENGINE_SUBSTITUTION,STRICT_TRANS_TABLES 
 
basedir = /dbdata/mysql                                 
datadir = /dbdata/dbdata                               
log-error = /dbdata/dblog/error.log                   
log-bin = /dbdata/dblog/mysql-bin                  
log-bin-index = /dbdata/dblog/mysql-bin.index    
max_binlog_size = 512M                         
binlog_cache_size = 16M                      
long_query_time = 2                               
socket = /dbdata/dbdata/mysql.sock               
pid-file = /dbdata/dbdata/mysql.pid             
port = 3306                                    
character-set-server=utf8                     
skip-grant-tables                            

back_log = 1024                             
max_connections = 3000                     
max_connect_errors = 50                   
table_open_cache = 4096                  
max_allowed_packet = 256M               
 
max_heap_table_size = 256M             
read_rnd_buffer_size = 32M            
sort_buffer_size = 256M              
join_buffer_size = 16M              
thread_cache_size = 16             
query_cache_size = 128M           
query_cache_limit = 4M       
ft_min_word_len = 2         
 
thread_stack = 512K        
transaction_isolation = READ-COMMITTED   
tmp_table_size = 256M                   
long_query_time = 6                    
 
server_id=1                           
 
innodb_buffer_pool_size = 2G         
innodb_thread_concurrency = 16
innodb_log_buffer_size = 64M        
 
innodb_log_file_size = 512M        
innodb_log_files_in_group = 3     
innodb_max_dirty_pages_pct = 90  
innodb_lock_wait_timeout = 120  
innodb_file_per_table = on     
 
[mysqldump]
quick                         
 
max_allowed_packet = 128M    
 
[mysql]
auto-rehash                 
default-character-set=utf8  
socket = /dbdata/dbdata/mysql.sock
 
[myisamchk]
key_buffer = 16M
sort_buffer_size = 16M
read_buffer = 8M
write_buffer = 8M
 
[mysqld_safe]
open-files-limit = 10240
EOF

cp -a /dbdata/mysql/support-files/mysql.server /etc/init.d/mysqld  
sed -i 's/^basedir=$/basedir=\/dbdata\/mysql/g' /etc/init.d/mysqld 
sed -i 's/^datadir=$/datadir=\/dbdata\/dbdata/g' /etc/init.d/mysqld
sed -i 's/^mysqld_pid_file_path=$/mysqld_pid_file_path=\/dbdata\/dbdata\/mysql.pid/g' /etc/init.d/mysqld
service mysqld start
service mysqld status

if [[ $(ps aux|grep mysql.pid|wc -l) -eq 0 ]];
then
    echo "mysql service start <failed>"
    exit
fi

mysql -uroot -p'\n' <<EOF
use mysql;
UPDATE mysql.user SET authentication_string = PASSWORD('root') WHERE user = 'root' AND host = 'localhost';
FLUSH PRIVILEGES;
EOF

sed -i '/skip-grant-tables/d' /etc/my.cnf

service mysqld restart

echo "============================== MySQL-5.7.24 has been installed successfully =============================="
echo "ROOT account password: root"
echo "Login command: mysql -uroot -proot"
