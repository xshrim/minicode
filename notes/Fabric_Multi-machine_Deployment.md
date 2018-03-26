# Fabric v1.1.0多机部署

本文将详细介绍Hyperledger Fabric v1.1.0版的多机部署过程并安装链码进行简单测试, 我们依然采用Fabric-Sample演示的1 Orderer(solo) + 4 Peer(2 org)的架构, 将orderer和peer的服务以源码编译的方式直接部署在五台物理机操作系统上(实验环境使用VirtualBox虚拟机代替), 而不再依赖已经部署好相关服务的Docker容器.

------

## 环境准备

### 硬件环境

一共需要使用6台物理机/虚拟机, 网络划分与用途分配如下:
|主机名                  |操作系统              |IPv4地址   |用途                                                                         |备注              |
| --------------------- | ------------------- | --------- | -------------------------------------------------------------------------- | ----------------- |
|orderer.example.com    |CentOS-7.4-x64(1708) |10.0.2.5   |orderer共识节点, 接收交易背书, 共识排序, 生成并传递区块                          |必需                |
|peer0.org1.example.com |CentOS-7.4-x64(1708) |10.0.2.6   |组织1的背书锚节点, 处理交易请求, 生成背书, 验证区块并记账, 实现通道内各组织互相发现  |必需                |
|peer1.org1.example.com |CentOS-7.4-x64(1708) |10.0.2.7   |组织1的背书节点, 处理交易请求, 生成背书, 验证区块并记账                           |非必需, 防止单点故障 |
|peer0.org2.example.com |CentOS-7.4-x64(1708) |10.0.2.8   |组织2的背书锚节点, 处理交易请求, 生成背书, 验证区块并记账, 实现通道内各组织互相发现  |必需                |
|peer1.org2.example.com |CentOS-7.4-x64(1708) |10.0.2.9   |组织2的背书节点, 处理交易请求, 生成背书, 验证区块并记账                           |非必需, 防止单点故障 |

*[注] 生产环境下orderer节点和各组织的peer节点应处于不同网段但可互通, 实验环境下直接使用VirtualBox的NatNetwork实现多虚拟机同网段互通(直接使用NAT虚拟机之间是无法互通的), 使用NatNetwork端口转发实现宿主机登录和调试虚拟机(NatNetwork默认支持虚拟机ping通宿主机, 不支持宿主机ping通虚拟机).*

### 软件环境

> 以下基础环境设置步骤需要在所有节点执行.

#### 系统设置

设置主机名, 关闭防火墙, 设置主机名-IP映射.

```bash
systemctl stop firewalld                        #开启防火墙,只开放fabric节点间通信的相关端口也可以
systemctl disable firewalld
hostnamectl set-hostname orderer.example.com    #其他节点类似
vim /etc/hosts
    ......
    10.0.2.5    orderer.example.com
    10.0.2.6    peer0.org1.example.com
    10.0.2.7    peer1.org1.example.com
    10.0.2.8    peer0.org2.example.com
    10.0.2.9    peer1.org2.example.com
```

#### 依赖包

用于编译相关源码.

```bash
yum install snappy-devel zlib-devel bzip2-devel libtool-ltdl-devel -y
```

#### Git

版本要求v1.8+, yum安装或者官网下载, 用于从github下载源码和版本选择.

```bash
yum install git -y
```

#### Go

版本要求 v1.9+, 官网下载, 用于源码编译和链码支持.

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

#### Fabric

版本要求 v1.1.0, git下载源码, 用于编译生成Fabric服务环境.

```bash
go get -u github.com/hyperledger/fabric
#也可以git clone:
#mkdir -p $GOPATH/src/github.com/hyperledger
#cd $GOPATH/src/github.com/hyperledger
#git clone https://github.com/hyperledger/fabric.git
cd $GOPATH/src/github.com/hyperledger/fabric
git tag                  #查看可用的fabric版本
git checkout v1.1.0      #切换为使用最新的v1.1.0版,重要!!!
```

#### Docker

版本要求 v17.00+, yum安装, 用于源码编译和链码运行.

```bash
yum-config-manager --add-repo https://download.daocloud.io/docker/linux/centos/docker-ce.repo       #添加repo
yum install docker-ce -y                                                                            #安装docker
systemctl enable docker
systemctl start docker                                                                              #启动docker并设为开机运行
curl -sSL https://get.daocloud.io/daotools/set_mirror.sh | sh -s http://64ff7d79.m.daocloud.io      #加速docker镜像下载
curl -L https://get.daocloud.io/docker/compose/releases/download/1.12.0/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose                                                              #安装docker-compose, 用于从baseimage和yaml配置文件生成docker镜像并运行
```

