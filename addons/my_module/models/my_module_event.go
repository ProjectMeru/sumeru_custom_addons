package models

import (
	"sumeru/core/sdk"
)

type MyModuleEvent struct {
	sdk.Model `sumeru:"model=my.module.event"`

	Name      sdk.String                `sumeru:"required,unique,index,string=Name"`
	ModuleID  sdk.Many2One[MyModule]    `sumeru:"string=Cookbook Record,ondelete=set_null"`
	UserID    sdk.Many2One[CoreUser]    `sumeru:"string=Responsible"`
	DateStart sdk.Date                  `sumeru:"required,index,string=Start Date"`
	DateStop  sdk.Date                  `sumeru:"string=Stop Date"`
	State     sdk.Selection[EventState] `sumeru:"required,string=Status,default=draft"`
}
