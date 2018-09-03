package models

type KnowledgeGroup struct {
	GroupID		string	`json:"groupid"`
	GroupName	string	`json:"groupname"`
	DocType		string	`json:"doctype"`
}

type KnowledgeGroupMember struct {
	ID			string	`json:"id"`
	GroupID		string	`json:"groupid"`
	MemberType	string	`json:"membertype"`
	UserID		string	`json:"userid"`
	DocType		string	`json:"doctype"`
}