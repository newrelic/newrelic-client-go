#
# Makefile fragment for Linting
#

GO           ?= go
MISSPELL     ?= misspell
GOFMT        ?= gofmt
GOIMPORTS    ?= goimports

GOLINTER      = golangci-lint

EXCLUDEDIR      ?= .git
SRCDIR          ?= .
GO_PKGS         ?= $(shell ${GO} list ./... | grep -v -e "/vendor/" -e "/example")
FILES           ?= $(shell find ${SRCDIR} -type f | grep -v -e '.git/' -e '/vendor/' -e 'go.sum')
GO_FILES        ?= $(shell find $(SRCDIR) -type f -name "*.go" | grep -v -e ".git/" -e '/vendor/' -e '/example/')
PROJECT_MODULE  ?= $(shell $(GO) list -m)

GO_MOD_OUTDATED ?= go-mod-outdated

lint: deps spell-check gofmt lint-commit golangci goimports outdated
lint-fix: deps spell-check-fix gofmt-fix goimports

#
# Check spelling on all the files, not just source code
#
spell-check: deps
	@echo "=== $(PROJECT_NAME) === [ spell-check      ]: Checking for spelling mistakes with $(MISSPELL)..."
	@$(MISSPELL) -source text $(FILES)

spell-check-fix: deps
	@echo "=== $(PROJECT_NAME) === [ spell-check-fix  ]: Fixing spelling mistakes with $(MISSPELL)..."
	@$(MISSPELL) -source text -w $(FILES)

gofmt: deps
	@echo "=== $(PROJECT_NAME) === [ gofmt            ]: Checking file format with $(GOFMT)..."
	@find . -path "$(EXCLUDEDIR)" -prune -print0 | xargs -0 $(GOFMT) -e -l -s -d ${SRCDIR}

gofmt-fix: deps
	@echo "=== $(PROJECT_NAME) === [ gofmt-fix        ]: Fixing file format with $(GOFMT)..."
	@find . -path "$(EXCLUDEDIR)" -prune -print0 | xargs -0 $(GOFMT) -e -l -s -w ${SRCDIR}

goimports: deps
	@echo "=== $(PROJECT_NAME) === [ goimports        ]: Checking imports with $(GOIMPORTS)..."
	@$(GOIMPORTS) -w -local $(PROJECT_MODULE) $(GO_FILES)

golangci: deps
	@echo "=== $(PROJECT_NAME) === [ golangci-lint    ]: Linting using $(GOLINTER) ($(COMMIT_LINT_CMD))..."
	@$(GOLINTER) run

outdated: deps tools-outdated
	@echo "=== $(PROJECT_NAME) === [ outdated         ]: Finding outdated deps with $(GO_MOD_OUTDATED)..."
	@$(GO) list -u -m -json all | $(GO_MOD_OUTDATED) -direct -update

.PHONY: lint spell-check spell-check-fix gofmt gofmt-fix lint-fix lint-commit outdated goimports
