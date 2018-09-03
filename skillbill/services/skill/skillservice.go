package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// SkillChaincode define the Smart Contract structure
type SkillChaincode struct {
}

// Init method is called when the Smart Contract "Skill" is instantiated by the blockchain network
func (s *SkillChaincode) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

// Invoke method is called as a result of an application request to run the Smart Contract "Skill"
func (s *SkillChaincode) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()

	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "getAll" {
		return s.getAll(APIstub)
	} else if function == "getAllByQuery" {
		return s.getAllByQuery(APIstub)
	} else if function == "create" {
		return s.create(APIstub, args)
	} else if function == "getByID" {
		return s.getByID(APIstub, args)
	} else if function == "delete" {
		return s.delete(APIstub, args)
	} else if function == "update" {
		return s.update(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SkillChaincode) create(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// level, errlevel := strconv.Atoi(args[5])
	// if errlevel != nil {
	// 	fmt.Println(errlevel)
	// }

	// est, errest := strconv.ParseFloat(args[8], 64)
	// if errest != nil {
	// 	fmt.Println(errest)
	// }

	// ver, errver := strconv.Atoi(args[9])
	// if errver != nil {
	// 	fmt.Println(errver)
	// }

	var skill = Skill{BackwardCompatibleTo: args[1],
		DescriptionTranslationID: args[2],
		ImageID:                  args[3],
		KnowledgeGroupID:         args[4],
		//Level:                    level,
		NameTranslationID: args[6],
		SkillID:           args[7],
		//TimeEstimationInHours:    est,
		//Version:                  ver
	}

	skillAsBytes, _ := json.Marshal(skill)
	APIstub.PutState(args[0], skillAsBytes)

	return shim.Success(nil)
}

func (s *SkillChaincode) getAllByQuery(APIstub shim.ChaincodeStubInterface) sc.Response {

	resultsIterator, err := APIstub.GetQueryResult(`{"selector":{"DocType":"skill"}}`)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing
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

	fmt.Printf("- getAllByQuery:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SkillChaincode) getAll(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "SKILL0"
	endKey := "SKILL999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing
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

	fmt.Printf("- GetAllskills:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SkillChaincode) getByID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	skillAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(skillAsBytes)
}

func (s *SkillChaincode) delete(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	APIstub.DelState(args[0])
	return shim.Success(nil)
}

func (s *SkillChaincode) update(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	data := &Skill{}

	skillAsBytes, _ := APIstub.GetState(args[0])
	err := json.Unmarshal(skillAsBytes, data)
	if err != nil {
		fmt.Printf("Error creating new Skill Chaincode: %s", err)
	}
	//data.Level, err = strconv.Atoi(args[1])
	if err != nil {
		fmt.Printf("Error creating new Skill Chaincode: %s", err)
	}
	skill2AsBytes, _ := json.Marshal(data)
	APIstub.PutState(args[0], skill2AsBytes)

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(SkillChaincode))
	if err != nil {
		fmt.Printf("Error creating new Skill Chaincode: %s", err)
	}
}
