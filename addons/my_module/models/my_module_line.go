package models

import (
	"sumeru/core/sdk"
)

type MyModuleLine struct {
	sdk.Model `sumeru:"model=my.module.line"`

	ModuleID  sdk.Many2One[MyModule] `sumeru:"required,index,string=Module,ondelete=cascade"`
	Name      sdk.String             `sumeru:"required,unique,string=Description"`
	Quantity  sdk.Integer            `sumeru:"string=Quantity,default=1"`
	UnitPrice sdk.Numeric            `sumeru:"string=Unit Price,default=0"`
	Note      sdk.Text               `sumeru:"string=Note"`
}
