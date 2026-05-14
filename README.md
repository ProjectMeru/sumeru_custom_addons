# sumeru_custom_addons

Thin **workspace** next to the standard **`sumeru`** Go module: own `go.mod`, generated **`addonimports/`**, and a local **`sumeru.conf`** so you can **`git pull`** in **`../sumeru`** without committing generated glue here.

### Architecture: 3-Tier Split
To ensure future updates remain seamless, the system is divided into three layers:
1. **Tier 1: Core Framework (`sumeru`)** — Standard engine and base models. **READ-ONLY**.
2. **Tier 2: Standard Addons (`sumeru_addons`)** — Core business modules (CRM, Sales, Inventory). **READ-ONLY**.
3. **Tier 3: Custom Workspace (`sumeru_custom_addons`)** — Your development area. Custom modules go in `addons/`.

#### Updating:
- `cd ../sumeru && git pull`
- `cd ../sumeru_addons && git pull`
- `cd ../sumeru_custom_addons && make generate`

---

## Prerequisites

- **Go** (same version as `../sumeru/go.mod`)
- **PostgreSQL** and an empty database matching `db_name` in your INI
- Standard **Sumeru** checkout at **`SUMERU_ROOT`** (default `../sumeru`)

---

## Commands to run

Work from **this directory** (`sumeru_custom_addons`) unless you use only absolute paths in the INI.

| Step | Command | What it does |
| ---- | ------- | -------------- |
| 1. Config | `cp sumeru.conf.example sumeru.conf` then edit | Local INI (gitignored if you use the repo’s `.gitignore` for `sumeru.conf`). |
| 2. Point `go.mod` at standard Sumeru | `make replace-sumeru SUMERU_ROOT=../sumeru` | Writes `replace sumeru => …` and runs `go mod tidy`. Use an **absolute** path on servers, e.g. `SUMERU_ROOT=/opt/sumeru`. |
| 3. Regenerate addon blank-imports | `make generate` | Runs `sumeru-import-gen` → **`addonimports/zimports.go`** (gitignored). Required after changing **`addons_path`** or adding/removing addons. |
| 4. Start HTTP server | `make run` | Runs **`make generate`** then **`go run . -- -c sumeru.conf`**. |
| | `go run . -- -c sumeru.conf` | Same server without regenerating imports (only if `addonimports/zimports.go` is already fresh). |
| 5. Install / update apps (optional) | `go run . -- -c sumeru.conf -i acme_demo --stop-after-init` | Install module **`acme_demo`** then exit (no HTTP). |
| | `go run . -- -c sumeru.conf -u acme_demo,workspace_notes` | Reload XML/metadata for those modules. |
| Build binary | `go build -o sumeru-workspace .` | Produces a binary; run with `./sumeru-workspace -c sumeru.conf`. |

Inspect paths:

```bash
make show-sumeru    # SUMERU_ROOT, resolved path, current go.mod replace line
make help           # short target list
```

---

## Makefile variables and targets

Variables (override on the command line, e.g. `make run CONF=/path/to/other.conf`, or put them in **`config.mk`** copied from **`config.mk.example`**):

| Variable | Default | Purpose |
| -------- | ------- | ------- |
| **`SUMERU_ROOT`** | `../sumeru` | Filesystem path to the standard **`sumeru`** repo (must contain **`go.mod`** `module sumeru`). Used for **`go run …/cmd/sumeru-import-gen -root`** and for **`make replace-sumeru`**. |
| **`CONF`** | `sumeru.conf` | INI passed to **`go run . -- -c`** and used by **`make generate`** (resolved to **`$(CURDIR)/$(CONF)`** when not absolute). |
| **`OUT`** | `$(CURDIR)/addonimports/zimports.go` | Absolute output path for generated imports so files are **never** written under **`../sumeru`** by mistake. |
| **`EXTRA_RUN_FLAGS`** | *(empty)* | Appended to **`go run . --`** in **`make run`** (e.g. `EXTRA_RUN_FLAGS='-p 9090 -i sales'`). |

Targets:

| Target | Purpose |
| ------ | ------- |
| **`make generate`** | Refresh **`addonimports/zimports.go`** from **`CONF`** (`addons_path`, etc.). |
| **`make run`** | **`generate`** then start the server with **`-c $(CONF)`**. |
| **`make replace-sumeru`** | Set **`go.mod`** `replace sumeru => …` from **`SUMERU_ROOT`**, then **`go mod tidy`**. |
| **`make show-sumeru`** | Print **`SUMERU_ROOT`**, resolved path, and the **`replace`** line. |
| **`make help`** | Summarize variables and targets. |

---

## Server CLI flags (`go run . -- …`)

Parsed by **`sumeru/core/server`**. All are optional except you must pass **`-c`** (or rely on defaults only if your process cwd and config layout match).

