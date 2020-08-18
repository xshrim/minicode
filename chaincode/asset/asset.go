package main

import (
	"bytes"
	"crypto/x509"
	"encoding/asn1"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

const (
	permSwitch = "off"
)

type AssetManage struct {
}

// Attributes contains attribute names and values
type Attributes struct {
	Attrs map[string]string `json:"attrs"`
}

/*
// Names returns the names of the attributes
func (a *Attributes) Names() []string {
	i := 0
	names := make([]string, len(a.Attrs))
	for name := range a.Attrs {
		names[i] = name
		i++
	}
	return names
}

// Contains returns true if the named attribute is found
func (a *Attributes) Contains(name string) bool {
	_, ok := a.Attrs[name]
	return ok
}

// Value returns an attribute's value
func (a *Attributes) Value(name string) (string, bool, error) {
	attr, ok := a.Attrs[name]
	return attr, ok, nil
}
*/

type User struct {
	Name  string `json:"name"`  // 用户名
	Hash  string `json:"hash"`  // 用户链外信息哈希
	Mspid string `json:"mspid"` // 用户组织
	//Vest    string    `json:"vest"`    // 用户归属(集团.公司.部门.团队)
	//Posts   [4]string `json:"posts"`   // 用户职位(集团, 公司, 部门, 团队)
	Role    string `json:"role"`    // 用户角色
	Balance int64  `json:"balance"` // 用户积分余额
	Date    int64  `json:"date"`    // 用户创建时间
	Ddate   int64  `json:"ddate"`   // 用户离职时间
	Status  int    `json:"status"`  // 用户状态(在职, 离职)
}

type Asset struct {
	Hash string `json:"hash"` // 资源哈希
	// Name  string            `json:"name"`  // 资源名
	Link  string `json:"link"`  // 资源链接
	Owner string `json:"owner"` // 资源上传者
	Mspid string `json:"mspid"` // 资源组织
	Value int64  `json:"value"` // 资源价格
	//Vest  string            `json:"vest"`  // 资源归属(集团.公司.部门.团队)
	Times [2]int64 `json:"times"` // 资源浏览, 下载次数
	Rate  [2]int64 `json:"rate"`  // 资源评分(评分总数/评分次数)
	//Perm  [5]string         `json:"perm"`  // 资源访问权限(联盟内, 集团内, 公司内, 部门内, 团队内)
	Plist  map[string]string `json:"plist"`  // 资源访问特殊权限名单
	Date   int64             `json:"date"`   // 资源上传时间戳
	Status int               `json:"status"` // 资源状态(正常, 冻结, 不可用)
}

type Comment struct {
	Hash      string `json:"hash"`      // 评论哈希
	AssetHash string `json:"assethash"` // 对应资源哈希
	// Content    string `json:"content"`    // 评论正文
	Reviewer string   `json:"reviewer"` // 评论者
	Floor    string   `json:"floor"`    // 评论楼层
	Rate     [2]int64 `json:"rate"`     // 评论赞成反对数(赞成数/反对数)
	Grade    int64    `json:"grade"`    // 对资源的评分
	Date     int64    `json:"date"`     // 评论时间戳
	Status   int64    `json:"status"`   // 评论状态(正常, 不显示, 不存在)

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

// 判断leader是否是staff的上级(组织管理员, 同级领导或上一级领导)
func isLeader(staff string, leader *User) bool {
	if permSwitch == "off" {
		return true
	}

	return false
}

// Is the object ID equal to the attribute info object ID?
func isAttrOID(oid asn1.ObjectIdentifier) bool {
	AttrOID := asn1.ObjectIdentifier{1, 2, 3, 4, 5, 6, 7, 8, 1}
	if len(oid) != len(AttrOID) {
		return false
	}
	for idx, val := range oid {
		if val != AttrOID[idx] {
			return false
		}
	}
	return true
}

// Get the attribute info from a certificate extension, or return nil if not found
func getAttributesFromCert(cert *x509.Certificate) ([]byte, error) {
	for _, ext := range cert.Extensions {
		if isAttrOID(ext.Id) {
			return ext.Value, nil
		}
	}
	return nil, nil
}

/*
func getVest(name string) (string, error) {
	uname := strings.Split(name, "@")
	if len(uname) != 2 {
		return "", errors.New("Illegal user name")
	}

	var uvest string

	utail := strings.Split(uname[1], ".")
	for i, j := 0, len(utail)-1; i < j; i, j = i+1, j-1 {
		utail[i], utail[j] = utail[j], utail[i]
	}
	uvest = strings.Join(utail, ".")

	uhead := strings.Split(uname[0], ".")

	if len(uhead) > 1 {
		uteam := uhead[1]
		if uteam != "" {
			if uvest != "" {
				uvest = uvest + "." + uteam
			} else {
				uvest = uteam
			}
		}
	}

	return uvest, nil
}
*/

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

// 解析stub获取交易请求时间戳, 交易请求者mspid和userid
func stubParse(stub shim.ChaincodeStubInterface) (string, string, string, int64, error) {

	var txtime int64
	var mspid string

	role := "user"

	txid := stub.GetTxID()

	timestamp, err := stub.GetTxTimestamp()
	if err != nil {
		txtime = int64(0)
	} else {
		// 毫秒级时间戳
		// timestr = strconv.Itoa(int(timestamp.Seconds*1000 + int64(timestamp.Nanos/1000000)))
		txtime = timestamp.Seconds
		// timestr = timestamp.String()

		// timestr = time.Unix(timestamp.Seconds, int64(timestamp.Nanos)).String()
	}

	creatorBytes, err := stub.GetCreator()
	if err != nil || creatorBytes == nil || len(creatorBytes) == 0 {
		creatorBytes = []byte("")
	}

	// 获取证书部分
	certStart := bytes.IndexAny(creatorBytes, "-----BEGIN")
	if certStart == -1 {
		return "", role, txid, txtime, errors.New("Certificate not found")
	}

	// 使用正则表达式过滤掉creatorBytes中mspid部分的不可见字符(protobuf编码)
	mspid = regexp.MustCompile(`\w+`).FindString(string(creatorBytes[:certStart]))

	// 直接从creatorBytes获取mspid
	//msp := creatorBytes[2 : certStart-3] // mspid两端存在不可见字符

	certText := creatorBytes[certStart:]

	// 获取证书的pem格式编码
	bl, _ := pem.Decode(certText)
	if bl == nil {
		return mspid, role, txid, txtime, errors.New("Certificate decode error")
	}

	// 生成证书对象
	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return mspid, role, txid, txtime, errors.New("Certificate parse error")
	}

	// 获取证书中包含的各种数据

	cname := cert.Subject.CommonName // 用户名, 如Admin@org1.example.com
	if strings.ToLower(strings.Split(cname, "@")[0]) == "admin" {
		role = "admin"
	}
	/*
		user.Vest, err = getVest(user.Name)
		if err != nil {
			user.Vest = ""
		}
	*/

	buff, err := getAttributesFromCert(cert)
	if err != nil {
		return mspid, role, txid, txtime, nil
	}

	attrs := &Attributes{}
	if buff != nil {
		err := json.Unmarshal(buff, attrs)
		if err != nil {
			return mspid, role, txid, txtime, nil
		}
	}

	/*
		if v, ok := attrs.Attrs["bpost"]; ok {
			user.Posts[0] = v
		}

		if v, ok := attrs.Attrs["cpost"]; ok {
			user.Posts[1] = v
		}

		if v, ok := attrs.Attrs["dpost"]; ok {
			user.Posts[2] = v
		}

		if v, ok := attrs.Attrs["epost"]; ok {
			user.Posts[3] = v
		}
	*/

	if v, ok := attrs.Attrs["role"]; ok {
		role = strings.ToLower(v)
	}

	return mspid, role, txid, txtime, nil
}

/*
func IDCounter(stub shim.ChaincodeStubInterface, domain string) string {
	nextID := 0
	idBytes, err := stub.GetState(domain)
	if err == nil && len(idBytes) > 0 {
		if curtID, err := strconv.Atoi(string(idBytes)); err == nil {
			nextID = curtID + 1
		}
	}
	return fmt.Sprintf("%013s", nextID)
}
*/

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

/*
func parseIterator(stub shim.ChaincodeStubInterface, resultIterator shim.StateQueryIteratorInterface) []byte {
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

		kv := make(map[string][]byte)
		kv[key] = value

		resp, err := json.Marshal(kv)
		if err != nil {
			continue
		}

		buffer.Write(resp)

		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")
	return buffer.Bytes()
}
*/

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

func userGet(stub shim.ChaincodeStubInterface, uname string) (*User, error) {
	userBytes, err := stub.GetState("user:" + uname)

	if err != nil {
		return nil, errors.New("User not found")
	}

	user := new(User)
	if err := json.Unmarshal(userBytes, user); err != nil {
		return nil, errors.New("Unmarshal user error: " + err.Error())
	}

	return user, nil
}

func userPut(stub shim.ChaincodeStubInterface, user *User) error {
	userBytes, err := json.Marshal(user)
	if err != nil || userBytes == nil || len(userBytes) == 0 {
		return err
	}

	if err := stub.PutState("user:"+user.Name, userBytes); err != nil {
		return err
	}

	return nil
}

/*
// 获取余额
func balanceGet(stub shim.ChaincodeStubInterface, uname string) (int64, error) {

	user, err := userGet(stub, uname)
	if err != nil {
		return 0, err
	}

	return user.Balance, nil
}

// 增减余额
func balanceAdd(stub shim.ChaincodeStubInterface, uname string, value int64) (int64, error) {
	user, err := userGet(stub, uname)
	if err != nil {
		return 0, err
	}

	user.Balance += value

	if user.Balance < 0 {
		return -1, errors.New("Balance is less than 0")
	}

	err = userPut(stub, user)
	if err != nil {
		return -1, err
	}

	return user.Balance, nil
}
*/

func balanceTransfer(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("balanceTransfer")

	if res := paramCheck(args, 3); res != "" {
		return shim.Error(res)
	}

	/*
		if len(args) != 2 {
			return shim.Error("Incorrect number of arguments. Expecting 2")
		}

		if args[0] == "" || args[1] == "" {
			return shim.Error("Arguments can not be empty string")
		}
	*/

	mspid, role, txid, txtime, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	if role != "admin" {
		return shim.Error("Only admin can do this")
	}

	sender, err := userGet(stub, args[0])
	if err != nil {
		return shim.Error("Sender not found")
	}

	if mspid != sender.Mspid {
		return shim.Error("Only users belong you can be operated")
	}

	receiver, err := userGet(stub, args[1])
	if err != nil {
		return shim.Error("Receiver not found")
	}

	if sender.Status != 0 || receiver.Status != 0 {
		return shim.Error("User is invalid")
	}

	value, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil || value < 0 {
		return shim.Error("Balance value is illegal")
	}

	if value <= 0 {
		return shim.Error("Balance value must more than 0")
	}

	if value > sender.Balance {
		return shim.Error("Balance not enough")
	}

	sender.Balance -= value

	receiver.Balance += value

	if err := userPut(stub, sender); err != nil {
		return shim.Error("Balance transform failed")
	}

	if err := userPut(stub, receiver); err != nil {
		return shim.Error("Balance transform failed")
	}

	/*
		if _, err := balanceAdd(stub, sender.ID, 0-value); err != nil {
			return shim.Error("Balance transform failed")
		}

		if _, err := balanceAdd(stub, receiver.ID, value); err != nil {
			return shim.Error("Balance transform failed")
		}
	*/

	if ok := recordPut(stub, txid, sender.Name, []string{receiver.Name}, "balanceTransfer", txtime, "value:"+args[2]); !ok {
		return shim.Error("Put record failure")
	}

	if ok := setEvent(stub, txid, sender.Name, "balanceTransfer", "OK", txtime); !ok {
		return shim.Error("Triggle event failure")
	}

	//logger.Info("OK")

	return shim.Success(nil)
}

// 积分充值
func balanceDeposit(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("balanceDeposit")

	if res := paramCheck(args, 2); res != "" {
		return shim.Error(res)
	}

	/*
		if len(args) != 2 {
			return shim.Error("Incorrect number of arguments. Expecting 2")
		}

		if args[0] == "" || args[1] == "" {
			return shim.Error("Arguments can not be empty string")
		}
	*/

	mspid, role, txid, txtime, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	if role != "admin" {
		return shim.Error("Only admin can do this")
	}

	user, err := userGet(stub, args[0])
	if err != nil {
		return shim.Error("User not found")
	}

	if user.Status != 0 {
		return shim.Error("User is invalid")
	}

	if mspid != user.Mspid {
		return shim.Error("Only users belong you can be operated")
	}

	value, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil || value < 0 {
		return shim.Error("Balance value is illegal")
	}

	if value <= 0 {
		return shim.Error("Balance value must more than 0")
	}

	user.Balance += value

	if err := userPut(stub, user); err != nil {
		return shim.Error("Balance deposit failed")
	}

	if ok := recordPut(stub, txid, role+"@"+mspid, []string{user.Name}, "balanceDeposit", txtime, "value:"+args[1]); !ok {
		return shim.Error("Put record failure")
	}

	if ok := setEvent(stub, txid, user.Name, "balanceDeposit", "OK", txtime); !ok {
		return shim.Error("Triggle event failure")
	}

	return shim.Success([]byte(strconv.FormatInt(user.Balance, 10)))

	//logger.Info("OK")

	return shim.Error("Balance deposit failure")
}

func userInit(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("userInit")

	if res := paramCheck(args, 3); res != "" {
		return shim.Error(res)
	}

	mspid, role, txid, txtime, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	if role != "admin" {
		return shim.Error("Only admin can do this")
	}

	if ouser, err := userGet(stub, args[0]); err == nil && ouser != nil {
		if ouser.Status == 0 || ouser.Ddate == 0 { // 还在职的用户信息不允许被覆盖
			return shim.Error("User already exist")
		}
	}

	user := &User{
		Name:    args[0] + "@" + mspid,
		Hash:    args[1],
		Mspid:   mspid,
		Role:    args[2],
		Balance: 100, // 初始积分
		Date:    txtime,
		Ddate:   int64(0),
		Status:  0,
	}

	if err := userPut(stub, user); err != nil {
		return shim.Error("Save user failed")
	}

	if ok := recordPut(stub, txid, role+"@"+mspid, []string{user.Name}, "userInit", txtime, "value:"+strconv.FormatInt(user.Balance, 10)); !ok {
		return shim.Error("Put record failure")
	}

	if ok := setEvent(stub, txid, user.Name, "userInit", "OK", txtime); !ok {
		return shim.Error("Triggle event failure")
	}

	//logger.Info("OK")

	return shim.Success(nil)
}

// 用户离职
func userDismiss(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("userDismiss")

	if res := paramCheck(args, 3); res != "" {
		return shim.Error(res)
	}

	mspid, role, txid, txtime, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	if role != "admin" {
		return shim.Error("Only admin can do this")
	}

	user, err := userGet(stub, args[0])
	if err != nil {
		return shim.Error("User not found")
	}

	if mspid != user.Mspid {
		return shim.Error("Only users belong you can be operated")
	}

	if user.Status != 0 || user.Ddate != 0 {
		return shim.Error("User is already dissmissed")
	}

	user.Ddate = txtime // 用户离职时间

	user.Status = 1 // 标记用户离职

	if err := userPut(stub, user); err != nil {
		return shim.Error("Save user failed")
	}

	if ok := recordPut(stub, txid, role+"@"+mspid, []string{user.Name}, "userDismiss", txtime, ""); !ok {
		return shim.Error("Put record failure")
	}

	if ok := setEvent(stub, txid, user.Name, "userDismiss", "OK", txtime); !ok {
		return shim.Error("Triggle event failure")
	}

	//logger.Info("OK")

	return shim.Success(nil)
}

// 通过asset hash获取asset对象
func assetGet(stub shim.ChaincodeStubInterface, assetHash string) (*Asset, error) {
	assetBytes, err := stub.GetState("asset:" + assetHash)
	if err != nil || assetBytes == nil || len(assetBytes) == 0 {
		return nil, errors.New("Asset not found error: " + err.Error())
	}

	asset := new(Asset)
	if err := json.Unmarshal(assetBytes, asset); err != nil {
		return nil, errors.New("Unmarshal asset error: " + err.Error())
	}

	return asset, nil
}

func assetPut(stub shim.ChaincodeStubInterface, asset *Asset) error {
	assetBytes, err := json.Marshal(asset)
	if err != nil || assetBytes == nil || len(assetBytes) == 0 {
		return err
	}

	if err := stub.PutState("asset:"+asset.Hash, assetBytes); err != nil {
		return err
	}

	return nil
}

// 上传asset
func assetUpload(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("assetUpload")

	if res := paramCheck(args, 4); res != "" {
		return shim.Error(res)
	}

	mspid, role, txid, txtime, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	if role != "admin" {
		return shim.Error("Only admin can do this")
	}

	assetBytes, err := stub.GetState("asset:" + args[0])
	if err == nil && assetBytes != nil && len(assetBytes) > 0 {
		return shim.Error("Asset already exist")
	}

	user, err := userGet(stub, args[2])
	if err != nil && user == nil {
		return shim.Error("User not found")
	}

	if mspid != user.Mspid {
		return shim.Error("Only users belong you can be operated")
	}

	if user.Status != 0 {
		return shim.Error("User is invalid")
	}

	value, err := strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		return shim.Error("Asset value is illegal")
	}

	asset := &Asset{
		Hash: args[0],
		// Name:  args[1],
		Link:  args[1],
		Owner: user.Name,
		Mspid: mspid,
		Value: value,
		//Vest:  user.Vest,
		Times: [2]int64{0, 0},
		Rate:  [2]int64{0, 0},
		//Perm:  [5]string{"-1", "-1", "-1", "-1", "-1"},
		Plist:  make(map[string]string),
		Date:   txtime,
		Status: 0,
	}

	/*
		// 将"-1,-1,300,100,0"字符串的访问权限参数转换为数组形式, 范围从大到小, 缺省值为"-1"
		if args[2] != "" {
			for idx, val := range strings.Split(args[2], ",") {
				asset.Perm[idx] = strings.TrimSpace(val)
			}
		}
	*/

	/*
		// 设置特殊权限, 参数形式: "ebg.ebtech.oms.team1-zhangsan:-1,ebg.ebtech.oms:0"
		// TODO 名单中各个元素的范围冲突问题待定
		if args[3] != "" {
			for _, val := range strings.Split(args[3], ",") {
				kv := strings.Split(val, ":")
				if len(kv) != 2 {
					shim.Error("Illegal plist argument")
				}
				asset.Plist[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
			}
		}
	*/

	if err := assetPut(stub, asset); err != nil {
		return shim.Error("Save asset error: " + err.Error())
	}

	user.Balance += 5 // 上传奖励积分

	if err := userPut(stub, user); err != nil {
		return shim.Error(err.Error())
	}

	// 用户和资产关联键
	uaKey, err := stub.CreateCompositeKey("user~asset", []string{asset.Owner, asset.Hash})
	if err != nil {
		return shim.Error("Create key error: " + err.Error())
	}

	// 用户和资产关联键值为关系-时间戳
	if err := stub.PutState(uaKey, []byte("owner-"+strconv.FormatInt(txtime, 10))); err != nil {
		return shim.Error("Associate user and asset error: " + err.Error())
	}

	if ok := recordPut(stub, txid, asset.Owner, []string{asset.Hash}, "assetUpload", asset.Date, ""); !ok {
		return shim.Error("Put record failure")
	}

	if ok := setEvent(stub, txid, user.Name, "assetUpload", "OK", txtime); !ok {
		return shim.Error("Triggle event failure")
	}

	//logger.Info("OK")

	return shim.Success(nil)
}

