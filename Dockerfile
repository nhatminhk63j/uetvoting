FROM golang:1.14 AS builder

# Set the Current Working Directory inside the container.
WORKDIR /app

COPY go.mod go.sum /app/

RUN set -eux; \
    GOSUMDB=off go mod download

# Coppy source and build project.
COPY . /app/
RUN go build -o uetvoting cmd/main.go

FROM ubuntu:20.04

# Install ca-certificates to use sendgrid.
RUN apt-get update -y
RUN apt-get install -y ca-certificates

# Copy timezone to run image.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
# Copy birary & migration folder.
COPY --from=builder /app/uetvoting /app/
COPY migration ./migration

EXPOSE 10080 10433

RUN chmod +x /app/uetvoting
ENTRYPOINT ["/app/uetvoting"]


