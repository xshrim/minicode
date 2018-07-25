const fs = require('fs');
const Web3 = require('web3');
const solc = require('solc');
// const async = require('async');
// const subprocess = require('child_process');

// 注意: web3.js 1.0 API 不再支持同步调用. ethereum console支持同步方法, web3 API支持异步方法
// 文档: http://web3js.readthedocs.io/en/1.0/index.html
/* 安装所需node modules
npm install web3
npm install solc
*/

/* 运行一个开发者节点
var cmd = 'nohup geth --identity=myeth --networkid=111 --nodiscover --maxpeers=10 --port=30303 --dev --dev.period=0 --gasprice=20000000000 --targetgaslimit=4712388 --rpc --rpcapi=db,eth,net,web3,personal,miner,admin,debug --rpcaddr=0.0.0.0 --rpccorsdomain=* --rpcport=8545 --ws --wsapi=db,eth,net,web3,personal,miner,admin,debug --wsaddr=0.0.0.0 --wsorigins=* --wsport=8546 --datadir=/tmp & 2>/dev/null'
//subprocess.exec("rm -rf /home/xshrim/test/*")
//subprocess.exec(cmd);
//区块gas和API gas:
//区块gas的--gasprice参数表示区块所能接受的一个gas的最小价格(多少Wei, 默认0.05e12 Wei)
//区块gas的--targetgaslimit参数表示一个区块最多能打包多少个gas(区块打包的所有交易消耗的gas个数相加不能高于此值, 此值可以由矿工调整)
//API gas的gas参数表示一次交易允许消耗的gas个数的上限(如果交易实际需要消耗的gas数量高于此值, 交易会被回滚, 消耗的gas不会退回)
//API gas的gasPrice参数表示交易消耗的gas的单价(一个gas价格是多少Wei)
//综上, API中的gas参数值应当低于区块--targetgaslimit参数的值, 但不能太低, 否则会导致交易失败回滚; API中的gasPrice参数值应当稍微高于区块gas的gasprice参数的值, 否则区块链会拒绝将区块打包到区块中
*/

// var gownerAddress, gcontractName, gcontractAbi, gcontractCode, gcontractAddress;

function sleep(milliSeconds) {
    var startTime = new Date().getTime();
    while (new Date().getTime() < startTime + milliSeconds);
};

function display(ownerAddress, contractName, contractAbi, contractCode, contractAddress) {
    var jsonObj = {}
    console.log("================================================================================");
    if (ownerAddress !== 'undefined' && ownerAddress !== '') {
        jsonObj['ownerAddress'] = ownerAddress;
        console.log("------------------------------- Owner Address Begin ----------------------------");
        console.log(ownerAddress);
        console.log("-------------------------------- Owner Address End -----------------------------");
    }
    if (contractName !== 'undefined' && contractName !== '') {
        jsonObj['contractName'] = contractName;
        console.log("------------------------------- Contract Name Begin ----------------------------");
        console.log(contractName);
        console.log("-------------------------------- Contract Name End -----------------------------");
    }
    if (contractAbi !== 'undefined' && contractAbi !== '') {
        jsonObj['contractAbi'] = contractAbi;
        console.log("------------------------------- Contract ABI Begin -----------------------------");
        console.log(contractAbi);
        console.log("-------------------------------- Contract ABI End ------------------------------");
    }
    if (contractCode !== 'undefined' && contractCode !== '') {
        jsonObj['contractCode'] = contractCode;
        console.log("------------------------------- Contract Code Begin ----------------------------");
        console.log(contractCode);
        console.log("-------------------------------- Contract Code End -----------------------------");
    }
    if (contractAddress !== 'undefined' && contractAddress !== '') {
        jsonObj['contractAddress'] = contractAddress;
        console.log("------------------------------ Contract Address Begin --------------------------");
        console.log(contractAddress);
        console.log("------------------------------- Contract Address End ---------------------------");
    }
    console.log("================================================================================");

    fs.writeFile('./contractInfo.json', JSON.stringify(jsonObj), function (err) {
        console.log(err);
    });
}

// 初始化web3实例
function init(provider) {
    var web3 = new Web3(Web3.givenProvider || provider);
    if (web3.currentProvider.connection._readyState !== 0) {
        console.log('<web3 connection failed!>');
        // console.log(web3.currentProvider.connection);
        process.exit();
    } else {
        console.log('<web3 connection success!>');
        // console.log(web3.currentProvider.connection);
    }
    return web3;
}

/* 变异步为同步
const fs = require('fs/promises');

const readConfig = async (configFile) => {
  try {
      const config = await fs.readFile(configFile, 'utf-8');
      // do something with file contents
  } catch (e) {
    throw new Error(`Failed to load config file at ${configFile}`);
  }
} 
*/

