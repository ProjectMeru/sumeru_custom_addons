# sumeru_custom_addons

Thin **workspace** for **your** custom / client Sumeru addons. This repo sits next to the standard **`sumeru`** and **`sumeru_addons`** trees: own `go.mod`, generated **`addonimports/`**, and a local **`sumeru.conf`**, so you can **`git pull`** upstream core without committing generated glue here.

Put modules you create for a client or project under **`addons/`**. Keep **`sumeru`** and **`sumeru_addons`** read-only.

## License

This repository is licensed under the **Apache License 2.0**. See [LICENSE](LICENSE) for the full text.

Custom addons you add under `addons/` are part of this workspace and are covered by the same license unless you change licensing for your derived work.

## Architecture: 3-Tier Split

1. **Tier 1: Core Framework (`sumeru`)** — Standard engine and base models. **READ-ONLY**.
2. **Tier 2: Standard Addons (`sumeru_addons`)** — Core business modules (CRM, Sales, Inventory). **READ-ONLY**.
3. **Tier 3: Custom Workspace (`sumeru_custom_addons`)** — Your development area. Custom modules go in `addons/`.

Expected sibling layout:

```text
parent/
  sumeru/                 # Tier 1 (read-only)
  sumeru_addons/          # Tier 2 (read-only)
  sumeru_custom_addons/   # Tier 3 (this repo — your custom addons)
```

## Prerequisites

- **Go** (same version as `../sumeru/go.mod`)
- **PostgreSQL** and an empty database matching `db_name` in your INI
- Checkouts of **`sumeru`** and **`sumeru_addons`** as siblings (defaults: `../sumeru`, `../sumeru_addons`)

## Quick start (how to use)

Work from **this directory** unless you use absolute paths in the INI.

```bash
# Clone as siblings (adjust URLs to your org)
git clone <sumeru-remote> sumeru
git clone <sumeru-addons-remote> sumeru_addons
git clone git@github.com:ProjectMeru/sumeru_custom_addons.git
cd sumeru_custom_addons

# Local config (gitignored)
cp sumeru.conf.example sumeru.conf
# Edit db_* and paths in sumeru.conf

# Bootstrap (wire go.mod, generate imports, create sumeru.conf if missing)
make setup
make run
```

| Step | Command | What it does |
| --- | --- | --- |
| 1. Config | `cp sumeru.conf.example sumeru.conf` then edit | Local INI (gitignored). |
| 2. Bootstrap | `make setup` | Wire `go.mod`, generate `addonimports/zimports.go`, create `sumeru.conf` if missing. |
| 3. Start server | `make run` | Generate imports then start HTTP. |
| 4. Install | `make install MODULES=my_module` | `-i` then exit (no HTTP). |
| 5. Update | `make update MODULES=all` or `MODULES=mod1,mod2` | `-u` then exit; `all` = installed only; explicit list skips uninstalled. |
| Other DB | `make run DB=sumeru_staging` | `-d` overrides INI `db_name` for this run. |
| Build binary | `make build` | `./bin/sumeru-erp -c sumeru.conf` |

Inspect paths:

```bash
make paths
make help
```

Optional: copy **`config.mk.example`** → **`config.mk`** (gitignored) to pin `SUMERU_ROOT`, `ADDONS_ROOT`, or default `DB`.

## Make this repo yours (re-point origin)

This repository is a **template workspace**. After cloning, remove the template `origin` and attach **your** project remote so your client/custom addons stay on your branch.

```bash
# 1) Clone the template
git clone git@github.com:ProjectMeru/sumeru_custom_addons.git
cd sumeru_custom_addons

# 2) Ready the workspace (siblings sumeru / sumeru_addons must exist)
cp sumeru.conf.example sumeru.conf
# edit db_* and paths in sumeru.conf
make setup

# 3) Detach from template upstream; attach your project remote
git remote remove origin
git remote add origin git@github.com:YOUR_ORG/YOUR_CLIENT_ADDONS.git
git push -u origin main
# or: git checkout -b your-branch && git push -u origin your-branch

# 4) Day-to-day: keep core updated, develop and push only here
cd ../sumeru && git pull
cd ../sumeru_addons && git pull
cd ../sumeru_custom_addons && make generate
# commit and push this repo to your origin
```

Do **not** commit changes into `sumeru` or `sumeru_addons` for client work — those stay read-only and updatable via `git pull`.

## Developing custom addons

Create modules under `addons/<technical_name>/` with the usual layout (`manifest.json`, `init.go`, models, views, security). Sample modules in this repo:

