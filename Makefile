.PHONY: infrastructure build all

all: build infrastructure

infrastructure:
	@echo "Running infrastructure setup..."
	./infrastructure/run.sh

build:
	@echo "Building project..."
	./scripts/build.sh
