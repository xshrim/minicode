# 链码说明

*[注]以下所有函数的返回值均为字节数组. 需转换为字符串后再进行返回值解析.*

## 链上对象

![bcobj](./bcobj.png)

## 链码函数

### 初始化

- 调用函数: **init**

- 函数说明:

> 链码初始化函数, 仅可在链码实例化和升级时自动调用, 目前无需参数传递

- 函数参数:


- 请求示例:

  - 命令行:
  ```bash
  peer chaincode instantiate -n example -v 1 -c '{"Args":["init"]}' -C mych -P   "AND ('Org1MSP.member')"
  ```

### 用户初始化

- 调用函数: **userInit**

- 函数说明: 

> 对用户进行链上初始化, 未初始化的用户不可对区块链进行其他请求. 初始化无需额外参数, 用户信息将从用户证书中读取.

- 函数参数:

- 请求示例:

  - 命令行:
  ```bash
  peer chaincode invoke -C mych -n example -c '{"Args":["userInit"]}'
  ```
  - SDK(go):
  ```go
  req := channel.Request{ChaincodeID: "mych", Fcn: "userInit", Args: nil}
  respone, err := t.Client.Execute(req)
  ```

### 资源上传

- 调用函数: **assetUpload**

- 函数说明: 

> 为用户提供资源上传入链的接口. 用户上传资源后自动与用户账号绑定. 资源哈希既作为链上快速查询的key, 又可作为校验资源其他相关外部数据真实性的hash码. 资源外部链接用于定位资源实际存储位置. 资源价值表示资源交易产生的价值转移.

- 函数参数:

| 参数名称 | 参数类型 | 参数说明               | 备注   |
| -------- | -------- | ---------------------- | ------ |
| assetHash | string   | 资源哈希, 唯一确定资源 | 不可空 |
| assetLink | string   | 资源外部链接           | 不可空 |
| assetPerm | string   | 资源权限               |  可空      |
| assetPlist | string   | 资源白名单/黑名单      | 可空 |

- 请求示例:

  - 命令行:
  ```bash
  peer chaincode invoke -C mych -n example -c '{"Args":["assetUpload", "03a471b6c3d2...", "url...", "-1,-1,300,100,0", "ebg.ebtech.oms:true, ebg.ebtech:false"]}'
  ```
  - SDK(go):
  ```go
  assetHash := "03a471b6c3d2..."
  assetLink := "url..."
  assetPerm := ""
  assetPlist := ""
  req := channel.Request{ChaincodeID: "mych", Fcn: "assetUpload", Args: [][]byte{[]byte(assetHash), []byte(assetLink), []byte(assetPerm), []byte(assetPlist)}}
  respone, err := t.Client.Execute(req)
  ```

### 资源购买

- 函数调用: **assetBuy**

- 函数说明:

> 为用户提供资源购买的接口. 资源购买不涉及所有权转移, 只是将资源标记为已拥有状态, 并扣除购买者于资源等值的积分余额, 同时为资源所有者(上传者)增加相同的积分余额.

- 函数参数:

| 参数名称  | 参数类型 | 参数说明         | 备注   |
| --------- | -------- | ---------------- | ------ |
| assetHash | string   | 要购买的资源哈希 | 不可空 |

- 请求示例:

  - 命令行:
  ```bash
  peer chaincode invoke -C mych -n example -c '{"Args":["assetBuy", "03a471b6c3d2..."]}'
  ```
  - SDK(go):
  ```go
  assetHash := "03a471b6c3d2..."
  req := channel.Request{ChaincodeID: "mych", Fcn: "assetBuy", Args: [][]byte{[]byte(assetHash)}}
  respone, err := t.Client.Execute(req)
  ```

### 资源浏览/上传次数统计

- 函数调用: **assetTimesCount**

- 函数说明:

> 为用户提供资源浏览/上传次数统计的接口. 

- 函数参数:

