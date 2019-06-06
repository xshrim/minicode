package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type Token struct {
}

type Record struct {
	TxID    string   `json:"txid"`    // 交易ID
	Subject string   `json:"subject"` // 行为主体
	Object  []string `json:"object"`  // 行为客体
	Action  string   `json:"action"`  // 行为
	Date    int64    `json:"date"`    // 行为时间戳
	Notes   string   `json:"notes"`   // 行为备注
}

type Event struct { // Fabric的链码事件默认会返回区块号, 交易ID, 链码名, 事件名和事件内容, 因此事件内容也可以不必包含下述部分信息
	TxID     string `json:"txid"`
	UserName string `json:"username"`
	FcName   string `json:"fcname"`
	Message  string `json:"message"`
	Date     int64  `json:"date"`
}

// 参数检查
func paramCheck(args []string, expect int) string {
	if len(args) != expect {
		return "Incorrect number of arguments. Expecting " + strconv.Itoa(expect)
	}

	for _, arg := range args {
		if strings.TrimSpace(arg) == "" {
			return "The arguments can not be empty string"
		}
	}
	return ""
}

func getKey(stub shim.ChaincodeStubInterface, skey string) (key string) {

	defer func() {
		if err := recover(); err != nil {
			key = skey
		}
	}()

	if !strings.Contains(skey, "~") {
		key = skey
	} else {
		_, compositeKeyParts, err := stub.SplitCompositeKey(skey)
		if err != nil {
			key = skey
		} else {
			key = strings.Join(compositeKeyParts, ":")
		}
	}

	return
}

func addPaginationMetadataToQueryResults(buffer *bytes.Buffer, responseMetadata *peer.QueryResponseMetadata) {

	buffer.WriteString("[{\"Meta\":{\"Count\":")
	buffer.WriteString("\"")
	buffer.WriteString(fmt.Sprintf("%v", responseMetadata.FetchedRecordsCount))
	buffer.WriteString("\"")
	buffer.WriteString(", \"Mark\":")
	buffer.WriteString("\"")
	buffer.WriteString(responseMetadata.Bookmark)
	buffer.WriteString("\"}}]")

}

func constructQueryResponseFromIterator(stub shim.ChaincodeStubInterface, resultIterator shim.StateQueryIteratorInterface) *bytes.Buffer {
	var buffer bytes.Buffer

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			continue
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		key := getKey(stub, queryResponse.Key)
		value := queryResponse.Value

		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(value))
		buffer.WriteString("}")

		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")
	return &buffer
}

// 发起链码事件
func setEvent(stub shim.ChaincodeStubInterface, txid, username, fcname, message string, date int64) bool {
	event := &Event{
		TxID:     txid,
		UserName: username,
		FcName:   fcname,
		Message:  message,
		Date:     date,
	}

	logger.Info(event)

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return false
	}

	err = stub.SetEvent(fcname, eventBytes)
	if err != nil {
		return false
	}

	//logger.Info(err)

	return true
}

func recordPut(stub shim.ChaincodeStubInterface, txid, subject string, object []string, action string, date int64, notes string) bool {
	/*
		stKey, err := stub.CreateCompositeKey("subject~timestamp", []string{subject, date})
		if err != nil {
			return false
		}
	*/

	record := &Record{
		TxID:    txid,
		Subject: subject,
		Object:  object,
		Action:  action,
		Date:    date,
		Notes:   notes,
	}

	recordBytes, err := json.Marshal(record)
	if err != nil || recordBytes == nil || len(recordBytes) == 0 {
		return false
	}

	if err := stub.PutState("record:"+txid, recordBytes); err != nil {
		return false
	}

	urKey, err := stub.CreateCompositeKey("user~record", []string{record.Subject, record.TxID})
	if err != nil {
		return false
	}

	if err := stub.PutState(urKey, []byte(record.Action+"-"+strconv.FormatInt(record.Date, 10))); err != nil {
		return false
	}

	return true
}

