import os, time, json

from web3 import Web3
from web3.middleware import geth_poa_middleware
from flask import Flask, jsonify, request, abort


class Block:
    def __init__(self, difficulty, gasLimit, gasUsed, ihash, logsBloom, miner, mixHash, nonce, number, parentHash, stateRoot,
                 transactionsRoot, receiptsRoot, sha3Uncles, size, timestamp, totalDifficulty, transactions, uncles):
        self.difficulty = difficulty
        self.gasLimit = gasLimit
        self.gasUsed = gasUsed
        self.hash = ihash
        self.logsBloom = logsBloom
        self.miner = miner
        self.mixHash = mixHash
        self.nonce = nonce
        self.number = number
        self.parentHash = parentHash
        self.stateRoot = stateRoot
        self.transactionsRoot = transactionsRoot
        self.receiptsRoot = receiptsRoot
        self.sha3Uncles = sha3Uncles
        self.size = size
        self.timestamp = timestamp
        self.totalDifficulty = totalDifficulty
        self.transactions = transactions
        self.uncles = uncles

    def __str__(self):
        return str({
            "difficulty": self.difficulty,
            "gasLimit": self.gasLimit,
            "gasUsed": self.gasUsed,
            "hash": self.hash,
            "logsBloom": self.logsBloom,
            "miner": self.miner,
            "mixHash": self.mixHash,
            "nonce": self.nonce,
            "number": self.number,
            "parentHash": self.parentHash,
            "stateRoot": self.stateRoot,
            "transactionsRoot": self.transactionsRoot,
            "receiptsRoot": self.receiptsRoot,
            "sha3Uncles": self.sha3Uncles,
            "size": self.size,
            "timestamp": self.timestamp,
            "totalDifficulty": self.totalDifficulty,
            "transactions": self.transactions,
            "uncles": self.uncles
        })

    def show(self):
        return {
            "difficulty": self.difficulty,
            "gasLimit": self.gasLimit,
            "gasUsed": self.gasUsed,
            "hash": self.hash,
            "logsBloom": self.logsBloom,
            "miner": self.miner,
            "mixHash": self.mixHash,
            "nonce": self.nonce,
            "number": self.number,
            "parentHash": self.parentHash,
            "stateRoot": self.stateRoot,
            "transactionsRoot": self.transactionsRoot,
            "receiptsRoot": self.receiptsRoot,
            "sha3Uncles": self.sha3Uncles,
            "size": self.size,
            "timestamp": self.timestamp,
            "totalDifficulty": self.totalDifficulty,
            "transactions": self.transactions,
            "uncles": self.uncles
        }


class Transaction:
    def __init__(self, blockHash, blockNumber, ifrom, to, value, gas, gasPrice, ihash, iinput, nonce, transactionIndex, r, s, v):
        self.blockHash = blockHash
        self.blockNumber = blockNumber
        self.fm = ifrom
        self.to = to
        self.value = value
        self.gas = gas
        self.gasPrice = gasPrice
        self.hash = ihash
        self.input = iinput
        self.nonce = nonce
        self.transactionIndex = transactionIndex
        self.r = r
        self.s = s
        self.v = v

    def __str__(self):
        return str({
            "blockHash": self.blockHash,
            "blockNumber": self.blockNumber,
            "from": self.fm,
            "to": self.to,
            "value": self.value,
            "gas": self.gas,
            "gasPrice": self.gasPrice,
            "hash": self.hash,
            "input": self.input,
            "nonce": self.nonce,
            "transactionIndex": self.transactionIndex,
            "r": self.r,
            "s": self.s,
            "v": self.v
        })

    def show(self):
        return {
            "blockHash": self.blockHash,
            "blockNumber": self.blockNumber,
            "from": self.fm,
            "to": self.to,
            "value": self.value,
            "gas": self.gas,
            "gasPrice": self.gasPrice,
            "hash": self.hash,
            "input": self.input,
            "nonce": self.nonce,
            "transactionIndex": self.transactionIndex,
            "r": self.r,
            "s": self.s,
            "v": self.v
        }


