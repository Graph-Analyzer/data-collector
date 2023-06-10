# GO Repo base repo
FROM golang:1.20-alpine3.17 as builder
RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .
# Download all the dependencies
RUN go mod download

COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -installsuffix cgo -o data-collector .

#########################################################################
# GO Repo base repo
FROM alpine:3.17 as prod
RUN apk add tzdata
RUN addgroup -S data-collector && adduser -S data-collector -G data-collector
RUN mkdir /app
WORKDIR /app/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/data-collector .

# Run as user
USER data-collector

# Run Executable
CMD ["./data-collector"]