// asset状态变更(冻结不可购买, 废弃不可见)
func assetStatusSet(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("assetStatusSet")

	if res := paramCheck(args, 2); res != "" {
		return shim.Error(res)
	}

	mspid, role, txid, txtime, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	if role != "admin" {
		return shim.Error("Only admin can do this")
	}

	hashs := strings.Replace(args[0], " ", "", -1)

	status, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("status is illegal")
	}

	if status < 0 || status > 2 {
		return shim.Error("status is illegal")
	}

	for _, assetHash := range strings.Split(hashs, ",") {
		asset, err := assetGet(stub, assetHash)
		if err != nil || asset == nil {
			return shim.Error("Asset " + assetHash + " no found")
		}
		if mspid != asset.Mspid {
			return shim.Error("Asset " + assetHash + " is not belong to you ")
		}
		asset.Status = status

		if err := assetPut(stub, asset); err != nil {
			return shim.Error("Save asset error: " + err.Error())
		}
	}

	if ok := recordPut(stub, txid, role+"@"+mspid, args, "assetStatusSet", txtime, ""); !ok {
		return shim.Error("Put record failure")
	}

	if ok := setEvent(stub, txid, role+"@"+mspid, "assetStatusSet", "OK", txtime); !ok {
		return shim.Error("Triggle event failure")
	}

	//logger.Info("OK")

	return shim.Success(nil)
}

