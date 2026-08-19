BINARY = terraform-provider-disassembly
VERSION ?= 0.1.0

build:
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY)

tidy:
	go mod tidy

# Install into the local plugin dir for dev overrides (see README).
install: build
	@echo "Add the dev_overrides block from README.md to ~/.terraformrc, then run terraform plan."

.PHONY: build tidy install
