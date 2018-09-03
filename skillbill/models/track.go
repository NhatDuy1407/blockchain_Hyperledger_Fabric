package models

type Track struct {
	TrackID					string 	`json:"trackid"`
	TrackTranslationID		string	`json:"tracktranslationid"`
	Version					string	`json:"version"`
	DocType					string	`json:"doctype"`
}