// 购买asset
func assetBuy(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("assetBuy")

	if res := paramCheck(args, 2); res != "" {
		return shim.Error(res)
	}

	mspid, role, txid, txtime, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	if role != "admin" {
		return shim.Error("Only admin can do this")
	}

	user, err := userGet(stub, args[0])
	if err != nil || user == nil {
		return shim.Error("User not found")
	}

	asset, err := assetGet(stub, args[1])
	if err != nil || asset == nil {
		return shim.Error("Asset not found")
	}

	if mspid != asset.Mspid {
		return shim.Error("Asset not belong to you")
	}

	/*
		if asset.Owner == user.Name {
			return shim.Error("User is the owner of the asset")
		}

		if isLeader(asset.Owner, user) {
			return shim.Error("User is the leader of the owner of the asset")
		}

		uaKey, err := stub.CreateCompositeKey("user~asset", []string{user.Name, asset.Hash})
		if err != nil {
			return shim.Error("Create key error: " + err.Error())
		}

		if uaBytes, err := stub.GetState(uaKey); err == nil && len(uaBytes) > 0 {
			return shim.Error("User have bought the asset")
		}
	*/

	if user.Name == asset.Owner {
		return shim.Error("User is owner of the asset")
	}

	if ok := userAssetAccess(stub, asset, user); ok {
		return shim.Error("User already has access to the asset")
	}

	if asset.Value < 0 {
		return shim.Error("Asset can not be bought")
	}

	if asset.Value == 0 {
		return shim.Error("Asset no need to buy")
	}

	if user.Balance < asset.Value {
		return shim.Error("Balance is insufficient")
	}

	owner, err := userGet(stub, asset.Owner)
	if err != nil {
		return shim.Error("The owner of the asset not found")
	}

	user.Balance -= asset.Value

	if err := userPut(stub, user); err != nil {
		return shim.Error("Payment error: " + err.Error())
	}

	if owner.Status == 0 && owner.Ddate == 0 && owner.Mspid == asset.Mspid { // 资源所有者没离职且状态正常才接收积分
		owner.Balance += asset.Value

		if err := userPut(stub, owner); err != nil {
			return shim.Error("Payment error: " + err.Error())
		}
	} else {
		// TODO 由谁接收积分
	}

	uaKey, err := stub.CreateCompositeKey("user~asset", []string{user.Name, asset.Hash})
	if err != nil {
		return shim.Error("Create key error: " + err.Error())
	}

	if err := stub.PutState(uaKey, []byte("user-"+strconv.FormatInt(txtime, 10))); err != nil {
		return shim.Error("Associate user and asset error: " + err.Error())
	}

	if ok := recordPut(stub, txid, user.Name, []string{asset.Hash}, "assetBuy", txtime, "value:"+strconv.FormatInt(asset.Value, 10)); !ok {
		return shim.Error("Put record failure")
	}

	if ok := setEvent(stub, txid, user.Name, "assetBuy", "OK", txtime); !ok {
		return shim.Error("Triggle event failure")
	}

	//logger.Info("OK")

	return shim.Success(nil)
}

