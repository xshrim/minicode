package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	// "github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	"github.com/hyperledger/fabric/protos/peer"
)

type PaperContract struct {
	PaperPrefix string
	TxPrefix    string
}

type CommercialPaper struct {
	Issuer       string `json:"issuer"`
	Paper        string `json:"paper"`
	Owner        string `json:"owner"`
	IssueDate    string `json:"idate"`
	MaturityDate string `json:"mdate"`
	FaceValue    string `json:"value"`
	State        string `json:"state"`
}

type Transaction struct {
}

func (s *PaperContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	s.PaperPrefix = "org.papernet.paper"
	s.TxPrefix = "org.papernet.tx"
	return shim.Success(nil)
}

func (s *PaperContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()
	switch function {
	case "initPaper":
		return s.initPaper(stub)
	case "invokerInfo":
		return s.invokerInfo(stub)
	case "queryPaper":
		return s.queryPaper(stub, args)
	case "queryPaperByRange":
		return s.queryPaperByRange(stub, args)
	case "queryPaperByCompositeKey":
		return s.queryPaperByCompositeKey(stub, args)
	case "issuePaper":
		return s.issuePaper(stub, args)
	case "tradePaper":
		return s.tradePaper(stub, args)
	case "redeemPaper":
		return s.redeemPaper(stub, args)
	default:
		return shim.Error("Invalid Smart Contract Function Name")
	}
}

func (s *PaperContract) initPaper(stub shim.ChaincodeStubInterface) peer.Response {
	papers := []CommercialPaper{
		CommercialPaper{Issuer: "Org1MSP", Paper: "00001", Owner: "Org1MSP", IssueDate: "2019-01-01", MaturityDate: "2019-02-01", FaceValue: "10000", State: "issued"},
		CommercialPaper{Issuer: "Org1MSP", Paper: "00002", Owner: "Org1MSP", IssueDate: "2019-02-01", MaturityDate: "2019-03-01", FaceValue: "20000", State: "issued"},
		CommercialPaper{Issuer: "Org1MSP", Paper: "00003", Owner: "Org1MSP", IssueDate: "2019-03-01", MaturityDate: "2019-04-01", FaceValue: "30000", State: "issued"},
		CommercialPaper{Issuer: "Org1MSP", Paper: "00004", Owner: "Org1MSP", IssueDate: "2019-04-01", MaturityDate: "2019-05-01", FaceValue: "40000", State: "issued"},
		CommercialPaper{Issuer: "Org1MSP", Paper: "00005", Owner: "Org1MSP", IssueDate: "2019-05-01", MaturityDate: "2019-06-01", FaceValue: "50000", State: "issued"},
		CommercialPaper{Issuer: "Org1MSP", Paper: "00006", Owner: "Org1MSP", IssueDate: "2019-06-01", MaturityDate: "2019-07-01", FaceValue: "60000", State: "issued"},
		CommercialPaper{Issuer: "Org1MSP", Paper: "00007", Owner: "Org1MSP", IssueDate: "2019-07-01", MaturityDate: "2019-08-01", FaceValue: "70000", State: "issued"},
		CommercialPaper{Issuer: "Org1MSP", Paper: "00008", Owner: "Org1MSP", IssueDate: "2019-08-01", MaturityDate: "2019-09-01", FaceValue: "80000", State: "issued"},
		CommercialPaper{Issuer: "Org1MSP", Paper: "00009", Owner: "Org1MSP", IssueDate: "2019-09-01", MaturityDate: "2019-10-01", FaceValue: "90000", State: "issued"},
	}

	for i := 0; i < len(papers); i++ {
		paper := papers[i]

		ikey, err := stub.CreateCompositeKey("prefix~issuer~paper", []string{s.PaperPrefix, paper.Issuer, paper.Paper})
		if err != nil {
			continue
		}

		paperBytes, err := json.Marshal(paper)
		if err != nil {
			continue
		}

		stub.PutState(ikey, paperBytes)
	}

	return shim.Success(nil)
}

