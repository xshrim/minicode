import os, time, json, web3, timeit, asyncio, uuid, random, hashlib, string, sqlite3

from web3 import Web3
from threading import Thread
from solc import compile_source
from pymongo import MongoClient
from web3.contract import ConciseContract
from web3.middleware import geth_poa_middleware


class Token:
    idhash = "0x" + "0" * 32
    parenthash = "0x" + "0" * 32
    origincode = "0x" + "0" * 32
    statuscode = "0x" + "0" * 32
    typecode = "0x" + "0" * 32
    itemcode = "0x" + "0" * 32
    unitcode = "yuan"
    chancode = "0x" + "0" * 32
    keycode = "0x" + "0" * 32
    gaddr = "0x" + "0" * 40
    saddr = "0x" + "0" * 40
    caddr = "0x" + "0" * 40
    paddr = "0x" + "0" * 40
    snetwork = 0
    cnetwork = 0
    amount = 0
    gtimestamp = 0
    ftimestamp = 0
    etimestamp = 0
    ctimestamp = 0
    stimestamp = 0
    status = 0
    reserve = 0

    def __init__(self, idhash, parenthash, origincode, statuscode, typecode, itemcode, unitcode, chancode, keycode, gaddr, saddr,
                 caddr, paddr, snetwork, cnetwork, amount, gtimestamp, ftimestamp, etimestamp, ctimestamp, stimestamp, status,
                 reserve):
        self.idhash = idhash
        self.parenthash = parenthash
        self.origincode = origincode
        self.statuscode = statuscode
        self.typecode = typecode
        self.itemcode = itemcode
        self.unitcode = unitcode
        self.chancode = chancode
        self.keycode = keycode
        self.gaddr = gaddr
        self.saddr = saddr
        self.caddr = caddr
        self.paddr = paddr
        self.snetwork = snetwork
        self.cnetwork = cnetwork
        self.amount = amount
        self.gtimestamp = gtimestamp
        self.ftimestamp = ftimestamp
        self.etimestamp = etimestamp
        self.ctimestamp = ctimestamp
        self.stimestamp = stimestamp
        self.status = status
        self.reserve = reserve

    def __eq__(self, other):
        return self.idhash == other.idhash and self.parenthash == other.parenthash and self.origincode == other.origincode and self.statuscode == other.statuscode and self.typecode == other.typecode and self.itemcode == other.itemcode and self.unitcode == other.unitcode and self.chancode == other.chancode and self.keycode == other.keycode and self.gaddr == other.gaddr and self.saddr == other.saddr and self.caddr == other.caddr and self.paddr == other.paddr and self.snetwork == other.cnetwork and self.amount == other.amount and self.gtimestamp == other.gtimestamp and self.ftimestamp == other.ftimestamp and self.etimestamp == other.etimestamp and self.ctimestamp == other.ctimestamp and self.stimestamp == other.stimestamp and self.status == other.status and self.reserve == other.reserve

    def __str__(self):
        return str({
            "idhash": self.idhash,
            "parenthash": self.parenthash,
            "origincode": self.origincode,
            "statuscode": self.statuscode,
            "typecode": self.typecode,
            "itemcode": self.itemcode,
            "unitcode": self.unitcode,
            "chancode": self.chancode,
            "keycode": self.keycode,
            "gaddr": self.gaddr,
            "saddr": self.saddr,
            "caddr": self.caddr,
            "paddr": self.paddr,
            "snetwork": self.snetwork,
            "cnetwork": self.cnetwork,
            "amount": self.amount,
            "gtimestamp": self.gtimestamp,
            "ftimestamp": self.ftimestamp,
            "etimestamp": self.etimestamp,
            "ctimestamp": self.ctimestamp,
            "stimestamp": self.stimestamp,
            "status": self.status,
            "reserve": self.reserve
        })

    def show(self):
        return {
            "idhash": self.idhash,
            "parenthash": self.parenthash,
            "origincode": self.origincode,
            "statuscode": self.statuscode,
            "typecode": self.typecode,
            "itemcode": self.itemcode,
            "unitcode": self.unitcode,
            "chancode": self.chancode,
            "keycode": self.keycode,
            "gaddr": self.gaddr,
            "saddr": self.saddr,
            "caddr": self.caddr,
            "paddr": self.paddr,
            "snetwork": self.snetwork,
            "cnetwork": self.cnetwork,
            "amount": self.amount,
            "gtimestamp": self.gtimestamp,
            "ftimestamp": self.ftimestamp,
            "etimestamp": self.etimestamp,
            "ctimestamp": self.ctimestamp,
            "stimestamp": self.stimestamp,
            "status": self.status,
            "reserve": self.reserve
        }


