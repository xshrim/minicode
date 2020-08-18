package fabcore

import "crypto/x509"

type CA struct {
	CaName                    string `json:"caname"`
	CaChain                   string `json:"cachain"`
	Version                   string `json:"version"`
	Affiliations              string `json:"affiliations"`
	Identities                string `json:"identities"`
	IssuerPublicKey           string `json:"issuerpublickey"`
	IssuerRevocationPublicKey string `json:"issuerrevocationpublickey"`
}
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
}

type Config struct {
	ChannelID       string   `json:"channelid"`
	AnchorPeers     []string `json:"anchorpeers"`
	LastConfigBlock uint64   `json:"lastconfigblock"`
	Msps            []string `json:"msps"`
	Orderers        []string `json:"orderers"`
	Consensus       string   `json:"consensus"`
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

type Peer struct {
	Name   string            `json:"name"`
	Url    string            `json:"url"`
	Cert   *x509.Certificate `json:"cert"`
	MspID  string            `json:"mspid"`
	Status string            `json:"status"`
}
