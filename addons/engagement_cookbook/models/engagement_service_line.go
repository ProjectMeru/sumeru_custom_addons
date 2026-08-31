package models

import (
	"sumeru/core/sdk"
)

type EngagementServiceLine struct {
	sdk.Model `sumeru:"model=engagement.service.line"`

	Name        sdk.String  `sumeru:"required,unique,index,string=Service Line"`
	Sequence    sdk.Integer `sumeru:"string=Sequence,default=10"`
	Active      sdk.Boolean `sumeru:"string=Active,default=true"`
	Description sdk.Text    `sumeru:"string=Description"`
}