class Record:
    sn = ""
    txtype = ""
    callfunc = ""
    txhash = "0x" + "0" * 32
    ifhash = "0x" + "0" * 32
    channel = "0x" + "0" * 32
    tonetwork = 0
    timestamp = 0
    sender = "0x" + "0" * 40
    payer = "0x" + "0" * 40
    payee = "0x" + "0" * 40
    splitids = []
    payids = []
    backids = []
    infos = ""

    def __init__(self, sn, txtype, callfunc, txhash, ifhash, channel, tonetwork, timestamp, sender, payer, payee, splitids,
                 payids, backids, infos):
        self.sn = sn
        self.txtype = txtype
        self.callfunc = callfunc
        self.txhash = txhash
        self.ifhash = ifhash
        self.channel = channel
        self.tonetwork = tonetwork
        self.timestamp = timestamp
        self.sender = sender
        self.payer = payer
        self.payee = payee
        self.splitids = splitids
        self.payids = payids
        self.backids = backids
        self.infos = infos

    def __str__(self):
        return str({
            "sn": self.sn,
            "txtype": self.txtype,
            "callfunc": self.callfunc,
            "txhash": self.txhash,
            "ifhash": self.ifhash,
            "channel": self.channel,
            "tonetwork": self.tonetwork,
            "timestamp": self.timestamp,
            "sender": self.sender,
            "payer": self.payer,
            "payee": self.payee,
            "splitids": self.splitids,
            "payids": self.payids,
            "backids": self.backids,
            "infos": self.infos
        })

    def show(self):
        return {
            "sn": self.sn,
            "txtype": self.txtype,
            "callfunc": self.callfunc,
            "txhash": self.txhash,
            "ifhash": self.ifhash,
            "channel": self.channel,
            "tonetwork": self.tonetwork,
            "timestamp": self.timestamp,
            "sender": self.sender,
            "payer": self.payer,
            "payee": self.payee,
            "splitids": self.splitids,
            "payids": self.payids,
            "backids": self.backids,
            "infos": self.infos
        }


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

    def query(self, dic):
        self.display("find...")
        return self.dbset.find(dic)


# 智能合约代码
contractSourceCode = '''
pragma solidity ^0.4.21;

contract Greeter {
    uint public count;
    string public greeting;

    function Greeter() public {
        count = 0;
        greeting = 'Hello';
    }

    function setGreeting(string _greeting) public {
        greeting = _greeting;
        count++;
    }

    function greet() view public returns (string) {
        return greeting;
    }

    function getCount() view public returns (uint) {
        return count;
    }
}
'''


def myProvider(provider_url, passphrase="test"):
    # web3.py instance
    # w3 = Web3(Web3.EthereumTesterProvider())
    # PoA共识机制下api需要注入PoA中间件
    w3 = Web3(Web3.WebsocketProvider(provider_url))
    w3.middleware_stack.inject(geth_poa_middleware, layer=0)
    if passphrase is not None and passphrase != '':
        unlockAccount(w3, passphrase)
    return w3


# 从文件获取合约信息(address, abi)
def getContract(filepath):
    cinfo = None
    with open(filepath, 'r') as rf:
        cinfo = json.loads(rf.readline())
    return cinfo


# 解锁默认账号
def unlockAccount(w3, passphrase):
    # set pre-funded account as sender
    w3.eth.defaultAccount = w3.eth.accounts[0]
    w3.personal.unlockAccount(w3.eth.accounts[0], passphrase)


# 部署合约
def deploy(contract_source_code, w3):
    # 编译合约源码
    compiledSol = compile_source(contract_source_code)
    contractInterface = compiledSol['<stdin>:Exchange']

    # 生成合约对象
    Contract = w3.eth.contract(abi=contractInterface['abi'], bytecode=contractInterface['bin'])

    # with open('./test.json', 'w') as wf:
    #    wf.write(contractInterface['bin'])
    # print(contractInterface['bin'])
    # print(Contract.constructor().estimateGas(params))
    # 部署合约
    txhash = Contract.constructor().transact(params)

    # 等待合约部署交易完成
    txreceipt = w3.eth.waitForTransactionReceipt(txhash)
    print(txreceipt)
    if txreceipt['status'] != 1:
        return None

    # 将合约地址和abi进行json化保存到本地文件
    contractInfo = {'address': txreceipt.contractAddress, 'abi': contractInterface['abi']}
    with open('./cinfo.json', 'w') as wf:
        wf.write(json.dumps(contractInfo))

    return contractInfo


