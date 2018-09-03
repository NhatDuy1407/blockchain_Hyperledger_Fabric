package models

type User struct {
	ADLogin        string `json:"adlogin"`
	HashedPassword string `json:"hashedpassword"`
	PublicKey      string `json:"publickey"`
	RoleID         string `json:"roleid"`
	DocType        string `json:"doctype"`
}
