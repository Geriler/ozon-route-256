build-all:
	cd cart && GOOS=linux GOARCH=amd64 make build


run-all: build-all
	docker-compose up --force-recreate --build -d

cover:
	go test -cover ./cart/...

lint:
	cd cart && golangci-lint run --config=./config/linter-config.yml
