FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /todo .

FROM alpine:3.21

RUN adduser -D -h /home/todo todo
USER todo
WORKDIR /home/todo

COPY --from=builder /todo /usr/local/bin/todo

# Default data directory inside the container
RUN mkdir -p /home/todo/.local/share/todo

VOLUME ["/home/todo/.local/share/todo"]

ENTRYPOINT ["todo"]
CMD ["daemon"]
