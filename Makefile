@GOBIN=$(GOPATH)/bin

run-all:
	docker-compose up --force-recreate --build -d

test:
	go test -cover -race ./cart/... && \
	go test -cover -race ./loms/...
