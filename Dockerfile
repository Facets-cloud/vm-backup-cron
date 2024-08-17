FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go .

RUN go build -o /vmbackup-prod main.go

FROM victoriametrics/vmbackup:v1.96.0

WORKDIR /job

COPY --from=builder /vmbackup-prod /vmbackup-prod

RUN chmod +x /vmbackup-prod

ENTRYPOINT ["/vmbackup-prod"]