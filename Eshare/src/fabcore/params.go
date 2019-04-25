package fabcore

type Chain struct {
	CurBlockHash string `json:"curblockhash"`
	PreBlockHash string `json:"preblockhash"`
	Height       uint64 `json:"height"`
	Endorser     string `json:"endorser"`
	Status       string `json:"status"`
}

type Asn1Head struct {
	Number       int64
	PreviousHash []byte
	DataHash     []byte
}

type Block struct {
	Number       uint64         `json:"number"`
	Hash         string         `json:"hash"`
	Size         int            `json:"size"`
	PreviousHash string         `json:"previoushash"`
	DataHash     string         `json:"datahash"`
	Txs          []*Transaction `json:"txs"`
	Timestamp    string         `json:"timestamp"`
}

type Transaction struct {
	TxID          string   `json:"txid"`
	BlkNum        uint64   `json:"blknum"`
	TxType        string   `json:"txtype"`
	Size          int      `json:"size"`
	MspID         string   `json:"mspid"`
	Endorsers     []string `json:"endorsers"`
	Args          []string `json:"args"`
	ChannelID     string   `json:"channelid"`
	ChaincodeID   string   `json:"chaincodeid"`
	ChaincodeLang string   `json:"chaincodelang"`
	Timestamp     string   `json:"timestamp"`
}

type BlockEvent struct {
	Number       uint64         `json:"number"`
	Hash         string         `json:"hash"`
	Size         int            `json:"size"`
	PreviousHash string         `json:"previoushash"`
	DataHash     string         `json:"datahash"`
	Txs          []*Transaction `json:"txs"`
	Timestamp    string         `json:"timestamp"`
	SourceURL    string         `json:"sourceurl"`
	SourceAlias  string         `json:"sourcealias"`
}

type Config struct {
	ChannelID       string   `json:"channelid"`
	AnchorPeers     []string `json:"anchorpeers"`
	LastConfigBlock uint64   `json:"lastconfigblock"`
	Msps            []string `json:"msps"`
	Orderers        []string `json:"orderers"`
	Consensus       string   `json:"consensus"`
}

type CChannel struct {
	Name      string                 `json:"name"`
	PeerPerms map[string]interface{} `json:"peerperms"`
	Policies  map[string]interface{} `json:"polices"`
}

type Chaincode struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Path    string `json:"path"`
	Version string `json:"version"`
	Input   string `json:"input"`
	Escc    string `json:"escc"`
	Vscc    string `json:"vscc"`
	Status  string `json:"status"`
}

type COrganization struct {
	Name        string        `json:"name"`
	MspID       string        `json:"mspid"`
	Crypto      string        `json:"crypto"`
	Nodes       []interface{} `json:"nodes"`
	CertAuthors []interface{} `json:"certauthors"`
	Users       []interface{} `json:"users"`
}

type CPeer struct {
	Name        string                 `json:"name"`
	Url         string                 `json:"url"`
	GrpcOptions map[string]interface{} `json:"grpcoptions"`
	TlsCaCerts  map[string]interface{} `json:"tlscacerts"`
}

type COrderer struct {
	Name        string                 `json:"name"`
	Url         string                 `json:"url"`
	GrpcOptions map[string]interface{} `json:"grpcoptions"`
	TlsCaCerts  map[string]interface{} `json:"tlscacerts"`
}

type Peer struct {
	Name   string `json:"name"`
	Url    string `json:"url"`
	MspID  string `json:"mspid"`
	Status string `json:"status"`
}

type OverView struct {
	Orgs       []string              `json:"orgs"`
	Peers      map[string]*Peer      `json:"peers"`
	Channels   []string              `json:"channels"`
	Chaincodes map[string]*Chaincode `json:"chaincodes"`
	BlockNum   uint64                `json:"blocknum"`
	TxNum      uint64                `json:"txnum"`
}

type FabNet struct {
	Channels      map[string]*CChannel      `json:"channels"`
	Organizations map[string]*COrganization `json:"organizations"`
	Peers         map[string]*CPeer         `json:"peers"`
	Orderers      map[string]*COrderer      `json:"orderers"`
}