// 获取余额
func balanceGet(stub shim.ChaincodeStubInterface, uname string) (int64, error) {

	balanceBytes, err := stub.GetState("user:" + uname)
	if err == nil && balanceBytes != nil && len(balanceBytes) > 0 {
		return strconv.ParseInt(string(balanceBytes), 10, 64)
	}
	return -1, fmt.Errorf("User not found")
}

// 增减余额
func balanceAdd(stub shim.ChaincodeStubInterface, uname string, value int64) (int64, error) {
	balance, err := balanceGet(stub, uname)
	if err != nil {
		return -1, err
	}

	balance += value

	if balance < 0 {
		return -1, fmt.Errorf("Balance is less than 0")
	}

	if err := stub.PutState("user:"+uname, []byte(strconv.FormatInt(balance, 10))); err != nil {
		return -1, fmt.Errorf("Add balance error: " + err.Error())
	}

	return balance, nil
}

func balanceTransfer(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("balanceTransfer")

	if res := paramCheck(args, 3); res != "" {
		return shim.Error(res)
	}

	var txtime int64
	txid := stub.GetTxID()
	timestamp, err := stub.GetTxTimestamp()
	if err == nil {
		txtime = timestamp.Seconds
	}

	sender := args[0]

	receiver := args[1]

	value, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil || value < 0 {
		return shim.Error("Balance value is illegal")
	}

	if value <= 0 {
		return shim.Error("Balance value must more than 0")
	}

	if _, err := balanceAdd(stub, sender, 0-value); err != nil {
		return shim.Error("Balance set failed: " + err.Error())
	}

	if _, err := balanceAdd(stub, receiver, value); err != nil {
		return shim.Error("Balance set failed: " + err.Error())
	}

	if ok := recordPut(stub, txid, sender, []string{receiver}, "balanceTransfer", txtime, "value:"+args[2]); !ok {
		return shim.Error("Put record failure")
	}

	if ok := setEvent(stub, txid, sender, "balanceTransfer", "OK", txtime); !ok {
		return shim.Error("Triggle event failure")
	}

	//logger.Info("OK")

	return shim.Success(nil)
}

// 积分充值(可负)
func balanceDeposit(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("balanceDeposit")

	if res := paramCheck(args, 2); res != "" {
		return shim.Error(res)
	}

	var txtime int64
	txid := stub.GetTxID()
	timestamp, err := stub.GetTxTimestamp()
	if err == nil {
		txtime = timestamp.Seconds
	}

	uname := args[0]

	value, err := strconv.ParseInt(args[1], 10, 64) // 允许value为负, 保留扣减积分窗口
	if err != nil || value < 0 {
		return shim.Error("Balance value is illegal")
	}

	balance, err := balanceGet(stub, uname)
	if err != nil { // 积分信息不存在,
		balance = value
	} else {
		balance += value
	}

	if err := stub.PutState("user:"+uname, []byte(strconv.FormatInt(balance, 10))); err != nil {
		return shim.Error("Deposit balance error: " + err.Error())
	}

	if ok := recordPut(stub, txid, uname, []string{uname}, "balanceDeposit", txtime, "value:"+args[1]); !ok {
		return shim.Error("Put record failure")
	}

	if ok := setEvent(stub, txid, uname, "balanceDeposit", "OK", txtime); !ok {
		return shim.Error("Triggle event failure")
	}

	return shim.Success(nil)
}

