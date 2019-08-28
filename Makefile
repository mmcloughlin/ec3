REPO = github.com/mmcloughlin/ec3

.PHONY: fmt
fmt:
	find . -name '*.go' | xargs sed -i.fmtbackup '/^import (/,/)/ { /^$$/ d; }'
	find . -name '*.fmtbackup' -delete
	find . -name '*.go' | xargs gofumports -w -local $(REPO)
	find . -name '*.go' | xargs mathfmt -w
	find . -name '*.go' | xargs bib -bib docs/bibliography.bib -w

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: generate
generate:
	go generate -x ./...

.PHONY: bootstrap
bootstrap:
	go get -v -t ./...
	go get -u \
		mvdan.cc/gofumpt/gofumports \
		github.com/mna/pigeon
	go install \
		./tools/mathfmt \
		./tools/bib \
		./tools/assets
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b ${GOPATH}/bin v1.17.1
