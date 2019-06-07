#!/usr/bin/bash

rm -rf ./code/bin/*
solc --bin --abi --optimize --overwrite -o ./code/bin ./code/Test.sol
if [[ ! -f ./code/bin/Test.abi ]] || [[ ! -f ./code/bin/Test.bin ]];then
    echo 'compile failed!'
    exit
fi

abi=$(cat ./code/bin/Test.abi)
bin="0x"$(cat ./code/bin/Test.bin)

echo $abi
echo "================================="
echo $bin

geth --identity "myeth" --networkid 111 --nodiscover --maxpeers 10 --port "30303" --mine --rpc --rpcapi "db,eth,net,web3,personal,miner,admin,debug" --rpcaddr 0.0.0.0 --rpcport "8545" --wsport "8546" --datadir "./ethbase" console 2> geth.log << EOF
personal.unlockAccount(web3.eth.accounts[0], 'bcfan', 10000)
eth.defaultAccount=eth.coinbase
abi = $abi
testContract = web3.eth.contract(abi)
test = testContract.new({from: web3.eth.accounts[0], data: "$bin", gas: "300000"})
miner.start(2);admin.sleepBlocks(1);miner.stop()
test.address
MyContract = web3.eth.contract(abi)
myContract = MyContract.at(test.address)
miner.start(2)
myContract.invoke("multiply", 10)
miner.stop()
EOF
