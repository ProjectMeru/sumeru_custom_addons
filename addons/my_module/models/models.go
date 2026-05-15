// Package models contains the business domain objects for my_module.
package models

import (
	"context"
	"sumeru/core/base"
)

// MyModule represents the primary record for My Module.
type MyModule struct {
	base.BaseModel
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
func (m *MyModule) Fields() []base.FieldDefinition {
	return []base.FieldDefinition{
		{Name: "name", Type: base.Char, String: "Name", Required: true, Index: true},
		{Name: "description", Type: base.Text, String: "Description"},
		{Name: "active", Type: base.Boolean, String: "Active", DefaultVal: "true"},
		{Name: "sequence", Type: base.Integer, String: "Sequence", DefaultVal: "10"},
	}
}

// --- Standard CRUD Methods ---
// These methods demonstrate how to use the 'base' API while respecting security and context.

// CreateNew inserts a new record with the given values.
func (m *MyModule) CreateNew(ctx context.Context, vals map[string]interface{}) (int, error) {
	// base.Create handles RBAC/ABAC automatically based on the 'ctx' user.
	return base.Create(ctx, base.CreateInput{
		Model:  m,
		Values: vals,
	})
}

// FetchAll retrieves all active records.
func (m *MyModule) FetchAll(ctx context.Context) ([]map[string]interface{}, error) {
	return base.Search(ctx, base.SearchInput{
		ModelName: m.ModelName(),
		Domain:    [][]interface{}{ {"active", "=", true} },
	})
}

// UpdateRecord modifies an existing record.
func (m *MyModule) UpdateRecord(ctx context.Context, id int, vals map[string]interface{}) (int, error) {
	vals["id"] = id
	return base.Upsert(ctx, base.UpsertInput{
		Model:       m,
		Values:      vals,
		ConflictCol: "id",
	})
}

func init() {
	// Register the model in the global Sumeru registry (Module must match manifest.json "name").
	base.RegisterModel(base.RegisterModelInput{Model: &MyModule{}, Module: "my_module"})
}
