#
# Makefile fragment for Generate
#

GO           ?= go

generate:
	@go run tools/typegen.go -v
	@go fmt pkg/*/types.go
