package core

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// CreateBase is constructor
func CreateBase() IBase {
	return Base{}
}

type Base struct {
}

func toChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

// ValidateLogin to check user can login (check public key of caller)
func (t Base) ValidateLogin(stub shim.ChaincodeStubInterface, password string) sc.Response {

	channelName := ""
	chaincodeName := "security"
	functionName := "ValidateLogin"

	queryArgs := toChaincodeArgs(functionName, password)

	response := stub.InvokeChaincode(chaincodeName, queryArgs, channelName)
	return response
}

// CheckUserPermission to check user can access feature
// arg[0] : featureID, arg[1] : accessLevel (0- readonly, 1- write and read)
func (t Base) CheckUserPermission(stub shim.ChaincodeStubInterface, featureID string, accessLevel string) sc.Response {

	channelName := ""
	chaincodeName := "security"
	functionName := "CheckUserPermission"

	queryArgs := toChaincodeArgs(functionName, featureID, accessLevel)

	response := stub.InvokeChaincode(chaincodeName, queryArgs, channelName)
	return response
}
