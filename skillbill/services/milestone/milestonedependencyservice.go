package main

import (
	"encoding/json"

	. "github.com/ahmetb/go-linq"
	"github.com/beevik/guid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	"github.com/skillbill/models"
)

func (m MilestoneChaincode) CreateMilestoneDependency(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// Check the existing both of milestone before adding the dependent.
	result, err := getCountOfMilestone(APIstub, func(item interface{}) bool {
		return item.(models.Milestone).MilestoneID == args[0] || item.(models.Milestone).MilestoneID == args[1]
	})

	if err != nil {
		return shim.Error("Failed to add depending the milestone " + args[1] + " to " + args[0] + " due to " + err.Error())
	}

	if result != 2 {
		return shim.Error("Failed to add depending the milestone " + args[1] + " to " + args[0] + " due to not exit milestones.")
	}

	isExisted, err := checkExistForMstDependency(APIstub, func(item interface{}) bool {
		return item.(models.MilestoneDependency).DependingMilestone == args[0] && item.(models.MilestoneDependency).MilestoneID == args[1]
	})

	if err != nil {
		return shim.Error("Failed to add depending the milestone " + args[1] + " to " + args[0] + " due to " + err.Error())
	}

	if isExisted {
		return shim.Error("The milestone " + args[1] + "has been depended to " + args[0])
	}

	var mstDependency = models.MilestoneDependency{
		ID:                 guid.New().StringUpper(),
		DependingMilestone: args[0],
		MilestoneID:        args[1],
		DocType:            MilestoneDependencyDocType}

	data, err := json.Marshal(mstDependency)
	if err != nil {
		return shim.Error("Failed to add depending the milestone " + args[1] + " to " + args[0] + "due to " + err.Error())
	}

	err = milestoneDependencyRepo.Save(APIstub, mstDependency.ID, data)

	if err != nil {
		return shim.Error("Failed to add depending the milestone " + args[1] + " to " + args[0])
	}

	return shim.Success([]byte(mstDependency.ID))
}

func (m MilestoneChaincode) GetDependingsByID(APIstub shim.ChaincodeStubInterface, milestoneID string) sc.Response {

	var query = `{"selector":{"doctype": "` + MilestoneDependencyDocType + `", "dependingmilestone":"` + milestoneID + `"}}`

	data, err := milestoneDependencyRepo.GetByQuery(APIstub, query)

	if err != nil {
		return shim.Error("Failed to get milestone depending for milestone: " + milestoneID)
	}

	return shim.Success(data)
}

func (m MilestoneChaincode) UpdateMilestoneDependency(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	var mstDependencyID = args[0]
	mstDepending, err := milestoneDependencyRepo.GetByKey(APIstub, mstDependencyID)

	if err != nil {
		return shim.Error("Failed to update the milestone depending for key: " + mstDependencyID)
	}

	if len(string(mstDepending)) == 0 {
		return shim.Error("Failed to update the milestone depending, because it does not exist.")
	}

	err = milestoneDependencyRepo.Delete(APIstub, mstDependencyID)

	if err != nil {
		return shim.Error("Failed to update the milestone depending for key: " + mstDependencyID)
	}

	newMstDepending := models.MilestoneDependency{
		ID:                 mstDependencyID,
		DependingMilestone: args[1],
		MilestoneID:        args[2],
		DocType:            MilestoneDependencyDocType}

	data, _ := json.Marshal(newMstDepending)

	err = milestoneDependencyRepo.Save(APIstub, newMstDepending.ID, data)

	if err != nil {
		return shim.Error("Failed to update depending the milestone " + args[2] + " to " + args[1])
	}

	return shim.Success([]byte(newMstDepending.ID))

}

func checkExistForMstDependency(APIstub shim.ChaincodeStubInterface, predicate func(interface{}) bool) (bool, error) {
	var entities, _ = milestoneDependencyRepo.GetAll(APIstub)
	var jsonEntity []models.MilestoneDependency
	err := json.Unmarshal(entities, &jsonEntity)
	if err != nil {
		return nil, err
	}

	isAny := From(&jsonEntity).AnyWith(predicate)

	return isAny, nil
}

func getCountOfMilestone(APIstub shim.ChaincodeStubInterface, predicate func(interface{}) bool) (int, error) {

	var entities, _ = milestoneRepo.GetAll(APIstub)
	var jsonEntity []models.Milestone
	err := json.Unmarshal(entities, &jsonEntity)
	if err != nil {
		return 0, err
	}
	count := From(jsonEntity).CountWithT(predicate)

	return count, nil
}
