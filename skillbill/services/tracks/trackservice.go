package main


import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// RoleChaincode define the Smart Contract structure
type TrackChaincode struct {
	repo 	TrackRepo
}

// Init method is called when the Smart Contract "Feature" is instantiated by the blockchain network
func (t *TrackChaincode) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	InitRepo()

	return shim.Success(nil)
}

// Invoke method is called as a result of an application request to run the Smart Contract "Feature"
func (t *TrackChaincode) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()

	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "GetAllByQuery" {
		return t.GetAllByQuery(APIstub, args)
	} else if function == "GetTrackByID" {
		return t.GetTrackByID(APIstub, args)
	} else if function == "CreateTrack" {
		return t.CreateTrack(APIstub, args)
	} else if function == "UpdateTrack" {
		return t.UpdateTrack(APIstub, args)
	} else if function == "DeleteTrack" {
		return t.DeleteTrack(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

// args[0].. args[n] are pair column and value
// e.g: args['DocType,track', 'TrackID,E1ED5DAD-B286-4522-8A93-926E6D5DC9C9', ....]
func (t *TrackChaincode) GetAllByQuery(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	
	result, err := t.repo.GetByQuery(APIstub, args)

	if err != nil {
		return shim.Error("Failed to query track due to " + err.Error())
	}

	return shim.Success(result)
}

func (t *TrackChaincode) GetTrackByID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var id = args[0]
	value, err := t.repo.GetByKey(APIstub, id)

	if err != nil{
		return shim.Error("Failed to delete track " + id + " due to " + err.Error())
	}

	if string(value) == "" {
		return shim.Error("Failed to delete track, because the track " + id + " does not exist.")
	}

	return shim.Success([]byte(value))
}

func (t *TrackChaincode) CreateTrack(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	id, err := t.repo.CreateTrack(APIstub, args[0], args[1])

	if err != nil {
		return shim.Error("Failed to create track due to " + err.Error())
	}

	return shim.Success([]byte(id))
}

func (t *TrackChaincode) DeleteTrack(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var trackId = args[0]

	data, e := t.repo.GetByKey(APIstub, trackId)

	if e != nil{
		return shim.Error("Failed to delete track " + trackId + " due to " + e.Error())
	}

	if string(data) == "" {
		return shim.Error("Failed to delete track, because the track " + trackId + " does not exist.")
	}

	err := t.repo.DeleteTrack(APIstub, trackId)

	if err != nil{
		return shim.Error("Failed to delete track " + trackId + " due to " + err.Error())
	}

	return shim.Success(nil)
}

func (t *TrackChaincode) UpdateTrack(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	var trackId = args[0]

	trk, err := t.repo.GetByKey(APIstub, trackId)

	if err != nil{
		return shim.Error("Failed to update track " + trackId + ": " + err.Error())
	}

	if string(trk) == "" {
		return shim.Error("Failed to update track, because the track does not exist.")
	}

	err = t.repo.UpdateTrack(APIstub, trackId, args[1], args[2])

	if err != nil{
		return shim.Error("Failed to update track " + trackId + ": " + err.Error())
	}

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(TrackChaincode))
	if err != nil {
		fmt.Printf("Error creating new TrackChaincode: %s", err)
	}
}
