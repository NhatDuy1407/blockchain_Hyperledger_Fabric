package core

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type IBase interface {
	ValidateLogin(shim.ChaincodeStubInterface, string) sc.Response
	CheckUserPermission(shim.ChaincodeStubInterface, string, string) sc.Response
}
