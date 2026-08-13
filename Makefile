# Sumeru custom-addons workspace — wire go.mod, generate imports, install/update, run ERP.
# Local overrides: cp config.mk.example config.mk  (gitignored)
# CLI mapping: MODULES → -i/-u, DB → -d, EXTRA_RUN_FLAGS → -p and other flags

-include config.mk

.DEFAULT_GOAL := help

# --- Configurable (override: make run DB=sumeru_test MODULES=…) ---
SUMERU_ROOT     ?= ../sumeru
ADDONS_ROOT     ?= ../sumeru_addons
CONF            ?= sumeru.conf
OUT             ?= $(CURDIR)/addonimports/zimports.go
DB              ?=
MODULES         ?=
EXTRA_RUN_FLAGS ?=

# --- Resolved paths ---
CONF_ABS := $(if $(filter /%,$(CONF)),$(CONF),$(CURDIR)/$(CONF))
SUMERU_ABS := $(shell cd "$(SUMERU_ROOT)" 2>/dev/null && pwd)
ADDONS_ABS := $(shell cd "$(ADDONS_ROOT)" 2>/dev/null && pwd)
REPLACE_SUMERU        := $(if $(filter /%,$(SUMERU_ROOT)),$(SUMERU_ABS),$(SUMERU_ROOT))
REPLACE_SUMERU_ADDONS := $(if $(filter /%,$(ADDONS_ROOT)),$(ADDONS_ABS),$(ADDONS_ROOT))

RUN := go run . -- -c $(CONF)

# Optional -d and extra flags (-p, etc.)
RUN_FLAGS :=
ifneq ($(strip $(DB)),)
RUN_FLAGS += -d $(DB)
endif
RUN_FLAGS += $(EXTRA_RUN_FLAGS)

.PHONY: help setup conf paths show-sumeru wire replace-sumeru replace-sumeru-addons \
        generate run build install update sync i u cli tidy

help:
	@echo "Variables: SUMERU_ROOT ADDONS_ROOT CONF OUT DB MODULES EXTRA_RUN_FLAGS"
	@echo "           (optional config.mk — see config.mk.example)"
	@echo ""
	@echo "Flag mapping (same as go run . -- …):"
	@echo "  MODULES=x     →  -i x (install) or -u x (update)"
	@echo "  DB=name       →  -d name  (override db_name in INI for this run)"
	@echo "  EXTRA_RUN_FLAGS='-p 9090'  →  other CLI flags"
	@echo ""
	@echo "Bootstrap:"
	@echo "  make conf              copy sumeru.conf.example → sumeru.conf if missing"
	@echo "  make setup             conf + wire go.mod + generate (first-time / after clone)"
	@echo ""
	@echo "Day-to-day:"
	@echo "  make run [DB=name]               generate + start HTTP server"
	@echo "  make generate                    refresh addonimports/zimports.go"
	@echo "  make build                       generate + build bin/sumeru-erp"
	@echo "  make install MODULES=x [DB=name] install (-i), then exit"
	@echo "  make update  MODULES=x [DB=name] update (-u), then exit"
	@echo "  make i / make u                  short aliases for install / update"
	@echo "  make sync  MODULES=x [DB=name]   install then update same modules"
	@echo "  make cli [EXTRA_RUN_FLAGS='…']   arbitrary flags (add --stop-after-init for -i/-u)"
	@echo ""
	@echo "Update (-u) semantics:"
	@echo "  MODULES=all          every installed module on disk"
	@echo "  MODULES=mod1,mod2    only those names; not installed → skip silently"
	@echo ""
	@echo "go.mod:"
	@echo "  make wire              replace sumeru + sumeru_addons, then tidy"
	@echo "  make replace-sumeru    replace sumeru only"
	@echo "  make replace-sumeru-addons  replace sumeru_addons only"
	@echo "  make tidy              go mod tidy"
	@echo ""
	@echo "Inspect:  make paths"
	@echo ""
	@echo "Examples:"
	@echo "  make setup"
	@echo "  make run DB=sumeru_staging"
	@echo "  make install MODULES=my_module,student"
	@echo "  make update  MODULES=all"
	@echo "  make update  MODULES=contacts,student DB=sumeru_dev"
	@echo ""
	@echo "Raw CLI (after make generate):"
	@echo "  go run . -- -c $(CONF) -i my_module --stop-after-init"
	@echo "  go run . -- -c $(CONF) -u all --stop-after-init"
	@echo "  go run . -- -c $(CONF) -d sumeru_test -p 9090"

# --- Checks ---
check-sumeru:
	@test -n "$(SUMERU_ABS)" || (echo "SUMERU_ROOT=$(SUMERU_ROOT): directory not found" >&2; exit 1)

check-addons:
	@test -n "$(ADDONS_ABS)" || (echo "ADDONS_ROOT=$(ADDONS_ROOT): directory not found" >&2; exit 1)

# --- Bootstrap ---
conf:
	@test -f "$(CONF)" || cp sumeru.conf.example "$(CONF)"
	@test -f "$(CONF)" && echo "using $(CONF)" || (echo "failed to create $(CONF)" >&2; exit 1)

setup: conf wire generate

# --- go.mod wiring ---
wire: replace-sumeru replace-sumeru-addons

replace-sumeru: check-sumeru
	go mod edit -replace sumeru=$(REPLACE_SUMERU)
	go mod tidy
	@echo "go.mod: replace sumeru => $(REPLACE_SUMERU)"

replace-sumeru-addons: check-addons
	go mod edit -replace sumeru_addons=$(REPLACE_SUMERU_ADDONS)
	go mod tidy
	@echo "go.mod: replace sumeru_addons => $(REPLACE_SUMERU_ADDONS)"

tidy:
	go mod tidy

# --- Generate & run ---
generate: check-sumeru
	go run $(SUMERU_ABS)/cmd/sumeru-import-gen \
		-root $(SUMERU_ABS) -config $(CONF_ABS) -out $(OUT) -package addonimports

run: generate
	$(RUN) $(RUN_FLAGS)

build: generate
	@mkdir -p bin
	go build -o bin/sumeru-erp .

install i: generate
	@test -n "$(MODULES)" || (echo 'usage: make install MODULES=my_module[,other]  (alias: make i)' >&2; exit 1)
	$(RUN) -i $(MODULES) --stop-after-init $(RUN_FLAGS)

update u: generate
	@test -n "$(MODULES)" || (echo 'usage: make update MODULES=my_module  or  MODULES=all  (alias: make u)' >&2; exit 1)
	$(RUN) -u $(MODULES) --stop-after-init $(RUN_FLAGS)

sync: generate
	@test -n "$(MODULES)" || (echo 'usage: make sync MODULES=my_module[,other]' >&2; exit 1)
	$(RUN) -i $(MODULES) -u $(MODULES) --stop-after-init $(RUN_FLAGS)

cli: generate
	$(RUN) $(RUN_FLAGS)

# --- Inspect ---
paths show-sumeru:
	@echo "CONF=$(CONF)"
	@echo "DB=$(if $(strip $(DB)),$(DB),(from INI db_name))"
	@echo "SUMERU_ROOT=$(SUMERU_ROOT)  →  $(SUMERU_ABS)"
	@echo "ADDONS_ROOT=$(ADDONS_ROOT)  →  $(ADDONS_ABS)"
	@echo "replace sumeru        => $(REPLACE_SUMERU)"
	@echo "replace sumeru_addons => $(REPLACE_SUMERU_ADDONS)"
	@grep '^replace sumeru' go.mod 2>/dev/null || true
