#
# Makefile fragment for Generate
#

GO           ?= go

PACKAGES ?= $(shell $(GO) list ./...)

generate: tools-compile
	@echo "=== $(PROJECT_NAME) === [ generate         ]: Running generate..."
	@for p in $(PACKAGES); do \
		echo "=== $(PROJECT_NAME) === [ generate         ]:     $$p"; \
			$(GO) generate -x $$p ; \
	done
