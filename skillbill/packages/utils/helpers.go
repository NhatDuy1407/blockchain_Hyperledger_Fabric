package utils

import (
	"crypto/x509"
	"encoding/pem"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/msp"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword to hash the password with cost is 10
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compare password and hashed password match or not
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetCurrentUser get username
func GetCurrentUser(stub shim.ChaincodeStubInterface) (string, error) {
	cert, err := GetCreatorCert(stub)

	return cert.Subject.CommonName, err
}

// GetPublicKey to get the public key of current user from their certificate
func GetPublicKey(stub shim.ChaincodeStubInterface) (string, error) {
	cert, err := GetCreatorCert(stub)

	bytePublicKey, _ := x509.MarshalPKIXPublicKey(cert.PublicKey)
	return string(bytePublicKey), err
}

// GetCreatorCert to get certificate
func GetCreatorCert(stub shim.ChaincodeStubInterface) (*x509.Certificate, error) {
	creator, err := stub.GetCreator()
	if err != nil {
		return nil, err
	}
	id := &msp.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(id.IdBytes)
	cert, err := x509.ParseCertificate(block.Bytes)
	return cert, err
}
