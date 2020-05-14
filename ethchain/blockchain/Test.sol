pragma solidity ^0.4.22;


contract Test {
    mapping(address => uint)public balances;
    bytes32 public message;

    event WatchDog(address sender, address receiver, bytes32 themsg, uint number);

    constructor(uint amount) public {
        balances[msg.sender] = amount;
    }

    function getBalance() public view returns (uint) {
        return balances[msg.sender];
    }

    function sendMessage(bytes32 themsg) public {
        message = themsg;
        emit WatchDog(msg.sender, msg.sender, themsg, 1);
    }

    function getMessage() public view returns (bytes32) {
        return message;
    }

    function sendBalance(address receiver, uint amount) public {
        if (balances[msg.sender] < amount) {
            emit WatchDog(msg.sender, receiver, "Failure", amount);
            return;
        }
        balances[msg.sender] -= amount;
        balances[receiver] += amount;
        emit WatchDog(msg.sender, receiver, "Success", amount);
    }
}