| 参数名称  | 参数类型 | 参数说明                             | 备注   |
| --------- | -------- | ------------------------------------ | ------ |
| assetHash | string   | 要计数的资源哈希                     | 不可空 |
| countType | string   | 计数类型(view: 浏览; download: 下载) | 不可空 |

- 请求示例:

  - 命令行:
  ```bash
  peer chaincode invoke -C mych -n example -c '{"Args":["assetTimesCount", "03a471b6c3d2...", "view"]}'
  ```
  - SDK(go):
  ```go
  assetHash := "03a471b6c3d2..."
  countType := "view"
  req := channel.Request{ChaincodeID: "mych", Fcn: "assetTimesCount", Args: [][]byte{[]byte(assetHash), []byte(countType)}}
  respone, err := t.Client.Execute(req)
  ```

### 资源冻结

- 函数调用: **assetFreeze**

- 函数说明:

> 为用户提供资源冻结的接口. 可一次性冻结一个或多个资源. 被冻结的资源不会移除, 只是不允许再被购买.

- 函数参数:

| 参数名称   | 参数类型 | 参数说明         | 备注   |
| ---------- | -------- | ---------------- | ------ |
| assetHashs | []string | 要冻结的资源哈希 | 不可空 |

- 请求示例:
  - 命令行:
  ```bash
  peer chaincode invoke -C mych -n example -c '{"Args":["assetFreeze", "03a471b6c3d2...", "26b902c813a1..."]}'
  ```
  - SDK(go):
  ```go
  assetHash1 := "03a471b6c3d2..."
  assetHash2 := "26b902c813a1..."
  req := channel.Request{ChaincodeID: "mych", Fcn: "assetFreeze", Args: [][]byte{[]byte(assetHash1), []byte(assetHash2)}}}
  respone, err := t.Client.Execute(req)
  ```

### 资源废弃

- 函数调用: **assetDiscard**

- 函数说明:

> 为用户提供资源废弃的接口. 可一次性废弃一个或多个资源. 被废弃的资源不会移除, 只是不再可见.

- 函数参数:

| 参数名称   | 参数类型 | 参数说明         | 备注   |
| ---------- | -------- | ---------------- | ------ |
| assetHashs | []string | 要废弃的资源哈希 | 不可空 |

- 请求示例:
  - 命令行:
  ```bash
  peer chaincode invoke -C mych -n example -c '{"Args":["assetDiscard", "03a471b6c3d2...", "26b902c813a1..."]}'
  ```
  - SDK(go):
  ```go
  assetHash1 := "03a471b6c3d2..."
  assetHash2 := "26b902c813a1..."
  req := channel.Request{ChaincodeID: "mych", Fcn: "assetDiscard", Args: [][]byte{[]byte(assetHash1), []byte(assetHash2)}}}
  respone, err := t.Client.Execute(req)
  ```

### 积分充值

- 函数调用: **balanceRecharge**

- 函数说明:

> 为用户提供积分充值接口. 

- 函数参数:

| 参数名称 | 参数类型 | 参数说明                           | 备注   |
| -------- | -------- | ---------------------------------- | ------ |
| userid   | string   | 用户ID, 形如User1@org1.example.com | 不可空 |
| value    | string   | 需要充值的积分值, 不允许为负       | 不可空 |

- 请求示例:
  - 命令行:
  ```bash
  peer chaincode invoke -C mych -n example -c '{"Args":["balanceRecharge", "User1@org1.example.com", "100"]}'
  ```
  - SDK(go):
  ```go
  userid := "User1@org1.example.com"
  value := "100"
  req := channel.Request{ChaincodeID: "mych", Fcn: "balanceRecharge", Args: [][]byte{[]byte(userid), []byte(value)}}
  respone, err := t.Client.Execute(req)
  ```

### 积分转让

- 函数调用: **balanceTransfer**

- 函数说明:

> 为用户提供积分转让接口. 

- 函数参数:

