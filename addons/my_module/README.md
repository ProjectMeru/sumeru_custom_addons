# my_module — Developer UI/ORM Cookbook

Reference addon for Sumeru custom module authors. The `my.module` form demonstrates every supported SDK field marker, struct tag, relation pattern, compute, and related field.

## Models

| Model | Purpose |
|-------|---------|
| `my.module` | Full field-type cookbook |
| `my.module.line` | One2Many child lines |
| `my.module.tag` | Many2Many tags |
| `my.module.category` | Many2Many with explicit rel table |

## Cross-module relations

Depends on **`base`** and **`hr`**. Run from workspace root:

```bash
make generate
```

This generates **`models/zrefs.go`** with typed handles (`CoreCompany`, `HrEmployee`, `CorePartner`, …). Use them in relation fields — no engine edits, no `sdk.CoreCompany` markers:

```go
CompanyID  sdk.Many2One[CoreCompany]
EmployeeID sdk.Many2One[HrEmployee]
PartnerID  sdk.Many2One[CorePartner]   // Partner record type is basemodels.Partner
```

Add `"salary"` (or any module) to `manifest.json` `depends`, then `make generate` again.

## Field markers demonstrated

| Marker | Example field | Notes |
|--------|---------------|-------|
| `sdk.String` | `name`, `reference` | `column=`, `unique`, `size=` |
| `sdk.Text` | `description` | |
| `sdk.HTML` | `notes` | widget `html` |
| `sdk.URL` | `website` | widget `url` |
| `sdk.Boolean` | `active`, `verified` | |
| `sdk.Integer` | `quantity`, `rating` | `min=`, `max=` |
| `sdk.Float` | `amount` | |
| `sdk.Numeric` | `price`, `balance` | `precision=`, `scale=` |
| `sdk.Money` | `subtotal`, `salary` | `currency=CurrencyID` |
| `sdk.Selection[T]` | `priority`, `state` | const options auto-discovered |
| `sdk.Date` / `DateTime` / `Time` | `date_start`, `datetime_due`, `opening_time` | |
| `sdk.Duration` | `processing_time` | |
| `sdk.Email` / `Phone` | `email`, `phone` | auto widgets |
| `sdk.UUID` | `public_id` | `default=uuid` |
| `sdk.Json` / `Binary` / `Image` | `settings`, `document`, `avatar` | |
| `sdk.Reference` | `resource_ref` | |
| `sdk.Many2OneReference` | `resource_id` | `model_field=ResourceModel` |
| `sdk.Many2One` / `One2Many` / `Many2Many` | relations | `ondelete=`, `inverse=`, `table=` |
| related | `company_name` | `related=company_id.name` |
| compute | `computed_amount`, `stored_line_count` | handlers in `computed.go` |

## Object actions

- `action_confirm` → state `confirmed`
- `action_done` → state `done`
- `action_reset_draft` → state `draft`

## Reload after changes

```bash
make generate
make update MODULES='my_module'
make run
```

## Conventions

- **Selection:** `type Foo string` + `const (...)` — no `init()` registration.
- **Same-module relations:** `Many2One[MyModule]`, `One2Many[MyModuleLine]`.
- **Cross-module relations:** declare `depends`, run `make generate`, use types from `zrefs.go`.
- **Optional soft deps:** `Many2One[sdk.Any]` + `comodel=technical.model` when the module may not be installed.
- Register compute handlers in `computed.go` via `orm.RegisterCompute` (names must match `compute=` tags).
- List views live in `views/list_view.xml` with `type="list"`.