#### Go Tools

版本要求 v1.9+, git下载源码, 用于生成相关Go工具.

```bash
cd $GOPATH/src/github.com/hyperledger/fabric
mkdir -p gotools/build/gopath/src/golang.org/x/
cd gotools/build/gopath/src/golang.org/x/
git clone https://github.com/golang/tools.git
git clone https://github.com/golang/lint.git
```

------

## 源码编译

> 以下编译操作需要在所有节点执行.

### 编译gotools

```bash
cd $GOPATH/src/github.com/hyperledger/fabric
make gotools                                          #相关二进制go tools会生成于$GOPATH/bin
```

### 编译fabric

```bash
cd $GOPATH/src/github.com/hyperledger/fabric
mkdir -p build/docker/gotools/bin/
cp $GOPATH/bin/* build/docker/gotools/bin/          #将上一步生成的gotools二进制文件复制到docker相应位置
make buildenv                                       #编译fabric基础环境, 会下载docker baseimage
make orderer                                        #编译fabric相关服务工具
make peer
make configtxgen
make cryptogen                                      #生成的工具位于$GOPATH/src/github.com/hyperledger/fabric/build/bin
mkdir -p /etc/hyperledger/fabric                    #fabric默认的配置文件HOME目录
vim /etc/profile                                    #将工具和配置文件目录加入系统环境变量中
    ......
    #Fabric
    export FABRIC_CFG_PATH=/etc/hyperledger/fabric
    export PATH=$GOPATH/src/github.com/hyperledger/fabric/build/bin:$PATH
source /etc/profile
```

### 重启系统

> 建议完成上述操作后重启一次所有节点, 再执行下面的步骤. 如果是虚拟机实验环境, 建议为每个节点做一次快照.

------

## Fabric网络

> 开始搭建Fabric多节点网络, 以下操作均在/etc/hyperledger/fabric目录下执行, 且不再需要在所有节点均执行一遍,每一步操作需要在哪些节点执行会明确指出.

### 组织结构与安全配置

在yaml配置文件中描述Fabric网络中组织结构, 然后使用cryptogen工具(gotools编译生成)根据yaml配置生成各节点,组织和用户的密钥,数字签名,数字证书和数据传输加密证书.
> 此部分操作只在orderer节点执行(任一节点均可).

#### 配置crypto-config.yaml

本实验使用一个orderer组织, 两个peer组织, orderer组织拥有一个单orderer节点,每个peer组织拥有两个peer节点(一个锚节点一个普通节点), 每个组织除默认的Admin用户外,增加一个普通用户. cli节点并不属于fabric网络的一部分,无需作安全配置. 具体配置内容如下:

```yaml
OrdererOrgs:
  - Name: Orderer
    Domain: example.com
    Specs:
      - Hostname: orderer          #一个orderer节点
PeerOrgs:
  - Name: Org1
    Domain: org1.example.com
    Template:
      Count: 2                     #两个节点
    Users:
      Count: 1                     #一个普通用户
  - Name: Org2
    Domain: org2.example.com
    Template:
      Count: 2
    Users:
      Count: 1
```

#### 生成安全认证相关文件

使用gotools编译生成的cryptogen工具生成相关密钥证书和签名文件, 然后将其分发到其他节点.
> 事实上cryptogen会根据yaml文件的描述生成fabric网络中所有节点的安全文件, 每一个节点只需要与自身相关的那部分安全文件, 实验中直接将全部安全文件分发给了所有节点, 生产环境中不可能如此.

```bash
cd /etc/hyperledger/fabric
cryptogen generate --config=crypto-config.yaml --output crypto-config        #根据yaml文件配置内容生成相关安全文件放入crypto-config目录中
tree ./crypto-config                                                         #查看安全相关文件的结构
scp -r crypto-config root@10.0.2.6:/etc/hyperledger/fabric/                  #将安全文件分发到所有节点
scp -r crypto-config root@10.0.2.7:/etc/hyperledger/fabric/
scp -r crypto-config root@10.0.2.8:/etc/hyperledger/fabric/
scp -r crypto-config root@10.0.2.9:/etc/hyperledger/fabric/
# rm -rf crypto-config.yaml                                                  #删除不再需要的crypto-config.yaml配置文件
```

### 创世区块与通道配置