| Flag | Purpose |
| ---- | ------- |
| **`-c <path>`** | Path to the INI file (e.g. **`sumeru.conf`**). |
| **`-i mod`** or **`-i mod1,mod2`** | **Install** listed modules after startup init. |
| **`-u mod`** or **`-u mod1,mod2`** or **`-u all`** | **Update** installed modules from disk (reload XML / metadata). |
| **`-d <name>`** | Override **`db_name`** from the INI for this run. |
| **`--database <name>`** | Same as **`-d`**; if both are set, **`--database`** wins. |
| **`-p <port>`** | Shorthand for HTTP port; overrides **`http_port`** in the INI. |
| **`--http-port <port>`** | Same; if both **`-p`** and **`--http-port`** are set, **`-p`** wins. |
| **`--stop-after-init`** | After **`-i`** / **`-u`**, exit without starting HTTP (only if module flags were used). |

Examples:

```bash
go run . -- -c sumeru.conf -p 9090
go run . -- -c sumeru.conf -d sumeru_staging
go run . -- -c sumeru.conf -i company,user,sales --stop-after-init
go run . -- -c /etc/sumeru/prod.conf --http-port 443
```

---

## INI `[options]` keys (`sumeru.conf`)

Section header: **`[options]`**. Format: **`key = value`**. Lines starting with **`#`** or **`;`** are comments.

Path keys (**`addons_path`**, **`sumeru_home`**, **`assets_path`**, **`templates_path`**, **`logo_path`**, **`brand_css`**, **`log_file`**) resolve **relative values from the INI file’s directory** (unless the value is already absolute).

### Required

| Key | Purpose |
| --- | ------- |
| **`db_host`** | PostgreSQL host. |
| **`db_port`** | PostgreSQL port. |
| **`db_user`** | Database user. |
| **`db_password`** | Database password. |
| **`db_name`** | Database name (overridable at runtime with **`-d`** / **`--database`**). |
| **`http_port`** | HTTP listen port (overridable with **`-p`** / **`--http-port`**). |
| **`addons_path`** | Comma-separated directories; each immediate subfolder with **`manifest.json`** is an addon. **`core/base`** (platform addons) is always prepended; **later roots override** the same technical module name. |

### Optional — database

| Key | Default | Purpose |
| --- | ------- | ------- |
| **`db_sslmode`** | `disable` | PostgreSQL **`sslmode`** (e.g. `require`). |

### Optional — workspace / paths

| Key | Purpose |
| --- | ------- |
| **`sumeru_home`** | Directory of the standard **`sumeru`** checkout. When set, **`core/base`** loads from here; if **`assets_path`** / **`templates_path`** are omitted, defaults are under this tree. Relative → resolved from the INI directory. |
| **`assets_path`** | Static files (CSS/JS). Default: **`core/engine/assets`** (under **`sumeru_home`** if set, else next to INI if INI sits in a tree with **`go.mod`**, else cwd semantics — see upstream docs). |
| **`templates_path`** | HTML templates. Default: **`core/engine/templates`** (same rules). |

### Optional — branding

| Key | Purpose |
| --- | ------- |
| **`logo_path`** | Image served at **`/static/app-logo`**. |
| **`company_display_name`** | Header chip; if empty and **`company`** is installed, first **`res.company`** name is used. |
| **`user_display_name`** | Header label; if empty and **`user`** is installed, first **`res.users`** display is used. |
| **`brand_css`** | Extra CSS linked as **`/static/brand.css`** after view stylesheets. |

### Optional — logging

| Key | Purpose |
| --- | ------- |
| **`log_file`** | If set, **`log`** output is **appended** to this file **and** still written to **stderr**. Parent directories are created. |

Full annotated list: **`../sumeru/sumeru.conf`**.

---

## `sumeru-import-gen` (used by `make generate`)

Invoked as:

`go run $(SUMERU_ROOT)/cmd/sumeru-import-gen -root $(SUMERU_ROOT) -config <absolute-CONF> -out $(OUT) -package addonimports`

| Flag | Purpose |
| ---- | ------- |
| **`-root`** | Standard **`sumeru`** repo root (directory whose **`go.mod`** is **`module sumeru`**). |
| **`-config`** | Absolute path to the INI whose **`addons_path`** / **`sumeru_home`** define discovery. |
| **`-out`** | Generated `.go` file path (absolute **`OUT`** in this Makefile). |
| **`-package`** | Go package name inside that file (here **`addonimports`**). |

---

## Custom addons in this repo

Sample modules live under **`addons/`**:

- **`acme_demo`** — example app with a root menu entry.
- **`workspace_notes`** — second example (`application: false` in manifest).

List them on **`addons_path`** after the standard tree, e.g. **`../sumeru/addons,./addons`**, and run **`make generate`** so **`addonimports/zimports.go`** picks up **`sumeru_custom_addons/addons/...`** imports alongside **`sumeru/...`**.
