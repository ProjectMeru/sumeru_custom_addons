package models

import "sumeru/core/sdk"

// MyModuleTag is a Many2Many tag on my.module.
type MyModuleTag struct {
	sdk.BaseModel
}

func (MyModuleTag) ModelName() string { return "my.module.tag" }

func (MyModuleTag) Fields() []sdk.FieldDefinition {
	return []sdk.FieldDefinition{
		{Name: "name", Type: sdk.Char, String: "Name", Required: true, Index: true, Unique: true},
		{Name: "color", Type: sdk.Selection, String: "Color", DefaultVal: "gray", Selection: [][]string{
			{"gray", "Gray"},
			{"blue", "Blue"},
			{"green", "Green"},
			{"orange", "Orange"},
			{"red", "Red"},
		}},
	}
}

func init() {
	sdk.RegisterModel(sdk.RegisterModelInput{Model: &MyModuleTag{}, Module: "my_module"})
}
