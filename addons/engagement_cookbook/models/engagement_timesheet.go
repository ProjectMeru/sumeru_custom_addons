package models

import (
	"sumeru/core/sdk"
)

type EngagementTimesheet struct {
	sdk.Model `sumeru:"model=engagement.timesheet"`

	ProjectID   sdk.Many2One[EngagementProject]      `sumeru:"required,index,string=Engagement,ondelete=cascade"`
	EmployeeID  sdk.Many2One[HrEmployee]              `sumeru:"string=Consultant"`
	Date        sdk.Date                              `sumeru:"required,index,string=Date"`
	Hours       sdk.Float                             `sumeru:"string=Hours,default=0"`
	Billable    sdk.Boolean                           `sumeru:"string=Billable,default=true"`
	Category    sdk.Selection[TimesheetCategory]      `sumeru:"string=Category,default=dev"`
	Rate        sdk.Numeric                           `sumeru:"string=Hourly Rate,precision=18,scale=2,default=0"`
	Description sdk.Text                              `sumeru:"string=Description"`
}
