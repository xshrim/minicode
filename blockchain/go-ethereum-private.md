# 私有链环境部署

## 目录

<!-- TOC -->

- [私有链环境部署](#私有链环境部署)
    - [目录](#目录)
    - [前置准备](#前置准备)
        - [go安装](#go安装)
        - [go-ethereum安装](#go-ethereum安装)
        - [关闭防火墙](#关闭防火墙)
    - [环境部署](#环境部署)
        - [环境说明](#环境说明)
            - [PoA概述](#poa概述)
            - [网络结构](#网络结构)
        - [部署过程](#部署过程)
            - [建立账号](#建立账号)
            - [生成密码文件](#生成密码文件)
            - [生成创世json文件](#生成创世json文件)
            - [编辑创世json文件](#编辑创世json文件)
            - [初始化节点](#初始化节点)
            - [启动boot节点](#启动boot节点)
            - [启动各节点](#启动各节点)
            - [验证节点间连通性](#验证节点间连通性)
            - [交易测试](#交易测试)
            - [api调用测试](#api调用测试)
            - [新增认证节点](#新增认证节点)
            - [删除认证节点](#删除认证节点)
    - [压力测试](#压力测试)
        - [压测脚本](#压测脚本)
        - [压测步骤](#压测步骤)
    - [Sample](#sample)

<!-- /TOC -->

## 前置准备

### go安装

```bash
wget https://studygolang.com/dl/golang/go1.10.linux-amd64.tar.gz
tar -zxvf go1.10.linux-amd64.tar.gz -C /usr/local/share/    #解压
mkdir -p /root/go         #创建go本地运行目录
vim /etc/profile          #创建go环境变量
    ......
    #Go
    export GOROOT=/usr/local/share/go
    export GOPATH=/root/go
    export PATH=$GOROOT/bin:$GOPATH/bin:$PATH
source /etc/profile
```

### go-ethereum安装

需要预先安装Git, Go和C编译器.

```bash
git clone https://github.com/ethereum/go-ethereum.git
cd  go-ethereum
git checkout v1.8.9  (git tag查看最新版本)
make geth
make all
```

### 关闭防火墙

如果是多主机环境, 防火墙可能导致节点间无法通信, 需要开放相应节点端口(默认30301, 30313)或关闭防火墙.

```bash
# centos 6
service iptables stop
chkconfig iptables off

#centos 7
systemctl stop firewalld
systemctl disable firewalld
```

---

## 环境部署

### 环境说明

#### PoA概述

go-ethereum自1.6版起开始支持可插拔的共识机制, 当前go-ethereum客户端内置的共识机制有两种: 一种是总所周知的用于公有链的ethash共识(即PoW, 工作量证明), 另一种是用于私有链和联盟链的clique共识(即PoA, 基于认证的证明). PoA共识机制在go-ethereum下称为clique, 在Parity下称为aura. 本共识机制利用区块上的extraData字段, 指定允许进行挖矿生成区块的节点(即指定节点上的矿工账号), 这些节点称为认证节点, 不允许进行挖矿和区块生成的节点称为普通节点. 初始的认证节点由创世json文件指定, 后续可以动态增减认证节点. clique共识下不再依赖矿工消耗CPU和内存算力来动态出块, 而是直接指定区块生成的间隔时间(即挖矿几乎无难度, 瞬间完成), 所有认证节点达到生成区块的时间时同时挖矿形成区块, 但是认证节点间存在轮换机制(当前区块号除以认证节点个数的余数为新的出块节点, 即inturn节点), 每轮出块只有一个inturn节点,inturn节点会马上向其他节点广播区块请求验证, 而其他节点的区块会延迟广播, 保证生效的区块是inturn节点的区块. 生成区块的间隔时间默认是15s(实验证明即使设置为1s仍可稳定运行), 将此间隔设置为0s表示只有在出现pending的交易的时候才生成区块(相当于dev模式的dev.period参数, 事实上dev模式即是使用clique共识, 只是dev模式下不支持多节点网络). 动态增减认证节点是通过投票机制实现的, 当节点状态的改变获得超过半数认证节点的验证即可生效. 目前PoA共识机制的ethereum区块链网络已经作为ethereum公用测试网络稳定运行中.

**参考文档**: [Clique共识算法](https://ethfans.org/posts/Clique-Consensus-Algorithm)

#### 网络结构

直接在一台物理机上模拟五节点(同IP不同端口)ethereum网络. 具体如下:

| 节点  |  类型  |     数据存储     |  IP地址   | 节点端口 | RPC端口 | WS端口 | 账号口令 |          账号地址(账号生成后获取)          |
|:-----:|:------:|:----------------:|:---------:|:--------:|:-------:|:------:|:--------:|:------------------------------------------:|
| node1 | signer | /tmp/priv/node1/ | 127.0.0.1 |  30313   |  8515   |  8516  |   test   | 2ab766077d074e0b4976213f3b836e40b04cbc71 | |
| node2 | signer | /tmp/priv/node2/ | 127.0.0.1 |  30323   |  8525   |  8526  |   test   | 470a0985df071f964c9ea1a0428cc00f341106c0 | |
| node3 | signer | /tmp/priv/node3/ | 127.0.0.1 |  30333   |  8535   |  8536  |   test   | 1290cc0186210482522921c8813615a3c7024f6b | |
| node4 | normal | /tmp/priv/node4/ | 127.0.0.1 |  30343   |  8545   |  8546  |   test   | 88efb100cfd993390ac06d9a117d918725d07dfc | |
| node5 | normal | /tmp/priv/node5/ | 127.0.0.1 |  30353   |  8555   |  8556  |   test   | 5960563fab4d74398c3cc9673dab93c250b7cd70 | |

### 部署过程

#### 建立账号

为每一个节点建立一个挖矿账号

```bash
cd /tmp
mkdir priv
cd /tmp/priv
for i in 1 2 3 4 5; do geth --datadir ./node$i account new; done    # 为每一个账号输入口令
```

#### 生成密码文件

为每一个节点的初始账号生成一个密码文件以便解锁账号
```bash
for i in 1 2 3 4 5; do echo "node"$i > /tmp/priv/node$i/node$i.pass; done
```


#### 生成创世json文件

使用go-ethereum自带的puppeth工具生成genesis文件.

```bash
puppeth
# Please specify a network name to administer (no spaces or hyphens, please)     # 指定网络名称
> priv
# What would you like to do? (default = stats)                                    # 选择操作
#  1. Show network stats
#  2. Configure new genesis
#  3. Track new remote server
#  4. Deploy network components
> 2
# Which consensus engine to use? (default = clique)                                # 选择共识机制
#  1. Ethash - proof-of-work
#  2. Clique - proof-of-authority
> 2
# How many seconds should blocks take? (default = 15)                               # 选择出块间隔(0表示有交易挂起才出块)
> 0
# Which accounts are allowed to seal? (mandatory at least one)                      # 允许哪些节点挖矿(即哪些节点作为初始认证节点)
> 0x2ab766077d074e0b4976213f3b836e40b04cbc71
> 0x470a0985df071f964c9ea1a0428cc00f341106c0
> 0x1290cc0186210482522921c8813615a3c7024f6b
> 0x
# Which accounts should be pre-funded? (advisable at least one)                     # 为哪些账号地址预分配以太币(几乎无限)
> 0x2ab766077d074e0b4976213f3b836e40b04cbc71
> 0x470a0985df071f964c9ea1a0428cc00f341106c0
> 0x1290cc0186210482522921c8813615a3c7024f6b
> 0x88efb100cfd993390ac06d9a117d918725d07dfc
> 0x5960563fab4d74398c3cc9673dab93c250b7cd70
> 0x
# Specify your chain/network ID if you want an explicit one (default = random)      # 设置网络ID
> 11111
# What would you like to do? (default = stats)                                      # genesis相关参数设置完成后保存到本地
#  1. Show network stats
#  2. Manage existing genesis
#  3. Track new remote server
#  4. Deploy network components
> 2
# 1. Modify existing fork rules                                                     # 导出配置文件
# 2. Export genesis configuration
# 3. Remove genesis configuration
> 2
# Which file to save the genesis into? (default = priv.json)                        # 命名配置文件(保持默认)
> (回车)
cat priv.json
```

#### 编辑创世json文件

编辑创世json文件的gasLimit参数以保证大交易(如合约部署)可以顺利执行.

```bash
vi priv.json
"gasLimit": "0x8000000"      # 默认值为 0x47b760
```

#### 初始化节点

使用创世json文件初始化各节点

```bash
for i in 1 2 3 4 5; do geth --datadir ./node$i init priv.json; done
```

#### 启动boot节点

boot节点用于各对等节点间互相发现. 各节点互相发现后, boot节点即可取消(节点退出后重新发现仍需要boot节点). 当然也可以不用boot节点, 各节点启动后上手动使用admin.addPeer()添加节点即可.

```bash
bootnode --genkey=boot.key                                                          # 生成boot节点key文件
nohup bootnode --nodekey=boot.key &> /tmp/priv/node1/boot.log &                                                      # 启动boot节点
```

**注**: boot节点启动后会提供一个enode连接码供ethereum节点连接:
> cat boot.log
> enode://5e9cc210572f70853f2bc700ae50d65f10cad4910aa5785e0a9f4bf5c2097720f160348718c7b72ae9bd8ede7190cd681f971f6790647369c22d31eb1ee7c1f2@[::]:30301
启动ethereum节点使用此连接时需要将[::]部分改为boot节点的IP.

#### 启动各节点

节点启动时的identity和networkid需要与生成创世json文件时一致. gasprice和targetgaslimit 参数规定了节点生成区块的gas消耗上限和所能接受的单个gas的最小以太币价格. gasprice建议设置极小(因私有链无须考虑交易费), 但不要设置为0(因go-ethereum中依然有计算gas消耗的内容, 设为0可能出现意外). targetgaslimit可以设置大一些, 以便支持更大的交易和合约, 以及更大的交易吞吐量(交易消耗的gas数量固定的情况下, targetgaslimit越大, 一个区块能打包的交易数就越多, 高并发交易下性能更高). rpcapi和wsapi开启的功能模块在生产环境中应该有所删减. rpccorsdomain和wsorigins在生产环境中应该指定ip范围. node4和node5可以去掉mine和minerthreads参数, 但后续如果变成认证节点则无法挖矿了, 除非重启节点.

```bash
# 分别开启一个单独的终端窗口执行节点服务启动的命令, 启动后需要输入初始矿工账号的口令以解锁账号
# 窗口1
nohup geth --identity=node1 --networkid=11111 --maxpeers=50 --port=30313 --syncmode=full --gasprice=1 --targetgaslimit=471238800 --rpc --rpcapi=db,eth,net,web3,personal,miner,admin,debug,txpool --rpcaddr=0.0.0.0 --rpccorsdomain=* --rpcport=8515 --ws --wsapi=db,eth,net,web3,personal,miner,admin,debug,txpool --wsaddr=0.0.0.0 --wsorigins=* --wsport=8516 --mine --minerthreads=1 --etherbase=0x2ab766077d074e0b4976213f3b836e40b04cbc71 --unlock=0x2ab766077d074e0b4976213f3b836e40b04cbc71 --password=/tmp/priv/node1/node1.pass --datadir=/tmp/priv/node1 --bootnodes=enode://5e9cc210572f70853f2bc700ae50d65f10cad4910aa5785e0a9f4bf5c2097720f160348718c7b72ae9bd8ede7190cd681f971f6790647369c22d31eb1ee7c1f2@127.0.0.1:30301 &> /tmp/priv/node1/node1.log &
# 窗口2
nohup geth --identity=node2 --networkid=11111 --maxpeers=50 --port=30323 --syncmode=full --gasprice=1 --targetgaslimit=471238800 --rpc --rpcapi=db,eth,net,web3,personal,miner,admin,debug,txpool --rpcaddr=0.0.0.0 --rpccorsdomain=* --rpcport=8525 --ws --wsapi=db,eth,net,web3,personal,miner,admin,debug,txpool --wsaddr=0.0.0.0 --wsorigins=* --wsport=8526 --mine --minerthreads=1 --etherbase=0x470a0985df071f964c9ea1a0428cc00f341106c0 --unlock=0x470a0985df071f964c9ea1a0428cc00f341106c0 --password=/tmp/priv/node2/node2.pass --datadir=/tmp/priv/node2 --bootnodes=enode://5e9cc210572f70853f2bc700ae50d65f10cad4910aa5785e0a9f4bf5c2097720f160348718c7b72ae9bd8ede7190cd681f971f6790647369c22d31eb1ee7c1f2@127.0.0.1:30301 &> /tmp/priv/node2/node2.log &
# 窗口3
nohup geth --identity=node3 --networkid=11111 --maxpeers=50 --port=30333 --syncmode=full --gasprice=1 --targetgaslimit=471238800 --rpc --rpcapi=db,eth,net,web3,personal,miner,admin,debug,txpool --rpcaddr=0.0.0.0 --rpccorsdomain=* --rpcport=8535 --ws --wsapi=db,eth,net,web3,personal,miner,admin,debug,txpool --wsaddr=0.0.0.0 --wsorigins=* --wsport=8536 --mine --minerthreads=1 --etherbase=0x1290cc0186210482522921c8813615a3c7024f6b --unlock=0x1290cc0186210482522921c8813615a3c7024f6b --password=/tmp/priv/node3/node3.pass --datadir=/tmp/priv/node3 --bootnodes=enode://5e9cc210572f70853f2bc700ae50d65f10cad4910aa5785e0a9f4bf5c2097720f160348718c7b72ae9bd8ede7190cd681f971f6790647369c22d31eb1ee7c1f2@127.0.0.1:30301 &> /tmp/priv/node3/node3.log &
# 窗口4
nohup geth --identity=node4 --networkid=11111 --maxpeers=50 --port=30343 --syncmode=full --gasprice=1 --targetgaslimit=471238800 --rpc --rpcapi=db,eth,net,web3,personal,miner,admin,debug,txpool --rpcaddr=0.0.0.0 --rpccorsdomain=* --rpcport=8545 --ws --wsapi=db,eth,net,web3,personal,miner,admin,debug,txpool --wsaddr=0.0.0.0 --wsorigins=* --wsport=8546 --mine --minerthreads=1 --etherbase=0x88efb100cfd993390ac06d9a117d918725d07dfc --unlock=0x88efb100cfd993390ac06d9a117d918725d07dfc --password=/tmp/priv/node4/node4.pass --datadir=/tmp/priv/node4 --bootnodes=enode://5e9cc210572f70853f2bc700ae50d65f10cad4910aa5785e0a9f4bf5c2097720f160348718c7b72ae9bd8ede7190cd681f971f6790647369c22d31eb1ee7c1f2@127.0.0.1:30301 &> /tmp/priv/node4/node4.log &
# 窗口5
nohup geth --identity=node5 --networkid=11111 --maxpeers=50 --port=30353 --syncmode=full --gasprice=1 --targetgaslimit=471238800 --rpc --rpcapi=db,eth,net,web3,personal,miner,admin,debug,txpool --rpcaddr=0.0.0.0 --rpccorsdomain=* --rpcport=8555 --ws --wsapi=db,eth,net,web3,personal,miner,admin,debug,txpool --wsaddr=0.0.0.0 --wsorigins=* --wsport=8556 --mine --minerthreads=1 --etherbase=0x5960563fab4d74398c3cc9673dab93c250b7cd70 --unlock=0x5960563fab4d74398c3cc9673dab93c250b7cd70 --password=/tmp/priv/node5/node5.pass --datadir=/tmp/priv/node5 --bootnodes=enode://5e9cc210572f70853f2bc700ae50d65f10cad4910aa5785e0a9f4bf5c2097720f160348718c7b72ae9bd8ede7190cd681f971f6790647369c22d31eb1ee7c1f2@127.0.0.1:30301 &> /tmp/priv/node5/node5.log &
```

#### 验证节点间连通性

```bash
> admin.peers                 # 正常应该显示除自身以外的所有节点
```

#### 交易测试

在node5的ethereum控制台下新建应该账号, 并发起一笔转账交易, 测试私有链网络的可用性

```bash
> personal.newAccount('test')
> eth.accounts
> eth.sendTransaction({from:eth.accounts[0], to:eth.accounts[1], value: web3.toWei(3, 'ether')})
> txpool.status
> eth.blockNumber
```

**注**: 如果发现该笔交易出现"Discarded bad propagated block"错误, 则重启node1节点即可解决.

#### api调用测试

使用Web3.py进行api调用测试, 由于使用clique共识, 需要引入geth_poa_middleware模块.

```bash
sudo pip3 install web3
```

简单调用代码

```python
# from web3.auto import w3      标准端口下可以自动连接
from web3 import Web3
from web3.middleware import geth_poa_middleware
w3 = Web3(Web3.WebsocketProvider("ws://127.0.0.1:8546")
w3.middleware_stack.inject(geth_poa_middleware, layer=0)

w3.isConnected()
w3.version.node

w3.eth.getBlock('latest')
w3.eth.getTransaction(w3.eth.getBlock(1).transactions[0])
```

#### 新增认证节点

在超过半数的当前认证节点上执行认证节点新增命令, 等新生成几个区块, 各认证节点投票后, 超过半数的认证节点便可以认可新的认证节点了, 新认证节点也会加入区块生成的行列. 建议在新增认证节点时在当前所有认证节点上都执行新增命令.

```bash
# 在node1和node2上均执行, 然后在任意节点发起几个交易即可看到node4也参与挖矿生成区块了
> clique.propose('0x88efb100cfd993390ac06d9a117d918725d07dfc', true)
> eth.sendTransaction({from:eth.accounts[0], to:eth.accounts[1], value: web3.toWei(3, 'ether')})
> eth.sendTransaction({from:eth.accounts[0], to:eth.accounts[1], value: web3.toWei(3, 'ether')})
> clique.getSnapshot()
```

#### 删除认证节点

同理在超过半数的当前认证节点上执行节点删除命令, 等新生成几个区块, 各认证节点投票后, 超过半数认证节点便可以剔除指定的认证节点了.

```bash
> clique.discard('0x88efb100cfd993390ac06d9a117d918725d07dfc')
> eth.sendTransaction({from:eth.accounts[0], to:eth.accounts[1], value: web3.toWei(3, 'ether')})
> eth.sendTransaction({from:eth.accounts[0], to:eth.accounts[1], value: web3.toWei(3, 'ether')})
> clique.getSnapshot()
```

## 压力测试

### 压测脚本

使用ethereum的python api(web3.py)进行压力测试

```python
# core.py: 主要函数库
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
```

```python
# tx.py: 多线程模拟大量并发请求
import core

w3 = core.myProvider("ws://127.0.0.1:8546", "test")
#contractInfo = core.deploy(core.contractSourceCode, w3)
contractInfo = core.getContract('./cinfo.json')
core.benchmark(w3, contractInfo, 2000, 'thread')
```

```python
# tps.py: 监控每秒交易量(TPS)
import core
# from web3.auto import w3

w3 = core.myProvider("ws://127.0.0.1:8546")
# contractInfo = deploy(contractSourceCode, w3)
# contractInfo = core.getContract('./cinfo.json')
# core.invoke(w3, contractInfo, 1000, 'thread')
core.monitor(w3)
```

```python
# demo.py: 查看区块链统计信息
import core
from web3.auto import w3
from web3 import Web3
from web3.middleware import geth_poa_middleware

# 直接使用web3自动的自动连接功能
w3.middleware_stack.inject(geth_poa_middleware, layer=0)
num = w3.eth.blockNumber
tnum = 0
btlist = []
# 计算每个区块的交易数量
for i in range(0, num + 1):
    trans = w3.eth.getBlock(i).transactions
    btlist.append((i, len(trans)))
    for tran in trans:
        tnum += 1
        tranhash = Web3.toHex(tran)
        res = w3.eth.getTransactionReceipt(tranhash)
        print(str(i) + ':' + str(res))
        #info = w3.eth.getTransaction(trans)
        #pint(info)

# 调用合约内计数函数获取交易计数
contractInfo = core.getContract('./cinfo.json')
contract = w3.eth.contract(
    address=contractInfo['address'],
    abi=contractInfo['abi'],
)
print(contract.functions.getCount().call())

# 查看每个区块交易数
for bt in btlist:
    print(bt)
# 查看区块总数和总交易数
print(num, tnum)
```

### 压测步骤

1. 启动五节点ethereum网络.
2. 修改tx.py中的注释进行合约部署(仅首次运行需要)
3. 运行tps.py脚本开始监控.
4. 修改tx.py中的注释直接调用合约模拟大量并发交易.
5. 完成后运行demo.py脚本确认交易完成情况.

**注**:

- 建议运行tps.py脚本后马上运行tx.py脚本(或者调换两个脚本运行顺序), 以免计时偏差过大.
- clique共识的即时出块策略会导致大量交易并发时各个区块内的交易数量较少且不均匀, 产生的区块量较多. 比较适合交易并发小的情况. 对于交易并发量大的情况, 建议使用固定出块策略, 能够保证每个区块内有较多的交易并且相对均匀, 产生的区块量也会较少.
- 本压测过程因在即时出块策略下模拟高交易并发, 因此TPS并不高, 大概在60TPS左右; 如果在高并发情景下使用固定出块策略, 同时调整targetgaslimit使一个区块能打包更多交易(使用默认值4712388的情况下, 一个区块可以打包合约产生的交易个数约为120个), 则TPS应该会有所提高.
- geth本身还支持一些调优参数, 后续再研究.
- 本压测受环境所限, 硬件配置与生产环境相差较大, 网络环境却比生产环境更好, 交易本身也比较简单, 这些因素都会影响实际的性能表现.

## Sample

物理机部署三节点的命令样例.

```bash
nohup geth --identity=node1 --networkid=111 --maxpeers=50 --port=30313 --syncmode=full --gasprice=1 --targetgaslimit=471238800 --rpc --rpcapi=db,eth,net,web3,personal,miner,admin,debug,txpool --rpcaddr=0.0.0.0 --rpccorsdomain=* --rpcport=8545 --ws --wsapi=db,eth,net,web3,personal,miner,admin,debug,txpool --wsaddr=0.0.0.0 --wsorigins=* --wsport=8546 --mine --minerthreads=1 --etherbase=0xa8150525e6e90a16c4fd3881c651093c618dc20f --unlock=0xa8150525e6e90a16c4fd3881c651093c618dc20f --password=/root/node1/node1.pass --datadir=/root/node1 --bootnodes=enode://5e06027be450833c1783f69c2a57393ef101806fe4175e85b3edc0ef6d368ce8de1fca625503e2b9daa4d98f1b1871acebeedcefab5cdb8d9ffa8b70f0e34a60@192.168.115.17:30301 &> /root/node1/node1.log &

nohup geth --identity=node2 --networkid=111 --maxpeers=50 --port=30313 --syncmode=full --gasprice=1 --targetgaslimit=471238800 --rpc --rpcapi=db,eth,net,web3,personal,miner,admin,debug,txpool --rpcaddr=0.0.0.0 --rpccorsdomain=* --rpcport=8545 --ws --wsapi=db,eth,net,web3,personal,miner,admin,debug,txpool --wsaddr=0.0.0.0 --wsorigins=* --wsport=8546 --mine --minerthreads=1 --etherbase=0x39b760038ff2faad4f4939ddd44aa46ca3ee89b6 --unlock=0x39b760038ff2faad4f4939ddd44aa46ca3ee89b6 --password=/root/node2/node2.pass --datadir=/root/node2 --bootnodes=enode://5e06027be450833c1783f69c2a57393ef101806fe4175e85b3edc0ef6d368ce8de1fca625503e2b9daa4d98f1b1871acebeedcefab5cdb8d9ffa8b70f0e34a60@192.168.115.17:30301 &> /root/node2/node2.log &

nohup geth --identity=node3 --networkid=111 --maxpeers=50 --port=30313 --syncmode=full --gasprice=1 --targetgaslimit=471238800 --rpc --rpcapi=db,eth,net,web3,personal,miner,admin,debug,txpool --rpcaddr=0.0.0.0 --rpccorsdomain=* --rpcport=8545 --ws --wsapi=db,eth,net,web3,personal,miner,admin,debug,txpool --wsaddr=0.0.0.0 --wsorigins=* --wsport=8546 --mine --minerthreads=1 --etherbase=0xaa35a703bdf7f874ca5f02366dbe0490051275db --unlock=0xaa35a703bdf7f874ca5f02366dbe0490051275db --password=/root/node3/node3.pass --datadir=/root/node3 --bootnodes=enode://5e06027be450833c1783f69c2a57393ef101806fe4175e85b3edc0ef6d368ce8de1fca625503e2b9daa4d98f1b1871acebeedcefab5cdb8d9ffa8b70f0e34a60@192.168.115.17:30301 &> /root/node3/node3.log &

```