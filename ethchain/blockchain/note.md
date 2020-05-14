#!/usr/bin/bash
# geth --datadir "/home/xshrim/lab/ethlab/ethbase" init genesis.json
# 初始化说明:
# 初始化操作只需要进行一次
# genesis.json 参数说明:
# mixhash:	与nonce配合用于挖矿，由上一个区块的一部分生成的hash。注意他和nonce的设置需要满足以太坊的Yellow paper, 4.3.4. Block Header Validity,
# nonce:	nonce就是一个64位随机数，用于挖矿，注意他和mixhash的设置需要满足以太坊的Yellow paper, 4.3.4. Block Header Validity,
# difficulty:	设置当前区块的难度，如果难度过大，cpu挖矿就很难，这里设置较小难度
# alloc:	用来预置账号以及账号的以太币数量，因为私有链挖矿比较容易，所以我们不需要预置有币的账号，需要的时候自己创建即可以。
# coinbase:	矿工的账号，随便填
# timestamp:	设置创世块的时间戳
# parentHash:	上一个区块的hash值，因为是创世块，所以这个值是0
# extraData:	附加信息，随便填，可以填你的个性信息
# gasLimit:	该值设置对GAS的消耗总量限制，用来限制区块能包含的交易信息总和，因为我们是私有链，所以填最大。

# 启动说明:
# --identity:  本节点身份
# --nodiscover: 本节点不会被其他节点发现(除非其他节点手动添加本节点)
# --maxpeers: 最多允许多少个节点连接本节点
# --networkid: 本节点的网络ID(公有链网络ID为1)
# --unlock: 本节点自动解锁账号(需配合--password选项)
# --password: 本节点账号密码文件路径(不能是密码本身)
# --mine: 本节点自动挖矿(需要事先至少有一个账号充当矿工)
# --verbosity: 本节点日志详细程度(默认3)
# --gasprice: 本节点处理交易矿工需要收取的交易费, 为0则无交易费依然会处理交易
# --rpc: 激活本节点的RPC接口(geth默认激活)
# --rpcaddr: 本节点RPC接口监听的地址(geth默认localhost)
# --rpcport: 本节点HTTP-RPC接口(短链接)监听的端口(geth默认8545)
# --wsport: 本节点WS-RPC接口(长链接)监听的短裤(geth默认8546)
# --rpccorsdomain: 本节点允许哪些URL连接并执行RPC任务(geth默认允许所有)
# --rpcapi: 本节点允许RPC接入者使用什么API
# --datadir: 本节点数据存放目录(区块数据和账号数据)
# --port: 本节点网络监听端口(geth默认30303)

# 开发者模式说明:
# 启动console时使用--dev即表示启用开发者模式, 开发者模式并不需要手动进行创世块的初始化(如果初始化创世块后再启用开发者模式则会导致服务起不来), 此模式下会默认预分配一个开发者账户(有巨额余额)并自动开启挖矿, 开发者账号无需解锁即可进行交易(其他账号仍需解锁). 如果再额外指定 --dev.period 0 参数则表示仅当出现交易的时候才自动挖一次矿(保证交易立即得到确认)