class Receipt:
    def __init__(self, blockHash, blockNumber, contractAddress, transactionIndex, transactionHash, root, ifrom, to,
                 cumulativeGasUsed, gasUsed, logsBloom, logs, status):
        self.blockHash = blockHash
        self.blockNumber = blockNumber
        self.contractAddress = contractAddress
        self.transactionIndex = transactionIndex
        self.transactionHash = transactionHash
        self.root = root
        self.fm = ifrom
        self.to = to
        self.cumulativeGasUsed = cumulativeGasUsed
        self.gasUsed = gasUsed
        self.logsBloom = logsBloom
        self.logs = logs
        self.status = status

    def __str__(self):
        return str({
            "blockHash": self.blockHash,
            "blockNumber": self.blockNumber,
            "contractAddress": self.contractAddress,
            "transactionIndex": self.transactionIndex,
            "transactionHash": self.transactionHash,
            "root": self.root,
            "from": self.fm,
            "to": self.to,
            "cumulativeGasUsed": self.cumulativeGasUsed,
            "gasUsed": self.gasUsed,
            "logsBloom": self.logsBloom,
            "logs": self.logs,
            "status": self.status
        })

    def show(self):
        return {
            "blockHash": self.blockHash,
            "blockNumber": self.blockNumber,
            "contractAddress": self.contractAddress,
            "transactionIndex": self.transactionIndex,
            "transactionHash": self.transactionHash,
            "root": self.root,
            "from": self.fm,
            "to": self.to,
            "cumulativeGasUsed": self.cumulativeGasUsed,
            "gasUsed": self.gasUsed,
            "logsBloom": self.logsBloom,
            "logs": self.logs,
            "status": self.status
        }


class Log:
    def __init__(self, address, blockHash, blockNumber, logIndex, transactionIndex, transactionHash, data, topics, removed):
        self.address = address
        self.blockHash = blockHash
        self.blockNumber = blockNumber
        self.logIndex = logIndex
        self.transactionIndex = transactionIndex
        self.transactionHash = transactionHash
        self.data = data
        self.topics = topics
        self.removed = removed

    def __str__(self):
        return str({
            "address": self.address,
            "blockHash": self.blockHash,
            "blockNumber": self.blockNumber,
            "logIndex": self.logIndex,
            "transactionIndex": self.transactionIndex,
            "transactionHash": self.transactionHash,
            "data": self.data,
            "topics": self.topics,
            "removed": self.removed
        })

    def show(self):
        return {
            "address": self.address,
            "blockHash": self.blockHash,
            "blockNumber": self.blockNumber,
            "logIndex": self.logIndex,
            "transactionIndex": self.transactionIndex,
            "transactionHash": self.transactionHash,
            "data": self.data,
            "topics": self.topics,
            "removed": self.removed
        }


def myProvider(provider_url, passphrase="test"):
    # web3.py instance
    # w3 = Web3(Web3.EthereumTesterProvider())
    # PoA共识机制下api需要注入PoA中间件
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


def handle_event(event):
    print(event)
    # and whatever


def log_loop(event_filter, poll_interval):
    while True:
        for event in event_filter.get_new_entries():
            handle_event(event)
        time.sleep(poll_interval)


def createBlock(iblock):
    transactions = []
    uncles = []
    for transaction in iblock.transactions:
        transactions.append(str(w3.toHex(transaction)))
    for uncle in iblock.uncles:
        uncles.append(str(w3.toHex(uncle)))
    block = Block(
        int(iblock['difficulty']), int(iblock['gasLimit']), int(iblock['gasUsed']), str(w3.toHex(iblock['hash'])),
        str(w3.toHex(iblock['logsBloom'])), str(iblock['miner']),
        str(w3.toHex(iblock['mixHash'])), str(w3.toHex(iblock['nonce'])), int(iblock['number']),
        str(w3.toHex(iblock['parentHash'])), str(w3.toHex(iblock['stateRoot'])), str(w3.toHex(iblock['transactionsRoot'])),
        str(w3.toHex(iblock['receiptsRoot'])), str(w3.toHex(iblock['sha3Uncles'])), int(iblock['size']), int(iblock['timestamp']),
        int(iblock['totalDifficulty']), transactions, uncles)
    return block


