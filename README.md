# sumeru_custom_addons

**Your Sumeru development workspace** — run the ERP server, build custom addons, and keep core engines read-only.

> **Pre-alpha software** — not for production. See the [sumeru README](https://github.com/ProjectMeru/sumeru/blob/main/README.md) for the full status notice.

This repository is **Tier 3** in the Sumeru stack: your client or project workspace with its own `go.mod`, generated **`addonimports/`**, and local **`sumeru.conf`**. Pull **`sumeru`** and **`sumeru_addons`** as read-only siblings; develop and commit only here.

## Three-tier layout

| Tier | Repository | Role |
| ---- | ---------- | ---- |
| 1 | **[sumeru](https://github.com/ProjectMeru/sumeru)** | Core engine + kernel apps — **read-only** |
| 2 | **[sumeru_addons](https://github.com/ProjectMeru/sumeru_addons)** | Standard business apps (CRM, Sales, …) — **read-only** |
| 3 | **sumeru_custom_addons** (this repo) | Custom modules under `addons/`, config, `make run` |

```text
parent/
  sumeru/                 # Tier 1
  sumeru_addons/          # Tier 2
  sumeru_custom_addons/   # Tier 3 — you work here
```

## What you get when you run the server

Starting **`make run`** launches the full Sumeru web stack:

- **Settings hub** and **app launcher** with multi-company support
- **Workspace views** — list, form, kanban, graph, pivot, calendar, gantt, map, cohort
- **Collection control bar** — search, filters, multi group-by, custom domain rules, saved-search favorites
- **Standard apps** from `sumeru_addons` plus **your modules** under `addons/`
- **JSON-RPC API** at `POST /api/rpc` for integrations

## Prerequisites

- [Go](https://go.dev/dl/) (same version as `../sumeru/go.mod`, currently 1.26.2+)
- [Node.js](https://nodejs.org/) (npm — builds the SWC workspace UI in `../sumeru`)
- [PostgreSQL](https://www.postgresql.org/) and an empty database matching `db_name` in your INI
- Git checkouts of **`sumeru`** and **`sumeru_addons`** as siblings (defaults: `../sumeru`, `../sumeru_addons`)

## Quick start

Work from **this directory**:

```bash
mkdir -p ~/sumeru_erp && cd ~/sumeru_erp
git clone git@github.com:ProjectMeru/sumeru.git
git clone git@github.com:ProjectMeru/sumeru_addons.git
git clone git@github.com:ProjectMeru/sumeru_custom_addons.git

# Create a PostgreSQL database matching db_name in your INI, e.g.:
#   psql -c "CREATE DATABASE sumeru;"

cd sumeru_custom_addons
cp sumeru.conf.example sumeru.conf   # edit db_*, http_port, addons_path
make setup    # go.mod replaces, imports, SWC + login JS bundles
make run      # regenerate imports; rebuild assets when missing or stale
```

Open **`http://localhost:8080`** (or your `http_port`). `/` redirects to **`/web/apps`**.

Optional: copy **`config.mk.example`** → **`config.mk`** (gitignored) to pin `SUMERU_ROOT`, `ADDONS_ROOT`, or default `DB`.

## Daily developer workflow

| Task | Command |
| ---- | ------- |
| First-time bootstrap | `make setup` |
| Start dev server | `make run` or `make dev` |
| After editing SWC in `../sumeru` | `make swc` (or `make run` rebuilds when sources are stale) |
| New custom module | `make new MODULE=my_app` then `make install MODULES=my_app` |
| Reload views / manifest data | `make update MODULES=my_app` or `make update MODULES=all` |
| Pull upstream core | see [Keeping upstream updated](#keeping-upstream-updated) |
| Tests | `make check` |
| Production binary | `make build` → `bin/sumeru-erp` |
| Other database | `make run DB=sumeru_staging` |
| Custom port | `make run EXTRA_RUN_FLAGS='-p 9090'` |

```text
make run  →  generate (addonimports)  →  assets (../sumeru SWC)  →  HTTP server
```

Run **`make help`** for a short target list.

## Client assets

The browser UI is built in **`../sumeru`**, not stored in git. Sources live under `../sumeru/core/swc/src/`; outputs are gitignored:

| Output | Purpose |
| ------ | ------- |
| `../sumeru/core/engine/assets/swc/swc.js` | Workspace UI |
| `../sumeru/core/engine/assets/js/sumeru-password-*.js` | Login / setup helpers |

**`make setup`** and **`make run`** call `make -C ../sumeru assets`, which builds bundles when missing or when TypeScript sources changed. Use **`make swc`** to force a full rebuild. Node.js is required on the machine where you run setup/run.

## Developing custom addons

Put modules under **`addons/<technical_name>/`** with the usual layout (`manifest.json`, `init.go`, models, views, security`). Reference tutorial: **`addons/engagement_cookbook`**.

```bash
make new MODULE=my_app
make install MODULES=my_app
```

Ensure **`./addons`** is on **`addons_path`** in `sumeru.conf`:

```ini
addons_path = ../sumeru/addons,../sumeru_addons,./addons
```

After adding or removing an addon, run **`make generate`** so `addonimports/zimports.go` and per-addon **`zmodels.go`** / **`zrefs.go`** stay in sync.

### Cross-module relations (`zrefs.go`)

Never edit `sumeru/core` for comodel types. Instead:

1. List upstream modules in **`manifest.json`**: `"depends": ["base", "hr"]`
2. Run **`make generate`**
3. Use generated types in your models:

```go
CompanyID  sdk.Many2One[CoreCompany]  // from zrefs.go
EmployeeID sdk.Many2One[HrEmployee]   // after depends includes hr
```

Add a new dependency → run **`make generate`** again.

## Make this repo yours

This repo is a **template workspace**. After cloning, point **`origin`** at your project remote:

```bash
git clone git@github.com:ProjectMeru/sumeru_custom_addons.git
cd sumeru_custom_addons
cp sumeru.conf.example sumeru.conf
make setup

git remote remove origin
git remote add origin git@github.com:YOUR_ORG/YOUR_CLIENT_ADDONS.git
git push -u origin main
```

Do **not** commit client work into `sumeru` or `sumeru_addons` — keep those read-only and updatable via `git pull`.

## Keeping upstream updated

```bash
cd ../sumeru && git pull
cd ../sumeru_addons && git pull
cd ../sumeru_custom_addons && make run
```

`make run` regenerates imports and rebuilds SWC assets when needed. Commit and push **only** this repo to your project `origin`.

## Makefile reference

### Variables

Set on the command line or in **`config.mk`** (from **`config.mk.example`**):

| Variable | Default | Purpose |
| -------- | ------- | ------- |
| `SUMERU_ROOT` | `../sumeru` | Core checkout |
| `ADDONS_ROOT` | `../sumeru_addons` | Standard addons checkout |
| `CONF` | `sumeru.conf` | INI path |
| `DB` | _(empty)_ | Override `db_name` for this run |
| `MODULES` | _(empty)_ | Comma-separated names for install/update |
| `EXTRA_RUN_FLAGS` | _(empty)_ | Extra CLI flags (e.g. `-p 9090`) |
| `OUT` | `addonimports/zimports.go` | Generated imports file |

### Targets

| Target | Purpose |
| ------ | ------- |
| `setup` | Config, go.mod replaces, `generate`, SWC assets |
| `run` / `dev` | `generate` + `assets` + HTTP server |
| `generate` | Refresh `addonimports/zimports.go`, `init.go`, `zmodels.go`, `zrefs.go` |
| `assets` | Build SWC + login JS when missing or stale (delegates to `../sumeru`) |
| `swc` | Force rebuild client bundles |
| `new MODULE=x` | Scaffold addon under `addons/` |
| `install MODULES=x` | Install module(s), exit without HTTP |
| `update MODULES=x` | Update module(s) or `all`, exit without HTTP |
| `build` | `generate` + `assets` + `bin/sumeru-erp` |
| `check` | SWC typecheck + `go test ./...` |
| `replace-sumeru` | Re-wire `go.mod` replace for core |
| `replace-sumeru-addons` | Re-wire replace for standard addons |
| `swc-check` / `swc-test` | Delegate TypeScript check / vitest to core |
| `help` | Print common targets |

Configuration keys, server CLI flags (`-c`, `-d`, `-p`, `-i`, `-u`), and import-gen details: **[sumeru README](https://github.com/ProjectMeru/sumeru/blob/main/README.md)**, [Configuration guide](https://projectmeru.github.io/sumeru/docs/guides/start/configuration.html), [Tooling docs](https://projectmeru.github.io/sumeru/docs/guides/build/tooling.html).

## Documentation

| Resource | Contents |
| -------- | -------- |
| This README | Workspace runner, daily workflow, custom addons |
| [sumeru README](https://github.com/ProjectMeru/sumeru/blob/main/README.md) | Core engine, client assets, architecture |
| [sumeru_addons README](https://github.com/ProjectMeru/sumeru_addons/blob/main/README.md) | Standard business apps |
| [Documentation home](https://projectmeru.github.io/sumeru/docs/) | Guides and API reference |
| [Creating an addon](https://projectmeru.github.io/sumeru/docs/guides/build/creating-an-addon.html) | Module authoring |

## License

Apache License 2.0 — see [LICENSE](LICENSE). Custom addons under `addons/` follow this license unless you specify otherwise for derived work.