def handle_event(event):
    print(event)
    # and whatever


def log_loop(event_filter, poll_interval):
    while True:
        for event in event_filter.get_new_entries():
            handle_event(event)
        time.sleep(poll_interval)


def getDBConnect(dbfile):
    conn = sqlite3.connect(dbfile)
    conn.row_factory = dict_factory
    if os.path.exists(dbfile) and os.path.isfile(dbfile):
        return conn
    else:
        return None


def dbCreateTable(conn, table):
    cursor = conn.cursor()
    if table == 'token':
        cursor.execute('''
        CREATE TABLE IF NOT EXISTS TOKEN(
            IDHASH TEXT PRIMARY KEY NOT NULL,
            PARENTHASH  TEXT    NOT     NULL,
            ORIGINCODE  TEXT    NOT     NULL,
            STATUSCODE  TEXT    NOT     NULL,
            TYPECODE    TEXT    NOT     NULL,
            ITEMCODE    TEXT    NOT     NULL,
            UNITCODE    TEXT    NOT     NULL,
            SADDR   TEXT    NOT     NULL,
            CADDR   TEXT    NOT     NULL,
            snetwork  INT     NOT     NULL,
            cnetwork  INT     NOT     NULL,
            AMOUNT  INT     NOT     NULL,
            GTIMESTAMP  INT     NOT     NULL,
            CTIMESTAMP  INT     NOT     NULL,
            GTIMESTAMP  INT     NOT     NULL,
            STATUS  INT     NOT     NULL
        );
        ''')
    conn.commit()
    cursor.close()
    if len(dbQuery(conn, "SELECT count(*) FROM sqlite_master WHERE type='table' AND name='TOKEN'")) > 0:
        return True
    else:
        return False


def dbQuery(conn, sql):
    rows = []
    cursor = conn.cursor()
    csr = cursor.execute(sql)
    for row in csr:
        rows.append(row)
    cursor.close()
    return rows


def dbInsert(conn, sqls):
    cursor = conn.cursor()
    for sql in sqls:
        if sql is not None and sql != '':
            cursor.execute(sql)
    cursor.commit()


def getVersion():
    return w3.version


def getAccounts():
    return w3.eth.accounts


def getTokenSum(addr):
    return contract.functions.getTokenSum(addr).call()


def randata(data, mode="default", num=10):
    if mode == "str":
        data.append(''.join(random.sample(string.ascii_letters + string.digits, num)))
    if mode == "bytes":
        data.append(''.join(random.sample(string.ascii_letters + string.digits, num)).encode('utf8'))
    elif mode == "addr":
        data.append("0x" + str(hashlib.sha1(uuid.uuid1().bytes).hexdigest()))
    elif mode == "hash":
        # data.append("0x" + str(hashlib.sha256(uuid.uuid1().bytes).hexdigest()))
        data.append("0x" + str(uuid.uuid1().hex))
    elif mode == "small":
        data.append(random.randint(100, 999))
    elif mode == "large":
        data.append(random.randint(0, 9999999999))
    elif mode == "timestamp":
        data.append(int(time.time()) + random.randint(-1000000, 1000000))
    else:
        random.shuffle(data)
    return random.choice(data)
    # return random.sample(data, 1)[0]


def getTokenHolder(arg, mode=0):
    holder = None
    try:
        if mode == 0:
            holder = contract.functions.getTokenHolderByIdx(arg).call()
        else:
            holder = contract.functions.getTokenHolderById(arg).call()
    except Exception as ex:
        pass
    return holder