// 资源浏览和下载次数统计
func assetTimesCount(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("assetTimesCount")

	if res := paramCheck(args, 2); res != "" {
		return shim.Error(res)
	}

	mspid, role, txid, txtime, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	if role != "admin" {
		return shim.Error("Only admin can do this")
	}

	asset, err := assetGet(stub, args[0])
	if err != nil || asset == nil {
		return shim.Error("Asset not found")
	}

	if mspid != asset.Mspid {
		return shim.Error("Asset not belong to you")
	}

	if asset.Status != 0 && asset.Status != 1 { // 资源不可用
		return shim.Error("Asset is invalid")
	}

	if strings.ToLower(args[1]) == "view" {
		asset.Times[0] += 1
	} else {
		asset.Times[1] += 1
	}

	// TODO 是否需要给上传者积分奖励(奖励可能导致恶意刷新)

	owner, err := userGet(stub, asset.Owner)
	if err != nil {
		return shim.Error("The owner of the asset not found")
	}

	if owner.Status == 0 && owner.Ddate == 0 && owner.Mspid == asset.Mspid { // 用户未离职且状态正常才可得到积分奖励
		owner.Balance += 1 // 奖励积分

		if err := userPut(stub, owner); err != nil {
			return shim.Error("Save user error: " + err.Error())
		}
	} else {
		// TODO 由谁获得积分
	}

	if err := assetPut(stub, asset); err != nil {
		return shim.Error("Count asset error: " + err.Error())
	}

	if ok := recordPut(stub, txid, role+"@"+mspid, []string{asset.Hash}, "assetTimesCount", txtime, args[1]+":1"); !ok {
		return shim.Error("Put record failure")
	}

	if ok := setEvent(stub, txid, role+"@"+mspid, "assetTimesCount", "OK", txtime); !ok {
		return shim.Error("Triggle event failure")
	}

	//logger.Info("OK")

	return shim.Success(nil)
}

