PACK := cockroach
ORG := nellisauction
PROJECT := github.com/$(ORG)/pulumi-$(PACK)
PROVIDER_PATH := provider
VERSION_PATH := $(PROVIDER_PATH)/pkg/version.Version
CODEGEN := pulumi-tfgen-$(PACK)
PROVIDER := pulumi-resource-$(PACK)
TESTPARALLELISM := 10
WORKING_DIR := $(shell pwd)

ifndef PROVIDER_VERSION
$(error PROVIDER_VERSION must be set explicitly)
endif

LDFLAGS_STRIP_SYMBOLS=-s -w
LDFLAGS_PROJ_VERSION=-X $(PROJECT)/$(VERSION_PATH)=$(PROVIDER_VERSION)
LDFLAGS=$(LDFLAGS_PROJ_VERSION) $(LDFLAGS_STRIP_SYMBOLS)

build: provider build_sdks

generate: schema generate_sdks

provider: bin/$(PROVIDER)

bin/$(PROVIDER): bin/$(CODEGEN)
	./bin/$(CODEGEN) schema --out provider/cmd/$(PROVIDER)
	cd provider && go build -o $(WORKING_DIR)/bin/$(PROVIDER) -ldflags "$(LDFLAGS)" $(PROJECT)/$(PROVIDER_PATH)/cmd/$(PROVIDER)

provider_no_deps:
	cd provider && go build -o $(WORKING_DIR)/bin/$(PROVIDER) -ldflags "$(LDFLAGS)" $(PROJECT)/$(PROVIDER_PATH)/cmd/$(PROVIDER)

schema: bin/$(CODEGEN)
	./bin/$(CODEGEN) schema --out provider/cmd/$(PROVIDER)

bin/$(CODEGEN):
	cd provider && go build -o $(WORKING_DIR)/bin/$(CODEGEN) -ldflags "$(LDFLAGS_PROJ_VERSION)" $(PROJECT)/$(PROVIDER_PATH)/cmd/$(CODEGEN)

tfgen_build_only: bin/$(CODEGEN)

generate_sdks: generate_nodejs generate_python generate_go generate_dotnet

build_sdks: build_nodejs build_python build_go build_dotnet

generate_nodejs: bin/$(CODEGEN)
	./bin/$(CODEGEN) nodejs --out sdk/nodejs/
	printf "module fake_nodejs_module // Exclude this directory from Go tools\n\ngo 1.17\n" > sdk/nodejs/go.mod

build_nodejs: generate_nodejs
	cd sdk/nodejs/ && npm install && \
		node -e "let p = require('./package.json'); \
		         p.version = '$(PROVIDER_VERSION)'; \
		         p.pulumi.version = '$(PROVIDER_VERSION)'; \
		         p.repository = { type: 'git', url: 'git+https://github.com/nellisauction/pulumi-cockroach.git' }; \
		         require('fs').writeFileSync('./package.json', JSON.stringify(p, null, 4));" && \
		npm run build && \
		cp ../../README.md ../../LICENSE package.json package-lock.json ./bin/

generate_python: bin/$(CODEGEN)
	./bin/$(CODEGEN) python --out sdk/python/
	printf "module fake_python_module // Exclude this directory from Go tools\n\ngo 1.17\n" > sdk/python/go.mod

build_python: generate_python
	cd sdk/python/ && rm -rf ./bin/ ../python.bin/ && cp -R . ../python.bin && mv ../python.bin ./bin && rm ./bin/go.mod && python3 -m venv venv && ./venv/bin/python -m pip install build==1.2.1 && cd ./bin && ../venv/bin/python -m build .

generate_go: bin/$(CODEGEN)
	./bin/$(CODEGEN) go --out sdk/go/

build_go: generate_go
	cd sdk && go build ./go/...

generate_dotnet: bin/$(CODEGEN)
	./bin/$(CODEGEN) dotnet --out sdk/dotnet/

build_dotnet: generate_dotnet
	cd sdk/dotnet/ && \
		dotnet build -c Release /p:Version=$(PROVIDER_VERSION) && \
		dotnet pack -c Release /p:Version=$(PROVIDER_VERSION) -o ./bin

test_provider:
	cd provider && go test -v ./...

clean:
	rm -rf sdk/{dotnet,nodejs,go,python} bin/* dist/*

# Cross-platform targets
provider-linux-amd64: bin/$(CODEGEN)
	cd provider && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(WORKING_DIR)/bin/linux-amd64/$(PROVIDER) -ldflags "$(LDFLAGS)" $(PROJECT)/$(PROVIDER_PATH)/cmd/$(PROVIDER)
provider-linux-arm64: bin/$(CODEGEN)
	cd provider && GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o $(WORKING_DIR)/bin/linux-arm64/$(PROVIDER) -ldflags "$(LDFLAGS)" $(PROJECT)/$(PROVIDER_PATH)/cmd/$(PROVIDER)
provider-darwin-amd64: bin/$(CODEGEN)
	cd provider && GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o $(WORKING_DIR)/bin/darwin-amd64/$(PROVIDER) -ldflags "$(LDFLAGS)" $(PROJECT)/$(PROVIDER_PATH)/cmd/$(PROVIDER)
provider-darwin-arm64: bin/$(CODEGEN)
	cd provider && GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o $(WORKING_DIR)/bin/darwin-arm64/$(PROVIDER) -ldflags "$(LDFLAGS)" $(PROJECT)/$(PROVIDER_PATH)/cmd/$(PROVIDER)
provider-windows-amd64: bin/$(CODEGEN)
	cd provider && GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o $(WORKING_DIR)/bin/windows-amd64/$(PROVIDER).exe -ldflags "$(LDFLAGS)" $(PROJECT)/$(PROVIDER_PATH)/cmd/$(PROVIDER)

provider_dist-linux-amd64: provider-linux-amd64
	mkdir -p dist
	tar -zcf dist/$(PROVIDER)-v$(PROVIDER_VERSION)-linux-amd64.tar.gz -C bin/linux-amd64 .
provider_dist-linux-arm64: provider-linux-arm64
	mkdir -p dist
	tar -zcf dist/$(PROVIDER)-v$(PROVIDER_VERSION)-linux-arm64.tar.gz -C bin/linux-arm64 .
provider_dist-darwin-amd64: provider-darwin-amd64
	mkdir -p dist
	tar -zcf dist/$(PROVIDER)-v$(PROVIDER_VERSION)-darwin-amd64.tar.gz -C bin/darwin-amd64 .
provider_dist-darwin-arm64: provider-darwin-arm64
	mkdir -p dist
	tar -zcf dist/$(PROVIDER)-v$(PROVIDER_VERSION)-darwin-arm64.tar.gz -C bin/darwin-arm64 .
provider_dist-windows-amd64: provider-windows-amd64
	mkdir -p dist
	tar -zcf dist/$(PROVIDER)-v$(PROVIDER_VERSION)-windows-amd64.tar.gz -C bin/windows-amd64 .

.PHONY: build generate provider schema tfgen_build_only generate_sdks build_sdks generate_nodejs build_nodejs generate_python build_python generate_go build_go generate_dotnet build_dotnet test_provider clean provider-linux-amd64 provider-linux-arm64 provider-darwin-amd64 provider-darwin-arm64 provider-windows-amd64 provider_dist-linux-amd64 provider_dist-linux-arm64 provider_dist-darwin-amd64 provider_dist-darwin-arm64 provider_dist-windows-amd64