| 参数名称 | 参数类型 | 参数说明                                   | 备注   |
| -------- | -------- | ------------------------------------------ | ------ |
| userid   | string   | 需转让的用户ID, 形如User1@org1.example.com | 不可空 |
| value    | string   | 需转让的积分值, 不允许为负                 | 不可空 |

- 请求示例:
  - 命令行:
  ```bash
  peer chaincode invoke -C mych -n example -c '{"Args":["balanceTransfer", "User1@org1.example.com", "100"]}'
  ```
  - SDK(go):
  ```go
  userid := "User1@org1.example.com"
  value := "100"
  req := channel.Request{ChaincodeID: "mych", Fcn: "balanceTransfer", Args: [][]byte{[]byte(userid), []byte(value)}}
  respone, err := t.Client.Execute(req)
  ```
  
### 资源评论

- 函数调用: **commentAdd**

- 函数说明:

> 为用户提供资源评论接口. 

- 函数参数:

| 参数名称     | 参数类型 | 参数说明       | 备注   |
| ------------ | -------- | -------------- | ------ |
| commentHash  | string   | 评论哈希       | 不可空 |
| assetHash    | string   | 评论的资源哈希 | 不可空 |
| commentFloor | string   | 评论楼层       | 可空   |
| assetGrade   | string   | 对资源的评分   | 可空   |

- 请求示例:
  - 命令行:
  ```bash
  peer chaincode invoke -C mych -n example -c '{"Args":["commentAdd", "32b71e46a5c8...", "03a471b6c3d2...", "", ""]}'
  ```
  - SDK(go):
  ```go
  commentHash := "32b71e46a5c8..."
  assetHash := "03a471b6c3d2..."
  commentFloor := ""
  assetGrade := ""
  req := channel.Request{ChaincodeID: "mych", Fcn: "commentAdd", Args: [][]byte{[]byte(commentHash), []byte(assetHash), []byte(commentFloor), []byte(assetGrade)}}
  respone, err := t.Client.Execute(req)
  ```

### 评论赞成/反对

- 函数调用: **commentRateSet**

- 函数说明:

> 为用户提供资源评论的赞成和反对接口. 

- 函数参数:

| 参数名称    | 参数类型 | 参数说明                   | 备注   |
| ----------- | -------- | -------------------------- | ------ |
| commentHash | string   | 评论哈希                   | 不可空 |
| ptype       | string   | agree: 赞成; against: 反对 | 不可空 |

- 请求示例:
  - 命令行:
  ```bash
  peer chaincode invoke -C mych -n example -c '{"Args":["commentRateSet", "32b71e46a5c8...", "agree"]}'
  ```
  - SDK(go):
  ```go
  commentHash := "32b71e46a5c8..."
  ptype := "agree"
  req := channel.Request{ChaincodeID: "mych", Fcn: "commentRateSet", Args: [][]byte{[]byte(commentHash), []byte(ptype)}}
  respone, err := t.Client.Execute(req)
  ```
  

### 评论废弃

- 函数调用: **commentDiscard**

- 函数说明:

> 为用户提供资源评论废弃接口. 可一次性废弃一个或多个评论. 被废弃的评论不会删除, 只是不可见.

- 函数参数:

| 参数名称     | 参数类型 | 参数说明 | 备注   |
| ------------ | -------- | -------- | ------ |
| commentHashs | []string | 评论哈希 | 不可空 |

- 请求示例:
  - 命令行:
  ```bash
  peer chaincode invoke -C mych -n example -c '{"Args":["commentRateSet", "32b71e46a5c8...", "03a471b6c3d2..."]}'
  ```
  - SDK(go):
  ```go
  commentHash1 := "32b71e46a5c8..."
  commentHash2 := "03a471b6c3d2..."
  req := channel.Request{ChaincodeID: "mych", Fcn: "commentDiscard", Args: [][]byte{[]byte(commentHash1), []byte(commentHash2)}}
  respone, err := t.Client.Execute(req)
  ```

### 资源查询

