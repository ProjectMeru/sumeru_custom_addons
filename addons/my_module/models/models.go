// Package models contains the business domain objects for my_module.
package models

import (
	"sumeru/core/sdk"
)

// MyModule represents the primary record for My Module.
type MyModule struct {
	sdk.BaseModel
	Name        string `db:"name"`
	Description string `db:"description"`
	Active      bool   `db:"active"`
	Sequence    int    `db:"sequence"`
}

// ModelName returns the technical name of the model.
func (m *MyModule) ModelName() string {
	return "my.module"
}

// Fields returns the field definitions for the ORM and UI.
func (m *MyModule) Fields() []sdk.FieldDefinition {
	return []sdk.FieldDefinition{
		{Name: "name", Type: sdk.Char, String: "Name", Required: true, Index: true},
		{Name: "description", Type: sdk.Text, String: "Description"},
		{Name: "active", Type: sdk.Boolean, String: "Active", DefaultVal: "true"},
		{Name: "sequence", Type: sdk.Integer, String: "Sequence", DefaultVal: "10"},
	}
}

func init() {
	// Register the model in the global Sumeru registry (Module must match manifest.json "name").
	sdk.RegisterModel(sdk.RegisterModelInput{Model: &MyModule{}, Module: "my_module"})
}
