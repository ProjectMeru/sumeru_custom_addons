package models

import (
	"sumeru/core/sdk"
)

type MyModuleCategory struct {
	sdk.Model `sumeru:"model=my.module.category"`

	Name     sdk.String  `sumeru:"required,unique,index,string=Name"`
	Sequence sdk.Integer `sumeru:"string=Sequence,default=10"`
	Active   sdk.Boolean `sumeru:"string=Active,default=true"`
}