def createTransaction(itransaction):
    to = itransaction['to'] if itransaction['to'] is not None else ""
    transaction = Transaction(
        str(w3.toHex(itransaction['blockHash'])), int(itransaction['blockNumber']), str(itransaction['from']), str(to),
        int(itransaction['value']), int(itransaction['gas']), int(itransaction['gasPrice']), str(w3.toHex(itransaction['hash'])),
        str(itransaction['input']), int(itransaction['nonce']), int(itransaction['transactionIndex']),
        str(w3.toHex(itransaction['r'])), str(w3.toHex(itransaction['s'])), int(itransaction['v']))
    return transaction


def createLog(ilog):
    topics = []
    for topic in ilog['topics']:
        topics.append(str(w3.toHex(topic)))
    if 'address' in ilog.keys():
        address = ilog['address'] if ilog['address'] is not None else ""
    else:
        address = ""

    log = Log(
        str(address), str(w3.toHex(ilog['blockHash'])), int(ilog['blockNumber']), int(ilog['logIndex']),
        int(ilog['transactionIndex']), str(w3.toHex(ilog['transactionHash'])), str(ilog['data']), topics, str(ilog['removed']))
    return log


def createReceipt(ireceipt):
    logs = []
    for log in ireceipt.logs:
        logs.append(createLog(log).show())

    contractAddress = ireceipt['contractAddress'] if ireceipt['contractAddress'] is not None else ""
    if 'root' in ireceipt.keys():
        root = ireceipt['root'] if ireceipt['root'] is not None else ""
    else:
        root = ""
    to = ireceipt['to'] if ireceipt['to'] is not None else ""
    if 'status' in ireceipt.keys():
        status = ireceipt['status'] if ireceipt['status'] is not None else ""
    else:
        status = 1

    receipt = Receipt(
        str(w3.toHex(ireceipt['blockHash'])), int(ireceipt['blockNumber']), str(contractAddress), int(
            ireceipt['transactionIndex']), str(w3.toHex(ireceipt['transactionHash'])), str(root), str(ireceipt['from']), str(to),
        int(ireceipt['cumulativeGasUsed']), int(ireceipt['gasUsed']), str(w3.toHex(ireceipt['logsBloom'])), logs, str(status))
    return receipt


def getVersion():
    return w3.version


def getBlockChainStatus():
    status = 'normal'
    if w3.isConnected():
        if w3.eth.syncing:
            status = 'syncing'
        else:
            status = 'normal'
    else:
        status = 'down'
    return status


def getBlockNumber():
    return w3.eth.blockNumber


def getTUNumberByRange(fb=0, tb='latest'):
    num = [0, 0]
    if tb == 'latest':
        tb = getBlockNumber()
    for i in range(fb + 1, tb + 1, 1):
        num[0] += w3.eth.getBlockTransactionCount(i)
        num[1] += w3.eth.getUncleCount(i)
    return num


def getTUNumber():
    latest = getBlockNumber()
    current = tus['blockNumber']
    tcount = tus['transactionCount']
    ucount = tus['uncleCount']
    rnum = getTUNumberByRange(current, latest)
    tcount += rnum[0]
    ucount += rnum[1]
    tus['blockNumber'] = latest
    tus['transactionCount'] = tcount
    tus['uncleCount'] = ucount
    return [tcount, ucount]


def getRecentBlocks(seconds):
    blocks = []
    timestamp = time.time() - seconds
    blockNumber = getBlockNumber()
    for i in range(blockNumber, 0, -1):
        block = createBlock(getBlockByNum(i))
        if block.timestamp < timestamp:
            break
        else:
            blocks.append(block)
    return blocks


