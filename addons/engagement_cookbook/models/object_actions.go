package models

import (
	"context"
	"fmt"

	"sumeru/core/orm"
)

func init() {
	orm.RegisterObjectAction("engagement.project", "action_confirm", actionActivateProject)
	orm.RegisterObjectAction("engagement.project", "action_done", actionDoneProject)
	orm.RegisterObjectAction("engagement.project", "action_reset_draft", actionResetDraftProject)

	orm.RegisterObjectAction("engagement.milestone", "action_plan", actionPlanMilestone)
	orm.RegisterObjectAction("engagement.milestone", "action_done", actionDoneMilestone)
	orm.RegisterObjectAction("engagement.milestone", "action_reset_draft", actionResetDraftMilestone)
}

func actionActivateProject(ctx context.Context, model string, id int, _ map[string]string) (string, error) {
	return setProjectState(ctx, model, id, "active")
}

func actionDoneProject(ctx context.Context, model string, id int, _ map[string]string) (string, error) {
	return setProjectState(ctx, model, id, "done")
}

func actionResetDraftProject(ctx context.Context, model string, id int, _ map[string]string) (string, error) {
	return setProjectState(ctx, model, id, "draft")
}

func setProjectState(ctx context.Context, model string, id int, state string) (string, error) {
	_ = model
	if err := orm.CheckModelAccess(ctx, orm.SecurityUID(ctx), "engagement.project", "write"); err != nil {
		return "", err
	}
	if _, err := orm.SearchOne(ctx, "engagement.project", map[string]interface{}{"id": id}); err != nil {
		return "", fmt.Errorf("record not found")
	}
	if err := orm.UpdateRecordByID(ctx, "engagement.project", id, map[string]interface{}{"state": state}); err != nil {
		return "", err
	}
	return "", nil
}

func actionPlanMilestone(ctx context.Context, model string, id int, _ map[string]string) (string, error) {
	return setMilestoneState(ctx, model, id, "planned")
}

func actionDoneMilestone(ctx context.Context, model string, id int, _ map[string]string) (string, error) {
	return setMilestoneState(ctx, model, id, "done")
}

func actionResetDraftMilestone(ctx context.Context, model string, id int, _ map[string]string) (string, error) {
	return setMilestoneState(ctx, model, id, "draft")
}

func setMilestoneState(ctx context.Context, model string, id int, state string) (string, error) {
	_ = model
	if err := orm.CheckModelAccess(ctx, orm.SecurityUID(ctx), "engagement.milestone", "write"); err != nil {
		return "", err
	}
	if _, err := orm.SearchOne(ctx, "engagement.milestone", map[string]interface{}{"id": id}); err != nil {
		return "", fmt.Errorf("record not found")
	}
	if err := orm.UpdateRecordByID(ctx, "engagement.milestone", id, map[string]interface{}{"state": state}); err != nil {
		return "", err
	}
	return "", nil
}
