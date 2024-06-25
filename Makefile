@GOBIN=$(GOPATH)/bin

run-all:
	docker-compose up --force-recreate --build -d

cover:
	go test -cover ./cart/... && \
	go test -cover ./loms/...

race:
	go test -race ./cart/... && \
	go test -race ./loms/...

test-all: cover race
