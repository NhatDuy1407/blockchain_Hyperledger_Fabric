package main


const (
	Planned 			string = "planned"
	InProgress			string = "inprogress"
	Completed			string = "completed"
	AssessmentRequest	string = "assessmentrequest"
)

type SkillPlanPlannedSkill struct {
	ID				string	`json:"id"`
	PlannedFrom		string	`json:"plannedfrom"`
	PlannedTo		string	`json:"plannedto"`
	Priority		string	`json:"priority"`
	SkillID			string	`json:"skillid"`
	UserID			string	`json:"userid"`
	DocType			string	`json:"doctype"`
}

type SkillPlanInProgressSkill struct {
	ID					string	`json:"id"`
	SkillACID			string	`json:"skillacid"`
	SkillACStartdate	string	`json:"skillacstartdate"`
	SkillID				string	`json:"skillid"`
	UserID				string	`json:"userid"`
	DocType				string	`json:"doctype"`
}

type SkillPlanCompletedSkill struct {
	ID				string	`json:"id"`
	AccessedBy		string	`json:"assessedby"`
	CompletedOn		string 	`json:"completedon"`
	SkillID			string	`json:"skillid"`
	UserID			string	`json:"userid"`
	DocType			string	`json:"doctype"`
}

type SkillPlanAssessmentRequest struct {
	ID				string	`json:"id"`
	AssesseeID		string	`json:"assesseeid"`
	AssessorID		string	`json:"assessorid"`
	SkillID			string	`json:"skillid"`
	DocType			string	`json:"doctype"`
}
