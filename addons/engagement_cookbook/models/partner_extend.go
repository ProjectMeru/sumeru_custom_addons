package models

import "sumeru/core/sdk"

// PartnerEngagement extends core.partner with engagement-specific CRM fields.
type PartnerEngagement struct {
	sdk.Model `sumeru:"inherit=core.partner"`

	EngagementTier   sdk.Selection[EngagementTier] `sumeru:"string=Engagement Tier,default=standard"`
	StrategicAccount sdk.Boolean                 `sumeru:"string=Strategic Account,default=false"`
}
