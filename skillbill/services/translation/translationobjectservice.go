package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/msp"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// TranslationObjectChaincode define the Smart Contract structure
type TranslationObjectChaincode struct {
}

// Init method is called when the Smart Contract "translationobject" is instantiated by the blockchain network
func (s *TranslationObjectChaincode) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

// Invoke method is called as a result of an application request to run the Smart Contract "translationobject"
func (s *TranslationObjectChaincode) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()

	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "getAll" {
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

func (s *TranslationObjectChaincode) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	// GetCreator returns marshaled serialized identity of the client
	serializedID, _ := APIstub.GetCreator()
	sID := &msp.SerializedIdentity{}
	err := proto.Unmarshal(serializedID, sID)
	if err != nil {
		return shim.Error(fmt.Sprintf("Could not deserialize a SerializedIdentity, err %s", err))
	}

	bl, _ := pem.Decode(sID.IdBytes)
	if bl == nil {
		return shim.Error(fmt.Sprintf("Failed to decode PEM structure"))
	}

	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Unable to parse certificate %s", err))
	}
	rsaPublicKey := cert.PublicKey.(*rsa.PublicKey)
	return shim.Error(fmt.Sprintf("Public Key %s ", rsaPublicKey.N))

	translations := []TranslationObject{
		TranslationObject{LanguageID: "en", Translation: "the hyperledger blockchain", TranslationObjectID: "1"},
		TranslationObject{LanguageID: "en", Translation: "the go language", TranslationObjectID: "2"}}

	i := 0
	for i < len(translations) {
		fmt.Println("i is ", i)
		translations[i].DocType = "Translation"
		translationAsBytes, _ := json.Marshal(translations[i])
		APIstub.PutState("TRANS"+strconv.Itoa(i), translationAsBytes)
		fmt.Println("Added", translations[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *TranslationObjectChaincode) create(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	var translation = TranslationObject{LanguageID: args[1], Translation: args[2], DocType: args[3]}

	translationAsBytes, _ := json.Marshal(translation)
	APIstub.PutState(args[0], translationAsBytes)

	return shim.Success(nil)
}

func (s *TranslationObjectChaincode) getAllByQuery(APIstub shim.ChaincodeStubInterface) sc.Response {

	resultsIterator, err := APIstub.GetQueryResult(`{"selector":{"DocType":"translation"}}`)
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

	fmt.Printf("- GetAllTranslationsByQuery:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *TranslationObjectChaincode) getAll(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "TRANS0"
	endKey := "TRANS999"

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

	fmt.Printf("- GetAllTranslations:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *TranslationObjectChaincode) getByID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	translationAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(translationAsBytes)
}

func (s *TranslationObjectChaincode) delete(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	APIstub.DelState(args[0])
	return shim.Success(nil)
}

func (s *TranslationObjectChaincode) update(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	data := &TranslationObject{}

	translationAsBytes, _ := APIstub.GetState(args[0])
	err := json.Unmarshal(translationAsBytes, data)
	if err != nil {
		fmt.Printf("Error creating new TranslationObject Chaincode: %s", err)
	}
	data.LanguageID = args[1]
	data.Translation = args[2]

	translation2AsBytes, _ := json.Marshal(data)
	APIstub.PutState(args[0], translation2AsBytes)

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(TranslationObjectChaincode))
	if err != nil {
		fmt.Printf("Error creating new TranslationObject Chaincode: %s", err)
	}
}
