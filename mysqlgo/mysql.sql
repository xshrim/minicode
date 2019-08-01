create database blockchain;
use blockchain;
create table user ( id int(11) not null auto_increment, name varchar(45) default null, password varchar(45) default null, dt datetime default null, flag bit(1) default b'0', primary key(id), info varchar(500) default null )engine=innodb auto_increment=74 default charset =utf8; 
