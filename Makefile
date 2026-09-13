# Sumeru Tier-3 workspace — develop custom addons here; ../sumeru stays read-only.
-include config.mk

.DEFAULT_GOAL := help

SUMERU_ROOT ?= ../sumeru
ADDONS_ROOT ?= ../sumeru_addons
CONF        ?= sumeru.conf
OUT         ?= addonimports/zimports.go
DB          ?=
MODULES     ?=
EXTRA_RUN_FLAGS ?=

CONF_ABS := $(if $(filter /%,$(CONF)),$(CONF),$(CURDIR)/$(CONF))
SUMERU_ABS := $(shell cd "$(SUMERU_ROOT)" 2>/dev/null && pwd)
ADDONS_ABS := $(shell cd "$(ADDONS_ROOT)" 2>/dev/null && pwd)
REPLACE_SUMERU := $(if $(filter /%,$(SUMERU_ROOT)),$(SUMERU_ABS),$(SUMERU_ROOT))
REPLACE_SUMERU_ADDONS := $(if $(filter /%,$(ADDONS_ROOT)),$(ADDONS_ABS),$(ADDONS_ROOT))

RUN := go run . -- -c $(CONF)
RUN_FLAGS :=
ifneq ($(strip $(DB)),)
RUN_FLAGS += -d $(DB)
endif
RUN_FLAGS += $(EXTRA_RUN_FLAGS)

IMPORT_GEN := go run $(SUMERU_ABS)/cmd/sumeru-import-gen \
	-root $(SUMERU_ABS) -workspace $(CURDIR) -addons $(ADDONS_ABS) \
	-config $(CONF_ABS) -out $(OUT) -package addonimports

SUMERU_MAKE := $(MAKE) -C $(SUMERU_ABS)

.PHONY: help info setup generate assets swc swc-check swc-test \
	run run-fast dev build install install-base update new check test clean \
	replace-sumeru replace-sumeru-addons pull-upstream \
	module shell db-check \
	kernel-help kernel-check kernel-lint \
	check-sumeru check-addons

# =============================================================================
# Help
# =============================================================================

help:
	@echo "Sumeru custom workspace — run from sumeru_custom_addons/"
	@echo "Develop here; ../sumeru and ../sumeru_addons stay read-only siblings."
	@echo ""
	@echo "QUICK START"
	@echo "  make setup          first time: config, go.mod replaces, generate, assets"
	@echo "  make run            generate + assets (if stale) + HTTP server"
	@echo "  make dev            alias for run"
	@echo ""
	@echo "DEV SERVER"
	@echo "  make run-fast       generate + server (skip asset rebuild — Go/views only)"
	@echo "  make run DB=name    override database for this run"
	@echo "  make run EXTRA_RUN_FLAGS='-p 9090'   custom port or CLI flags"
	@echo ""
	@echo "MODULES"
	@echo "  make new MODULE=x           scaffold under addons/"
	@echo "  make install MODULES=x      install module(s), exit without HTTP"
	@echo "  make update MODULES=x|all   update module(s) or all"
	@echo "  make install-base           install base on fresh database"
	@echo "  make module ARGS='list'     module CLI (list, uninstall, …)"
	@echo ""
	@echo "ASSETS (delegates to ../sumeru — kernel Makefile untouched)"
	@echo "  make assets         build SWC + login JS when missing or stale"
	@echo "  make swc            force rebuild client bundles"
	@echo "  make swc-check      TypeScript typecheck in ../sumeru/core/swc"
	@echo "  make swc-test       vitest in ../sumeru/core/swc"
	@echo ""
	@echo "TOOLS"
	@echo "  make shell          ORM REPL (uses $(CONF))"
	@echo "  make db-check       test PostgreSQL connectivity"
	@echo "  make info           show resolved paths and sibling checkouts"
	@echo ""
	@echo "TEST & BUILD"
	@echo "  make check          generate + swc-check + go test ./..."
	@echo "  make test           go test ./... only"
	@echo "  make build          generate + assets + bin/sumeru-erp"
	@echo "  make clean          remove bin/"
	@echo ""
	@echo "UPSTREAM"
	@echo "  make pull-upstream          git pull ../sumeru and ../sumeru_addons"
	@echo "  make replace-sumeru         re-wire go.mod replace for core"
	@echo "  make replace-sumeru-addons  re-wire replace for standard addons"
	@echo ""
	@echo "KERNEL (optional — invoke ../sumeru targets without leaving this dir)"
	@echo "  make kernel-help    ../sumeru make help"
	@echo "  make kernel-check   ../sumeru make (lint + test + build)"
	@echo "  make kernel-lint    ../sumeru make lint"
	@echo ""
	@echo "VARIABLES (override on CLI or in config.mk)"
	@echo "  SUMERU_ROOT=$(SUMERU_ROOT)"
	@echo "  ADDONS_ROOT=$(ADDONS_ROOT)"
	@echo "  CONF=$(CONF)  DB=$(DB)  MODULES=$(MODULES)  EXTRA_RUN_FLAGS=$(EXTRA_RUN_FLAGS)"
	@echo ""
	@echo "See README.md for daily workflow and addon authoring."

