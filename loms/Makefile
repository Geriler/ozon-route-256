@GOBIN=$(GOPATH)/bin

buf-generate:
	PATH=$(@GOBIN):$(PATH) $(@GOBIN)/buf dep update
	PATH=$(@GOBIN):$(PATH) $(@GOBIN)/buf generate

run:
	go run ./cmd/loms/loms.go --config=./config/config.yml

lint:
	$(@GOBIN)/golangci-lint run --config=./config/.golangci.yml