- 函数调用: **assetQuery**

- 函数说明:

> 根据资源哈希提供资源查询接口. 可以一次查询一个或多个资源哈希.

- 函数参数:

| 参数名称   | 参数类型 | 参数说明         | 备注   |
| ---------- | -------- | ---------------- | ------ |
| assetHashs | []string | 要查询的资源哈希 | 不可空 |

- 请求示例:
  - 命令行:
  ```bash
  peer chaincode query -C mych -n example -c '{"Args":["assetQuery", "03a471b6c3d2...", "26b902c813a1..."]}'
  ```
  - SDK(go):
  ```go
  assetHash1 := "03a471b6c3d2..."
  assetHash2 := "26b902c813a1..."
  req := channel.Request{ChaincodeID: "mych", Fcn: "assetQuery", Args: [][]byte{[]byte(assetHash1), []byte(assetHash2)}}
  respone, err := t.Client.Query(req)
  ```
  
### 资源评论查询

- 函数调用: **assetCommentQuery**

- 函数说明:

> 根据资源哈希提供资源评论查询接口.

- 函数参数:

| 参数名称  | 参数类型 | 参数说明         | 备注   |
| --------- | -------- | ---------------- | ------ |
| assetHash | string   | 要查询的资源哈希 | 不可空 |
| pageSize  | string   | 分页页号         | 可选   |
| bookmark  | string   | 分页标记         | 可选   |

- 请求示例:
  - 命令行:
  ```bash
  peer chaincode query -C mych -n example -c '{"Args":["assetCommentQuery", "03a471b6c3d2..."]}'
  ```
  - SDK(go):
  ```go
  assetHash := "03a471b6c3d2..."
  req := channel.Request{ChaincodeID: "mych", Fcn: "assetCommentQuery", Args: [][]byte{[]byte(assetHash)}}
  respone, err := t.Client.Query(req)
  ```

### 用户评论查询

- 函数调用: **userCommentQuery**

- 函数说明:

> 为用户提供其所有资源评论查询接口.

- 函数参数:

| 参数名称  | 参数类型 | 参数说明         | 备注   |
| --------- | -------- | ---------------- | ------ |
| assetHash | string   | 要查询的资源哈希 | 不可空 |
| pageSize  | string   | 分页页号         | 可选   |
| bookmark  | string   | 分页标记         | 可选   |

- 请求示例:
  - 命令行:
  ```bash
  peer chaincode query -C mych -n example -c '{"Args":["userCommentQuery"]}'
  ```
  - SDK(go):
  ```go
  req := channel.Request{ChaincodeID: "mych", Fcn: "userCommentQuery", Args: nil}
  respone, err := t.Client.Query(req)
  ```

### 评论查询

- 函数调用: **commentQuery**

- 函数说明:

> 资源评论查询和用户评论查询只能查询到评论哈希, 需要使用评论哈希通过评论查询获取评论具体信息. 可一次查询一个或多个评论信息.

- 函数参数:

| 参数名称     | 参数类型 | 参数说明         | 备注   |
| ------------ | -------- | ---------------- | ------ |
| commentHashs | []string | 要查询的评论哈希 | 不可空 |
| pageSize     | string   | 分页页号         | 可选   |
| bookmark     | string   | 分页标记         | 可选   |

- 请求示例:
  - 命令行:
  ```bash
  peer chaincode query -C mych -n example -c '{"Args":["commentQuery", "03a471b6c3d2...", "26b902c813a1..."]}'
  ```
  - SDK(go):
  ```go
  commentHash1 := "03a471b6c3d2..."
  commentHash2 := "26b902c813a1..."
  req := channel.Request{ChaincodeID: "mych", Fcn: "commentQuery", Args: [][]byte{[]byte(commentHash1), []byte(commentHash2)}}
  respone, err := t.Client.Query(req)
  ```


### 日志查询

- 函数调用: **recordQuery**

- 函数说明:

> 对区块链数据的写操作都将自动记录操作日志, 本函数支持查询指定日志key或者指定用户的操作记录.

- 函数参数:

| 参数名称    | 参数类型 | 参数说明                             | 备注   |
| ----------- | -------- | ------------------------------------ | ------ |
| recordStart | string   | 要查询的日志key或者日志范围的Start点 | 不可空 |
| recordEnd   | string   | 要查询的日志范围的End点              | 可选   |
| pageSize    | string   | 分页页号                             | 可选   |
| bookmark    | string   | 分页标记                             | 可选   |

- 请求示例:
  - 命令行:
  ```bash
  peer chaincode query -C mych -n example -c '{"Args":["recordQuery", "User1@org1.example.com:1553743949"]}'
  ```
  - SDK(go):
  ```go
  keyStart := "User1@org1.example.com:0"
  keyEnd := "User1@org1.example.com:a"
  req := channel.Request{ChaincodeID: "mych", Fcn: "recordQuery", Args: [][]byte{[]byte(keyStart), []byte(keyEnd)}}
  respone, err := t.Client.Query(req)
  ```

### 用户日志查询

- 函数调用: **userRecordQuery**

- 函数说明:

> 为用户提供其所有操作日志查询接口.

- 函数参数:

| 参数名称 | 参数类型 | 参数说明 | 备注 |
| -------- | -------- | -------- | ---- |
| pageSize | string   | 分页页号 | 可选 |
| bookmark | string   | 分页标记 | 可选 |

- 请求示例:
  - 命令行:
  ```bash
  peer chaincode query -C mych -n example -c '{"Args":["userRecordQuery"]}'
  ```
  - SDK(go):
  ```go
  req := channel.Request{ChaincodeID: "mych", Fcn: "userRecordQuery", Args: nil}
  respone, err := t.Client.Query(req)
  ```

### 用户资源查询

- 函数调用: **userAssetQuery**

- 函数说明:

> 用于查询与用户相关的所有资源(包括上传和购买). 无需传入参数, 直接以函数调用者作为要查询的用户, 返回用户ID, 资源Hash和关联关系(0:上传, 1:购买).

- 函数参数:

| 参数名称 | 参数类型 | 参数说明 | 备注 |
| -------- | -------- | -------- | ---- |
| pageSize | string   | 分页页号 | 可选 |
| bookmark | string   | 分页标记 | 可选 |

- 请求示例:
  - 命令行:
  ```bash
  peer chaincode query -C mych -n example -c '{"Args":["userAssetQuery"]}'
  ```
  - SDK(go):
  ```go
  req := channel.Request{ChaincodeID: "mych", Fcn: "userAssetQuery", Args: nil}
  respone, err := t.Client.Query(req)
  ```

### 积分查询

- 函数调用: **balanceQuery**

- 函数说明:

> 查询用户的积分余额. 无需传入参数, 直接将函数调用者作为要查询的用户.

- 函数参数:

- 请求示例:
  - 命令行:
  ```bash
  peer chaincode query -C mych -n example -c '{"Args":["balanceQuery"]}'
  ```
  - SDK(go):
  ```go
  req := channel.Request{ChaincodeID: "mych", Fcn: "balanceQuery", Args: nil}
  respone, err := t.Client.Query(req)
  ```
  
### 用户查询

- 函数调用: **userQuery**
  
- 函数说明:
  
> 用户仅有ID会入链, 因此用户查询目前仅返回查询时间戳, 用户所属组织的mspid和用户ID. 无需传入参数, 直接将函数调用者作为要查询的用户.

- 函数参数:
  
- 请求示例:
  - 命令行:
  ```bash
  peer chaincode query -C mych -n example -c '{"Args":["userQuery"]}'
  ```
  - SDK(go):
  ```go
  req := channel.Request{ChaincodeID: "mych", Fcn: "userQuery", Args: nil}
  respone, err := t.Client.Query(req)
  ```

### 历史查询

- 函数调用: **historyQuery**