- **`my_module`** — example application module
- **`student`** — Student Management sample app

Ensure `./addons` is on `addons_path` (see `sumeru.conf.example`), e.g.:

```ini
addons_path = ../sumeru/addons,../sumeru_addons,./addons
```

After adding or removing an addon, run **`make generate`** so `addonimports/zimports.go` picks up the blank imports.

Install / update (Make or raw CLI):

```bash
make install MODULES=my_module,student
make update  MODULES=all
make update  MODULES=my_module DB=sumeru_dev

# equivalent go run (run make generate first)
go run . -- -c sumeru.conf -i my_module,student --stop-after-init
go run . -- -c sumeru.conf -u all --stop-after-init
go run . -- -c sumeru.conf -d sumeru_test -u my_module --stop-after-init
```

## Keeping upstream updated

```bash
cd ../sumeru && git pull
cd ../sumeru_addons && git pull
cd ../sumeru_custom_addons && make generate
```

Then continue development and push **only** this repo to your project `origin`.

---

## Makefile variables and targets

Variables (command line or **`config.mk`** from **`config.mk.example`**):

| Variable | Default | Maps to | Purpose |
| --- | --- | --- | --- |
| **`SUMERU_ROOT`** | `../sumeru` | — | Core checkout for import-gen and **`make replace-sumeru`**. |
| **`ADDONS_ROOT`** | `../sumeru_addons` | — | Standard addons checkout for **`make replace-sumeru-addons`**. |
| **`CONF`** | `sumeru.conf` | `-c` | INI path. |
| **`DB`** | _(empty)_ | `-d` | Override INI **`db_name`** for this run. |
| **`MODULES`** | _(empty)_ | `-i` / `-u` | Comma-separated module names; **`all`** valid for update only. |
| **`OUT`** | `addonimports/zimports.go` | — | Generated imports path. |
| **`EXTRA_RUN_FLAGS`** | _(empty)_ | `-p`, etc. | Other CLI flags (e.g. **`EXTRA_RUN_FLAGS='-p 9090'`**). |

Targets:

| Target | CLI equivalent |
| --- | --- |
| **`make setup`** | conf + wire + generate |
| **`make run [DB=…]`** | `go run . -- -c $(CONF) [-d …]` |
| **`make install MODULES=…`** / **`make i`** | `-i … --stop-after-init` |
| **`make update MODULES=…`** / **`make u`** | `-u … --stop-after-init` |
| **`make sync MODULES=…`** | `-i … -u … --stop-after-init` |
| **`make cli EXTRA_RUN_FLAGS='…'`** | arbitrary flags |
| **`make generate`** | refresh **`addonimports/zimports.go`** |
| **`make build`** | **`bin/sumeru-erp`** binary |
| **`make wire`** | both **`replace-*`** + tidy |
| **`make paths`** | print resolved paths |
| **`make help`** | full reference |

---

## Server CLI flags (`go run . -- …`)

Parsed by **`sumeru/core/server`**. Pass **`-c`** unless your cwd and config layout match defaults.

| Flag                                               | Purpose                                                              |
| -------------------------------------------------- | -------------------------------------------------------------------- |
| **`-c <path>`**                                    | Path to the INI file (e.g. **`sumeru.conf`**).                       |
| **`-i mod`** or **`-i mod1,mod2`**                 | **Install** listed modules after startup init.                       |
| **`-u mod`** or **`-u mod1,mod2`** or **`-u all`** | **Update** installed modules from disk (reload XML / metadata).      |
| **`-d <name>`**                                    | Override **`db_name`** from the INI for this run.                    |
| **`--database <name>`**                            | Same as **`-d`**; if both are set, **`--database`** wins.            |
| **`-p <port>`**                                    | HTTP port; overrides **`http_port`** in the INI.                     |
| **`--http-port <port>`**                           | Same; if both **`-p`** and **`--http-port`** are set, **`-p`** wins. |
| **`--stop-after-init`**                            | After **`-i`** / **`-u`**, exit without starting HTTP.               |

Examples:

```bash
go run . -- -c sumeru.conf -p 9090
go run . -- -c sumeru.conf -d sumeru_staging
go run . -- -c sumeru.conf -i company,user,sales --stop-after-init
go run . -- -c /etc/sumeru/prod.conf --http-port 443
go run . -- -c sumeru.conf -u my_module --stop-after-init
```

---

## INI `[options]` keys (`sumeru.conf`)

Section header: **`[options]`**. Format: **`key = value`**. Lines starting with **`#`** or **`;`** are comments.

