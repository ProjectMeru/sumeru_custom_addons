package models

import (
	"sumeru/core/sdk"
)

// EngagementProject is a client engagement / professional services project.
type EngagementProject struct {
	sdk.Model `sumeru:"model=engagement.project"`

	Name         sdk.String `sumeru:"required,unique,index,string=Engagement Name"`
	Reference    sdk.String `sumeru:"index,string=Reference"`
	ExternalCode sdk.String `sumeru:"column=external_ref,index,string=External Code"`
	SerialNumber sdk.String `sumeru:"unique,index,string=Serial Number"`
	ShortCode    sdk.String `sumeru:"size=32,string=Short Code"`
	Description  sdk.Text   `sumeru:"string=Description"`
	Notes        sdk.HTML   `sumeru:"string=Internal Notes"`
	Website      sdk.URL    `sumeru:"string=Website"`

	Active   sdk.Boolean `sumeru:"string=Active,default=true"`
	Verified sdk.Boolean `sumeru:"string=Verified,default=false"`
	Archived sdk.Boolean `sumeru:"string=Archived,default=false,index"`

	Sequence    sdk.Integer `sumeru:"string=Sequence,default=10"`
	Quantity    sdk.Integer `sumeru:"string=Quantity,default=1"`
	ProgressPct sdk.Integer `sumeru:"string=Progress %,default=0,min=0,max=100"`
	Rating      sdk.Integer `sumeru:"string=Rating,min=0,max=5"`
	Amount      sdk.Float   `sumeru:"string=Budget Amount,default=0"`
	Price       sdk.Numeric `sumeru:"string=Price,precision=18,scale=2,default=0"`
	TaxRate     sdk.Numeric `sumeru:"string=Tax Rate,precision=20,scale=8,default=0"`
	Balance     sdk.Numeric `sumeru:"string=Balance,precision=40,scale=18,default=0"`

	Subtotal   sdk.Money                  `sumeru:"string=Subtotal,currency=CurrencyID"`
	TaxAmount  sdk.Money                  `sumeru:"string=Tax,currency=CurrencyID"`
	CurrencyID sdk.Many2One[CoreCurrency] `sumeru:"string=Currency"`

	Priority  sdk.Selection[Priority]     `sumeru:"string=Priority,default=normal"`
	State     sdk.Selection[ProjectState] `sumeru:"required,string=Status,default=draft"`
	Kind      sdk.Selection[ProjectKind]    `sumeru:"string=Kind,default=engagement"`
	ColorCode sdk.String                  `sumeru:"string=Color"`

	DateStart   sdk.Date     `sumeru:"string=Start Date"`
	DateStop    sdk.Date     `sumeru:"string=End Date"`
	DatetimeDue sdk.DateTime `sumeru:"string=Due Date,index"`
	OpeningTime sdk.Time     `sumeru:"string=Opening Time"`
	ClosingTime sdk.Time     `sumeru:"string=Closing Time"`
	CompletedAt sdk.DateTime `sumeru:"string=Completed At"`

	ProcessingTime sdk.Duration `sumeru:"string=Processing Time"`

	Email    sdk.Email `sumeru:"string=Email"`
	Phone    sdk.Phone `sumeru:"string=Phone"`
	PublicID sdk.UUID  `sumeru:"string=Public ID,default=uuid,unique"`

	Settings     sdk.Json   `sumeru:"string=Settings"`
	Document     sdk.Binary `sumeru:"string=Document"`
	Avatar       sdk.Image  `sumeru:"string=Avatar"`
	MetadataJson sdk.Json   `sumeru:"string=Metadata"`

	CompanyID  sdk.Many2One[CoreCompany]      `sumeru:"string=Company"`
	UserID     sdk.Many2One[CoreUser]         `sumeru:"string=Account Manager"`
	EmployeeID sdk.Many2One[HrEmployee]       `sumeru:"string=Lead Consultant"`
	PartnerID  sdk.Many2One[CorePartner]      `sumeru:"string=Client"`
	CountryID  sdk.Many2One[CoreCountry]      `sumeru:"string=Country"`
	StateID    sdk.Many2One[CoreCountryState] `sumeru:"string=State"`
	CityID     sdk.Many2One[CoreCity]         `sumeru:"string=City"`
	ParentID   sdk.Many2One[EngagementProject] `sumeru:"column=parent_id,string=Parent Engagement,ondelete=set_null"`
	ChildIds   sdk.One2Many[EngagementProject] `sumeru:"inverse=ParentID,string=Sub-engagements"`

	TagIds          sdk.Many2Many[EngagementTag]         `sumeru:"table=engagement_tag_rel,left=project_id,right=tag_id,string=Tags"`
	ServiceLineIds  sdk.Many2Many[EngagementServiceLine] `sumeru:"table=engagement_service_line_rel,left=project_id,right=service_line_id,string=Service Lines"`
	DeliverableIds  sdk.One2Many[EngagementDeliverable]  `sumeru:"string=Deliverables"`
	SiteIds         sdk.One2Many[EngagementSite]         `sumeru:"string=Client Sites"`
	MilestoneIds    sdk.One2Many[EngagementMilestone]    `sumeru:"string=Milestones"`

	ResourceRef   sdk.Reference         `sumeru:"string=Resource Ref"`
	ResourceModel sdk.String            `sumeru:"string=Resource Model"`
	ResourceID    sdk.Many2OneReference `sumeru:"model_field=ResourceModel,string=Resource ID"`

	CompanyName     sdk.String  `sumeru:"related=company_id.name,string=Company Name"`
	ComputedAmount  sdk.Float   `sumeru:"compute=computed_amount,string=Computed Amount"`
	StoredLineCount sdk.Integer `sumeru:"compute=stored_line_count,store,readonly,string=Deliverable Count"`

	Salary sdk.Money `sumeru:"string=Consulting Rate,currency=CurrencyID,groups=engagement_cookbook.group_engagement_cookbook_manager"`
}
