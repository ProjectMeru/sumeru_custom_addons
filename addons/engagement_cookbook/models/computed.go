package models

import (
	"context"

	"sumeru/core/sdk"
)

func computed_amount(_ context.Context, rec map[string]interface{}) (interface{}, error) {
	amount := asFloat(rec["amount"])
	qty := asFloat(rec["quantity"])
	if 0 == qty {
		qty = 1
	}
	return amount * qty, nil
}

func stored_line_count(ctx context.Context, rec map[string]interface{}) (interface{}, error) {
	id, ok := sdk.CoerceInt64(rec["id"])
	if !ok || 0 >= id {
		return 0, nil
	}
	lines, err := sdk.Search(ctx, "engagement.deliverable", [][]interface{}{
		{"project_id", "=", int(id)},
	})
	if err != nil {
		return nil, err
	}
	return len(lines), nil
}

func line_total(_ context.Context, rec map[string]interface{}) (interface{}, error) {
	qty := asFloat(rec["quantity"])
	if 0 == qty {
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
