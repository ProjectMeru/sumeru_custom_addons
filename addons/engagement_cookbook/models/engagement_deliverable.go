package models

import (
	"sumeru/core/sdk"
)

type EngagementDeliverable struct {
	sdk.Model `sumeru:"model=engagement.deliverable"`

	ProjectID  sdk.Many2One[EngagementProject] `sumeru:"required,index,string=Engagement,ondelete=cascade"`
	Name       sdk.String                      `sumeru:"required,string=Description"`
	Sequence   sdk.Integer                     `sumeru:"string=Sequence,default=10"`
	Quantity   sdk.Integer                     `sumeru:"string=Quantity,default=1"`
	UnitPrice  sdk.Numeric                     `sumeru:"string=Unit Price,default=0,precision=18,scale=2"`
	LineTotal  sdk.Numeric                     `sumeru:"compute=line_total,string=Line Total,readonly,precision=18,scale=2"`
	Note       sdk.Text                        `sumeru:"string=Note"`
}
