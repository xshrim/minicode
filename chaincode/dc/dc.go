package main

import (
	"bytes"
	"crypto/x509"
	"encoding/asn1"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type DC struct {
	Admin []string `json:"admin"` // 管理员ID(数字证书ski)
}

// Attributes contains attribute names and values
type CertAttrs struct {
	Attrs map[string]string `json:"attrs"`
}

type User struct {
	Uid     string `json:"uid"`     // 用户唯一标识(数字证书ski)
	Mspid   string `json:"mspid"`   // 用户组织ID
	Role    string `json:"role"`    // 用户角色
	Balance int64  `json:"balance"` // 用户余额
}

type Asset struct {
	Dci      string `json:"dci"`      // 数字版权唯一标识符
	Name     string `json:"name"`     // 数字版权名
	Category string `json:"category"` // 数字版权类别
	Describe string `json:"describe"` // 数字版权描述信息
	Content  string `json:"content"`  // 数字版权内容(可以是一系列电子材料的哈希)
	Author   string `json:"author"`   // 数字版权创作者ID(用户数字证书ski)
	Owner    string `json:"owner"`    // 数字版权所有者ID(用户数字证书ski)
	Sprice   int64  `json:"sprice"`   // 数字版权授权价格
	Tprice   int64  `json:"tprice"`   // 数字版权转让价格
	Ctime    int64  `json:"ctime"`    // 数字版权创建时间
	Ttime    int64  `json:"ttime"`    // 数字版权最后一次转让时间
	Etime    int64  `json:"etime"`    // 数字版权过期时间(-1:不过期)
	Status   int64  `json:"status"`   // 数字版权状态(0:正常; 1:废弃)
}

type Perm struct {
	Uid    string `json:"uid"`    // 被授权用户ID
	Dci    string `json:"dci"`    // 版权标识符
	Ptype  string `josn:"ptype"`  // 版权权限类型:AuthorAndOwner, Owner, AuthorAndUser, User, Author, Invalid
	Clause string `json:"clause"` // 版权权限条款细则
	Price  int64  `json:"price"`  // 版权权限实际授权价格
	Vtime  int64  `json:"vtime"`  // 版权权限生效时间
	Etime  int64  `json:"etime"`  // 版权权限过期时间(-1:不过期)
	Ptime  int64  `json:"ptime"`  // 版权权限授权时间
	Status int64  `json:"status"` // 版权权限状态(0:已批准; 1: 申请中; 2: 已撤销; 3: 已驳回; 4: 已失效)
}

type Record struct {
	Txid     string   `json:"txid"`     // 交易编号(区块链交易哈希)
	Sender   string   `json:"sender"`   // 交易发起者(用户数字证书ski)
	Receiver string   `json:"receiver"` // 交易接收者(用户数字证书ski)
	Object   []string `json:"object"`   // 交易对象
	Action   string   `json:"action"`   // 交易类型
	Time     int64    `json:"time"`     // 交易时间
	Value    int64    `json:"value"`    // 交易涉及金额
	Note     string   `json:"note"`     // 交易备注
}

type Event struct { // Fabric的链码事件默认会返回区块号, 交易ID, 链码名, 事件名和事件内容, 因此事件内容也可以不必包含下述部分信息
	Txid    string `json:"txid"`    // 交易哈希
	Uid     string `json:"uid"`     // 交易发起者
	Fname   string `json:"fname"`   // 调用函数名
	Message string `json:"message"` // 事件消息内容
	Time    int64  `json:"time"`    // 事件产生时间
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

// 获取结构体所有字段
func getFieldNames(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		//log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Name)
	}
	return result
}

// 获取结构体所有字段tag
func getTagNames(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		//log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		tagName := t.Field(i).Name
		tags := strings.Split(string(t.Field(i).Tag), "\"")
		if len(tags) > 1 {
			tagName = tags[1]
		}
		result = append(result, tagName)
	}
	return result
}

// 获取给定key字符串中的所有key, 识别组合key并拆分
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

		if bytes.Compare(value, []byte{0x00}) == 0 { // 表示value部分为空值 (当前仅user~perm组合键存在value为空的情况)
			keys := strings.Split(key, ":")
			if len(keys) != 2 {
				continue
			}
			permBytes, err := stub.GetState("perm:" + keys[1] + ":" + keys[0]) // 根据key中的参数获取perm的具体信息
			if err != nil || permBytes == nil || len(permBytes) == 0 {
				continue
			}
			value = permBytes
		}
		buffer.WriteString(", \"Value\":")
		// value is a JSON object, so we write as-is
		buffer.WriteString(string(value))
		buffer.WriteString("}")

		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")
	return &buffer
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

