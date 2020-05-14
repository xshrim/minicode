import pymysql
import tokengen
'''
# mariadb

systemctl start mariadb

mysql -uroot -p     (passwd: root)

> explain select idhash, count(*) as count from token group by idhash having count > 1;

> select count(1) from token;

> create index caddr on token (caddr);

> show index from token;

> alter table token add id BIGINT(40);

> alter table token change id id BIGINT(40) not null auto_increment primary key;

> select * from token where caddr='0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa';
'''

settings = {
    "host": '127.0.0.1',
    "port": 3306,
    "dbname": "test",
    "charset": "utf8",
    "username": "root",
    "passwd": "root",
}


class MySqlDB(object):
    def __init__(self, log=False, host=None, port=None, dbname=None, charset=None, username=None, passwd=None):
        if host is None:
            host = settings["host"]
        if port is None:
            port = settings["port"]
        if dbname is None:
            dbname = settings["dbname"]
        if charset is None:
            charset = settings["charset"]
        if username is None:
            username = settings["username"]
        if passwd is None:
            passwd = settings["passwd"]
        try:
            self.log = log
            self.db = pymysql.connect(
                host=host, port=port, user=username, passwd=passwd, db=dbname,
                charset=charset)  # cursorclass=pymysql.cursors.DictCursor
            self.cursor = self.db.cursor()
        except Exception as ex:
            print(str(ex))

    def display(self, msg):
        if self.log:
            print(msg)

    def create(self, sql):
        res = -1
        self.display("create...")
        try:
            res = self.cursor.execute(sql)
            self.db.commit()
        except Exception as ex:
            print(str(ex))
            self.db.rollback()
        return res

    def insert(self, sql):
        res = -1
        self.display("insert...")
        try:
            res = self.cursor.execute(sql)
            self.db.commit()
        except Exception as ex:
            self.db.rollback()
        return res

    def insertMany(self, sql, datas):
        res = -1
        self.display("insert...")
        try:
            res = self.cursor.executemany(sql, datas)
            self.db.commit()
        except Exception as ex:
            print(str(ex))
            self.db.rollback()
        return res

    def update(self, sql):
        self.display("update...")
        res = -1
        try:
            res = self.cursor.execute(sql)
            self.db.commit()
        except Exception as ex:
            self.db.rollback()
        return res

    def delete(self, sql):
        self.display("delete...")
        res = -1
        try:
            res = self.cursor.execute(sql)
            self.db.commit()
        except Exception as ex:
            self.db.rollback()
        return res

    def drop(self, sql):
        self.display("clear...")
        res = -1
        try:
            res = self.cursor.execute(sql)
            self.db.commit()
        except Exception as ex:
            self.db.rollback()
        return res

    def select(self, sql):
        self.display("select...")
        datas = []
        try:
            self.cursor.execute(sql)
            for data in self.cursor:
                datas.append(data)
        except Exception as ex:
            self.db.rollback()
        return datas


mysql = MySqlDB()

drop_sql = '''
    DROP TABLE IF EXISTS `token`;
    '''

# print(mysql.drop(drop_sql))

create_sql = '''
    CREATE TABLE IF NOT EXISTS `token` (
   `idhash` VARCHAR(40) NOT NULL DEFAULT '0x00000000000000000000000000000000',
   `parenthash` VARCHAR(40) NOT NULL DEFAULT '0x00000000000000000000000000000000',
   `origincode` VARCHAR(40) NOT NULL DEFAULT '0x00000000000000000000000000000000',
   `statuscode` VARCHAR(40) NOT NULL DEFAULT '0x00000000000000000000000000000000',
   `typecode` VARCHAR(40) NOT NULL DEFAULT '0x00000000000000000000000000000000',
   `itemcode` VARCHAR(40) NOT NULL DEFAULT '0x00000000000000000000000000000000',
   `unitcode` VARCHAR(40) NOT NULL DEFAULT '0x00000000000000000000000000000000',
   `saddr` VARCHAR(50) NOT NULL DEFAULT '0x0000000000000000000000000000000000000000',
   `caddr` VARCHAR(50) NOT NULL DEFAULT '0x0000000000000000000000000000000000000000',
   `snetwork` BIGINT(40) NOT NULL DEFAULT 0,
   `cnetwork` BIGINT(40) NOT NULL DEFAULT 0,
   `amount` BIGINT(40) NOT NULL DEFAULT 0,
   `gtimestamp` BIGINT(40) NOT NULL DEFAULT 0,
   `ftimestamp` BIGINT(40) NOT NULL DEFAULT 0,
   `etimestamp` BIGINT(40) NOT NULL DEFAULT 0,
   `ctimestamp` BIGINT(40) NOT NULL DEFAULT 0,
   `stimestamp` BIGINT(40) NOT NULL DEFAULT 0,
   `status` BIGINT(40) NOT NULL DEFAULT 0
);
'''

# print(mysql.create(create_sql))

count = 9000000
tklist = []
insert_sql = 'insert into token values(%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)'
for i in range(0, count):
    tk = eval(str(tokengen.newToken()))
    print(str(i).rjust(8) + " ==> " + tk['idhash'])
    tklist.append((tk['idhash'], tk['parenthash'], tk['origincode'], tk['statuscode'], tk['typecode'], tk['itemcode'],
                   tk['unitcode'], tk['saddr'], tk['caddr'], tk['snetwork'], tk['cnetwork'], tk['amount'], tk['gtimestamp'],
                   tk['ftimestamp'], tk['etimestamp'], tk['ctimestamp'], tk['stimestamp'], tk['status']))
    if len(tklist) >= 50 or i == count - 1:
        mysql.insertMany(insert_sql, tklist)
        # print(mysql.insertMany(insert_sql, tklist))
        tklist = []
