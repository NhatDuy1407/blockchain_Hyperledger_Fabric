package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	"github.com/skillbill/models"
	"github.com/skillbill/packages/logs"
	"github.com/skillbill/packages/utils"
)

// ============================================================================================================================
// Internal Functions - Helpers, Utilities for the chaincode
// ============================================================================================================================

func toChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

// validateUser to get user that exists in the ledger by validationg the public key
func validateUser(stub shim.ChaincodeStubInterface, payload []byte, ecdsaPublicKey *ecdsa.PublicKey, password string) (*models.User, error) {
	var user *models.User

	err := json.Unmarshal(payload, &user)
	if reflect.DeepEqual(models.User{}, user) == false {
		publicKeyParser, _ := x509.MarshalPKIXPublicKey(ecdsaPublicKey)
		publicKey := hex.EncodeToString(publicKeyParser)
		hasPassword := false
		if password != "" {
			hasPassword = utils.CheckPasswordHash(password, user.HashedPassword)
		}

		var matchPub = strings.Compare(publicKey, user.PublicKey)

		if matchPub == 0 && hasPassword {
			return user, err
		}

		err = errors.New(user.ADLogin + " matchPub: " + string(user.PublicKey) + " hasPassword" + user.HashedPassword)
	}

	return nil, err
}

func getFeaturesByRoleIDs(stub shim.ChaincodeStubInterface, roleID string) sc.Response {

	channelName := ""

	// Check user existing or not && get UserID from public key
	chaincodeName := "role"
	functionName := "getFeaturesByRoleIDs"
	queryArgs := toChaincodeArgs(functionName, roleID)

	response := stub.InvokeChaincode(chaincodeName, queryArgs, channelName)
	return response
}

// ============================================================================================================================
// ChainCode Functions - Define functions for the chaincode
// ============================================================================================================================

// ValidateLogin to check user can login (check public key of caller)
func (s *SecurityChaincode) ValidateLogin(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	logs.LogInfo("Validate login")
	password := args[0]
	cert, err := utils.GetCreatorCert(stub)
	ecdsaPublicKey := cert.PublicKey.(*ecdsa.PublicKey)

	if err != nil {
		return shim.Error(fmt.Sprintf("Could not get Certificate, err %s", err))
	}

	// var query = strings.Join([]string{"adlogin", cert.Subject.CommonName}, ",")
	// logs.LogInfo("Query: " + query)

	response, err := userRepo.FirstOrDefault(stub, func(item interface{}) bool { return item.(models.User).ADLogin == cert.Subject.CommonName })

	if err != nil {
		errStr := fmt.Sprintf("Failed to get user. Got error: %s", err)
		logs.LogError(errStr)
		return shim.Error(errStr)
	}

	if response == nil {
		return shim.Error("Could not find any user with name " + cert.Subject.CommonName)
	}

	currentUser, err := validateUser(stub, response, ecdsaPublicKey, password)
	if err != nil {
		return shim.Error(fmt.Sprintf("Could not parse json to user object, err %s", err))
	}

	if currentUser == nil || reflect.DeepEqual(models.User{}, currentUser) {
		return shim.Success([]byte(strconv.FormatBool(false)))
	}

	return shim.Success([]byte(strconv.FormatBool(true)))
}

// CheckUserPermission to check user can access feature
// arg[0] : featureID, arg[1] : accessLevel (0- readonly, 1- write and read)
func (s *SecurityChaincode) CheckUserPermission(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	//featureID := args[0]
	//accessLevel := args[1]
	cert, err := utils.GetCreatorCert(stub)
	ecdsaPublicKey := cert.PublicKey.(*ecdsa.PublicKey)

	if err != nil {
		return shim.Error(fmt.Sprintf("Could not get Certificate, err %s", err))
	}

	response, err := userRepo.FirstOrDefault(stub, func(item interface{}) bool { return item.(models.User).ADLogin == cert.Subject.CommonName })

	if err != nil {
		errStr := fmt.Sprintf("Failed to get user. Got error: %s", err)
		logs.LogError(errStr)
		return shim.Error(errStr)
	}

	if response == nil {
		return shim.Error("Could not find any user with name " + cert.Subject.CommonName)
	}

	currentUser, err := validateUser(stub, response, ecdsaPublicKey, "")
	if err != nil {
		return shim.Error(fmt.Sprintf("Could not parse json to user object, err %s", err))
	}

	if reflect.DeepEqual(models.User{}, currentUser) {
		return shim.Success([]byte(strconv.FormatBool(false)))
	}

	// Check user can access feature : Get features role has
	// response = getFeaturesByRoleIDs(stub, currentUser.RoleID)
	// if response.Status != shim.OK {
	// 	errStr := fmt.Sprintf("Failed to query chaincode. Got error: %s", response.Payload)
	// 	fmt.Printf(errStr)
	// 	return shim.Error(errStr)
	// }

	// roleFeatures := make([]models.RoleFeature, 0)
	// err = json.Unmarshal(response.Payload, &roleFeatures)
	// if err != nil {
	// 	return shim.Error(fmt.Sprintf("Could not parse json to feature object, err %s", err))
	// }

	canAccess := false
	// if reflect.DeepEqual([]models.RoleFeature{}, roleFeatures) == false {
	// 	for _, roleFeature := range roleFeatures {
	// 		accessLevelInt, err := strconv.Atoi(accessLevel)

	// 		if err == nil && roleFeature.FeatureID == featureID && roleFeature.AccessLevel == accessLevelInt {
	// 			canAccess = true
	// 			break
	// 		}
	// 	}
	// }

	return shim.Success([]byte(strconv.FormatBool(canAccess)))
}
