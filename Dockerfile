# Set default GO version 
ARG GO_VERSION=1.22.2

FROM golang:${GO_VERSION}-alpine

# Authority certificates
RUN apk add --update --no-cache ca-certificates git

# These steps will be cached
RUN mkdir /golang-setup
WORKDIR /golang-setup

# Copy go.mod and go.sum files to the workspace
COPY go.mod .
COPY go.sum .

# Get dependencies which will be cached if we donn't change mod/sum
RUN go mod download

# Copy the source code
COPY . .

# run unit tests
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go test -v -cover ./...

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/golang-setup

# Build minimal image
FROM alpine
COPY --from=0 /go/bin/golang-setup /go/bin/golang-setup
COPY --from=0 /golang-setup/store/mysql/migrations/ ./store/mysql/migrations

# Run
ENTRYPOINT [ "/go/bin/golang-setup" ]