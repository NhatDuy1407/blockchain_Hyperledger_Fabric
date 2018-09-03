package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type ITrackRepo interface {
	GetByQuery(APIstub shim.ChaincodeStubInterface, args []string) 	([]byte, error)
	GetByKey(APIstub shim.ChaincodeStubInterface, key string) 	([]byte, error)
	CreateTrack(APIstub shim.ChaincodeStubInterface, trackTranslation string, version string) 	(string, error)
	UpdateTrack(APIstub shim.ChaincodeStubInterface, trackID string, trackTranslation string, version string) 	error
	DeleteTrack(APIstub shim.ChaincodeStubInterface, key string) 	error
}