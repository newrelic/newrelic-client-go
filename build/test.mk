#
# Makefile fragment for Testing
#

GO           ?= go
TEST_RUNNER  ?= gotestsum

COVERAGE_DIR ?= ./coverage/
COVERMODE    ?= atomic
SRCDIR       ?= .
GO_PKGS      ?= $(shell $(GO) list ./... | grep -v -e "/vendor/" -e "/example" -e "/testhelpers")

PROJECT_MODULE ?= $(shell $(GO) list -m)

LDFLAGS_UNIT ?= '-X $(PROJECT_MODULE)/internal/version.GitTag=$(PROJECT_VER_TAGGED)'

test: test-only
test-only: test-unit test-integration

test-unit: tools
	@echo "=== $(PROJECT_NAME) === [ test-unit        ]: running unit tests..."
	@mkdir -p $(COVERAGE_DIR)
	@$(TEST_RUNNER) -f pkgname --rerun-fails=3 --packages "$(GO_PKGS)" -- -v -ldflags=$(LDFLAGS_UNIT) -parallel 10 -tags unit -covermode=$(COVERMODE) -coverprofile $(COVERAGE_DIR)/unit.tmp $(GO_PKGS)
test-integration: tools
	@echo "=== $(PROJECT_NAME) === [ test-integration ]: running integration tests..."
	@mkdir -p $(COVERAGE_DIR)
	$(TEST_RUNNER) -f pkgname --rerun-fails=4 --packages "$(GO_PKGS)" -- -v -parallel 6 -tags integration -covermode=$(COVERMODE) -coverprofile $(COVERAGE_DIR)/integration.tmp $(GO_PKGS)

#
# Coverage
#
cover-clean:
	@echo "=== $(PROJECT_NAME) === [ cover-clean      ]: removing coverage files..."
	@rm -rfv $(COVERAGE_DIR)/*

cover-report:
	@echo "=== $(PROJECT_NAME) === [ cover-report     ]: generating coverage results..."
	@mkdir -p $(COVERAGE_DIR)
	@echo 'mode: $(COVERMODE)' > $(COVERAGE_DIR)/coverage.out
	@cat $(COVERAGE_DIR)/*.tmp | grep -v 'mode: $(COVERMODE)' >> $(COVERAGE_DIR)/coverage.out || true
	@$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "=== $(PROJECT_NAME) === [ cover-report     ]:     $(COVERAGE_DIR)coverage.html"

coverage-report-unit:
	@echo "=== $(PROJECT_NAME) === [ cover-report     ]: generating coverage results..."
	@mkdir -p $(COVERAGE_DIR)
	@echo 'mode: $(COVERMODE)' > $(COVERAGE_DIR)/unit.out
	@cat $(COVERAGE_DIR)/unit.tmp | grep -v 'mode: $(COVERMODE)' | grep -v -e "types.go" >> $(COVERAGE_DIR)/unit.out || true
	@$(GO) tool cover -html=$(COVERAGE_DIR)/unit.out -o $(COVERAGE_DIR)/unit.html
	@echo "=== $(PROJECT_NAME) === [ cover-report     ]:     $(COVERAGE_DIR)unit.html"

coverage-report-integration:
	@echo "=== $(PROJECT_NAME) === [ cover-report     ]: generating coverage results..."
	@mkdir -p $(COVERAGE_DIR)
	@echo 'mode: $(COVERMODE)' > $(COVERAGE_DIR)/integration.out
	@cat $(COVERAGE_DIR)/integration.tmp | grep -v 'mode: $(COVERMODE)' | grep -v -e "types.go" >> $(COVERAGE_DIR)/integration.out
	@$(GO) tool cover -html=$(COVERAGE_DIR)/integration.out -o $(COVERAGE_DIR)/integration.html
	@echo "=== $(PROJECT_NAME) === [ cover-report     ]:     $(COVERAGE_DIR)integration.html"

cover-view: cover-report
	@$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out

.PHONY: test test-only test-unit test-integration cover-report cover-view
