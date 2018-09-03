package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/common/tools/protolator"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/msp"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type RightService struct {
}

func (t *RightService) addRight(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of columns. Expecting 2")
	}

	var right = Right{RightID: args[0], RightName: args[1], DocType: RightTableName}

	bytes, _ := json.Marshal(right)
	APIstub.PutState(RightTableName+args[0], bytes)

	return shim.Success(nil)
}

func (t *RightService) getAllRights(APIstub shim.ChaincodeStubInterface) pb.Response {

	query := `{"selector":{"doctype":"` + RightTableName + `"}}`
	resultsIterator, err := APIstub.GetQueryResult(query)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getAllRights:\n%s\n", buffer.String())

	creator := t.getCreator(APIstub)
	fmt.Printf("- Creator:\n%s\n", creator.String())

	return shim.Success(buffer.Bytes())
}

// ============================================================================================================================
// Init - initialize the chaincode
// ============================================================================================================================

func (t *RightService) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Main
func main() {
	err := shim.Start(new(RightService))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================

func (t *RightService) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println(" ")
	fmt.Println("starting invoke, for - " + function)
	if function == "addRight" {
		return t.addRight(stub, args)
	} else if function == "getAllRights" {
		return t.getAllRights(stub)
	}

	return shim.Error("Function with the name " + function + " does not exist.")
}

//===================================================================================================
// functon getCreator
// You can verify by calling getCreator during initMarble and checking fot the value
// during a transferMarble say
//===================================================================================================

func (t *RightService) getCreator(stub shim.ChaincodeStubInterface) pb.Response {

	fmt.Printf("\nBegin*** getCreator \n")
	creator, err := stub.GetCreator()
	if err != nil {
		fmt.Printf("GetCreator Error")
		return shim.Error(err.Error())
	}

	si := &msp.SerializedIdentity{}
	err2 := proto.Unmarshal(creator, si)
	if err2 != nil {
		fmt.Printf("Proto Unmarshal Error")
		return shim.Error(err2.Error())
	}
	buf := &bytes.Buffer{}
	protolator.DeepMarshalJSON(buf, si)
	fmt.Printf("End*** getCreator \n")
	fmt.Printf(string(buf.Bytes()))

	return shim.Success([]byte(buf.Bytes()))
}
