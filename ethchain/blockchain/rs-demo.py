import os, sys, time, json, web3, timeit, asyncio, uuid, random, hashlib, string, sqlite3

from web3 import Web3
from threading import Thread
from solc import compile_source
from pymongo import MongoClient
from web3.contract import ConciseContract
from web3.middleware import geth_poa_middleware

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
    if 'http' in provider_url:
        w3 = Web3(Web3.HTTPProvider(provider_url))
    else:
        w3 = Web3(Web3.WebsocketProvider(provider_url))
    w3.middleware_stack.inject(geth_poa_middleware, layer=0)
    w3.eth.defaultAccount = w3.eth.accounts[0]
    if passphrase is not None and passphrase != '':
        w3.personal.unlockAccount(w3.eth.defaultAccount, passphrase)
    return w3


# 从文件获取合约信息(address, abi)
def getContract(filepath):
    cinfo = None
    with open(filepath, 'r') as rf:
        cinfo = json.loads(rf.readline())
    return cinfo


# 解锁默认账号
def unlockAccount(addr, passphrase):
    # set pre-funded account as sender
    w3.personal.unlockAccount(addr, passphrase)


# 根据keystore文件和密码获取私钥, 公钥和地址, 并可以通过私钥和新密码生成新的keystore文件和地址. 区块链会自动感知
# geth提供geth --datadir <somepath> account update <address>的方式更新密码, 但不能在控制台这样操作
def updatePasswd():
    prikey = ""
    oldpasswd = "test"
    newpasswd = "node"
    with open('/tmp/keystore/UTC--2018-07-24T08-45-18.136630972Z--dcf9c96b3f9aa9818a1a2543bad39ae4b9fb9616') as keyfile:
        encrypted_key = keyfile.read()
        private_key = w3.eth.account.decrypt(encrypted_key, oldpasswd)
        prikey = private_key
        print(w3.toHex(private_key))
        acct = w3.eth.account.privateKeyToAccount(private_key)
        print(acct.address)
        print(acct._key_obj.public_key)
    # TODO 删除keystore文件
    print(w3.personal.importRawKey(prikey, newpasswd))  # 地址不会改变


# 根据合约和交易hash获取交易日志的详细信息
def parseReceipt(txhash):
    txreceipt = w3.eth.getTransactionReceipt(txhash)
    receipt = contract.events.TXReceipt().processReceipt(txreceipt)
    print(receipt)
    return receipt


# 部署合约
def deploy(contract_source_code, w3):
    # 编译合约源码
    compiledSol = compile_source(contract_source_code)
    contractInterface = compiledSol['<stdin>:RS']

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


def getVersion():
    return w3.version


def getAccounts():
    return w3.eth.accounts


def getAddr(uid):
    uaddr, uid, utype, uname, uhash, ulevel, utimestamp, ustatus = contract.functions.getUserById(uid).call()
    print(w3.toHex(uaddr))
    return str(w3.toHex(uaddr))


def getUser(idx):
    uaddr, uid, utype, uname, uhash, ulevel, utimestamp, ustatus = contract.functions.getUserByIdx(idx).call()
    # print(uid)
    print(w3.toHex(uid))
    return str(w3.toHex(uid))


def getItem(idx):
    ids, itype, title, ihash, path, level, timestamp, ustatus = contract.functions.getItemByIdx(idx).call()
    # print(uid)
    # print(w3.toHex(uid))
    print(w3.toHex(ids[0]), w3.toHex(ids[1]), w3.toText(title))
    return str(w3.toHex(ids[0]))


def getPerm(idx):
    ids, ptype, ptimestamp, ptime, ptimes, status = contract.functions.getPermByIdx(idx).call()
    print(w3.toHex(ids[0]), w3.toHex(ids[1]), w3.toHex(ids[2]), w3.toHex(ids[3]))
    return str(w3.toHex(ids[0]))


def queryEvent():
    # b'^a\x82L\x9a\xcbg\xfd\xa9\xc5\x00\x10\x184\xea\n\xf2x{\xa6\x06\xca3\x01\xce\xfd\xd2}\xf4a\x17\xff'
    # str(w3.sha3(text="uMYHF94C6"))
    # myfilter = contract.eventFilter('TXReceipt', {'fromBlock': 175, 'toBlock': 'latest', 'filter': {'payee': getAccounts()[1]}})
    myfilter = contract.eventFilter(
        'TXReceipt',
        {
            'fromBlock': 0,
            'toBlock': 'latest',
            'filter': {
                'sn': b'3cd2f5ddf47f4267998274c8fdb597fa'
                #    'sn': "TJah0LveE".encode('utf8')
            }
        })  # sn must be bytes32
    eventlist = myfilter.get_all_entries()
    print(eventlist)


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


