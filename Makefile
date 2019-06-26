REPO = github.com/mmcloughlin/ec3

.PHONY: fmt
fmt:
	find . -name '*.go' | xargs gofumports -w -local $(REPO)
	find . -name '*.go' | xargs mathfmt -w

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: bootstrap
bootstrap:
	go get mvdan.cc/gofumpt/gofumports
	go install ./tools/mathfmt
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b ${GOPATH}/bin v1.17.1