def fetchToken(arg, src="bc", mode="id"):
    tklist = []
    try:
        if src == "bc":
            if mode == "id":
                hashs, addrs, ints = contract.functions.getTokenById(arg).call()
            elif mode == "idx":
                hashs, addrs, ints = contract.functions.getTokenByIdx(arg).call()
            else:
                return token
            idhash = str(w3.toHex(hashs[0]))
            if int(idhash, 16) != 0:
                parenthash = str(w3.toHex(hashs[1]))
                origincode = str(w3.toHex(hashs[2]))
                statuscode = str(w3.toHex(hashs[3]))
                typecode = str(w3.toHex(hashs[4]))
                itemcode = str(w3.toHex(hashs[5]))
                # unitcode = str(w3.toText(hashs[6]))
                unitcode = str(w3.toHex(hashs[6]))
                chancode = str(w3.toHex(hashs[7]))
                keycode = str(w3.toHex(hashs[8]))
                gaddr = str(addrs[0])
                saddr = str(addrs[1])  # str(w3.toHex(item[2]))
                caddr = str(addrs[2])  # str(w3.toHex(item[3]))
                paddr = str(addrs[3])
                snetwork = int(ints[0])
                cnetwork = int(ints[1])
                amount = int(ints[2])
                gtimestamp = int(ints[3])
                ftimestamp = int(ints[4])
                etimestamp = int(ints[5])
                ctimestamp = int(ints[6])
                stimestamp = int(ints[7])
                status = int(ints[8])
                reserve = int(ints[9])
                tklist.append(
                    Token(idhash, parenthash, origincode, statuscode, typecode, itemcode, unitcode, chancode, keycode, gaddr,
                          saddr, caddr, paddr, snetwork, cnetwork, amount, gtimestamp, ftimestamp, etimestamp, ctimestamp,
                          stimestamp, status, reserve))
        else:
            if mode == "id" or mode == "idhash":
                dic = {"idhash": arg}
            elif mode == "caddr":
                dic = {"caddr": arg}
            else:
                return tklist
            for doc in tmongo.query(dic):
                tklist.append(doc)
    except Exception as ex:
        pass
    return tklist


def getReceipt(sn):
    txhash = contract.functions.getSTXByNum(sn).call()
    txreceipt = w3.eth.getTransactionReceipt(txhash)
    # receipt = contract.events.PayReceipt().processReceipt(tx_receipt)[0]['args']
    receipt = contract.events.PayReceipt().processReceipt(txreceipt)[0]
    print(receipt)
    return receipt


def getReceipts(payer="0x" + "0" * 40, payee="0x" + "0" * 40, pay=None):
    receipts = []
    sns = contract.functions.getSnSum().call()
    for i in range(0, sns):
        txhash = contract.functions.getSTXByIdx(i).call()
        txreceipt = w3.eth.getTransactionReceipt(txhash)
        receipts.append(contract.events.PayReceipt().processReceipt(txreceipt)[0])
    if pay is None:
        if payer != 0x00 and payer != "0x" + "0" * 40:
            receipts = [receipt for receipt in receipts if receipt['args']['payer'] == str(payer)]
        if payee != 0x00 and payee != "0x" + "0" * 40:
            receipts = [receipt for receipt in receipts if receipt['args']['payee'] == str(payee)]
    else:
        receipts = [
            receipt for receipt in receipts if receipt['args']['payer'] == str(pay) or receipt['args']['payee'] == str(pay)
        ]
    return receipts


def setRTX(sn, txhashlist):
    if sn == "0":
        # sn = "0x" + str(hashlib.sha256(uuid.uuid1().bytes).hexdigest())
        sn = "0x" + str(uuid.uuid1().hex)
    txhashs = ','.join(txhashlist)
    txreceipthash = contract.functions.setSTX(sn, txhashs).transact(params)
    txreceipt = w3.eth.waitForTransactionReceipt(txreceipthash)
    return True
    # if txreceipt['status'] == 1:
    #    return True
    # else:
    #    return False


def recordDB(record):
    print(rmongo.insert([record.show()]).inserted_ids[0])


def tokenDB(token):
    print(tmongo.insert([token.show()]).inserted_ids[0])


def fetchAddress(addr):
    custaddr, custno, custhash, custinfo, custlevel, custstatus = contract.functions.getCustByAddr(addr).call()
    return (custaddr, custno, custhash, custinfo, custlevel, custstatus)


def enrollAddress(custaddr, custno, custhash, custinfo, custlevel, custstatus):
    record = createRecord(txtype="registeCustomer", callfunc="setCustAddr", payee=custaddr)
    '''
    details = {
        "sn": record.sn,
        "txtype": record.txtype,
        "callfunc": record.callfunc,
        "txhash": record.txhash,
        "ifhash": record.ifhash,
        "channel": record.channel,
        "tonetwork": record.tonetwork,
        "timestamp": record.timestamp,
        "sender": record.sender,
        "payer": record.payer,
        "payee": record.payee,
        "splitids": record.splitids,
        "payids": record.payids,
        "backids": record.backids,
        "infos": record.infos
    }
    '''
    txhash = contract.functions.setCustAddr(
        str(record.sn).encode('utf8'), str(record), str(custaddr), str(custno), str(custhash), str(custinfo), int(custlevel),
        int(custstatus), int(record.timestamp)).transact(params)
    txreceipt = w3.eth.waitForTransactionReceipt(txhash)
    if txreceipt is not None and txreceipt['status'] == 1:
        # details["txhash"] = str(w3.toHex(txhash))
        record.txhash = str(w3.toHex(txhash))
        recordDB(record)
        print(txreceipt)


