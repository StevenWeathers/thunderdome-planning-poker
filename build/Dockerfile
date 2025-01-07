############################
# STEP 1 build the ui
############################
FROM node:22.12 as builderNode

RUN mkdir /webapp
COPY ./ui/ /webapp/ui/
WORKDIR /webapp/ui
# install node packages
RUN npm set progress=false
RUN npm ci --no-optional
# Build the web app
RUN npm run build
############################
# STEP 2 build executable binary
############################
FROM golang:1.23-alpine as builderGo
ARG BUILD_VERSION
# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates
# Create appuser
RUN adduser -D -g '' appuser
# Copy the go source
COPY ./docs/ $GOPATH/src/github.com/stevenweathers/thunderdome-planning-poker/docs/
COPY ./internal/ $GOPATH/src/github.com/stevenweathers/thunderdome-planning-poker/internal/
# COPY ./db/ $GOPATH/src/github.com/stevenweathers/thunderdome-planning-poker/db/
COPY ./thunderdome/ $GOPATH/src/github.com/stevenweathers/thunderdome-planning-poker/thunderdome/
COPY ./*.go $GOPATH/src/github.com/stevenweathers/thunderdome-planning-poker/
COPY ./go.mod $GOPATH/src/github.com/stevenweathers/thunderdome-planning-poker/
COPY ./go.sum $GOPATH/src/github.com/stevenweathers/thunderdome-planning-poker/
COPY ./ui/*.go $GOPATH/src/github.com/stevenweathers/thunderdome-planning-poker/ui/
# Copy our static assets
COPY --from=builderNode /webapp/ui/dist $GOPATH/src/github.com/stevenweathers/thunderdome-planning-poker/ui/dist
# Set working dir
WORKDIR $GOPATH/src/github.com/stevenweathers/thunderdome-planning-poker/
# Fetch dependencies.
RUN go mod download
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -ldflags "-X main.version=$BUILD_VERSION" -o /go/bin/thunderdome
############################
# STEP 3 build a small image
############################
FROM scratch
# Import from builder.
COPY --from=builderGo /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builderGo /etc/passwd /etc/passwd
# Copy our static executable
COPY --from=builderGo /go/bin/thunderdome /go/bin/thunderdome
# Use an unprivileged user.
USER appuser

# Run the thunderdome binary.
ENTRYPOINT ["/go/bin/thunderdome"]