# Sumeru custom-addons workspace
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

.PHONY: help setup generate assets swc swc-check swc-test dev \
	replace-sumeru replace-sumeru-addons new run build install update check

help:
	@echo "Sumeru custom workspace — common dev flow:"
	@echo "  make setup   - go.mod replaces, generate imports, SWC assets, sumeru.conf"
	@echo "  make run     - generate + assets + start HTTP (alias: make dev)"
	@echo "  make build   - generate + assets + bin/sumeru-erp"
	@echo ""
	@echo "Client bundles (delegates to ../sumeru):"
	@echo "  make assets  - build SWC + login JS when missing or stale"
	@echo "  make swc     - always rebuild client bundles"
	@echo ""
	@echo "Modules:"
	@echo "  make new MODULE=x     - scaffold addon under addons/"
	@echo "  make install MODULES=x - install module(s), no HTTP"
	@echo "  make update MODULES=x  - update module(s) or all"
	@echo ""
	@echo "Vars: SUMERU_ROOT ADDONS_ROOT CONF DB MODULES EXTRA_RUN_FLAGS"

check-sumeru:
	@test -n "$(SUMERU_ABS)" || (echo "SUMERU_ROOT not found ($(SUMERU_ROOT))" >&2; exit 1)

check-addons:
	@test -n "$(ADDONS_ABS)" || (echo "ADDONS_ROOT not found ($(ADDONS_ROOT))" >&2; exit 1)

replace-sumeru: check-sumeru
	go mod edit -replace sumeru=$(REPLACE_SUMERU)
	go mod tidy

replace-sumeru-addons: check-addons
	go mod edit -replace sumeru_addons=$(REPLACE_SUMERU_ADDONS)
	go get sumeru_addons
	go mod tidy

assets: check-sumeru
	$(SUMERU_MAKE) assets

swc: check-sumeru
	$(SUMERU_MAKE) swc

swc-check: check-sumeru
	$(SUMERU_MAKE) swc-check

swc-test: check-sumeru
	$(SUMERU_MAKE) swc-test

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

new: check-sumeru
	@test -n "$(MODULE)" || (echo 'usage: make new MODULE=my_app' >&2; exit 1)
	go run $(SUMERU_ABS)/cmd/sumeru-bp -bp $(MODULE) -out addons
	$(MAKE) generate

run: generate assets
	$(RUN) $(RUN_FLAGS)

dev: run

build: generate assets
	@mkdir -p bin
	go build -o bin/sumeru-erp .

install: generate
	@test -n "$(MODULES)" || (echo 'usage: make install MODULES=my_app' >&2; exit 1)
	$(RUN) -i $(MODULES) --stop-after-init $(RUN_FLAGS)

update: generate
	@test -n "$(MODULES)" || (echo 'usage: make update MODULES=my_app|all' >&2; exit 1)
	$(RUN) -u $(MODULES) --stop-after-init $(RUN_FLAGS)

check: generate swc-check
	go test ./...
