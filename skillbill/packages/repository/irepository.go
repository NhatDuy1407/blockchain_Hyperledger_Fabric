package repository

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// IRepo is an interface to wrap generic functions to communicate with database
type IRepo interface {
	GetAll(APIstub shim.ChaincodeStubInterface) ([]byte, error)
	FirstOrDefault(APIstub shim.ChaincodeStubInterface, predicate func(interface{}) bool) ([]byte, error)
	Any(APIstub shim.ChaincodeStubInterface, predicate func(interface{}) bool) ([]byte, error)
	Count(APIstub shim.ChaincodeStubInterface, predicate func(interface{}) bool) ([]byte, error)
	GetByQuery(APIstub shim.ChaincodeStubInterface, query string) ([]byte, error)
	GetByKey(APIstub shim.ChaincodeStubInterface, key string) ([]byte, error)
	Save(APIstub shim.ChaincodeStubInterface, key string, value []byte) error
	Delete(APIstub shim.ChaincodeStubInterface, id string) error
}
