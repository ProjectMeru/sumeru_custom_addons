# my_module — Developer UI/ORM Cookbook

Reference addon for Sumeru custom module authors. Open the form, list, and kanban views to copy patterns for fields, widgets, and layout tags.

## Models

| Model | Purpose |
|-------|---------|
| `my.module` | All 13 ORM field types + relations |
| `my.module.line` | One2Many child lines |
| `my.module.tag` | Many2Many tags |

## Field types on `my.module`

| Field | Type | Widget / notes |
|-------|------|----------------|
| `name` | Char | Title in form h1 |
| `reference` | Char | |
| `description` | Text | Notebook tab |
| `active` | Boolean | List toggle; form radio |
| `sequence` | Integer | List handle |
| `quantity`, `progress_pct` | Integer | |
| `amount` | Float | |
| `price` | Numeric | |
| `priority`, `state` | Selection | Statusbar on `state` |
| `date_start` | Date | |
| `datetime_due` | DateTime | |
| `email` | Char | `widget="email"` |
| `phone` | Char | `widget="phone"` |
| `image` | Text | `widget="image"` |
| `metadata_json` | Json | Raw JSON |
| `company_id`, `user_id` | Many2One | `widget="selection"` |
| `country_id`, `state_id`, `city_id` | Many2One | Geo cascade |
| `tag_ids` | Many2Many | `widget="many2many_tags"` |
| `line_ids` | One2Many | Embedded lines table |

## Object actions

- `action_confirm` → state `confirmed`
- `action_done` → state `done`
- `action_reset_draft` → state `draft`

## Reload after changes

```bash
make build
./bin/sumeru-erp -c sumeru.conf -u my_module -stop-after-init=true
make run
```

Or from this workspace: `make update MODULES='my_module'`

## Conventions

- List views live in `views/list_view.xml` with `type="list"` and id suffix `_list`.
- XML `id` values must be **globally unique** across all installed modules (see student vs my_module menu collision fix).
