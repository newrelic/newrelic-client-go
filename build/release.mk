RELEASE_SCRIPT ?= ./scripts/release.sh

REL_CMD ?= $(GOBIN)/goreleaser
DIST_DIR ?= ./dist

# Versioning info
VER_CMD  ?= $(GOBIN)/svu
VER_BUMP ?= $(GOBIN)/gobump
VER_PKG  ?= internal/version

# Technically relies on tools, but we don't want the status output
version: tools
	@echo "=== $(PROJECT_NAME) === [ version          ]: Versions:"
	@printf "Next: "
	@$(VER_CMD) next
	@printf "Tag:  "
	@$(VER_CMD) current
	@printf "Code: v"
	@$(VER_BUMP) show -r $(VER_PKG)

# Example usage: make release
release: clean tools
	@echo "=== $(PROJECT_NAME) === [ release          ]: Generating release..."
	@$(RELEASE_SCRIPT)

release-clean:
	@echo "=== $(PROJECT_NAME) === [ release-clean    ]: distribution files..."
	@rm -rfv $(DIST_DIR) $(SRCDIR)/tmp

release-build: clean tools
	@echo "=== $(PROJECT_NAME) === [ release-build    ]: Building release..."
	$(REL_CMD) build

release-package: clean tools
	@echo "=== $(PROJECT_NAME) === [ release-publish  ]: Packaging release..."
	$(REL_CMD) release --skip-publish

# Local Snapshot
snapshot: clean tools
	@echo "=== $(PROJECT_NAME) === [ snapshot         ]: Creating release snapshot..."
	@echo "=== $(PROJECT_NAME) === [ snapshot         ]:   THIS WILL NOT BE PUBLISHED!"
	@$(REL_CMD) --skip-publish --snapshot

.PHONY: release release-clean release-publish snapshot
