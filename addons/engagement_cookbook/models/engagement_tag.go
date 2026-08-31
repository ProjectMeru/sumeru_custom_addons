package models

import (
	"sumeru/core/sdk"
)

type EngagementTag struct {
	sdk.Model `sumeru:"model=engagement.tag"`

	Name        sdk.String `sumeru:"required,unique,index,string=Tag"`
	Color       sdk.String `sumeru:"string=Color"`
	Sequence    sdk.Integer `sumeru:"string=Sequence,default=10"`
	Description sdk.Text   `sumeru:"string=Description"`
}
