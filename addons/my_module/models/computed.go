
package models

import (
	"context"

	"sumeru/core/orm"
)

func init() {
	orm.RegisterCompute("my.module", "computed_amount", []string{"amount", "quantity"}, computeAmount)
	orm.RegisterCompute("my.module", "stored_line_count", []string{"line_ids"}, computeLineCount)
}

func computeAmount(_ context.Context, rec map[string]interface{}) (interface{}, error) {
	amount := asFloat(rec["amount"])
	qty := asFloat(rec["quantity"])
	if qty == 0 {
		qty = 1
	}
	return amount * qty, nil
}

func computeLineCount(ctx context.Context, rec map[string]interface{}) (interface{}, error) {
	id, ok := orm.CoerceInt64(rec["id"])
	if !ok || id <= 0 {
		return 0, nil
	}
	lines, err := orm.Search(ctx, "my.module.line", [][]interface{}{
		{"module_id", "=", int(id)},
	})
	if err != nil {
		return nil, err
	}
	return len(lines), nil
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