- 函数说明:

> 查询指定key的所有变更记录. 支持简单key和复合key. 但一次仅支持查询一个key.

- 函数参数:

| 参数名称 | 参数类型 | 参数说明                         | 备注   |
| -------- | -------- | -------------------------------- | ------ |
| key      | []string | 单元素的简单key和多元素的复合key | 不可空 |

- 请求示例:
  - 命令行:
  ```bash
  peer chaincode query -C mych -n example -c '{"Args":["historyQuery"], "User1@org1.example.com"}' # 简单key
  ```
  - SDK(go):
  ```go
  keyType := "user~asset"
  userid := "User1@org1.example.com"
  assetHash := "03a471b6c3d2..."
  req := channel.Request{ChaincodeID: "mych", Fcn: "historyQuery", Args: [][]byte{[]byte(keyType), []byte(userid), []byte(assetHash)}}  // 复合key
  respone, err := t.Client.Query(req)
  ```

### 全查询
- 函数调用: **fullQuery**

- 函数说明:

> 查询指定对象类型(user, asset, comment, record)的全部数据.

- 函数参数:

| 参数名称     | 参数类型 | 参数说明                         | 备注                       |
| ----------- | -------- | -------------------------------- | -------------------------------- |
| qtype | string | 要查询的对象类型(user, asset, comment, record) | 不可空 |
| pageSize | string | 分页页号 | 可选 |
| bookmark | string | 分页标记 | 可选 |

- 请求示例:
  - 命令行:
  ```bash
  peer chaincode query -C mych -n example -c '{"Args":["fullQuery", "asset"]}'
  ```
  - SDK(go):
  ```go
  qtype := "asset"
  req := channel.Request{ChaincodeID: "mych", Fcn: "fullQuery", Args: [][]byte{[]byte(qtype)}}
  respone, err := t.Client.Query(req)
  ```

### 用户资源验证
- 函数调用: **userAssetVerify**

- 函数说明:

> 查询用户对指定资源的访问权限. 以函数调用者作为要查询的用户.

- 函数参数:

| 参数名称     | 参数类型 | 参数说明                         | 备注                       |
| ----------- | -------- | -------------------------------- | -------------------------------- |
| assetHash   | string | 要查询的资源哈希 | 不可空 |

- 请求示例:
  - 命令行:
  ```bash
  peer chaincode query -C mych -n example -c '{"Args":["userAssetVerify", "03a471b6c3d2..."]}'
  ```
  - SDK(go):
  ```go
  assetHash := "03a471b6c3d2..."
  req := channel.Request{ChaincodeID: "mych", Fcn: "userAssetVerify", Args: [][]byte{[]byte(assetHash)}}
  respone, err := t.Client.Query(req)
  ```

### 富查询

- 函数调用: **richQueryResult**

- 函数说明:

| 参数名称     | 参数类型 | 参数说明                         | 备注                       |
| ----------- | -------- | -------------------------------- | -------------------------------- |
| querystr | string | 富查询选择器 | 不可空 |
| pageSize | string | 分页页号 | 可选 |
| bookmark | string | 分页标记 | 可选 |

> 针对使用couchdb作为后台状态数据库的Fabric网络, 以couchdb支持的语法格式进行高级查询. 链码暂时直接以传入参数作为选择器, 未支持链码分析传入参数后自动构建选择器.

- 请求示例:
  - 命令行:
  ```bash
  peer chaincode query -C mych -n example -c '{"Args":["richQueryResult", "{\"selector\":{\"owner\":\"tom\"}}","3",""]}'
  ```
  - SDK(go):
  ```go
  querystr := "{\"selector\":{\"owner\":\"tom\"}}"
  pageSize := "3"
  bookmark := ""
  req := channel.Request{ChaincodeID: "mych", Fcn: "userAssetVerify", Args: [][]byte{[]byte(querystr), []byte(pageSize), []byte(bookmark)}}
  respone, err := t.Client.Query(req)
  ```