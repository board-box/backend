FROM golang:1.24.3-alpine AS builder

WORKDIR /build

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY docs ./docs
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/main.go

FROM scratch
COPY --from=builder /bin/server /server

EXPOSE 8080

ENTRYPOINT ["/server"]