// 编译solidity源码
function compile(filepath) {
    var arr = new Array();
    var compiledContract = solc.compile(fs.readFileSync(filepath, "utf8"));
    if (JSON.stringify(compiledContract).indexOf("errors") >= 0) {
        console.log(compiledContract);
        if (JSON.stringify(compiledContract).indexOf("Error") >= 0) {
            process.exit();
        }
    }
    for (var contractName in compiledContract.contracts) {
        var contractCode = "0x" + compiledContract.contracts[contractName].bytecode;
        var contractAbi = JSON.parse(compiledContract.contracts[contractName].interface);
        //var contractAbi = compiledContract.contracts[contractName].interface;
        arr.push([contractName, contractAbi, contractCode]);
        //var gasEstimate = web.eth.estimateGas({
        //    data: '0x' + bytecode
        //});
        //var con = new web3.eth.Contract(abi);
        //console.log(con.options);
        //console.log(gasEstimate);

    }
    /* eth_compile is deprecated, use solc instead
    web3.eth.compile.solidity(source).then((compiled) => {
        console.log(compile);
        for (var contractName in compiled) {
            var abi = compiled[contractName].info.abiDefinition;
            var bcd = compiled[contractName].code;
        }
        console.log(abi);
    });
    */
    return arr;
}

// 向节点部署合约
function deploy(web3, ownerAddress, contractName, contractAbi, contractCode, contractArgs, contractGas, contractGasPrice) {
    // console.log("deploy contract " + contractName + "...");
    return new web3.eth.Contract(contractAbi).deploy({
            data: contractCode,
            arguments: contractArgs
        }).send({
            from: ownerAddress,
            gas: contractGas,
            gasPrice: contractGasPrice,
        })
        .on('error', function (error) {
            console.log("error: " + error);
            console.log("deploy contract failed");
        })
        .on('transactionHash', function (transactionHash) {
            console.log("transactionHash: " + transactionHash);
        })
        .on('receipt', function (receipt) {
            console.log("receipt: " + JSON.stringify(receipt));
        })
        .on('confirmation', function (confirmationNumber, receipt) {
            console.log("confirmation: " + confirmationNumber);
            // console.log("confirmation: " + JSON.stringify(receipt));
        })
        .then((instance) => {
            var cownerAddress = ownerAddress;
            var ccontractName = contractName;
            var ccontractAbi = JSON.stringify(instance.options.jsonInterface);
            var ccontractCode = contractCode;
            var ccontractAddress = instance.options.address;
            display(cownerAddress, ccontractName, ccontractAbi, ccontractCode, ccontractAddress);
            // console.log("deploy contract " + contractName + " done");
            return [cownerAddress, ccontractName, ccontractAbi, ccontractCode, ccontractAddress];
        });
}

// 获取合约函数集
function methods(web3, contractAbi, contractAddress) {
    return new web3.eth.Contract(JSON.parse(contractAbi), contractAddress).methods;
}

// 调用合约示例
function invoke(web3, ownerAddress, contractAbi, contractAddress, contractMethods, contractArgs, contractGas, contractGasPrice) {

    var invokeContract = new web3.eth.Contract(JSON.parse(contractAbi), contractAddress)
    // var cmethods = methods(web3, contractAbi, contractAddress);

    // 事件保存到区块链成为日志, 由于只有send调用才会写入区块(call则不会), 因此只有在send相关的函数中发起的事件才会成为日志记入区块链
    var subscription = web3.eth.subscribe('logs', { // 订阅日志, 订阅只能看到事件的日志ID
            address: contractAddress,
        }, function (error, result) {
            if (!error)
                console.log(result);
            else
                console.log(error);
        })
        .on("data", function (log) {
            console.log(log);
        })
        .on("changed", function (log) {});

    invokeContract.events.WatchDog({ // 事件监控(支持多种写法), 能够看到事件日志的详细信息, 可自定义监控范围(指定区块, 事件名等)
        fromBlock: 0,
        toBlock: 'latest'
    }, (error, data) => {
        if (error)
            console.log("Error: " + error);
        else
            console.log("Log data: " + web3.utils.hexToAscii(data["returnValues"]["themsg"]));
    });

    /* 监控指定事件一次
    invokeContract.once('WatchDog', {
        fromBlock: 0,
        toBlock: 'latest'
    }, (error, data) => {
        if (error)
            console.log("Error: " + error);
        else
            console.log("Log data: " + web3.utils.hexToAscii(data["returnValues"]["themsg"]));
    });
    */

    /* 监控全部事件
    invokeContract.events.allEvents({
        fromBlock: 0,
        toBlock: 'latest'
    }, (error, data) => {
        if (error)
            console.log("Error: " + error);
        else
            console.log("Log data: " + web3.utils.hexToAscii(data["returnValues"]["themsg"]));
    });
    */

    /* 获取以往事件
    myContract.getPastEvents('WatchDog', {
        filter: {myIndexedParam: [20,23], myOtherIndexedParam: '0x123456789...'}, // Using an array means OR: e.g. 20 or 23
        fromBlock: 0,
        toBlock: 'latest'
    }, function(error, events){ console.log(events); })
    .then(function(events){
        console.log(events) // same results as the optional callback above
    });
     */

    var options = {
        from: ownerAddress,
        gas: contractGas,
        gasPrice: contractGasPrice
    }

    const message = web3.utils.asciiToHex(contractArgs[0]);
    /* 合约函数调用的三种形式
    The name: myContract.methods.myMethod(123)
    The name with parameters: myContract.methods['myMethod(uint256)'](123)
    The signature: myContract.methods['0x58cf5f10'](123)
    */

    return invokeContract.methods[contractMethods[0]](message)
        .send(options, (err, hash) => {
            if (err) {
                console.log(err);
            }
            console.log(`TxHash: ${hash}`);
        })
        .on('receipt', function (receipt) {
            // 监控回单(包括事件), 也可以得到事件日志的详细信息, 还能得到很多事件之外的区块信息
            console.log('Receipt: ' + JSON.stringify(receipt));
        })
        .then((result) => {
            // 取消订阅
            subscription.unsubscribe(function (error, success) {
                if (success)
                    console.log('Successfully unsubscribed!');
            });

            return invokeContract.methods[contractMethods[1]]().call()
                .then((result) => {
                    return web3.utils.hexToAscii(result);
                });
        });
}

