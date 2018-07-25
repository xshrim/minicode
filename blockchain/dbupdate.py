import pymysql
from pymongo import MongoClient


def updateMySql(start, interval, newaddr):

    db = pymysql.connect(
        host='127.0.0.1',
        port=3306,
        user='root',
        passwd='root',
        db='test',
        charset='utf8',
        cursorclass=pymysql.cursors.DictCursor)
    cursor = db.cursor()

    # count = 10000000
    idx = start
    sql = "select count(1) as count from token"
    cursor.execute(sql)
    count = cursor.fetchone()['count']

    while idx <= count:
        # sql = "select idhash from token limit " + str(current) + ", 1000"  # has index on idhash and caddr columns

        # the table token should have a auto increment primary key named id to run the follow code
        # > alter table token add id BIGINT(40);
        # > alter table token change id id BIGINT(40) not null auto_increment primary key;
        sql = "select id, idhash, caddr from token where `id`=" + str(idx)
        cursor.execute(sql)
        item = cursor.fetchone()  # cursor._rows[0]
        print(str(idx).rjust(8, '0') + ' => ' + item['idhash'] + ' : ' + item['caddr'] + "/" + newaddr, end=' # ')
        # print(str(current).rjust(8, '0') + ' => ' + item['idhash'] + ' : ' + '0x' + 'a' * 40, end=' # ')
        sql = "update token set caddr = %s where idhash = %s"
        print(cursor.execute(sql, (newaddr, item['idhash'])))
        idx += interval
    db.commit()


def updateMongo(start, interval, newaddr):
    client = MongoClient('127.0.0.1', 27017)
    db = client['test']
    db.authenticate('tester', '123')
    dbset = db['tc']

    idx = start
    count = dbset.find().count()
    '''
    first = dbset.find_one()
    fid = first['_id']
    print(fid)
    fidint = int(str(fid), 16)
    nid = ObjectId(str(hex(fidint + 100))[2:])
    print(nid)
    print(dbset.find_one({"_id": nid}))
    '''

    cid = dbset.find().skip(idx - 1).limit(1)[0]['_id']
    while idx < count:
        item = db.tc.find({"_id": {"$gte": cid}}).limit(interval)
        print(str(idx - 1).rjust(8, '0') + ' => ' + item[0]['idhash'] + ' : ' + item[0]['caddr'] + "/" + newaddr, end=' # ')
        print(dbset.update_one({"_id": item[0]['_id']}, {"$set": {"caddr": newaddr}}).modified_count)
        cid = item[interval - 1]['_id']
        idx += interval
        '''
        item = dbset.find().skip(idx - 1).limit(1)
        cid = item[0]['_id']
        print(str(idx - 1).rjust(8, '0') + ' => ' + item[0]['idhash'] + ' : ' + item[0]['caddr'] + "/" + newaddr, end=' # ')
        print(dbset.update_one({"_id": cid}, {"$set": {"caddr": newaddr}}).modified_count)
        idx += interval
        '''


# the limit keyword in mysql and the skip keyword in mongo is not efficient for batch updates

# updateMySql(start=1, interval=1000, newaddr='0x' + 'a' * 40)

# updateMongo(start=1, interval=1000, newaddr='0x' + 'a' * 40)
