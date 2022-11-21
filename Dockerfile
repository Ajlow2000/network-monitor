# Specifies a parent image
FROM golang:1.19-bullseye

# Creates an app directory to hold your app’s source code
WORKDIR /app
 
# Copy over all files
COPY * ./
 
# Install any required modules
RUN go mod download

# Builds your app with optional configuration
RUN go build -o /network-monitor
 
# Tells Docker which network port your container listens on
EXPOSE 587
 
# Specifies the executable command that runs when the container starts
ENTRYPOINT [ "/network-monitor" ]