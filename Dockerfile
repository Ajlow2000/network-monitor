# https://www.docker.com/blog/developing-go-apps-docker/

# Specifies a parent image
FROM golang:1.18-alpine3.16 

# Install Ping Command
RUN apk add --update iputils

# Creates an app directory to hold your app’s source code
WORKDIR /app
 
# Copy over all go config (go.mod, go.sum etc.)
COPY go.* ./

# Install any required modules
RUN go mod download

# Copy over Go source code
COPY *.go ./
 
# Builds your app with optional configuration
RUN go build -o /network-monitor
 
# Tells Docker which network port your container listens on
EXPOSE 587
 
# Specifies the executable command that runs when the container starts
ENTRYPOINT [ "/network-monitor" ]