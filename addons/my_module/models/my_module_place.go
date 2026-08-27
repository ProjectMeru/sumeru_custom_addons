package models

import (
	"sumeru/core/sdk"
)

type MyModulePlace struct {
	sdk.Model `sumeru:"model=my.module.place"`

	Name      sdk.String                `sumeru:"required,unique,index,string=Name"`
	Latitude  sdk.Float                 `sumeru:"string=Latitude"`
	Longitude sdk.Float                 `sumeru:"string=Longitude"`
	Color     sdk.String                `sumeru:"string=Color"`
	PartnerID sdk.Many2One[CorePartner] `sumeru:"string=Partner"`
	Active    sdk.Boolean               `sumeru:"string=Active,default=true"`
}
