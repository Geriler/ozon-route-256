@GOBIN=$(GOPATH)/bin

buf-generate:
	PATH=$(@GOBIN):$(PATH) $(@GOBIN)/buf dep update
	PATH=$(@GOBIN):$(PATH) $(@GOBIN)/buf generate

lint:
	$(@GOBIN)/golangci-lint run --config=./config/.golangci.yml
