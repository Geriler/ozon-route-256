@GOBIN=$(GOPATH)/bin

run-all:
	docker-compose up --force-recreate --build -d

cover:
	go test -cover ./cart/... && \
	go test -cover ./loms/...