本实验仅新建一条通道,并将两个peer组织的四个peer节点均加入通道中. 在yaml配置文件中描述创世区块和通道的相关配置信息, 然后使用configtxgen工具(gotools编译生成)根据yaml配置生成创世区块,通道配置文件以及锚节点更新配置文件.
> 此部分操作只在orderer节点执行, 然后将相关配置文件分发给其他节点.

#### 知识说明

* 通道与区块链一一对应, 通道分为系统通道和应用通道, 应用通道供用户使用, 负责承载各种交易, 系统通道则负责管理应用通道.
* 这里的创世区块就是用于创建并连接到系统通道的, 创世区块和系统通道只需存在于orderer节点间. 此创世区块也称orderer创世区块.
* 这里的通道配置文件才是用于创建承载交易的应用通道, 根据通道配置文件可以生成通道所对应的区块链的创世区块. 这个区块链是存在于加入了该通道的所有节点中的(包括orderer节点).
* Fabric的应用通道创世区块中已经包含了组织信息, 已包含的组织的新节点想要加入该通道很容易, 但是想要加入创世区块中未包含的新组织(及其节点)到该通道中就非常困难, 需要使用configtxlator修改应用通道创世区块.

#### 配置configtx.yaml

configtx.yaml文件中描述了系统通道的创世区块和应用通道的相关配置信息, 具体内容如下:

```yaml
Profiles:
    TwoOrgsOrdererGenesis:              #orderer创世区块
        Orderer:
            <<: *OrdererDefaults
            Organizations:
                - *OrdererOrg
        Consortiums:                    #组织的集合体, 可以定义多个集合体
            SampleConsortium:           #Sample集合体中包含两个组织
                Organizations:
                    - *Org1             #每个组织可以处于多个组织集合体中
                    - *Org2
    TwoOrgsChannel:                     #定义一个应用通道
        Consortium: SampleConsortium    #每个通道都对应一个组织集合体, 也对应一个区块链, 意味着每个组织都可以加入多个子链中
        Application:
            <<: *ApplicationDefaults
            Organizations:
                - *Org1                 #允许参与本通道相关应用的组织
                - *Org2
Organizations:
    - &OrdererOrg
        Name: OrdererOrg
        ID: OrdererMSP
        MSPDir: crypto-config/ordererOrganizations/example.com/msp     #指定各组织安全认证文件的位置, 可以根据需要修改
        AdminPrincipal: Role.ADMIN
    - &Org1
        Name: Org1MSP
        ID: Org1MSP
        MSPDir: crypto-config/peerOrganizations/org1.example.com/msp
        AdminPrincipal: Role.ADMIN
        AnchorPeers:                   #每个peer组织可以指定多个锚节点, 普通节点无需指定
            - Host: peer0.org1.example.com
              Port: 7051
    - &Org2
        Name: Org2MSP
        ID: Org2MSP
        MSPDir: crypto-config/peerOrganizations/org2.example.com/msp
        AdminPrincipal: Role.ADMIN
        AnchorPeers:
            - Host: peer0.org2.example.com
              Port: 7051
Orderer: &OrdererDefaults              #共识/排序节点的配置
    OrdererType: solo                  #单orderer节点的solo模式
    Addresses:
        - orderer.example.com:7050
    BatchTimeout: 2s                   #将交易打包到区块前等待的时间, 通常时间越长一个区块里包含的交易就会越多, 生成区块就越不频繁, 但交易完成确认就会越慢
    BatchSize:
        MaxMessageCount: 10            #一个区块里最大的交易数
        AbsoluteMaxBytes: 98 MB        #一个区块的最大大小
        PreferredMaxBytes: 512 KB      #一个区块的建议大小, 如果一个交易消息的大小超过了这个值, 就会被放入更大的区块中
    Kafka:
        Brokers:                       #kafka集群的每一个节点称为一个broker
            - 127.0.0.1:9092
    Organizations:
Application: &ApplicationDefaults
    Organizations:
```

#### 生成创世区块和应用通道配置文件

```bash
cd /etc/hyperledger/fabric
configtxgen -profile TwoOrgsOrdererGenesis -outputBlock genesis.block    #生成orderer创世块
configtxgen -profile TwoOrgsChannel -outputCreateChannelTx mych.tx -channelID mych  #生成交易通道文件mych.tx, 用于生成通道对应的区块链的创世区块
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate Org1MSPanchors.tx -channelID mych -asOrg Org1MSP  #生成更新组织锚节点的配置文件
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate Org2MSPanchors.tx -channelID mych -asOrg Org2MSP
scp Org1MSPanchors.tx root@10.0.2.6:/etc/hyperledger/fabric/       #将锚节点更新配置文件分别分发到对应组织的锚节点(本例中即两个peer0)
scp Org2MSPanchors.tx root@10.0.2.8:/etc/hyperledger/fabric/
# rm -rf configtx.yaml                                             #删除不再需要的configtx.yaml文件
```

