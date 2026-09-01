package models

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityNormal Priority = "normal"
	PriorityHigh   Priority = "high"
)

type ProjectState string

const (
	ProjectStateDraft  ProjectState = "draft"
	ProjectStateActive ProjectState = "active"
	ProjectStateDone   ProjectState = "done"
)

type ProjectKind string

const (
	ProjectKindEngagement ProjectKind = "engagement"
	ProjectKindDemo       ProjectKind = "demo"
)

type MilestoneState string

const (
	MilestoneDraft   MilestoneState = "draft"
	MilestonePlanned MilestoneState = "planned"
	MilestoneDone    MilestoneState = "done"
)

type EngagementTier string

const (
	TierStandard  EngagementTier = "standard"
	TierPremium   EngagementTier = "premium"
	TierStrategic EngagementTier = "strategic"
)

type TimesheetCategory string

const (
	TimesheetCategoryDev     TimesheetCategory = "dev"
	TimesheetCategoryQA      TimesheetCategory = "qa"
	TimesheetCategoryPM      TimesheetCategory = "pm"
	TimesheetCategorySupport TimesheetCategory = "support"
)
