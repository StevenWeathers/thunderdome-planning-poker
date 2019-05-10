############################
# STEP 1 build the ui
############################
FROM node:10.14.1-alpine as builderNode

RUN mkdir /webapp
COPY ./src/ /webapp/src/
COPY ./public/ /webapp/public/
COPY ./*.json /webapp/
COPY ./*.js /webapp/
WORKDIR /webapp
# install node packages
RUN npm set progress=false
RUN npm install
# Build the web app
RUN npm run build
############################
# STEP 2 build executable binary
############################
FROM golang:alpine as builderGo
# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates
# Create appuser
RUN adduser -D -g '' appuser
# Create the data dir
RUN mkdir /data
# Copy the go source
COPY ./*.go $GOPATH/src/mypackage/myapp/
# Copy our static assets
COPY --from=builderNode /webapp/dist $GOPATH/src/mypackage/myapp/dist
# Set working dir
WORKDIR $GOPATH/src/mypackage/myapp/
# Fetch dependencies.
# Using go mod with go 1.11
#RUN GO111MODULE=on go mod download
# Using go get.
RUN go get -d -v
RUN go get -u github.com/gobuffalo/packr/packr
# Bundle the static assets
RUN packr
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/thunderdome
############################
# STEP 3 build a small image
############################
FROM scratch
# Import from builder.
COPY --from=builderGo /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builderGo /etc/passwd /etc/passwd
# Copy our static executable
COPY --from=builderGo /go/bin/thunderdome /go/bin/thunderdome
# Copy our data dir
COPY --from=builderGo /data /data
# Use an unprivileged user.
USER appuser

VOLUME ["/data"]

# Run the thunderdome binary.
ENTRYPOINT ["/go/bin/thunderdome"]