### 节点服务配置与启动

依次配置各节点的fabric服务配置文件, 然后启动各节点.

#### Orderer节点配置与启动

> 此部分操作只在orderer节点执行.

##### Orderer节点配置

通过orderer.yaml配置文件配置节点服务参数. yaml配置文件内容如下:

```yaml
General:                        #通用配置
    LedgerType: file            #区块链账本类型, 支持内存账本ram, 简单文件账本json和基于文件的账本file三种, 只有file类型是产品级的
    ListenAddress: 0.0.0.0      #监听IP地址
    ListenPort: 7050            #监听端口
    TLS:                        #grpc消息传输的tls加密设置
        Enabled: true
        PrivateKey: /etc/hyperledger/fabric/tls/server.key
        Certificate: /etc/hyperledger/fabric/tls/server.crt
        RootCAs:
          - /etc/hyperledger/fabric/tls/ca.crt
        ClientAuthRequired: false
        ClientRootCAs:
    LogLevel: debug             #日志记录级别
    LogFormat: '%{color}%{time:2006-01-02 15:04:05.000 MST} [%{module}] %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}'                  #日志记录格式
    GenesisMethod: file                      #创世方法, 支持provisional和file两种类型, 前者在初始化orderer系统通道时根据创世配置参数动态生成创世区块, 后者指定一个创世文件作为创世区块
    GenesisProfile: SampleInsecureSolo                     #由于GenesiseMethod为file, 故此项不生效
    GenesisFile: /etc/hyperledger/fabric/genesis.block     #创世区块的位置
    LocalMSPDir: /etc/hyperledger/fabric/msp        #orderer节点所需的安全认证文件的位置
    LocalMSPID: OrdererMSP    #MSP管理器用于注册安全认证文件的ID, 此ID必须与配置系统通道和创世区块时(configtx.yaml的OrdererGenesis部分)指定的组织中的某一个组织的ID一致
    Profile:                     #为Go pprof性能优化工具启用一个HTTP服务以便作性能分析(https://golang.org/pkg/net/http/pprof)
        Enabled: false           #不启用
        Address: 0.0.0.0:6060    #Go pprof的HTTP服务监听的地址和端口
    BCCSP:                       #区块链密码服务提供者
        Default: SW              #默认使用的密码服务提供者, 支持基于软件(SW)和基于PKCS11硬件(PKCS11)两种实现方式
        SW:                      #基于软件的区块链加密服务的配置参数
            Hash: SHA2           #加密算法
            Security: 256
            FileKeyStore:        #密钥存储位置,如果不指定,默认使用'LocalMSPDir'/keystore
                KeyStore:
FileLedger:                      #账本配置(file和json两种类型)
    Location: /data/hyperledger/production/orderer       #区块存储位置, 如不指定, 每次节点重启都将使用一个新的临时位置
    # The prefix to use when generating a ledger directory in temporary space.
    # Otherwise, this value is ignored.
    Prefix: test-ledger                                  #作为使用临时位置存储账本的目录, 如果指定了区块存储位置, 此项就不会生效
RAMLedger:                      #账本配置(ram类型)
    HistorySize: 1000           #内存账本所支持存储的区块的数量, 如果内存中存储的区块达到上限, 继续追加区块会导致最旧的区块被丢弃
Kafka:                          #基于kafka集群的orderer节点配置, 本实验不采用
# kafka是一种基于发布/订阅模式的分布式消息系统
# fabric网络中, orderer节点集群组成kafka集群, 客户端是kafka集群的Producer(消息生产者), peer是kafka集群的Consumer(消息消费者)
# kafka集群使用ZooKeeper(分布式应用协调服务)管理集群节点, 选举leader.
    Retry:                      #orderer节点到kafka集群的连接失败或者重复请求的重试策略
        # 当一个新通道建立或者一个已存在的通道重载(如orderer节点重启)时,orderer节点会创建一个productor和一个consumer, 由productor向kafka集群发送一条no-op CONNECT消息, 由consumer接收, 只有发送和接收都正常完成, orderer节点才能在该通道上进行读写, 如果失败, orderer会重试, ShortTotal时间内每隔ShortInterval时间重试一次, 之后的LongTotal时间内每隔LongInterval时间重试一次.
        ShortInterval: 5s
        ShortTotal: 10m
        LongInterval: 5m
        LongTotal: 12h
        NetworkTimeouts:                 #连接和读写超时
            DialTimeout: 10s
            ReadTimeout: 10s
            WriteTimeout: 10s
        Metadata:
            RetryBackoff: 250ms          #leader选举过程中元数据请求失败的重试间隔
            RetryMax: 3                  #最大重试次数
        Producer:
            RetryBackoff: 100ms          #向kafka集群发送消息失败后的重试间隔
            RetryMax: 3                  #最大重试次数
        Consumer:
            RetryBackoff: 2s             #从kafka集群拉取消息失败后的重试间隔
    Verbose: false                       #orderer与kafka集群交互是否生成日志
    TLS:                                 #orderer与kafka集群连接的tls加密设置
      Enabled: false                     #不启用加密连接
      PrivateKey:                        #加密连接的密钥证书的存储位置
      Certificate:
      RootCAs:
    Version: 0.10.2.0                    #kafka集群的版本
Debug:                        #orderer节点的调试参数
    BroadcastTraceDir:                   #该orderer节点的广播服务请求保存的位置
    DeliverTraceDir:                     #该orderer节点的传递服务请求保存的位置
```

