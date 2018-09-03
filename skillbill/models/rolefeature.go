package models

const (
	ReadOnly  int = 0
	ReadWrite int = 1
	Unknown   int = 99
)

type RoleFeature struct {
	ID          string `json:"id"`
	AccessLevel int    `json:"accesslevel"`
	FeatureID   string `json:"featureid"`
	RoleID      string `json:"roleid"`
	DocType     string `json:"doctype"`
}
