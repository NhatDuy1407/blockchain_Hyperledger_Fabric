package main

import (
	"github.com/skillbill/services/core"
)


type Feature struct{
	FeatureID		string	`json:"featureid"`
	FeatureName		string	`json:"featurename"`
	DocType			string	`json:"doctype"`
	repo 			core.EntityChaincode 
}