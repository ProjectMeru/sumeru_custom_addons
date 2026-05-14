# Optional local overrides (gitignored): SUMERU_ROOT, CONF, OUT, etc.
-include config.mk

.PHONY: generate run help replace-sumeru show-sumeru

# Path to the standard sumeru Go module (directory containing go.mod).
# Default: sibling ../sumeru.
#
# Prod (absolute install, stable path):
#   make replace-sumeru SUMERU_ROOT=/opt/sumeru
#
# Dev (clone anywhere; keep a relative segment in go.mod):
#   make replace-sumeru SUMERU_ROOT=../sumeru
#   make replace-sumeru SUMERU_ROOT=../../projects/sumeru
#
# Or set SUMERU_ROOT once in config.mk (see config.mk.example).
SUMERU_ROOT ?= ../sumeru

CONF ?= sumeru.conf
# Absolute path: -out is relative to SUMERU_ROOT when not absolute (see sumeru-import-gen).
OUT ?= $(CURDIR)/addonimports/zimports.go

# import-gen resolves a relative -config against SUMERU_ROOT; workspace INI lives here → absolute.
CONF_FOR_GEN := $(if $(filter /%,$(CONF)),$(CONF),$(CURDIR)/$(CONF))

# Resolved standard tree (must exist) for import-gen -root and go run …/cmd.
SUMERU_ROOT_ABS := $(shell cd "$(SUMERU_ROOT)" 2>/dev/null && pwd || echo "")

# go.mod replace target: absolute for prod-style SUMERU_ROOT, else keep the path literal (portable dev).
REPLACE_SUMERU := $(if $(filter /%,$(SUMERU_ROOT)),$(SUMERU_ROOT_ABS),$(SUMERU_ROOT))

# SCSS paths relative to SUMERU_ROOT
SCSS_SRC := $(SUMERU_ROOT_ABS)/core/engine/assets/scss/main.scss
CSS_OUT  := $(SUMERU_ROOT_ABS)/core/engine/assets/css/sumeru.css

# Regenerate blank imports from addons_path in CONF (see sumeru.conf.example).
generate:
	@test -n "$(SUMERU_ROOT_ABS)" || (echo "SUMERU_ROOT=$(SUMERU_ROOT) is not a directory; fix path or see: make help" >&2 && exit 1)
	go run $(SUMERU_ROOT_ABS)/cmd/sumeru-import-gen -root $(SUMERU_ROOT_ABS) -config $(CONF_FOR_GEN) -out $(OUT) -package addonimports

	@echo "go.mod: replace sumeru => $(REPLACE_SUMERU)"

# Compile SCSS to CSS using Sass
sass:
	@test -n "$(SUMERU_ROOT_ABS)" || (echo "SUMERU_ROOT not set" >&2 && exit 1)
	@echo "Compiling SCSS: $(SCSS_SRC) -> $(CSS_OUT)..."
	@mkdir -p $(dir $(CSS_OUT))
	sass $(SCSS_SRC) $(CSS_OUT) --style compressed --no-source-map

# Watch SCSS for changes (run in a separate terminal)
sass-watch:
	@test -n "$(SUMERU_ROOT_ABS)" || (echo "SUMERU_ROOT not set" >&2 && exit 1)
	@echo "Watching SCSS in $(SUMERU_ROOT_ABS)/core/engine/assets/scss/ ..."
	sass $(SCSS_SRC):$(CSS_OUT) --watch --style expanded

# Set go.mod: replace sumeru => $(REPLACE_SUMERU), then go mod tidy.
replace-sumeru:
	@test -n "$(SUMERU_ROOT_ABS)" || (echo "SUMERU_ROOT=$(SUMERU_ROOT): no such directory" >&2 && exit 1)
	go mod edit -replace sumeru=$(REPLACE_SUMERU)
	go mod tidy
	@echo "go.mod: replace sumeru => $(REPLACE_SUMERU)"

show-sumeru:
	@echo "SUMERU_ROOT=$(SUMERU_ROOT)"
	@echo "resolved tree: $(SUMERU_ROOT_ABS)"
	@echo "replace segment: $(REPLACE_SUMERU)"
	@grep '^replace sumeru' go.mod || true

# Dev server: regenerate imports then run this module’s main.
run: generate sass
	go run . -- -c $(CONF) $(EXTRA_RUN_FLAGS)

help:
	@echo "Variables: SUMERU_ROOT, CONF, OUT. Optional config.mk — see config.mk.example"
	@echo "Targets:"
	@echo "  make replace-sumeru  - write go.mod replace sumeru => path (abs if SUMERU_ROOT is absolute, else literal); tidy"
	@echo "  make show-sumeru     - print SUMERU_ROOT, resolved path, and go.mod replace line"
	@echo "  make generate        - refresh addonimports/zimports.go (uses CONF in this dir by default)"
	@echo "  make run             - generate then go run . -- -c $(CONF)"
	@echo "  copy sumeru.conf.example to sumeru.conf and adjust paths first."
