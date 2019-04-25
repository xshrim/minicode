package fabcore

import (
	"bytes"
	"context"
	"encoding/asn1"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/discovery"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric/common/util"

	//disclient "github.com/hyperledger/fabric-sdk-go/internal/github.com/hyperledger/fabric/discovery/client"
	cb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
	putils "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/utils"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type FabClient struct {
	userName      string
	orgName       string
	channelID     string
	connNumber    int
	fabnet        *FabNet
	overViews     map[string]*OverView
	fabricSDK     *fabsdk.FabricSDK
	resMgmtClient *resmgmt.Client
	discClient    *discovery.Client
	mspClient     *mspclient.Client
	channelClient *channel.Client
	eventClient   *event.Client
	ledgerClient  *ledger.Client
}

// const ChaincodeVersion = "1.4"

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

// instantiate fabric sdk with config file
func Init(configFile string) (*fabsdk.FabricSDK, error) {
	sdk, err := fabsdk.New(config.FromFile(configFile))
	if err != nil {
		return nil, fmt.Errorf("Instantiate fabric sdk error: %v", err)
	}

	return sdk, nil
}

// build fabric client
func New(sdk *fabsdk.FabricSDK, userName, orgName, channelID string) (*FabClient, error) {
	if userName == "" || orgName == "" {
		return nil, fmt.Errorf("userName or orgName is empty")
	}

	client := &FabClient{
		userName:  userName,
		orgName:   orgName,
		channelID: channelID,
		overViews: make(map[string]*OverView),
	}

	if sdk == nil {
		sdk, err := Init(fmt.Sprintf("./config/%s-network.yaml", strings.ToLower(orgName)))
		if err != nil {
			return nil, err
		}
		client.fabricSDK = sdk
	} else {
		client.fabricSDK = sdk
	}

	client.fabnet = client.getFabNetInfo()

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

	//fab.ClientContext
	//disclient, err := discovery.New(client.fabricSDK.Context())
	if channelID != "" {

		channelClientContext := client.fabricSDK.ChannelContext(channelID, fabsdk.WithUser(userName), fabsdk.WithOrg(orgName))

		channelClient, err := channel.New(channelClientContext)
		if err == nil {
			client.channelClient = channelClient
		}

		eventClient, err := event.New(channelClientContext, event.WithBlockEvents())
		if err == nil {
			client.eventClient = eventClient
		}

		ledgerClient, err := ledger.New(channelClientContext)
		if err == nil {
			client.ledgerClient = ledgerClient
		}

	}

	// fmt.Println("Instantiate fabric sdk successfully")
	return client, nil
}

func (fc *FabClient) SetConnNumber(num int) {
	fc.connNumber += num
	if fc.connNumber <= 0 {
		fc.Close()
	}
}

func (fc *FabClient) Close() {
	fc.fabricSDK.Close()
	fc.fabricSDK = nil
	//log.Println(fc.fabricSDK)
}

func (fc *FabClient) IsClosed() bool {
	if fc == nil || fc.fabricSDK == nil {
		return true
	}
	return false
}

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
	} else if clientType == "ledgerClient" {
		if fc.ledgerClient == nil || fc.channelID != channelID {
			clientChannelContext := fc.fabricSDK.ChannelContext(channelID, fabsdk.WithUser(fc.userName), fabsdk.WithOrg(fc.orgName))

			ledgerClient, err := ledger.New(clientChannelContext)
			if err != nil {
				return fmt.Errorf("Create channel client error: %s", err.Error())
			}

			fc.channelID = channelID
			fc.ledgerClient = ledgerClient
		}
	} else if clientType == "channelClient" {
		if fc.channelClient == nil || fc.channelID != channelID {
			clientChannelContext := fc.fabricSDK.ChannelContext(channelID, fabsdk.WithUser(fc.userName), fabsdk.WithOrg(fc.orgName))

			channelClient, err := channel.New(clientChannelContext)
			if err != nil {
				return fmt.Errorf("Create channel client error: %s", err.Error())
			}

			fc.channelID = channelID
			fc.channelClient = channelClient
		}
	} else if clientType == "eventClient" {
		if fc.eventClient == nil || fc.channelID != channelID {
			clientChannelContext := fc.fabricSDK.ChannelContext(channelID, fabsdk.WithUser(fc.userName), fabsdk.WithOrg(fc.orgName))

			eventClient, err := event.New(clientChannelContext, event.WithBlockEvents())
			if err != nil {
				return fmt.Errorf("Create event client error: %s", err.Error())
			}

			fc.channelID = channelID
			fc.eventClient = eventClient
		}
	}

	return nil
}

