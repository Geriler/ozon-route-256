@GOBIN=$(GOPATH)/bin

run-all:
	docker-compose up --force-recreate --build -d && \
	cd migrations && \
    $(@GOBIN)/goose postgres "postgres://postgres:postgres@localhost:5432/postgres" up

cover:
	go test -cover ./cart/... && \
	go test -cover ./loms/...
