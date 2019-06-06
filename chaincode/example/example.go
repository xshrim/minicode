/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Init")
	_, args := stub.GetFunctionAndParameters()
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var err error

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	B = args[2]
	Bval, err = strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("example Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "info" {
		return t.info(stub)
	} else if function == "invoke" {
		// Make payment of X units from A to B
		return t.invoke(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

func (t *SimpleChaincode) info(stub shim.ChaincodeStubInterface) pb.Response {
	var buffer bytes.Buffer

	timestamp, _ := stub.GetTxTimestamp()

	timestring := time.Unix(timestamp.Seconds, int64(timestamp.Nanos)).String()

	creatorBytes, err := stub.GetCreator()
	if err != nil {
		creatorBytes = []byte("")
	}

	// 获取证书部分
	certStart := bytes.IndexAny(creatorBytes, "-----BEGIN")
	if certStart == -1 {
		return shim.Error("No certificate found")
	}

	// 使用正则表达式过滤掉creatorBytes中mspid部分的不可见字符(protobuf编码)
	msp := regexp.MustCompile(`\w+`).FindString(string(creatorBytes[:certStart]))

	// 直接从creatorBytes获取mspid
	//msp := creatorBytes[2 : certStart-3] // mspid两端存在不可见字符

	certText := creatorBytes[certStart:]

	// 获取证书的pem格式编码
	bl, _ := pem.Decode(certText)
	if bl == nil {
		return shim.Error("Could not decode the PEM structure")
	}

	// 生成证书对象
	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return shim.Error("ParseCertificate failed")
	}

	// 获取证书中包含的各种数据
	uname := cert.Subject.CommonName // 用户名, 如Admin@org1.example.com
	fmt.Println("Name:" + uname)

	// 直接解析证书的话, 就不需要使用cid包了(即可以无需引入外部包)
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

	buffer.WriteString("\"timestamp\":")
	buffer.WriteString("\"")
	buffer.WriteString(timestring)
	buffer.WriteString("\"")
	buffer.WriteString(",")

	buffer.WriteString("\"creator\":")
	buffer.WriteString("\"")
	buffer.WriteString(uname)
	buffer.WriteString("\"")
	buffer.WriteString(",")

	buffer.WriteString("\"id\":")
	buffer.WriteString("\"")
	buffer.WriteString(id)
	buffer.WriteString("\"")
	buffer.WriteString(",")

	buffer.WriteString("\"msp\":")
	buffer.WriteString("\"")
	buffer.WriteString(msp)
	buffer.WriteString("\"")
	buffer.WriteString(",")

	buffer.WriteString("\"mspid\":")
	buffer.WriteString("\"")
	buffer.WriteString(mspid)
	buffer.WriteString("\"")

	buffer.WriteString("}")

	return shim.Success(buffer.Bytes())
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var X int          // Transaction value
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	A = args[0]
	B = args[1]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Bvalbytes == nil {
		return shim.Error("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	Aval = Aval - X
	Bval = Bval + X
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

// query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