// create channel
func (fc *FabClient) CreateChannel(channelID, channelConfigPath, ordererName string) error {
	if fc.resMgmtClient == nil || fc.mspClient == nil {
		return fmt.Errorf("ResMgmtClient or MspClient is nil")
	}

	adminIdentity, err := fc.mspClient.GetSigningIdentity(fc.userName) // should be admin
	if err != nil {
		return fmt.Errorf("Get signing identity error: %v", err)
	}

	channelReq := resmgmt.SaveChannelRequest{
		ChannelID:         channelID,
		ChannelConfigPath: channelConfigPath,
		SigningIdentities: []msp.SigningIdentity{adminIdentity},
	}

	_, err = fc.resMgmtClient.SaveChannel(channelReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(ordererName))
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

	return nil
}

// join channel
func (fc *FabClient) JoinChannel(channelID, ordererName, peerNames string) error {
	if fc.resMgmtClient == nil {
		return fmt.Errorf("ResMgmtClient is nil")
	}

	if peerNames != "" {
		return fc.resMgmtClient.JoinChannel(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(ordererName), resmgmt.WithTargetEndpoints(strings.Split(peerNames, ",")...))
	}

	return fc.resMgmtClient.JoinChannel(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(ordererName))
}

// install chaincode
func (fc *FabClient) InstallChaincode(chaincodePath, chaincodeGoPath, chaincodeID, chaincodeVersion, peerNames string) error {
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

	if peerNames != "" {
		_, err = fc.resMgmtClient.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithTargetEndpoints(strings.Split(peerNames, ",")...))
	} else {
		_, err = fc.resMgmtClient.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	}

	if err != nil {
		return fmt.Errorf("Install chaincode error: %s", err.Error())
	}

	return nil
	// fmt.Println("Install chaincode successfully")
}

