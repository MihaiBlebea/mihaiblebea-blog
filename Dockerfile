FROM golang:1.14-alpine AS build_base

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /tmp/go-sample-app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

# Install the go-bindata package to bundle the assets
RUN go get -u github.com/jteeuwen/go-bindata/...

RUN go mod download

COPY . .

# Unit tests
RUN CGO_ENABLED=0 go test -v

ENV HTTP_PORT=$HTTP_PORT
ENV MAILCHIMP_API_KEY=$MAILCHIMP_API_KEY
ENV MAILCHIMP_LIST_ID=$MAILCHIMP_LIST_ID

RUN env

# Build the Go app
RUN go generate

RUN go build -o ./out/go-sample-app .

# Start fresh from a smaller image
FROM alpine:3.9

RUN apk add ca-certificates

WORKDIR /app

COPY --from=build_base /tmp/go-sample-app/out/go-sample-app /app/blog

# Copy the static folder and get all the static files inside it
COPY ./static /app/static

ENV HTTP_PORT=$HTTP_PORT
ENV MAILCHIMP_API_KEY=$MAILCHIMP_API_KEY
ENV MAILCHIMP_LIST_ID=$MAILCHIMP_LIST_ID

EXPOSE $HTTP_PORT

RUN chmod +x /app/blog

# Run the binary program produced by `go install`
CMD ["/app/blog"]