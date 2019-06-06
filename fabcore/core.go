package fabcore

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/cryptosuite"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/comm"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/discovery"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric/common/util"

	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	contextImpl "github.com/hyperledger/fabric-sdk-go/pkg/context"
	fabImpl "github.com/hyperledger/fabric-sdk-go/pkg/fab"
	mspImpl "github.com/hyperledger/fabric-sdk-go/pkg/msp"
	cb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	putils "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/utils"
)

type FabClient struct {
	userName      string
	orgName       string
	channels      map[string]*ChanClient
	connNumber    int
	fabricSDK     *fabsdk.FabricSDK
	resMgmtClient *resmgmt.Client
	discClient    *discovery.Client
	mspClient     *mspclient.Client
}

type ChanClient struct {
	channelID     string
	channelClient *channel.Client
	eventClient   *event.Client
	ledgerClient  *ledger.Client
	transactor    *fab.Transactor
}

/*
func regSplit(text string, delimeter string) []string {
	reg := regexp.MustCompile(delimeter)
	indexes := reg.FindAllStringIndex(text, -1)
	laststart := 0
	result := make([]string, len(indexes)+1)
	for i, element := range indexes {
		result[i] = text[laststart:element[0]]
		laststart = element[1]
	}
	result[len(indexes)] = text[laststart:len(text)]
	return result
}
*/

/**
 * 拷贝文件夹,同时拷贝文件夹中的文件
 * @param srcPath  		需要拷贝的文件夹路径: D:/test
 * @param destPath		拷贝到的位置: D:/backup/
 */