# 发起交易
def transit(contract, params):
    contract.functions.setGreeting('Nihao').transact(params)


def createToken(idhash=None,
                parenthash=None,
                origincode=None,
                statuscode=None,
                typecode=None,
                itemcode=None,
                unitcode=None,
                chancode=None,
                keycode=None,
                gaddr=None,
                saddr=None,
                caddr=None,
                paddr=None,
                snetwork=None,
                cnetwork=None,
                amount=None,
                gtimestamp=None,
                ftimestamp=None,
                etimestamp=None,
                ctimestamp=None,
                stimestamp=None,
                status=None,
                reserve=None):
    try:
        if idhash is None:
            idhash = randata([], "hash")
        if parenthash is None:
            parenthash = "0x00" + "0" * 30
        if origincode is None:
            origincode = randata(["0x00" + "0" * 30, "0x01" + "0" * 30, "0x10" + "0" * 30, "0x11" + "0" * 30], "hash")
        if statuscode is None:
            statuscode = randata(["0x00" + "0" * 30, "0x01" + "0" * 30, "0x10" + "0" * 30, "0x11" + "0" * 30], "hash")
        if typecode is None:
            # typecode = str(hashlib.sha256(randata(metadata[4:6], "str", 6).encode('utf8')).hexdigest())
            # typecode = randata(metadata[4:6], "str", 6).encode('utf8')
            typecode = randata(["0x00" + "0" * 30, "0x01" + "0" * 30, "0xff" + randata([], "hash")[4:]], "None")
        if itemcode is None:
            itemcode = randata([], "hash")
        if unitcode is None:
            # unitcode = str(hashlib.sha256(randata(metadata[6:], "str", 6).encode('utf8')).hexdigest())
            # unitcode = randata(["yuan", "ge", "jian", "xiang"], "str", 6)
            unitcode = randata(["0x" + "0" * 30 + "01", "0x" + "0" * 30 + "02", "0x" + "0" * 30 + "03", "0x" + "0" * 30 + "04"],
                               "hash")
        if chancode is None:
            chancode = randata(["0x" + "0" * 30 + "01", "0x" + "0" * 30 + "02", "0x" + "0" * 30 + "03", "0x" + "0" * 30 + "04"],
                               "None")
        if keycode is None:
            keycode = "0x00" + "0" * 30
        if gaddr is None:
            gaddr = getAccounts()[0]
        if saddr is None:
            saddr = randata(getAccounts()[1:], "None")
        if caddr is None:
            caddr = saddr
        if paddr is None:
            paddr = "0x" + "0" * 40
        if snetwork is None:
            snetwork = int(getVersion().network)
        if cnetwork is None:
            cnetwork = snetwork
        if amount is None:
            amount = randata([], "large")
        if gtimestamp is None:
            gtimestamp = randata([], "timestamp")
        if ftimestamp is None:
            ftimestamp = gtimestamp
        if etimestamp is None:
            etimestamp = gtimestamp + 50000000
        if ctimestamp is None:
            ctimestamp = gtimestamp
        if stimestamp is None:
            stimestamp = gtimestamp
        if status is None:
            status = 0
        if reserve is None:
            reserve = 0
        # idhash = "0x" + str(uuid.uuid1().bytes.hexdigest())  # sha1
        # idhash = "0x" + str(uuid.uuid1().hex)  # sha1
        # idhash = Web3.sha3(uuid.uuid1().bytes)
        # idhash = "0x" + str(uuid.uuid1()).replace("-", "") + "00000000000000000000000000000000"
        return Token(
            str(idhash), str(parenthash), str(origincode), str(statuscode), str(typecode), str(itemcode), str(unitcode),
            str(chancode), str(keycode), str(gaddr), str(saddr), str(caddr), str(paddr), int(snetwork), int(cnetwork),
            int(amount), int(gtimestamp), int(ftimestamp), int(etimestamp), int(ctimestamp), int(stimestamp), int(status),
            int(reserve))
    except Exception as ex:
        return None


