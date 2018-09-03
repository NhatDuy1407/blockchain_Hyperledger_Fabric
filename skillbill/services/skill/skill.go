package main

const (
	Seflstudy	string = "Seflstudy"
	Assessment	string = "Assessment"

	FileLink	string = "FileLink"
	WebLink		string = "WebLink"
)

// Skill model
type Skill struct {
	AssessmentType			 string	`json:"assessmenttype"`
	DocType                  string `json:"doctype"`
	BackwardCompatibleTo     string	`json:"backwardcompatibleto"`
	DescriptionTranslationID string	`json:"descriptiontranslationid"`
	ImageID                  string	`json:"imageid"`
	KnowledgeGroupID         string	`json:"knowledgegroupid"`
	Level                    string	`json:"level"`
	NameTranslationID        string	`json:"nametranslationid"`
	SkillID                  string	`json:"skillid"`
	TimeEstimationInHours    string	`json:"timeestimationinhours"`		
	Version					 string `json:"version"`	
}

type SkillResource struct {
	ID							string `json:"id"`
	DocType						string `json:"doctype"`
	ResourceLink				string `json:"resourcelink"`
	ResourceTranslationID		string `json:"resourcetranslationid"`
	SkillID                  	string `json:"skillid"`
}

type SkillDenpendency struct {
	ID							string `json:"id"`
	DocType						string `json:"doctype"`
	DependingOnSkill			string `json:"denpendingonskill"`
	SkillID                  	string `json:"skillid"`
}

type SkillAcceptanceCriteria struct {
	ID							string `json:"id"`
	DocType						string `json:"doctype"`
	DescriptionTranslationID	string `json:"descriptiontranslationid"`
	SkillACID					string `json:"skillacid"`
	SkillID						string `json:"skillid"`
}
