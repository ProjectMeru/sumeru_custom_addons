// Package models contains the business domain objects for my_module.
package models

import (
	"sumeru/core/sdk"
)

// MyModule is the primary showcase record exercising all ORM field types.
type MyModule struct {
	sdk.BaseModel
}

// ModelName returns the technical name of the model.
func (m *MyModule) ModelName() string {
	return "my.module"
}

// Fields returns the field definitions for the ORM and UI.
func (m *MyModule) Fields() []sdk.FieldDefinition {
	return []sdk.FieldDefinition{
		{Name: "name", Type: sdk.Char, String: "Name", Required: true, Index: true},
		{Name: "reference", Type: sdk.Char, String: "Reference", Index: true},
		{Name: "description", Type: sdk.Text, String: "Description"},
		{Name: "active", Type: sdk.Boolean, String: "Active", DefaultVal: "true"},
		{Name: "sequence", Type: sdk.Integer, String: "Sequence", DefaultVal: "10"},
		{Name: "quantity", Type: sdk.Integer, String: "Quantity", DefaultVal: "1"},
		{Name: "progress_pct", Type: sdk.Integer, String: "Progress %", DefaultVal: "0"},
		{Name: "amount", Type: sdk.Float, String: "Amount", DefaultVal: 0},
		{Name: "price", Type: sdk.Numeric, String: "Price", DefaultVal: 0},
		{Name: "priority", Type: sdk.Selection, String: "Priority", DefaultVal: "normal", Selection: [][]string{
			{"low", "Low"},
			{"normal", "Normal"},
			{"high", "High"},
		}},
		{Name: "state", Type: sdk.Selection, String: "Status", DefaultVal: "draft", Selection: [][]string{
			{"draft", "Draft"},
			{"confirmed", "Confirmed"},
			{"done", "Done"},
		}},
		{Name: "date_start", Type: sdk.Date, String: "Start Date"},
		{Name: "datetime_due", Type: sdk.DateTime, String: "Due Date"},
		{Name: "email", Type: sdk.Char, String: "Email"},
		{Name: "phone", Type: sdk.Char, String: "Phone"},
		{Name: "image", Type: sdk.Text, String: "Image"},
		{Name: "metadata_json", Type: sdk.Json, String: "Metadata"},
		{Name: "company_id", Type: sdk.Many2One, Relation: "core.company", String: "Company"},
		{Name: "user_id", Type: sdk.Many2One, Relation: "core.user", String: "Responsible"},
		{Name: "country_id", Type: sdk.Many2One, Relation: "core.country", String: "Country"},
		{Name: "state_id", Type: sdk.Many2One, Relation: "core.country.state", String: "State"},
		{Name: "city_id", Type: sdk.Many2One, Relation: "core.city", String: "City"},
		{Name: "tag_ids", Type: sdk.Many2Many, Relation: "my.module.tag", RelationTable: "my_module_tag_rel", Column1: "module_id", Column2: "tag_id", String: "Tags"},
		{Name: "line_ids", Type: sdk.One2Many, Relation: "my.module.line", String: "Lines"},
	}
}

func init() {
	sdk.RegisterModel(sdk.RegisterModelInput{Model: &MyModule{}, Module: "my_module"})
}
