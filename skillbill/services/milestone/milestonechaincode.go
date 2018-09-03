package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	log "github.com/skillbill/packages/logs"
	"github.com/skillbill/packages/repository"
)

const MilestoneDocType string = "milestone"
const MilestoneDependencyDocType string = "milestonedependency"
const MilestoneSkillDocType string = "milestoneskill"

var milestoneRepo repository.IRepo
var milestoneDependencyRepo repository.IRepo
var milestoneSkillRepo repository.IRepo

type MilestoneChaincode struct {
}

// Init method is called when the Smart Contract "Feature" is instantiated by the blockchain network
func (m *MilestoneChaincode) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	log.SetUpLogging("var/log/milestone.log")

	milestoneRepo = repository.InitRepo(MilestoneDocType)
	milestoneDependencyRepo = repository.InitRepo(MilestoneDependencyDocType)
	// milestoneSkillRepo = repository.InitRepo(MilestoneSkillDocType)

	return shim.Success(nil)
}

// Invoke method is called as a result of an application request to run the Smart Contract "Feature"
func (m *MilestoneChaincode) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()

	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "GetAllByQuery" {
		return m.GetAllMilestones(APIstub)
	} else if function == "CreateMilestone" {
		return m.CreateMilestone(APIstub, args)
	} else if function == "GetMilestoneByID" {
		return m.GetMilestoneByID(APIstub, args[0])
	} else if function == "UpdateMilestone" {
		return m.UpdateMilestone(APIstub, args)
	} else if function == "DeleteRecord" {
		return m.DeleteRecord(APIstub, args[0])
	} else if function == "CreateMilestoneDependency" {
		return m.CreateMilestoneDependency(APIstub, args)
	} else if function == "GetDependingsByID" {
		return m.GetDependingsByID(APIstub, args[0])
	} else if function == "UpdateMilestoneDependency" {
		return m.UpdateMilestoneDependency(APIstub, args)
	}

	// switch strings.ToUpper(function) {
	// case "ADDMILESTONEDEPENDENCY":
	// 	return m.CreateMilestoneDependency(APIstub, args)
	// case "GetDependingsByID":
	// 	return m.GetDependingsByID(APIstub, args[0])
	// case "UpdateMilestoneDependency":
	// 	return m.UpdateMilestoneDependency(APIstub, args)
	// default:
	// 	return shim.Error("Invalid Smart Contract function name: " + function)
	// }

	return shim.Error("Invalid Smart Contract function name: " + function + "sadsadsad")
}

// func (m *MilestoneChaincode) GetAllMilestones(APIstub shim.ChaincodeStubInterface) sc.Response {

// 	result, err := m.GetAll(APIstub)

// 	if err != nil {
// 		return shim.Error("Failed to query milestone due to " + err.Error())
// 	}

// 	return shim.Success(result)
// }

// func (m *MilestoneChaincode) GetMilestoneByID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
// 	if len(args) != 1 {
// 		return shim.Error("Incorrect number of arguments. Expecting 1")
// 	}
// 	var id = args[0]
// 	value, err := m.mstRepo.GetByKey(APIstub, id)

// 	if err != nil {
// 		return shim.Error("Failed to get milestone by id: " + id)
// 	}

// 	if string(value) == "" {
// 		return shim.Error("Failed to get milestone because the milestone " + id + " does not exist.")
// 	}

// 	return shim.Success(value)
// }

// func (m *MilestoneChaincode) CreateMilestone(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

// 	if len(args) != 3 {
// 		return shim.Error("Incorrect number of arguments. Expecting 3")
// 	}

// 	var mst = models.Milestone{
// 		MilestoneID:            guid.New().StringUpper(),
// 		MilestoneTranslationID: args[0],
// 		TrackID:                args[1],
// 		Version:                args[2],
// 		DocType:                MilestoneDocType}

// 	id, err := m.mstRepo.CreateMilestone(mst)

// 	if err != nil {
// 		log.LogError(err.Error())
// 		return shim.Error("Failed to create milestone due to " + err.Error())
// 	}

// 	return shim.Success([]byte(id))
// }

// func (m *MilestoneChaincode) DeleteRecord(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

