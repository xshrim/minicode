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
	Name    string    `json:"name"`    // 用户名
	Mspid   string    `json:"mspid"`   // 用户组织
	Vest    string    `json:"vest"`    // 用户归属(集团.公司.部门.团队)
	Posts   [4]string `json:"posts"`   // 用户职位(集团, 公司, 部门, 团队)
	Role    string    `json:"role"`    // 用户角色
	Balance int       `json:"balance"` // 用户积分余额
}

type Asset struct {
	Hash string `json:"hash"` // 资源哈希
	// Name  string            `json:"name"`  // 资源名
	Link  string            `json:"link"`  // 资源链接
	Owner string            `json:"owner"` // 资源上传者
	Mspid string            `json:"mspid"` // 资源组织
	Vest  string            `json:"vest"`  // 资源归属(集团.公司.部门.团队)
	Times [2]int            `json:"times"` // 资源浏览, 下载次数
	Rate  [2]int            `json:"rate"`  // 资源评分(评分总数/评分次数)
	Perm  [5]string         `json:"perm"`  // 资源访问权限(联盟内, 集团内, 公司内, 部门内, 团队内)
	Plist map[string]string `json:"plist"` // 资源访问特殊权限名单
	Date  string            `json:"date"`  // 资源上传时间戳
	State string            `json:"state"` // 资源状态(正常, 冻结, 不可用)
}

type Comment struct {
	Hash      string `json:"hash"`      // 评论哈希
	AssetHash string `json:"assethash"` // 对应资源哈希
	// Content    string `json:"content"`    // 评论正文
	Reviewer   string `json:"reviewer"`   // 评论者
	Floor      string `json:"floor"`      // 评论楼层
	Rate       [2]int `json:"rate"`       // 评论赞成反对数(赞成数/反对数)
	AssetGrade string `json:"assetgrade"` // 对资源的评分
	Date       string `json:"date"`       // 评论时间戳
	State      string `json:"state"`      // 评论状态

}

type Record struct {
	TxID    string   `json:"txid"`    // 交易ID
	Subject string   `json:"subject"` // 行为主体
	Object  []string `json:"object"`  // 行为客体
	Action  string   `json:"action"`  // 行为
	Date    string   `json:"date"`    // 行为时间戳
	Notes   string   `json:"notes"`   // 行为备注
}

type Event struct { // Fabric的链码事件默认会返回区块号, 交易ID, 链码名, 事件名和事件内容, 因此事件内容也可以不必包含下述部分信息
	TxID     string `json:"txid"`
	UserName string `json:"username"`
	FcName   string `json:"fcname"`
	Message  string `json:"message"`
	Date     string `json:"date"`
}