// 判断obj是否在target中，target支持的类型arrary,slice,map
func contain(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

// 字符串首字母大写, 其他字母小写
func firstToUpper(str string) string {
	if len(str) < 1 {
		return ""
	}
	rs := []rune(str)
	for i := range rs {
		if i == 0 {
			if rs[i] >= 97 && rs[i] <= 122 {
				rs[i] -= 32
			}
		} else {
			if rs[i] >= 65 && rs[i] <= 90 {
				rs[i] += 32
			}
		}
	}

	return string(rs)
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

// 解析stub获取交易请求者信息, 交易请求时间戳和交易请求哈希
func stubParse(stub shim.ChaincodeStubInterface) (*User, int64, string, error) {

	var txtime int64

	user := &User{
		Uid:     "",
		Mspid:   "",
		Role:    "",
		Balance: int64(0),
	}

	txid := stub.GetTxID()

	timestamp, err := stub.GetTxTimestamp()
	if err != nil {
		txtime = int64(0)
	} else {
		// 毫秒级时间戳
		// txtime = strconv.FormatInt(timestamp.Seconds*1000+int64(timestamp.Nanos/1000000), 10)

		// 秒级时间戳
		// txtime = strconv.FormatInt(timestamp.Seconds, 10)
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
		return user, txtime, txid, errors.New("Certificate not found")
	}

	// 使用正则表达式过滤掉creatorBytes中mspid部分的不可见字符(protobuf编码)
	user.Mspid = regexp.MustCompile(`\w+`).FindString(string(creatorBytes[:certStart]))

	// 直接从creatorBytes获取mspid
	//msp := creatorBytes[2 : certStart-3] // mspid两端存在不可见字符

	certText := creatorBytes[certStart:]

	// 获取证书的pem格式编码
	bl, _ := pem.Decode(certText)
	if bl == nil {
		return user, txtime, txid, errors.New("Certificate decode error")
	}

	// 生成证书对象
	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return user, txtime, txid, errors.New("Certificate parse error")
	}

	// 获取证书中包含的各种数据
	// user.Name = cert.Subject.CommonName // 用户名, 如Admin@org1.example.com

	// 获取证书中的SubjectKeyId, 是publickey的hash, ski类似以太坊中的用户地址
	// user.Uid = hex.EncodeToString(cert.SubjectKeyId) // 并非所有证书都包含此字段

	// SerialNumber是证书的唯一识别码
	// user.Uid = hex.EncodeToString(cert.SerialNumber.Bytes())

	// CommonName是向CA注册用户时指定的用户名
	cn := cert.Subject.CommonName
	user.Uid = strings.Split(cn, "@")[0] + "@" + user.Mspid

	// 三种Uid表示方式
	// 方式1: 不可取
	// 方式2: 可取, 唯一性高, 但可读性差, 需要外部保存实际用户到SerialNumber的映射关系
	// 方式3: 可取, 可读性高, 正常而言唯一性也可以保证

	buff, err := getAttributesFromCert(cert)
	if err != nil {
		return user, txtime, txid, nil
	}

	attrs := &CertAttrs{}
	if buff != nil {
		err := json.Unmarshal(buff, attrs)
		if err != nil {
			return user, txtime, txid, nil
		}
	}

	if v, ok := attrs.Attrs["role"]; ok {
		user.Role = v
	}

	user.Balance, err = balanceGet(stub, user.Uid)
	if err != nil {
		return user, txtime, txid, err
	}

	return user, txtime, txid, nil
}

// 发起链码事件
func setEvent(stub shim.ChaincodeStubInterface, txid, uid, fname, message string, txtime int64) bool {
	event := &Event{
		Txid:    txid,
		Uid:     uid,
		Fname:   fname,
		Message: message,
		Time:    txtime,
	}

	logger.Info(event)

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return false
	}

	err = stub.SetEvent(fname, eventBytes) // 函数名作为事件过滤条件
	if err != nil {
		return false
	}

	logger.Info(err)

	return true
}

// 权限时间相关检查
func permTimeCheck(asset *Asset, txtime, vtime, etime int64) error {
	if vtime < txtime {
		return fmt.Errorf("The effective time of perm must be larger than current time")
	}

	if etime <= txtime {
		return fmt.Errorf("The expire time of perm must be larger than current time")
	}

	if etime <= vtime {
		return fmt.Errorf("The expire time of perm must be larger than effective time")
	}

	if vtime > asset.Etime {
		return fmt.Errorf("The effective time of perm must be less than the expire time of asset")
	}

	if etime > asset.Etime {
		return fmt.Errorf("The expire time of perm must be less than the expire time of asset")
	}

	return nil
}

// 版权权限合法性检查
func permCheck(stub shim.ChaincodeStubInterface, user *User, asset *Asset, ptype, pricestr, vtimestr, etimestr string, txtime int64) (string, int64, int64, int64, error) {

	ptype = firstToUpper(ptype)

	if ptype != "Owner" && ptype != "User" { // 只允许请求这两种权限
		return "", 0, 0, 0, fmt.Errorf("Requested perm is not allowed")
	}

	if asset.Owner == user.Uid {
		return "", 0, 0, 0, fmt.Errorf("Owner of the asset do not need request perm of it")
	}

	if asset.Author == user.Uid {
		ptype = "AuthorAnd" + ptype
	}

	if asset.Etime <= txtime {
		return "", 0, 0, 0, fmt.Errorf("Asset cryptoright is expired")
	}

	if asset.Status != 0 {
		return "", 0, 0, 0, fmt.Errorf("Asset is invalid")
	}

	price, err := strconv.ParseInt(pricestr, 10, 64)
	if err != nil {
		return "", 0, 0, 0, fmt.Errorf("Illegal argument")
	}

	vtime, err := strconv.ParseInt(vtimestr, 10, 64)
	if err != nil {
		return "", 0, 0, 0, fmt.Errorf("Illegal argument")
	}
	etime, err := strconv.ParseInt(etimestr, 10, 64)
	if err != nil {
		return "", 0, 0, 0, fmt.Errorf("Illegal argument")
	}

	if vtime < txtime {
		vtime = txtime
	}

	if etime <= 0 {
		etime = asset.Etime
	}

	if price < 0 { // 原价购买
		if ptype == "Owner" {
			price = asset.Tprice
		} else {
			price = asset.Sprice
		}
	}

	if user.Balance < price {
		return "", 0, 0, 0, fmt.Errorf("Balance not enough")
	}

	if err := permTimeCheck(asset, txtime, vtime, etime); err != nil {
		return "", 0, 0, 0, err
	}

	return ptype, price, vtime, etime, nil
}

