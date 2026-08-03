// Package models contains the business domain objects for student.
package models

import (
	"sumeru/core/sdk"
)

// Student represents the primary record for Student Management.
type Student struct {
	sdk.BaseModel
	Name        string `db:"name"`
	Email       string `db:"email"`
	Active      bool   `db:"active"`
	Phone       int    `db:"phone"`
}

// ModelName returns the technical name of the model.
func (m *Student) ModelName() string {
	return "student"
}

// Fields returns the field definitions for the ORM and UI.
func (m *Student) Fields() []sdk.FieldDefinition {
	return []sdk.FieldDefinition{
		{Name: "name", Type: sdk.Char, String: "Name", Required: true, Index: true},
		{Name: "email", Type: sdk.Char, String: "Email", Required: true, Index: true},
		{Name: "active", Type: sdk.Boolean, String: "Active", DefaultVal: "true"},
		{Name: "phone", Type: sdk.Integer, String: "Phone", Required: true},
	}
}

func init() {
	// Register the model in the global Sumeru registry (Module must match manifest.json "name").
	sdk.RegisterModel(sdk.RegisterModelInput{Model: &Student{}, Module: "student"})
}