// 运行合约部署与调用主函数
function main(provider, contractPath, deployFlag, invokeFlag) {
    var contractsInfo = [];
    console.log("<compile contracts ...>");
    var contracts = compile(contractPath);
    console.log("<compile contracts done>");

    if (!deployFlag) {
        contracts.forEach(contract => {
            display('', contract[0], JSON.stringify(contract[1]), contract[2], '');
        });
    } else {
        var web3 = init(provider);
        console.log("<get account address ...>");
        web3.eth.getAccounts()
            .then((accounts) => {
                // console.log("ownerAddress: " + accounts[0]);
                console.log("<deploy contracts ...>");

                contracts.forEach(contract => {
                    deploy(web3, accounts[0], contract[0], contract[1], contract[2], [10000000], 2000000, '22000000000')
                        .then((result) => {
                            console.log("<deploy contracts done>");
                            if (!invokeFlag) {
                                process.exit();
                            } else {
                                console.log("<invoke contract ...>");
                                invoke(web3, result[0], result[2], result[4], ['sendMessage(bytes32)', 'getMessage()'], ['hi'], 2000000, '22000000000')
                                    .then((result) => {
                                        console.log("invoke result: " + result);
                                        console.log("<invoke contract done>");
                                        process.exit();
                                    });
                            }
                        });
                });
                /* 变异步为同步
                async.series(tasks, function (error, result) {
                    if (error) {
                        console.log(error);
                        process.exit();
                    }
                    contractsInfo.push(result);
                    console.log(result);
                });
                */
            });
    }
}

/*
var flag = 0;
var spawn = require('child_process').spawn;
free = spawn('geth', ['--identity=myeth', '--networkid=111', '--nodiscover', '--maxpeers=10', '--port=30303', '--dev', '--dev.period=0', '--gasprice=20000000000', '--targetgaslimit=4712388', '--rpc', '--rpcapi=db,eth,net,web3,personal,miner,admin,debug', '--rpcaddr=0.0.0.0', '--rpccorsdomain=*', '--rpcport=8545', '--ws', '--wsapi=db,eth,net,web3,personal,miner,admin,debug', '--wsaddr=0.0.0.0', '--wsorigins=*', '--wsport=8546', '--datadir=/tmp']);
// 捕获标准输出并将其打印到控制台 
free.stdout.on('data', function (data) {
    console.log('standard output:\n' + data);
});
// 捕获标准错误输出并将其打印到控制台 
free.stderr.on('data', function (data) {
    console.log('standard error output:\n' + data);
    if (flag == 0 && data.indexOf('waiting for transactions') > 0) {
        flag = 1;
        console.log('GO');
        var deployFlag = true; // 是否部署合约
        var invokeFlag = false; // 是否进行函数调用测试
        var contractPath = './Echo.sol';
        var provider = "ws://localhost:8546"

        main(provider, contractPath, deployFlag, invokeFlag);
    }
});
// 注册子进程关闭事件 
free.on('exit', function (code, signal) {
    console.log('child process eixt ,exit:' + code);
});
*/

// var provider = "ws://localhost:7545"

var deployFlag = true; // 是否部署合约
var invokeFlag = true; // 是否进行函数调用测试
var contractPath = './Test.sol';
var provider = "ws://localhost:8546"

main(provider, contractPath, deployFlag, invokeFlag);