// instantiate chaincode
func (fc *FabClient) InstantiateChaincode(channelID, chaincodeID, chaincodePath, chaincodeVersion, policy, peerNames string, args [][]byte) error {
	if err := fc.clientCheck("resMgmtClient", ""); err != nil {
		return err
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
		Args:   args,
		Policy: ccPolicy,
	}

	if peerNames != "" {
		_, err = fc.resMgmtClient.InstantiateCC(channelID, instantiateCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithTargetEndpoints(strings.Split(peerNames, ",")...))
	} else {
		_, err = fc.resMgmtClient.InstantiateCC(channelID, instantiateCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	}

	if err != nil {
		return fmt.Errorf("Instantiate chaincode error: %s", err.Error())
	}

	// fmt.Println("Instantiate chaincode successfully")

	return nil
}

/*
func InstallAndInstantiateCC(sdk *fabsdk.FabricSDK, info *FabInfo) error {
	clientContext := sdk.Context(fabsdk.WithUser(info.OrgAdmin), fabsdk.WithOrg(info.OrgName))
	if clientContext == nil {
		return fmt.Errorf("Create client context error")
	}

	fmt.Println("Start to install chainocde ...")
	ccPkg, err := gopackager.NewCCPackage(info.ChaincodePath, info.ChaincodeGoPath)
	if err != nil {
		return fmt.Errorf("Build chaincode package error: %s", err.Error())
	}

	installCCReq := resmgmt.InstallCCRequest{
		Name:    info.ChaincodeID,
		Path:    info.ChaincodePath,
		Version: info.ChaincodeVersion,
		Package: ccPkg,
	}

	if info.OrgResMgmt == nil {
		resMgmtClient, err := resmgmt.New(clientContext)
		if err != nil {
			return fmt.Errorf("Create resource management client error: %v", err)
		}
		info.OrgResMgmt = resMgmtClient
	}

	_, err = info.OrgResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return fmt.Errorf("Install chaincode error: %s", err.Error())
	}

	fmt.Println("Install chaincode successfully")

	fmt.Println("Start to instantiate chaincode ...")

	//ccPolicy := cauthdsl.SignedByAnyMember([]string{"Org1MSP"}) // 参数是组织MSPID的集合
	// ccPolicy := cauthdsl.AcceptAllPolicy  // 接受所有策略
	// ccPolicy := cauthdsl.Or(lhs, rhs)
	ccPolicy, _ := cauthdsl.FromString("or('Org1MSP.member')") // 参数是与命令行相同背书策略语法的字符串

	instantiateCCReq := resmgmt.InstantiateCCRequest{
		Name:    info.ChaincodeID,
		Path:    info.ChaincodePath,
		Version: info.ChaincodeVersion,
		Args:    [][]byte{[]byte("init")},
		Policy:  ccPolicy,
	}

	_, err = info.OrgResMgmt.InstantiateCC(info.ChannelID, instantiateCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return fmt.Errorf("Instantiate chaincode error: %s", err.Error())
	}

	fmt.Println("Instantiate chaincode successfully")

	return nil
}
*/

/*
func CreateChannelClient(sdk *fabsdk.FabricSDK, info *FabInfo) (*channel.Client, error) {
	clientChannelContext := sdk.ChannelContext(info.ChannelID, fabsdk.WithUser(info.UserName), fabsdk.WithOrg(info.OrgName))

	channelClient, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("Create channel client error: %s", err.Error())
	}

	fmt.Println("Create chaincode client successfully")

	return channelClient, nil
}
*/

// invoke chaincode
func (fc *FabClient) Invoke(channelID, chaincodeID, funcName string, args [][]byte) (string, error) {
	if err := fc.clientCheck("channelClient", channelID); err != nil {
		return "", err
	}

	req := channel.Request{
		ChaincodeID: chaincodeID,
		Fcn:         funcName,
		Args:        args,
	}

	resp, err := fc.channelClient.Execute(req, channel.WithRetry(retry.DefaultChannelOpts))
	//resp, err := fc.channelClient.Execute(req)
	//resp, err := fc.channelClient.InvokeHandler(invoke.NewExecuteHandler(), req)
	if err != nil {
		return "", err
	}
	//fmt.Println(resp.TxValidationCode)
	//fmt.Println(resp.Responses[0].Status)
	return string(resp.TransactionID) + ":" + resp.TxValidationCode.String(), nil
}

// query chaincode
func (fc *FabClient) Query(channelID, chaincodeID, funcName string, args [][]byte) ([]byte, error) {
	if err := fc.clientCheck("channelClient", channelID); err != nil {
		return nil, err
	}

	req := channel.Request{
		ChaincodeID: chaincodeID,
		Fcn:         funcName,
		Args:        args,
	}

	resp, err := fc.channelClient.Query(req, channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		return nil, err
	}

	return resp.Payload, nil

	// return string(resp.Payload), nil
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

func convertPeerName(sourceurl string, peers map[string]*CPeer) string {
	var fix string
	su := strings.Split(sourceurl, ":")
	if len(su) > 1 {
		fix = su[1]
	} else {
		fix = su[0]
	}
	for _, peer := range peers {
		if len(su) > 1 {
			if strings.HasSuffix(peer.Url, fix) {
				return peer.Name
			}
		} else {
			if strings.HasPrefix(peer.Url, fix) {
				return peer.Name
			}
		}
	}
	return sourceurl
}

func parseBlockEvent(ev *fab.BlockEvent, peers map[string]*CPeer) ([]byte, error) {
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
		SourceAlias:  convertPeerName(ev.SourceURL, peers),
	}

	return json.Marshal(blockEvent)
}

func (fc *FabClient) getBlockTxNum(channelID string, blkNum uint64) (uint64, error) {
	if err := fc.clientCheck("ledgerClient", channelID); err != nil {
		return 0, err
	}

	blk, err := fc.ledgerClient.QueryBlock(blkNum)
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

func (fc *FabClient) Discover(channelID string) {
	if err := fc.clientCheck("discClient", ""); err != nil {
		return
	}

	req := discovery.NewRequest().AddLocalPeersQuery().OfChannel(channelID).AddPeersQuery()

	var targets []fab.PeerConfig

	for _, v := range fc.fabnet.Peers {
		target := fab.PeerConfig{
			URL:         v.Url,
			GRPCOptions: v.GrpcOptions,
		}
		targets = append(targets, target)
	}

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

	responses, err := fc.discClient.Send(ctx, req, targets...)
	if err != nil {
		return
	}
	//responses, err := client.Send(ctx, req)
	response := responses[0]
	chResp := response.ForChannel(channelID)
	peers, _ := chResp.Peers()
	for _, peer := range peers {
		fmt.Println(peer)
		//payload := peer.AliveMessage.Envelope.Payload
		//log.Println(string(payload))
		/*
			endPoint := peer.AliveMessage.GetAliveMsg().Membership.Endpoint
			mspid := peer.MSPID
			name := strings.Split(endPoint, ":")[0]

			flag := false

			for _, fp := range fc.fabnet.Peers {
				if fp.Url == endPoint {
					fc.fabnet.Channels[channelID].Config.Peers[name] = &Peer{
						Name:        name,
						Url:         endPoint,
						MspID:       mspid,
						GrpcOptions: fp.GrpcOptions,
						Status:      "Online",
					}
					flag = true
					break
				}
			}

			if !flag {
				fc.fabnet.Channels[channelID].Config.Peers[name] = &Peer{
					Name:  name,
					Url:   endPoint,
					MspID: mspid,
					GrpcOptions: map[string]interface{}{
						"allow-insecure": true,
					},
					Status: "Online",
				}
			}
		*/
	}
}

/*
func (fc *FabClient) ProbeInfo(channelID string) {
	fc.fabnet.Channels[channelID].Chain = fc.getChainInfo(channelID)
	fc.fabnet.Channels[channelID].Config = fc.getConfigInfo(channelID)
}
*/

func (fc *FabClient) FabNetInfo(channelID string) ([]byte, error) {
	return json.Marshal(fc.fabnet)
}

func (fc *FabClient) ChannelInfo(channelID string) ([]byte, error) {
	return nil, nil
}

func (fc *FabClient) ChainInfo(channelID string) ([]byte, error) {
	chain, err := fc.GetChainInfo(channelID)
	if err != nil && chain != nil {
		return nil, err
	}

	// TODO  TxNum

	return json.Marshal(chain)
}

func (fc *FabClient) GetChainInfo(channelID string) (*Chain, error) {
	if err := fc.clientCheck("ledgerClient", channelID); err != nil {
		//return nil, err
		return nil, err
	}

	lgc, err := fc.ledgerClient.QueryInfo(ledger.WithTimeout(fab.ChannelConfigRefresh, 2*time.Second))
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

func (fc *FabClient) ConfigInfo(channelID string) ([]byte, error) {
	config, err := fc.GetConfigInfo(channelID)
	if err != nil && config != nil {
		return nil, err
	}

	return json.Marshal(config)
}

func (fc *FabClient) GetConfigInfo(channelID string) (*Config, error) {
	if err := fc.clientCheck("ledgerClient", channelID); err != nil {
		//return nil, err
		return nil, err
	}

	cfg, err := fc.ledgerClient.QueryConfig()
	if err != nil {
		//return nil, fmt.Errorf("Get config info error: %s", err.Error())
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
	//return json.Marshal(config)
	return config, nil
}

func (fc *FabClient) GetPeerInfo(channelID string) (map[string]*Peer, error) {
	peers := make(map[string]*Peer)
	if err := fc.clientCheck("discClient", ""); err != nil {
		return nil, err
	}

	req := discovery.NewRequest().AddLocalPeersQuery().OfChannel(channelID).AddPeersQuery()

	var targets []fab.PeerConfig

	for _, v := range fc.fabnet.Peers {
		target := fab.PeerConfig{
			URL:         v.Url,
			GRPCOptions: v.GrpcOptions,
		}
		targets = append(targets, target)
	}

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

	responses, err := fc.discClient.Send(ctx, req, targets...)
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

		endPoint := pr.AliveMessage.GetAliveMsg().Membership.Endpoint
		mspid := pr.MSPID
		name := strings.Split(endPoint, ":")[0]

		peer := &Peer{
			Name:   name,
			Url:    endPoint,
			MspID:  mspid,
			Status: "Online",
		}

		peers[peer.Name] = peer
	}

	return peers, nil
}

func (fc *FabClient) ChaincodeInfo(channelID string, peer string) ([]byte, error) {
	chaincodes, err := fc.GetChaincodeInfo(channelID, peer)
	if err != nil && chaincodes != nil {
		return nil, err
	}

	return json.Marshal(chaincodes)
}

func (fc *FabClient) GetChaincodeInfo(channelID string, peer string) (map[string]*Chaincode, error) {
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

func (fc *FabClient) BlocksInfo(channelID, endstr, numstr string) ([]byte, error) {
	var end, num uint64
	num, err := strconv.ParseUint(numstr, 10, 64)
	if err != nil {
		return nil, err
	}

	var blocks []*Block
	if err := fc.clientCheck("ledgerClient", channelID); err != nil {
		return nil, err
	}

	lgc, err := fc.ledgerClient.QueryInfo(ledger.WithTimeout(fab.ChannelConfigRefresh, 2*time.Second))
	if err != nil {
		//return nil, fmt.Errorf("Get ledger info error: %s", err.Error())
		return nil, err
	}

	if endstr == "-1" {
		end = lgc.BCI.Height
	} else {
		end, err = strconv.ParseUint(endstr, 10, 64)
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
		block, err := fc.GetBlockInfo(channelID, "", strconv.FormatUint(end-1, 10))
		if err == nil && block != nil {
			blocks = append(blocks, block)
		}
		end--
	}

	return json.Marshal(blocks)
}

func (fc *FabClient) BlockInfo(channelID, by, arg string) ([]byte, error) {
	block, err := fc.GetBlockInfo(channelID, by, arg)
	if err != nil && block != nil {
		return nil, err
	}

	return json.Marshal(block)
}

func (fc *FabClient) GetBlockInfo(channelID, by, arg string) (*Block, error) {
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
		blk, err = fc.ledgerClient.QueryBlock(num)
	} else if by == "hash" {
		blk, err = fc.ledgerClient.QueryBlockByHash([]byte(arg))
	} else if by == "txhash" {
		blk, err = fc.ledgerClient.QueryBlockByTxID(fab.TransactionID(arg))
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

func (fc *FabClient) TransactionsInfo(channelID, endstr, numstr string) ([]byte, error) {
	num, err := strconv.ParseUint(numstr, 10, 64)
	if err != nil {
		return nil, err
	}

	if err := fc.clientCheck("ledgerClient", channelID); err != nil {
		return nil, err
	}

	var txs []*Transaction

	lgc, err := fc.ledgerClient.QueryInfo(ledger.WithTimeout(fab.ChannelConfigRefresh, 2*time.Second))
	if err != nil {
		//return nil, fmt.Errorf("Get ledger info error: %s", err.Error())
		return nil, err
	}

	for i := lgc.BCI.Height; i > 0; i-- {
		block, err := fc.GetBlockInfo(channelID, "", strconv.FormatUint(i-1, 10))
		if err == nil && block != nil {
			for j := len(block.Txs); j > 0; j-- {
				tx := block.Txs[j-1]
				if tx.TxID != endstr && uint64(len(txs)) < num {
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

func (fc *FabClient) TransactionInfo(channelID, arg string) ([]byte, error) {
	transaction, err := fc.GetTransactionInfo(channelID, arg)
	if err != nil && transaction != nil {
		return nil, err
	}

	return json.Marshal(transaction)
}

func (fc *FabClient) GetTransactionInfo(channelID, arg string) (*Transaction, error) {
	if err := fc.clientCheck("ledgerClient", channelID); err != nil {
		return nil, err
	}

	blk, err := fc.ledgerClient.QueryBlockByTxID(fab.TransactionID(arg))
	if err != nil {
		return nil, err
	}

	blkNum := blk.Header.Number

	ptx, err := fc.ledgerClient.QueryTransaction(fab.TransactionID(arg))
	if err != nil {
		return nil, fmt.Errorf("Get transaction info error: %s", err.Error())
	}

	tx, err := getTransaction(ptx.GetTransactionEnvelope(), blkNum)
	if err != nil {
		return nil, fmt.Errorf("Get transaction info error: %s", err.Error())
	}

	return tx, nil
}

/*
	Orgs       []string `json:"orgs"`
	Peers      []string `json:"peers"`
	Channels   []string `json:"channels"`
	Chaincodes []string `json:"chaincodes"`
	BlockNum   uint64   `json:"blocknum"`
	TxNum      uint64   `json:"txnum"`
*/

func (fc *FabClient) BlkTxInfo(channelID string) ([]byte, error) {
	var buff bytes.Buffer

	/*
		chain, err := fc.GetChainInfo(channelID)
		if err != nil {
			return nil, err
		}
		blkNum = chain.Height

		if fc.overViews[channelID] == nil {
			txNum = fc.getTxNum(channelID, uint64(0), uint64(0), blkNum)
		} else {
			txNum = fc.getTxNum(channelID, fc.overViews[channelID].TxNum, fc.overViews[channelID].BlockNum, blkNum)
		}

		fc.overViews[channelID].BlockNum = blkNum
		fc.overViews[channelID].TxNum = txNum
	*/

	if fc.overViews[channelID] == nil {
		fc.OverViewInfo(channelID)
	}

	blkNum := fc.overViews[channelID].BlockNum
	txNum := fc.overViews[channelID].TxNum

	buff.WriteString("{")
	buff.WriteString("\"blocknum\":")
	buff.WriteString("\"" + strconv.FormatUint(blkNum, 10) + "\"")
	buff.WriteString(",")
	buff.WriteString("\"txnum\":")
	buff.WriteString("\"" + strconv.FormatUint(txNum, 10) + "\"")
	buff.WriteString("}")

	return buff.Bytes(), nil
}

func (fc *FabClient) OverViewInfo(channelID string) ([]byte, error) {
	overview := new(OverView)

	for chname, _ := range fc.fabnet.Channels {
		overview.Channels = append(overview.Channels, chname)
	}

	chain, err := fc.GetChainInfo(channelID)
	if err != nil {
		return nil, err
	}
	overview.BlockNum = chain.Height

	if fc.overViews[channelID] == nil {
		overview.TxNum = fc.getTxNum(channelID, uint64(0), uint64(0), overview.BlockNum)
	} else {
		overview.TxNum = fc.getTxNum(channelID, fc.overViews[channelID].TxNum, fc.overViews[channelID].BlockNum, overview.BlockNum)
	}

	config, err := fc.GetConfigInfo(channelID)
	if err != nil {
		return nil, err
	}

	var orgs []string
	for _, org := range config.Msps {
		if !strings.HasPrefix(org, "Orderer") {
			orgs = append(orgs, org)
		}
	}
	overview.Orgs = orgs

	peers, err := fc.GetPeerInfo(channelID)
	if err != nil {
		return nil, err
	}
	overview.Peers = peers

	var speer string
	for pname, p := range peers {
		if strings.HasPrefix(p.MspID, fc.orgName) {
			speer = pname
			break
		}
	}

	if speer != "" {
		chaincodes, err := fc.GetChaincodeInfo(channelID, speer)
		if err != nil {
			return nil, err
		}

		overview.Chaincodes = chaincodes
	}

	fc.overViews[channelID] = overview

	return json.Marshal(overview)
}

// register chaincode event
func (fc *FabClient) RegisterChaincodeEvent(channelID, chaincodeID, eventFilter string) (fab.Registration, <-chan *fab.CCEvent) {
	if err := fc.clientCheck("channelClient", channelID); err != nil {
		return nil, nil
	}

	reg, notifier, err := fc.channelClient.RegisterChaincodeEvent(chaincodeID, eventFilter)
	if err != nil {
		return nil, nil
	}
	return reg, notifier
}

func eventMonitor(notifier <-chan *fab.CCEvent, eventFilter string) {
	for {
		select {
		case ccEvent := <-notifier:
			fmt.Println(ccEvent)
			//case <-time.After(time.Second * 10):
			//	return nil, fmt.Errorf("Can not get event for Filter: %s", eventFilter)
		}
	}
}

// register various event
func (fc *FabClient) EventGo(channelID string) {
	if err := fc.clientCheck("eventClient", channelID); err != nil {
		return
	}

	// fmt.Println("event client created")
	registration, eventch, err := fc.eventClient.RegisterBlockEvent() //RegisterChaincodeEvent, RegisterFilteredBlockEvent, RegisterTxStatusEvent
	if err != nil {
		fmt.Println("Event Register error: ", err)
	}
	defer fc.eventClient.Unregister(registration)

	for {
		select {
		case ev, ok := <-eventch:
			if !ok {
				fmt.Println("Event Fetch error: ", err)
			} else {
				fmt.Println("Event: ", ev)
			}
		case <-time.After(5 * time.Second):
			fmt.Println("timed out waiting for block event")
		}
	}
}

func (fc *FabClient) BlockMonitor(channelID string, ctx context.Context, conn interface{}, handler func(interface{}, []byte, error)) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("block monitor error")
		}
	}()

	if err := fc.clientCheck("eventClient", channelID); err != nil {
		return
	}

	// fmt.Println("event client created")
	registration, eventch, err := fc.eventClient.RegisterBlockEvent() //RegisterChaincodeEvent, RegisterFilteredBlockEvent, RegisterTxStatusEvent
	if err != nil {
		fmt.Println("Event Register error: ", err)
		return
	}
	defer fc.eventClient.Unregister(registration)

	fmt.Println("Start monitor block event...")

	for {
		if fc.IsClosed() {
			return
		}
		select {
		case ev, ok := <-eventch:
			if ok {
				//fmt.Println("Event: ", ev)
				if fc.overViews[channelID] == nil {
					fc.OverViewInfo(channelID)
				}
				fc.overViews[channelID].BlockNum += 1
				fc.overViews[channelID].TxNum += uint64(len(ev.Block.Data.Data))
				eventBytes, err := parseBlockEvent(ev, fc.fabnet.Peers)
				handler(conn, eventBytes, err)
			}
		case <-ctx.Done():
			fmt.Println("Monitor stopped")
			return
		}
	}
}

/*
func (fc *FabClient) ChaincodeMonitor(channelID string, ctx context.Context, conn interface{}, handler func(interface{}, []byte, error)) {
	if err := fc.clientCheck("eventClient", channelID); err != nil {
		return
	}

	// fmt.Println("event client created")
	registration, eventch, err := fc.eventClient.RegisterChaincodeEvent(ccID, eventFilter) //RegisterChaincodeEvent, RegisterFilteredBlockEvent, RegisterTxStatusEvent
	if err != nil {
		fmt.Println("Event Register error: ", err)
		return
	}
	defer fc.eventClient.Unregister(registration)

	fmt.Println("Start monitor block event...")

	for {
		if fc.IsClosed() {
			return
		}
		select {
		case ev, ok := <-eventch:
			if ok {
				//fmt.Println("Event: ", ev)
				eventBytes, err := parseBlockEvent(ev)
				handler(conn, eventBytes, err)
			}
		case <-ctx.Done():
			fmt.Println("Monitor stopped")
			return
		}
	}
}
*/

// query channel
func (fc *FabClient) QueryChannel(peer string) ([]string, error) {
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

// query config
func (fc *FabClient) QueryConfig(key string) (map[string]interface{}, error) {
	if err := fc.clientCheck("", ""); err != nil {
		return nil, err
	}

	config, err := fc.fabricSDK.Config()
	if err != nil {
		return nil, err
	}

	if attr, ok := config.Lookup(key); ok {
		return attr.(map[string]interface{}), nil
		/*
			for k, v := range attr.(map[string]interface{}) {
				fmt.Println(k)
				fmt.Println(v)
			}
		*/
	}
	return nil, fmt.Errorf("Key not found")
}

func (fc *FabClient) getChannels() map[string]*CChannel {
	channels := make(map[string]*CChannel)
	channelMap, err := fc.QueryConfig("channels")
	if err != nil {
		return channels
	}

	for k, v := range channelMap {
		channel := new(CChannel)
		channel.Name = k
		attrs := v.(map[string]interface{})
		if v, ok := attrs["peers"]; ok && v != nil {
			channel.PeerPerms = v.(map[string]interface{})
		}
		if v, ok := attrs["policies"]; ok && v != nil {
			channel.Policies = v.(map[string]interface{})
		}
		channels[channel.Name] = channel
	}
	return channels
}

func (fc *FabClient) getOrganizations() map[string]*COrganization {
	orgs := make(map[string]*COrganization)
	orgMap, err := fc.QueryConfig("organizations")
	if err != nil {
		return orgs
	}

	for k, v := range orgMap {
		org := new(COrganization)
		org.Name = k
		attrs := v.(map[string]interface{})
		if v, ok := attrs["mspid"]; ok && v != nil {
			org.MspID = v.(string)
		}
		if v, ok := attrs["cryptopath"]; ok && v != nil {
			org.Crypto = v.(string)
		}
		if v, ok := attrs["peers"]; ok && v != nil {
			org.Nodes = v.([]interface{})
		}
		if v, ok := attrs["orderers"]; ok && v != nil {
			org.Nodes = v.([]interface{})
		}
		if v, ok := attrs["certificateauthorities"]; ok && v != nil {
			org.CertAuthors = v.([]interface{})
		}
		if v, ok := attrs["users"]; ok && v != nil {
			org.Users = v.([]interface{})
		}
		orgs[org.Name] = org
	}
	return orgs
}

func (fc *FabClient) getPeers() map[string]*CPeer {
	peers := make(map[string]*CPeer)
	peerMap, err := fc.QueryConfig("peers")
	if err != nil {
		return peers
	}

	for k, v := range peerMap {
		peer := new(CPeer)
		peer.Name = k
		attrs := v.(map[string]interface{})
		if v, ok := attrs["url"]; ok && v != nil {
			peer.Url = v.(string)
		}
		if v, ok := attrs["grpcoptions"]; ok && v != nil {
			peer.GrpcOptions = v.(map[string]interface{})
		}
		if v, ok := attrs["tlscacerts"]; ok && v != nil {
			peer.TlsCaCerts = v.(map[string]interface{})
		}
		peers[peer.Name] = peer
	}
	return peers
}

func (fc *FabClient) getOrderers() map[string]*COrderer {
	orderers := make(map[string]*COrderer)
	ordererMap, err := fc.QueryConfig("orderers")
	if err != nil {
		return orderers
	}

	for k, v := range ordererMap {
		orderer := new(COrderer)
		orderer.Name = k
		attrs := v.(map[string]interface{})
		if v, ok := attrs["url"]; ok && v != nil {
			orderer.Url = v.(string)
		}
		if v, ok := attrs["grpcoptions"]; ok && v != nil {
			orderer.GrpcOptions = v.(map[string]interface{})
		}
		if v, ok := attrs["tlscacerts"]; ok && v != nil {
			orderer.TlsCaCerts = v.(map[string]interface{})
		}
		orderers[orderer.Name] = orderer
	}
	return orderers
}

func (fc *FabClient) getFabNetInfo() *FabNet {
	configFile := new(FabNet)
	configFile.Channels = fc.getChannels()
	configFile.Organizations = fc.getOrganizations()
	configFile.Peers = fc.getPeers()
	//return json.Marshal(configFile)
	return configFile
}
