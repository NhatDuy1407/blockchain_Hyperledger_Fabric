package main


import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/beevik/guid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"

	log "github.com/sirupsen/logrus"
)

type SkillPlanChaincode struct {

}

// Init method is called when the Smart Contract "Feature" is instantiated by the blockchain network
func (s *SkillPlanChaincode) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	setUpLogging(filepath)

	return shim.Success(nil)
}

// Invoke method is called as a result of an application request to run the Smart Contract "skill plan"
func (s *SkillPlanChaincode) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()

	log.Info("Function is calling: " , function)
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "getAllByQuery" {
		return getAllByQuery(APIstub, args)
	} else if function == "createSkillPlan" {
		return createSkillPlan(APIstub, args)
	} else if function == "deleteSkillPlan" {
		return deleteSkillPlan(APIstub, args)
	} else if function == "updatePlannedSkill" {
		return updatePlannedSkill(APIstub, args)
	} else if function == "updateCompletedSkill" {
		return updateCompletedSkill(APIstub, args)
	} else if function == "updateAssessmentRequest" {
		return updateAssessmentRequest(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func createSkillPlan(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	// Build skill plan object base on type
	id, data, errMsg := buildSkillPlanObject(args)

	if errMsg != "" {
		log.Info(errMsg)

		return shim.Error(errMsg)
	}

	log.Info("Skill Plan Id: %s", id)
	APIstub.PutState(id, data)

	return shim.Success([]byte(id))
}

func deleteSkillPlan(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Skill id are empty. Please specify the skill id to delete.")
	}

	var id = args[0]
	data, e := APIstub.GetState(id)

	if e != nil{
		return shim.Error("Failed to delete skill plan" + id + " due to " + e.Error() + "\n")
	}

	if len(string(data)) == 0 {
		return shim.Error("Failed to delete skill plan, because the skill " + id + " does not exist.")
	}
	
	err := APIstub.DelState(id)

	if err != nil{
		return shim.Error("Failed to delete skill plan " + id + " due to " + err.Error())
	}

	return shim.Success(nil)
}

func updatePlannedSkill(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	data, err := APIstub.GetState(args[0])

	if err != nil{
		return shim.Error("Failed to update planned skill " + args[0] + ": " + err.Error())
	}

	if len(string(data)) == 0 {
		return shim.Error("Failed to update planned skill, because the skill does not exist.")
	}

	plannedskill := SkillPlanPlannedSkill{}
	json.Unmarshal(data, &plannedskill)

	plannedskill.PlannedFrom = args[1]
	plannedskill.PlannedTo = args[2]
	plannedskill.Priority = args[3]
	plannedskill.SkillID = args[4]
	plannedskill.UserID = args[5]

	data, _ = json.Marshal(plannedskill)

	APIstub.PutState(args[0], data)

	return shim.Success(nil)
}

func updateCompletedSkill(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	data, err := APIstub.GetState(args[0])

	if err != nil{
		return shim.Error("Failed to update completed skill " + args[0] + ": " + err.Error())
	}

	if len(string(data)) == 0 {
		return shim.Error("Failed to update completed skill, because the skill does not exist.")
	}

	skill := SkillPlanCompletedSkill{}
	json.Unmarshal(data, &skill)
	skill.AccessedBy = args[1]
	skill.CompletedOn = args[2]
	skill.SkillID = args[3]
	skill.UserID = args[4]

	data, _ = json.Marshal(skill)

	APIstub.PutState(args[0], data)

	return shim.Success(nil)
}

func updateAssessmentRequest(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	data, err := APIstub.GetState(args[0])

	if err != nil{
		return shim.Error("Failed to update completed skill " + args[0] + ": " + err.Error())
	}

	if len(string(data)) == 0 {
		return shim.Error("Failed to update planned skill, because the skill does not exist.")
	}

	skill := SkillPlanAssessmentRequest{}
	json.Unmarshal(data, &skill)
	skill.AssesseeID = args[1]
	skill.AssessorID = args[2]
	skill.SkillID = args[3]

	data, _ = json.Marshal(skill)

	APIstub.PutState(args[0], data)

	return shim.Success(nil)
}

// args[0].. args[n] are pair column and value
// e.g: args['doctype,milestone', 'milestoneid,E1ED5DAD-B286-4522-8A93-926E6D5DC9C9', ....]
func getAllByQuery(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response {
	
	var b bytes.Buffer
	i := 0
	for i < len(args) {
		var params = strings.Split(args[i], ",")
		b.WriteString("\""+ params[0] +"\":\""+ params[1] +"\"")
		i = i + 1
		if i != len(args){
			b.WriteString(",");
		}
	}

	var query = `{"selector":{`+ b.String() +`}}`

	resultsIterator, err := APIstub.GetQueryResult(query)

	if err != nil {
		return shim.Error(err.Error() + "query: \n"+ query)
	}
	
	defer resultsIterator.Close()
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("Data: \n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func buildSkillPlanObject(args []string) (string, []byte, string) {
	objType := strings.ToLower(args[0])

	log.Info("Action Type: ", objType)
	
	switch objType {
	case Planned:
		var skill = SkillPlanPlannedSkill { 
			ID: guid.New().StringUpper(), 
			PlannedFrom: args[1], 
			PlannedTo: args[2], 
			Priority: args[3], 
			SkillID: args[4], 
			UserID: args[5], 
			DocType: "plannedskill"}
	
		data, _ := json.Marshal(skill)

		return skill.ID, data, ""
	
	case InProgress:
		var skill = SkillPlanInProgressSkill { 
			ID: guid.New().StringUpper(), 
			SkillACID: args[1], 
			SkillACStartdate: args[2],
			SkillID: args[3], 
			UserID: args[4], 
			DocType: "inprogressskill"}
	
		data, _ := json.Marshal(skill)

		return skill.ID, data, ""

	case Completed:
		var skill = SkillPlanCompletedSkill { 
			ID: guid.New().StringUpper(), 
			AccessedBy: args[1], 
			CompletedOn: args[2],
			SkillID: args[3], 
			UserID: args[4], 
			DocType: "completedskill"}
	
		data, _ := json.Marshal(skill)

		return skill.ID, data, ""

	case AssessmentRequest:
		var skill = SkillPlanAssessmentRequest { 
			ID: guid.New().StringUpper(), 
			AssesseeID: args[1], 
			AssessorID: args[2],
			SkillID: args[3], 
			DocType: "assessmentrequest"}
	
		data, _ := json.Marshal(skill)

		return skill.ID, data, ""

	default:
			return "Invalid skill plan type: " + objType, nil, ""
	}
}

func main() {
	err := shim.Start(new(SkillPlanChaincode))
	if err != nil {
		fmt.Printf("Error creating new skill plan chaincode: %s", err)
	}
}
