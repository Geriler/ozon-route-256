FROM golang:1.23rc2-alpine3.20 as builder

WORKDIR /build

COPY go.mod ./

RUN go mod download

COPY . .

COPY ./config/config.yml /config.yml

RUN CGO_ENABLED=0 GOOS=linux go build -o /notifier ./cmd/notifier

FROM scratch
COPY --from=builder notifier /bin/notifier
COPY --from=builder config.yml /bin/config.yml

ENTRYPOINT ["/bin/notifier", "--config=/bin/config.yml"]
