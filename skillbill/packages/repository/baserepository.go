package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	. "github.com/ahmetb/go-linq"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	logs "github.com/skillbill/packages/logs"
)

// BaseRepo to implement basic functions
type BaseRepo struct {
	DocType string
	// Entity  interface{}
}

// InitRepo to create Repo
func InitRepo(doctype string) IRepo {
	return BaseRepo{DocType: doctype}
}

// GetAll to get all data by doctype
func (r BaseRepo) GetAll(APIstub shim.ChaincodeStubInterface) ([]byte, error) {
	query := `{"selector":{"doctype":"` + r.DocType + `"}}`
	return r.GetByQuery(APIstub, query)
}

// Count is len of list by doctype and predicate func
func (r BaseRepo) Count(APIstub shim.ChaincodeStubInterface, predicate func(interface{}) bool) ([]byte, error) {

	var entities, _ = r.GetAll(APIstub)
	var jsonEntity []*interface{}
	err := json.Unmarshal(entities, &jsonEntity)
	if err != nil {
		return nil, err
	}
	count := From(jsonEntity).CountWith(predicate)

	return []byte(strconv.Itoa(count)), nil
}

// FirstOrDefault is get first data
func (r BaseRepo) FirstOrDefault(APIstub shim.ChaincodeStubInterface, predicate func(interface{}) bool) ([]byte, error) {
	var entities, _ = r.GetAll(APIstub)
	var jsonEntity []interface{}
	err := json.Unmarshal(entities, jsonEntity)
	if err != nil {
		return nil, err
	}

	entity := From(jsonEntity).FirstWith(predicate)
	if entity == nil {
		return nil, nil
	}

	entityByte, err := json.Marshal(entity)
	if err != nil {
		return nil, err
	}

	return entityByte, nil
}

// Any is get first data
func (r BaseRepo) Any(APIstub shim.ChaincodeStubInterface, predicate func(interface{}) bool) ([]byte, error) {
	var entities, _ = r.GetAll(APIstub)
	var jsonEntity []*interface{}
	err := json.Unmarshal(entities, &jsonEntity)
	if err != nil {
		return nil, err
	}

	isAny := From(jsonEntity).AnyWith(predicate)

	return []byte(strconv.FormatBool(isAny)), nil
}

// GetByQuery is return list of query entity (by doctype)
func (r BaseRepo) GetByQuery(APIstub shim.ChaincodeStubInterface, query string) ([]byte, error) {

	logs.LogInfo("Calling GetByQuery in the base repo.")

	resultsIterator, err := APIstub.GetQueryResult(query)

	if err != nil {
		return nil, fmt.Errorf(err.Error() + "query: \n" + query)
	}

	defer resultsIterator.Close()
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return buffer.Bytes(), nil
}

// GetByKey is get entity by key
func (r BaseRepo) GetByKey(APIstub shim.ChaincodeStubInterface, key string) ([]byte, error) {

	value, err := APIstub.GetState(key)

	if err != nil {
		return nil, fmt.Errorf("Failed to get record %s", key)
	}

	if len(string(value)) == 0 {
		return nil, fmt.Errorf("The record has been not found.")
	}

	return value, nil
}

// Save is store data into ledger
func (r BaseRepo) Save(APIstub shim.ChaincodeStubInterface, key string, value []byte) error {
	err := APIstub.PutState(key, value)

	return err
}

// Delete is remove data in ledger
func (r BaseRepo) Delete(APIstub shim.ChaincodeStubInterface, key string) error {

	if len(key) < 1 {
		return fmt.Errorf("Entity Ids are empty. Please specify the entity id to delete.")
	}

	obj, ex := APIstub.GetState(key)

	if ex != nil {
		return fmt.Errorf("Failed to delete %s due to %s", key, ex.Error())
	}

	if len(string(obj)) == 0 {
		return fmt.Errorf("Failed to delete, because the key %s does not exist.", key)
	}

	err := APIstub.DelState(key)

	return err
}
