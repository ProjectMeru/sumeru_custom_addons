# Engagements Cookbook

**Professional services tutorial app** for custom-module authors in `sumeru_custom_addons`. It models client engagements end-to-end while demonstrating every major Sumeru addon pattern: multi-model apps, field widgets, view types, **`depends`**, **model inherit**, **view xpath inherit**, cross-module relations, and **event bridges**.

## Install

From `sumeru_custom_addons/`:

```bash
make setup
make install MODULES=engagement_cookbook
make run
```

Open **Engagements Cookbook** in the Apps launcher. Demo data loads from `data/demo.xml` (`noupdate="1"`).

**Migrating from `my_module`:** technical names and tables changed in v2. Use a fresh database or uninstall the old module before installing `engagement_cookbook`.

## Module dependencies

| Module | Why |
|--------|-----|
| `base` | Companies, users, partners, security groups |
| `contacts` | Partner form UI + target for view inherit |
| `hr` | Lead consultant (`hr.employee`) on engagements |
| `mail` | Chatter on the engagement form |

```json
"depends": ["base", "contacts", "hr", "mail"]
```

After changing `depends` or cross-module field types, run **`make generate`** so `models/zrefs.go` stays in sync.

## Business models

| Model | Role |
|-------|------|
| `engagement.project` | Client engagement — workflow, budget, relations, widgets |
| `engagement.deliverable` | Billable lines with computed `line_total` |
| `engagement.milestone` | Calendar / gantt milestones with object actions |
| `engagement.site` | Client sites on the map view |
| `engagement.tag` | Engagement labels |
| `engagement.service_line` | Consulting / implementation / support lines |

## Menus

| Section | Menu | View modes |
|---------|------|------------|
| Operations | Engagements | kanban, list, form, calendar, gantt |
| Operations | Deliverables | list, form |
| Planning | Milestones | calendar, gantt, list, form |
| Planning | Client sites | map, list, form |
| Analytics | Engagement analysis | graph, pivot, cohort, list |
| Clients | Contacts (extended) | list, form, kanban |
| Configuration | Tags, Service lines, About / tutorial | manager only |

## Model inherit (`inherit=`)

`models/partner_extend.go` adds fields to **`core.partner`**:

- `engagement_tier` — standard / premium / strategic
- `strategic_account` — boolean flag

Run `make generate`, then `make update MODULES=engagement_cookbook` for schema sync.

## View inherit (xpath)

`views/partner_inherit_views.xml` patches **`contacts.view_core_partner_form`** to show inherited partner fields after **Phone**.

## Event bridge

`hooks.go` subscribes to `record.updated`. When an **`engagement.project`** reaches **`done`**, the client's partner is promoted to **premium** tier (idempotent tutorial behavior).

Production teams often move this into a small **`application: false` bridge module** with `depends` on two apps — see `sale_crm` in `sumeru_addons`.

## Object actions

- **Project:** `action_confirm` (→ active), `action_done`, `action_reset_draft`
- **Milestone:** `action_plan`, `action_done`, `action_reset_draft`

## File map

| Area | Path |
|------|------|
| Models | `models/engagement_*.go`, `models/partner_extend.go` |
| Computed | `models/computed.go` |
| Actions | `models/object_actions.go` |
| Hooks | `hooks.go` |
| Views | `views/*.xml` |
| Demo | `data/demo.xml` |
| Security | `security/security.xml`, `security/sys.access.csv` |

## Docs

Published guide: [Engagements Cookbook](https://projectmeru.github.io/sumeru/docs/addons/howtos/engagement-cookbook.html) (mirror in `sumeru_docs/addons/howtos/engagement-cookbook.md`).
