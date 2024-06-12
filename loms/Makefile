@GOBIN=$(GOPATH)/bin

vendor-proto/validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
	https://github.com/bufbuild/protovalidate vendor.protogen/tmp && \
	cd vendor.protogen/tmp && \
	git sparse-checkout set --no-cone validate &&\
	git checkout
	mkdir -p vendor.protogen/buf
	mv vendor.protogen/tmp/proto/protovalidate/buf vendor.protogen/
	rm -rf vendor.protogen/tmp

.vendor-proto: vendor-proto/validate

buf-generate:
	PATH=$(@GOBIN):$(PATH) $(@GOBIN)/buf dep update
	PATH=$(@GOBIN):$(PATH) $(@GOBIN)/buf generate

run:
	go run ./cmd/loms/loms.go --config=./config/config.yml

lint:
	$(@GOBIN)/golangci-lint run --config=./config/.golangci.yml
