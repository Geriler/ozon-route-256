FROM golang:1.22-alpine as builder

WORKDIR /build

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /cart ./cmd/cart

FROM scratch
COPY --from=builder cart /bin/cart

ENTRYPOINT ["/bin/cart"]
