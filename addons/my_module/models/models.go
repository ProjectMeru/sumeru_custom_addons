package models

import (
	"sumeru/core/sdk"
)

type MyModule struct {
	sdk.Model `sumeru:"model=my.module"`

	// String / text
	Name         sdk.String `sumeru:"required,index,string=Name"`
	Reference    sdk.String `sumeru:"index,string=Reference"`
	ExternalCode sdk.String `sumeru:"column=external_ref,index,string=External Code"`
	SerialNumber sdk.String `sumeru:"unique,index,string=Serial Number"`
	ShortCode    sdk.String `sumeru:"size=32,string=Short Code"`
	Description  sdk.Text   `sumeru:"string=Description"`
	Notes        sdk.HTML   `sumeru:"string=Notes"`
	Website      sdk.URL    `sumeru:"string=Website"`

	// Boolean
	Active   sdk.Boolean `sumeru:"string=Active,default=true"`
	Verified sdk.Boolean `sumeru:"string=Verified,default=false"`
	Archived sdk.Boolean `sumeru:"string=Archived,default=false,index"`

	// Integer / numeric
	Sequence    sdk.Integer `sumeru:"string=Sequence,default=10"`
	Quantity    sdk.Integer `sumeru:"string=Quantity,default=1"`
	ProgressPct sdk.Integer `sumeru:"string=Progress %,default=0,min=0,max=100"`
	Rating      sdk.Integer `sumeru:"string=Rating,min=0,max=5"`
	Amount      sdk.Float   `sumeru:"string=Amount,default=0"`
	Price       sdk.Numeric `sumeru:"string=Price,precision=18,scale=2,default=0"`
	TaxRate     sdk.Numeric `sumeru:"string=Tax Rate,precision=20,scale=8,default=0"`
	Balance     sdk.Numeric `sumeru:"string=Balance,precision=40,scale=18,default=0"`

	// Money
	Subtotal   sdk.Money          `sumeru:"string=Subtotal,currency=CurrencyID"`
	TaxAmount  sdk.Money          `sumeru:"string=Tax,currency=CurrencyID"`
	CurrencyID sdk.Many2One[CoreCurrency] `sumeru:"string=Currency"`

	// Selection — options discovered from Priority/State consts (no init registration)
	Priority sdk.Selection[Priority] `sumeru:"string=Priority,default=normal"`
	State    sdk.Selection[State]    `sumeru:"required,string=Status,default=draft"`

	// Date / time
	DateStart   sdk.Date     `sumeru:"string=Start Date"`
	DatetimeDue sdk.DateTime `sumeru:"string=Due Date,index"`
	OpeningTime sdk.Time     `sumeru:"string=Opening Time"`
	ClosingTime sdk.Time     `sumeru:"string=Closing Time"`
	CompletedAt sdk.DateTime `sumeru:"string=Completed At"`

	// Duration
	ProcessingTime sdk.Duration `sumeru:"string=Processing Time"`

	// Contact / identity
	Email    sdk.Email `sumeru:"string=Email"`
	Phone    sdk.Phone `sumeru:"string=Phone"`
	PublicID sdk.UUID  `sumeru:"string=Public ID,default=uuid,unique"`

	// Media / structured
	Settings     sdk.Json   `sumeru:"string=Settings"`
	Document     sdk.Binary `sumeru:"string=Document"`
	Avatar       sdk.Image  `sumeru:"string=Avatar"`
	MetadataJson sdk.Json   `sumeru:"string=Metadata"`

	// Relations — same-module generics; cross-module types from generated zrefs.go (make generate)
	CompanyID  sdk.Many2One[CoreCompany]        `sumeru:"string=Company"`
	UserID     sdk.Many2One[CoreUser]            `sumeru:"string=Responsible"`
	EmployeeID sdk.Many2One[HrEmployee]          `sumeru:"string=Employee"`
	PartnerID  sdk.Many2One[CorePartner]         `sumeru:"string=Partner"`
	CountryID  sdk.Many2One[CoreCountry]         `sumeru:"string=Country"`
	StateID    sdk.Many2One[CoreCountryState]    `sumeru:"string=State"`
	CityID     sdk.Many2One[CoreCity]            `sumeru:"string=City"`
	ParentID  sdk.Many2One[MyModule]              `sumeru:"column=parent_id,string=Parent,ondelete=set_null"`
	ChildIds  sdk.One2Many[MyModule]              `sumeru:"inverse=ParentID,string=Children"`

	TagIds      sdk.Many2Many[MyModuleTag]      `sumeru:"table=my_module_tag_rel,left=module_id,right=tag_id,string=Tags"`
	CategoryIds sdk.Many2Many[MyModuleCategory] `sumeru:"table=my_module_category_rel,left=my_module_id,right=category_id,string=Categories"`
	LineIds     sdk.One2Many[MyModuleLine]      `sumeru:"string=Lines"`

	// Reference
	ResourceRef   sdk.Reference         `sumeru:"string=Resource Ref"`
	ResourceModel sdk.String            `sumeru:"string=Resource Model"`
	ResourceID    sdk.Many2OneReference `sumeru:"model_field=ResourceModel,string=Resource ID"`

	// Related / computed
	CompanyName     sdk.String  `sumeru:"related=company_id.name,string=Company Name"`
	ComputedAmount  sdk.Float   `sumeru:"compute=computed_amount,string=Computed Amount"`
	StoredLineCount sdk.Integer `sumeru:"compute=stored_line_count,store,readonly,string=Line Count"`

	// Security demo
	Salary sdk.Money `sumeru:"string=Salary,currency=CurrencyID,groups=my_module.group_my_module_manager"`
}
