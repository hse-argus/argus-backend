FROM golang:alpine AS builder

WORKDIR /argus-backend
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY docs ./docs
RUN CGO_ENABLED=0 GOOS=linux go build -o argus ./cmd/argus/main.go

CMD ["./argus"]