// 获取权限
func permGet(stub shim.ChaincodeStubInterface, dci, uid string) (*Perm, error) {

	permBytes, err := stub.GetState("perm:" + dci + ":" + uid)
	if err != nil || permBytes == nil || len(permBytes) == 0 {
		return nil, errors.New("Perm not found error: " + err.Error())
	}

	perm := new(Perm)
	if err := json.Unmarshal(permBytes, perm); err != nil {
		return nil, errors.New("Unmarshal perm error: " + err.Error())
	}

	return perm, nil
}

// 保存权限
func permPut(stub shim.ChaincodeStubInterface, uid, dci, ptype, clause string, price, vtime, etime, ptime, status int64) bool {

	perm := &Perm{
		Uid:    uid,
		Dci:    dci,
		Ptype:  ptype,
		Clause: clause,
		Price:  price,
		Vtime:  vtime,
		Etime:  etime,
		Ptime:  ptime,
		Status: status,
	}

	permBytes, err := json.Marshal(perm)
	if err != nil || permBytes == nil || len(permBytes) == 0 {
		return false
	}

	// key中包含授权时间是为了避免多次授权产生数据覆盖
	/*
			if err := stub.PutState("perm:"+perm.Dci+":"+perm.Uid+":"+strconv.FormatInt(perm.Ptime, 10), permBytes); err != nil { // 版权和哪些用户有关联
				return false
			}

			upKey, err := stub.CreateCompositeKey("user~perm~time", []string{perm.Uid, perm.Dci, strconv.FormatInt(perm.Ptime, 10)}) // 用户在哪些版权上有权限
			if err != nil {
				return false
		    }
	*/

	// 对同一用户进行同一版权的重复授权将导致旧的授权信息被覆盖, 这是允许的, 因为fabric会自行维护同一个key的全部历史变更记录
	if err := stub.PutState("perm:"+perm.Dci+":"+perm.Uid, permBytes); err != nil { // 版权和哪些用户有关联
		return false
	}

	// 即便权限处于申请状态, 也需要建立用户和权限的联系, 因为用户需要检索自己哪些权限还在申请中
	upKey, err := stub.CreateCompositeKey("user~perm", []string{perm.Uid, perm.Dci}) // 用户在哪些版权上有权限
	if err != nil {
		return false
	}

	// 因为权限信息已经保存在上一个kv里了, 此kv仅仅用于表示用户和权限之间的联系, 便于检索, 无需重复保存权限信息. 因此value是无意义的, 仅保存为null字符即可. value不能保存为nil, 因为value为nil会导致kv被删除
	if err := stub.PutState(upKey, []byte{0x00}); err != nil {
		return false
	}

	return true
}

// 保存日志
func recordPut(stub shim.ChaincodeStubInterface, txid, sender, receiver string, object []string, action string, txtime, value int64, note string) bool {
	/*
		stKey, err := stub.CreateCompositeKey("subject~timestamp", []string{subject, date})
		if err != nil {
			return false
		}
	*/

	record := &Record{
		Txid:     txid,
		Sender:   sender,
		Receiver: receiver,
		Object:   object,
		Action:   action,
		Time:     txtime,
		Value:    value,
		Note:     note,
	}

	recordBytes, err := json.Marshal(record)
	if err != nil || recordBytes == nil || len(recordBytes) == 0 {
		return false
	}

	if err := stub.PutState("record:"+txid, recordBytes); err != nil {
		return false
	}

	urKey, err := stub.CreateCompositeKey("user~record", []string{record.Sender, record.Txid})
	if err != nil {
		return false
	}

	if err := stub.PutState(urKey, []byte(record.Action+"-"+strconv.FormatInt(record.Time, 10))); err != nil {
		return false
	}

	return true
}

// 获取余额
func balanceGet(stub shim.ChaincodeStubInterface, uid string) (int64, error) {

	balanceBytes, err := stub.GetState("user:" + uid)
	if err == nil && balanceBytes != nil && len(balanceBytes) > 0 {
		return strconv.ParseInt(string(balanceBytes), 10, 64)
	}
	return 0, errors.New("User not found")
}

// 增减余额
func balanceAdd(stub shim.ChaincodeStubInterface, uid string, value int64) (int64, error) {
	balance, err := balanceGet(stub, uid)
	if err != nil {
		return 0, err
	}

	balance += value

	if balance < 0 {
		return -1, errors.New("Balance could not be less than 0")
	}

	if err := stub.PutState("user:"+uid, []byte(strconv.FormatInt(balance, 10))); err != nil {
		return -1, errors.New("Add balance error: " + err.Error())
	}

	return balance, nil
}

