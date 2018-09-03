package main

import (
	"fmt"

	"github.com/skillbill/packages/repository"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

const UserTableName string = "user"

const UserPublicKeyColumnName string = "publickey"

const UserADLoginColumnName string = "publickey"

// SecurityChaincode provides functions to manage authorization
type SecurityChaincode struct {
}

var userRepo repository.IRepo

// ============================================================================================================================
// Base Functions - Invoke | Init
// ============================================================================================================================

// Init method is called when the Smart Contract is instantiated by the blockchain network
func (s *SecurityChaincode) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	userRepo = repository.InitRepo("user")
	return shim.Success(nil)
}

// Invoke method is called as a result of an application request to run the Smart Contract
func (s *SecurityChaincode) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()

	if function == "ValidateLogin" {
		return s.ValidateLogin(APIstub, args)
	} else if function == "CheckUserPermission" {
		return s.CheckUserPermission(APIstub, args)
	} else if function == "RegisterUser" {
		return s.RegisterUser(APIstub, args)
	} else if function == "AddUser" {
		return s.AddUser(APIstub, args)
	} else if function == "GetAllUsers" {
		return s.GetAllUsers(APIstub)
	} else if function == "GetUserByPublicKey" {
		return s.GetUserByPublicKey(APIstub, args)
	}

	return shim.Error("Function with the name " + function + " does not exist.")
}

func main() {
	err := shim.Start(new(SecurityChaincode))
	if err != nil {
		fmt.Printf("Error creating new Authorization Chaincode: %s", err)
	}
}
