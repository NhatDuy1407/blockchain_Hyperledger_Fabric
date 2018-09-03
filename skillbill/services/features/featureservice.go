package main

import (
	"encoding/json"
	"fmt"

	"github.com/beevik/guid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	//repo "github.com/skillbill/services/core"
)

var ccInstance repo.IChaincode

// Init method is called when the Smart Contract "Feature" is instantiated by the blockchain network
func (s *Feature) Init(APIstub shim.ChaincodeStubInterface) sc.Response {

	//ccInstance = repo.NewChaincode(APIstub)
	return shim.Success(nil)
}

// Invoke method is called as a result of an application request to run the Smart Contract "Feature"
func (s *Feature) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()

	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "initData" {
		return initData(APIstub)
	} else if function == "getByQuery" {
		return s.getByQuery(APIstub, args)
	} else if function == "createFeature" {
		return createFeature(APIstub, args)
	} else if function == "deleteFeature" {
		return deleteFeature(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name: " + function)
}

func initData(APIstub shim.ChaincodeStubInterface) sc.Response {
	features := []Feature{
		Feature{ FeatureID: "025D1E9A-9B52-E811-AA17-FCAA145000C2", FeatureName: "SkillPlanManagement", DocType: "feature" },
		Feature{ FeatureID: "035D1E9A-9B52-E811-AA17-FCAA145000C2", FeatureName: "SkillManagement", DocType: "feature" },
		Feature{ FeatureID: "045D1E9A-9B52-E811-AA17-FCAA145000C2", FeatureName: "TrackManagement", DocType: "feature" },
		Feature{ FeatureID: "055D1E9A-9B52-E811-AA17-FCAA145000C2", FeatureName: "MilestoneManagement", DocType: "feature" },
		Feature{ FeatureID: "065D1E9A-9B52-E811-AA17-FCAA145000C2", FeatureName: "UserManagement", DocType: "feature" },
		Feature{ FeatureID: "075D1E9A-9B52-E811-AA17-FCAA145000C2", FeatureName: "RoleManagement", DocType: "feature" },
		Feature{ FeatureID: "085D1E9A-9B52-E811-AA17-FCAA145000C2", FeatureName: "KnowledgeGroup", DocType: "feature" },
		Feature{ FeatureID: "095D1E9A-9B52-E811-AA17-FCAA145000C2", FeatureName: "FeatureManagement", DocType: "feature" },
		Feature{ FeatureID: "0A5D1E9A-9B52-E811-AA17-FCAA145000C2", FeatureName: "TranslationManagement", DocType: "feature" },
	}

	i := 0
	for i < len(features) {
		data, _ := json.Marshal(features[i])
		ccInstance.Save(APIstub, features[i].FeatureID, data)

		i = i + 1
	}

	return shim.Success(nil)
}

func (s *Feature) getByQuery(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var result = s.repo.GetByQuery(args)

	return shim.Success([]byte(result))
}

func createFeature(APIstub shim.ChaincodeStubInterface, args[] string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var feature = Feature{ FeatureID: guid.New().StringUpper(), FeatureName: args[0], DocType: "feature" }
	data, _ := json.Marshal(feature)
	err := ccInstance.Save(APIstub, feature.FeatureID, data)

	if err != nil {
		return shim.Error("Failed to create feature due to " + err.Error())
	}
	
	return shim.Success([]byte(feature.FeatureID))
}

func deleteFeature(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var featureId = args[0]
	err := ccInstance.Delete(APIstub, featureId)

	if err != nil {
		shim.Error("Failed to delete feature " + featureId + " due to " + err.Error())
	}

	return shim.Success(nil)
}

func main() {
	
	err := shim.Start(new(Feature))
	
	if err != nil {
		fmt.Printf("Error creating new Skill Chaincode: %s", err)
	}
}
