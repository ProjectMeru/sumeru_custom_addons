
package models

// Priority is a type-safe selection value for my.module.priority.
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityNormal Priority = "normal"
	PriorityHigh   Priority = "high"
)

// State is a type-safe workflow selection for my.module.state.
type State string

const (
	StateDraft     State = "draft"
	StateConfirmed State = "confirmed"
	StateDone      State = "done"
)
