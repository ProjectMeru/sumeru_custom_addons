package models

import (
	"sumeru/core/sdk"
)

type EngagementSite struct {
	sdk.Model `sumeru:"model=engagement.site"`

	Name      sdk.String                    `sumeru:"required,index,string=Site Name"`
	ProjectID sdk.Many2One[EngagementProject] `sumeru:"string=Engagement,ondelete=set_null"`
	Latitude  sdk.Float                     `sumeru:"string=Latitude"`
	Longitude sdk.Float                     `sumeru:"string=Longitude"`
	Color     sdk.String                    `sumeru:"string=Color"`
	PartnerID sdk.Many2One[CorePartner]     `sumeru:"string=Client Contact"`
	Active    sdk.Boolean                   `sumeru:"string=Active,default=true"`
}