def createRecord(sn=None,
                 txtype=None,
                 callfunc=None,
                 txhash=None,
                 ifhash=None,
                 channel=None,
                 tonetwork=None,
                 timestamp=None,
                 sender=None,
                 payer=None,
                 payee=None,
                 splitids=None,
                 payids=None,
                 backids=None,
                 infos=None):
    try:
        if sn is None:
            sn = randata([], "str", 9)
        if txtype is None:
            txtype = "assignToken"
        if callfunc is None:
            callfunc = ""
        if txhash is None:
            txhash = "0x00" + "0" * 30
        if ifhash is None:
            ifhash = randata([], "hash")
        if channel is None:
            channel = randata(["0x" + "0" * 30 + "01", "0x" + "0" * 30 + "02", "0x" + "0" * 30 + "03", "0x" + "0" * 30 + "04"],
                              "None")
        if tonetwork is None:
            tonetwork = int(getVersion().network)
        if timestamp is None:
            timestamp = randata([], "timestamp")
        if sender is None:
            sender = getAccounts()[0]
        if payer is None:
            payer = sender
        if payee is None:
            payee = randata(getAccounts()[1:], "None")
        if splitids is None:
            splitids = []
        if payids is None:
            payids = []
        if backids is None:
            backids = []
        if infos is None:
            infos = ifhash
        return Record(sn, txtype, callfunc, txhash, ifhash, channel, tonetwork, timestamp, sender, payer, payee, splitids, payids,
                      backids, infos)
    except Exception as ex:
        return None


def assignTokenToDB(tkid):
    token = fetchToken(tkid, src="bc", mode="id")
    holder = getTokenHolder(tkid, mode=1)
    if token is None or holder is None or int(holder, 16) == 0 or token.caddr != holder:
        return False
    else:
        pass


def assignToken(sn, token, record):
    hashs = [
        str(token.idhash),
        str(token.parenthash),
        str(token.origincode),
        str(token.statuscode),
        str(token.typecode),
        str(token.itemcode),
        # str(token.unitcode).encode('utf8'),
        str(token.unitcode),
        str(token.chancode),
        str(token.keycode)
    ]
    addrs = [str(token.gaddr), str(token.saddr), str(token.caddr), str(token.paddr)]
    ints = [
        int(token.snetwork),
        int(token.cnetwork),
        int(token.amount),
        int(token.gtimestamp),
        int(token.ftimestamp),
        int(token.etimestamp),
        int(token.ctimestamp),
        int(token.stimestamp),
        int(token.status),
        int(token.reserve)
    ]
    txhash = contract.functions.assignToken(str(sn).encode('utf8'), str(record), str(token.caddr), hashs, addrs,
                                            ints).transact(params)
    txreceipt = w3.eth.waitForTransactionReceipt(txhash)
    if txreceipt is not None and txreceipt['status'] == 1:
        print(txreceipt)
        return (True, txhash)
    else:
        return (False, txhash)
    # if txreceipt['status'] == 1:
    # print("WatchDog", contract.events.WatchDog().processReceipt(txreceipt))
    # print("PayReceipt", contract.events.PayReceipt().processReceipt(txreceipt))
    # payreceipt = contract.events.PayReceipt().processReceipt(txreceipt)
    # sn = random.randint(1, 10000000000000000)
    # return (True, txhash)
    # return (True, "success", sn, w3.toHex(txhash))
    # else:
    # return (False, txhash)
    # return (False, "failed", sn, w3.toHex(txhash))


def grantToken(payee="0x" + "0" * 40, tktype="0x00" + "0" * 30):
    token = createToken(saddr=payee, typecode=tktype)
    record = createRecord(
        txtype="assignToken",
        callfunc="assignToken",
        channel=token.chancode,
        timestamp=token.gtimestamp,
        sender=token.gaddr,
        payer=token.gaddr,
        payee=token.saddr,
        payids=[token.idhash])
    '''
    details = {
        "sn": record.sn,
        "txtype": record.txtype,
        "callfunc": record.callfunc,
        "txhash": record.txhash,
        "ifhash": record.ifhash,
        "channel": record.channel,
        "tonetwork": record.tonetwork,
        "timestamp": record.timestamp,
        "sender": record.sender,
        "payer": record.payer,
        "payee": record.payee,
        "splitids": record.splitids,
        "payids": record.payids,
        "backids": record.backids,
        "infos": record.infos
    }
    '''
    # print(token)
    # print(record)
    res, txhash = assignToken(record.sn, token, record)
    if res is True:
        tokenDB(token)
        # details["txhash"] = str(w3.toHex(txhash))
        record.txhash = str(w3.toHex(txhash))
        recordDB(record)

    # print(contract.events.PayReceipt().processReceipt(receipt)[0])


