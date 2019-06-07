import sys, time, json, web3, timeit, hashlib, uuid, random

from web3 import Web3
from threading import Thread
from solc import compile_source
from web3.contract import ConciseContract
from web3.middleware import geth_poa_middleware
from flask import Flask, jsonify, request, abort


class Token:
    idx = 0
    idhash = "0x00"
    parenthash = "0x00"
    saddr = "0x00"
    caddr = "0x00"
    amount = 0
    price = 1
    timestamp = 0
    stimestamp = 0
    status = 0

    def __init__(self, idx, idhash, parenthash, saddr, caddr, amount, price, timestamp, stimestamp, status):
        self.idx = idx
        self.idhash = idhash
        self.parenthash = parenthash
        self.saddr = saddr
        self.caddr = caddr
        self.amount = amount
        self.price = price
        self.timestamp = timestamp
        self.stimestamp = stimestamp
        self.status = status


# 智能合约代码
contractSourceCode = ""
with open('./Exchange.sol', 'r') as rf:
    contractSourceCode = rf.read()
#print(contractSourceCode)
#sys.exit(0)
'''
contractSourceCode = 
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
    for i in range(0, len(w3.eth.accounts)):
        w3.personal.unlockAccount(w3.eth.accounts[i], passphrase)


# 部署合约
def deploy(contract_source_code, w3):
    # 编译合约源码
    compiled_sol = compile_source(contract_source_code)
    contract_interface = compiled_sol['<stdin>:Exchange']

    # 生成合约对象
    Contract = w3.eth.contract(abi=contract_interface['abi'], bytecode=contract_interface['bin'])

    # 部署合约
    tx_hash = Contract.constructor().transact()

    # 等待合约部署交易完成
    tx_receipt = w3.eth.waitForTransactionReceipt(tx_hash)
    # print(tx_receipt)

    # 将合约地址和abi进行json化保存到本地文件
    contractInfo = {'address': tx_receipt.contractAddress, 'abi': contract_interface['abi']}
    with open('./cinfo.json', 'w') as wf:
        wf.write(json.dumps(contractInfo))

    return contractInfo


# 发起交易
def transit(contract, params):
    contract.functions.setGreeting('Nihao').transact(params)


# 多线程模拟大量并发交易
def tranthread(w3, contractInfo, count):
    print("send %d transactions, multi-threaded, one thread per tx:\n" % (count))

    threads = []

    contract = w3.eth.contract(
        address=contractInfo['address'],
        abi=contractInfo['abi'],
    )
    for i in range(count):
        t = Thread(target=transit, args=[contract, {'from': w3.eth.coinbase, 'gas': 100000}])
        threads.append(t)
        print(".", end="")
    print("\n%d transaction threads created." % len(threads))

    for t in threads:
        t.start()
        print(".", end="")
        sys.stdout.flush()
    print("\nall threads started.")

    for t in threads:
        t.join()
    print("all threads ended.")


# 压力测试
def benchmark(w3, contractInfo, count, type):
    if type == 'thread':
        tranthread(w3, contractInfo, count)


# 一定间隔内交易量
def txsum(w3, blockNumber, newBlockNumber, txCount, start_time):

    # 统计间隔交易量
    txCount_new = 0
    for bl in range(blockNumber + 1, newBlockNumber + 1):
        txCount_new += w3.eth.getBlockTransactionCount(bl)

    ts_blockNumber = w3.eth.getBlock(blockNumber).timestamp
    ts_newBlockNumber = w3.eth.getBlock(newBlockNumber).timestamp
    duration = ts_newBlockNumber - ts_blockNumber
    if duration == 0:
        duration = 0.5
    tps_current = txCount_new / duration

    txCount += txCount_new
    elapsed = timeit.default_timer() - start_time
    tps = txCount / elapsed
    line = "block %d | new #TX %3d / %4.0f ms = %5.1f TPS_current | total: #TX %4d / %4.1f s = %5.1f TPS_average"
    line = line % (blockNumber, txCount_new, duration * 1000, tps_current, txCount, elapsed, tps)
    print(line)

    return txCount


# 实时监测tps
def monitor(w3, pauseBetweenQueries=0.5):

    blockNumber = w3.eth.blockNumber
    txCount = w3.eth.getBlockTransactionCount(blockNumber)
    start_time = timeit.default_timer()
    print('starting timer, at block ', blockNumber, ' which has ', txCount, ' transactions, start time ', start_time)
    while (True):
        newBlockNumber = w3.eth.blockNumber
        if (blockNumber != newBlockNumber):
            txCount = txsum(w3, blockNumber, newBlockNumber, txCount, start_time)
            blockNumber = newBlockNumber

        time.sleep(pauseBetweenQueries)

    print("end")


def getAddress(w3):
    return w3.eth.accounts


def newCompany(addr, desc):
    hexdesc = hashlib.sha256(desc.encode('utf8')).hexdigest()
    print(addr)
    tx_hash = contract.functions.addCompany(str(addr), "0x" + str(hexdesc)).transact(params)
    tx_receipt = w3.eth.waitForTransactionReceipt(tx_hash)
    print(tx_receipt)


def issueToken(addr, parenthash, saddr, caddr, amount, price):
    idhash = "0x" + str(hashlib.sha256(uuid.uuid1().bytes).hexdigest())
    # idhash = "0x" + str(uuid.uuid1()).replace("-", "") + "00000000000000000000000000000000"
    timestamp = int(time.time())
    #print(idhash, timestamp)
    tx_hash = contract.functions.issueToken(str(addr), idhash, parenthash, str(addr), str(addr), amount, 1,
                                            timestamp).transact(params)
    tx_receipt = w3.eth.waitForTransactionReceipt(tx_hash)
    print(tx_receipt)


def splitToken(payer, payee, token, subamount, timestamp):
    rtidx = -1
    rtidhash = "0x" + str(hashlib.sha256(uuid.uuid1().bytes).hexdigest())
    rtparenthash = token.idhash
    rtsaddr = token.saddr
    rtcaddr = payer
    rtamount = int(subamount / token.price)
    rtprice = token.price
    rtstatus = token.status
    toPayerToken = Token(rtidx, rtidhash, rtparenthash, rtsaddr, rtcaddr, rtamount, rtprice, timestamp, timestamp, rtstatus)

    etidx = -1
    etidhash = "0x" + str(hashlib.sha256(uuid.uuid1().bytes).hexdigest())
    etparenthash = token.idhash
    etsaddr = token.saddr
    etcaddr = payee
    etamount = token.amount - int(subamount / token.price)
    etprice = token.price
    etstatus = token.status
    toPayeeToken = Token(etidx, etidhash, etparenthash, etsaddr, etcaddr, etamount, etprice, timestamp, timestamp, etstatus)

    return (toPayerToken, toPayeeToken)


def transact(payer, payee, amount):
    if getCompanyBalance(payer, 'quick') < amount:
        print("not sufficient funds")
        return (False, "not sufficient funds", -1, "0x00")
    else:
        currentSum = 0
        chooseTokens = []
        tokens = getCompanyTokens(payer, 0)
        tokens = [token for token in tokens if token.status == 0]
        tokens = sorted(tokens, key=lambda x: x.amount * x.price, reverse=True)
        ctokens = [token for token in tokens if token.amount * token.price == amount]
        if len(ctokens) > 0:
            chooseTokens.append(ctokens[0])
            currentSum = amount
        else:
            for idx, token in enumerate(tokens):
                #print(token.idx, token.idhash, token.saddr, token.amount, token.timestamp, token.status)
                currentSum += token.amount * token.price
                if currentSum >= amount:
                    chooseTokens = tokens[:idx + 1]
                    break
        #print("=" * 20)
        totkidxs = []
        splittkidx = -1
        hashs = []
        addrs = []
        aps = []

        details = ""
        detaildict = {"amount": amount}

        fullPayToken = ""
        for token in chooseTokens[:-1]:
            totkidxs.append(token.idx)
            fullPayToken += token.idhash + ","

        timestamp = int(time.time())

        if currentSum == amount:
            totkidxs.append(chooseTokens[-1].idx)
            fullPayToken += chooseTokens[-1].idhash

            splittkidx = -1
            hashs = ["0x00", "0x00", "0x00", "0x00"]
            addrs = [payer, payer, payer, payer]
            aps = [0, 0, 1, 1]
        else:
            fullPayToken = fullPayToken[:-1]
            splittkidx = chooseTokens[-1].idx

            toPayerToken, toPayeeToken = splitToken(payer, payee, chooseTokens[-1], currentSum - amount, timestamp)
            hashs = [toPayerToken.idhash, toPayeeToken.idhash, toPayerToken.parenthash, toPayeeToken.parenthash]
            addrs = [toPayerToken.saddr, toPayeeToken.saddr, toPayerToken.caddr, toPayeeToken.caddr]
            aps = [toPayerToken.amount, toPayeeToken.amount, toPayerToken.price, toPayeeToken.price]

            detaildict["splitToken"] = chooseTokens[-1].idhash
            detaildict["toPayerToken"] = toPayerToken.idhash
            detaildict["toPayeeToken"] = toPayeeToken.idhash

        detaildict["fullPayToken"] = fullPayToken
        details = str(detaildict)

        print(details)
        return transactFull(payer, payee, totkidxs, details, splittkidx, hashs, addrs, aps, timestamp)
        #for token in chooseTokens:
        #    print(token.idx, token.idhash, token.saddr, token.amount, token.timestamp, token.status)


def transactFull(payer, payee, totkidxs, details, splittkidx, hashs, addrs, aps, timestamp):
    tx_hash = contract.functions.transactFull(payer, payee, totkidxs, details, splittkidx, hashs, addrs, aps,
                                              timestamp).transact(params)
    tx_receipt = w3.eth.waitForTransactionReceipt(tx_hash)
    if tx_receipt['status'] == 1:
        print(contract.events.PayReceipt().processReceipt(tx_receipt))
        sn = random.randint(1, 10000000000000000)
        tx_receipt_hash = contract.functions.setSTX(sn, tx_hash).transact(params)
        tx_receipt_receipt = w3.eth.waitForTransactionReceipt(tx_receipt_hash)
        return (True, "success", sn, w3.toHex(tx_hash))
    else:
        return (False, "failed", sn, w3.toHex(tx_hash))

    #print(dict(tx_receipt))
    #print(w3.toText(tx_receipt['logs'][0]['data']))


def getReceipt(sn):
    tx_hash = contract.functions.getSTXByNum(sn).call()
    tx_receipt = w3.eth.getTransactionReceipt(tx_hash)
    #receipt = contract.events.PayReceipt().processReceipt(tx_receipt)[0]['args']
    receipt = contract.events.PayReceipt().processReceipt(tx_receipt)[0]
    print(receipt)
    return receipt


def getReceipts(payer="0x00", payee="0x00", pay=None):
    receipts = []
    sns = contract.functions.getSnNum().call()
    for i in range(0, sns):
        tx_hash = contract.functions.getSTXByIdx(i).call()
        tx_receipt = w3.eth.getTransactionReceipt(tx_hash)
        receipts.append(contract.events.PayReceipt().processReceipt(tx_receipt)[0])
    if pay is None:
        if payer != 0x00 and payer != "0x00":
            receipts = [receipt for receipt in receipts if receipt['args']['payer'] == str(payer)]
        if payee != 0x00 and payee != "0x00":
            receipts = [receipt for receipt in receipts if receipt['args']['payee'] == str(payee)]
    else:
        receipts = [
            receipt for receipt in receipts if receipt['args']['payer'] == str(pay) or receipt['args']['payee'] == str(pay)
        ]
    return receipts


def getCompanyHistory(addr):
    receiptjsons = {}
    receipts = getReceipts(pay=addr)
    for receipt in receipts:
        receiptitem = {}
        receiptitem['payer'] = str(receipt['args']['payer'])
        receiptitem['payee'] = str(receipt['args']['payee'])
        receiptitem['amount'] = str(eval(receipt['args']['detail'])['amount'])
        receiptitem['timestamp'] = str(receipt['args']['timestamp'])
        receiptjsons[str(w3.toHex(receipt['transactionHash']))] = receiptitem
    print(receiptjsons)
    return receiptjsons


def getCompanyNum(status):
    num = contract.functions.getCompanyNum(status).call()
    return num


def getCompanyBalance(addr, mode):
    balance = 0
    if mode == 'quick':
        balance = contract.functions.getCompanyBalance(str(addr)).call()
    else:
        num = getCompanyTokenNum(addr)
        for i in range(0, num):
            token = getCompanyToken(addr, i)
            if int(str(token.idhash), 16) != 0 and token.status == 0:
                balance += token.amount * token.price
    #print(balance)
    return balance


def getCompanyTokenNum(addr):
    num = contract.functions.getCompanyTokenNum(str(addr), -1).call()
    return num


def getCompanyToken(addr, idx):
    item = contract.functions.getCompanyToken(str(addr), idx).call()
    idhash = str(w3.toHex(item[0]))
    parenthash = str(w3.toHex(item[1]))
    saddr = str(item[2])  #str(w3.toHex(item[2]))
    caddr = str(item[3])  #str(w3.toHex(item[3]))
    amount = int(item[4])
    price = int(item[5])
    timestamp = int(item[6])
    stimestamp = int(item[7])
    status = int(item[8])

    token = Token(idx, idhash, parenthash, saddr, caddr, amount, price, timestamp, stimestamp, status)
    return token


def getCompanyTokens(addr, status):
    tokens = []
    num = getCompanyTokenNum(addr)
    for i in range(0, num):
        token = getCompanyToken(addr, i)
        if int(str(token.idhash), 16) != 0:
            if status == 8 or token.status == status:
                tokens.append(token)
    return tokens


def addCompany():
    newCompany(accounts[0], "OK")
    newCompany(accounts[1], "HI")
    newCompany(accounts[2], "GD")
    newCompany(accounts[3], "SC")


def injectData(addr):
    #newCompany(addr, "OK")
    issueToken(addr, "0x00", addr, addr, 1000, 1)
    issueToken(addr, "0x00", addr, addr, 500, 1)
    issueToken(addr, "0x00", addr, addr, 1200, 1)
    issueToken(addr, "0x00", addr, addr, 700, 1)
    issueToken(addr, "0x00", addr, addr, 100, 1)
    issueToken(addr, "0x00", addr, addr, 1400, 1)
    issueToken(addr, "0x00", addr, addr, 300, 1)


def getCompanys():
    comps = {}
    for i in range(0, getCompanyNum(0)):
        addr, _, _, _ = contract.functions.getCompanyByIdx(i).call()
        comps[addr] = getCompanyBalance(addr, 'quick')
    return comps


def getCompanyJsonTokens(addr, status):
    jsontokens = {}
    tokens = getCompanyTokens(addr, status)
    for token in tokens:
        tokenitem = {}
        tokenitem['idhash'] = token.idhash
        tokenitem['parenthash'] = token.parenthash
        tokenitem['saddr'] = token.saddr
        tokenitem['caddr'] = token.caddr
        tokenitem['amount'] = token.amount
        tokenitem['price'] = token.price
        tokenitem['timestamp'] = token.timestamp
        tokenitem['stimestamp'] = token.stimestamp
        tokenitem['status'] = token.status
        jsontokens[token.idhash] = tokenitem
    return jsontokens


def setCompanyAmount(addr, amount):
    totkidxs = [0, 0]
    splittkidx = -1
    hashs = ["0x00", "0x" + str(hashlib.sha256(uuid.uuid1().bytes).hexdigest()), "0x00", "0x00"]
    addrs = [addr, addr, addr, addr]
    aps = [0, int(amount), 1, 1]
    timestamp = int(time.time())

    details = ""
    detaildict = {"amount": amount}
    details = detaildict["toPayeeToken"] = hashs[1]
    details = str(detaildict)
    transactFull(accounts[0], addr, totkidxs, details, splittkidx, hashs, addrs, aps, timestamp)
    return getCompanyBalance(addr, 'quick')


w3 = myProvider("ws://127.0.0.1:8546", "test")
#contractInfo = deploy(contractSourceCode, w3)
contractInfo = getContract('./cinfo.json')
params = {'from': w3.eth.coinbase, 'gas': 4700000}  # gas: 21000 - 4712388
contract = w3.eth.contract(
    address=contractInfo['address'],
    abi=contractInfo['abi'],
)
accounts = getAddress(w3)
#addCompany()
#getReceipt(6258427331344578)
#getCompanyHistory('0x65195Ab9867AeCb07Ac69fa6DBFB72eACadCa9Fe')
#print(w3.eth.getTransactionReceipt('0x99a69be583dc0d4ee447875156057daf8e6af76a510b654ba03bc8d59ef41d0d'))
#injectData(accounts[1])
#transact(accounts[1], accounts[2], 1000)

print(accounts)
print(getCompanyBalance(accounts[1], 'quick'))
print(getCompanyBalance(accounts[2], 'quick'))
'''
#addCompany()
#injectData(accounts[1])
# issueToken(accounts[0], None, None, None, None, 200, None)
# getCompanyTokenNum(accounts[0], -1)
#transact(accounts[1], accounts[2], 1500)
'''

app = Flask(__name__)


# http://127.0.0.1:5000/demo/api/v1.0/accounts
@app.route('/demo/api/v1.0/accounts/', methods=['GET', 'POST'])
def FlaskGetAccounts():
    comps = getCompanys()
    # return "successCallback"+"("+jsonify(comps)+")" #并返回这个添加的task内容，和状态码
    return "callback" + "(" + json.dumps(comps) + ")"


# http://127.0.0.1:5000/demo/api/v1.0/balance/0x90292b8fEf35De29a95874ac04872F34C0B3704C
@app.route('/demo/api/v1.0/balance/<string:addr>', methods=['GET', 'POST'])
def FlaskGetAccountBalance(addr):
    balance = getCompanyBalance(addr, 'quick')
    # return jsonify({'addr': addr, 'balance': balance}), 201  #并返回这个添加的task内容，和状态码
    return "callback" + "(" + json.dumps({'addr': addr, 'balance': balance}) + ")"


# http://127.0.0.1:5000/demo/api/v1.0/tokens/0x90292b8fEf35De29a95874ac04872F34C0B3704C/0
@app.route('/demo/api/v1.0/tokens/<string:addr>/<int:status>', methods=['GET', 'POST'])
def FlaskGetTokens(addr, status):
    jsontokens = getCompanyJsonTokens(addr, status)
    # return jsonify({addr: jsontokens}), 201  #并返回这个添加的task内容，和状态码
    return "callback" + "(" + json.dumps({addr: jsontokens}) + ")"


# http://127.0.0.1:5000/demo/api/v1.0/history/0x90292b8fEf35De29a95874ac04872F34C0B3704C
@app.route('/demo/api/v1.0/history/<string:addr>', methods=['GET', 'POST'])
def FlaskGetAccountHistory(addr):
    jsonhistory = getCompanyHistory(addr)
    # return jsonify(jsonhistory), 201
    return "callback" + "(" + json.dumps(jsonhistory) + ")"


# http://127.0.0.1:5000/demo/api/v1.0/setamount/0x90292b8fEf35De29a95874ac04872F34C0B3704C/10000
@app.route('/demo/api/v1.0/setamount/<string:addr>/<int:amount>', methods=['GET', 'POST'])
def FlaskSetAmount(addr, amount):
    balance = setCompanyAmount(addr, amount)
    # return jsonify({addr: balance}), 201  #并返回这个添加的task内容，和状态码
    return "callback" + "(" + json.dumps({addr: balance}) + ")"


# http://127.0.0.1:5000/demo/api/v1.0/transfer/0x90292b8fEf35De29a95874ac04872F34C0B3704C/0x1B07C35cD769010eD986c9B6E751670620A7c340/1000
@app.route('/demo/api/v1.0/transfer/<string:payer>/<string:payee>/<int:amount>', methods=['GET', 'POST'])
def FlaskTransfer(payer, payee, amount):
    res, msg, sn, tx_hash = transact(payer, payee, amount)
    receipt= {
        'res': res,
        'msg': msg,
        'sn': sn,
        'txhash': tx_hash,
        'detail': {
            'payer': payer,
            'payee': payee,
            'amount': amount
        }
    }
    '''
    return jsonify({
        'res': res,
        'msg': msg,
        'sn': sn,
        'txhash': tx_hash,
        'detail': {
            'payer': payer,
            'payee': payee,
            'amount': amount
        }
    }), 201
    '''
    return "callback" + "(" + json.dumps(receipt) + ")"


# curl -i -H "Content-Type: application/json" -X POST -d '{"addr":"0x84389572984723656435", "amount": 10000}' http://localhost:5000/demo/api/v1.0/setamount
@app.route('/demo/api/v1.0/setamount', methods=['POST'])
def FlaskSetAmountPost():
    if not request.json or not 'addr' in request.json or not 'amount' in request.json:  #如果请求里面没有JSON数据，或者在JSON数据里面，title的内容是空的
        abort(404)  #返回404报错
    addr = request.json['addr']
    amount = request.json['amount']
    receipt = {'addr': addr, 'amount': amount}
    # return jsonify({'addr': addr, 'amount': amount}), 201  #并返回这个添加的task内容，和状态码
    return "callback" + "(" + json.dumps(receipt) + ")"


if __name__ == '__main__':
    app.run(host='0.0.0.0')
'''
w3 = myProvider("ws://127.0.0.1:8546", "bcfan")
# contractInfo = deploy(contractSourceCode, w3)
contractInfo = getContract('./cinfo.json')
invoke(w3, contractInfo, 100, 'thread')
'''
'''
# Display the default greeting from the contract
print('Default contract greeting: {}'.format(greeter.functions.greet().call()))

print('Setting the greeting to Nihao...')
tx_hash = greeter.functions.setGreeting('Nihao').transact()

# Wait for transaction to be mined...
tx_receipt = w3.eth.waitForTransactionReceipt(tx_hash)
print(tx_receipt)

# Display the new greeting value
print('Updated contract greeting: {}'.format(greeter.functions.greet().call()))

# When issuing a lot of reads, try this more concise reader:

assert reader.greet() == "Nihao"
print(reader.getCount())
'''
'''
reader = ConciseContract(contract)
print(reader.greet())
txParameters = {'from': w3.eth.coinbase, 'gas': 10000000}
w3.eth.waitForTransactionReceipt(contract.functions.setGreeting('Nihao').transact(txParameters))
print(reader.greet())
'''

#w3.eth.waitForTransactionReceipt()
#print(reader.getCount())