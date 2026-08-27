package models

import (
	"context"
	"fmt"

	"sumeru/core/orm"
)

func init() {
	orm.RegisterObjectAction("my.module", "action_confirm", actionConfirmMyModule)
	orm.RegisterObjectAction("my.module", "action_done", actionDoneMyModule)
	orm.RegisterObjectAction("my.module", "action_reset_draft", actionResetDraftMyModule)
}

func actionConfirmMyModule(ctx context.Context, model string, id int, _ map[string]string) (string, error) {
	return setMyModuleState(ctx, model, id, "confirmed")
}

func actionDoneMyModule(ctx context.Context, model string, id int, _ map[string]string) (string, error) {
	return setMyModuleState(ctx, model, id, "done")
}

func actionResetDraftMyModule(ctx context.Context, model string, id int, _ map[string]string) (string, error) {
	return setMyModuleState(ctx, model, id, "draft")
}

func setMyModuleState(ctx context.Context, model string, id int, state string) (string, error) {
	_ = model
	if err := orm.CheckModelAccess(ctx, orm.SecurityUID(ctx), "my.module", "write"); err != nil {
		return "", err
	}
	if _, err := orm.SearchOne(ctx, "my.module", map[string]interface{}{"id": id}); err != nil {
		return "", fmt.Errorf("record not found")
	}
	if err := orm.UpdateRecordByID(ctx, "my.module", id, map[string]interface{}{"state": state}); err != nil {
		return "", err
	}
	return "", nil
}
