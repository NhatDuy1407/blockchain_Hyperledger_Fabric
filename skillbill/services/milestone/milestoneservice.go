package main

import (
	"encoding/json"

	"github.com/beevik/guid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	"github.com/skillbill/models"
)

func (m MilestoneChaincode) GetAllMilestones(APIstub shim.ChaincodeStubInterface) sc.Response {

	result, err := milestoneRepo.GetAll(APIstub)

	if err != nil {
		return shim.Error("Failed to query milestone due to " + err.Error())
	}
	return shim.Success(result)
}

func (m MilestoneChaincode) GetMilestoneByID(APIstub shim.ChaincodeStubInterface, key string) sc.Response {
	value, err := milestoneRepo.GetByKey(APIstub, key)

	if err != nil {
		return shim.Error("Failed to get milestone due to: " + err.Error())
	}

	if string(value) == "" {
		return shim.Error("Failed to get milestone because the milestone " + key + " does not exist.")
	}

	return shim.Success(value)
}

func (m MilestoneChaincode) CreateMilestone(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	var mst = models.Milestone{
		MilestoneID:            guid.New().StringUpper(),
		MilestoneTranslationID: args[0],
		TrackID:                args[1],
		Version:                args[2],
		DocType:                MilestoneDocType}

	data, _ := json.Marshal(mst)
	err := milestoneRepo.Save(APIstub, mst.MilestoneID, data)

	if err != nil {
		return shim.Error("Failed to create milestone due to: " + err.Error())
	}

	return shim.Success([]byte(mst.MilestoneID))
}

func (m MilestoneChaincode) UpdateMilestone(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	data, err := milestoneRepo.GetByKey(APIstub, args[0])

	if err != nil {
		return shim.Error("Failed to update the milestone: " + args[0] + "due to: " + err.Error())
	}

	if len(string(data)) == 0 {
		return shim.Error("Failed to update the milestone, because it does not exist.")
	}

	mst := models.Milestone{}
	json.Unmarshal(data, &mst)
	mst.MilestoneTranslationID = args[1]
	mst.TrackID = args[2]
	mst.Version = args[3]
	data, _ = json.Marshal(mst)

	err = milestoneRepo.Save(APIstub, mst.MilestoneID, data)

	if err != nil {
		return shim.Error("Failed to update the milestone: " + args[0] + "due to: " + err.Error())
	}

	return shim.Success(nil)
}

func (m MilestoneChaincode) DeleteRecord(APIstub shim.ChaincodeStubInterface, key string) sc.Response {
	err := milestoneRepo.Delete(APIstub, key)

	if err != nil {
		return shim.Error("Failed to delete record with key : " + key + " due to " + err.Error())
	}

	return shim.Success(nil)
}
