package main

import (
	"fmt"
	"bytes"
	"strings"
	"encoding/json"
	"github.com/beevik/guid"
	"github.com/skillbill/models"
	"github.com/skillbill/packages/repository"
	"github.com/hyperledger/fabric/core/chaincode/shim"

	logs "github.com/skillbill/packages/logs"
)

type TrackRepo struct {
	repo	repository.BaseRepo
}

func InitRepo() ITrackRepo {
	return TrackRepo{}
}

func (t TrackRepo) GetByQuery(APIstub shim.ChaincodeStubInterface,  args []string) 	([]byte, error){
	var b bytes.Buffer
	i := 0
	for i < len(args) {

		var params = strings.Split(args[i], ",")
		b.WriteString("\""+params[0]+"\":\""+ params[1]+"\"")
		i = i + 1
		if i != len(args){
			b.WriteString(",");
		}
	}

	var query = `{"selector":{`+ b.String() +`}}`

	logs.LogInfo("Query string: " + query)

	return t.repo.GetByQuery(APIstub, query)
}

func (t TrackRepo) GetByKey(APIstub shim.ChaincodeStubInterface, key string) ([]byte, error) {
	value, err := t.repo.GetByKey(APIstub, key)
	
	if err != nil {
		return nil, fmt.Errorf("Failed to get record %s", key) 
	}

	return value, nil
}

func (t TrackRepo) CreateTrack(APIstub shim.ChaincodeStubInterface, trackTranslation string, version string) (string, error) {
	var track = models.Track{TrackID: guid.New().StringUpper(), TrackTranslationID: trackTranslation, Version: version, DocType: "track"}

	data, _ := json.Marshal(track)

	err := t.repo.Save(APIstub, track.TrackID, data)

	return track.TrackID, err
}

func (t TrackRepo) UpdateTrack(APIstub shim.ChaincodeStubInterface, trackID string, trackTranslation string, version string) error {
	value, err := t.GetByKey(APIstub, trackID)

	if err != nil{
		return fmt.Errorf("Failed to update the track %s: %s", trackID, err.Error())
	}

	if len(string(value)) == 0 {
		return fmt.Errorf("Failed to update knowledge group, because the knowledge group does not exist.")
	}

	track := models.Track{}
	json.Unmarshal(value, &track)
	track.TrackTranslationID = trackTranslation
	track.Version  = version
	value, _ = json.Marshal(track)

	return t.repo.Save(APIstub, trackID, value)
}

func (t TrackRepo) DeleteTrack(APIstub shim.ChaincodeStubInterface, key string) error {
	return t.repo.Delete(APIstub, key)
}