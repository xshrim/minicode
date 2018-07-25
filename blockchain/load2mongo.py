from pymongo import MongoClient
import tokengen
'''
# mongoDB
mongod --dbpath=/home/xshrim/data/db

mongo
> use admin
> db.createUser(
  {
    user: "root",
    pwd: "root",
    roles: [ { role: "root", db: "admin" } ]
  }
)
> db.createUser(
  {
    user: "admin",
    pwd: "123456",
    roles: [ { role: "userAdminAnyDatabase", db: "admin" } ]
  }
)
> use test
> db.createUser(
  {
    user: "tester",
    pwd: "123",
    roles: [ { role: "readWrite", db: "test" },
             { role: "read", db: "reporting" } ]
  }
)
> exit

mongod --auth --bind_ip_all --port 27017 --fork --dbpath /home/xshrim/data/db --logpath /root/mongod.log

mongo --port 27017 -u "tester" -p "123" --authenticationDatabase "test"

> db.tc.find({idhash: "0xa32affee707a11e895b528e347233abc"}).explain('executionStats')

> db.tc.find({gtimestamp: {$gt: 1529243511}}).explain('executionStats')

> db.tc.createIndex({caddr: 1})

> db.tc.getIndexes()

> db.tc.aggregate([{$group: {_id: "$idhash", count: {$sum: 1}}}, {$match: {count: {$gt: 1}}}], {allowDiskUse: true})

> db.tc.update({}, {$set: {id: 0}}, {multi: 1})

> db.tc.update({},{$unset:{'id': 0}},false, true)

> db.tc.explain('executionStats').find({"caddr": "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"})

> use admin
> db.shutdownServer()
'''

# MongoDB
''' 本机环境
settings = {
    "host": '127.0.0.1',
    "port": 27017,
    "dbname": "test",
    "setname": "tc",
    "username": "tester",
    "passwd": "123",
}
'''
''' 甲方环境'''
settings = {
    "host": '192.168.0.98',
    "port": 27017,
    "dbname": "rmb",
    "setname": "token",
    "username": "rmber",
    "passwd": "rmber",
}


class MyMongoDB(object):
    def __init__(self, log=False, host=None, port=None, dbname=None, setname=None, username=None, passwd=None):
        if host is None:
            host = settings["host"]
        if port is None:
            port = settings["port"]
        if dbname is None:
            dbname = settings["dbname"]
        if setname is None:
            setname = settings["setname"]
        if username is None:
            username = settings["username"]
        if passwd is None:
            passwd = settings["passwd"]
        try:
            self.log = log
            self.client = MongoClient(host, port)
            self.db = self.client[dbname]
            self.db.authenticate(username, passwd)
            self.dbset = self.db[setname]
        except Exception as ex:
            print(str(ex))

    def display(self, msg):
        if self.log:
            print(msg)

    def insert(self, dic):
        self.display("insert...")
        return self.dbset.insert_many(dic)

    def update(self, dic, newdic):
        self.display("update...")
        return self.dbset.update_many(dic, newdic)

    def delete(self, dic):
        self.display("delete...")
        return self.dbset.delete_many(dic)

    def clear(self):
        self.display("clear...")
        return self.dbset.delete_many({})

    def dbfind(self, dic):
        self.display("find...")
        return self.dbset.find(dic)


'''
def getDbSet(host="127.0.0.1", port=27017, dbname="test", user="tester", passwd="123", collection="tc"):
    client = MongoClient("127.0.0.1", 27017)
    db = client[dbname]
    db.authenticate(user, passwd)
    dbset = db[collection]
    return dbset


def load(token, mongo):
    print(mongo.insert([eval(str(token))]).inserted_ids)
'''

mongo = MyMongoDB()

for i in range(0, 3):
    print(str(i).rjust(8) + " ==> " + str(mongo.insert([eval(str(tokengen.newToken()))]).inserted_ids[0]))
    # load(newToken(), mongo)

# print(mongo.clear().deleted_count)