FROM golang:1.21 AS builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download -x
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/gochat

FROM scratch
COPY --from=builder /app/bin/gochat /
CMD [ "/gochat" ]