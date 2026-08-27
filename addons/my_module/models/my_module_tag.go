package models

import (
	"sumeru/core/sdk"
)

type MyModuleTag struct {
	sdk.Model `sumeru:"model=my.module.tag"`

	Name     sdk.String  `sumeru:"required,unique,index,string=Name"`
	Color    sdk.String  `sumeru:"string=Color"`
	Sequence sdk.Integer `sumeru:"string=Sequence,default=10"`
}