func recordQuery(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("recordQuery")

	var psizestr string
	var bookmark string
	var buffer *bytes.Buffer

	length := len(args)

	if length < 2 || length > 4 {
		return shim.Error("Incorrect number of arguments. Expecting 2~4")
	}
	if length > 2 {
		psizestr = args[2]
		if length > 3 {
			bookmark = args[3]
		}
	}

	if args[0] == "" {
		return shim.Error("The first argument can not be empty string")
	}

	if args[1] == "" {
		recordBytes, err := stub.GetState("record:" + args[0])
		if err != nil || recordBytes == nil || len(recordBytes) == 0 {
			return shim.Error("Record not found")
		}

		return shim.Success(recordBytes)
	}

	if psizestr == "" {
		resultIterator, err := stub.GetStateByRange("record:"+args[0], "record:"+args[1])
		if err != nil {
			return shim.Error("Query record error: " + err.Error())
		}
		defer resultIterator.Close()

		buffer = constructQueryResponseFromIterator(stub, resultIterator)
	} else {
		pagesize, err := strconv.ParseInt(psizestr, 10, 32)
		if err != nil {
			return shim.Error("Rich query error: " + err.Error())
		}
		resultsIterator, responseMetadata, err := stub.GetStateByRangeWithPagination("record:"+args[0], "record:"+args[1], int32(pagesize), bookmark)
		if err != nil {
			return shim.Error("Rich query error: " + err.Error())
		}
		defer resultsIterator.Close()

		buffer = constructQueryResponseFromIterator(stub, resultsIterator)

		addPaginationMetadataToQueryResults(buffer, responseMetadata)
	}

	return shim.Success(buffer.Bytes())
}

// 查询用户的所有操作记录
func userRecordQuery(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("userRecordQuery")

	var psizestr string
	var bookmark string
	var buffer *bytes.Buffer

	length := len(args)
	if length < 1 || length > 3 {
		return shim.Error("Incorrect number of arguments. Expecting 1~3")
	}
	if length > 1 {
		psizestr = args[1]
		if length > 2 {
			bookmark = args[2]
		}
	}

	/*
			if len(args) > 1 {
				shim.Error("Incorrect number of arguments. Expecting 0 or 1")
			}
			if len(args) == 0 {
				_, _, userid = stubParse(stub)
			} else {
				userid = args[0]
		    }
	*/

	if psizestr == "" {
		resultIterator, err := stub.GetStateByPartialCompositeKey("user~record", []string{args[0]})
		if err != nil {
			return shim.Error("Query user record error: " + err.Error())
		}
		defer resultIterator.Close()

		buffer = constructQueryResponseFromIterator(stub, resultIterator)
	} else {
		pagesize, err := strconv.ParseInt(psizestr, 10, 32)
		if err != nil {
			return shim.Error("Query user record error: " + err.Error())
		}
		resultsIterator, responseMetadata, err := stub.GetStateByPartialCompositeKeyWithPagination("user~record", []string{args[0]}, int32(pagesize), bookmark)
		if err != nil {
			return shim.Error("Query user record error: " + err.Error())
		}
		defer resultsIterator.Close()

		buffer = constructQueryResponseFromIterator(stub, resultsIterator)

		addPaginationMetadataToQueryResults(buffer, responseMetadata)
	}

	return shim.Success(buffer.Bytes())
}

// 查询用户积分
func balanceQuery(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("balanceQuery")

	if res := paramCheck(args, 1); res != "" {
		return shim.Error(res)
	}

	balance, err := balanceGet(stub, args[0])
	if err != nil {
		return shim.Error("Balance of the user not found")
	}

	return shim.Success([]byte(strconv.FormatInt(balance, 0)))
}