// 通过comment hash获取comment对象
func commentGet(stub shim.ChaincodeStubInterface, commentHash string) (*Comment, error) {
	commentBytes, err := stub.GetState("comment:" + commentHash)
	if err != nil || commentBytes == nil || len(commentBytes) == 0 {
		return nil, errors.New("Comment not found error: " + err.Error())
	}

	comment := new(Comment)
	if err := json.Unmarshal(commentBytes, comment); err != nil {
		return nil, errors.New("Unmarshal comment error: " + err.Error())
	}

	return comment, nil
}

func commentPut(stub shim.ChaincodeStubInterface, comment *Comment) error {
	commentBytes, err := json.Marshal(comment)
	if err != nil || commentBytes == nil || len(commentBytes) == 0 {
		return err
	}

	if err := stub.PutState("comment:"+comment.Hash, commentBytes); err != nil {
		return err
	}

	return nil
}

// 添加评论
func commentAdd(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("commentAdd")

	if res := paramCheck(args, 5); res != "" {
		return shim.Error(res)
	}

	mspid, role, txid, txtime, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	if role != "admin" {
		return shim.Error("Only admin can do this")
	}

	commentBytes, err := stub.GetState("comment:" + args[0])
	if err == nil && commentBytes != nil && len(commentBytes) > 0 {
		return shim.Error("Comment already exist")
	}

	asset, err := assetGet(stub, args[1])
	if err != nil || asset == nil {
		return shim.Error("Asset not found")
	}

	if asset.Status != 0 {
		return shim.Error("Comment not allowed to this asset")
	}

	user, err := userGet(stub, args[2])
	if err != nil {
		return shim.Error("User not found")
	}

	if user.Status != 0 || user.Ddate != 0 || mspid != user.Mspid {
		return shim.Error("User not belong to you")
	}

	if res := userAssetAccess(stub, asset, user); !res {
		return shim.Error("Comment for the asset not allowed by the user")
	}

	grade, err := strconv.ParseInt(args[4], 10, 64)
	if err != nil {
		return shim.Error("Grade value is illegal")
	}

	comment := &Comment{
		Hash: args[0],
		// Name:  args[1],
		AssetHash: args[1],
		Reviewer:  user.Name,
		Floor:     args[3],
		Rate:      [2]int64{0, 0},
		Grade:     grade,
		Date:      txtime,
		Status:    int64(0),
	}

	if comment.Grade != -1 {

		asset.Rate[0] = asset.Rate[0] + comment.Grade
		asset.Rate[1] = asset.Rate[1] + 1

		if err := assetPut(stub, asset); err != nil {
			return shim.Error("Save asset error: " + err.Error())
		}
	}

	if err := commentPut(stub, comment); err != nil {
		return shim.Error("Add comment error: " + err.Error())
	}

	// 用户和评论关联键
	ucKey, err := stub.CreateCompositeKey("user~comment", []string{comment.Reviewer, comment.Hash})
	if err != nil {
		return shim.Error("Create key error: " + err.Error())
	}

	// 用户和评论关联键值为资产哈希-时间戳
	if err := stub.PutState(ucKey, []byte(comment.AssetHash+"-"+strconv.FormatInt(txtime, 10))); err != nil {
		return shim.Error("Associate user and comment error: " + err.Error())
	}

	// 资产和评论关联键
	acKey, err := stub.CreateCompositeKey("asset~comment", []string{comment.AssetHash, comment.Hash})
	if err != nil {
		return shim.Error("Create key error: " + err.Error())
	}

	// 用资产和评论关联键值为用户名-时间戳
	if err := stub.PutState(acKey, []byte(comment.Reviewer+"-"+strconv.FormatInt(txtime, 10))); err != nil {
		return shim.Error("Associate asset and comment error: " + err.Error())
	}

	if ok := recordPut(stub, txid, comment.Reviewer, []string{asset.Hash}, "commentAdd", comment.Date, ""); !ok {
		return shim.Error("Put record failure")
	}

	if ok := setEvent(stub, txid, user.Name, "commentAdd", "OK", txtime); !ok {
		return shim.Error("Triggle event failure")
	}
	// TODO 上传是否奖励积分

	//logger.Info("OK")

	return shim.Success(nil)
}