// 通过数字版权标识符获取版权对象
func assetGet(stub shim.ChaincodeStubInterface, dci string) (*Asset, error) {
	assetBytes, err := stub.GetState("asset:" + dci)
	if err != nil || assetBytes == nil || len(assetBytes) == 0 {
		return nil, errors.New("Asset not found error: " + err.Error())
	}

	asset := new(Asset)
	if err := json.Unmarshal(assetBytes, asset); err != nil {
		return nil, errors.New("Unmarshal asset error: " + err.Error())
	}

	return asset, nil
}

// 保存数字版权信息
func assetPut(stub shim.ChaincodeStubInterface, asset *Asset) error {
	assetBytes, err := json.Marshal(asset)
	if err != nil || assetBytes == nil || len(assetBytes) == 0 {
		return err
	}

	if err := stub.PutState("asset:"+asset.Dci, assetBytes); err != nil {
		return err
	}

	return nil
}

func trace(stub shim.ChaincodeStubInterface, txid, sender, receiver, object, action, note string, txtime, value int64) error {
	objs := strings.Split(object, ",")
	if ok := recordPut(stub, txid, sender, receiver, objs, action, txtime, value, note); !ok {
		return fmt.Errorf("Save record failure")
	}

	if ok := setEvent(stub, txid, sender, action, note, txtime); !ok {
		return fmt.Errorf("Triggle event failure")
	}

	return nil
}

