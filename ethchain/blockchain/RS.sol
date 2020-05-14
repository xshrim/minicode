pragma solidity ^0.4.23;


/*@ 文件描述
//@filename: RS.sol
//@author: 陈盼
//@date: 2018-07-13
//@project: 滚石
//@describe: 滚石资源权限管理智能合约
*/


/*@ 滚石合约
*/
contract RS {

    /*@ 用户结构体
    */
    struct User {
        address addr;           // 用户地址
        bytes32 id;             // 用户id
        bytes32 itype;          // 用户类型
        bytes32 name;           // 用户名
        bytes32 hash;           // 用户信息hash
        // bytes32 info;        // 用户主要信息
        uint64 level;           // 用户级别(0: 地址不存在, 1: 管理员, 2: 普通用户)
        uint64 timestamp;       // 用户创建时间
        uint64 status;          // 用户状态(0: 不存在, 1: 正常, 2：冻结， 6: 不可用)
    }

    /*@ 资源结构体
    */
    struct Item {
        bytes32 id;             // 资源id
        bytes32 userid;         // 资源所有者id
        bytes32 itype;          // 资源类型
        bytes32 title;          // 资源名称
        bytes32 hash;           // 资源信息hash
        // bytes32 info;        // 资源主要信息
        bytes32 path;           // 资源路径
        uint64 level;           // 资源保密级别
        uint64 timestamp;       // 资源生成时间戳
        uint64 status;          // 资源状态(0: 不存在, 1: 正常, 2：冻结， 6: 不可用)
    }

    /*@ 权限结构体
    */
    struct Perm {
        bytes32 id;             // 授权权限id
        bytes32 signer;         // 授权者id
        bytes32 userid;         // 授权用户id
        bytes32 itemid;         // 授权资源id
        bytes32 device;         // 授权查看设备
        uint64[3] ptype;        // 授权类型([是否可查看, 是否可下载, 是否可向下授权])
        uint64[2] ptimestamp;   // 授权时间戳([生成时间戳, 过期时间戳])
        uint64[2] ptime;        // 授权查看时间([总时间, 剩余时间])
        uint64[2] ptimes;       // 授权查看次数([总次数, 剩余次数])
        uint64 status;          // 授权权限状态(0: 不存在, 1: 正常, 2：冻结， 6: 不可用)
    }

    /*@ 日志结构体
    */
    struct Log {
        bytes32 id;             // 操作日志id
        bytes32 userid;         // 操作用户id
        bytes32 itemid;         // 操作资源id
        bytes32 action;         // 操作行为
        bytes32 hash;           // 操作hash
        bytes32 desc;           // 操作描述
        uint64 duration;        // 操作时长
        uint64 timestamp;       // 操作时间戳
    }

    
    // uint networkid = 111;

    address public owner;

    uint64 usernum;
    mapping (uint64 => bytes32) uidxids;            // user index => user id
    mapping (bytes32 => address) uidaddrs;          // user id => user address
    mapping (address => bytes32) uaddrids;          // user address => user id
    // mapping (uint64 => bytes32) idxuids;         // user index => user id
    mapping (bytes32 => User) users;                // user id => User
    // mapping (bytes32 => User) iusers;            // user id => User


    uint64 itemnum;
    mapping (uint64 => bytes32) iidxids;            // item index => item id
    mapping (bytes32 => bytes32) iiduids;           // item id => user id
    mapping (bytes32 => Item) items;                // item id => Item

    uint64 permnum;
    mapping (uint64 => bytes32) pidxids;            // perm index => perm id
    mapping (bytes32 => bytes32[2]) piduiids;       // perm id => [userid, itemid]
    mapping (bytes32 => bytes32) puiidids;          // keccak256([userid, itemid]) => perm id
    mapping (bytes32 => Perm) perms;                // perm id => Perm

    uint64 lognum;
    mapping (uint64 => bytes32) lidxids;            // log index => log id
    mapping (bytes32 => bytes32[2]) liduiids;       // log id => [userid, itemid]
    mapping (bytes32 => Log) logs;                  // log id => Log


    /*@ 测试事件
    //@param themsg: bytes32类型日志字段
    //@param thimsg: string类型日志字段
    //@param res: uint类型日志字段
    */
    event WatchDog (bytes32 themsg, string thimsg, uint indexed res);

    /*@ 交易凭据日志
    //@param sn: 交易流水号
    //@param sender: 交易发起者地址
    //@param userid: 交易发起者id
    //@param itemid: 交易资源id
    //@param cid: 交易授权/日志id
    //@param func: 交易调用函数
    //@param details: 交易详细信息
    //@param timestamp: 交易时间戳
    */
    event TXReceipt (bytes32 indexed sn, address sender, bytes32 indexed userid, bytes32 indexed itemid, bytes32 cid, bytes32 func, string details, uint64 timestamp);


    constructor() public payable{
        owner = msg.sender;
        uidxids[usernum++] = "0";
        uidaddrs["0"] = msg.sender;
        uaddrids[msg.sender] = "0";
        users["0"] = User(msg.sender, "0", "admin", "root", "owner", 1, uint64(block.timestamp), 1);
        owner.transfer(1000000000000000000000);
        // tkids.push(0x00);
        // tokens[0x00] = Token(0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0, 0, 0, 0, 0, 0, 0, 2);
    }

    // 修饰器
    modifier onlyOwner() {
        require (msg.sender == owner, "sender != owner");
        /*
        if (msg.sender != owner) {
            revert("sender != owner");
        }
        */
        _;
    }

    // 修饰器
    modifier onlyAdmin()
    {
        uint64[2] memory ls = getUserLSByAddr(msg.sender);
        require(ls[0] == 1 && ls[1] == 1, "sender level != 1 or status != 1");
        _;
    }

    // 修饰器
    modifier onlyValid()
    {
        uint64[2] memory ls = getUserLSByAddr(msg.sender);
        require(ls[1] == 1, "sender status != 1");
        _;
    }


    /*@ 将两个bytes32变量进行hash
    //@param x: 变量x
    //@param y: 变量y
    //@return 0: hash值
    */
    function concat(bytes32 x, bytes32 y) internal pure returns (bytes32) {
        /*
        bytes memory s = new bytes(64);
        for (uint i = 0; i < 32; i++){
            s[i] = x[i];
            s[i+32] = y[i];
        }
        return keccak256(s);
        */
        /*
        assembly {
            result := mload(add(s, 32))
        }
        */

        return keccak256(abi.encodePacked(x, y));
        
    }

    /*@ 判断用户是否是管理员By用户地址
    //@return 0: 是否是管理员
    */
    function isAdmin() internal view returns (bool) {
        uint64[2] memory ls = getUserLSByAddr(msg.sender);
        if (ls[0] == 1) {
            return true;
        } else {
            return false;
        }
    }

    /*@ 设置用户
    //@param addr: 用户地址
    //@param id: 用户id
    //@param itype: 用户类型
    //@param name: 用户名
    //@param hash: 用户hash
    //@param level: 用户级别
    //@param status: 用户状态
    //@param timestamp: 生成时间戳
    //@return 0: 是否设置成功
    */
    function setUser(address addr, bytes32 id, bytes32 itype, bytes32 name, bytes32 hash, uint64 level, uint64 status, uint64 timestamp) internal returns (bool) {
        uidxids[usernum++] = id;
        uidaddrs[id] = addr;
        uaddrids[addr] = id;
        users[id] = User(addr, id, itype, name, hash, level, timestamp, status);
        return true;
    }

    /*@ 添加用户
    //@param sn: 交易流水号
    //@param details: 交易详细信息
    //@param addr: 用户地址
    //@param id: 用户id
    //@param itype: 用户类型
    //@param name: 用户名
    //@param hash: 用户hash
    //@param level: 用户级别
    //@param status: 用户状态
    //@return 0: 是否设置成功
    */
    function addUser(bytes32 sn, string details, address addr, bytes32 id, bytes32 itype, bytes32 name, bytes32 hash, uint64 level, uint64 status) public payable onlyAdmin returns (bool) {
        if (setUser(addr, id, itype, name, hash, level, status, uint64(block.timestamp))) {
            addr.transfer(1000000000000000000000);
            emit TXReceipt(sn, msg.sender, id, 0x00, 0x00, "addUser", details, uint64(block.timestamp));
            // bytes32 lid = keccak256(block.difficulty,now);
            setLog([sn, uaddrids[msg.sender], 0x00], "addUser", 0x00, 0x00, 0, uint64(block.timestamp));
            return true;
        } else {
            revert();
        }
    }

    /*@ 设置用户状态
    //@param id: 用户id
    //@param newstatus: 用户新的状态
    //@return 0: 是否设置成功
    */
    function setUserStatus(bytes32 id, uint64 newstatus) internal returns (bool) {
        if (id == 0x00) {
            revert();
        } else {
            User storage user = users[id];
            user.status = newstatus;
            return true;
        }
    }

    /*@ 设置用户状态By地址
    //@param sn: 交易流水号
    //@param details: 交易详细信息
    //@param addr: 用户地址
    //@param newstatus: 用户新的状态
    //@return 0: 是否设置成功
    */
    function setUserStatusByAddr(bytes32 sn, string details, address addr, uint64 newstatus) public onlyAdmin returns (bool) {
        if (setUserStatus(uaddrids[addr], newstatus)) {
            emit TXReceipt(sn, msg.sender, uaddrids[addr], 0x00, 0x00, "setUserStatusByAddr", details, uint64(block.timestamp));
            setLog([sn, uaddrids[msg.sender], 0x00], "setUserStatusByAddr", 0x00, 0x00, 0, uint64(block.timestamp));
            return true;
        } else {
            revert();
        }
    }

    /*@ 设置用户状态By用户id
    //@param sn: 交易流水号
    //@param details: 交易详细信息
    //@param id: 用户id
    //@param newstatus: 用户新的状态
    //@return 0: 是否设置成功
    */
    function setUserStatusById(bytes32 sn, string details, bytes32 id, uint64 newstatus) public onlyAdmin returns (bool) {
        if (setUserStatus(id, newstatus)) {
            emit TXReceipt(sn, msg.sender, id, 0x00, 0x00, "setUserStatusById", details, uint64(block.timestamp));
            setLog([sn, uaddrids[msg.sender], 0x00], "setUserStatusById", 0x00, 0x00, 0, uint64(block.timestamp));
            return true;
        } else {
            revert();
        }
    }

    /*@ 设置用户状态By索引
    //@param sn: 交易流水号
    //@param details: 交易详细信息
    //@param idx: 用户索引
    //@param newstatus: 用户新的状态
    //@return 0: 是否设置成功
    */
    function setUserStatusByIdx(bytes32 sn, string details, uint64 idx, uint64 newstatus) public onlyAdmin returns (bool) {
        if (setUserStatus(uidxids[idx], newstatus)) {
            emit TXReceipt(sn, msg.sender, uidxids[idx], 0x00, 0x00, "setUserStatusByIdx", details, uint64(block.timestamp));
            setLog([sn, uaddrids[msg.sender], 0x00], "setUserStatusByIdx", 0x00, 0x00, 0, uint64(block.timestamp));
            return true;
        } else {
            revert();
        }
    }

    /*@ 获取用户总数
    //@return 0: 用户总数
    */
    function getUserNum() public view returns (uint64) {
        return usernum;
    }

    /*@ 获取用户信息By用户id
    //@param id: 用户id
    //@return 0: 用户地址
    //@return 1: 用户id
    //@return 2: 用户类型
    //@return 3: 用户名
    //@return 4: 用户hash
    //@return 5: 用户级别
    //@return 6: 生成时间戳
    //@return 7: 用户状态
    */
    function getUserById(bytes32 id) public view returns (address, bytes32, bytes32, bytes32, bytes32, uint64, uint64, uint64) {
        if (id == 0x00) {
            return (0x00, bytes32(0x00), 0x00, 0x00, 0x00, uint64(0), 0, 0);
        } else {
            User memory user = users[id];
            return (user.addr, user.id, user.itype, user.name, user.hash, user.level, user.timestamp, user.status);
        }
    }

    /*@ 获取用户信息By用户地址
    //@param addr: 用户地址
    //@return 0: 用户地址
    //@return 1: 用户id
    //@return 2: 用户类型
    //@return 3: 用户名
    //@return 4: 用户hash
    //@return 5: 用户级别
    //@return 6: 生成时间戳
    //@return 7: 用户状态
    */
    function getUserByAddr(address addr) public view returns (address, bytes32, bytes32, bytes32, bytes32, uint64, uint64, uint64) {
        return getUserById(uaddrids[addr]);
    }

    /*@ 获取用户信息By用户索引
    //@param idx: 用户索引
    //@return 0: 用户地址
    //@return 1: 用户id
    //@return 2: 用户类型
    //@return 3: 用户名
    //@return 4: 用户hash
    //@return 5: 用户级别
    //@return 6: 生成时间戳
    //@return 7: 用户状态
    */
    function getUserByIdx(uint64 idx) public view returns (address, bytes32, bytes32, bytes32, bytes32, uint64, uint64, uint64) {
        return getUserById(uidxids[idx]);
    }

    /*@ 获取用户等级状态 
    //@param id: 用户id
    //@return 0: [用户等级, 用户状态]
    */
    function getUserLS(bytes32 id) public view returns (uint64[2]) {
        if (id == 0x00) {
            return ([uint64(0), 0]);
        } else {
            User memory user = users[id];
            return [user.level, user.status];
        }
    }

    /*@ 获取用户等级状态By用户地址 
    //@param addr: 用户地址
    //@return 0: [用户等级, 用户状态]
    */
    function getUserLSByAddr(address addr) public view returns (uint64[2]) {
        return getUserLS(uaddrids[addr]);
    }
    

    /*@ 设置资源
    //@param ids: [资源id, 资源所有者id]
    //@param itype: 资源类型
    //@param title: 资源标题
    //@param hash: 资源hash
    //@param path: 资源路径
    //@param level: 资源安全级别
    //@param status: 资源状态
    //@param timestamp: 生成时间戳
    //@return 0: 是否设置成功
    */
    function setItem(bytes32[2] ids, bytes32 itype, bytes32 title, bytes32 hash, bytes32 path, uint64 level, uint64 timestamp, uint64 status) internal returns (bool) {
        iidxids[itemnum++] = ids[0];
        items[ids[0]] = Item(ids[0], ids[1], itype, title, hash, path, level, timestamp, status);
        iiduids[ids[0]] = ids[1];
        return true;
    }

    /*@ 添加资源
    //@param sn: 交易流水号
    //@param details: 交易详细信息
    //@param ids: [资源id, 资源所有者id]
    //@param itype: 资源类型
    //@param title: 资源标题
    //@param hash: 资源hash
    //@param path: 资源路径
    //@param level: 资源安全级别
    //@param status: 资源状态
    //@return 0: 是否设置成功
    */
    function addItem(bytes32 sn, string details, bytes32[] ids, bytes32 itype, bytes32 title, bytes32 hash, bytes32 path, uint64 level, uint64 status) public onlyAdmin returns (bool) {
        //bytes32 senderid = uaddrids[msg.sender];
        //if (ids.length == 2) {
        //    senderid = ids[1];
        //}
        if (ids.length == 2) {
            if (setItem([ids[0], ids[1]], itype, title, hash, path, level, uint64(block.timestamp), status)) {
                emit TXReceipt(sn, msg.sender, ids[1], ids[0], 0x00, "addItem", details, uint64(block.timestamp));
                setLog([sn, ids[1], ids[0]], "addItem", 0x00, 0x00, 0, uint64(block.timestamp));
                return true;
            } else {
                revert();
            }
        } else {
            if (setItem([ids[0], uaddrids[msg.sender]], itype, title, hash, path, level, uint64(block.timestamp), status)) {
                emit TXReceipt(sn, msg.sender, uaddrids[msg.sender], ids[0], 0x00, "addItem", details, uint64(block.timestamp));
                setLog([sn, uaddrids[msg.sender], ids[0]], "addItem", 0x00, 0x00, 0, uint64(block.timestamp));
                return true;
            } else {
                revert();
            }
        }

    }

    /*@ 上传资源
    //@param sn: 交易流水号
    //@param details: 交易详细信息
    //@param itemid: 资源id
    //@param itype: 资源类型
    //@param title: 资源标题
    //@param hash: 资源hash
    //@param path: 资源路径
    //@param level: 资源安全级别
    //@param status: 资源状态
    //@return 0: 是否设置成功
    */
    function uplItem(bytes32 sn, string details, bytes32 itemid, bytes32 itype, bytes32 title, bytes32 hash, bytes32 path, uint64 level, uint64 status) public onlyValid returns (bool) {
        //bytes32 senderid = uaddrids[msg.sender];
        if (setItem([itemid, uaddrids[msg.sender]], itype, title, hash, path, level, uint64(block.timestamp), status)) {
            emit TXReceipt(sn, msg.sender, uaddrids[msg.sender], itemid, 0x00, "uplItem", details, uint64(block.timestamp));
            setLog([sn, uaddrids[msg.sender], 0x00], "uplItem", 0x00, 0x00, 0, uint64(block.timestamp));
            return true;
        } else {
            revert();
        }
    }

    /*@ 设置资源状态
    //@param itemid: 资源id
    //@param newstatus: 资源新的状态
    //@return 0: 是否设置成功
    */
    function setItemStatus(bytes32 itemid, uint64 newstatus) internal returns (bool) {
        if (itemid == 0x00) {
            revert();
        } else {
            bytes32 senderid = uaddrids[msg.sender];
            User memory user = users[senderid];
            Item storage item = items[itemid];
            if (user.status == 0) {
                revert();
            } else {
                if (user.level == 1 || item.userid == senderid) {
                    item.status = newstatus;
                    return true;
                } else {
                    revert();
                }
            }
        }
    }

    /*@ 设置资源状态By资源id
    //@param sn: 交易流水号
    //@param details: 交易详细信息
    //@param id: 资源id
    //@param newstatus: 资源新的状态
    //@return 0: 是否设置成功
    */
    function setItemStatusById(bytes32 sn, string details, bytes32 id, uint64 newstatus) public onlyValid returns (bool) {
        if (setItemStatus(id, newstatus)) {
            emit TXReceipt(sn, msg.sender, 0x00, id, 0x00, "setItemStatusById", details, uint64(block.timestamp));
            setLog([sn, uaddrids[msg.sender], id], "setItemStatusById", 0x00, 0x00, 0, uint64(block.timestamp));
            return true;
        } else {
            revert();
        }
    }

    /*@ 设置资源状态By资源索引
    //@param sn: 交易流水号
    //@param details: 交易详细信息
    //@param idx: 资源索引
    //@param newstatus: 资源新的状态
    //@return 0: 是否设置成功
    */
    function setItemStatusByIdx(bytes32 sn, string details, uint64 idx, uint64 newstatus) public onlyValid returns (bool) {
        if (setItemStatus(iidxids[idx], newstatus)) {
            emit TXReceipt(sn, msg.sender, 0x00, iidxids[idx], 0x00, "setItemStatusByIdx", details, uint64(block.timestamp));
            setLog([sn, uaddrids[msg.sender], iidxids[idx]], "setItemStatusByIdx", 0x00, 0x00, 0, uint64(block.timestamp));
            return true;
        } else {
            revert();
        }
    }

    /*@ 获取资源总数
    //@return 0: 资源总数
    */
    function getItemNum() public view returns (uint64) {
        return itemnum;
    }

    /*@ 获取资源信息By资源id
    //@param id: 资源id
    //@return 0: [资源id, 资源所有者id]
    //@return 1: 资源类型
    //@return 2: 资源标题
    //@return 3: 资源hash
    //@return 4: 资源路径
    //@return 5: 资源安全级别
    //@return 6: 资源生成时间戳
    //@return 7: 资源状态
    */
    function getItemById(bytes32 id) public view returns (bytes32[2], bytes32, bytes32, bytes32, bytes32, uint64, uint64, uint64) {
        if (id == 0x00) {
            return ([bytes32(0x00), 0x00], bytes32(0x00), 0x00, 0x00, 0x00, uint64(0), 0, 0);
        } else {
            Item memory item = items[id];
            return ([item.id, item.userid], item.itype, item.title, item.hash, item.path, item.level, item.timestamp, item.status);
        }
    }

    /*@ 获取资源信息By资源索引
    //@param idx: 资源索引
    //@return 0: [资源id, 资源所有者id]
    //@return 1: 资源类型
    //@return 2: 资源标题
    //@return 3: 资源hash
    //@return 4: 资源路径
    //@return 5: 资源安全级别
    //@return 6: 资源生成时间戳
    //@return 7: 资源状态
    */
    function getItemByIdx(uint64 idx) public view returns (bytes32[2], bytes32, bytes32, bytes32, bytes32, uint64, uint64, uint64) {
        return getItemById(iidxids[idx]);
    }

    /*@ 获取资源等级状态
    //@param id: 资源id
    //@return 0: [资源安全级别, 资源状态]
    */
    function getItemLS(bytes32 id) public view returns (uint64[2]) {
        if (id == 0x00) {
            return ([uint64(0), 0]);
        } else {
            Item memory item = items[id];
            return [item.level, item.status];
        }
    }
    

    /*@ 设置权限
    //@param ids: [权限id, 授予用户id, 资源id]
    //@param device: 设备限制
    //@param ptype: [查看权限, 下载权限, 授予权限]
    //@param ptimestamp: [权限生成时间戳, 权限过期时间戳]
    //@param ptime: [权限总时间, 权限剩余时间]
    //@param ptimes: [权限可用次数, 权限剩余可用次数]
    //@param status: 权限状态
    //@return 0: 设置是否成功
    */
    function setPerm(bytes32[3] ids, bytes32 device, uint64[3] ptype, uint64[2] ptimestamp, uint64[2] ptime, uint64[2] ptimes, uint64 status) internal returns (bool) {
        // TODO user.level >= item.level
        bytes32 senderid = uaddrids[msg.sender];
        User memory user = users[ids[1]];
        Item memory item = items[ids[2]];
        if (item.userid == ids[1] || user.level == 1) {
            revert();
        }
        if (verPerm(senderid, ids[2], ptype, ptimestamp[1], ptime[0], ptimes[0])) {
            pidxids[permnum++] = ids[0];
            perms[ids[0]] = Perm(ids[0], senderid, ids[1], ids[2], device, [ptype[0], ptype[1], ptype[2]], [ptimestamp[0], ptimestamp[1]], [ptime[0], ptime[1]], [ptimes[0], ptimes[1]], status);
            piduiids[ids[0]] = [ids[1], ids[2]];
            puiidids[concat(ids[1], ids[2])] = ids[0];

            user = users[senderid];
            if (item.userid != senderid && user.level != 1) {
                if (setPermTimes(puiidids[concat(senderid, ids[2])], ptimes[0])) {
                    return true;
                } else {
                    revert();
                }
            } else {
                return true;
            }
        } else {
            revert();
        }
    }


    /*@ 添加权限
    //@param sn: 交易流水号
    //@param details: 交易详细信息
    //@param ids: [权限id, 授予用户id, 资源id]
    //@param device: 设备限制
    //@param ptype: [查看权限, 下载权限, 授予权限]
    //@param ptimestamp: [权限生成时间戳, 权限过期时间戳]
    //@param ptime: [权限总时间, 权限剩余时间]
    //@param ptimes: [权限可用次数, 权限剩余可用次数]
    //@param status: 权限状态
    //@return 0: 设置是否成功
    */
    function addPerm(bytes32 sn, string details, bytes32[3] ids, bytes32 device, uint64[3] ptype, uint64[2] ptimestamp, uint64[2] ptime, uint64[2] ptimes, uint64 status) public onlyValid returns (bool) {
        if (puiidids[concat(ids[1], ids[2])] == 0x00) {
            if (setPerm(ids, device, ptype, ptimestamp, ptime, ptimes, status)) {
                emit TXReceipt(sn, msg.sender, ids[1], ids[2], ids[0], "addPerm", details, ptimestamp[0]);
                setLog([sn, uaddrids[msg.sender], ids[2]], "addPerm", 0x00, 0x00, 0, uint64(block.timestamp));
                return true;
            } else {
                revert();
            }
        } else {
            if (apdPerm(puiidids[concat(ids[1], ids[2])], device, ptype, ptimestamp[1], ptime[1], ptimes[1])) {
                emit TXReceipt(sn, msg.sender, ids[1], ids[2], ids[0], "addPerm", details, ptimestamp[0]);
                setLog([sn, uaddrids[msg.sender], ids[2]], "addPerm", 0x00, 0x00, 0, uint64(block.timestamp));
                return true;
            }
        }
    }

    /*@ 追加权限
    //@param id: 权限id
    //@param device: 设备限制
    //@param ptype: [查看权限, 下载权限, 授予权限]
    //@param iptimestamp: 权限过期时间戳
    //@param iptime: 权限追加时间
    //@param iptimes: 权限追加次数
    //@return 0: 设置是否成功
    */
    function apdPerm(bytes32 id, bytes32 device, uint64[3] ptype, uint64 iptimestamp, uint64 iptime, uint64 iptimes) internal returns (bool) {
        if (id == 0x00) {
            revert();
        } else {
            bytes32 senderid = uaddrids[msg.sender];
            Perm storage perm = perms[id];
            if (verPerm(senderid, perm.itemid, ptype, iptimestamp, iptime, iptimes)) {
                perm.signer = senderid;
                perm.device = device;
                perm.ptype = [ptype[0], ptype[1], ptype[2]];
                perm.ptimestamp = [perm.ptimestamp[0], iptimestamp];
                perm.ptime = [perm.ptime[0] + iptime, perm.ptime[1] + iptime];
                perm.ptimes = [perm.ptimes[0] + iptimes, perm.ptimes[1] + iptimes];
                perm.status = 1;

                User memory user = users[senderid];
                Item memory item = items[perm.itemid];
                if (item.userid != senderid && user.level != 1) {
                    if (setPermTimes(puiidids[concat(senderid, perm.itemid)], iptimes)) {
                        return true;
                    } else {
                        revert();
                    }
                } else {
                    return true;
                }
            } else {
                revert();
            }
        }
    }

    /*@ 追加权限By权限id
    //@param sn: 交易流水号
    //@param details: 交易详细信息
    //@param id: 权限id
    //@param device: 设备限制
    //@param ptype: [查看权限, 下载权限, 授予权限]
    //@param iptimestamp: 权限过期时间戳
    //@param iptime: 权限追加时间
    //@param iptimes: 权限追加次数
    //@return 0: 设置是否成功
    */
    function apdPermById(bytes32 sn, string details, bytes32 id, bytes32 device, uint64[3] ptype, uint64 iptimestamp, uint64 iptime, uint64 iptimes) public onlyValid returns (bool) {
        if (apdPerm(id, device, ptype, iptimestamp, iptime, iptimes)) {
            emit TXReceipt(sn, msg.sender, piduiids[id][0], piduiids[id][1], id, "apdPermById", details, uint64(block.timestamp));
            setLog([sn, uaddrids[msg.sender], piduiids[id][1]], "adpPermById", 0x00, 0x00, 0, uint64(block.timestamp));
            return true;
        } else {
            revert();
        }
    }

    /*@ 追加权限By权限索引
    //@param sn: 交易流水号
    //@param details: 交易详细信息
    //@param idx: 权限索引
    //@param device: 设备限制
    //@param ptype: [查看权限, 下载权限, 授予权限]
    //@param iptimestamp: 权限过期时间戳
    //@param iptime: 权限追加时间
    //@param iptimes: 权限追加次数
    //@return 0: 设置是否成功
    */
    function apdPermById(bytes32 sn, string details, uint64 idx, bytes32 device, uint64[3] ptype, uint64 iptimestamp, uint64 iptime, uint64 iptimes) public onlyValid returns (bool) {
        if (apdPerm(pidxids[idx], device, ptype, iptimestamp, iptime, iptimes)) {
            emit TXReceipt(sn, msg.sender, piduiids[pidxids[idx]][0], piduiids[pidxids[idx]][1], pidxids[idx], "apdPermByIdx", details, uint64(block.timestamp));
            setLog([sn, uaddrids[msg.sender], piduiids[pidxids[idx]][1]], "adpPermByIdx", 0x00, 0x00, 0, uint64(block.timestamp));
            return true;
        } else {
            revert();
        }
    }

    /*@ 追加权限By权限用户和资源
    //@param sn: 交易流水号
    //@param details: 交易详细信息
    //@param ids: [用户id, 资源id]
    //@param device: 设备限制
    //@param ptype: [查看权限, 下载权限, 授予权限]
    //@param iptimestamp: 权限过期时间戳
    //@param iptime: 权限追加时间
    //@param iptimes: 权限追加次数
    //@return 0: 设置是否成功
    */
    function apdPermByUI(bytes32 sn, string details, bytes32[2] ids, bytes32 device, uint64[3] ptype, uint64 iptimestamp, uint64 iptime, uint64 iptimes) public onlyValid returns (bool) {
        bytes32 id = puiidids[concat(ids[0], ids[1])];
        if (apdPerm(id, device, ptype, iptimestamp, iptime, iptimes)) {
            emit TXReceipt(sn, msg.sender, ids[0], ids[1], id, "apdPermByUI", details, uint64(block.timestamp));
            setLog([sn, uaddrids[msg.sender], ids[1]], "adpPermByUI", 0x00, 0x00, 0, uint64(block.timestamp));
            return true;
        } else {
            revert();
        }
    }

    /*@ 设置权限时间次数
    //@param id: 权限id
    //@param utime: 消耗的时间
    //@return 0: 是否设置成功
    */
    function setPermTime(bytes32 id, uint64 utime) internal returns (bool) {
        if (id == 0x00) {
            revert();
        } else {
            Perm storage perm = perms[id];
            bytes32 senderid = uaddrids[msg.sender];
            if (verPerm(senderid, perm.itemid, [uint64(0), 0, 0], 0, 0, 0)) {
                uint64 prtime = perm.ptime[1];
                // uint64 prtimes = perm.ptimes[1];
                if (prtime == 0) {
                    return false;
                } else {
                    if (utime >= prtime) {
                        prtime = 0;
                    } else {
                        prtime -= utime;
                    }
                    /*
                    if (utime > 0) {
                        if (prtimes <= 1) {
                            prtimes = 0;
                        } else {
                            prtimes -= 1;
                        }
                    }
                    */
                    perm.ptime = [perm.ptime[0], prtime];
                    // perm.ptimes = [perm.ptimes[0], prtimes];
                    return true;
                }
            } else {
                revert();
            }
        }
    }

    /*@ 设置权限次数
    //@param id: 权限id
    //@param utimes: 消耗的次数
    //@return 0: 是否设置成功
    */
    function setPermTimes(bytes32 id, uint64 utimes) internal returns (bool) {
        if (id == 0x00) {
            revert();
        } else {
            Perm storage perm = perms[id];
            bytes32 senderid = uaddrids[msg.sender];
            if (verPerm(senderid, perm.itemid, [uint64(0), 0, 0], 0, 0, 0)) {
                uint64 prtimes = perm.ptimes[1];
                if (prtimes == 0) {
                    return false;
                } else {
                    if (prtimes < utimes) {
                        prtimes = 0;
                    } else {
                        prtimes -= utimes;
                    }
                    perm.ptimes = [perm.ptimes[0], prtimes];
                    return true;
                }
            } else {
                revert();
            }
        }
    }

    /*@ 设置权限时间次数By权限id
    //@param sn: 交易流水号
    //@param details: 交易详细信息
    //@param id: 权限id
    //@param utime: 消耗的时间
    //@return 0: 是否设置成功
    */
    function setPermTimeById(bytes32 sn, string details, bytes32 id, uint64 utime) public onlyValid returns (bool) {
        if (setPermTime(id, utime)) {
            emit TXReceipt(sn, msg.sender, piduiids[id][0], piduiids[id][1], id, "setPermTimeById", details, uint64(block.timestamp));
            setLog([sn, uaddrids[msg.sender], piduiids[id][1]], "setPermTimeById", 0x00, 0x00, 0, uint64(block.timestamp));
            return true;
        } else {
            revert();
        }
    }

    /*@ 设置权限时间次数By权限索引
    //@param sn: 交易流水号
    //@param details: 交易详细信息
    //@param idx: 权限索引
    //@param utime: 消耗的时间
    //@return 0: 是否设置成功
    */
    function setPermTimeByIdx(bytes32 sn, string details, uint64 idx, uint64 utime) public onlyValid returns (bool) {
        bytes32 id = pidxids[idx];
        if (setPermTime(id, utime)) {
            emit TXReceipt(sn, msg.sender, piduiids[id][0], piduiids[id][1], id, "setPermTimeByIdx", details, uint64(block.timestamp));
            setLog([sn, uaddrids[msg.sender], piduiids[id][1]], "setPermTimeByIdx", 0x00, 0x00, 0, uint64(block.timestamp));
            return true;
        } else {
            revert();
        }
    }

    /*@ 设置权限时间次数By权限用户和资源
    //@param sn: 交易流水号
    //@param details: 交易详细信息
    //@param ids: [用户id, 资源id]
    //@param utime: 消耗的时间
    //@return 0: 是否设置成功
    */
    function setPermTimeByUI(bytes32 sn, string details, bytes32[2] ids, uint64 utime) public onlyValid returns (bool) {
        bytes32 id = puiidids[concat(ids[0], ids[1])];
        if (setPermTime(id, utime)) {
            emit TXReceipt(sn, msg.sender, ids[0], ids[1], id, "setPermTimeByUI", details, uint64(block.timestamp));
            setLog([sn, uaddrids[msg.sender], ids[1]], "setPermTimeByUI", 0x00, 0x00, 0, uint64(block.timestamp));
            return true;
        } else {
            revert();
        }
    }

    /*@ 设置权限状态
    //@param id: 权限id
    //@param newstatus: 权限新的状态
    //@return 0: 是否设置成功
    */
    function setPermStatus(bytes32 id, uint64 newstatus) internal returns (bool) {
        if (id == 0x00) {
            revert();
        } else {
            Perm storage perm = perms[id];
            bytes32 senderid = uaddrids[msg.sender];
            if (verPerm(senderid, perm.itemid, [uint64(0), 0, 0], 0, 0, 0)) {
                perm.status = newstatus;
                return true;
            } else {
                revert();
            }
        }
    }

    /*@ 设置权限状态By权限id
    //@param sn: 交易流水号
    //@param details: 交易详细信息
    //@param id: 权限id
    //@param newstatus: 权限新的状态
    //@return 0: 是否设置成功
    */
    function setPermStatusById(bytes32 sn, string details, bytes32 id, uint64 newstatus) public onlyValid returns (bool) {
        if (setPermStatus(id, newstatus)) {
            emit TXReceipt(sn, msg.sender, piduiids[id][0], piduiids[id][1], id, "setPermStatusById", details, uint64(block.timestamp));
            setLog([sn, uaddrids[msg.sender], piduiids[id][1]], "setPermStatusById", 0x00, 0x00, 0, uint64(block.timestamp));
            return true;
        } else {
            revert();
        }
    }

    /*@ 设置权限状态By权限索引
    //@param sn: 交易流水号
    //@param details: 交易详细信息
    //@param idx: 权限索引
    //@param newstatus: 权限新的状态
    //@return 0: 是否设置成功
    */
    function setPermStatusByIdx(bytes32 sn, string details, uint64 idx, uint64 newstatus) public onlyValid returns (bool) {
        bytes32 id = pidxids[idx];
        if (setPermStatus(id, newstatus)) {
            emit TXReceipt(sn, msg.sender, piduiids[id][0], piduiids[id][1], id, "setPermStatusByIdx", details, uint64(block.timestamp));
            setLog([sn, uaddrids[msg.sender], piduiids[id][1]], "setPermStatusByIdx", 0x00, 0x00, 0, uint64(block.timestamp));
            return true;
        } else {
            revert();
        }
    }

    /*@ 设置权限状态By权限用户和资源
    //@param sn: 交易流水号
    //@param details: 交易详细信息
    //@param ids: [用户id, 资源id]
    //@param newstatus: 权限新的状态
    //@return 0: 是否设置成功
    */
    function setPermStatusByUI(bytes32 sn, string details, bytes32[2] ids, uint64 newstatus) public onlyValid returns (bool) {
        bytes32 id = puiidids[concat(ids[0], ids[1])];
        if (setPermStatus(id, newstatus)) {
            emit TXReceipt(sn, msg.sender, piduiids[id][0], piduiids[id][1], id, "setPermStatusByUI", details, uint64(block.timestamp));
            setLog([sn, uaddrids[msg.sender], piduiids[id][1]], "setPermStatusByUI", 0x00, 0x00, 0, uint64(block.timestamp));
            return true;
        } else {
            revert();
        }
    }

    /*@ 验证用户对资源的权限
    //@param userid: 用户id
    //@param itemid: 资源id
    //@param ptype: [查看权限, 下载权限, 授权权限]
    //@param ptimestamp: 过期时间戳限制
    //@param ptime: 剩余时间限制
    //@param ptimes: 剩余次数限制
    //@return 0: 用户是否满足对资源的所有权限限制
    */
    function verPerm(bytes32 userid, bytes32 itemid, uint64[3] ptype, uint64 ptimestamp, uint64 ptime, uint64 ptimes) internal view returns (bool) {
        // TODO
        if (userid == 0x00 || itemid == 0x00) {
            return false;
        } else {
            User memory user = users[userid];
            Item memory item = items[itemid];
            if (user.status == 0) {
                return false;
            } else {
                if (user.level == 1) {
                    return true;
                } else {
                    if (item.userid == userid) {
                        return true;
                    } else {
                        bytes32 pid = puiidids[concat(userid, itemid)];
                        if ( pid == 0x00) {
                            return false;
                        } else {
                            Perm memory perm = perms[pid];
                            if (perm.status != 1 || perm.ptype[2] != 1) {
                                return false;
                            } else {
                                if (perm.ptype[0] >= ptype[0] && perm.ptype[1] >= ptype[1] && perm.ptype[2] >= ptype[2] && perm.ptimestamp[1] >= ptimestamp && perm.ptimestamp[1] > 0 && perm.ptime[1] >= ptime && perm.ptime[1] > 0 && perm.ptimes[1] >= ptimes && perm.ptimes[1] > 0) {
                                    return true;
                                } else {
                                    return false;
                                }
                            }
                        }
                    }
                }
            }
        }
    }

    
    /*@ 获取权限总数
    //@return 0: 权限总数
    */
    function getPermNum() public view returns (uint64) {
        return permnum;
    }

    /*@ 获取权限信息By权限id
    //@param id: 权限id
    //@return 0: [权限id, 授予者id, 授予用户id, 资源id]
    //@return 1: 权限播放设备限制
    //@return 2: [查看权限, 下载权限, 授予权限]
    //@return 3: [权限生成时间戳, 权限过期时间戳]
    //@return 4: [权限总时间, 权限剩余时间]
    //@return 5: [权限可用次数, 权限剩余可用次数]
    //@return 6: 权限状态
    */
    function getPermById(bytes32 id) public view returns (bytes32[4], bytes32, uint64[3], uint64[2], uint64[2], uint64[2], uint64) {
        if (id == 0x00) {
            return ([bytes32(0x00), 0x00, 0x00, 0x00], bytes32(0x00), [uint64(0), 0, 0], [uint64(0), 0], [uint64(0), 0], [uint64(0), 0], uint64(0));
        } else {
            Perm memory perm = perms[id];
            return ([perm.id, perm.signer, perm.userid, perm.itemid], perm.device, [perm.ptype[0], perm.ptype[1], perm.ptype[2]], [perm.ptimestamp[0], perm.ptimestamp[1]], [perm.ptime[0], perm.ptime[1]], [perm.ptimes[0], perm.ptimes[1]], perm.status);
        }
    }

    /*@ 获取权限信息By权限索引
    //@param idx: 权限索引
    //@return 0: [权限id, 授予者id, 授予用户id, 资源id]
    //@return 1: 权限播放设备限制
    //@return 2: [查看权限, 下载权限, 授予权限]
    //@return 3: [权限生成时间戳, 权限过期时间戳]
    //@return 4: [权限总时间, 权限剩余时间]
    //@return 5: [权限可用次数, 权限剩余可用次数]
    //@return 6: 权限状态
    */
    function getPermByIdx(uint64 idx) public view returns (bytes32[4], bytes32, uint64[3], uint64[2], uint64[2], uint64[2], uint64) {
        return getPermById(pidxids[idx]);
    }

    /*@ 获取权限信息By权限用户和资源
    //@param ids: [用户id, 资源id]
    //@return 0: [权限id, 授予者id, 授予用户id, 资源id]
    //@return 1: 权限播放设备限制
    //@return 2: [查看权限, 下载权限, 授予权限]
    //@return 3: [权限生成时间戳, 权限过期时间戳]
    //@return 4: [权限总时间, 权限剩余时间]
    //@return 5: [权限可用次数, 权限剩余可用次数]
    //@return 6: 权限状态
    */
    function getPermByUI(bytes32[2] ids) public view returns (bytes32[4], bytes32, uint64[3], uint64[2], uint64[2], uint64[2], uint64) {
        return getPermById(puiidids[concat(ids[0], ids[1])]);
    }


    /*@ 获取用户对资源的权限(管理员用户, 资源所有者, 权限拥有者)
    //@param userid: 用户id
    //@param itemid: 资源id
    //@return 0: 用户角色(0=参数不合法, 1=用户不可用, 2=管理员, 3=资源所有者, 4=无权限, 5=有权限)
    //@return 1: [权限id, 授予者id, 授予用户id, 资源id]
    //@return 2: [查看权限, 下载权限, 授予权限]
    //@return 3: [权限生成时间戳, 权限过期时间戳]
    //@return 4: [权限总时间, 权限剩余时间]
    //@return 5: [权限可用次数, 权限剩余可用次数]
    //@return 6: 权限状态
    */
    function fetchPerm(bytes32 userid, bytes32 itemid) public view returns (uint64, bytes32[4], uint64[3], uint64[2], uint64[2], uint64[2], uint64) {
        if (userid == 0x00 || itemid == 0x00) {
            return (0, [bytes32(0x00), 0x00, 0x00, 0x00], [uint64(0), 0, 0], [uint64(0), 0], [uint64(0), 0], [uint64(0), 0], 0);
        } else {
            User memory user = users[userid];
            if (user.status == 0) {
                return (1, [bytes32(0x00), 0x00, 0x00, 0x00], [uint64(0), 0, 0], [uint64(0), 0], [uint64(0), 0], [uint64(0), 0], 0);
            } else {
                if (user.level == 1) {
                    return (2, [bytes32(0x00), 0x00, 0x00, 0x00], [uint64(0), 0, 0], [uint64(0), 0], [uint64(0), 0], [uint64(0), 0], 0);
                }
                Item memory item = items[itemid];
                if (item.userid == userid) {
                    return (3, [bytes32(0x00), 0x00, 0x00, 0x00], [uint64(0), 0, 0], [uint64(0), 0], [uint64(0), 0], [uint64(0), 0], 0);
                }
                bytes32 pid = puiidids[concat(userid, itemid)];
                if (pid == 0x00) {
                    return (4, [bytes32(0x00), 0x00, 0x00, 0x00], [uint64(0), 0, 0], [uint64(0), 0], [uint64(0), 0], [uint64(0), 0], 0);
                } else {
                    Perm memory perm = perms[pid];
                    return (5, [perm.id, perm.signer, perm.userid, perm.itemid], [perm.ptype[0], perm.ptype[1], perm.ptype[2]], [perm.ptimestamp[0], perm.ptimestamp[1]], [perm.ptime[0], perm.ptime[1]], [perm.ptimes[0], perm.ptimes[1]], perm.status);
                }
            }
        }
    }
    

    /*@ 设置日志
    //@param ids: [日志id, 用户id, 资源id]
    //@param action: 日志操作类型
    //@param hash: 日志hash
    //@param desc: 日志描述
    //@param duration: 日志操作时长
    //@param timestamp: 日志生成时间戳
    //@return 0: 是否设置成功
    */
    function setLog(bytes32[3] ids, bytes32 action, bytes32 hash, bytes32 desc, uint64 duration, uint64 timestamp) internal returns (bool) {
        lidxids[lognum++] = ids[0];
        logs[ids[0]] = Log(ids[0], ids[1], ids[2], action, hash, desc, duration, timestamp);
        liduiids[ids[0]] = [ids[1], ids[2]];
        return true;
    }

    /*@ 添加日志
    //@param sn: 交易流水号
    //@param details: 交易详细信息
    //@param ids: [日志id, 用户id, 资源id]
    //@param action: 日志操作类型
    //@param hash: 日志hash
    //@param desc: 日志描述
    //@param duration: 日志操作时长
    //@param mode: 日志操作类型(0=只记日志不修改权限次数和时间, 1=记日志的同时修改权限次数和时间)
    //@return 0: 是否设置成功
    */
    function addLog(bytes32 sn, string details, bytes32[3] ids, bytes32 action, bytes32 hash, bytes32 desc, uint64 duration, uint64 mode) public onlyAdmin returns (bool) {
        if (setLog(ids, action, hash, desc, duration, uint64(block.timestamp))) {
            if (mode != 0) {
                bytes32 pid = puiidids[concat(ids[1], ids[2])];
                if (duration == 0) {
                    setPermTimes(pid, 1);
                } else {
                    setPermTime(pid, duration);
                }
            }
            emit TXReceipt(sn, msg.sender, ids[1], ids[2], ids[0], "addLog", details, uint64(block.timestamp));
            return true;
        } else {
            revert();
        }
    }

    /*@ 获取日志总数
    //@return 0: 日志总数
    */
    function getLogNum() public view returns (uint64) {
        return lognum;
    }

    /*@ 获取日志信息By日志id
    //@param id: 日志id
    //@return 0: [日志id, 用户id, 资源id]
    //@return 1: 日志操作类型
    //@return 2: 日志hash
    //@return 3: 日志描述
    //@return 4: 日志操作时长
    //@return 5: 日志生成时间戳
    */
    function getLogById(bytes32 id) public view returns (bytes32[3], bytes32, bytes32, bytes32, uint64, uint64) {
        if (id == 0x00) {
            return ([bytes32(0x00), 0x00, 0x00], bytes32(0x00), 0x00, 0x00, uint64(0), 0);
        } else {
            Log memory log = logs[id];
            return ([log.id, log.userid, log.itemid], log.action, log.hash, log.desc, log.duration, log.timestamp);
        }
    }

    /*@ 获取日志信息By日志索引
    //@param idx: 日志索引
    //@return 0: [日志id, 用户id, 资源id]
    //@return 1: 日志操作类型
    //@return 2: 日志hash
    //@return 3: 日志描述
    //@return 4: 日志操作时长
    //@return 5: 日志生成时间戳
    */
    function getLogByIdx(uint64 idx) public view returns (bytes32[3], bytes32, bytes32, bytes32, uint64, uint64) {
        return getLogById(lidxids[idx]);
    }

    /*@ 合约自毁
    */
    function kill() public onlyOwner {
        selfdestruct(owner);
    }


}