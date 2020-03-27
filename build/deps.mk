#
# Makefile fragment for installing deps
#

GO           ?= go
VENDOR_CMD   ?= ${GO} mod tidy
BUILD_DIR    ?= ./bin/

# These should be mirrored in /tools.go to keep versions consistent
GOTOOLS      += github.com/client9/misspell/cmd/misspell


tools: check-version tools-compile
	@echo "=== $(PROJECT_NAME) === [ tools            ]: Installing tools required by the project..."
	@$(GO) install $(GOTOOLS)

tools-update: check-version
	@echo "=== $(PROJECT_NAME) === [ tools-update     ]: Updating tools required by the project..."
	@$(GO) get -u $(GOTOOLS)

deps: tools deps-only


# Determine commands by looking into cmd/*
TOOL_COMMANDS   ?= $(shell find ${SRCDIR}/tools -depth 1 -type d)
# Determine binary names by stripping out the dir names
TOOL_BINS       := $(foreach tool,${TOOL_COMMANDS},$(notdir ${tool}))

tools-compile: deps-only
	@echo "=== $(PROJECT_NAME) === [ tools-compile    ]: building tools:"
	@for b in $(TOOL_BINS); do \
		echo "=== $(PROJECT_NAME) === [ tools-compile    ]:     $$b"; \
		BUILD_FILES=`find $(SRCDIR)/tools/$$b -type f -name "*.go"` ; \
		GOOS=$(GOOS) $(GO) build -ldflags=$(LDFLAGS) -o $(SRCDIR)/$(BUILD_DIR)/$$b $$BUILD_FILES ; \
	done

deps-only:
	@echo "=== $(PROJECT_NAME) === [ deps             ]: Installing package dependencies required by the project..."
	@$(VENDOR_CMD)

.PHONY: deps deps-only tools tools-update
