package engagement_cookbook

import (
	"context"

	"sumeru/core/event"
	"sumeru/core/orm"
)

func init() {
	event.Subscribe("record.updated", onProjectDoneUpdatePartner)
}

func onProjectDoneUpdatePartner(ctx context.Context, ev event.Event) error {
	model, _ := ev.Payload["model"].(string)
	if model != "engagement.project" {
		return nil
	}
	id, ok := coerceEventID(ev.Payload["id"])
	if !ok {
		return nil
	}
	bypass := orm.ContextWithBypass(ctx, true)
	project, err := orm.SearchOne(bypass, "engagement.project", map[string]interface{}{"id": id})
	if err != nil {
		return err
	}
	if orm.AsString(project["state"]) != "done" {
		return nil
	}
	partnerID, _ := orm.CoerceInt64(project["partner_id"])
	if partnerID <= 0 {
		return nil
	}
	partner, err := orm.SearchOne(bypass, "core.partner", map[string]interface{}{"id": partnerID})
	if err != nil {
		return err
	}
	if orm.AsString(partner["engagement_tier"]) == "strategic" {
		return nil
	}
	return orm.UpdateRecordByID(bypass, "core.partner", int(partnerID), map[string]interface{}{
		"engagement_tier":   "premium",
		"strategic_account": true,
	})
}

func coerceEventID(v interface{}) (int, bool) {
	n, ok := orm.CoerceInt64(v)
	return int(n), ok && n > 0
}
