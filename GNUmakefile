NAME=sshkey
BINARY=packer-plugin-${NAME}

COUNT?=1
TEST?=$(shell go list ./...)

.PHONY: dev

# go install github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc@latest
generate:
	@cd sshkey && go generate

build: generate
	@go build -o ${BINARY}

dev: build
	@mkdir -p ~/.packer.d/plugins/
	@cp ${BINARY} ~/.packer.d/plugins/${BINARY}

run-example: dev
	@packer build sshkey/test-fixtures/rsa.pkr.hcl

test:
	@go test -count $(COUNT) $(TEST) -timeout=3m

testacc: dev
	@rm -rf sshkey/packer_cache
	@PACKER_ACC=1 go test -count $(COUNT) -v $(TEST) -timeout=120m
