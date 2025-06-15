.PHONY: infrastructure build all

all: infrastructure build

infrastructure:
	@echo "Running infrastructure setup..."
	./infrastructure/run.sh

build:
	@echo "Building project..."
	./scripts/build.sh
