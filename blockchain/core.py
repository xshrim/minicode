import sys, time, json, web3, timeit

from web3 import Web3
from threading import Thread
from solc import compile_source
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
    compiled_sol = compile_source(contract_source_code)
    contract_interface = compiled_sol['<stdin>:Greeter']

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