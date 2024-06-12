@GOBIN=$(GOPATH)/bin

run:
	go run ./cmd/cart/cart.go --config=./config/config.yml

lint:
	$(@GOBIN)/golangci-lint run --config=./config/.golangci.yml