##### Orderer节点启动

```bash
cd /etc/hyperledger/fabric
cp -r crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp .
cp -r crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/tls .
ls /etc/hyperledger/fabric                #保证该目录下有如下文件和目录: crypto-config、msp、tls、orderer.yaml、genesis.block
orderer start
# nohup orderer start &> orderer.log &    #后台运行
```

#### Peer节点配置与启动

> 此部分操作需要在各peer节点分别执行, 各peer节点的配置文件需要分别配置, 仅以org1的peer0节点(锚节点)为例.

##### Peer节点配置

> 本例中所有配置完成的yaml配置文件下载地址见文后说明

通过core.yaml配置文件配置节点服务参数. yaml配置文件内容如下:

```yaml
logging:                                  #日志配置
    level:       info                     #默认日志记录级别, 仅当--logging-level命令行参数和CORE_LOGGING_LEVEL配置参数均未设置时生效
    cauthdsl:   warning                   #下方指定的细化到各组件的日志记录级别将覆盖level, --logging-level和CORE_LOGGING_LEVEL指定的默认级别
    gossip:     warning
    grpc:       error
    ledger:     info
    msp:        warning
    policies:   warning
    peer:
        gossip: warning
    format: '%{color}%{time:2006-01-02 15:04:05.000 MST} [%{module}] %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}'                                                                #日志记录格式
peer:                                      #节点配置
    id: peer0.org1.example.com             #节点唯一ID
    networkId: example                     #节点所属的逻辑网络分区的ID
    listenAddress: 0.0.0.0:7051            #监听地址和端口
    chaincodeListenAddress: peer0.org1.example.com:7052    #节点链码的监听地址和端口, 如未指定, 则使用listenAddress的IP地址和7052端口
    address: peer0.org1.example.com:7051   #该项配置在peer节点上,表示peer与同组织内其他peer通信的地址; 配置在cli节点上表示peer与其通信的地址
    addressAutoDetect: false               #节点是以编程方式确定自身的address, 此项在docker容器方式的节点上非常有用
    gomaxprocs: -1                         #Go运行时使用的CPU核数
    gossip:                                #闲话算法/谣言传播算法/病毒感染算法, peer间同步信息的方法, 保证了最终各peer的信息将是一致的
        bootstrap: 127.0.0.1:7051          #peer服务启动时尝试连接的节点列表, 只能指定同组织内的其他peer, 连接到不同组织的peer时会被拒绝. 通常锚节点设置为连接自身, 非锚节点设置为连接本组织的锚节点
        useLeaderElection: true            #是否使用动态算法进行leader选举, 一个组织内可以有多个leader节点, leader节点是与排序服务建立连接并使用delivery协议从排序服务拉取账本区块,将区块散步到本组织其他peer的peer节点, 在有大量节点的环境中推荐设置此项
        orgLeader: false                   #是否是leader节点, 此项与useLeaderElection互斥, 不能同时为true, 当useLeaderElection为false时, 组织内必须至少有一个peer的orgLeader值为true.
        endpoint: 0.0.0.0:7051             #公布给本组织内其他peer的地址和端口
        maxBlockCountToStore: 100          #驻留在内存中的最大区块数
        maxPropagationBurstLatency: 10ms   #连续推送消息的最短间隔
        maxPropagationBurstSize: 10        #触发消息推送前允许驻留的最大消息数
        propagateIterations: 1             #一条消息推送到其他peer的次数
        propagatePeerNum: 3                #向多少个peer推送消息
        pullInterval: 4s                   #拉取消息的频率
        pullPeerNum: 3                     #向多少各peer拉取消息
        requestStateInfoInterval: 4s       #拉取peer状态信息的频率
        publishStateInfoInterval: 4s       #推送peer状态信息的频率
        stateInfoRetentionInterval:        #peer状态信息的过期时间
        publishCertPeriod: 10s             #发布消息的生命周期
        skipBlockVerification: false       #是否忽略区块验证消息
        dialTimeout: 3s                    #拨号超时时间
        connTimeout: 2s                    #连接超时时间
        recvBuffSize: 20                   #接收消息的缓冲大小
        sendBuffSize: 200                  #发送消息的缓冲大小
        digestWaitTime: 1s
        requestWaitTime: 1s
        responseWaitTime: 2s               #拉取引擎结束拉取的等待时间
        aliveTimeInterval: 5s              #检查连接活跃度的间隔
        aliveExpirationTimeout: 25s        #连接过期时间(多长时间不活跃视为过期)
        reconnectInterval: 25s             #重连间隔
        externalEndpoint: peer0.org1.example.com:7051          #公布给其他组织的地址和端口, 如果不指定, 其他组织将无法知道本peer的存在
        election:                           #leader选举配置
            startupGracePeriod: 15s         #等待稳定成员参与选举的最长时间
            membershipSampleInterval: 1s    #测试peer稳定性的时间间隔
            leaderAliveThreshold: 10s       #自上一次选举以来再次选举的时间间隔
            leaderElectionDuration: 5s      #选举时间

        pvtData:
            pullRetryThreshold: 60s
            minPeers: 1
            maxPeers: 1
            transientstoreMaxBlockRetention: 1000
    events:                                              #事件配置
        address: peer0.org1.example.com:7053             #事件服务监听的地址和端口
        buffersize: 100                                  #保证无阻塞发送的事件缓存数量
        timeout: 10ms     #事件发送超时事件, 如果事件缓存已满, timeout < 0, 事件被丢弃; timeout > 0, 阻塞直到超时丢弃, timeout = 0, 阻塞但保证事件一定会被发送出去
        timewindow: 15m                                  #允许peer和client时间不一致的最大时间差
    tls:                               #数据传输加密设置
        enabled: true                  #peer间数据传输是否使用加密
        clientAuthRequired: false      #客户端连接到peer是否需要使用加密
        cert:                          #证书密钥的位置, 各peer应该填写各自相应的路径
            file: /etc/hyperledger/fabric/tls/server.crt
        key:
            file: /etc/hyperledger/fabric/tls/server.key
        rootcert:
            file: /etc/hyperledger/fabric/tls/ca.crt
        clientRootCAs:
            files:
              - /etc/hyperledger/fabric/tls/ca.crt
        serverhostoverride:
    fileSystemPath: /data/hyperledger/production      #peer数据存储位置(包括账本,状态数据库等)
    BCCSP:                                            #区块链加密提供者, 与orderer设置相同
        Default: SW
        SW:
            Hash: SHA2
            Security: 256
            FileKeyStore:
                KeyStore:
    mspConfigPath: /etc/hyperledger/fabric/msp
    localMspId: Org1MSP                               #各peer填写各自组织的MSP ID
    deliveryclient:
        reconnectTotalTimeThreshold: 3600s            #交付服务交付失败后尝试重连的时间
    profile:
        enabled:     false
        listenAddress: 0.0.0.0:6060
    handlers:                                         #自定义信息过滤和处理模块
        authFilters:
          -
            name: DefaultAuth
        decorators:
          -
            name: DefaultDecorator
    validatorPoolSize:                                #处理交易验证的并发数, 默认是CPU的核数
vm:                                                   #docker环境配置
    endpoint: unix:///var/run/docker.sock
    docker:
        tls:
            enabled: false
            ca:
                file: docker/ca.crt
            cert:
                file: docker/tls.crt
            key:
                file: docker/tls.key
        attachStdout: false
        hostConfig:
            NetworkMode: host
            Dns:
            LogConfig:
                Type: json-file
                Config:
                    max-size: "50m"
                    max-file: "5"
            Memory: 2147483648
chaincode:                                       #链码设置
    peerAddress:
    id:
        path:
        name:
    builder: $(DOCKER_NS)/fabric-ccenv:$(ARCH)-$(PROJECT_VERSION)
    golang:
        runtime: $(BASE_DOCKER_NS)/fabric-baseos:$(ARCH)-$(BASE_VERSION)
        dynamicLink: false
    car:
        runtime: $(BASE_DOCKER_NS)/fabric-baseos:$(ARCH)-$(BASE_VERSION)
    java:
        Dockerfile:  |
            from $(DOCKER_NS)/fabric-javaenv:$(ARCH)-$(PROJECT_VERSION)
    node:
        runtime: $(BASE_DOCKER_NS)/fabric-baseimage:$(ARCH)-$(BASE_VERSION)
    startuptimeout: 300s
    executetimeout: 30s
    mode: net
    keepalive: 0
    system:                                       #系统链码白名单
        cscc: enable
        lscc: enable
        escc: enable
        vscc: enable
        qscc: enable
        rscc: disable
    systemPlugins:
    logging:
      level:  info
      shim:   warning
      format: '%{color}%{time:2006-01-02 15:04:05.000 MST} [%{module}] %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}'
ledger:                                             #账本设置
  blockchain:
  state:
    stateDatabase: goleveldb                        #区块链状态数据库选择, 支持leveldb和CouchDB
    couchDBConfig:                                  #CouchDB配置, CouchDB是非关系模型的文档数据库, 数据以json文档的方式存储, 相比levelDB, 数据查看更直观
       couchDBAddress: 127.0.0.1:5984
       username:
       password:
       maxRetries: 3
       maxRetriesOnStartup: 10
       requestTimeout: 35s
       queryLimit: 10000
       maxBatchUpdateSize: 1000
  history:
    enableHistoryDatabase: true                        #是否存储key update历史, 仅对levelDB生效
metrics:                                               #服务度量监控
        enabled: false
        reporter: statsd
        interval: 1s
        statsdReporter:
              address: 0.0.0.0:8125
              flushInterval: 2s
              flushBytes: 1432
        promReporter:
              listenAddress: 0.0.0.0:8080
```

