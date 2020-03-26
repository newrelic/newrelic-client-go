#
# Makefile fragment for Generate
#

GO           ?= go

generate:
	@go run cmd/typegen.go -v
	@go fmt pkg/*/types.go