def getLatestBlocks(number):
    blocks = []
    blockNumber = getBlockNumber()
    for i in range(blockNumber, blockNumber - number, -1):
        block = createBlock(getBlockByNum(i))
        blocks.append(block)
    return blocks


def getViewBlocks(start, offset):
    blocks = []
    if start == 'latest':
        start = getBlockNumber()
    else:
        start = int(start)
    offset = int(offset)
    for i in range(start, start - offset, -1):
        block = createBlock(getBlockByNum(i))
        blocks.append(block)
    return blocks


def getRecentTransactions(seconds):
    transactions = []
    blocks = getRecentBlocks(seconds)
    for block in blocks:
        for transaction in reversed(block.transactions):
            transactions.append(createTransaction(getTransactionByHash(transaction)))
    return transactions


def getLatestTransactions(number):
    transactions = []
    cblockNumber = getBlockNumber()
    while True:
        block = createBlock(getBlockByNum(cblockNumber))
        for transaction in reversed(block.transactions):
            if number > 0:
                transactions.append(createTransaction(getTransactionByHash(transaction)))
                number -= 1
            else:
                break
        if number > 0:
            cblockNumber -= 1
        else:
            break
    return transactions


def getTxPool():
    txp = w3.txpool.status
    pending = int(txp['pending'], 16)
    queued = int(txp['queued'], 16)
    return [pending, queued]


def getAccounts():
    return w3.eth.accounts


def getAccountNum():
    return len(w3.eth.accounts)


def getBlockByNum(num):
    return w3.eth.getBlock(num)


def getBlockByHash(ihash):
    return w3.eth.getBlock(ihash)


def getTransactionByHash(ihash):
    return w3.eth.getTransaction(ihash)


def getReceiptByHash(ihash):
    return w3.eth.getTransactionReceipt(ihash)


def getBaseInfo():
    blockChainStatus = getBlockChainStatus()
    contractAddress = "0x331Ca136561C53A18590FE4fdAcFda8F8abCBBeb"
    # contractAddress = str(contractInfo['address'])
    blockNumber = getBlockNumber()
    tuNumber = getTUNumber()
    uncleBlockNumber = tuNumber[1]
    transactionNumber = tuNumber[0]
    addressNumber = getAccountNum()
    txpool = getTxPool()
    pendingTxNumber = txpool[0]
    queuedTxNumber = txpool[1]

    return {
        "blockChainStatus": blockChainStatus,
        "contractAddress": contractAddress,
        "blockNumber": blockNumber,
        "uncleBlockNumber": uncleBlockNumber,
        "transactionNumber": transactionNumber,
        "addressNumber": addressNumber,
        "pendingTxNumber": pendingTxNumber,
        "queuedTxNumber": queuedTxNumber
    }


def queryEvent():
    # b'^a\x82L\x9a\xcbg\xfd\xa9\xc5\x00\x10\x184\xea\n\xf2x{\xa6\x06\xca3\x01\xce\xfd\xd2}\xf4a\x17\xff'
    # str(w3.sha3(text="uMYHF94C6"))
    # myfilter = contract.eventFilter('TXReceipt', {'fromBlock': 175, 'toBlock': 'latest', 'filter': {'payee': getAccounts()[1]}})
    myfilter = contract.eventFilter(
        'TXReceipt',
        {
            'fromBlock': 0,
            'toBlock': 'latest',
            #'filter': {
            #    'sn': "TJah0LveE".encode('utf8')
            #}
        })  # sn must be bytes32
    eventlist = myfilter.get_all_entries()
    print(eventlist)