##### Peer节点启动

每个节点都要相应修改core.yaml配置文件, 依次登陆各节点,启动fabric服务

```bash
cd /etc/hyperledger/fabric
cp -r crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/msp .   #各节点复制对应的msp和tls目录
cp -r crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls .
ls /etc/hyperledger/fabric   #保证该目录下有如下文件和目录: crypto-config、msp、tls、core.yaml、Org1MSPanchors.tx(Org1的peer0)、Org2MSPanchors.tx(Org2的peer0)
peer node start
# nohup peer node start &> peer1.org1.log &       #后台运行
```

### 应用通道创建与使用

使用之前生成的应用通道配置文件生成应用通道对应的区块链的创世区块, 然后将此区块分发到所有需要加入此通道的peer节点, 使用此区块加入通道中.
> 以下步骤需要保证orderer和peer节点的fabric服务均处于运行状态

#### 配置管理员环境变量

管理员环境变量只能临时使用, 不能添加到/etc/profile等配置文件中, 否则可能导致节点重启fabric服务时出现身份冲突(一个组织的多个节点都使用管理员身份启动fabric服务会报错).
> 此操作需要在所有peer节点执行
> 创建通道, 加入通道, 更新锚节点, 安装链码, 实例化链码, 更新链码等关键性操作都必须以管理员身份执行