# 账号说明:
# 进入控制台后创建三个账号:
# personal.newAccount() 或 personal.newAccount('账号密码')
# 账号地址/公钥最后20个字节                     账号密码
# 0xa5a811b5c35468f88120a10ddb72b4bc7b54cd7e    bcfan
# 0x059297dc5c6870fde7b56840eb04b99f2c6c2ff9    maple
# 0x768b91f31713d7867f688c7683af734ea98a01bb    shrim
# 创建账号时, 提供的账号密码和一系列参数经由kdf生成函数计算得出解密密钥. 以太坊客户端使用ECDSA算法自动为账号生成非对称密钥对(公钥+私钥), 私钥即是账号的根本, 公钥是可以根据私钥计算得出的, 以太坊取公钥的最后20个字节作为账号的地址(公钥是64个十六进制数组成的字符串, 即32字节, 后20字节即公钥字符串的后40个十六进制字符, 以太坊在此40个十六进制字符前加0x标示为十六进制数), 账号地址也就是账号的钱包. 对于私钥, 以太坊客户端使用强对称加密算法cipher将私钥, 加密算法参数(系统固定的)和由账号密码生成的解密密钥三者一起进行加密, 生成加密的私钥, 加密的私钥和账号地址, 对称加密算法和非对称加密算法及其参数, 以及mac值等信息一起写入一个json文件中, 也就是以太坊的keystore. 因此, 账号创建者必须备份keystore并牢记账号密码, 二者缺一便无法提取出账号的私钥. 没有私钥便无法进行交易.
# 将解密密钥和加密的私钥连在一起, 然后使用sha3-256算法计算出的hash(即摘要)就是mac值, 当使用错误的账号密码时, 仍然可以生成错误的解密密钥, 使用keystore中加密的私钥和错误的解密密钥进行sha3算法后得出的hash值肯定与mac值不符, 由此可以验证账号密码的真假.  
# 矿工的账号名为coinbase, 它默认是本地账号中的第一个账号(即eth.accounts[0]), 挖矿的奖励收入会发放给coinbase账号. 使用miner.setEtherbase()可以修改coinbase账号

# 其他说明:
# geth 命令可以直接启动ethereum服务, 而gethconsole命令则是启动服务后再进入控制台进行交互
# 可以使用nohup geth --identity "myeth" --networkid 111 --nodiscover --maxpeers 10 --port "30303"  --dev --dev.period 0 --rpc --rpcapi "db,eth,net,web3" --rpcport "8545" --wsport "8546" --datadir "/home/xshrim/lab/ethlab/ethbase"  2> geth.log & 让服务后台运行
# 当节点的ethereum服务正在运行的时候, 可以使用ipc或者rpc方式attach到该服务的控制台中(geth attach http://localhost:8545 或 geth attach ethbase/geth.ipc)
# 还可以使用JSON API(即REST API)与ethereum服务进行交互, 如查询账号:curl -X POST -H "Content-Type":application/json --data '{"jsonrpc":"2.0", "method":"eth_accounts","params":[],"id":9}' localhost:8545
# ethereum的区块链以及状态数据库(全局账户信息)等所有数据均保存在leveldb中, 使用geth removedb可以删除leveldb数据库, 即区块链和状态库均被删除, 也可以通过删除--datadir目录下的相应文件删除leveldb和本地账户(keystore)
# 使用miner.start()进行挖矿时, 可忽视null返回值, 实际上log文件和eth.blockNumber都能看到区块在增加
# 使用sendTransaction()进行转账交易时, 交易不会马上生效, 而是处在挂起状态(使用txpool.status查看交易状态), 需要等待下一个区块产生, 交易才能被写入确认. 此时再次开始挖矿直到产生区块, 挂起的交易便可以完成.
# 可以将挖矿规则设置为生成一个区块后就停止挖矿, 命令为(注意这是一条命令): miner.start(1);admin.sleepBlocks(1);miner.stop()
# ethereum出于安全考虑, 账号一段时间后会被锁定(默认600s), 因此发起交易前需要使用密码解锁(personal.unlockAccount('账号地址') 或 personal.unlockAccount('账号地址','账号密码',解锁持续时间))
# 常用命令: eth.getBalance(账号地址)查看账号余额, eth.accounts查看账号列表, personal.listAccounts查看账号, personal.listWallet查看钱包, eth.blockNumer查看块数量/链高度, eth.syncing进行同步, eth.getBlock(块号)获取块信息, eth.getTransaction(交易号)获取交易信息, web3.toWei(数量, 单位)和web3.fromWei(数量, 单位)转换币单位, txpool.status查看交易状态, miner.start(线程数), admin.sleepBlocks(块数)和miner.stop()开始和停止挖矿, eth.getTransactionCount(账号地址, 块状态[latest,earliest,pending])获取账号交易次数, 也就是nonce.
