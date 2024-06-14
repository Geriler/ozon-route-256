@GOBIN=$(GOPATH)/bin

buf-generate:
	PATH=$(@GOBIN):$(PATH) $(@GOBIN)/buf dep update ../loms
	PATH=$(@GOBIN):$(PATH) $(@GOBIN)/buf generate ../loms

run:
	go run ./cmd/cart/cart.go --config=./config/config.yml

lint:
	$(@GOBIN)/golangci-lint run --config=./config/.golangci.yml