func copyDir(srcPath string, destPath string) error {
	//检测目录正确性
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return err
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("srcPath is not a directory")
	}

	destInfo, err := os.Stat(destPath)
	if err != nil {
		return err
	}

	if !destInfo.IsDir() {
		return fmt.Errorf("destInfo is not a directory")

	}

	// 加上拷贝时间:不用可以去掉
	// destPath = destPath + "_" + time.Now().Format("20060102150405")

	return filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if !f.IsDir() {
			path := strings.Replace(path, "\\", "/", -1)
			destNewPath := strings.Replace(path, srcPath, destPath, -1)
			// fmt.Println("copy file " + path + " to " + destNewPath)
			_, err := copyFile(path, destNewPath)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

//生成目录并拷贝文件
func copyFile(src, dest string) (w int64, err error) {
	w = 0
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()
	//分割path目录
	destSplitPathDirs := strings.Split(dest, "/")

	//检测时候存在目录
	destSplitPath := ""
	for index, dir := range destSplitPathDirs {
		if index < len(destSplitPathDirs)-1 {
			destSplitPath = destSplitPath + dir + "/"
			b, _ := pathExists(destSplitPath)
			if b == false {
				//fmt.Println("创建目录:" + destSplitPath)
				//创建目录
				err = os.Mkdir(destSplitPath, os.ModePerm)
				if err != nil {
					return
				}
			}
		}
	}
	dstFile, err := os.Create(dest)
	if err != nil {
		return
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

//检测文件夹路径是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 使用配置文件实例化fabric sdk
func NewFabSdk(configFile string) (*fabsdk.FabricSDK, error) {
	sdk, err := fabsdk.New(config.FromFile(configFile))
	if err != nil {
		return nil, fmt.Errorf("Instantiate fabric sdk error: %v", err)
	}

	return sdk, nil
}

// 创建fabric客户端
func NewFabClient(sdk *fabsdk.FabricSDK, orgName, userName string) (*FabClient, error) {
	if userName == "" || orgName == "" {
		return nil, fmt.Errorf("userName or orgName is empty")
	}

	client := &FabClient{
		userName: userName,
		orgName:  orgName,
		channels: make(map[string]*ChanClient),
	}

	if sdk == nil {
		sdk, err := NewFabSdk(fmt.Sprintf("./config/%s-network.yaml", strings.ToLower(orgName)))
		if err != nil {
			return nil, err
		}
		client.fabricSDK = sdk
	} else {
		client.fabricSDK = sdk
	}

	// fabsdk.New(integration.ConfigBackend)

	resClientContext := client.fabricSDK.Context(fabsdk.WithUser(userName), fabsdk.WithOrg(orgName))
	if resClientContext != nil {
		resMgmtClient, err := resmgmt.New(resClientContext)
		if err == nil {
			client.resMgmtClient = resMgmtClient
		}

		clientContext, err := resClientContext()
		if err == nil {
			discClient, err := discovery.New(clientContext)
			if err == nil {
				client.discClient = discClient
			}
		}
	}

	mspClient, err := mspclient.New(client.fabricSDK.Context(), mspclient.WithOrg(orgName))
	if err == nil {
		client.mspClient = mspClient
	}

	// fmt.Println("Instantiate fabric sdk successfully")
	return client, nil
}

// 创建channel客户端
func (fc *FabClient) NewChanClient(channelID string) {
	if fc.channels[channelID] == nil || fc.channels[channelID].channelClient == nil {
		fc.channels[channelID] = new(ChanClient)
		fc.channels[channelID].channelID = channelID
		resClientContext := fc.fabricSDK.Context(fabsdk.WithUser(fc.userName), fabsdk.WithOrg(fc.orgName))
		channelClientContext := fc.fabricSDK.ChannelContext(channelID, fabsdk.WithUser(fc.userName), fabsdk.WithOrg(fc.orgName))
		channelContext, err := channelClientContext()

		if err == nil {
			//contextComm.Client
			resctx, err := resClientContext()
			if err == nil {
				reqCtx, cancel := contextImpl.NewRequest(resctx, contextImpl.WithTimeoutType(fab.ResMgmt))
				defer cancel()
				transactor, err := channelContext.ChannelService().Transactor(reqCtx)
				if err == nil {
					fc.channels[channelID].transactor = &transactor
				}
			}

		}

		channelClient, err := channel.New(channelClientContext)
		if err == nil {
			fc.channels[channelID].channelClient = channelClient
		}

		eventClient, err := event.New(channelClientContext, event.WithBlockEvents())
		if err == nil {
			fc.channels[channelID].eventClient = eventClient
		}

		ledgerClient, err := ledger.New(channelClientContext)
		if err == nil {
			fc.channels[channelID].ledgerClient = ledgerClient
		}
	}
}

// 客户端检查
func (fc *FabClient) clientCheck(clientType, channelID string) error {
	if fc.fabricSDK == nil {
		return fmt.Errorf("FabricSDK is nil")
	}
	if clientType == "" {
		return nil
	}
	if clientType == "resMgmtClient" {
		if fc.resMgmtClient == nil {
			return fmt.Errorf("ResMgmtClient is nil")
		}
	} else if clientType == "discClient" {
		if fc.discClient == nil {
			return fmt.Errorf("discClient is nil")
		}
	} else if clientType == "mspClient" {
		if fc.mspClient == nil {
			return fmt.Errorf("mspClient is nil")
		}
	} else if clientType == "ledgerClient" || clientType == "channelClient" || clientType == "eventClient" || clientType == "transactor" {
		if fc.channels[channelID] == nil || fc.channels[channelID].ledgerClient == nil || fc.channels[channelID].channelClient == nil || fc.channels[channelID].eventClient == nil || fc.channels[channelID].transactor == nil {
			fc.NewChanClient(channelID)
		}

		if fc.channels[channelID] == nil || fc.channels[channelID].ledgerClient == nil || fc.channels[channelID].channelClient == nil || fc.channels[channelID].eventClient == nil || fc.channels[channelID].transactor == nil {
			return fmt.Errorf("Channel Client is nil")
		}
	}
	return nil
}

// 从配置文件获取配置信息
func (fc *FabClient) getConfig() (core.CryptoSuiteConfig, msp.IdentityConfig, fab.EndpointConfig) {
	backend, err := fc.fabricSDK.Config()
	if err != nil {
		return nil, nil, nil
	}

	// set cryptoSuiteConfig
	cryptoSuiteConfig := cryptosuite.ConfigFromBackend(backend)

	// set identityConfig
	identityConfig, err := mspImpl.ConfigFromBackend(backend)
	if err != nil {
		identityConfig = nil
	}

	// set endpointConfig
	endpointConfig, err := fabImpl.ConfigFromBackend(backend)
	if err != nil {
		endpointConfig = nil
	}

	return cryptoSuiteConfig, identityConfig, endpointConfig
}

func getBlock(blk *cb.Block) (*Block, error) {
	header := blk.Header
	data := blk.Data
	block := &Block{
		Number:       header.Number,
		Size:         blk.XXX_Size(),
		PreviousHash: fmt.Sprintf("%x", header.PreviousHash),
		DataHash:     fmt.Sprintf("%x", header.DataHash),
	}

	asn1Head := Asn1Head{
		Number:       int64(header.Number), // asn1不支持uint64类型
		PreviousHash: header.PreviousHash,
		DataHash:     header.DataHash,
	}

	asn1HeadBytes, err := asn1.Marshal(asn1Head)
	fmt.Println(err)
	hashBytes := util.ComputeSHA256(asn1HeadBytes)
	hash := fmt.Sprintf("%x", hashBytes)

	block.Hash = hash

	// common.Transaction
	// protostr := proto.CompactTextString(blk)
	// fmt.Println(proto.MarshalTextString(data))

	for _, txBytes := range data.Data {
		/* 笨办法
		certStart := bytes.IndexAny(txBytes, "-----BEGIN")
		fmt.Println(string(txBytes[:certStart]))
		res := regexp.MustCompile(`\w+`).FindAllString(string(txBytes[:certStart]), -1)
		fmt.Println(res)
		*/

		env, err := putils.GetEnvelopeFromBlock(txBytes)
		if err != nil || env == nil {
			continue
		}

		tx, err := getTransaction(env, block.Number)
		if err != nil || tx == nil {
			continue
		}

		block.Timestamp = tx.Timestamp
		block.Txs = append(block.Txs, tx)

		/*
			d := regSplit(string(txBytes[:certStart]), `\s`)
			fmt.Println(len(d))
			fmt.Println(d)
			dd := regexp.MustCompile(`\w+`).FindAllString(string(txBytes[:certStart]), -1)
			fmt.Println(dd)
		*/
	}

	return block, nil
}

func getTransaction(env *cb.Envelope, blkNum uint64) (*Transaction, error) {

	// s, _ := putils.GetSignatureHeader(env.Signature)
	// fmt.Println(string(s.Creator))
	// sd, _ := env.AsSignedData()

	payload, err := putils.GetPayload(env)
	if err != nil {
		return nil, fmt.Errorf("Extract payload from envelope error: %s", err.Error())
	}

	channelHeaderBytes := payload.Header.ChannelHeader

	channelHeader := &cb.ChannelHeader{}
	if err := proto.Unmarshal(channelHeaderBytes, channelHeader); err != nil {
		return nil, fmt.Errorf("Extract channelHeader from payload error: %s", err.Error())
	}

	tx, err := putils.GetTransaction(payload.Data)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal transaction payload error: %s", err.Error())
	}

	var txType string
	var mspid string

	//log.Println(tx.Actions[0].String())

	headBytes := tx.Actions[0].Header
	if headBytes == nil || len(headBytes) == 0 {
		txType = "config"
		mspid = "OrdererMSP"
		/*
			result := &Transaction{
				TxID: channelHeader.TxId,
				//TxType:        strconv.Itoa(int(channelHeader.GetType())),
				ChannelID: channelHeader.ChannelId,
				//Timestamp:     time.Unix(channelHeader.Timestamp.Seconds, 0).Format("2006-01-02 15:04:05"),
				Timestamp: fmt.Sprintf("%d", time.Unix(channelHeader.Timestamp.Seconds, 0).Unix()),
			}
			return result, nil
		*/
	} else {
		txType = "invoke"
		certStart := bytes.IndexAny(headBytes, "-----BEGIN")
		mspids := regexp.MustCompile(`\w+`).FindAllString(string(headBytes[:certStart]), -1)
		mspid = mspids[len(mspids)-1]

		//certStart := bytes.IndexAny(headBytes, "-----BEGIN")
		//mspid = regexp.MustCompile(`\w+`).FindString(string(headBytes[:certStart]))
	}

	chaincodeActionPayload, err := putils.GetChaincodeActionPayload(tx.Actions[0].Payload)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal chaincode action payload error: %s", err.Error())
	}

	var endorsers []string
	action := chaincodeActionPayload.Action

	for _, endorsement := range action.Endorsements {
		var itemBytes []byte
		if txType == "config" {
			itemBytes = endorsement.Signature
		} else {
			itemBytes = endorsement.Endorser
		}
		sign, err := putils.GetSignatureHeader(itemBytes)
		if err != nil {
			continue
		}
		//certStart := bytes.IndexAny(itemBytes, "-----BEGIN")
		//mspids := regexp.MustCompile(`\w+`).FindAllString(string(itemBytes[:certStart]), -1)
		endorser := string(sign.Creator)
		endorser = strings.Replace(endorser, "OrdererOrg", "OrdererMSP", -1)
		endorsers = append(endorsers, endorser)
	}

	propPayload := &pb.ChaincodeProposalPayload{}
	if err := proto.Unmarshal(chaincodeActionPayload.ChaincodeProposalPayload, propPayload); err != nil {
		return nil, fmt.Errorf("Extract channelHeader from payload error: %s", err.Error())
	}

	invokeSpec := &pb.ChaincodeInvocationSpec{}
	err = proto.Unmarshal(propPayload.Input, invokeSpec)
	if err != nil {
		return nil, fmt.Errorf("Extract channelHeader from payload error: %s", err.Error())
	}

	var args []string
	var chaincodeID string
	var chaincodeLang string
	if invokeSpec.ChaincodeSpec != nil {
		chaincodeID = invokeSpec.ChaincodeSpec.ChaincodeId.Name
		chaincodeLang = invokeSpec.ChaincodeSpec.Type.String()
		txType = string(invokeSpec.ChaincodeSpec.Input.Args[0])
		for _, v := range invokeSpec.ChaincodeSpec.Input.Args {
			args = append(args, string(v))
		}
	}

	/*
		TxID        string `json:"txid"`
		TxType      string `json:"txtype"`
		ChannelID   string `json:"channelid"`
		ChaincodeID string `json:"chaincodeid"`
		MspID       string `json:"mspid"`
		Timestamp   string `json:"timestamp`
	*/

	result := &Transaction{
		TxID:      channelHeader.TxId,
		BlkNum:    blkNum,
		TxType:    txType,
		Size:      env.XXX_Size(),
		MspID:     mspid,
		Endorsers: endorsers,
		//TxType:        strconv.Itoa(int(channelHeader.GetType())),
		Args:          args,
		ChannelID:     channelHeader.ChannelId,
		ChaincodeID:   chaincodeID,
		ChaincodeLang: chaincodeLang,
		//Timestamp:     time.Unix(channelHeader.Timestamp.Seconds, 0).Format("2006-01-02 15:04:05"),
		Timestamp: fmt.Sprintf("%d", time.Unix(channelHeader.Timestamp.Seconds, 0).Unix()),
	}

	return result, nil
}

// TODO
func convertPeerName(sourceurl string, endpointConfig fab.EndpointConfig) string {
	var fix string
	su := strings.Split(sourceurl, ":")
	if len(su) > 1 {
		fix = su[1]
	} else {
		fix = su[0]
	}
	for _, peer := range endpointConfig.ChannelPeers("mychannel") {
		if len(su) > 1 {
			if strings.HasSuffix(peer.URL, fix) {
				return peer.MSPID
			}
		} else {
			if strings.HasPrefix(peer.URL, fix) {
				return peer.MSPID
			}
		}
	}
	return sourceurl
}

func parseBlockEvent(ev *fab.BlockEvent) ([]byte, error) {
	block, err := getBlock(ev.Block)
	if err != nil {
		return nil, err
	}

	blockEvent := &BlockEvent{
		Number:       block.Number,
		Hash:         block.Hash,
		Size:         block.Size,
		PreviousHash: block.PreviousHash,
		DataHash:     block.DataHash,
		Txs:          block.Txs,
		Timestamp:    block.Timestamp,
		SourceURL:    ev.SourceURL,
		//SourceAlias:  convertPeerName(ev.SourceURL, endpointConfig),
	}

	return json.Marshal(blockEvent)
}

func (fc *FabClient) getBlockTxNum(channelID string, blkNum uint64) (uint64, error) {
	if err := fc.clientCheck("ledgerClient", channelID); err != nil {
		return 0, err
	}

	blk, err := fc.channels[channelID].ledgerClient.QueryBlock(blkNum)
	if err != nil {
		return 0, fmt.Errorf("Get block info error: %s", err.Error())
	}
	if blk == nil {
		return 0, fmt.Errorf("Get block info error: %s", "Block Not Found")
	}

	return uint64(len(blk.Data.Data)), nil
}

func (fc *FabClient) getTxNum(channelID string, txnum, sblk, eblk uint64) uint64 {
	var i uint64
	if sblk == 0 {
		i = sblk
	} else {
		i = sblk + 1
	}
	for i < eblk {
		num, _ := fc.getBlockTxNum(channelID, i)
		txnum += num
		i++
	}
	return txnum
}

func getAffiliation(affis []mspclient.AffiliationInfo) string {
	var buffer bytes.Buffer

	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	for _, affi := range affis {
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("{")

		buffer.WriteString("\"name\":")
		buffer.WriteString("\"")
		buffer.WriteString(affi.Name)
		buffer.WriteString("\"")
		if len(affi.Affiliations) > 0 {
			buffer.WriteString(",\"child\":")
			buffer.WriteString(getAffiliation(affi.Affiliations))
		}

		buffer.WriteString("}")

		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")

	return buffer.String()
}

func getAttribute(attrs []mspclient.Attribute) string {
	var buffer bytes.Buffer

	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	for _, attr := range attrs {
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("{")

		buffer.WriteString("\"Name\":")
		buffer.WriteString("\"")
		buffer.WriteString(attr.Name)
		buffer.WriteString("\"")
		buffer.WriteString(",")

		buffer.WriteString("\"Value\":")
		buffer.WriteString("\"")
		buffer.WriteString(attr.Value)
		buffer.WriteString("\"")
		buffer.WriteString(",")

		buffer.WriteString("\"ECert\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(attr.ECert))
		buffer.WriteString("\"")

		buffer.WriteString("}")

		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")

	return buffer.String()
}

func getIdentity(ids []*mspclient.IdentityResponse) string {
	var buffer bytes.Buffer

	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	for _, id := range ids {
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("{")

		buffer.WriteString("\"ID\":")
		buffer.WriteString("\"")
		buffer.WriteString(id.ID)
		buffer.WriteString("\"")
		buffer.WriteString(",")

		buffer.WriteString("\"Secret\":")
		buffer.WriteString("\"")
		buffer.WriteString(id.Secret)
		buffer.WriteString("\"")
		buffer.WriteString(",")

		buffer.WriteString("\"Type\":")
		buffer.WriteString("\"")
		buffer.WriteString(id.Type)
		buffer.WriteString("\"")
		buffer.WriteString(",")

		buffer.WriteString("\"Affiliation\":")
		buffer.WriteString("\"")
		buffer.WriteString(id.Affiliation)
		buffer.WriteString("\"")
		buffer.WriteString(",")

		buffer.WriteString("\"CAName\":")
		buffer.WriteString("\"")
		buffer.WriteString(id.CAName)
		buffer.WriteString("\"")
		buffer.WriteString(",")

		buffer.WriteString("\"MaxEnrollments\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.Itoa(id.MaxEnrollments))
		buffer.WriteString("\"")
		buffer.WriteString(",")

		buffer.WriteString("\"Attributes\":")
		// buffer.WriteString("\"")
		buffer.WriteString(getAttribute(id.Attributes))
		// buffer.WriteString("\"")

		buffer.WriteString("}")

		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")

	return buffer.String()
}

// 获取连接客户端的组织和用户名
func (fc *FabClient) GetClientInfo() (string, string) {
	return fc.orgName, fc.userName
}

// 查询CA信息
func (fc *FabClient) QueryCa() ([]byte, error) {
	if err := fc.clientCheck("mspClient", ""); err != nil {
		return nil, err
	}

	cainfo, err := fc.mspClient.GetCAInfo()
	if err != nil {
		return nil, err
	}

	affis, _ := fc.mspClient.GetAllAffiliations()

	affiliations := getAffiliation(affis.Affiliations)

	ids, _ := fc.mspClient.GetAllIdentities()

	identities := getIdentity(ids)

	ca := &CA{
		CaName:                    cainfo.CAName,
		CaChain:                   base64.StdEncoding.EncodeToString(cainfo.CAChain),
		Version:                   cainfo.Version,
		Affiliations:              affiliations,
		Identities:                identities,
		IssuerPublicKey:           base64.StdEncoding.EncodeToString(cainfo.IssuerPublicKey),
		IssuerRevocationPublicKey: base64.StdEncoding.EncodeToString(cainfo.IssuerRevocationPublicKey),
	}

	/*
		var buffer bytes.Buffer

		buffer.WriteString("{")

		buffer.WriteString("\"CAName\":")
		buffer.WriteString("\"")
		buffer.WriteString(cainfo.CAName)
		buffer.WriteString("\"")
		buffer.WriteString(",")

		buffer.WriteString("\"CAChain\":")
		buffer.WriteString("\"")
		buffer.WriteString(base64.StdEncoding.EncodeToString(cainfo.CAChain))
		buffer.WriteString("\"")
		buffer.WriteString(",")

		buffer.WriteString("\"IssuerPublicKey\":")
		buffer.WriteString("\"")
		buffer.WriteString(base64.StdEncoding.EncodeToString(cainfo.IssuerPublicKey))
		buffer.WriteString("\"")
		buffer.WriteString(",")

		buffer.WriteString("\"IssuerRevocationPublicKey\":")
		buffer.WriteString("\"")
		buffer.WriteString(base64.StdEncoding.EncodeToString(cainfo.IssuerRevocationPublicKey))
		buffer.WriteString("\"")
		buffer.WriteString(",")

		buffer.WriteString("\"Version\":")
		buffer.WriteString("\"")
		buffer.WriteString(cainfo.Version)
		buffer.WriteString("\"")
		buffer.WriteString(",")

		buffer.WriteString("\"Affiliations\":")
		// buffer.WriteString("\"")
		buffer.WriteString(affiliations)
		// buffer.WriteString("\"")
		buffer.WriteString(",")

		buffer.WriteString("\"Identities\":")
		// buffer.WriteString("\"")
		buffer.WriteString(identities)
		// buffer.WriteString("\"")

		buffer.WriteString("}")
	*/

	return json.Marshal(ca)

}

// 通过ca注册用户
func (fc *FabClient) Register(iname, itype, affi, scrt string, maxenroll int, attrs ...[3]string) (string, error) {

	if err := fc.clientCheck("mspClient", ""); err != nil {
		return "", err
	}
	/*
		_, _, endpointConfig := fc.getConfig()
		if endpointConfig == nil {
			return fmt.Errorf("get end porint config failed")
		}
		orgs := endpointConfig.NetworkConfig().Organizations
		fmt.Println(orgs)
		for _, ca := range orgs["org1"].CertificateAuthorities {
			fmt.Println(ca)
		}
	*/
	cainfo, err := fc.mspClient.GetCAInfo()
	if err != nil {
		return "", err
	}

	var attributes []mspclient.Attribute
	for _, attr := range attrs {
		ecert, err := strconv.ParseBool(attr[2])
		if err != nil {
			return "", err
		}
		attributes = append(attributes, mspclient.Attribute{Name: attr[0], Value: attr[1], ECert: ecert})
	}

	request := &mspclient.RegistrationRequest{
		// Name is the unique name of the identity
		Name: iname,
		// Type of identity being registered (e.g. "peer, app, user")
		Type: itype,
		// MaxEnrollments is the number of times the secret can  be reused to enroll.
		// if omitted, this defaults to max_enrollments configured on the server
		// The identity's affiliation e.g. org1.department1
		Affiliation: affi,
		// Optional attributes associated with this identity
		Attributes: attributes,
		// CAName is the name of the CA to connect to
		CAName: cainfo.CAName,
		// Secret is an optional password.  If not specified,
		// a random secret is generated.  In both cases, the secret
		// is returned from registration.
		Secret: scrt,
	}

	if maxenroll > 0 {
		request.MaxEnrollments = maxenroll
	}

	return fc.mspClient.Register(request)
}

// 登记用户
func (fc *FabClient) Enroll(enrollID, scrt string) (string, error) {

	if err := fc.clientCheck("mspClient", ""); err != nil {
		return "", err
	}

	// WithOrg 指定msp绑定的org
	// WithCA 指定使用哪个CA服务器进行enroll
	// WithProfile 指定是否启用tls通信 WithProfile("tls")
	// WithAttributeRequests 指定用户的哪些属性需要写入证书中
	//fc.mspClient.Enroll(enrollmentID, mspclient.WithSecret(scrt), mspclient.WithOrg(orgname), mspclient.WithCA(caname), mspclient.WithAttributeRequests(attrReqs))
	err := fc.mspClient.Enroll(enrollID, mspclient.WithSecret(scrt))

	if err == nil {
		return enrollID, fc.credentialIntegrate(enrollID)
	}

	return enrollID, err
}

// 将ca生成的密钥和证书整合到yaml配置文件中指定的目录中, 方便统一调用
func (fc *FabClient) credentialIntegrate(enrollID string) error {
	cryptoSuiteConfig, identityConfig, endpointConfig := fc.getConfig()
	keyStore := cryptoSuiteConfig.KeyStorePath()      // 私钥目录
	userStore := identityConfig.CredentialStorePath() // 证书目录
	// id, _ := fc.mspClient.GetIdentity(enrollID) // 需要从CA读取
	sid, _ := fc.mspClient.GetSigningIdentity(enrollID)

	mspid := sid.Identifier().MSPID

	certBytes := sid.EnrollmentCertificate()
	blk, _ := pem.Decode(certBytes)
	cert, _ := x509.ParseCertificate(blk.Bytes)

	// affis := strings.Split("com.example.org1", ".")
	affi := cert.Subject.OrganizationalUnit[1:] // 从证书中读出affilliation (fabric sdk go 会将affiliation拼凑在OU中)
	for i, j := 0, len(affi)-1; i < j; i, j = i+1, j-1 {
		affi[i], affi[j] = affi[j], affi[i]
	}

	cryptoPath := endpointConfig.CryptoConfigPath()

	userPath := enrollID + "@" + strings.Join(affi, ".")

	keyFilename := hex.EncodeToString(sid.PrivateKey().SKI()) + "_sk"

	certFilename := enrollID + "@" + mspid + "-cert.pem"

	keyFullname := filepath.Join(keyStore, keyFilename)

	certFullname := filepath.Join(userStore, certFilename)
	// fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++++")
	// fmt.Println(userPath, keyFullname, certFullname)
	var orgCryptoPath string
	for _, v := range endpointConfig.NetworkConfig().Organizations {
		if v.MSPID == mspid {
			orgCryptoPath = v.CryptoPath
			break
		}
	}

	compile := regexp.MustCompile(`^.*/users/`)

	orgRootPath := strings.TrimSuffix(filepath.Join(cryptoPath, compile.FindString(orgCryptoPath)), "/users")

	orgCAPath := filepath.Join(orgRootPath, "/msp/cacerts")
	orgTlsCAPath := filepath.Join(orgRootPath, "/msp/tlscacerts")

	userCryptoPath := filepath.Join(orgRootPath, "users", userPath)
	if !filepath.IsAbs(userCryptoPath) {
		return fmt.Errorf("user crypto path not exist")
	}

	caStore := filepath.Join(userCryptoPath, "cacerts")
	iKeyStore := filepath.Join(userCryptoPath, "keystore")
	aUserStore := filepath.Join(userCryptoPath, "admincerts")
	sUserStore := filepath.Join(userCryptoPath, "signcerts")
	tlsCaStore := filepath.Join(userCryptoPath, "tlscacerts")

	// fmt.Println(iKeyStore, aUserStore, sUserStore, orgCAPath, orgTlsCAPath)

	/*
		if _, err := exec.Command("mkdir", userCryptoPath).Output(); err != nil {
			return
		}
	*/
	if err := os.MkdirAll(caStore, os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(tlsCaStore, os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(iKeyStore, os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(aUserStore, os.ModePerm); err != nil {
		return err
	}

	if err := os.MkdirAll(sUserStore, os.ModePerm); err != nil {
		return err
	}

	/*
		_, err := exec.Command("cp", "-r", orgCAPath, userCryptoPath).Output()
		if err != nil {
			return
		}
	*/

	if err := copyDir(orgCAPath, caStore); err != nil {
		return err
	}

	if err := copyDir(orgTlsCAPath, tlsCaStore); err != nil {
		return err
	}

	if _, err := copyFile(keyFullname, filepath.Join(iKeyStore, keyFilename)); err != nil {
		return err
	}

	if _, err := copyFile(certFullname, filepath.Join(aUserStore, certFilename)); err != nil {
		return err
	}

	if _, err := copyFile(certFullname, filepath.Join(sUserStore, certFilename)); err != nil {
		return err
	}

	return nil
}

// 获取链信息
func (fc *FabClient) GetChain(channelID string) (*Chain, error) {
	if err := fc.clientCheck("ledgerClient", channelID); err != nil {
		//return nil, err
		return nil, err
	}

	lgc, err := fc.channels[channelID].ledgerClient.QueryInfo(ledger.WithTimeout(fab.ChannelConfigRefresh, 2*time.Second))
	if err != nil {
		//return nil, fmt.Errorf("Get ledger info error: %s", err.Error())
		return nil, err
	}

	chain := &Chain{
		// CurBlockHash: base64.StdEncoding.EncodeToString(lgc.BCI.CurrentBlockHash),
		CurBlockHash: fmt.Sprintf("%x", lgc.BCI.CurrentBlockHash),
		PreBlockHash: fmt.Sprintf("%x", lgc.BCI.PreviousBlockHash),
		Height:       lgc.BCI.Height,
		Endorser:     lgc.Endorser,
		Status:       strconv.Itoa(int(lgc.Status)),
	}

	return chain, nil
}

// 查询链信息
func (fc *FabClient) QueryChain(channelID string) ([]byte, error) {
	chain, err := fc.GetChain(channelID)
	if err != nil && chain != nil {
		return nil, err
	}

	// TODO  TxNum

	return json.Marshal(chain)
}

// 获取通道信息
func (fc *FabClient) GetChannel(channelID, orderer string) (fab.ChannelCfg, error) {
	if err := fc.clientCheck("resMgmtClient", ""); err != nil {
		return nil, err
	}

	//_, _, endpointConfig := fc.getConfig()

	/*
		orderersconfig := endpointConfig.ChannelOrderers(channelID)
		fmt.Println(orderersconfig)
		for _, orderercfg := range orderersconfig {
			fmt.Println(orderercfg.URL)
		}
	*/
	//fmt.Println(endpointConfig.ChannelConfig(channelID).Orderers)

	var cfg fab.ChannelCfg
	var err error

	// get channel config by orderer
	if orderer == "" {
		cfg, err = fc.resMgmtClient.QueryConfigFromOrderer(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	} else {
		cfg, err = fc.resMgmtClient.QueryConfigFromOrderer(channelID, resmgmt.WithOrdererEndpoint(orderer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	}

	// get channel config by ledger
	if err != nil || cfg == nil {
		if err := fc.clientCheck("ledgerClient", channelID); err != nil {
			//return nil, err
			return nil, err
		}

		cfg, err = fc.channels[channelID].ledgerClient.QueryConfig()
		if err != nil {
			//return nil, fmt.Errorf("Get config info error: %s", err.Error())
			return nil, err
		}
	}

	return cfg, nil
}

// 查询通道信息
func (fc *FabClient) QueryChannel(channelID, orderer string) ([]byte, error) {
	cfg, err := fc.GetChannel(channelID, orderer)
	if err != nil {
		return nil, err
	}

	var anchorPeers []string
	for _, peer := range cfg.AnchorPeers() { // 获取锚节点
		anchorPeers = append(anchorPeers, peer.Host+"/"+peer.Org+":"+strconv.Itoa(int(peer.Port)))
	}

	var msps []string
	for _, msp := range cfg.MSPs() { // 获取加入通道的组织ID
		//mspstr := string(msp.Config)
		certStart := bytes.IndexAny(msp.Config, "-----BEGIN")
		mspid := regexp.MustCompile(`\w+`).FindString(string(msp.Config[:certStart]))
		msps = append(msps, mspid)
	}

	var orderers []string // 获取排序节点
	for _, orderer := range cfg.Orderers() {
		orderers = append(orderers, orderer)
	}
	// orderers = strings.TrimRight(orderers, ",")

	ordererGroup := cfg.Versions().Channel.GetGroups()["Orderer"]
	consensus := ordererGroup.GetValues()["ConsensusType"].String()

	config := &Config{
		ChannelID:       cfg.ID(),
		AnchorPeers:     anchorPeers,
		LastConfigBlock: cfg.BlockNumber(),
		Msps:            msps,
		Orderers:        orderers,
		Consensus:       consensus,
	}
	return json.Marshal(config)
	//return config, nil
}

// create channel
func (fc *FabClient) CreateChannel(channelID, channelConfigPath, orderer string) error {

	if err := fc.clientCheck("resMgmtClient", ""); err != nil {
		return err
	}

	if err := fc.clientCheck("mspClient", ""); err != nil {
		return err
	}

	//var lastConfigBlock uint64

	adminIdentity, err := fc.mspClient.GetSigningIdentity(fc.userName) // should be admin
	if err != nil {
		return fmt.Errorf("Get signing identity error: %v", err)
	}

	channelReq := resmgmt.SaveChannelRequest{
		ChannelID:         channelID,
		ChannelConfigPath: channelConfigPath,
		SigningIdentities: []msp.SigningIdentity{adminIdentity},
	}

	if orderer == "" {
		_, err = fc.resMgmtClient.SaveChannel(channelReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	} else {
		_, err = fc.resMgmtClient.SaveChannel(channelReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(orderer))
	}
	if err != nil {
		return fmt.Errorf("Create application channel error: %v", err)
	}

	// fmt.Println("Create channel successfully")

	// info.OrgResMgmt = resMgmtClient
	/*
		err = info.OrgResMgmt.JoinChannel(info.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererName))
		if err != nil {
			return fmt.Errorf("Join channel error: %v", err)
		}
	*/
	//_, err = fc.WaitForOrdererConfigUpdate(channelID, ordererName, true, lastConfigBlock)
	//if err != nil {
	//	return err
	//}

	return nil
}

// 加入通道
// 仅当前组织的节点可以被加入通道, 其他组织的节点需通过相应组织的管理员操作加入
func (fc *FabClient) JoinChannel(channelID, orderer, peers string) error {

	if err := fc.clientCheck("resMgmtClient", ""); err != nil {
		return err
	}

	peers = strings.Replace(peers, " ", "", -1) // 去除字符串中空格

	if peers != "" {
		// 仅本组织的指定节点加入通道, 无论连接配置文件中有多少组织和节点, 当前组织实例化出来的client只能操作本组织的节点
		if orderer != "" {
			return fc.resMgmtClient.JoinChannel(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(orderer), resmgmt.WithTargetEndpoints(strings.Split(peers, ",")...))
		}
		return fc.resMgmtClient.JoinChannel(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithTargetEndpoints(strings.Split(peers, ",")...))
	}

	// 本组织的所有节点加入通道
	if orderer != "" {
		return fc.resMgmtClient.JoinChannel(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(orderer))
	}

	return fc.resMgmtClient.JoinChannel(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
}

// 设置组织在通道中的锚节点(用于节点发现)
// 锚节点设置类似通道创建, 均是提交配置更新提案
func (fc *FabClient) SetAnchorPeer(channelID, channelAnchorPath, orderer string) error {
	if err := fc.clientCheck("resMgmtClient", ""); err != nil {
		return err
	}
	if err := fc.clientCheck("mspClient", ""); err != nil {
		return err
	}

	/*
		var lastConfigBlock uint64
		lastConfigBlock, err := fc.WaitForOrdererConfigUpdate(channelID, ordererName, false, lastConfigBlock)
		if err != nil {
			return err
		}
	*/

	adminIdentity, err := fc.mspClient.GetSigningIdentity(fc.userName) // should be admin
	if err != nil {
		return fmt.Errorf("Get signing identity error: %v", err)
	}

	channelReq := resmgmt.SaveChannelRequest{
		ChannelID:         channelID,
		ChannelConfigPath: channelAnchorPath,
		SigningIdentities: []msp.SigningIdentity{adminIdentity},
	}

	_, err = fc.resMgmtClient.SaveChannel(channelReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(orderer))
	if err != nil {
		return fmt.Errorf("Create application channel error: %v", err)
	}

	// fmt.Println("Create channel successfully")

	// info.OrgResMgmt = resMgmtClient
	/*
		err = info.OrgResMgmt.JoinChannel(info.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererName))
		if err != nil {
			return fmt.Errorf("Join channel error: %v", err)
		}
	*/

	/*
		_, err = fc.WaitForOrdererConfigUpdate(channelID, ordererName, false, lastConfigBlock)
		if err != nil {
			return err
		}
	*/

	return nil
}

// 获取节点信息
func (fc *FabClient) GetPeers(channelID, orgName string) (map[string]*Peer, error) {
	peers := make(map[string]*Peer)
	if err := fc.clientCheck("discClient", ""); err != nil {
		return nil, err
	}

	// 动态发现节点, fabric-sdk-go中默认不对外暴露NewRequest函数, 需要手动修改, 添加该函数
	req := discovery.NewRequest().AddLocalPeersQuery().OfChannel(channelID).AddPeersQuery()

	_, _, endpointConfig := fc.getConfig()
	targets, ok := endpointConfig.PeersConfig(orgName)
	if !ok {
		return nil, fmt.Errorf("Peers not found")
	}

	/*
		var targets []fab.PeerConfig
			for _, v := range fc.fabnet.Peers {
				target := fab.PeerConfig{
					URL:         v.Url,
					GRPCOptions: v.GrpcOptions,
				}
				targets = append(targets, target)
			}
	*/

	/*
		grpcOptions := map[string]interface{}{
			"allow-insecure": true,
		}

		target1 := fab.PeerConfig{
			URL:         "peer1.org1.example.com:7051",
			GRPCOptions: grpcOptions,
		}
		target2 := fab.PeerConfig{
			URL:         "peer1.org2.example.com:7051",
			GRPCOptions: grpcOptions,
		}
		target3 := fab.PeerConfig{
			URL:         "peer1.org2.example.com:7051",
			GRPCOptions: grpcOptions,
		}
	*/

	ctx, cancle := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancle()

	responses, err := fc.discClient.Send(ctx, req, targets...) // 服务发现机制可以发现通道内所有节点(包括其他组织的节点), 前提是组织在通道上有锚节点
	if err != nil {
		return nil, err
	}
	//responses, err := client.Send(ctx, req)
	response := responses[0]
	chResp := response.ForChannel(channelID)
	prs, _ := chResp.Peers()
	for _, pr := range prs {
		//payload := peer.AliveMessage.Envelope.Payload
		//log.Println(string(payload))
		var cert *x509.Certificate
		certBytes := pr.Identity
		// 获取证书部分
		certStart := bytes.IndexAny(certBytes, "-----BEGIN")
		if certStart >= 0 {
			certText := certBytes[certStart:]
			// 获取证书的pem格式编码
			bl, _ := pem.Decode(certText)
			if bl != nil {
				// 生成证书对象
				if ct, err := x509.ParseCertificate(bl.Bytes); err == nil {
					cert = ct
				}
			}

		}

		endPoint := pr.AliveMessage.GetAliveMsg().Membership.Endpoint
		mspid := pr.MSPID
		name := strings.Split(endPoint, ":")[0]

		peer := &Peer{
			Name:   name,
			Url:    endPoint,
			Cert:   cert,
			MspID:  mspid,
			Status: "Online",
		}

		peers[peer.Name] = peer
	}

	return peers, nil
}

// 查询节点信息
func (fc *FabClient) QueryPeers(channelID, orgName string) ([]byte, error) {
	if orgName == "" {
		orgName = fc.orgName
	}

	peers, err := fc.GetPeers(channelID, orgName)
	if err != nil && peers != nil {
		return nil, err
	}

	return json.Marshal(peers) // 返回通道中的所有链码信息
}

// 查询节点加入的通道
func (fc *FabClient) QueryPeerChannel(peer string) ([]string, error) {
	if err := fc.clientCheck("resMgmtClient", ""); err != nil {
		return nil, err
	}

	channels := []string{}

	resp, err := fc.resMgmtClient.QueryChannels(resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithTargetEndpoints(peer))
	if err != nil {
		return nil, err
	}

	for _, channel := range resp.Channels {
		channels = append(channels, channel.ChannelId)
	}
	return channels, nil
}

// 获取链码信息
func (fc *FabClient) GetChaincode(channelID string, peer string) (map[string]*Chaincode, error) {
	chaincodes := make(map[string]*Chaincode)
	if err := fc.clientCheck("resMgmtClient", channelID); err != nil {
		return chaincodes, err
	}

	installResp, err := fc.resMgmtClient.QueryInstalledChaincodes(resmgmt.WithTargetEndpoints(peer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return chaincodes, err
	}

	for _, ccinfo := range installResp.Chaincodes {
		chaincode := &Chaincode{
			ID:      fmt.Sprintf("%x", ccinfo.Id),
			Name:    ccinfo.Name,
			Path:    ccinfo.Path,
			Version: ccinfo.Version,
			Input:   ccinfo.Input,
			Escc:    ccinfo.Escc,
			Vscc:    ccinfo.Vscc,
			Status:  "Installed",
		}

		chaincodes[chaincode.Name] = chaincode
	}

	instantResp, err := fc.resMgmtClient.QueryInstantiatedChaincodes(channelID, resmgmt.WithTargetEndpoints(peer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return chaincodes, err
	}

	for _, ccinfo := range instantResp.Chaincodes {
		chaincode := &Chaincode{
			ID:      fmt.Sprintf("%x", ccinfo.Id),
			Name:    ccinfo.Name,
			Path:    ccinfo.Path,
			Version: ccinfo.Version,
			Input:   ccinfo.Input,
			Escc:    ccinfo.Escc,
			Vscc:    ccinfo.Vscc,
			Status:  "Instantiated",
		}
		chaincodes[chaincode.Name] = chaincode
	}

	return chaincodes, nil
}

// 查询链码信息
func (fc *FabClient) QueryChaincode(channelID, chaincodeID, peer string) ([]byte, error) {
	chaincodes, err := fc.GetChaincode(channelID, peer)
	if err != nil && chaincodes != nil {
		return nil, err
	}

	if chaincodeID != "" {
		if chaincode, ok := chaincodes[chaincodeID]; ok {
			return json.Marshal(chaincode)
		}
		return nil, fmt.Errorf("Chaincode not found")
	}

	return json.Marshal(chaincodes) // 返回通道中的所有链码信息

}

// 安装链码
func (fc *FabClient) InstallChaincode(chaincodeID, chaincodeVersion, chaincodeGoPath, chaincodePath, peers string) error {
	if err := fc.clientCheck("resMgmtClient", ""); err != nil {
		return err
	}
	if chaincodeGoPath == "" {
		chaincodeGoPath = os.Getenv("GOPATH")
	}

	// fmt.Println("Start to install chainocde ...")
	ccPkg, err := gopackager.NewCCPackage(chaincodePath, chaincodeGoPath)
	if err != nil {
		return fmt.Errorf("Build chaincode package error: %s", err.Error())
	}

	installCCReq := resmgmt.InstallCCRequest{
		Name:    chaincodeID,
		Path:    chaincodePath,
		Version: chaincodeVersion,
		Package: ccPkg,
	}

	peers = strings.Replace(peers, " ", "", -1) // 去除字符串中空格

	if peers != "" {
		// 仅本组织的指定节点安装链码
		_, err = fc.resMgmtClient.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithTargetEndpoints(strings.Split(peers, ",")...))
	} else {
		// 本组织的所有节点安装链码
		_, err = fc.resMgmtClient.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	}

	if err != nil {
		return fmt.Errorf("Install chaincode error: %s", err.Error())
	}

	return nil
	// fmt.Println("Install chaincode successfully")
}

// 实例化链码
func (fc *FabClient) InstantiateChaincode(channelID, chaincodeID, chaincodeVersion, chaincodePath, policy, peers string, args []string) error {
	if err := fc.clientCheck("resMgmtClient", ""); err != nil {
		return err
	}

	if len(args) < 1 {
		return fmt.Errorf("Incorrect number of arguments. Expecting >0") // 第一个参数应该是init
	}

	var argBytes [][]byte

	for _, arg := range args {
		argBytes = append(argBytes, []byte(arg))
	}

	// ccPolicy, _ := cauthdsl.FromString("or('Org1MSP.member')") // 参数是与命令行相同背书策略语法的字符串
	ccPolicy, err := cauthdsl.FromString(policy)
	if err != nil {
		return fmt.Errorf("Endorsement policy is incorrect: %s", err)
	}

	instantiateCCReq := resmgmt.InstantiateCCRequest{
		Name:    chaincodeID,
		Path:    chaincodePath,
		Version: chaincodeVersion,
		// Args:    [][]byte{[]byte("init")},
		Args:   argBytes,
		Policy: ccPolicy,
	}

	peers = strings.Replace(peers, " ", "", -1) // 去除字符串中空格

	if peers != "" {
		// 仅本组织的指定节点实例化链码
		_, err = fc.resMgmtClient.InstantiateCC(channelID, instantiateCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithTargetEndpoints(strings.Split(peers, ",")...))
	} else {
		// 本组织的所有节点实例化链码
		_, err = fc.resMgmtClient.InstantiateCC(channelID, instantiateCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	}

	if err != nil {
		return fmt.Errorf("Instantiate chaincode error: %s", err.Error())
	}

	// fmt.Println("Instantiate chaincode successfully")

	return nil
}

// 升级链码
func (fc *FabClient) UpgradeChaincode(channelID, chaincodeID, chaincodeVersion, chaincodePath, policy, peers string, args []string) error {
	if err := fc.clientCheck("resMgmtClient", ""); err != nil {
		return err
	}

	if len(args) < 1 {
		return fmt.Errorf("Incorrect number of arguments. Expecting >0") // 第一个参数应该是init
	}

	var argBytes [][]byte

	for _, arg := range args {
		argBytes = append(argBytes, []byte(arg))
	}

	// ccPolicy, _ := cauthdsl.FromString("or('Org1MSP.member')") // 参数是与命令行相同背书策略语法的字符串
	ccPolicy, err := cauthdsl.FromString(policy)
	if err != nil {
		return fmt.Errorf("Endorsement policy is incorrect: %s", err)
	}

	upgradeCCReq := resmgmt.UpgradeCCRequest{
		Name:    chaincodeID,
		Path:    chaincodePath,
		Version: chaincodeVersion,
		// Args:    [][]byte{[]byte("init")},
		Args:   argBytes,
		Policy: ccPolicy,
	}

	peers = strings.Replace(peers, " ", "", -1) // 去除字符串中空格

	if peers != "" {
		// 仅本组织的指定节点实例化链码
		_, err = fc.resMgmtClient.UpgradeCC(channelID, upgradeCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithTargetEndpoints(strings.Split(peers, ",")...))
	} else {
		// 本组织的所有节点实例化链码
		_, err = fc.resMgmtClient.UpgradeCC(channelID, upgradeCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	}

	if err != nil {
		return fmt.Errorf("Instantiate chaincode error: %s", err.Error())
	}

	// fmt.Println("Instantiate chaincode successfully")

	return nil
}

// 获取块信息
func (fc *FabClient) GetBlock(channelID, by, arg string) (*Block, error) {
	if err := fc.clientCheck("ledgerClient", channelID); err != nil {
		return nil, err
	}

	var blk *cb.Block
	var err error
	by = strings.ToLower(by)
	if by == "" || by == "number" {
		num, err := strconv.ParseUint(arg, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Illegal argument error: %s", err.Error())
		}
		blk, err = fc.channels[channelID].ledgerClient.QueryBlock(num)
	} else if by == "hash" {
		hash, err := hex.DecodeString(arg)
		if err != nil {
			return nil, fmt.Errorf("Illegal argument error: %s", err.Error())
		}
		blk, err = fc.channels[channelID].ledgerClient.QueryBlockByHash(hash)
	} else if by == "txhash" {
		blk, err = fc.channels[channelID].ledgerClient.QueryBlockByTxID(fab.TransactionID(arg))
	} else {
		return nil, fmt.Errorf("Illegal argument error: %s", err.Error())
	}

	if err != nil {
		return nil, fmt.Errorf("Get block info error: %s", err.Error())
	}

	if blk == nil {
		return nil, fmt.Errorf("Get block info error: %s", "Block Not Found")
	}

	//log.Println(string(data.Data[0]))

	return getBlock(blk)
}

// 查询指定块信息(按块号,按块哈希, 按交易哈希查询)
func (fc *FabClient) QueryBlock(channelID, by, arg string) ([]byte, error) {
	block, err := fc.GetBlock(channelID, by, arg)
	if err != nil && block != nil {
		return nil, err
	}

	return json.Marshal(block)
}

// 查询多个块信息(块高度height之前的quantity个块)
func (fc *FabClient) QueryBlocks(channelID, height, quantity string) ([]byte, error) {
	var end, num uint64
	num, err := strconv.ParseUint(quantity, 10, 64)
	if err != nil {
		return nil, err
	}

	var blocks []*Block
	if err := fc.clientCheck("ledgerClient", channelID); err != nil {
		return nil, err
	}

	lgc, err := fc.channels[channelID].ledgerClient.QueryInfo(ledger.WithTimeout(fab.ChannelConfigRefresh, 2*time.Second))
	if err != nil {
		//return nil, fmt.Errorf("Get ledger info error: %s", err.Error())
		return nil, err
	}

	if height == "-1" {
		end = lgc.BCI.Height
	} else {
		end, err = strconv.ParseUint(height, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	var start uint64
	if end == 0 {
		end = 1
	}
	if end > lgc.BCI.Height {
		end = lgc.BCI.Height
	}
	if end >= num {
		start = end - num
	}

	for end > start {
		block, err := fc.GetBlock(channelID, "", strconv.FormatUint(end-1, 10))
		if err == nil && block != nil {
			blocks = append(blocks, block)
		}
		end--
	}

	return json.Marshal(blocks)
}

// 获取交易信息
func (fc *FabClient) GetTransaction(channelID, arg string) (*Transaction, error) {
	if err := fc.clientCheck("ledgerClient", channelID); err != nil {
		return nil, err
	}

	blk, err := fc.channels[channelID].ledgerClient.QueryBlockByTxID(fab.TransactionID(arg))
	if err != nil {
		return nil, err
	}

	blkNum := blk.Header.Number

	ptx, err := fc.channels[channelID].ledgerClient.QueryTransaction(fab.TransactionID(arg))
	if err != nil {
		return nil, fmt.Errorf("Get transaction info error: %s", err.Error())
	}

	tx, err := getTransaction(ptx.GetTransactionEnvelope(), blkNum)
	if err != nil {
		return nil, fmt.Errorf("Get transaction info error: %s", err.Error())
	}

	return tx, nil
}

// 查询交易信息
func (fc *FabClient) QueryTransaction(channelID, arg string) ([]byte, error) {
	transaction, err := fc.GetTransaction(channelID, arg)
	if err != nil && transaction != nil {
		return nil, err
	}

	return json.Marshal(transaction)
}

// 查询多个交易信息
func (fc *FabClient) QueryTransactions(channelID, lasttx, quantity string) ([]byte, error) {
	num, err := strconv.ParseUint(quantity, 10, 64)
	if err != nil {
		return nil, err
	}

	if err := fc.clientCheck("ledgerClient", channelID); err != nil {
		return nil, err
	}

	var txs []*Transaction

	lgc, err := fc.channels[channelID].ledgerClient.QueryInfo(ledger.WithTimeout(fab.ChannelConfigRefresh, 2*time.Second))
	if err != nil {
		//return nil, fmt.Errorf("Get ledger info error: %s", err.Error())
		return nil, err
	}

	for i := lgc.BCI.Height; i > 0; i-- {
		block, err := fc.GetBlock(channelID, "", strconv.FormatUint(i-1, 10))
		if err == nil && block != nil {
			for j := len(block.Txs); j > 0; j-- {
				tx := block.Txs[j-1]
				if tx.TxID != lasttx && uint64(len(txs)) < num {
					txs = append(txs, tx)
				} else {
					break
				}
			}
			if uint64(len(txs)) == num {
				break
			}
		}
	}

	return json.Marshal(txs)
}

// 调用链码
func (fc *FabClient) Invoke(channelID, chaincodeID, peers, funcName string, args []string) ([]byte, error) {
	if err := fc.clientCheck("channelClient", channelID); err != nil {
		return nil, err
	}

	var argBytes [][]byte

	for _, arg := range args {
		argBytes = append(argBytes, []byte(strings.TrimSpace(arg)))
	}

	peers = strings.Replace(peers, " ", "", -1) // 去除字符串中空格, peer可以是节点名或节点URL, sdk会自动构建节点对象

	req := channel.Request{
		ChaincodeID: chaincodeID,
		Fcn:         funcName,
		Args:        argBytes,
	}

	var resp channel.Response
	var err error

	if peers != "" { // 如果网络连接配置文件yaml中没有交易背书策略要求的所有组织的节点配置信息, 则必须手动指定交易提案的分发节点
		// 指定交易提案分发节点
		// _, _, endpointConfig := fc.getConfig()
		/*
			resClientContext := fc.fabricSDK.Context(fabsdk.WithUser(fc.userName), fabsdk.WithOrg(fc.orgName))
			ctx, _ := resClientContext()
		*/
		ctx, _ := fc.fabricSDK.Context()()
		/*
			grpcOptions := map[string]interface{}{
				"allow-insecure": true,
			}
		*/
		var targets []fab.Peer

		for _, url := range strings.Split(peers, ",") { // 可以作为手动配置节点连接信息的范例
			var peerCfg *fab.NetworkPeer
			var err error

			peerCfg, err = comm.NetworkPeerConfig(ctx.EndpointConfig(), url)
			if err != nil {
				return nil, err
			}
			//fmt.Println(peerCfg.URL)
			//fmt.Println(peerCfg.GRPCOptions)

			/*
				peerCfg = &fab.PeerConfig{
					URL:         peerCfg.URL,
					GRPCOptions: grpcOptions,
				}
			*/

			peerCfg.GRPCOptions["allow-insecure"] = true // 允许非tls连接
			/*
				peerCfg := fab.PeerConfig{
					URL:         peerCfg.URL,
					GRPCOptions: grpcOptions,
				}
			*/
			peer, _ := ctx.InfraProvider().CreatePeerFromConfig(peerCfg)
			targets = append(targets, peer)
		}
		// 此处使用WithTargets而不使用WithTargetEndpoints
		// WithTargetEndpoints本质上是是自动构建target后调用WithTarget的
		// WithTargetEndpoints自动构建的target使用的都是默认值, 其中peerCfg的GRPCOptions中allow-insecure项的默认值是false, 即要求调用peer时使用tls连接才行, 为了在禁用tls的fabric网络中正常使用, 所以手动配置了target
		// 如果网络配置文件yaml中包含了peers中指定的所有节点的配置信息, 在该配置文件中禁用了tls, 则可以直接使用WithTargetEndpoints, 否则自动生成的节点配置信息使用的都是默认值
		// 真实生产环境中, 不同组织只拥有本组织的凭证和节点信息, 因此sdk是无法从本组织的网络配置文件中读取到其他组织的节点的配置信息的, 当链码背书策略中要求交易提案必须有其他组织节点的签名时, sdk为了向其他组织的节点发送交易提案, 就需要指定对应的节点名或节点url, 同时指定是否通过tls进行通信(如果通过tls通信则还需要提供对方的tlsca证书)
		// 如链码背书策略为"and('Org1MSP.member','Org2MSP.member)"时, 如果网络连接配置文件仅包含本组织(Org1)的配置信息, 则需要向sdk传递peers参数指定Org1和Org2的节点名或节点url, 不可省略本组织的节点

		resp, err = fc.channels[channelID].channelClient.Execute(req, channel.WithRetry(retry.DefaultChannelOpts), channel.WithTargets(targets...))
		// resp, err = fc.channels[channelID].channelClient.Execute(req, channel.WithRetry(retry.DefaultChannelOpts), channel.WithTargetEndpoints(strings.Split(peers, ",")...))
	} else {
		// 交易提案自动分发
		// 自动分发的前提是网络连接配置文件yaml中包含了链码背书策略中指定的组织的相应节点的配置信息, 否则背书失败
		resp, err = fc.channels[channelID].channelClient.Execute(req, channel.WithRetry(retry.DefaultChannelOpts))
	}
	// sdk封装的execute函数会自动收集背书, 等待链码背书策略条件达成, 无需手动指定提案分发节点

	//resp, err := fc.channelClient.Execute(req)
	//resp, err := fc.channelClient.InvokeHandler(invoke.NewExecuteHandler(), req)
	if err != nil {
		return nil, err
	}
	//fmt.Println(resp.TxValidationCode)
	//fmt.Println(resp.Responses[0].Status)
	//return string(resp.TransactionID) + ":" + resp.TxValidationCode.String(), nil
	return json.Marshal(resp)
}

// 查询链码
func (fc *FabClient) Query(channelID, chaincodeID, peers, funcName string, args []string) ([]byte, error) {
	if err := fc.clientCheck("channelClient", channelID); err != nil {
		return nil, err
	}

	var argBytes [][]byte

	for _, arg := range args {
		argBytes = append(argBytes, []byte(strings.TrimSpace(arg)))
	}

	peers = strings.Replace(peers, " ", "", -1) // 去除字符串中空格, peer可以是节点名或节点URL, sdk会自动构建节点对象

	req := channel.Request{
		ChaincodeID: chaincodeID,
		Fcn:         funcName,
		Args:        argBytes,
	}

	var resp channel.Response
	var err error

	if peers != "" { // 查询操作通常只需要在本组织的节点上查询即可, 无关背书策略, 一般而言本组织的网络连接配置文件中都会包含组织内节点的连接配置信息, 因此可以直接使用WithTargetEndpoints进行自动处理
		resp, err = fc.channels[channelID].channelClient.Query(req, channel.WithRetry(retry.DefaultChannelOpts), channel.WithTargetEndpoints(strings.Split(peers, ",")...))
	} else {
		resp, err = fc.channels[channelID].channelClient.Query(req, channel.WithRetry(retry.DefaultChannelOpts))
	}
	if err != nil {
		return nil, err
	}

	return json.Marshal(resp)

	// return string(resp.Payload), nil
}

// 块监视
func (fc *FabClient) BlockMonitor(channelID string, ctx context.Context, conn interface{}, handler func(interface{}, []byte, error)) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("block monitor error")
		}
	}()

	if err := fc.clientCheck("eventClient", channelID); err != nil {
		return
	}

	// _, _, endpointConfig := fc.getConfig()

	registration, eventch, err := fc.channels[channelID].eventClient.RegisterBlockEvent() //RegisterChaincodeEvent, RegisterFilteredBlockEvent, RegisterTxStatusEvent

	// fmt.Println("event client created")

	if err != nil {
		fmt.Println("Event Register error: ", err)
		return
	}
	defer fc.channels[channelID].eventClient.Unregister(registration)

	fmt.Println("Start monitor block event...")

	for {
		select {
		case ev, ok := <-eventch:
			if ok {
				//eventBytes, err := parseBlockEvent(ev, endpointConfig)
				eventBytes, err := parseBlockEvent(ev)
				handler(conn, eventBytes, err)
			}
		case <-ctx.Done():
			fmt.Println("Monitor stopped")
			return
		}
	}
}

// 链码监视
func (fc *FabClient) ChaincodeMonitor(channelID, chaincodeID, filter string, ctx context.Context, conn interface{}, handler func(interface{}, []byte, error)) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("chaincode monitor error")
		}
	}()

	if err := fc.clientCheck("eventClient", channelID); err != nil {
		return
	}

	registration, eventch, err := fc.channels[channelID].eventClient.RegisterChaincodeEvent(chaincodeID, filter) //RegisterChaincodeEvent, RegisterFilteredBlockEvent, RegisterTxStatusEvent

	// fmt.Println("event client created")

	if err != nil {
		fmt.Println("Event Register error: ", err)
		return
	}
	defer fc.channels[channelID].eventClient.Unregister(registration)

	fmt.Println("Start monitor chaincode event...")

	for {
		select {
		case ev, ok := <-eventch:
			if ok {
				eventBytes, err := json.Marshal(ev)
				handler(conn, eventBytes, err)
			}
		case <-ctx.Done():
			fmt.Println("Monitor stopped")
			return
		}
	}
}

// WaitForOrdererConfigUpdate waits until the config block update has been committed.
// In Fabric 1.0 there is a bug that panics the orderer if more than one config update is added to the same block.
// This function may be invoked after each config update as a workaround.
func (fc *FabClient) WaitForOrdererConfigUpdate(channelID, ordererName string, genesis bool, lastConfigBlock uint64) (uint64, error) {

	if fc.resMgmtClient == nil {
		return uint64(0), fmt.Errorf("ResMgmtClient is nil")
	}

	blockNum, err := retry.NewInvoker(retry.New(retry.TestRetryOpts)).Invoke(
		func() (interface{}, error) {
			chConfig, err := fc.resMgmtClient.QueryConfigFromOrderer(channelID, resmgmt.WithOrdererEndpoint(ordererName))
			if err != nil {
				return nil, status.New(status.TestStatus, status.GenericTransient.ToInt32(), err.Error(), nil)
			}

			currentBlock := chConfig.BlockNumber()
			if currentBlock <= lastConfigBlock && !genesis {
				return nil, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("Block number was not incremented [%d, %d]", currentBlock, lastConfigBlock), nil)
			}
			return &currentBlock, nil
		},
	)

	if err != nil {
		return uint64(0), err
	}
	return *blockNum.(*uint64), nil

	// 两种retry的使用方式

	/*
		chConfig, err := fc.resMgmtClient.QueryConfigFromOrderer(channelID, resmgmt.WithOrdererEndpoint(ordererName), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
		if err != nil {
			return uint64(0), err
		}

		currentBlock := chConfig.BlockNumber()
		if currentBlock <= lastConfigBlock && !genesis {
			return uint64(0), fmt.Errorf("Block number was not incremented [%d, %d]", currentBlock, lastConfigBlock)
		}
		return *blockNum.(*uint64), nil
	*/
}

/* 手动发起交易(有问题)
func (fc *FabClient) Test(channelID, chaincodeID, fcn string, args []string) (string, error) {
	if err := fc.clientCheck("transactor", channelID); err != nil {
		return "", err
	}

	var argBytes [][]byte

	for _, arg := range args {
		argBytes = append(argBytes, []byte(arg))
	}

	request := fab.ChaincodeInvokeRequest{
		ChaincodeID: chaincodeID,
		Fcn:         fcn,
		Args:        argBytes,
		//TransientMap: chrequest.TransientMap,
	}

	//(*fc.transactor).CreateTransactionHeader()
	//(*fc.transactor).SendTransactionProposal(*fab.TransactionProposal, []fab.ProposalProcessor)

	txh, err := (*fc.transactor).CreateTransactionHeader()
	if err != nil {
		return "", err
	}

	tpreq, err := txn.CreateChaincodeInvokeProposal(txh, request)
	if err != nil {
		return "", err
	}

	_, _, endpointConfig := fc.getConfig()

	peer1, err := peer.New(endpointConfig, peer.WithMSPID("Org1MSP"), peer.WithURL("localhost:7051"), peer.WithInsecure())
	if err != nil {
		return "", err
	}

	peer2, err := peer.New(endpointConfig, peer.WithMSPID("Org1MSP"), peer.WithURL("localhost:8051"), peer.WithInsecure())

	if err != nil {
		return "", err
	}
	peers := []fab.Peer{peer1, peer2}
	tpr, err := (*fc.transactor).SendTransactionProposal(tpreq, peer.PeersToTxnProcessors(peers))
	if err != nil {
		return "", err
	}

	fmt.Println(tpr, tpreq.TxnID)
	return "", err

	// tx, err := (*fc.transactor).CreateTransaction(request)
	// resp, err := (*fc.transactor).SendTransaction(tx)

}
*/