def addCustomerDemo():
    enrollAddress(getAccounts()[1], "00000001", "admina", "admina", 1, 0)
    enrollAddress(getAccounts()[2], "00000002", "customerb", "customerb", 2, 0)
    enrollAddress(getAccounts()[3], "00000003", "customerc", "customerc", 2, 0)
    enrollAddress(getAccounts()[4], "00000004", "customerd", "customerd", 2, 0)
    enrollAddress(getAccounts()[5], "00000005", "customere", "customere", 2, 0)
    enrollAddress(getAccounts()[6], "00000006", "customerf", "customerf", 2, 0)
    enrollAddress(getAccounts()[7], "00000007", "customerg", "customerg", 2, 0)
    enrollAddress(getAccounts()[8], "00000008", "customerh", "customerh", 2, 0)
    enrollAddress(getAccounts()[9], "00000009", "customeri", "customeri", 2, 0)
    enrollAddress(getAccounts()[10], "00000010", "customerj", "customerj", 2, 0)

    enrollAddress("0x98db3E851f67538AE07686D0B45a76BFEcD9b2Ca", "00000011", "adminb", "adminb", 1, 0)
    enrollAddress("0x534E69511694f7BE128597BAb38f5F6f262f2436", "00000012", "customerk", "customerk", 2, 0)

    enrollAddress("0xCEa677127C6b8885eF4fb7616E2875380BD6a97a", "00000013", "adminc", "adminc", 1, 0)
    enrollAddress("0x2B797f302b6e6C19fcFb52a9674dB91C2d999b29", "00000014", "customerl", "customerl", 2, 0)

    # print(len(getAccounts()))
    for addr in getAccounts():
        print(addr + ":" + str(fetchAddress(addr)))


def addTokenDemo():
    for k in range(1, 11):
        for i in range(10):
            grantToken(getAccounts()[k], randata(["0x00" + "0" * 30, "0x01" + "0" * 30, "0xff" + randata([], "hash")[4:]],
                                                 "None"))
    for i in range(getTokenSum("0x" + "0" * 40)[0]):
        print(fetchToken(i, src="bc", mode="idx"))


def queryTokenDemo():
    tklist = fetchToken("0x415635aa805511e8807428e347233abc", src="bc", mode="id")
    # tklist = fetchToken(getAccounts()[2], src="bc", mode="id")
    for token in tklist:
        print(token)


def queryRecordDemo():
    # b'^a\x82L\x9a\xcbg\xfd\xa9\xc5\x00\x10\x184\xea\n\xf2x{\xa6\x06\xca3\x01\xce\xfd\xd2}\xf4a\x17\xff'
    # str(w3.sha3(text="uMYHF94C6"))
    # myfilter = contract.eventFilter('TXReceipt', {'fromBlock': 175, 'toBlock': 'latest', 'filter': {'payee': getAccounts()[1]}})
    myfilter = contract.eventFilter('TXReceipt', {
        'fromBlock': 290,
        'toBlock': 'latest',
        'filter': {
            'sn': "TJah0LveE".encode('utf8')
        }
    })  # sn must be bytes32
    eventlist = myfilter.get_all_entries()
    print(eventlist)


'''
geth --identity=myeth --networkid=111 --nodiscover --maxpeers=10 --port=30303 --dev --dev.period=0 --gasprice=1000 --targetgaslimit=4712388000 --rpc --rpcapi=db,eth,net,web3,personal,miner,admin,debug --rpcaddr=0.0.0.0 --rpccorsdomain=* --rpcport=8545 --ws --wsapi=db,eth,net,web3,personal,miner,admin,debug --wsaddr=0.0.0.0 --wsorigins=* --wsport=8546 --datadir=/tmp console
'''

rmongo = MyMongoDB(setname="record")
tmongo = MyMongoDB(setname="token")

w3 = myProvider("ws://192.168.0.93:8546", "test")
params = {'from': w3.eth.coinbase, 'gas': 30000000}
if os.path.exists('./cinfo.json'):
    contractInfo = getContract('./cinfo.json')
else:
    estimate = w3.eth.getBlock('latest').gasLimit
    while True:
        with open('./Exchange.sol') as rf:
            contractSourceCode = rf.read()
        contractInfo = deploy(contractSourceCode, w3)
        if contractInfo is None:
            print('deploy failed: ' + str(estimate))
            estimate = w3.eth.getBlock('latest').gasLimit
            params = {'from': w3.eth.coinbase, 'gas': estimate}
        else:
            break

contract = w3.eth.contract(
    address=contractInfo['address'],
    abi=contractInfo['abi'],
)

queryTokenDemo()

