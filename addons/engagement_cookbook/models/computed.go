package models

import (
	"context"

	"sumeru/core/orm"
)

func init() {
	orm.RegisterCompute("engagement.project", "computed_amount", []string{"amount", "quantity"}, computeProjectAmount)
	orm.RegisterCompute("engagement.project", "stored_line_count", []string{"deliverable_ids"}, computeDeliverableCount)
	orm.RegisterCompute("engagement.deliverable", "line_total", []string{"quantity", "unit_price"}, computeLineTotal)
}

func computeProjectAmount(_ context.Context, rec map[string]interface{}) (interface{}, error) {
	amount := asFloat(rec["amount"])
	qty := asFloat(rec["quantity"])
	if qty == 0 {
		qty = 1
	}
	return amount * qty, nil
}

func computeDeliverableCount(ctx context.Context, rec map[string]interface{}) (interface{}, error) {
	id, ok := orm.CoerceInt64(rec["id"])
	if !ok || id <= 0 {
		return 0, nil
	}
	lines, err := orm.Search(ctx, "engagement.deliverable", [][]interface{}{
		{"project_id", "=", int(id)},
	})
	if err != nil {
		return nil, err
	}
	return len(lines), nil
}

func computeLineTotal(_ context.Context, rec map[string]interface{}) (interface{}, error) {
	qty := asFloat(rec["quantity"])
	if qty == 0 {
		qty = 1
	}
	return qty * asFloat(rec["unit_price"]), nil
}

func asFloat(v interface{}) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case float32:
		return float64(t)
	case int:
		return float64(t)
	case int64:
		return float64(t)
	default:
		return 0
	}
}
