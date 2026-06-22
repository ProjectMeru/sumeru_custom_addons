// Package models contains the business domain objects for student.
package models

import (
	"sumeru/core/base"
)

// Student represents the primary record for Student Management.
type Student struct {
	base.BaseModel
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
func (m *Student) Fields() []base.FieldDefinition {
	return []base.FieldDefinition{
		{Name: "name", Type: base.Char, String: "Name", Required: true, Index: true},
		{Name: "email", Type: base.Char, String: "Email", Required: true, Index: true},
		{Name: "active", Type: base.Boolean, String: "Active", DefaultVal: "true"},
		{Name: "phone", Type: base.Integer, String: "Phone", Required: true},
	}
}

func init() {
	// Register the model in the global Sumeru registry (Module must match manifest.json "name").
	base.RegisterModel(base.RegisterModelInput{Model: &Student{}, Module: "student"})
}
