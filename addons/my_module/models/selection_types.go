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

// Kind distinguishes cookbook rows from seeded demo rows.
type Kind string

const (
	KindRecord Kind = "record"
	KindDemo   Kind = "demo"
)

// EventState is the workflow for my.module.event.
type EventState string

const (
	EventDraft   EventState = "draft"
	EventPlanned EventState = "planned"
	EventDone    EventState = "done"
)
