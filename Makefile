.PHONY: default
.PHONY: build
.PHONY: hsm
.PHONY: test
.PHONY: testclean
.PHONY: vet
.PHONY: staticcheck
.PHONY: lint
.PHONY: clean
.PHONY: count

default: vet lint staticcheck test build

build: hsm

hsm:
	@cd cmd/$@ && go build -o ../../bin/$@

test:
	@echo "*** $@"
	@go test ./...

testclean:
	@go clean -testcache

vet:
	@echo "*** $@"
	@go vet ./...

staticcheck:
	@staticcheck ./...

lint:
	@echo "*** $@"
	@revive ./...

clean:
	@rm -rf bin

count:
	@gocloc .

install-deps:
	@go install github.com/mgechev/revive@latest
