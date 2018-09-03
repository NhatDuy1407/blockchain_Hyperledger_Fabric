package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/skillbill/models"
	"github.com/skillbill/packages/utils"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser to anonymous user can register a user
func (s *SecurityChaincode) RegisterUser(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	var adLogin = args[0]
	var roleID = args[1]
	var password = args[2]
	response, _ := utils.GetPublicKey(APIstub)
	publicKey := hex.EncodeToString([]byte(response))

	passbytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return shim.Error(fmt.Sprintf("Could not hash the password, err %s", err))
	}

	return s.AddUser(APIstub, []string{adLogin, publicKey, roleID, string(passbytes)})
}

// AddUser is add new user
func (s *SecurityChaincode) AddUser(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of columns. Expecting 4 but " + strconv.Itoa(len(args)) + " is " + args[0])
	}
	var adLogin = args[0]
	var publicKey = args[1]
	var roleID = args[2]
	var hashedPassword = args[3]

	userAdLoginRes, err := userRepo.Any(APIstub, func(item interface{}) bool { return item.(models.User).ADLogin == adLogin })
	if err != nil {
		return shim.Error(fmt.Sprintf("Get error %s", err))
	}

	b, _ := strconv.ParseBool(string(userAdLoginRes))
	if b == true {
		return shim.Error("User login " + adLogin + " existed already")
	}

	userPublicKeyRes, err := userRepo.Any(APIstub, func(item interface{}) bool { return item.(models.User).PublicKey == publicKey })
	if err != nil {
		return shim.Error(fmt.Sprintf("Get error %s", err))
	}

	b, _ = strconv.ParseBool(string(userPublicKeyRes))
	if b == true {
		return shim.Error("User publickey existed already")
	}

	var user = models.User{ADLogin: adLogin, PublicKey: publicKey, RoleID: roleID, HashedPassword: hashedPassword, DocType: UserTableName}

	data, _ := json.Marshal(user)

	err = userRepo.Save(APIstub, user.ADLogin, data)

	if err != nil {
		return shim.Error("Failed to create knowledge group " + args[0] + " due to " + err.Error())
	}

	return shim.Success([]byte(user.ADLogin))
}

// GetAllUsers is get all users by doctype (user)
func (s *SecurityChaincode) GetAllUsers(APIstub shim.ChaincodeStubInterface) pb.Response {

	users, err := userRepo.GetAll(APIstub)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to get all users %s", err))
	}

	return shim.Success(users)
}

// GetUserByPublicKey is
func (s *SecurityChaincode) GetUserByPublicKey(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	publicKey := args[0]
	query := `{"selector":{"` + UserPublicKeyColumnName + `":"` + publicKey + `"}}`
	resultsIterator, err := APIstub.GetQueryResult(query)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	var users []models.User
	json.Unmarshal([]byte(buffer.String()), &users)
	if len(users) > 1 {
		return shim.Error("multiple users with the same public key " + publicKey)
	}

	fmt.Printf("- getUserByPublicKey:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
