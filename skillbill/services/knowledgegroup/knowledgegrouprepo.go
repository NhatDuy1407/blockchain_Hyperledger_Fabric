package main

import (
	"fmt"
	"bytes"
	"strings"
	"encoding/json"
	"github.com/beevik/guid"
	"github.com/skillbill/models"
	"github.com/skillbill/packages/repository"
	"github.com/hyperledger/fabric/core/chaincode/shim"

	logs "github.com/skillbill/packages/logs"
)

type KnowledgeGroupRepo struct {
	repo	repository.BaseRepo
}

func InitKnowledgeGrpRepo() IKnowledgeGrp {
	return KnowledgeGroupRepo{}
}

func (k KnowledgeGroupRepo) GetByQuery(APIstub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	logs.LogInfo("Calling GetByQuery in the knowledgegrp repo.")

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

	logs.LogInfo("Query string: " + query)

	return k.repo.GetByQuery(APIstub, query)
}

func (k KnowledgeGroupRepo) GetByKey(APIstub shim.ChaincodeStubInterface, key string) ([]byte, error) {

	value, err := k.repo.GetByKey(APIstub, key)
	
	if err != nil {
		return nil, fmt.Errorf("Failed to get record %s", key) 
	}

	return value, nil
}

// func (k KnowledgeGroupRepo) CheckIsExisted(APIstub shim.ChaincodeStubInterface, id string, doctype string)	(bool, error) {
// 	var query = `{"selector":{"GroupID":"` + groupID + `", "DocType": "knowledgegroupmember"}, "fields": ["UserID", "MemberType"]}`

// 	result, err := k.repo.Query()
// }

func (k KnowledgeGroupRepo) CreateKnowledgeGrp(APIstub shim.ChaincodeStubInterface, groupName string) (string, error) {

	var kgroup = models.KnowledgeGroup{ GroupID: guid.New().StringUpper(), GroupName: groupName, DocType: "knowledgegroup" }

	data, _ := json.Marshal(kgroup)

	err := k.repo.Save(APIstub, kgroup.GroupID, data)

	return kgroup.GroupID, err
}

func (k KnowledgeGroupRepo) UpdateKnowledgeGrp(APIstub shim.ChaincodeStubInterface, params []string) error {

	value, err := k.GetByKey(APIstub, params[0])

	if err != nil{
		return fmt.Errorf("Failed to update knowledge group %s: %s", params[0], err.Error())
	}

	if len(string(value)) == 0 {
		return fmt.Errorf("Failed to update knowledge group, because the knowledge group does not exist.")
	}

	kngroup := models.KnowledgeGroup{}
	json.Unmarshal(value, &kngroup)
	kngroup.GroupName = params[1]
	logs.LogInfo(string(value))
	value, _ = json.Marshal(kngroup)

	return k.repo.Save(APIstub, params[0], value)
}

func (k KnowledgeGroupRepo) DeleteRecord(APIstub shim.ChaincodeStubInterface, groupID string) error {
	return k.repo.Delete(APIstub, groupID)
}

func (k KnowledgeGroupRepo) AddMembersToKnowledgeGrp(APIstub shim.ChaincodeStubInterface, groupID string, memberType string, userID string) (string, error) {
	var kngroupMem = models.KnowledgeGroupMember { 
		ID: guid.New().StringUpper(), 
		GroupID: groupID, 
		MemberType: memberType, 
		UserID: userID, 
		DocType: "knowledgegroupmember",
	}
	dataAsByte, _ := json.Marshal(kngroupMem)
	err := k.repo.Save(APIstub, kngroupMem.ID, dataAsByte)

	return kngroupMem.ID, err
}

func (k KnowledgeGroupRepo) GetMembersByGroupID(APIstub shim.ChaincodeStubInterface, groupID string) ([]byte, error) {
	var query = `{"selector":{"GroupID":"` + groupID + `", "DocType": "knowledgegroupmember"}, "fields": ["UserID", "MemberType"]}`

	return k.repo.GetByQuery(APIstub, query)
}