```bash
export CORE_PEER_LOCALMSPID="Org1MSP"         #不同组织的peer节点设置相应的环境变量
export CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
```

#### 创建通道

> 此操作只在组织1的锚节点(peer0.org1.example.com)执行

```bash
cd /etc/hyperledger/fabric
peer channel create -o orderer.example.com:7050 -c mych -f mych.tx --tls true --cafile /etc/hyperledger/fabric/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem                                 #使用管理员身份生成mych.block
scp mych.block root@10.0.2.7:/etc/hyperledger/fabric/        #将应用通道创世区块分发到其他peer节点
scp mych.block root@10.0.2.8:/etc/hyperledger/fabric/
scp mych.block root@10.0.2.9:/etc/hyperledger/fabric/
```

#### 加入通道

使用通道创世区块将各peer节点加入到通道中.
> 此操作需要在所有peer节点执行

```bash
peer channel join -b mych.block      #使用管理员身份加入通道
peer channel lsit                    #查看当前节点加入的通道
```

#### 更新锚节点

##### 清除非锚节点管理员环境变量

> 此操作需要在各组织的非锚节点(peer1.org1.example.com、peer1.org2.example.com)上执行

```bash
unset CORE_PEER_LOCALMSPID
unset CORE_PEER_MSPCONFIGPATH
```

