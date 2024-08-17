FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go .

RUN go build -o /cron-vmbackup main.go

FROM victoriametrics/vmbackup:v1.96.0

WORKDIR /job

COPY --from=builder /cron-vmbackup /cron-vmbackup

RUN chmod +x /cron-vmbackup

ENTRYPOINT ["/cron-vmbackup"]