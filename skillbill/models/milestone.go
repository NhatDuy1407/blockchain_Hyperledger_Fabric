package models

type Milestone struct {
	DocType						string	`json:"doctype"`
	MilestoneID					string	`json:"milestoneid"`
	MilestoneTranslationID		string	`json:"milestonetranslationid"`
	TrackID						string 	`json:"trackid"`
	Version						string	`json:"version"`
}

type MilestoneDependency struct {
	ID						string	`json:"id"`
	DependingMilestone		string	`json:"dependingmilestone"`
	MilestoneID				string	`json:"milestoneid"`
	DocType					string	`json:"doctype"`
}

type MilestoneSkill struct {
	ID						string	`json:"id"`	
	MilestoneID				string	`json:"milestoneid"`
	SkillID					string	`json:"skillid"`
	DocType					string	`json:"doctype"`
}