@GOBIN=$(GOPATH)/bin

buf-generate:
	PATH=$(@GOBIN):$(PATH) $(@GOBIN)/buf dep update ../loms
	PATH=$(@GOBIN):$(PATH) $(@GOBIN)/buf generate ../loms

lint:
	$(@GOBIN)/golangci-lint run --config=./config/.golangci.yml
