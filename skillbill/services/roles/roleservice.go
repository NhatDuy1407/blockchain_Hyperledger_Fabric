package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/beevik/guid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	repo "github.com/skillbill/services/core"

	log "github.com/sirupsen/logrus"
)

var ccInstance repo.IChaincode

// Init method is called when the Smart Contract "Role" is instantiated by the blockchain network
func (s *Role) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	setUpLogging(filepath)

	return shim.Success(nil)
}

// Invoke method is called as a result of an application request to run the Smart Contract "Role"
func (s *Role) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()

	log.Info("Route to function based on function name")
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "initData" {
		return initRolesAndFeature(APIstub)
	} else if function == "getAllByQuery" {
		return getAllByQuery(APIstub, args)
	} else if function == "createRole" {
		return createRole(APIstub, args)
	} else if function == "deleteRole" {
		return deleteRole(APIstub, args)
	} else if function == "assignFeature" {
		return assignFeatureRole(APIstub, args)
	} else if function == "removeFeature" {
		return removeFeatureFromRole(APIstub, args)
	} else if function == "getFeaturesByRoleIDs" {
		return getFeaturesByRoleIDs(APIstub, args)
	}
	var errMsg = "Invalid Smart Contract function name: "+ function
	log.Error(errMsg)

	return shim.Error(errMsg)
}

func assignFeatureRole(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	accessLevel, err := parseStr2Enum(args[0])

	if err != nil {
		return shim.Error(err.Error())
	}	

	var roleFeature = RoleFeature{ ID: guid.New().StringUpper(), AccessLevel: accessLevel, RoleID: args[1], FeatureID: args[2], DocType: "rolefeature" }
	data, _ := json.Marshal(roleFeature)

	ccInstance.Save(APIstub, roleFeature.ID, data)

	log.Info("Assigned feature %s to role %s." , roleFeature.FeatureID, roleFeature.RoleID)

	return shim.Success(nil)
}

func removeFeatureFromRole(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	roleFeatures, err := APIstub.GetQueryResult(`{"selector":
		{"DocType":"rolefeature", "RoleId": "` + args[0] + `", "FeatureId": "` + args[1] + `"}}`)

	if roleFeatures == nil {
		return shim.Error("The role " + args[0] + "does not have feature " + args[1])
	}

	if err != nil {
		return shim.Error("Failed to remove the feature from role.")
	}

	for roleFeatures.HasNext(){
		reponse, err := roleFeatures.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		ccInstance.Delete(APIstub, reponse.Key)

		log.Info("Removed feature %s out of %s", args[1], args[0])
	}

	return shim.Success(nil)
}

func initRolesAndFeature(APIstub shim.ChaincodeStubInterface) sc.Response {
	
	// Init roles
	log.Info("Initializing role data.")
	roles := []Role{
		Role{ RoleID: guid.New().StringUpper(), RoleName: "Administrators", DocType: "role" },
		Role{ RoleID: guid.New().StringUpper(), RoleName: "ProfessionalGroupAdministrators", DocType: "role" },
		Role{ RoleID: guid.New().StringUpper(), RoleName: "SkillAdministartors", DocType: "role" },
		Role{ RoleID: guid.New().StringUpper(), RoleName: "Users", DocType: "role" },
	}

	i := 0
	for i < len(roles) {
		data, _ := json.Marshal(roles[i])
		APIstub.PutState(roles[i].RoleID, data)
		i = i + 1
	}

	// Init features for Administrator role.
	log.Info("Initializing role feature data with Administrator and User.")
	roleFeatures := []RoleFeature{
		RoleFeature{ ID: guid.New().StringUpper(), AccessLevel: ReadWrite, RoleID: roles[0].RoleID, FeatureID: "025D1E9A-9B52-E811-AA17-FCAA145000C2", DocType: "rolefeature"},
		RoleFeature{ ID: guid.New().StringUpper(), AccessLevel: ReadWrite, RoleID: roles[0].RoleID, FeatureID: "035D1E9A-9B52-E811-AA17-FCAA145000C2", DocType: "rolefeature"},
		RoleFeature{ ID: guid.New().StringUpper(), AccessLevel: ReadWrite, RoleID: roles[0].RoleID, FeatureID: "045D1E9A-9B52-E811-AA17-FCAA145000C2", DocType: "rolefeature"},
		RoleFeature{ ID: guid.New().StringUpper(), AccessLevel: ReadWrite, RoleID: roles[0].RoleID, FeatureID: "055D1E9A-9B52-E811-AA17-FCAA145000C2", DocType: "rolefeature"},
		RoleFeature{ ID: guid.New().StringUpper(), AccessLevel: ReadWrite, RoleID: roles[0].RoleID, FeatureID: "065D1E9A-9B52-E811-AA17-FCAA145000C2", DocType: "rolefeature"},
		RoleFeature{ ID: guid.New().StringUpper(), AccessLevel: ReadWrite, RoleID: roles[0].RoleID, FeatureID: "075D1E9A-9B52-E811-AA17-FCAA145000C2", DocType: "rolefeature"},
		RoleFeature{ ID: guid.New().StringUpper(), AccessLevel: ReadWrite, RoleID: roles[0].RoleID, FeatureID: "085D1E9A-9B52-E811-AA17-FCAA145000C2", DocType: "rolefeature"},
		RoleFeature{ ID: guid.New().StringUpper(), AccessLevel: ReadWrite, RoleID: roles[0].RoleID, FeatureID: "095D1E9A-9B52-E811-AA17-FCAA145000C2", DocType: "rolefeature"},
		RoleFeature{ ID: guid.New().StringUpper(), AccessLevel: ReadWrite, RoleID: roles[0].RoleID, FeatureID: "0A5D1E9A-9B52-E811-AA17-FCAA145000C2", DocType: "rolefeature"},
		RoleFeature{ ID: guid.New().StringUpper(), AccessLevel: ReadWrite, RoleID: roles[3].RoleID, FeatureID: "075D1E9A-9B52-E811-AA17-FCAA145000C2", DocType: "rolefeature"},
		RoleFeature{ ID: guid.New().StringUpper(), AccessLevel: ReadWrite, RoleID: roles[3].RoleID, FeatureID: "085D1E9A-9B52-E811-AA17-FCAA145000C2", DocType: "rolefeature"},
	}

	j := 0
	for j < len(roleFeatures) {
		data, _ := json.Marshal(roleFeatures[j])
		APIstub.PutState(roleFeatures[j].ID, data)
		j = j + 1
	}

	return shim.Success(nil)
}