// 判断leader是否是staff的上级(组织管理员, 同级领导或上一级领导)
func isLeader(staff string, leader *User) bool {
	if permSwitch == "off" {
		return true
	}

	staffVest, err := getVest(staff)
	if err != nil {
		return false
	}

	if staffVest == leader.Vest || strings.HasPrefix(staffVest, leader.Vest+".") {
		if strings.Split(leader.Name, "@")[0] == "admin" {
			return true
		}

		idx := strings.Count(staffVest, ".")
		if leader.Posts[idx] == "leader" {
			return true
		}
		if idx > 0 && leader.Posts[idx-1] == "leader" {
			return true
		}
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

// 解析stub获取交易请求时间戳, 交易请求者mspid和userid
func stubParse(stub shim.ChaincodeStubInterface) (*User, string, string, error) {

	var timestr string

	user := &User{
		Name:    "",
		Mspid:   "",
		Posts:   [4]string{"staff", "staff", "staff", "staff"},
		Role:    "",
		Balance: 0,
	}

	txid := stub.GetTxID()

	timestamp, err := stub.GetTxTimestamp()
	if err != nil {
		timestr = ""
	} else {
		// 毫秒级时间戳
		timestr = strconv.Itoa(int(timestamp.Seconds*1000 + int64(timestamp.Nanos/1000000)))

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
		return user, timestr, txid, errors.New("Certificate not found")
	}

	// 使用正则表达式过滤掉creatorBytes中mspid部分的不可见字符(protobuf编码)
	user.Mspid = regexp.MustCompile(`\w+`).FindString(string(creatorBytes[:certStart]))

	// 直接从creatorBytes获取mspid
	//msp := creatorBytes[2 : certStart-3] // mspid两端存在不可见字符

	certText := creatorBytes[certStart:]

	// 获取证书的pem格式编码
	bl, _ := pem.Decode(certText)
	if bl == nil {
		return user, timestr, txid, errors.New("Certificate decode error")
	}

	// 生成证书对象
	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return user, timestr, txid, errors.New("Certificate parse error")
	}

	// 获取证书中包含的各种数据
	user.Name = cert.Subject.CommonName // 用户名, 如Admin@org1.example.com

	user.Vest, err = getVest(user.Name)
	if err != nil {
		user.Vest = ""
	}

	buff, err := getAttributesFromCert(cert)
	if err != nil {
		return user, timestr, txid, nil
	}

	attrs := &Attributes{}
	if buff != nil {
		err := json.Unmarshal(buff, attrs)
		if err != nil {
			return user, timestr, txid, nil
		}
	}

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

	if v, ok := attrs.Attrs["role"]; ok {
		user.Role = v
	}

	user.Balance, err = balanceGet(stub, user.Name)
	if err != nil {
		return user, timestr, txid, err
	}

	return user, timestr, txid, nil
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
func setEvent(stub shim.ChaincodeStubInterface, txid, username, fcname, message, date string) bool {
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

	logger.Info(err)

	return true
}

func recordPut(stub shim.ChaincodeStubInterface, txid, subject string, object []string, action, date, notes string) bool {
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

	if err := stub.PutState(urKey, []byte(record.Action+"-"+record.Date)); err != nil {
		return false
	}

	return true
}

// 获取余额
func balanceGet(stub shim.ChaincodeStubInterface, uname string) (int, error) {

	balanceBytes, err := stub.GetState("user:" + uname)
	if err == nil && balanceBytes != nil && len(balanceBytes) > 0 {
		return strconv.Atoi(string(balanceBytes))
	}
	return 0, errors.New("User not found")
}

// 增减余额
func balanceAdd(stub shim.ChaincodeStubInterface, uname string, value int) (int, error) {
	balance, err := balanceGet(stub, uname)
	if err != nil {
		return 0, err
	}

	balance += value

	if balance < 0 {
		return -1, errors.New("Balance is less than 0")
	}

	if err := stub.PutState("user:"+uname, []byte(strconv.Itoa(balance))); err != nil {
		return -1, errors.New("Add balance error: " + err.Error())
	}

	return balance, nil
}

func balanceTransfer(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("balanceTransfer")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	if args[0] == "" || args[1] == "" {
		return shim.Error("Arguments can not be empty string")
	}

	user, timestr, txid, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	value, err := strconv.Atoi(args[1])
	if err != nil || value < 0 {
		return shim.Error("Balance value is illegal")
	}

	if value <= 0 {
		return shim.Error("Balance value must more than 0")
	}

	if value > user.Balance {
		return shim.Error("Balance not enough")
	}

	if _, err := balanceAdd(stub, user.Name, 0-value); err != nil {
		return shim.Error("Balance transform failed")
	}

	if _, err := balanceAdd(stub, args[0], value); err != nil {
		return shim.Error("Balance transform failed")
	}

	if ok := recordPut(stub, txid, user.Name, []string{args[0]}, "balanceTransfer", timestr, "value:"+args[1]); !ok {
		return shim.Error("Put record failure")
	}

	if ok := setEvent(stub, txid, user.Name, "balanceTransfer", "hello", timestr); !ok {
		return shim.Error("Triggle event failure")
	}

	return shim.Success(nil)
}

// 积分充值
func balanceRecharge(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("balanceRecharge")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	if args[0] == "" || args[1] == "" {
		return shim.Error("Arguments can not be empty string")
	}

	user, timestr, txid, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	value, err := strconv.Atoi(args[1])
	if err != nil || value < 0 {
		return shim.Error("Balance value is illegal")
	}

	// 指定特定用户可以为其他用户充值 (仅组织管理员可以为其下的用户充值)
	uprefix := strings.Split(args[0], "@")[0]
	//if strings.HasPrefix(user.Name, "Admin@") {
	if strings.Replace(args[0], uprefix, "admin", 1) == user.Name {
		if balance, err := balanceAdd(stub, args[0], value); err == nil && balance >= 0 {
			if ok := recordPut(stub, txid, user.Name, []string{args[0]}, "balanceRecharge", timestr, "value:"+strconv.Itoa(value)); !ok {
				return shim.Error("Put record failure")
			}

			if ok := setEvent(stub, txid, user.Name, "balanceRecharge", "hello", timestr); !ok {
				return shim.Error("Triggle event failure")
			}

			return shim.Success([]byte(strconv.Itoa(balance)))
		}
	}

	return shim.Error("Balance recharge failure")
}

func userInit(stub shim.ChaincodeStubInterface) peer.Response {

	logger.Info("userInit")

	user, timestr, txid, err := stubParse(stub)
	if err != nil && err.Error() != "User not found" {
		return shim.Error("Parse stub error: " + err.Error())
	}

	if err == nil {
		return shim.Error("User already exist")
	}

	initBalance := 100

	if err := stub.PutState("user:"+user.Name, []byte(strconv.Itoa(initBalance))); err != nil {
		return shim.Error("Add balance error: " + err.Error())
	}

	if ok := recordPut(stub, txid, user.Name, []string{user.Name}, "userInit", timestr, "value:"+strconv.Itoa(initBalance)); !ok {
		return shim.Error("Put record failure")
	}

	if ok := setEvent(stub, txid, user.Name, "userInit", "hello", timestr); !ok {
		return shim.Error("Triggle event failure")
	}

	logger.Info("OKOK")

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

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	if args[0] == "" || args[1] == "" {
		return shim.Error("The first two arguments can not be empty string")
	}

	assetBytes, err := stub.GetState("asset:" + args[0])
	if err == nil && assetBytes != nil && len(assetBytes) > 0 {
		return shim.Error("Asset already exist")
	}

	// nextid := IDCounter(stub, "asset")
	user, timestr, txid, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	asset := &Asset{
		Hash: args[0],
		// Name:  args[1],
		Link:  args[1],
		Owner: user.Name,
		Mspid: user.Mspid,
		Vest:  user.Vest,
		Times: [2]int{0, 0},
		Rate:  [2]int{0, 0},
		Perm:  [5]string{"-1", "-1", "-1", "-1", "-1"},
		Plist: make(map[string]string),
		Date:  timestr,
		State: "NORMAL",
	}

	// 将"-1,-1,300,100,0"字符串的访问权限参数转换为数组形式, 范围从大到小, 缺省值为"-1"
	if args[2] != "" {
		for idx, val := range strings.Split(args[2], ",") {
			asset.Perm[idx] = strings.TrimSpace(val)
		}
	}

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

	if err := assetPut(stub, asset); err != nil {
		return shim.Error("Save asset error: " + err.Error())
	}

	if _, err := balanceAdd(stub, user.Name, 5); err != nil {
		return shim.Error("Add balance error: " + err.Error())
	}

	// 用户和资产关联键
	uaKey, err := stub.CreateCompositeKey("user~asset", []string{asset.Owner, asset.Hash})
	if err != nil {
		return shim.Error("Create key error: " + err.Error())
	}

	// 用户和资产关联键值为关系-时间戳
	if err := stub.PutState(uaKey, []byte("owner-"+timestr)); err != nil {
		return shim.Error("Associate user and asset error: " + err.Error())
	}

	if ok := recordPut(stub, txid, asset.Owner, []string{asset.Hash}, "assetUpload", asset.Date, ""); !ok {
		return shim.Error("Put record failure")
	}

	if ok := setEvent(stub, txid, user.Name, "assetUpload", "hello", timestr); !ok {
		return shim.Error("Triggle event failure")
	}

	return shim.Success(nil)
}

// 冻结asset(被冻结的asset将不可购买)
func assetFreeze(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("assetFreeze")

	user, timestr, txid, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	for _, assetHash := range args {
		asset, err := assetGet(stub, assetHash)
		if err != nil || asset == nil {
			return shim.Error("Asset " + assetHash + " no found")
		}
		if user.Name != asset.Owner && !isLeader(asset.Owner, user) {
			return shim.Error("User is not the owner or leader of asset " + assetHash)
		}
		asset.State = "Freezed"

		if err := assetPut(stub, asset); err != nil {
			return shim.Error("Freeze asset error: " + err.Error())
		}
	}

	if ok := recordPut(stub, txid, user.Name, args, "assetFreeze", timestr, ""); !ok {
		return shim.Error("Put record failure")
	}

	if ok := setEvent(stub, txid, user.Name, "assetFreeze", "hello", timestr); !ok {
		return shim.Error("Triggle event failure")
	}

	return shim.Success(nil)
}

// 废弃asset(被废弃的asset将不可购买, 浏览和下载)
func assetDiscard(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("assetDiscard")

	user, timestr, txid, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	for _, assetHash := range args {
		asset, err := assetGet(stub, assetHash)
		if err != nil || asset == nil {
			return shim.Error("Asset " + assetHash + " no found")
		}
		if user.Name != asset.Owner && !isLeader(asset.Owner, user) {
			return shim.Error("User is not the owner or leader of asset " + assetHash)
		}
		asset.State = "Discard"

		if err := assetPut(stub, asset); err != nil {
			return shim.Error("Discard asset error: " + err.Error())
		}
	}

	if ok := recordPut(stub, txid, user.Name, args, "assetDiscard", timestr, ""); !ok {
		return shim.Error("Put record failure")
	}

	if ok := setEvent(stub, txid, user.Name, "assetDiscard", "hello", timestr); !ok {
		return shim.Error("Triggle event failure")
	}

	return shim.Success(nil)
}

// 购买asset
func assetBuy(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("assetBuy")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	if args[0] == "" {
		return shim.Error("Argument can not be empty string")
	}

	user, timestr, txid, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	asset, err := assetGet(stub, args[0])
	if err != nil || asset == nil {
		return shim.Error("Asset not found")
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

	perm, reason := userAssetPerm(stub, asset, user)

	value, err := strconv.Atoi(perm)
	if err != nil {
		return shim.Error("Asset value is illegal")
	}

	if value < 0 {
		return shim.Error("Asset can not be bought by the user")
	}

	if value == 0 {
		return shim.Error("Asset can be accessed by the user already. reason: " + reason)
	}

	if user.Balance < value {
		return shim.Error("Balance is insufficient")
	}

	if _, err = balanceAdd(stub, user.Name, 0-value); err != nil {
		return shim.Error("Payment faillure")
	}

	if _, err = balanceAdd(stub, asset.Owner, value); err != nil {
		return shim.Error("Payment faillure")
	}

	uaKey, err := stub.CreateCompositeKey("user~asset", []string{user.Name, asset.Hash})
	if err != nil {
		return shim.Error("Create key error: " + err.Error())
	}

	if err := stub.PutState(uaKey, []byte("user-"+timestr)); err != nil {
		return shim.Error("Associate user and asset error: " + err.Error())
	}

	if ok := recordPut(stub, txid, user.Name, []string{asset.Hash}, "assetBuy", timestr, "value:"+perm); !ok {
		return shim.Error("Put record failure")
	}

	if ok := setEvent(stub, txid, user.Name, "assetBuy", "hello", timestr); !ok {
		return shim.Error("Triggle event failure")
	}

	return shim.Success(nil)
}

func assetTimesCount(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("assetTimesCount")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	if args[0] == "" || args[1] == "" {
		return shim.Error("Arguments can not be empty string")
	}

	user, timestr, txid, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	asset, err := assetGet(stub, args[0])
	if err != nil || asset == nil {
		return shim.Error("Asset not found")
	}

	if asset.State == "Discard" {
		return shim.Error("Asset is discarded")
	}

	if args[1] == "view" {
		asset.Times[0] = asset.Times[0] + 1
	} else {
		asset.Times[1] = asset.Times[1] + 1
	}

	// TODO 是否需要给上传者积分奖励(奖励可能导致恶意刷新)

	if err := assetPut(stub, asset); err != nil {
		return shim.Error("Count asset error: " + err.Error())
	}

	// 给上传者奖励积分
	if asset.Owner != user.Name {
		if _, err := balanceAdd(stub, asset.Owner, 1); err != nil {
			return shim.Error("Add balance error: " + err.Error())
		}
	}

	if ok := recordPut(stub, txid, user.Name, []string{asset.Hash}, "assetTimesCount", timestr, args[1]+":1"); !ok {
		return shim.Error("Put record failure")
	}

	if ok := setEvent(stub, txid, user.Name, "assetTimesCount", "hello", timestr); !ok {
		return shim.Error("Triggle event failure")
	}

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

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	if args[0] == "" || args[1] == "" {
		return shim.Error("The first two arguments can not be empty string")
	}

	commentBytes, err := stub.GetState("comment:" + args[0])
	if err == nil && commentBytes != nil && len(commentBytes) > 0 {
		return shim.Error("Comment already exist")
	}

	asset, err := assetGet(stub, args[1])
	if err != nil || asset == nil {
		return shim.Error("Asset not found")
	}

	// nextid := IDCounter(stub, "asset")
	user, timestr, txid, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	perm, _ := userAssetPerm(stub, asset, user)
	if perm != "0" {
		return shim.Error("Comment for the asset not allowed by the user")
	}

	comment := &Comment{
		Hash: args[0],
		// Name:  args[1],
		AssetHash:  args[1],
		Reviewer:   user.Name,
		Floor:      args[2],
		Rate:       [2]int{0, 0},
		AssetGrade: args[3],
		Date:       timestr,
		State:      "NORMAL",
	}

	if comment.AssetGrade != "" {
		grade, err := strconv.Atoi(comment.AssetGrade)
		if err != nil {
			return shim.Error("Illegal argument: " + err.Error())
		}

		asset.Rate[0] = asset.Rate[0] + grade
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
	if err := stub.PutState(ucKey, []byte(comment.AssetHash+"-"+timestr)); err != nil {
		return shim.Error("Associate user and comment error: " + err.Error())
	}

	// 资产和评论关联键
	acKey, err := stub.CreateCompositeKey("asset~comment", []string{comment.AssetHash, comment.Hash})
	if err != nil {
		return shim.Error("Create key error: " + err.Error())
	}

	// 用资产和评论关联键值为用户名-时间戳
	if err := stub.PutState(acKey, []byte(comment.Reviewer+"-"+timestr)); err != nil {
		return shim.Error("Associate asset and comment error: " + err.Error())
	}

	if ok := recordPut(stub, txid, comment.Reviewer, []string{asset.Hash}, "commentAdd", comment.Date, ""); !ok {
		return shim.Error("Put record failure")
	}

	if ok := setEvent(stub, txid, user.Name, "commentAdd", "hello", timestr); !ok {
		return shim.Error("Triggle event failure")
	}
	// TODO 上传是否奖励积分

	return shim.Success(nil)
}

func commentRateSet(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("commentRateSet")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	if args[0] == "" || args[1] == "" {
		return shim.Error("Arguments can not be empty string")
	}

	comment, err := commentGet(stub, args[0])
	if err != nil || comment == nil {
		return shim.Error("Comment not found")
	}

	if comment.State == "Discard" {
		return shim.Error("Comment is Discarded")
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
func commentDiscard(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("commentDiscard")

	user, timestr, txid, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	for _, commentHash := range args {
		comment, err := commentGet(stub, commentHash)
		if err != nil || comment == nil {
			return shim.Error("Comment " + commentHash + " no found")
		}
		if user.Name != comment.Reviewer && !isLeader(comment.Reviewer, user) {
			return shim.Error("User is not the owner or leader of comment " + commentHash)
		}
		if comment.State == "Discard" {
			return shim.Error("Comment has been discarded")
		}

		comment.State = "Discard"

		if err := commentPut(stub, comment); err != nil {
			return shim.Error("Discard comment error: " + err.Error())
		}
	}

	if ok := recordPut(stub, txid, user.Name, args, "commentDiscard", timestr, ""); !ok {
		return shim.Error("Put record failure")
	}

	if ok := setEvent(stub, txid, user.Name, "commentDiscard", "hello", timestr); !ok {
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

	if length > 2 {
		return shim.Error("Incorrect number of arguments. Expecting 0~2")
	}
	if length > 0 {
		psizestr = args[0]
		if length > 1 {
			bookmark = args[1]
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
	user, _, _, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	if psizestr == "" {
		resultIterator, err := stub.GetStateByPartialCompositeKey("user~comment", []string{user.Name})
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
		resultsIterator, responseMetadata, err := stub.GetStateByPartialCompositeKeyWithPagination("user~comment", []string{user.Name}, int32(pagesize), bookmark)
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
	if length > 1 {
		psizestr = args[1]
		if length > 2 {
			bookmark = args[2]
		}
	}

	if args[0] == "" {
		return shim.Error("The first argument can not be empty string")
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
	user, _, _, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	asset, err := assetGet(stub, args[0])
	if err != nil || asset == nil {
		return shim.Error("Asset not found")
	}

	perm, _ := userAssetPerm(stub, asset, user)
	if perm != "0" {
		return shim.Error("Not allowed to view comments for the asset")
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

	if length < 0 || length > 2 {
		return shim.Error("Incorrect number of arguments. Expecting 0~2")
	}
	if length > 0 {
		psizestr = args[0]
		if length > 1 {
			bookmark = args[1]
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
	user, _, _, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	if psizestr == "" {
		resultIterator, err := stub.GetStateByPartialCompositeKey("user~asset", []string{user.Name})
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
		resultsIterator, responseMetadata, err := stub.GetStateByPartialCompositeKeyWithPagination("user~asset", []string{user.Name}, int32(pagesize), bookmark)
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
	if length < 0 || length > 2 {
		return shim.Error("Incorrect number of arguments. Expecting 0~2")
	}
	if length > 0 {
		psizestr = args[0]
		if length > 1 {
			bookmark = args[1]
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
	user, _, _, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	if psizestr == "" {
		resultIterator, err := stub.GetStateByPartialCompositeKey("user~record", []string{user.Name})
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
		resultsIterator, responseMetadata, err := stub.GetStateByPartialCompositeKeyWithPagination("user~record", []string{user.Name}, int32(pagesize), bookmark)
		if err != nil {
			return shim.Error("Query user record error: " + err.Error())
		}
		defer resultsIterator.Close()

		buffer = constructQueryResponseFromIterator(stub, resultsIterator)

		addPaginationMetadataToQueryResults(buffer, responseMetadata)
	}

	return shim.Success(buffer.Bytes())
}

// 返回用户对资源的访问权限以及获得权限的原因
// -1: 无权访问; 0: 可访问; >1: 需购买访问
// nil: 参数为空; error: 解析出错, leader: 管理员或上级领导; perm: 访问控制; user: 已购买者; owner: 所有者
func userAssetPerm(stub shim.ChaincodeStubInterface, asset *Asset, user *User) (string, string) {
	if asset == nil || user == nil {
		return "-1", "nil"
	}

	if isLeader(asset.Owner, user) { // 是管理员或上级领导
		return "0", "leader"
	}

	if asset.State != "Normal" {
		if asset.Owner == user.Name {
			return "0", "owner-" + asset.Date
		}
		return "-1", asset.State
	}

	uaKey, _ := stub.CreateCompositeKey("user~asset", []string{user.Name, asset.Hash})
	uaBytes, err := stub.GetState(uaKey)
	if err != nil || uaBytes == nil || len(uaBytes) == 0 {
		i := 0
		uvest := strings.Split(user.Vest, ".")
		avest := strings.Split(asset.Vest, ".")
		for i < 4 {
			if uvest[i] != avest[i] {
				break
			}
			i++
		}

		perm := asset.Perm[i] // 表示该用户对该资源的权限值

		// TODO 考虑黑名单和白名单的影响
		return perm, "perm"
	}

	// 已经是资源的购买者或所有者
	return "0", string(uaBytes)
}

// 查询请求发起者对指定资源的访问权限
func userAssetVerify(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("userAssetVerify")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	user, _, _, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	asset, err := assetGet(stub, args[0])
	if err != nil || asset == nil {
		return shim.Error("Asset not found")
	}

	perm, reason := userAssetPerm(stub, asset, user)

	return shim.Success([]byte(perm + ":" + reason))
}

// 查询用户积分
func balanceQuery(stub shim.ChaincodeStubInterface) peer.Response {

	logger.Info("balanceQuery")

	user, _, _, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	return shim.Success([]byte(strconv.Itoa(user.Balance)))
}

// 查询用户信息
func userQuery(stub shim.ChaincodeStubInterface) peer.Response {

	logger.Info("userQuery")

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
}

// 键的历史记录
func historyQuery(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("historyQuery")

	var key string
	var buffer bytes.Buffer

	length := len(args)

	if length < 1 {
		return shim.Error("Incorrect number of arguments. Expecting > 1")
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
	if length > 1 {
		psizestr = args[1]
		if length > 2 {
			bookmark = args[2]
		}
	}

	if args[0] == "" {
		return shim.Error("The first argument can not be empty string")
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
	case "assetFreeze":
		return assetFreeze(stub, args)
	case "assetDiscard":
		return assetDiscard(stub, args)
	case "commentAdd":
		return commentAdd(stub, args)
	case "commentRateSet":
		return commentRateSet(stub, args)
	case "commentDiscard":
		return commentDiscard(stub, args)
	case "userInit":
		return userInit(stub)
	case "balanceTransfer":
		return balanceTransfer(stub, args)
	case "balanceRecharge":
		return balanceRecharge(stub, args)
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
		return balanceQuery(stub)
	case "userQuery":
		return userQuery(stub)
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