# addUser
def addUser(idx):
    addr = w3.eth.accounts[idx]
    uid = "0x" + str(uuid.uuid1().hex)
    itype = "user".encode('utf8')
    name = randata([], "str", 6).encode('utf8')
    ihash = "0x" + str(uuid.uuid1().hex)
    level = 2
    status = 1

    sn = "0x" + str(uuid.uuid1().hex)
    details = "addUser"

    params = {'from': w3.eth.coinbase, 'gas': w3.eth.getBlock('latest').gasLimit, 'value': 5000000000000000000000}
    txreceipthash = contract.functions.addUser(sn, details, addr, uid, itype, name, ihash, level, status).transact(params)
    txreceipt = w3.eth.waitForTransactionReceipt(txreceipthash)
    print(txreceipt)
    return uid


# addItem
def addItem():
    userid = getUser(1)
    itemid = "0x" + str(uuid.uuid1().hex)
    ids = [itemid, userid]
    itype = "video".encode('utf8')
    title = "sample video".encode('utf8')
    ihash = "0x" + str(uuid.uuid1().hex)
    path = "/root/sample.mp4".encode('utf8')
    level = 2
    status = 1

    sn = "0x" + str(uuid.uuid1().hex)
    details = "addItem"

    params = {'from': w3.eth.coinbase, 'gas': w3.eth.getBlock('latest').gasLimit, 'value': 5000000000000000000000}
    txreceipthash = contract.functions.addItem(sn, details, ids, itype, title, ihash, path, level, status).transact(params)
    txreceipt = w3.eth.waitForTransactionReceipt(txreceipthash)
    print(txreceipt)
    return itemid


# addPerm
def addPerm(userid, itemid, addr=None, passphrase="node1"):
    permid = "0x" + str(uuid.uuid1().hex)
    # userid = getUser(2)
    # itemid = getItem(0)
    ids = [permid, userid, itemid]
    device = 'windows'.encode('utf8')
    ptype = [0, 1, 1]
    ptimestamp = [1531357394956, 1531789394956]
    ptime = [10, 10]
    ptimes = [22, 22]
    status = 1

    sn = "0x" + str(uuid.uuid1().hex)
    details = "addPerm"
    if addr is None:
        addr = w3.eth.coinbase
    unlockAccount(addr, passphrase)
    params = {'from': addr, 'gas': w3.eth.getBlock('latest').gasLimit, 'value': 5000000000000000000000}
    txreceipthash = contract.functions.addPerm(sn, details, ids, device, ptype, ptimestamp, ptime, ptimes,
                                               status).transact(params)
    txreceipt = w3.eth.waitForTransactionReceipt(txreceipthash)
    print(txreceipt)
    return permid


'''
./geth --identity "myeth" --networkid 111 --nodiscover --maxpeers 10 --port "30303" --dev --dev.period=0 --targetgaslimit 80000000  --rpc --rpcapi "db,eth,net,web3,personal,miner,admin,debug" --rpcaddr 0.0.0.0 --rpcport "8545" --ws --wsaddr 0.0.0.0 --wsorigins=* --wsapi "db,eth,net,web3,personal,miner,admin,debug" --wsport "8546" --datadir "/tmp" console
'''

# rmongo = MyMongoDB(setname="record")
# tmongo = MyMongoDB(setname="token")

w3 = myProvider("ws://172.16.201.191:8546", "node1")
# w3 = myProvider("ws://127.0.0.1:8546", "test")
params = {'from': w3.eth.coinbase, 'gas': w3.eth.getBlock('latest').gasLimit, 'value': 5000000000000000000000}
if os.path.exists('./cinfo.json'):
    contractInfo = getContract('./cinfo.json')
else:
    estimate = w3.eth.getBlock('latest').gasLimit
    while True:
        with open('./RS.sol') as rf:
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
'''
addUser(1)
addUser(2)
addUser(3)
addItem()
addPerm(getUser(2), getItem(0))

addPerm()
getUser(1)
getItem(0)
getPerm(0)
'''
'''
for i in range(contract.functions.getUserNum().call()):
    print(contract.functions.getUserByIdx(i).call())

for i in range(contract.functions.getItemNum().call()):
    print(contract.functions.getItemByIdx(i).call())

for i in range(contract.functions.getPermNum().call()):
    print(contract.functions.getPermByIdx(i).call())

for i in range(contract.functions.getLogNum().call()):
    print(contract.functions.getLogByIdx(i).call())
'''

# print(contract.functions.fetchPerm("cp".encode(), "cd9f0c87e42d434ca61ecb2486445c89".encode()).call())
'''
print(contract.functions.verPerm(
    "cp".encode(),
    "9358742f03754a40b9e5d273fb4e3a2f".encode(),
    2,
    153180317,
    10,
    10,
).call())
'''
# addUser(1)
# queryEvent()

# print(contract.functions.getItemByIdx(0).call())

# print(contract.functions.test(getUser(1), getItem(0)).call())

# print(contract.functions.testa(getUser(1), getItem(0)).call())

#_, a, _, _, _, _, _ = contract.functions.getPermById("a".encode()).call()
#print(w3.toText(a))

for i in range(contract.functions.getUserNum().call()):
    print(contract.functions.getUserByIdx(i).call())