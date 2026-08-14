package models

import "sumeru/core/sdk"

// MyModuleLine is a One2Many child line on my.module.
type MyModuleLine struct {
	sdk.BaseModel
}

func (MyModuleLine) ModelName() string { return "my.module.line" }

func (MyModuleLine) Fields() []sdk.FieldDefinition {
	return []sdk.FieldDefinition{
		{Name: "module_id", Type: sdk.Many2One, Relation: "my.module", String: "Module", Required: true, Index: true},
		{Name: "name", Type: sdk.Char, String: "Description", Required: true},
		{Name: "quantity", Type: sdk.Integer, String: "Quantity", DefaultVal: "1"},
		{Name: "unit_price", Type: sdk.Numeric, String: "Unit Price", DefaultVal: 0},
		{Name: "note", Type: sdk.Text, String: "Note"},
	}
}

func init() {
	sdk.RegisterModel(sdk.RegisterModelInput{Model: &MyModuleLine{}, Module: "my_module"})
}
