# my_module — kitchen-sink demo

Reference addon for custom-module authors. It shows every built-in field widget and workspace view type, split across a few models and a CRM-style menu tree.

## Models

| Model | Purpose |
|-------|---------|
| `my.module` | Field/widget cookbook (kanban, list, form, calendar, gantt, graph, pivot, cohort) |
| `my.module.line` | One2many child lines (also a standalone list/form) |
| `my.module.event` | Calendar + gantt planning (`date_start` / `date_stop`) |
| `my.module.place` | Map markers (`latitude` / `longitude`) |
| `my.module.tag` | Many2many tags + color picker |
| `my.module.category` | Many2many categories |

## Menus

| Menu | Action | View modes |
|------|--------|------------|
| Cookbook operations → Cookbook records | `action_my_module` | kanban, list, form, calendar, gantt |
| Cookbook operations → Cookbook lines | `action_my_module_line` | list, form |
| Cookbook planning → Cookbook events | `action_my_module_event` | calendar, gantt, list, form |
| Cookbook planning → Cookbook places | `action_my_module_place` | map, list, form |
| Cookbook reporting → Cookbook analysis | `action_my_module_analysis` | graph, pivot, cohort, list |
| Cookbook configuration → Tags / Categories | manager | list, form |

Search views attach to the collection toolbars (not a separate menu).

## Widgets on the cookbook form

`char` / `default`, `email`, `integer`, `float`, `numeric`, `text`, `json`, `date`, `datetime`, `boolean`, `boolean_toggle`, `radio`, `phone`, `url`, `html`, `binary`, `image`, `selection`, `many2one`, `many2many` / `many2many_tags`, `one2many`, `statusbar`, `priority`, `monetary`, `progress` / `progressbar`, `color`, `reference`, `handle` (list sequence).

## Cross-module relations

Depends on **`base`**, **`hr`**, and **`mail`** (chatter). From the custom-addons workspace:

```bash
make generate
make update MODULES='my_module'
make run
```

`make generate` rewrites **`models/zmodels.go`** and **`models/zrefs.go`**. Do not edit those files.

Demo rows in `data/demo.xml` upsert on unique `name` (and unique line descriptions).

## Object actions

- `action_confirm` → `confirmed`
- `action_done` → `done`
- `action_reset_draft` → `draft`

## Conventions

- **Selection:** typed `const` blocks in `selection_types.go`.
- **Same-module relations:** `Many2One[MyModule]`, `One2Many[MyModuleLine]`.
- **Cross-module relations:** types from `zrefs.go` after `make generate`.
- Compute handlers in `computed.go` must match `compute=` tags.
- List views use `type="list"`.