// 	if len(args) != 1 {
// 		return shim.Error("Incorrect number of arguments. Expecting 1")
// 	}

// 	var mstId = args[0]

// 	data, err := m.mstRepo.GetByKey(APIstub, mstId)

// 	if err != nil {
// 		return shim.Error("Failed to delete record " + mstId + " due to " + err.Error())
// 	}

// 	if string(data) == "" {
// 		return shim.Error("Failed to delete record, because the record " + mstId + " does not exist.")
// 	}

// 	err = m.mstRepo.DeleteRecord(APIstub, mstId)

// 	if err != nil {
// 		return shim.Error("Failed to delete record " + mstId + " due to " + err.Error())
// 	}

// 	return shim.Success(nil)
// }

// func (m *MilestoneChaincode) UpdateMilestone(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

// 	if len(args) != 4 {
// 		return shim.Error("Incorrect number of arguments. Expecting 3")
// 	}

// 	var mstId = args[0]
// 	mst, err := m.mstRepo.GetByKey(APIstub, mstId)

// 	if err != nil {
// 		return shim.Error("Failed to update milestone " + args[0] + ": " + err.Error())
// 	}

// 	if string(mst) == "" {
// 		return shim.Error("Failed to update milestone, because the milestone does not exist.")
// 	}

// 	err = m.mstRepo.UpdateMilestone(APIstub, mstId, args[1], args[2], args[3])

// 	if err != nil {
// 		return shim.Error("Failed to update milestone " + args[0] + ": " + err.Error())
// 	}

// 	return shim.Success(nil)
// }

// func (m *MilestoneChaincode) AddMilestoneDependency(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

// 	if len(args) != 2 {
// 		return shim.Error("Incorrect number of arguments. Expecting 2")
// 	}

// 	var dependingMst = args[0]
// 	var milestoneId = args[1]

// 	isDepending, err := m.mstDependency.CheckDependency(APIstub, dependingMst, milestoneId)

// 	if err != nil {
// 		return shim.Error("Faild to add milestone dependency due to " + err.Error())
// 	}

// 	if isDepending {
// 		return shim.Error("The milestone " + args[0] + " has already depended on milestone " + args[1])
// 	}

// 	id, err := m.mstDependency.CreateMilestoneDependency(APIstub, dependingMst, milestoneId)

// 	return shim.Success([]byte(id))
// }

// func (m *MilestoneChaincode) GetMilestoneDependenciesForMilestone(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

// 	if len(args) != 1 {
// 		return shim.Error("Incorrect number of arguments. Expecting 1")
// 	}

// 	var milestoneId = args[0]

// 	data, err := m.mstDependency.GetDependingMstByMstID(APIstub, milestoneId)

// 	if err != nil {
// 		return shim.Error("Failed to get milestone dependeny for " + milestoneId + " due to " + err.Error())
// 	}

// 	return shim.Success(data)
// }

// //args[0] is milestone id, args[1] is skill id
// func (m *MilestoneChaincode) AddSkillsToMilestone(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

// 	if len(args) != 2 {
// 		return shim.Error("Incorrect number of arguments. Expecting 2")
// 	}

// 	var milestoneId = args[0]
// 	var skillId = args[1]

// 	id, err := m.mstSkillRepo.CreateMilestoneSkill(APIstub, milestoneId, skillId)

// 	if err != nil {
// 		return shim.Error("Failed to add skill" + skillId + " to milestone " + milestoneId + " due to " + err.Error())
// 	}

// 	return shim.Success([]byte(id))
// }

// func (m *MilestoneChaincode) GetSkillsByMilestone(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
// 	if len(args) != 1 {
// 		return shim.Error("Incorrect number of arguments. Expecting 1")
// 	}

// 	var milestoneId = args[0]

// 	skillIds, err := m.mstSkillRepo.GetSkillsByMilestone(APIstub, milestoneId)

// 	if err != nil {
// 		return shim.Error("Failed to get skill from the milestone " + milestoneId + " due to " + err.Error())
// 	}

// 	return shim.Success(skillIds)
// }

func main() {
	err := shim.Start(new(MilestoneChaincode))
	if err != nil {
		fmt.Printf("Error creating new Milestone chaincode: %s", err)
	}
}
