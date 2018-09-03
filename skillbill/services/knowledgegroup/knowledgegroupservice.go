package main


import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc 	"github.com/hyperledger/fabric/protos/peer"
	logs "github.com/skillbill/packages/logs"
)

// RoleChaincode define the Smart Contract structure
type KnowledgeGroupChaincode struct {
	repo 	KnowledgeGroupRepo
}

// Init method is called when the Smart Contract "Feature" is instantiated by the blockchain network
func (s *KnowledgeGroupChaincode) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	logs.SetUpLogging("var/log/knowledge.log")
	InitKnowledgeGrpRepo()

	return shim.Success(nil)
}

// Invoke method is called as a result of an application request to run the Smart Contract "Feature"
func (s *KnowledgeGroupChaincode) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()

	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "GetByQuery" {
		return s.GetByQuery(APIstub, args)
	} else if function == "CreateGroup" {
		return s.CreateGroup(APIstub, args)
	} else if function == "UpdateGroup" {
		return s.UpdateGroup(APIstub, args)
	} else if function == "Delete" {
		return s.DeleteKnowledgeGrpOrGrpMember(APIstub, args)
	} else if function == "AddMembersToGroup" {
		return s.AddMembersToGroup(APIstub, args)
	} else if function == "GetMemberByGroupID" {
		return s.GetMemberByGroupID(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

// args[0].. args[n] are pair column and value
// e.g: args['DocType,knowledgegroup', 'GroupID,E1ED5DAD-B286-4522-8A93-926E6D5DC9C9', ....]
func (s *KnowledgeGroupChaincode) GetByQuery(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	logs.LogInfo("Before call GetByQuery in the knowledgegrp repo.")
	result, err := s.repo.GetByQuery(APIstub, args)

	if err != nil {
		return shim.Error("Failed to query knowledge group due to " + err.Error())
	}

	return shim.Success(result)
}

func (s *KnowledgeGroupChaincode) CreateGroup(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	id, err := s.repo.CreateKnowledgeGrp(APIstub, args[0])

	if err != nil {
		return shim.Error("Failed to create knowledge group " + args[0] + " due to " + err.Error())
	}

	return shim.Success([]byte(id))
}

func (s *KnowledgeGroupChaincode) DeleteKnowledgeGrpOrGrpMember(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var id = args[0]
	value, err := s.repo.GetByKey(APIstub, id)
	
	if err != nil {
		return shim.Error("Failed to delete record " + id + " due to " + err.Error())
	}

	if string(value) == "" {
		return shim.Error("Failed to delete record " + id + ", beacause it does not exist.")
	}

	err = s.repo.DeleteRecord(APIstub, id)

	if err != nil{
		return shim.Error("Failed to delete record " + id + " due to " + err.Error())
	}

	return shim.Success(nil)
}

func (s *KnowledgeGroupChaincode) UpdateGroup(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	err := s.repo.UpdateKnowledgeGrp(APIstub, args)

	if err != nil {
		return shim.Error("Failed to update the knowledge group due to")
	}

	return shim.Success(nil)
}

// args[0] is group id, args[1] is member type, args[2] is user id (member)
// e.x: ['groupId', 'Professional','userId']
func (s *KnowledgeGroupChaincode) AddMembersToGroup(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	memberType, err := getMemberType(args[1])

	if err != nil {
		return shim.Error("Failed to add member " + args[2] + "due to " + err.Error())
	}

	id, err := s.repo.AddMembersToKnowledgeGrp(APIstub, args[0], memberType, args[1])

	return shim.Success([]byte(id))
}

func (s *KnowledgeGroupChaincode) GetMemberByGroupID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	data, err := s.repo.GetMembersByGroupID(APIstub, args[0])

	if err != nil {
		return shim.Error("Failed to get list members by group id "+ args[0] + " due to " + err.Error())
	}

	return shim.Success(data)
}

const (
	Professional	string = "Professional"
	Assessor		string = "Assessor"
)

func getMemberType(input string) (string, error){

	switch input {
	case Professional:
		return Professional, nil

	case Assessor:
		return Assessor, nil

	default:
		return "", fmt.Errorf("Invalid member type " + input)
	}
}

func main() {
	err := shim.Start(new(KnowledgeGroupChaincode))
	if err != nil {
		fmt.Printf("Error creating new knowledge group Chaincode: %s", err)
	}
}
