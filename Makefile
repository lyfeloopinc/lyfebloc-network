BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')
APPNAME := lyfeblocnetwork

# do not override user values
ifeq (,$(VERSION))
  VERSION := $(shell git describe --exact-match 2>/dev/null)
  # if VERSION is empty, then populate it with branch name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif

# Update the ldflags with the app, client & server names
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=$(APPNAME) \
	-X github.com/cosmos/cosmos-sdk/version.AppName=$(APPNAME)d \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'

##############
###  Test  ###
##############

test-unit:
	@echo Running unit tests...
	@go test -mod=readonly -v -timeout 30m ./...

test-race:
	@echo Running unit tests with race condition reporting...
	@go test -mod=readonly -v -race -timeout 30m ./...

test-cover:
	@echo Running unit tests and creating coverage report...
	@go test -mod=readonly -v -timeout 30m -coverprofile=$(COVER_FILE) -covermode=atomic ./...
	@go tool cover -html=$(COVER_FILE) -o $(COVER_HTML_FILE)
	@rm $(COVER_FILE)

bench:
	@echo Running unit tests with benchmarking...
	@go test -mod=readonly -v -timeout 30m -bench=. ./...

test: govet govulncheck test-unit

.PHONY: test test-unit test-race test-cover bench

#################
###  Install  ###
#################

all: install

install:
	@echo "--> ensure dependencies have not been modified"
	@go mod verify
	@echo "--> installing $(APPNAME)d"
	@go install $(BUILD_FLAGS) -mod=readonly ./cmd/$(APPNAME)d

.PHONY: all install

##################
###  Protobuf  ###
##################

# Use this target if you do not want to use Ignite for generating proto files

proto-deps:
	@echo "Installing proto deps"
	@echo "Proto deps present, run 'go tool' to see them"

proto-gen:
	@echo "Generating protobuf files..."
	@ignite generate proto-go --yes

.PHONY: proto-gen

#################
###  Linting  ###
#################

lint:
	@echo "--> Running linter"
	@go tool github.com/golangci/golangci-lint/cmd/golangci-lint run ./... --timeout 15m

lint-fix:
	@echo "--> Running linter and fixing issues"
	@go tool github.com/golangci/golangci-lint/cmd/golangci-lint run ./... --fix --timeout 15m

.PHONY: lint lint-fix

###################
### Development ###
###################

govet:
	@echo Running go vet...
	@go vet ./...

govulncheck:
	@echo Running govulncheck...
	@go tool golang.org/x/vuln/cmd/govulncheck@latest
	@govulncheck ./...

.PHONY: govet govulncheck

###################
###   Docker    ###
###################

NETWORK ?= devnet
DOCKER_COMPOSE = docker/docker-compose.yml

ifeq ($(NETWORK),testnet)
  VALIDATOR_COUNT := 10
else ifeq ($(NETWORK),mainnet)
  VALIDATOR_COUNT := 30
else
  VALIDATOR_COUNT := 2
endif

build:
	@echo "🔨 Building Lyfebloc Network binary..."
	GOCACHE=$(PWD)/.gocache go build -o ./build/lyfebloc-networkd ./cmd/lyfebloc-networkd
	@echo "✅ Build complete: ./build/lyfebloc-networkd"

init:
	@echo "🌱 Initializing $(NETWORK)..."
	scripts/setup-network.sh $(NETWORK)

docker-build:
	@echo "🐳 Building Docker images..."
	NETWORK=$(NETWORK) docker compose -f $(DOCKER_COMPOSE) --profile $(NETWORK) build
	@echo "✅ Docker images built."

docker-up:
	@echo "🚀 Starting Lyfebloc $(NETWORK) with Docker..."
	@mkdir -p data/validators/$(NETWORK)
	@for i in $(shell seq 1 $(VALIDATOR_COUNT)); do mkdir -p data/validators/$(NETWORK)/validator$$i; done
	NETWORK=$(NETWORK) docker compose -f $(DOCKER_COMPOSE) --profile $(NETWORK) up -d
	@echo "🌍 Network running: http://localhost:26657 | Grafana http://localhost:3000"

docker-logs:
	@echo "📜 Streaming validator logs..."
	NETWORK=$(NETWORK) docker compose -f $(DOCKER_COMPOSE) logs -f validator1

docker-down:
	@echo "🛑 Stopping Docker network..."
	NETWORK=$(NETWORK) docker compose -f $(DOCKER_COMPOSE) --profile $(NETWORK) down

docker-clean:
	@echo "🧹 Removing containers, images, and volumes..."
	NETWORK=$(NETWORK) docker compose -f $(DOCKER_COMPOSE) --profile $(NETWORK) down -v --rmi all

status:
	@echo "🌐 Querying chain status..."
	curl -s http://localhost:26657/status | jq '.result.node_info.network, .result.sync_info.latest_block_height'

help:
	@echo "Available targets:"
	@echo "  build          - Build the Lyfebloc binary"
	@echo "  init           - Initialize the $(NETWORK) network"
	@echo "  docker-build   - Build Docker images"
	@echo "  docker-up      - Launch the full stack"
	@echo "  docker-logs    - Tail validator logs"
	@echo "  docker-down    - Stop the Docker stack"
	@echo "  docker-clean   - Purge all Docker resources"
	@echo "  status         - Check RPC node status"

.PHONY: build init docker-build docker-up docker-logs docker-down docker-clean status help
