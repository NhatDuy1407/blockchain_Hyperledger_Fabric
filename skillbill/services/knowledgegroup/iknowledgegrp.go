package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type IKnowledgeGrp interface {
	GetByQuery(APIstub shim.ChaincodeStubInterface, args []string)							([]byte, error)
	GetByKey(APIstub shim.ChaincodeStubInterface, id string)								([]byte, error)

	CreateKnowledgeGrp(APIstub shim.ChaincodeStubInterface, groupName string)				(string, error)
	UpdateKnowledgeGrp(APIstub shim.ChaincodeStubInterface, params []string)				error
	DeleteRecord(APIstub shim.ChaincodeStubInterface, id string)							error

	AddMembersToKnowledgeGrp(APIstub shim.ChaincodeStubInterface, groupID string, memberType string, userID string)			(string, error)
	GetMembersByGroupID(APIstub shim.ChaincodeStubInterface, groupID string) 				([]byte, error)
}