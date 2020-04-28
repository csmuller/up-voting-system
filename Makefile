PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

NODE_NAME := pbbd
ACLI_NAME := acli
VCLI_NAME := vcli

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=NameService \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=$(NODE_NAME) \
	-X github.com/cosmos/cosmos-sdk/version.VClientName=$(VCLI_NAME) \
	-X github.com/cosmos/cosmos-sdk/version.AClientName=$(ACLI_NAME) \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) 

BUILD_FLAGS := -ldflags '$(ldflags)'

all: install

install: go.sum
		@echo "--> Build and install pbbd."
		go install $(BUILD_FLAGS) ./cli/$(NODE_NAME)
		@echo "--> Build and install acli."
		go install $(BUILD_FLAGS) ./cli/$(ACLI_NAME)
		@echo "--> Build and install vcli."
		go install $(BUILD_FLAGS) ./cli/$(VCLI_NAME)

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

test:
	@go test $(PACKAGES)