# b'^a\x82L\x9a\xcbg\xfd\xa9\xc5\x00\x10\x184\xea\n\xf2x{\xa6\x06\xca3\x01\xce\xfd\xd2}\xf4a\x17\xff'
# str(w3.sha3(text="uMYHF94C6"))
# myfilter = contract.eventFilter('TXReceipt', {'fromBlock': 175, 'toBlock': 'latest', 'filter': {'payee': getAccounts()[1]}})
# myfilter = contract.eventFilter('TXReceipt', {'fromBlock': 175, 'toBlock': 'latest', 'filter': {'sn': "TJah0LveE".encode('utf8')}})   # sn must be bytes32
# eventlist = myfilter.get_all_entries()
# print(eventlist)

# hex_str = "0x{:064x}".format(10)

# txreceipthash = contract.functions.setHello("credit".encode('utf8')).transact(params)
# txreceipt = w3.eth.waitForTransactionReceipt(txreceipthash)
# print("credit".encode('utf8'))
# print(w3.toText(contract.functions.getHello().call()))

# block_filter = contract.events.PayReceipt.createFilter(fromBlock='latest')
# Thread(target=log_loop, args=(block_filter, 2), daemon=True).start()

# txhash = contract.functions.setGreeting('Hi').transact(params)
# print(w3.eth.waitForTransactionReceipt(txhash))

# addToken()
# print(getVersion().network)
# print(getAccounts())
# print(getTokenSum())
# tk = getToken(0, 0)
# print(tk.unitcode)
'''
txhash = contract.functions.setCustAddr(getAccounts()[1], "00000001", "custa", "custa", 1).transact(params)
txreceipt = w3.eth.waitForTransactionReceipt(txhash)
print(txreceipt)
txhash = contract.functions.setCustAddr(getAccounts()[2], "00000002", "custb", "custb", 2).transact(params)
txreceipt = w3.eth.waitForTransactionReceipt(txhash)
print(txreceipt)
'''

# txhash = contract.functions.assignToken(addr, str("test"), hashs, addrs, ints).transact(params)
# txreceipt = w3.eth.waitForTransactionReceipt(txhash)
# print(txreceipt)

# for i in range(8):
#    print(getRToken(getAccounts()[i]))
# print(random.choice(getAccounts()))
'''
for k in range(1, 11):
    for i in range(10):
        grantToken(getAccounts()[k], randata(["0x00" + "0" * 30, "0x01" + "0" * 30, "0xff" + randata([], "hash")[4:]], "None"))
for i in range(getTokenSum("0x" + "0" * 40)[0]):
    print(fetchToken(i, "bc", "idx"))
# print(getAccounts()[0])
'''
'''
enrollAddress(getAccounts()[1], "00000001", "admina", "admina", 1, 0)
enrollAddress(getAccounts()[2], "00000002", "customerb", "customerb", 2, 0)
enrollAddress(getAccounts()[3], "00000003", "customerc", "customerc", 2, 0)
enrollAddress(getAccounts()[4], "00000004", "customerd", "customerd", 2, 0)
enrollAddress(getAccounts()[5], "00000005", "customere", "customere", 2, 0)
enrollAddress(getAccounts()[6], "00000006", "customerf", "customerf", 2, 0)
enrollAddress(getAccounts()[7], "00000007", "customerg", "customerg", 2, 0)
enrollAddress(getAccounts()[8], "00000008", "customerh", "customerh", 2, 0)
enrollAddress(getAccounts()[9], "00000009", "customeri", "customeri", 2, 0)
enrollAddress(getAccounts()[10], "00000010", "customerj", "customerj", 2, 0)

enrollAddress("0x98db3E851f67538AE07686D0B45a76BFEcD9b2Ca", "00000011", "adminb", "adminb", 1, 0)
enrollAddress("0x534E69511694f7BE128597BAb38f5F6f262f2436", "00000012", "customerk", "customerk", 2, 0)

enrollAddress("0xCEa677127C6b8885eF4fb7616E2875380BD6a97a", "00000013", "adminc", "adminc", 1, 0)
enrollAddress("0x2B797f302b6e6C19fcFb52a9674dB91C2d999b29", "00000014", "customerl", "customerl", 2, 0)

# print(len(getAccounts()))
for addr in getAccounts():
    print(addr + ":" + str(fetchAddress(addr)))
'''
'''
txhash = contract.functions.logToken(0).transact(params)
txreceipt = w3.eth.waitForTransactionReceipt(txhash)
if txreceipt['status'] == 1:
    print("WatchDog", contract.events.WatchDog().processReceipt(txreceipt))
    # print("PayReceipt", contract.events.PayReceipt().processReceipt(txreceipt))
    # payreceipt = contract.events.PayReceipt().processReceipt(txreceipt)
'''
# input()
