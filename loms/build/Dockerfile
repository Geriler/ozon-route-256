FROM golang:1.22-alpine as builder

WORKDIR /build

COPY go.mod ./

RUN go mod download

COPY . .

COPY ./config/config.yml /config.yml

RUN CGO_ENABLED=0 GOOS=linux go build -o /loms ./cmd/loms

FROM scratch
COPY --from=builder loms /bin/loms
COPY --from=builder config.yml /bin/config.yml

ENTRYPOINT ["/bin/loms", "--config=/bin/config.yml"]