##### 更新锚节点信息

更新各组织的锚节点信息, 使通道内所有组织的锚节点能够互相发现.
> 此操作需要在各组织的锚节点(peer0.org1.example.com、peer0.org2.example.com)上执行

```bash
cd /etc/hyperledger/fabric
peer channel update -o orderer.example.com:7050 -c mych -f Org1MSPanchors.tx --tls true --cafile /etc/hyperledger/fabric/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem      #两个锚节点均要更新,需要清除同一org内其他节点的管理员环境变量
```

##### 重设非锚节点管理员环境变量

> 此操作需要在各组织的非锚节点(peer1.org1.example.com、peer1.org2.example.com)上执行

```bash
export CORE_PEER_LOCALMSPID="Org1MSP"               #不同组织的peer节点设置相应的环境变量
export CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/fabric/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
```

### 链码安装与实例化

链码即是Fabric的智能合约, 以docker容器的方式运行, Fabric中任何与区块链(账本)进行的交互都必须通过链码完成, 无论是用户的区块链交易还是系统对区块链的管理, 用户交易使用的是定制化开发的链码, 而系统对区块链的管理则是使用的系统链码, 系统链码内置于Fabric各节点中.

#### 链码打包

本例使用Fabric源码自带的go语言链码chaincode_example02演示链码的使用.
> 此操作只需要在任意一个peer节点执行

```bash
cp -r /root/git/fabric/examples /root/go/src/
cd /root/go/src/
peer chaincode package -n mycc -p /root/go/src/examples/chaincode/go/chaincode_example02 -v 1.0 mycc.pak
cp mycc.pak /etc/hyperledger/fabric/                    #使用管理员身份打包链码
cd /etc/hyperledger/fabric/
scp mycc.pak  root@10.0.2.7:/etc/hyperledger/fabric/    #将打包好的链码分发到其他peer节点
scp mycc.pak  root@10.0.2.8:/etc/hyperledger/fabric/
scp mycc.pak  root@10.0.2.9:/etc/hyperledger/fabric/
```

#### 链码安装

> 此操作需要在所有peer节点执行

```bash
cd /etc/hyperledger/fabric/
peer chaincode install mycc.pak               #使用管理员身份安装链码
peer chaincode list  -C mych --installed      #查看安装的链码
```

#### 链码实例化

>此操作只需要在任意一个peer节点执行, 其他节点会自动同步该实例化的链码信息并创建链码容器, 实例化链码时可以指定背书策略(相当于权限控制)

```bash
peer chaincode instantiate -o orderer.example.com:7050 -C mych -n mycc -v 1.0 -c '{"Args":["init", "a", "100", "b", "200"]}' -P "OR('Org1MSP.member', 'Org2MSP.member')" --tls true --cafile /etc/hyperledger/fabric/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem     #使用管理员身份实例化链码
docker ps                #可以看到链码容器已启动(但只有发起链码实例化的节点的容器会启动,因为链码容器只有在节点发起交易的时候才启动)
peer chaincode list -C mych --instantiated    #可以在各节点上查看已经实例化的链码(确认链码在节点上是否实例化)
```

#### 交易测试

在不同节点上分别发起invoke交易和query交易, 查看Fabric是否正常工作

```bash
peer chaincode invoke -o orderer.example.com:7050 -C mych -n mycc -c '{"Args":["invoke", "a", "b", "10"]}' --tls true --cafile /etc/hyperledger/fabric/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem    #在另一个节点上发起一个交易,执行时间会稍长(因为需要启动容器)
peer chaincode query -C mych -n mycc -c '{"Args":["query", "a"]}'   #再在第三个节点上发起查询,查询也是交易的一种,仍然会启动容器,但查询交易并不需要其他节点进行同步
```

### 说明

* 全部配置完成的YAML配置文件下载地址:[链接: <https://pan.baidu.com/s/1EdTIROP0VbQ_L1u0KIjE4A> 密码: s2v7].
* 对于通道,链码和交易的相关操作会永久生效,节点重启后只需要启动服务即可(peer node start).
* 配置完成后, 请使用unset清除各peer节点的CORE_PEER_LOCALMSPID和CORE_PEER_MSPCONFIGPATH环境变量.
* 重启所有节点, 测试Fabric是否依然能够正常工作.
* 关于kafka模式的orderer集群, 链码的更新, 新组织的加入, Blockchain Explorer等内容, 待续...

------