Path keys (**`addons_path`**, **`sumeru_home`**, **`assets_path`**, **`templates_path`**, **`logo_path`**, **`brand_css`**, **`log_file`**, and similar) resolve **relative values from the INI file’s directory** (unless already absolute).

### Required

| Key               | Purpose                                                                                                                                          |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------ |
| **`db_host`**     | PostgreSQL host.                                                                                                                                 |
| **`db_port`**     | PostgreSQL port.                                                                                                                                 |
| **`db_user`**     | Database user.                                                                                                                                   |
| **`db_password`** | Database password.                                                                                                                               |
| **`db_name`**     | Database name (overridable with **`-d`** / **`--database`**).                                                                                    |
| **`http_port`**   | HTTP listen port (overridable with **`-p`** / **`--http-port`**).                                                                                |
| **`addons_path`** | Comma-separated directories; each immediate subfolder with **`manifest.json`** is an addon. Later roots override the same technical module name. |

### Optional — database

| Key              | Default   | Purpose                                    |
| ---------------- | --------- | ------------------------------------------ |
| **`db_sslmode`** | `disable` | PostgreSQL **`sslmode`** (e.g. `require`). |

### Optional — workspace / paths

| Key                  | Purpose                                                                                                                              |
| -------------------- | ------------------------------------------------------------------------------------------------------------------------------------ |
| **`sumeru_home`**    | Directory of the standard **`sumeru`** checkout. When set, omitted **`assets_path`** / **`templates_path`** default under this tree. |
| **`assets_path`**    | Static files (CSS/JS). Default: **`core/engine/assets`** under **`sumeru_home`** when set.                                           |
| **`templates_path`** | HTML templates. Default: **`core/engine/templates`** (same rules).                                                                   |

### Optional — branding

| Key                        | Purpose                                                                                      |
| -------------------------- | -------------------------------------------------------------------------------------------- |
| **`logo_path`**            | Image served at **`/static/app-logo`**.                                                      |
| **`company_display_name`** | Header chip; if empty and **`company`** is installed, first **`core.company`** name is used. |
| **`user_display_name`**    | Header label; if empty and **`user`** is installed, first **`core.user`** display is used.   |
| **`brand_css`**            | Extra CSS linked as **`/static/brand.css`** after view stylesheets.                          |

### Optional — logging (Zap JSON + optional lumberjack)

| Key                    | Purpose                                                                                    |
| ---------------------- | ------------------------------------------------------------------------------------------ |
| **`log_stdout`**       | Default **true** — emit JSON logs to **stdout**.                                           |
| **`log_file`**         | Optional second sink; path absolutized from the INI directory.                             |
| **`log_rolling`**      | **false** = append-only; **true** = size-based rotation (lumberjack).                      |
| **`log_max_size_mb`**  | Rotate after this many MB per file (default **100** when rolling is on and this is **0**). |
| **`log_max_backups`**  | Number of rotated files to keep.                                                           |
| **`log_max_age_days`** | Delete rotated files older than **N** days (**0** = no age-based pruning).                 |
| **`dev_mode`**         | **true** → Zap **debug** level and other dev-only behavior in core.                        |

Full annotated list (core-only defaults): **`../sumeru/sumeru.conf.example`**.

---

## `sumeru-import-gen` (used by `make generate`)

Invoked as:

`go run $(SUMERU_ROOT)/cmd/sumeru-import-gen -root $(SUMERU_ROOT) -config <absolute-CONF> -out $(OUT) -package addonimports`

| Flag           | Purpose                                                                                |
| -------------- | -------------------------------------------------------------------------------------- |
| **`-root`**    | Standard **`sumeru`** repo root (`module sumeru`).                                     |
| **`-config`**  | Absolute path to the INI whose **`addons_path`** / **`sumeru_home`** define discovery. |
| **`-out`**     | Generated `.go` file path (absolute **`OUT`** in this Makefile).                       |
| **`-package`** | Go package name inside that file (here **`addonimports`**).                            |

## Documentation

| Resource                                                                                      | Contents                                         |
| --------------------------------------------------------------------------------------------- | ------------------------------------------------ |
| This README                                                                                   | Workspace runner, `make generate`, custom addons |
| [`sumeru/README.md`](https://github.com/ProjectMeru/sumeru/blob/main/README.md)               | Core engine, config, CLI                         |
| [`sumeru_addons/README.md`](https://github.com/ProjectMeru/sumeru_addons/blob/main/README.md) | Standard business addon module                   |
