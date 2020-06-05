REPO = github.com/mmcloughlin/ec3

.PHONY: fmt
fmt:
	find . -name '*.go' | xargs grep -L '// Code generated' | xargs sed -i.fmtbackup '/^import (/,/)/ { /^$$/ d; }'
	find . -name '*.fmtbackup' -delete
	find . -name '*.go' | xargs grep -L '// Code generated' | xargs gofumports -w -local $(REPO)
	find . -name '*.go' | grep -v _test | xargs grep -L '// Code generated' | xargs mathfmt -w
	find . -name '*.go' | xargs grep -L '// Code generated' | xargs bib -bib docs/references.bib -w
	refs fmt -w docs/references.yml

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: generate
generate:
	make -C docs --always-make
	go generate -x ./...

.PHONY: bootstrap
bootstrap:
	go get -v -t ./...

	GO111MODULE=off go get -u github.com/myitcv/gobin
	gobin mvdan.cc/gofumpt/gofumports@v0.0.0-20200412215918-a91da47f375c
	gobin github.com/mna/pigeon@v1.0.1-0.20200224192238-18953b277063

	go install \
		./tools/mathfmt \
		./tools/refs \
		./tools/bib \
		./tools/assets
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b ${GOPATH}/bin v1.17.1