// 用户初始化
func userInit(stub shim.ChaincodeStubInterface) peer.Response {

	logger.Info("userInit")

	user, txtime, txid, err := stubParse(stub)
	if err != nil && err.Error() != "User not found" {
		return shim.Error("Parse stub error: " + err.Error())
	}

	if err == nil {
		return shim.Error("User already exist")
	}

	initBalance := int64(0)

	if err := stub.PutState("user:"+user.Uid, []byte(strconv.FormatInt(initBalance, 10))); err != nil {
		return shim.Error("Add balance error: " + err.Error())
	}

	if err := trace(stub, txid, user.Uid, user.Uid, user.Uid, "userInit", "", txtime, initBalance); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// 积分充值
func balanceDeposit(stub shim.ChaincodeStubInterface, admin, args []string) peer.Response {

	logger.Info("balanceDeposit")

	if resp := paramCheck(args, 2); resp != "" {
		return shim.Error(resp)
	}

	user, txtime, txid, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	value, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil || value <= 0 {
		return shim.Error("Balance value is illegal")
	}

	if !contain(user.Uid, admin) { // 只有管理员才可以给用户充值
		return shim.Error("User is not administrator")
	}

	if balance, err := balanceAdd(stub, args[0], value); err == nil && balance > 0 {
		if err := trace(stub, txid, user.Uid, args[0], args[0], "balanceDeposit", "", txtime, value); err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success([]byte(strconv.FormatInt(balance, 10)))
	}

	return shim.Error("Balance deposit failure")
}

// 积分转账
func balanceTransfer(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("balanceTransfer")

	if resp := paramCheck(args, 2); resp != "" {
		return shim.Error(resp)
	}

	user, txtime, txid, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	value, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil || value < 0 {
		return shim.Error("Balance value is illegal")
	}

	if value <= 0 {
		return shim.Error("Balance value must more than 0")
	}

	if value > user.Balance {
		return shim.Error("Balance not enough")
	}

	if _, err := balanceAdd(stub, user.Uid, 0-value); err != nil {
		return shim.Error("Balance transform failed")
	}

	if _, err := balanceAdd(stub, args[0], value); err != nil {
		return shim.Error("Balance transform failed")
	}

	if err := trace(stub, txid, user.Uid, args[0], args[0], "balanceTransfer", "", txtime, value); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// 数字版权注册
func assetRegister(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("assetRegister")

	if resp := paramCheck(args, 8); resp != "" {
		return shim.Error(resp)
	}

	assetBytes, err := stub.GetState("asset:" + args[0])
	if err == nil && assetBytes != nil && len(assetBytes) > 0 {
		return shim.Error("Asset already exist")
	}

	// nextid := IDCounter(stub, "asset")
	user, txtime, txid, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	sprice, err := strconv.ParseInt(args[5], 10, 64)
	if err != nil {
		return shim.Error("Illegal argument")
	}

	tprice, err := strconv.ParseInt(args[6], 10, 64)
	if err != nil {
		return shim.Error("Illegal argument")
	}

	etime, err := strconv.ParseInt(args[7], 10, 64)
	if err != nil {
		return shim.Error("Illegal argument")
	}

	if etime <= txtime {
		return shim.Error("Expire time must be larger than current time")
	}

	asset := &Asset{
		Dci: args[0],
		// Name:  args[1],
		Name:     args[1],
		Category: args[2],
		Describe: args[3],
		Content:  args[4],
		Author:   user.Uid,
		Owner:    user.Uid,
		Sprice:   sprice,
		Tprice:   tprice,
		Ctime:    txtime,
		Ttime:    txtime,
		Etime:    etime,
		Status:   int64(0),
	}

	if err := assetPut(stub, asset); err != nil {
		return shim.Error("Save asset error: " + err.Error())
	}

	// 用户和版权建立关联
	if ok := permPut(stub, user.Uid, asset.Dci, "AuthorAndOwner", "", 0, asset.Ctime, -1, asset.Ctime, 0); !ok {
		return shim.Error("Associate user and asset error: " + err.Error())
	}

	if err := trace(stub, txid, asset.Author, asset.Owner, asset.Dci, "assetRegister", "", asset.Ctime, 0); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// 设置版权字段
func assetSet(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	logger.Info("assetSet")

	if resp := paramCheck(args, 2); resp != "" {
		return shim.Error(resp)
	}

	user, txtime, txid, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	asset, err := assetGet(stub, args[0])
	if err != nil {
		return shim.Error("Get asset error: " + err.Error())
	}

	if user.Uid != asset.Owner {
		return shim.Error("Only the owner could modify asset value")
	}

	fields := make(map[string]string)

	pp := reflect.ValueOf(asset)
	//pt := reflect.TypeOf(user)
	pt := pp.Type()

	if pt.Kind() == reflect.Ptr {
		pt = pt.Elem()
	}
	if pt.Kind() == reflect.Struct {
		for i := 0; i < pt.NumField(); i++ {
			fname := pt.Field(i).Name
			fields[strings.ToLower(fname)] = fname
		}
	}

	if len(fields) < 1 {
		return shim.Error("No field could be modified")
	}

	kvs := strings.Split(args[1], ",")

	for _, kv := range kvs {
		kv := strings.Split(kv, "=")
		if len(kv) < 2 {
			return shim.Error("Illegal argument")
		}
		k := strings.TrimSpace(kv[0])
		v := strings.TrimSpace(kv[1])
		k = strings.ToLower(k)
		if k != "describe" && k != "sprice" && k != "tprice" && k != "etime" && k != "status" {
			return shim.Error("The field could not be modified")
		}

		if k == "status" {
			if v != "0" && v != "1" {
				return shim.Error("Status value must be 0 or 1")
			}
		}

		// TODO 是否允许修改版权有效期, 版权有效期的修改是否需要作出规定

		field := pp.Elem().FieldByName(fields[k])
		if field.CanSet() {
			switch field.Kind() {
			case reflect.String:
				field.SetString(v)
			case reflect.Int64, reflect.Int32, reflect.Int, reflect.Int16, reflect.Int8:
				iv, err := strconv.ParseInt(v, 10, 64)
				if err != nil || iv < 0 {
					return shim.Error("Illegal argument")
				}

				field.SetInt(iv)
			}

		} else {
			return shim.Error("The field could not be modified")
		}
	}

	if err := trace(stub, txid, user.Uid, asset.Owner, asset.Dci, "assetSet", "", txtime, 0); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// 版权转让, 修改当前所有者的版权信息
func assetTransfer(stub shim.ChaincodeStubInterface, asset *Asset, uid string, txtime int64) error {
	owner := asset.Owner
	asset.Owner = uid
	asset.Ttime = txtime
	assetBytes, err := json.Marshal(asset)
	if err != nil || assetBytes == nil || len(assetBytes) == 0 {
		return fmt.Errorf("Marshal asset failure")
	}

	if err := stub.PutState("asset:"+asset.Dci, assetBytes); err != nil {
		return fmt.Errorf("Save asset failure")
	}

	sperm, err := permGet(stub, asset.Dci, owner) // 用户自身与版权的关系也需要进行修改
	if err != nil || sperm == nil {
		return fmt.Errorf("Perm not found")
	}

	if sperm.Ptype == "AuthorAndOwner" {
		sperm.Ptype = "Author"
	} else {
		sperm.Ptype = "Invalid"
		// 这里也可以选择删除用户和版权的关联kv, 保留此kv可以表示用户和版权之间曾经有过关联
	}
	sperm.Status = 4 // 版权转让后用户对版权的权限将失效

	spermBytes, err := json.Marshal(sperm)
	if err != nil || spermBytes == nil || len(spermBytes) == 0 {
		return fmt.Errorf("Marshal perm failure")
	}

	if err := stub.PutState("perm:"+sperm.Dci+":"+sperm.Uid, spermBytes); err != nil {
		return fmt.Errorf("Save perm failure")
	}

	return nil
}

// 发起版权权限申请
func permRequest(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("permRequest")

	if resp := paramCheck(args, 6); resp != "" {
		return shim.Error(resp)
	}

	user, txtime, txid, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	dci := args[0]

	ptype := args[1]

	clause := args[2]

	pricestr := args[3]

	vtimestr := args[4]

	etimestr := args[5]

	asset, err := assetGet(stub, dci)
	if err != nil || asset == nil {
		return shim.Error("Asset not found")
	}

	ptype, price, vtime, etime, err := permCheck(stub, user, asset, ptype, pricestr, vtimestr, etimestr, txtime)
	if err != nil {
		return shim.Error(err.Error())
	}
	/* 允许还价
	if ptype == "Owner" && user.Balance < asset.Tprice { // 用户余额必须足够支付授权
		return shim.Error("Balance not enough")
	}
	if ptype == "User" && user.Balance < asset.Sprice {
		return shim.Error("Balance not enough")
	}
	*/

	// 发起权限申请
	if ok := permPut(stub, user.Uid, dci, ptype, clause, price, vtime, etime, txtime, 1); !ok {
		return shim.Error("Save perm failure")
	}
	// TODO 撤销权限申请

	/*
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
	*/

	if err := trace(stub, txid, user.Uid, asset.Owner, asset.Dci, "permRequest", "", txtime, 0); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// 修改版权权限字段, 仅可修改申请中和已撤销的权限
func permSet(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	logger.Info("permSet")

	if resp := paramCheck(args, 2); resp != "" {
		return shim.Error(resp)
	}

	user, txtime, txid, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	perm, err := permGet(stub, args[0], user.Uid)
	if err != nil {
		return shim.Error("Get perm error: " + err.Error())
	}

	if user.Uid != perm.Uid {
		return shim.Error("Only the requisitioner could modify perm information")
	}

	if perm.Status != 1 && perm.Status != 2 { // 只能修改申请中和已撤销的权限请求 (已撤销表示申请者放弃申请)
		return shim.Error("The modify is not allowed")
	}

	fields := make(map[string]string)

	pp := reflect.ValueOf(perm)
	//pt := reflect.TypeOf(user)
	pt := pp.Type()

	if pt.Kind() == reflect.Ptr {
		pt = pt.Elem()
	}
	if pt.Kind() == reflect.Struct {
		for i := 0; i < pt.NumField(); i++ {
			fname := pt.Field(i).Name
			fields[strings.ToLower(fname)] = fname
		}
	}

	if len(fields) < 1 {
		return shim.Error("No field could be modified")
	}

	kvs := strings.Split(args[1], ",")

	for _, kv := range kvs {
		kv := strings.Split(kv, "=")
		if len(kv) < 2 {
			return shim.Error("Illegal argument")
		}
		k := strings.TrimSpace(kv[0])
		v := strings.TrimSpace(kv[1])
		k = strings.ToLower(k)
		if k != "ptype" && k != "clause" && k != "price" && k != "vtime" && k != "etime" && k != "status" {
			return shim.Error("The field could not be modified")
		}

		if k == "status" {
			if v != "1" && v != "2" {
				return shim.Error("Modified status value must be 1 or 2")
			}
		}

		// TODO 是否允许修改版权有效期, 版权有效期的修改是否需要作出规定

		field := pp.Elem().FieldByName(fields[k])
		if field.CanSet() {
			switch field.Kind() {
			case reflect.String:
				field.SetString(v)
			case reflect.Int64, reflect.Int32, reflect.Int, reflect.Int16, reflect.Int8:
				iv, err := strconv.ParseInt(v, 10, 64)
				if err != nil || iv < 0 {
					return shim.Error("Illegal argument")
				}

				field.SetInt(iv)
			}

		} else {
			return shim.Error("The field could not be modified")
		}
	}

	asset, err := assetGet(stub, perm.Dci)
	if err != nil {
		return shim.Error("Asset not found")
	}

	if perm.Vtime < txtime {
		perm.Vtime = txtime
	}

	if perm.Etime <= 0 {
		perm.Etime = asset.Etime
	}

	if user.Balance < perm.Price {
		return shim.Error("Balance not enough")
	}

	if err := permTimeCheck(asset, txtime, perm.Vtime, perm.Etime); err != nil {
		return shim.Error(err.Error())
	}

	if err := trace(stub, txid, user.Uid, perm.Uid, perm.Dci, "permSet", "", txtime, 0); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// 版权交易(版权所有者主动进行无偿版权授权和转让)
func assetTrade(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	logger.Info("assetPermRequest")

	if resp := paramCheck(args, 6); resp != "" {
		return shim.Error(resp)
	}

	user, txtime, txid, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	dci := args[0]

	uid := args[1]

	ptype := args[2]

	clause := args[3]

	asset, err := assetGet(stub, dci)
	if err != nil || asset == nil {
		return shim.Error("Asset not found")
	}

	pricestr := "0"

	vtimestr := args[4]

	etimestr := args[5]

	ptype, price, vtime, etime, err := permCheck(stub, user, asset, ptype, pricestr, vtimestr, etimestr, txtime)
	if err != nil {
		return shim.Error(err.Error())
	}

	if ok := permPut(stub, user.Uid, dci, ptype, clause, price, vtime, etime, txtime, 0); !ok {
		return shim.Error("Save perm failure")
	}

	action := "assetTrade"

	if ptype == "Owner" || ptype == "AuthorAndOwner" { // 版权转让
		action = "assetTransfer"

		if err := assetTransfer(stub, asset, uid, txtime); err != nil {
			return shim.Error(err.Error())
		}
	} else {
		action = "assetGrant"
	}

	if err := trace(stub, txid, user.Uid, uid, dci, action, "", txtime, 0); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// 回复版权权限申请
func permResponse(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	logger.Info("permResponse")

	if resp := paramCheck(args, 3); resp != "" {
		return shim.Error(resp)
	}

	action := "permResponse"

	dci := args[0]

	uid := args[1]

	status, err := strconv.ParseInt(args[2], 10, 64)
	if err != nil {
		return shim.Error("Illegal argument")
	}

	user, txtime, txid, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	asset, err := assetGet(stub, dci)
	if err != nil || asset == nil {
		return shim.Error("Asset not found")
	}

	if asset.Owner != user.Uid {
		return shim.Error("Only Owner of the asset allowed")
	}

	if asset.Status != 0 {
		return shim.Error("Asset is invalid")
	}

	if asset.Etime <= txtime {
		return shim.Error("Asset cryptoright is expired")
	}

	perm, err := permGet(stub, dci, uid)
	if err != nil || perm == nil {
		return shim.Error("Perm not found")
	}

	if perm.Status != 1 {
		return shim.Error("Perm not in requesting status")
	}

	if perm.Price < 0 {
		return shim.Error("Perm price should be >=0")
	}

	if perm.Vtime < txtime {
		perm.Vtime = txtime
	}

	if err := permTimeCheck(asset, txtime, perm.Vtime, perm.Etime); err != nil {
		return shim.Error(err.Error())
	}

	perm.Status = status
	permBytes, err := json.Marshal(perm)
	if err != nil || permBytes == nil || len(permBytes) == 0 {
		return shim.Error("Marshal perm failure")
	}

	if err := stub.PutState("perm:"+perm.Dci+":"+perm.Uid, permBytes); err != nil {
		return shim.Error("Save perm failure")
	}

	if perm.Status == 0 { // 申请已批准
		if perm.Ptype == "Owner" || perm.Ptype == "AuthorAndOwner" { // 版权转让
			action = "assetTransfer"

			if err := assetTransfer(stub, asset, uid, txtime); err != nil {
				return shim.Error(err.Error())
			}

		} else {
			action = "assetGrant"
		}

		if _, err := balanceAdd(stub, uid, 0-perm.Price); err != nil {
			return shim.Error("Balance transfer failure")
		}

		if _, err := balanceAdd(stub, user.Uid, perm.Price); err != nil {
			return shim.Error("Balance transfer failure")
		}
	}

	if err := trace(stub, txid, user.Uid, perm.Uid, perm.Dci, action, "", txtime, perm.Price); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// 查询用户信息
func userQuery(stub shim.ChaincodeStubInterface) peer.Response {

	logger.Info("userQuery")

	user, txtime, txid, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}
	var buffer bytes.Buffer

	buffer.WriteString("{")

	buffer.WriteString("\"uid\":")
	buffer.WriteString("\"")
	buffer.WriteString(user.Uid)
	buffer.WriteString("\"")
	buffer.WriteString(",")

	buffer.WriteString("\"mspid\":")
	buffer.WriteString("\"")
	buffer.WriteString(user.Mspid)
	buffer.WriteString("\"")
	buffer.WriteString(",")

	buffer.WriteString("\"role\":")
	buffer.WriteString("\"")
	buffer.WriteString(user.Role)
	buffer.WriteString("\"")
	buffer.WriteString(",")

	buffer.WriteString("\"balance\":")
	buffer.WriteString("\"")
	buffer.WriteString(strconv.FormatInt(user.Balance, 10))
	buffer.WriteString("\"")
    buffer.WriteString(",")

	buffer.WriteString("\"txid\":")
	buffer.WriteString("\"")
	buffer.WriteString(txid)
	buffer.WriteString("\"")
	buffer.WriteString(",")

	buffer.WriteString("\"time\":")
	buffer.WriteString("\"")
	buffer.WriteString(strconv.FormatInt(txtime, 10))
	buffer.WriteString("\"")

	buffer.WriteString("}")

	return shim.Success(buffer.Bytes())
}

// 查询用户积分
func balanceQuery(stub shim.ChaincodeStubInterface) peer.Response {

	logger.Info("balanceQuery")

	user, _, _, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	return shim.Success([]byte(strconv.FormatInt(user.Balance, 10)))
}

// 通过版权标识符查询多个版权的详细信息
func assetQuery(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("assetQuery")

	var buffer bytes.Buffer

	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for _, dci := range args {

		assetBytes, err := stub.GetState("asset:" + dci)
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

// 查询和用户有关的所有版权(包括作者关系, 所有权关系, 授权关系)
func userAssetQuery(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("userAssetQuery")

	var psizestr string
	var bookmark string
	var buffer *bytes.Buffer

	length := len(args)

	if length > 2 {
		return shim.Error("Incorrect number of arguments. Expecting <3")
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
		resultIterator, err := stub.GetStateByPartialCompositeKey("user~perm", []string{user.Uid})
		if err != nil {
			return shim.Error("Query user perm error: " + err.Error())
		}
		defer resultIterator.Close()

		buffer = constructQueryResponseFromIterator(stub, resultIterator)
	} else {
		pagesize, err := strconv.ParseInt(psizestr, 10, 32)
		if err != nil {
			return shim.Error("Query user asset error: " + err.Error())
		}
		resultsIterator, responseMetadata, err := stub.GetStateByPartialCompositeKeyWithPagination("user~perm", []string{user.Uid}, int32(pagesize), bookmark)
		if err != nil {
			return shim.Error("Query user asset error: " + err.Error())
		}
		defer resultsIterator.Close()

		buffer = constructQueryResponseFromIterator(stub, resultsIterator)

		addPaginationMetadataToQueryResults(buffer, responseMetadata)
	}

	return shim.Success(buffer.Bytes())
}

// 查询版权关联的所有权限
func permQuery(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("permQuery")

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
			bookmark = args[2] // bookmark可空
		}
	}

	/*
		user, _, _, err := stubParse(stub)
		if err != nil {
			return shim.Error("Parse stub error: " + err.Error())
		}
	*/

	if psizestr == "" {
		resultIterator, err := stub.GetStateByRange("perm:"+args[0]+":", "perm:"+args[0]+":~")
		if err != nil {
			return shim.Error("Query asset error: " + err.Error())
		}
		defer resultIterator.Close()

		buffer = constructQueryResponseFromIterator(stub, resultIterator)
	} else {
		pagesize, err := strconv.ParseInt(psizestr, 10, 32)
		if err != nil {
			return shim.Error("Query asset error: " + err.Error())
		}
		resultsIterator, responseMetadata, err := stub.GetStateByRangeWithPagination("record:"+args[0]+":", "record:"+args[1]+":~", int32(pagesize), bookmark)
		if err != nil {
			return shim.Error("Query asset error: " + err.Error())
		}
		defer resultsIterator.Close()

		buffer = constructQueryResponseFromIterator(stub, resultsIterator)

		addPaginationMetadataToQueryResults(buffer, responseMetadata)
	}

	return shim.Success(buffer.Bytes())
}

// 查询用户和指定版权的权限关联
func userPermQuery(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("userPermQuery")

	var dci string
	var uid string

	if len(args) != 1 && len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 1 or 2")
	}

	user, _, _, err := stubParse(stub)
	if err != nil {
		return shim.Error("Parse stub error: " + err.Error())
	}

	dci = args[0]

	if len(args) == 1 {
		uid = user.Uid
	} else {
		// TODO 是否需要限定对非自身数据的查询只能由管理员执行
		/*
			if !contain(user.Uid, admin) { // 只有管理员才可以查询其他用户的权限
				return shim.Error("User is not administrator")
			}
		*/
		uid = args[1]
	}

	asset, err := assetGet(stub, args[0])
	if err != nil || asset == nil {
		return shim.Error("Asset not found")
	}

	permBytes, err := stub.GetState("perm:" + dci + ":" + uid)
	if err != nil || permBytes == nil || len(permBytes) == 0 {
		return shim.Success(nil)
	}

	return shim.Success(permBytes)
}

// 日志查询(支持单查询与多查询)
func recordQuery(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("recordQuery")

	var psizestr string
	var bookmark string
	var buffer *bytes.Buffer

	length := len(args)

	if length < 1 || length > 4 {
		return shim.Error("Incorrect number of arguments. Expecting 1~4")
	}
	if length > 2 {
		psizestr = args[2]
		if length > 3 {
			bookmark = args[3] // bookmark可空
		}
	}

	if args[0] == "" {
		return shim.Error("The first argument can not be empty string")
	}

	if length == 1 || args[1] == "" {
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
			return shim.Error("Query record error: " + err.Error())
		}
		resultsIterator, responseMetadata, err := stub.GetStateByRangeWithPagination("record:"+args[0], "record:"+args[1], int32(pagesize), bookmark)
		if err != nil {
			return shim.Error("Query record error: " + err.Error())
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
		resultIterator, err := stub.GetStateByPartialCompositeKey("user~record", []string{user.Uid})
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
		resultsIterator, responseMetadata, err := stub.GetStateByPartialCompositeKeyWithPagination("user~record", []string{user.Uid}, int32(pagesize), bookmark)
		if err != nil {
			return shim.Error("Query user record error: " + err.Error())
		}
		defer resultsIterator.Close()

		buffer = constructQueryResponseFromIterator(stub, resultsIterator)

		addPaginationMetadataToQueryResults(buffer, responseMetadata)
	}

	return shim.Success(buffer.Bytes())
}

// 键的历史记录
func historyQuery(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("historyQuery")

	var key string
	var buffer bytes.Buffer

	length := len(args)

	if length < 1 {
		return shim.Error("Incorrect number of arguments. Expecting >0")
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

// 全查询, 支持查询指定类型的简单键的所有键值
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
func richQuery(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	logger.Info("richQuery")

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
		return shim.Error("The first argument must be a selector statement")
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

func (dc *DC) Init(stub shim.ChaincodeStubInterface) peer.Response {

	_, args := stub.GetFunctionAndParameters()
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting >0")
	}

	for _, uid := range args {
		dc.Admin = append(dc.Admin, uid) // 链码初始化时指定管理员(参数为管理员ski)
	}

	return shim.Success(nil)
}

func (dc *DC) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()

	switch function {
	case "userInit":
		return userInit(stub)
	case "balanceDeposit":
		return balanceDeposit(stub, dc.Admin, args)
	case "balanceTransfer":
		return balanceTransfer(stub, args)
	case "assetRegister":
		return assetRegister(stub, args)
	case "assetSet":
		return assetSet(stub, args)
	case "assetTrade":
		return assetTrade(stub, args)
	case "permRequest":
		return permRequest(stub, args)
	case "permSet":
		return permSet(stub, args)
	case "permResponse":
		return permResponse(stub, args)
	case "balanceQuery":
		return balanceQuery(stub)
	case "userQuery":
		return userQuery(stub)
	case "assetQuery":
		return assetQuery(stub, args)
	case "userAssetQuery":
		return userAssetQuery(stub, args)
	case "permQuery":
		return permQuery(stub, args)
	case "userPermQuery":
		return userPermQuery(stub, args)
	case "recordQuery":
		return recordQuery(stub, args)
	case "userRecordQuery":
		return userRecordQuery(stub, args)
	case "historyQuery":
		return historyQuery(stub, args)
	case "fullQuery":
		return fullQuery(stub, args)
	case "richQuery":
		return richQuery(stub, args)
	default:
		return shim.Error("Unsupported function: " + function)
	}

	// return shim.Success(nil)
}

var logger = shim.NewLogger("ebclog")

func main() {

	logger.SetLevel(shim.LogDebug)

	err := shim.Start(new(DC))
	if err != nil {
		logger.Error("Error starting AssetManage chaincode: %s", err)
	}
}