info: check-sumeru check-addons
	@echo "Workspace:  $(CURDIR)"
	@echo "Config:     $(CONF_ABS)"
	@echo "Sumeru:     $(SUMERU_ABS)"
	@echo "Addons:     $(ADDONS_ABS)"
	@echo "Generated:  $(OUT)"
	@cd "$(SUMERU_ABS)" && echo "Core branch:  $$(git branch --show-current) @ $$(git rev-parse --short HEAD)"
	@cd "$(ADDONS_ABS)" && echo "Std branch:   $$(git branch --show-current) @ $$(git rev-parse --short HEAD)"

# =============================================================================
# Path guards
# =============================================================================

check-sumeru:
	@test -n "$(SUMERU_ABS)" || (echo "SUMERU_ROOT not found ($(SUMERU_ROOT))" >&2; exit 1)

check-addons:
	@test -n "$(ADDONS_ABS)" || (echo "ADDONS_ROOT not found ($(ADDONS_ROOT))" >&2; exit 1)

# =============================================================================
# Bootstrap & imports
# =============================================================================

replace-sumeru: check-sumeru
	go mod edit -replace sumeru=$(REPLACE_SUMERU)
	go mod tidy

replace-sumeru-addons: check-addons
	go mod edit -replace sumeru_addons=$(REPLACE_SUMERU_ADDONS)
	go get sumeru_addons
	go mod tidy

setup: check-sumeru check-addons
	@test -f "$(CONF)" || cp sumeru.conf.example "$(CONF)"
	go mod edit -replace sumeru=$(REPLACE_SUMERU)
	go mod edit -replace sumeru_addons=$(REPLACE_SUMERU_ADDONS)
	go get sumeru_addons
	go mod tidy
	$(MAKE) generate
	$(MAKE) assets

generate: check-sumeru check-addons
	$(IMPORT_GEN)

pull-upstream: check-sumeru check-addons
	cd "$(SUMERU_ABS)" && git pull
	cd "$(ADDONS_ABS)" && git pull
	@echo "Upstream pulled — run 'make run' if imports or assets may have changed."

# =============================================================================
# Assets (read-only delegate to ../sumeru)
# =============================================================================

assets: check-sumeru
	$(SUMERU_MAKE) assets

swc: check-sumeru
	$(SUMERU_MAKE) swc

swc-check: check-sumeru
	$(SUMERU_MAKE) swc-check

swc-test: check-sumeru
	$(SUMERU_MAKE) swc-test

# =============================================================================
# Server
# =============================================================================

run: generate assets
	$(RUN) $(RUN_FLAGS)

run-fast: generate
	$(RUN) $(RUN_FLAGS)

dev: run

build: generate assets
	@mkdir -p bin
	go build -o bin/sumeru-erp .

install: generate
	@test -n "$(MODULES)" || (echo 'usage: make install MODULES=my_app' >&2; exit 1)
	$(RUN) -i $(MODULES) --stop-after-init $(RUN_FLAGS)

install-base: generate
	$(RUN) -i base --stop-after-init $(RUN_FLAGS)

update: generate
	@test -n "$(MODULES)" || (echo 'usage: make update MODULES=my_app|all' >&2; exit 1)
	$(RUN) -u $(MODULES) --stop-after-init $(RUN_FLAGS)

# =============================================================================
# Modules & CLI tools (kernel binaries, workspace config)
# =============================================================================

new: check-sumeru
	@test -n "$(MODULE)" || (echo 'usage: make new MODULE=my_app' >&2; exit 1)
	go run $(SUMERU_ABS)/cmd/sumeru-bp -bp $(MODULE) -out addons
	$(MAKE) generate

module: check-sumeru generate
	go run $(SUMERU_ABS)/cmd/sumeru-module -- -c $(CONF_ABS) $(ARGS)

shell: check-sumeru
	go run $(SUMERU_ABS)/cmd/sumeru-shell -- -c $(CONF_ABS)

db-check: check-sumeru
	go run $(SUMERU_ABS)/cmd/sumeru-db-check -- -c $(CONF_ABS)

# =============================================================================
# Test & build hygiene
# =============================================================================

check: generate swc-check
	go test ./...

test:
	go test ./...

clean:
	rm -rf bin/

# =============================================================================
# Kernel delegate (optional — for verifying ../sumeru without cd)
# =============================================================================

kernel-help: check-sumeru
	$(SUMERU_MAKE) help

kernel-check: check-sumeru
	$(SUMERU_MAKE)

kernel-lint: check-sumeru
	$(SUMERU_MAKE) lint