func (s *PaperContract) invokerInfo(stub shim.ChaincodeStubInterface) peer.Response {
	var buffer bytes.Buffer

	creator, err := stub.GetCreator()
	if err != nil {
		creator = []byte("")
	}

	id, err := cid.GetID(stub)
	if err != nil {
		id = ""
	}

	mspid, err := cid.GetMSPID(stub)
	if err != nil {
		mspid = ""
	}

	// cert, err := cid.GetX509Certificate(stub)

	buffer.WriteString("{")

	buffer.WriteString("\"creator\":")
	buffer.WriteString("\"")
	buffer.WriteString(string(creator))
	buffer.WriteString("\"")
	buffer.WriteString(",")

	buffer.WriteString("\"id\":")
	buffer.WriteString("\"")
	buffer.WriteString(id)
	buffer.WriteString("\"")
	buffer.WriteString(",")

	buffer.WriteString("\"mspid\":")
	buffer.WriteString("\"")
	buffer.WriteString(mspid)
	buffer.WriteString("\"")

	buffer.WriteString("}")

	return shim.Success(buffer.Bytes())
}

func (s *PaperContract) queryPaper(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	ikey, err := stub.CreateCompositeKey("prefix~issuer~paper", []string{s.PaperPrefix, args[0], args[1]})
	if err != nil {
		return shim.Error("Failed to create composite key")
	}

	paperBytes, err := stub.GetState(ikey)
	if err != nil {
		return shim.Error("Can not find satisfied data")
	}

	return shim.Success(paperBytes)
}

// not work, range query collisions with composite keys, range query noly work with simple key
func (s *PaperContract) queryPaperByRange(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	startkey, err := stub.CreateCompositeKey("prefix~issuer~paper", []string{s.PaperPrefix, args[0], args[1]})
	if err != nil {
		return shim.Error("Failed to create composite key")
	}
	endkey, err := stub.CreateCompositeKey("prefix~issuer~paper", []string{s.PaperPrefix, args[0], args[2]})
	if err != nil {
		return shim.Error("Failed to create composite key")
	}

	resultIterator, err := stub.GetStateByRange(startkey, endkey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

func (s *PaperContract) queryPaperByCompositeKey(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	resultIterator, err := stub.GetStateByPartialCompositeKey("prefix~issuer~paper", []string{s.PaperPrefix, args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultIterator.HasNext() {
		queryResponse, err := resultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

func (s *PaperContract) issuePaper(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	paper := CommercialPaper{Issuer: args[0], Paper: args[1], Owner: args[0], IssueDate: args[2], MaturityDate: args[3], FaceValue: args[4], State: "issued"}
	ikey, err := stub.CreateCompositeKey("prefix~issuer~paper", []string{s.PaperPrefix, paper.Issuer, paper.Paper})
	if err != nil {
		return shim.Error("Failed to create composite key")
	}

	paperBytes, _ := json.Marshal(paper)
	stub.PutState(ikey, paperBytes)

	return shim.Success([]byte(ikey))
}

func (s *PaperContract) tradePaper(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	ikey, err := stub.CreateCompositeKey("prefix~issuer~paper", []string{s.PaperPrefix, args[0], args[1]})
	if err != nil {
		return shim.Error("Failed to create composite key")
	}

	paperBytes, err := stub.GetState(ikey)
	if err != nil {
		return shim.Error("Can not find satisfied data")
	}

	paper := CommercialPaper{}
	json.Unmarshal(paperBytes, &paper)

	if paper.State == "redeemed" {
		return shim.Error("Paper " + paper.Issuer + paper.Paper + " already redeemed")
	}

	if paper.State == "issued" {
		paper.State = "trading"
	}

	if paper.State == "trading" {
		paper.Owner = args[2]
		paperBytes, err = json.Marshal(paper)
		if err != nil {
			return shim.Error("Failed to create composite key")
		}

		stub.PutState(ikey, paperBytes)
	}

	return shim.Success([]byte(ikey))
}

func (s *PaperContract) redeemPaper(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	ikey, err := stub.CreateCompositeKey("prefix~issuer~paper", []string{s.PaperPrefix, args[0], args[1]})
	if err != nil {
		return shim.Error("Failed to create composite key")
	}

	paperBytes, err := stub.GetState(ikey)
	if err != nil {
		return shim.Error("Can not find satisfied data")
	}

	paper := CommercialPaper{}
	json.Unmarshal(paperBytes, &paper)

	if paper.State == "redeemed" {
		return shim.Error("Paper " + paper.Issuer + paper.Paper + " already redeemed")
	}

	if paper.State == "trading" {
		paper.Owner = paper.Issuer
		paper.State = "redeemed"

		paperBytes, err = json.Marshal(paper)
		if err != nil {
			return shim.Error("Failed to create composite key")
		}

		stub.PutState(ikey, paperBytes)
	}

	return shim.Success([]byte(ikey))
}

func main() {
	err := shim.Start(new(PaperContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
