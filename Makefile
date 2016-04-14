NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
OK_GREEN_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_CYN_COLOR=\033[33;01m


.PHONY: all clean build test release

default: help

help:
	@echo "These 'make' targets are available."
	@echo
	@echo "  run                Generates, runs the service locally in go"
	@echo "  all                Cleans, builds, runs tests"
	@echo "  clean              Removes all build output"
	@echo "  clean-all          Remove all build output and generated code"
	@echo "  generate           Generates both server and client"
	@echo "  build              generates swagger code and rebuilds the service only"
	@echo "  test               Run the unit tests"
	@echo "  coverage           Run the unit tests and produces a coverage report"
	@echo "  tools              Installs tools needed to run"
	@echo


run:	generate
	godep go run cmd/catalog-service-manager/catalog-service-manager.go


all: 	clean-all build test

clean:
	@echo "$(OK_COLOR)==> Removing build artifacts$(NO_COLOR)"
	rm -rf ${GOBIN}/catalog-service-manager
	rm -rf bin

clean-all: clean
	@echo "$(OK_COLOR)==> Removing generated code$(NO_COLOR)"
	rm -rf generated


generate:
	@echo "$(OK_COLOR)==> Generating code $(NO_COLOR)"
	scripts/generate-server.sh

coverage:
	@echo "$(OK_COLOR)==> Running tests with coverage tool$(NO_COLOR)"
	./scripts/testCoverage.sh

build:	generate
	@echo "$(OK_COLOR)==> Building Catalog Service Manager code $(NO_COLOR)"
	cd cmd/catalog-service-manager;\
        godep go install .
        
test:
	@echo "$(OK_COLOR)==> Running tests $(NO_COLOR)"
	godep go test ./... | grep -v generated | grep -v cmd/catalog-service-manager/handlers | grep -v scripts | grep -v examples

tools:
	@echo "$(OK_COLOR)==> Installing tools and go dependancies $(NO_COLOR)"
	go get golang.org/x/tools/cmd/cover
	go get github.com/tools/godep

	./scripts/tools/codegen.sh