w3 = myProvider("ws://172.16.201.191:8546", "node1")
# w3 = myProvider("ws://127.0.0.1:8546", "test")
'''
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

tus = {'blockNumber': 0, 'transactionCount': 0, 'uncleCount': 0}
getBaseInfo()
# for trans in getRecentTransactions(3600):
#     print(trans)
'''
for i in range(getBlockNumber()):
    block = createBlock(getBlockByNum(i))
    # block = getBlockByNum(10)
    # print(str(block['miner']))
    print(block)
    # print(getBlockByHash('0xbeeaead53d3db71f9068f04c3246e3fc71724c4efa7ccb46f3d1b00e29f7172c'))
    if len(block.transactions) > 0:
        trans = createTransaction(getTransactionByHash(block.transactions[0]))
        print(trans)

        receipt = createReceipt(getReceiptByHash(trans.hash))
        print(receipt)

        for log in receipt.logs:
            print(log)
'''

app = Flask(__name__)


# http://127.0.0.1:5000/rs/accounts
@app.route('/rs/accounts', methods=['GET', 'POST'])
def FlaskGetAccounts():
    try:
        account = getAccountNum()
    except Exception as ex:
        return "callback" + "(" + json.dumps({'error': str(ex)}) + ")"
    # return "successCallback"+"("+jsonify(comps)+")" #并返回这个添加的task内容，和状态码
    return "callback" + "(" + json.dumps({'count': account}) + ")"


# http://127.0.0.1:5000/rs/baseinfo
@app.route('/rs/baseinfo', methods=['GET', 'POST'])
def FlaskGetBaseInfo():
    try:
        baseinfo = getBaseInfo()
    except Exception as ex:
        return "callback" + "(" + json.dumps({'error': str(ex)}) + ")"
    # return "successCallback"+"("+jsonify(comps)+")" #并返回这个添加的task内容，和状态码
    return "callback" + "(" + json.dumps({'baseinfo': baseinfo}) + ")"


# http://127.0.0.1:5000/rs/block/5
# http://127.0.0.1:5000/rs/block/0xddc2a2ed8caeb5160d2b7bc77137c8e332e7a1b1e236d36f79af9ba8eb8d5272
@app.route('/rs/block/<string:idhash>', methods=['GET', 'POST'])
def FlaskGetBlock(idhash):
    try:
        if idhash.isdigit():
            block = createBlock(getBlockByNum(int(idhash)))
        else:
            block = createBlock(getBlockByHash(str(idhash)))
    except Exception as ex:
        return "callback" + "(" + json.dumps({'error': str(ex)}) + ")"
    # return jsonify({'addr': addr, 'balance': balance}), 201  #并返回这个添加的task内容，和状态码
    return "callback" + "(" + json.dumps({'block': block.show()}) + ")"


# http://127.0.0.1:5000/rs/transaction/0xe3557dd42ef48884550d751dd4727e56041e04c6d2dc764bc9146616cc2ccca0
@app.route('/rs/transaction/<string:idhash>', methods=['GET', 'POST'])
def FlaskGetTransaction(idhash):
    try:
        transaction = createTransaction(getTransactionByHash(str(idhash)))
    except Exception as ex:
        return "callback" + "(" + json.dumps({'error': str(ex)}) + ")"
    # return jsonify({addr: jsontokens}), 201  #并返回这个添加的task内容，和状态码
    return "callback" + "(" + json.dumps({'transaction': transaction.show()}) + ")"


# http://127.0.0.1:5000/rs/transactions/5
# http://127.0.0.1:5000/rs/transactions/0xddc2a2ed8caeb5160d2b7bc77137c8e332e7a1b1e236d36f79af9ba8eb8d5272
@app.route('/rs/transactions/<string:idhash>', methods=['GET', 'POST'])
def FlaskGetTransactions(idhash):
    transactions = {}
    try:
        if idhash.isdigit():
            block = createBlock(getBlockByNum(int(idhash)))
        else:
            block = createBlock(getBlockByHash(str(idhash)))
        for txhash in block.transactions:
            transaction = createTransaction(getTransactionByHash(str(txhash)))
            transactions[str(txhash)] = transaction.show()
    except Exception as ex:
        return "callback" + "(" + json.dumps({'error': str(ex)}) + ")"
    # return jsonify({addr: jsontokens}), 201  #并返回这个添加的task内容，和状态码
    return "callback" + "(" + json.dumps({'transactions': transactions}) + ")"


# http://127.0.0.1:5000/rs/receipt/0xe3557dd42ef48884550d751dd4727e56041e04c6d2dc764bc9146616cc2ccca0
@app.route('/rs/receipt/<string:idhash>', methods=['GET', 'POST'])
def FlaskGetReceipt(idhash):
    try:
        receipt = createReceipt(getReceiptByHash(str(idhash)))
    except Exception as ex:
        return "callback" + "(" + json.dumps({'error': str(ex)}) + ")"
    # return jsonify({addr: jsontokens}), 201  #并返回这个添加的task内容，和状态码
    return "callback" + "(" + json.dumps({'receipt': receipt.show()}) + ")"


# http://127.0.0.1:5000/rs/recent/block/7200
@app.route('/rs/recent/block/<string:seconds>', methods=['GET', 'POST'])
def FlaskGetRecentBlocks(seconds):
    blocks = {}
    try:
        for block in getRecentBlocks(int(seconds)):
            blocks[block.hash] = block.show()
    except Exception as ex:
        return "callback" + "(" + json.dumps({'error': str(ex)}) + ")"
    # return jsonify({addr: jsontokens}), 201  #并返回这个添加的task内容，和状态码
    return "callback" + "(" + json.dumps({'blocks': blocks}) + ")"


# http://127.0.0.1:5000/rs/recent/transaction/7200
@app.route('/rs/recent/transaction/<string:seconds>', methods=['GET', 'POST'])
def FlaskGetRecentTransactions(seconds):
    transactions = {}
    try:
        for transaction in getRecentTransactions(int(seconds)):
            transactions[transaction.hash] = transaction.show()
    except Exception as ex:
        return "callback" + "(" + json.dumps({'error': str(ex)}) + ")"
    # return jsonify({addr: jsontokens}), 201  #并返回这个添加的task内容，和状态码
    return "callback" + "(" + json.dumps({'transactions': transactions}) + ")"


# http://127.0.0.1:5000/rs/latest/block/10
@app.route('/rs/latest/block/<string:number>', methods=['GET', 'POST'])
def FlaskGetLatestBlocks(number):
    blocks = {}
    try:
        for block in getLatestBlocks(int(number)):
            blocks[block.hash] = block.show()
    except Exception as ex:
        return "callback" + "(" + json.dumps({'error': str(ex)}) + ")"
    # return jsonify({addr: jsontokens}), 201  #并返回这个添加的task内容，和状态码
    return "callback" + "(" + json.dumps({'blocks': blocks}) + ")"


# http://127.0.0.1:5000/rs/latest/transaction/10
@app.route('/rs/latest/transaction/<string:number>', methods=['GET', 'POST'])
def FlaskGetLatestTransactions(number):
    transactions = {}
    try:
        for transaction in getLatestTransactions(int(number)):
            transactions[transaction.hash] = transaction.show()
    except Exception as ex:
        return "callback" + "(" + json.dumps({'error': str(ex)}) + ")"
    # return jsonify({addr: jsontokens}), 201  #并返回这个添加的task内容，和状态码
    return "callback" + "(" + json.dumps({'transactions': transactions}) + ")"


# http://127.0.0.1:5000/rs/view/block/latest/20
@app.route('/rs/view/block/<string:start>/<string:offset>', methods=['GET', 'POST'])
def FlaskGetViewBlocks(start, offset):
    blocks = {}
    try:
        for block in getViewBlocks(str(start), int(offset)):
            blocks[block.number] = block.show()
    except Exception as ex:
        return "callback" + "(" + json.dumps({'error': str(ex)}) + ")"
    # return jsonify({addr: jsontokens}), 201  #并返回这个添加的task内容，和状态码
    return "callback" + "(" + json.dumps({'blocks': blocks}) + ")"


if __name__ == '__main__':
    app.run(host='0.0.0.0', debug=True)
