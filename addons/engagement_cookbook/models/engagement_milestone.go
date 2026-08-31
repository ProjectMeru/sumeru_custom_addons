package models

import (
	"sumeru/core/sdk"
)

type EngagementMilestone struct {
	sdk.Model `sumeru:"model=engagement.milestone"`

	Name       sdk.String                    `sumeru:"required,index,string=Milestone"`
	ProjectID  sdk.Many2One[EngagementProject] `sumeru:"required,string=Engagement,ondelete=cascade"`
	UserID     sdk.Many2One[CoreUser]        `sumeru:"string=Organizer"`
	EmployeeID sdk.Many2One[HrEmployee]      `sumeru:"string=Owner"`
	DateStart  sdk.Date                      `sumeru:"required,index,string=Start Date"`
	DateStop   sdk.Date                      `sumeru:"string=End Date"`
	State      sdk.Selection[MilestoneState] `sumeru:"required,string=Status,default=draft"`
}