func commentRateSet(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("commentRateSet")

	if res := paramCheck(args, 2); res != "" {
		return shim.Error(res)
	}

	_, role, _, _, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	if role != "admin" {
		return shim.Error("Only admin can do this")
	}

	comment, err := commentGet(stub, args[0])
	if err != nil || comment == nil {
		return shim.Error("Comment not found")
	}

	if comment.Status != 0 {
		return shim.Error("Comment is invalid")
	}

	if args[1] == "agree" {
		comment.Rate[0] = comment.Rate[0] + 1
	} else {
		comment.Rate[1] = comment.Rate[1] + 1
	}

	if err := commentPut(stub, comment); err != nil {
		return shim.Error("Save comment error: " + err.Error())
	}

	return shim.Success(nil)
}

// 废弃comment(废弃的comment将不可见)
func commentStatusSet(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("commentStatusSet")

	if res := paramCheck(args, 2); res != "" {
		return shim.Error(res)
	}

	mspid, role, txid, txtime, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	if role != "admin" {
		return shim.Error("Only admin can do this")
	}

	chashs := strings.Replace(args[0], " ", "", -1)

	status, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return shim.Error("Status value is illegal")
	}

	for _, commentHash := range strings.Split(chashs, ",") {
		comment, err := commentGet(stub, commentHash)
		if err != nil || comment == nil {
			return shim.Error("Comment " + commentHash + " no found")
		}
		if strings.HasSuffix(comment.Reviewer, mspid) {
			return shim.Error("The owner  of comment " + commentHash + " is not belong to you")
		}

		comment.Status = status

		if err := commentPut(stub, comment); err != nil {
			return shim.Error("Set comment status error: " + err.Error())
		}
	}

	if ok := recordPut(stub, txid, role+"@"+mspid, args, "commentStatusSet", txtime, ""); !ok {
		return shim.Error("Put record failure")
	}

	if ok := setEvent(stub, txid, role+"@"+mspid, "commentStatusSet", "OK", txtime); !ok {
		return shim.Error("Triggle event failure")
	}

	return shim.Success(nil)
}

