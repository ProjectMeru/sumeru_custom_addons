# Engagements Cookbook

**Professional services tutorial app** for custom-module authors in `sumeru_custom_addons`. It models client engagements end-to-end while demonstrating every major Sumeru addon pattern: multi-model apps, field widgets, view types, **reporting**, **bulk import**, **conditional fields**, **`depends`**, **model inherit**, **view xpath inherit**, cross-module relations, and **event bridges**.

## Install

From `sumeru_custom_addons/`:

```bash
make setup
make install MODULES=engagement_cookbook
make run
```

Open **Engagements Cookbook** in the Apps launcher. Demo data loads from `data/demo.xml` (`noupdate="1"`).

**Migrating from `my_module`:** technical names and tables changed in v2. Use a fresh database or uninstall the old module before installing `engagement_cookbook`.

**Upgrading to v2.1:** run `make generate` then `make update MODULES=engagement_cookbook` to add the `engagement.timesheet` model and new views.

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
| `engagement.timesheet` | Time entries — flat model for bulk CSV import |
| `engagement.site` | Client sites on the map view |
| `engagement.tag` | Engagement labels |
| `engagement.service_line` | Consulting / implementation / support lines |

## Menus

| Section | Menu | View modes |
|---------|------|------------|
| Operations | Engagements | kanban, list, form, calendar, gantt |
| Operations | Deliverables | list, form |
| Operations | Timesheets | list, form |
| Planning | Milestones | calendar, gantt, list, form |
| Planning | Client sites | map, list, form |
| Analytics | Engagement analysis | graph, pivot, cohort, list |
| Analytics | Timesheet analysis | graph, list, form |
| Clients | Contacts (extended) | list, form, kanban |
| Configuration | Tags, Service lines, About / tutorial | manager only |

## Dynamic field conditions

Sumeru evaluates **expression attributes** on fields at runtime (see `sumeru_docs/core/reference/field-attributes.md`):

```xml
<field name="completed_at" invisible="state != 'done'"/>
<field name="archived" readonly="state == 'done'"/>
<field name="description" required="kind == 'engagement'"/>
```

| Model | Field | Modifier | Expression |
|-------|-------|----------|------------|
| `engagement.project` | `completed_at` | invisible | `state != 'done'` |
| `engagement.project` | `archived` | readonly | `state == 'done'` |
| `engagement.project` | `description` | required | `kind == 'engagement'` |
| `engagement.project` | `date_stop` | required | `state == 'active'` |
| `engagement.deliverable` | `note` | required | `quantity > 1` |
| `engagement.milestone` | `date_stop` | required | `state == 'planned'` |
| `engagement.milestone` | `employee_id` | invisible | `state == 'draft'` |
| `engagement.timesheet` | `rate` | invisible | `billable == false` |
| `engagement.timesheet` | `description` | required | `hours > 8` |

Toggle **Status** and **Kind** on an engagement form (Overview → Conditions demo) to see fields appear, lock, or become required.

## Reporting (CSV, XLSX, PDF)

List and form views declare export via a `<report>` child — no addon Go code required:

```xml
<report download="csv,xlsx,pdf" pdf_sizes="a4,letter"/>
```

| View | Formats | Bulk upload |
|------|---------|-------------|
| Engagements list | CSV, XLSX, PDF | — |
| Engagement form | CSV, PDF | — |
| Deliverables list | CSV, XLSX, PDF | CSV |
| Timesheets list | CSV, XLSX, PDF | CSV |
| Milestones list | CSV, PDF | — |

Use the **Report** menu on the list toolbar. Exports respect the current search domain (max **500 rows**). PDF supports A4 and Letter page sizes.

**Charts vs exports:** graph views (bar, line, pie) are interactive Chart.js analytics in the **Analytics** menus. They are separate from tabular XLSX/CSV downloads — Excel files do not embed charts.

## Bulk import

Enable bulk CSV upload on a list view:

```xml
<report download="csv,xlsx,pdf" upload="bulk" modes="create,upsert"/>
```

Workflow on **Timesheets** or **Deliverables**:

1. Report → **Bulk upload CSV** (or download template first)
2. Map CSV columns to model fields in the mapping wizard
3. Preview first 10 rows, then confirm import

Import is **CSV only** (not XLSX upload). Max file size **8 MB**. See `sumeru_docs/business/platform/export-and-import.md`.

## Graph views

| Menu | Chart types | Dimensions |
|------|-------------|------------|
| Engagement analysis | bar, line, pie | state / priority / kind × amount |
| Timesheet analysis | bar, pie | category / project × hours |

Switch chart type via the view switcher when multiple graph views exist for the same model.

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