// 键的历史记录
func historyQuery(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("historyQuery")

	var key string
	var buffer bytes.Buffer

	length := len(args)

	if length < 1 {
		return shim.Error("Incorrect number of arguments. Expecting > 0")
	}

	if args[0] == "" {
		return shim.Error("The first argument can not be empty string")
	}

	if length == 1 {
		key = args[0]
	} else {
		var err error
		key, err = stub.CreateCompositeKey(args[0], args[1:])
		if err != nil {
			return shim.Error("Create composite key error: " + err.Error())
		}
	}

	resultIterator, err := stub.GetHistoryForKey(key)
	if err != nil {
		return shim.Error("Get history for key error: " + err.Error())
	}
	defer resultIterator.Close()

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			continue
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		item, err := json.Marshal(queryResponse)
		if err != nil {
			continue
		}

		buffer.WriteString(string(item))

		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

func fullQuery(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("fullQuery")

	var psizestr string
	var bookmark string
	var buffer *bytes.Buffer

	length := len(args)
	if length < 1 || length > 3 {
		return shim.Error("Incorrect number of arguments. Expecting 1~3")
	}

	if args[0] == "" {
		return shim.Error("The first argument can not be empty string")
	}

	if length > 1 {
		psizestr = args[1]
		if length > 2 {
			bookmark = args[2]
		}
	}

	//resultIterator, err := stub.GetStateByPartialCompositeKey("subject~timestamp", args)
	resultIterator, err := stub.GetStateByRange(args[0]+":", args[0]+":~")
	if err != nil {
		return shim.Error("Query error: " + err.Error())
	}
	defer resultIterator.Close()

	if psizestr == "" {
		resultIterator, err := stub.GetStateByRange(args[0]+":", args[0]+":~")
		if err != nil {
			return shim.Error("Query error: " + err.Error())
		}
		defer resultIterator.Close()

		buffer = constructQueryResponseFromIterator(stub, resultIterator)
	} else {
		pagesize, err := strconv.ParseInt(psizestr, 10, 32)
		if err != nil {
			return shim.Error("Query error: " + err.Error())
		}
		resultsIterator, responseMetadata, err := stub.GetStateByRangeWithPagination(args[0]+":", args[0]+":~", int32(pagesize), bookmark)
		if err != nil {
			return shim.Error("Query error: " + err.Error())
		}
		defer resultsIterator.Close()

		buffer = constructQueryResponseFromIterator(stub, resultsIterator)

		addPaginationMetadataToQueryResults(buffer, responseMetadata)
	}

	return shim.Success(buffer.Bytes())
}

func (t *Token) Init(stub shim.ChaincodeStubInterface) peer.Response {
	/*
			_, args := stub.GetFunctionAndParameters()
			if len(args) != 1 {
				return shim.Error("Incorrect number of arguments. Expecting 1")
		    }
	*/
	return shim.Success(nil)
}

func (t *Token) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "balanceTransfer":
		// 积分转账: [参数1] 转账人; [参数2] 被转账人; [参数3] 转账积分
		return balanceTransfer(stub, args)
	case "balanceDeposit":
		// 积分充值(允许充值负数以保留积分扣减通道): [参数1] 被充值人; [参数2] 充值积分
		return balanceDeposit(stub, args)
	case "recordQuery":
		// 记录查询(查询行为记录): [参数1] 起始交易号; [参数2] 结束交易号(可选); [参数3] 分页大小(可选); [参数4] 分页标记(可选)
		return recordQuery(stub, args)
	case "userRecordQuery":
		// 用户记录查询(查询指定用户行为记录): [参数1] 用户名; [参数2] 分页大小(可选); [参数3] 分页标记(可选)
		return userRecordQuery(stub, args)
	case "balanceQuery":
		// 积分查询: [参数1] 用户名;
		return balanceQuery(stub, args)
	case "historyQuery":
		// 历史记录查询: [参数1] 键值
		return historyQuery(stub, args)
	case "fullQuery":
		// 指定类型全部数据查询: [参数1] 类型(record, user); [参数2] 分页大小(可选); [参数3] 分页标记(可选)
		return fullQuery(stub, args)
	default:
		return shim.Error("Unsupported function: " + function)
	}

	// return shim.Success(nil)
}

var logger = shim.NewLogger("ebclog")

func main() {

	logger.SetLevel(shim.LogDebug)

	err := shim.Start(new(Token))
	if err != nil {
		logger.Error("Error starting AssetManage chaincode: %s", err)
	}
}