func createRole(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var role = Role{ RoleID: guid.New().StringUpper(), RoleName: args[0], DocType: "role" }
	data, _ := json.Marshal(role)

	err := ccInstance.Save(APIstub, role.RoleID, data)

	if err != nil {
		return shim.Error("Failed to create role due to " + err.Error())
	}

	log.Info("Created the role %s successfully.", role.RoleName)

	return shim.Success([]byte(role.RoleID))
}

// args[0].. args[n] are pair column and value
// e.g: args['DocType,role', 'RoleID,E1ED5DAD-B286-4522-8A93-926E6D5DC9C9', ....]
func getAllByQuery(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response {
	
	var b bytes.Buffer
	i := 0
	for i < len(args) {

		var params = strings.Split(args[i], ",")
		b.WriteString("\""+params[0]+"\":\""+ params[1]+"\"")
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

func deleteRole(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var roleId = args[0]

	ccInstance.GetByQuery(APIstub, [''])

	data, e := APIstub.GetState(roleId)

	if e != nil{
		return shim.Error("Failed to delete role " + roleId + " due to " + e.Error() + "\n")
	}

	if len(string(data)) == 0 {
		return shim.Error("Failed to delete role, because the role " + roleId + " does not exist.")
	}

	err := APIstub.DelState(roleId)

	if err != nil {
		shim.Error("Failed to delete role " + roleId + " due to " + err.Error())
	}

	return shim.Success(nil)
}

func getFeaturesByRoleIDs(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response{
	if len(args) == 0{
		return shim.Error("The args is empty, please specify the role ids.")
	}

	var roleIDs = strings.Split(args[0], ",")

	// Build mango query
	var b bytes.Buffer
	i := 0

	for i < len(roleIDs) {
		b.WriteString(`{"RoleID": "` + roleIDs[i] + `"}`)
		i = i + 1
		if i != len(roleIDs){
			b.WriteString(",");
		}
	}

	query := `{
				"selector": {
				"$or": [
					` + b.String() + `
				],
				"DocType": "rolefeature"
				},
				"fields": [
				"AccessLevel",
				"RoleID",
				"FeatureID"
				]}`
	
	log.Info(" query:\n%s\n", query)
	resultsIterator, err := APIstub.GetQueryResult(query)

	if err != nil {
		return shim.Error(err.Error())
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
	log.Info("Result :\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func parseStr2Enum(name string) (int, error) {
	val, err := strconv.Atoi(name)

	if err != nil {
		return Unknown, fmt.Errorf("Invalid type: failed to parse action %s", name)
	}

	if val == ReadOnly {
		return ReadOnly, nil
	} 
	
	if val == ReadWrite {
		return ReadWrite, nil
	}

	return Unknown, fmt.Errorf("%s is not a valid action", name)
}

func main() {
	ccInstance = repo.NewChaincode()
	err := shim.Start(new(Role))
	if err != nil {
		fmt.Printf("Error creating new role Chaincode: %s", err)
	}
}