func commentQuery(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("commentQuery")

	var buffer bytes.Buffer

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for _, commentHash := range args {

		commentBytes, err := stub.GetState("comment:" + commentHash)
		if err != nil || commentBytes == nil || len(commentBytes) == 0 {
			continue
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString(string(commentBytes))
		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

// 查询用户的评论
func userCommentQuery(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("userCommentQuery")

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
		resultIterator, err := stub.GetStateByPartialCompositeKey("user~comment", []string{args[0]})
		if err != nil {
			return shim.Error("Query user comment error: " + err.Error())
		}
		defer resultIterator.Close()

		buffer = constructQueryResponseFromIterator(stub, resultIterator)
	} else {
		pagesize, err := strconv.ParseInt(psizestr, 10, 32)
		if err != nil {
			return shim.Error("Query user comment error: " + err.Error())
		}
		resultsIterator, responseMetadata, err := stub.GetStateByPartialCompositeKeyWithPagination("user~comment", []string{args[0]}, int32(pagesize), bookmark)
		if err != nil {
			return shim.Error("Query user comment error: " + err.Error())
		}
		defer resultsIterator.Close()

		buffer = constructQueryResponseFromIterator(stub, resultsIterator)

		addPaginationMetadataToQueryResults(buffer, responseMetadata)
	}

	return shim.Success(buffer.Bytes())
}

// 查询资产的评论
func assetCommentQuery(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("assetCommentQuery")

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

	asset, err := assetGet(stub, args[0])
	if err != nil || asset == nil {
		return shim.Error("Asset not found")
	}

	if psizestr == "" {
		resultIterator, err := stub.GetStateByPartialCompositeKey("asset~comment", []string{asset.Hash})
		if err != nil {
			return shim.Error("Query asset comment error: " + err.Error())
		}
		defer resultIterator.Close()

		buffer = constructQueryResponseFromIterator(stub, resultIterator)
	} else {
		pagesize, err := strconv.ParseInt(psizestr, 10, 32)
		if err != nil {
			return shim.Error("Query asset comment error: " + err.Error())
		}
		resultsIterator, responseMetadata, err := stub.GetStateByPartialCompositeKeyWithPagination("asset~comment", []string{asset.Hash}, int32(pagesize), bookmark)
		if err != nil {
			return shim.Error("Query asset comment error: " + err.Error())
		}
		defer resultsIterator.Close()

		buffer = constructQueryResponseFromIterator(stub, resultsIterator)

		addPaginationMetadataToQueryResults(buffer, responseMetadata)
	}

	return shim.Success(buffer.Bytes())
}

// 通过asset hash查询多个asset的详细信息
func assetQuery(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("assetQuery")

	var buffer bytes.Buffer

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for _, assetHash := range args {

		assetBytes, err := stub.GetState("asset:" + assetHash)
		if err != nil || assetBytes == nil || len(assetBytes) == 0 {
			continue
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString(string(assetBytes))
		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
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

// 查询用户上传和购买的所有asset
func userAssetQuery(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("userAssetQuery")

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
		resultIterator, err := stub.GetStateByPartialCompositeKey("user~asset", []string{args[0]})
		if err != nil {
			return shim.Error("Query user asset error: " + err.Error())
		}
		defer resultIterator.Close()

		buffer = constructQueryResponseFromIterator(stub, resultIterator)
	} else {
		pagesize, err := strconv.ParseInt(psizestr, 10, 32)
		if err != nil {
			return shim.Error("Query user asset error: " + err.Error())
		}
		resultsIterator, responseMetadata, err := stub.GetStateByPartialCompositeKeyWithPagination("user~asset", []string{args[0]}, int32(pagesize), bookmark)
		if err != nil {
			return shim.Error("Query user asset error: " + err.Error())
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

// 返回用户是否可访问资源
func userAssetAccess(stub shim.ChaincodeStubInterface, asset *Asset, user *User) bool {
	if asset == nil || user == nil {
		return false
	}

	if asset.Owner == user.Name {
		return true
	}

	if asset.Status != 0 && asset.Status != 1 {
		return false
	}

	uaKey, _ := stub.CreateCompositeKey("user~asset", []string{user.Name, asset.Hash})
	uaBytes, err := stub.GetState(uaKey)
	if err != nil || uaBytes == nil || len(uaBytes) == 0 {
		return false
	}

	// 已经是资源的购买者或所有者
	return true
}

// 查询请求发起者对指定资源的访问权限
func userAssetVerify(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("userAssetVerify")

	if res := paramCheck(args, 2); res != "" {
		return shim.Error(res)
	}

	user, err := userGet(stub, args[0])
	if err != nil {
		return shim.Error("User not found")
	}

	asset, err := assetGet(stub, args[1])
	if err != nil || asset == nil {
		return shim.Error("Asset not found")
	}

	res := userAssetAccess(stub, asset, user)

	return shim.Success([]byte(strconv.FormatBool(res)))
}

// 查询用户积分
func balanceQuery(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("balanceQuery")

	if res := paramCheck(args, 1); res != "" {
		return shim.Error(res)
	}

	user, err := userGet(stub, args[0])
	if err != nil {
		return shim.Error("User not found")
	}

	return shim.Success([]byte(strconv.FormatInt(user.Balance, 0)))
}

// 查询用户信息
func userQuery(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("userQuery")

	if res := paramCheck(args, 1); res != "" {
		return shim.Error(res)
	}

	user, err := userGet(stub, args[0])
	if err != nil {
		return shim.Error("User not found")
	}

	/*
		user, timestr, _, err := stubParse(stub)

		if err != nil {
			return shim.Error("Parse stub error: " + err.Error())
		}
		var buffer bytes.Buffer

		buffer.WriteString("{")

		buffer.WriteString("\"time\":")
		buffer.WriteString("\"")
		buffer.WriteString(timestr)
		buffer.WriteString("\"")
		buffer.WriteString(",")

		buffer.WriteString("\"name\":")
		buffer.WriteString("\"")
		buffer.WriteString(user.Name)
		buffer.WriteString("\"")
		buffer.WriteString(",")

		buffer.WriteString("\"mspid\":")
		buffer.WriteString("\"")
		buffer.WriteString(user.Mspid)
		buffer.WriteString("\"")
		buffer.WriteString(",")

		buffer.WriteString("\"vest\":")
		buffer.WriteString("\"")
		buffer.WriteString(user.Vest)
		buffer.WriteString("\"")
		buffer.WriteString(",")

		buffer.WriteString("\"posts\":")
		buffer.WriteString("\"")
		buffer.WriteString(user.Posts[0] + "-" + user.Posts[1] + "-" + user.Posts[2] + "-" + user.Posts[3])
		buffer.WriteString("\"")
		buffer.WriteString(",")

		buffer.WriteString("\"role\":")
		buffer.WriteString("\"")
		buffer.WriteString(user.Role)
		buffer.WriteString("\"")
		buffer.WriteString(",")

		buffer.WriteString("\"balance\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.Itoa(user.Balance))
		buffer.WriteString("\"")

		buffer.WriteString("}")

		return shim.Success(buffer.Bytes())
	*/

	userBytes, err := json.Marshal(user)
	if err != nil {
		return shim.Error("Marshal user error: " + err.Error())
	}

	return shim.Success(userBytes)
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

// 通过富查询语句进行查询
func richQueryResult(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("richQueryResult")

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

	if args[0] == "" {
		return shim.Error("The first argument can not be empty string")
	}

	if !strings.HasPrefix(args[0], "selector") {
		return shim.Error("The first argument must be select statement")
	}

	if psizestr == "" {
		resultIterator, err := stub.GetQueryResult(args[0])
		if err != nil {
			return shim.Error("Rich query error: " + err.Error())
		}
		defer resultIterator.Close()

		buffer = constructQueryResponseFromIterator(stub, resultIterator)
	} else {
		pagesize, err := strconv.ParseInt(psizestr, 10, 32)
		if err != nil {
			return shim.Error("Rich query error: " + err.Error())
		}
		resultsIterator, responseMetadata, err := stub.GetQueryResultWithPagination(args[0], int32(pagesize), bookmark)
		if err != nil {
			return shim.Error("Rich query error: " + err.Error())
		}
		defer resultsIterator.Close()

		buffer = constructQueryResponseFromIterator(stub, resultsIterator)

		addPaginationMetadataToQueryResults(buffer, responseMetadata)
	}
	return shim.Success(buffer.Bytes())
}

func (t *AssetManage) Init(stub shim.ChaincodeStubInterface) peer.Response {
	/*
			_, args := stub.GetFunctionAndParameters()
			if len(args) != 1 {
				return shim.Error("Incorrect number of arguments. Expecting 1")
		    }
	*/
	return shim.Success(nil)
}

func (t *AssetManage) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "assetUpload":
		return assetUpload(stub, args)
	case "assetBuy":
		return assetBuy(stub, args)
	case "assetTimesCount":
		return assetTimesCount(stub, args)
	case "assetStatusSet":
		return assetStatusSet(stub, args)
	case "commentAdd":
		return commentAdd(stub, args)
	case "commentRateSet":
		return commentRateSet(stub, args)
	case "commentStatusSet":
		return commentStatusSet(stub, args)
	case "userInit":
		return userInit(stub, args)
	case "balanceTransfer":
		return balanceTransfer(stub, args)
	case "balanceDeposit":
		return balanceDeposit(stub, args)
	case "assetQuery":
		return assetQuery(stub, args) //queryAsset(stub, args)
	case "commentQuery":
		return commentQuery(stub, args)
	case "userCommentQuery":
		return userCommentQuery(stub, args)
	case "assetCommentQuery":
		return assetCommentQuery(stub, args)
	case "richQueryResult":
		return richQueryResult(stub, args) //queryAssetHistory(stub, args)
	case "recordQuery":
		return recordQuery(stub, args)
	case "userRecordQuery":
		return userRecordQuery(stub, args)
	case "userAssetQuery":
		return userAssetQuery(stub, args)
	case "userAssetVerify":
		return userAssetVerify(stub, args)
	case "balanceQuery":
		return balanceQuery(stub, args)
	case "userQuery":
		return userQuery(stub, args)
	case "historyQuery":
		return historyQuery(stub, args)
	case "fullQuery":
		return fullQuery(stub, args)
	default:
		return shim.Error("Unsupported function: " + function)
	}

	// return shim.Success(nil)
}

var logger = shim.NewLogger("ebclog")

func main() {

	logger.SetLevel(shim.LogDebug)

	err := shim.Start(new(AssetManage))
	if err != nil {
		logger.Error("Error starting AssetManage chaincode: %s", err)
	}
}
