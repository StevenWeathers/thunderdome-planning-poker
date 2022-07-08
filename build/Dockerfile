############################
# STEP 1 build the ui
############################
FROM node:18.5 as builderNode

RUN mkdir /webapp
COPY ./frontend/ /webapp/frontend/
COPY ./build/*.js /webapp/build/
COPY ./*.json /webapp/
COPY ./*.js /webapp/
WORKDIR /webapp
# install node packages
RUN npm set progress=false
RUN npm ci --no-optional
# Build the web app
RUN npm run build
############################
# STEP 2 build executable binary
############################
FROM golang:1.18-alpine as builderGo
ARG BUILD_VERSION
# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates
# Create appuser
RUN adduser -D -g '' appuser
# Copy the go source
COPY ./api/ $GOPATH/src/github.com/stevenweathers/thunderdome-planning-poker/api/
COPY ./db/ $GOPATH/src/github.com/stevenweathers/thunderdome-planning-poker/db/
COPY ./email/ $GOPATH/src/github.com/stevenweathers/thunderdome-planning-poker/email/
COPY ./model/ $GOPATH/src/github.com/stevenweathers/thunderdome-planning-poker/model/
COPY ./*.go $GOPATH/src/github.com/stevenweathers/thunderdome-planning-poker/
COPY ./go.mod $GOPATH/src/github.com/stevenweathers/thunderdome-planning-poker/
COPY ./go.sum $GOPATH/src/github.com/stevenweathers/thunderdome-planning-poker/
# Copy our static assets
COPY --from=builderNode /webapp/dist $GOPATH/src/github.com/stevenweathers/thunderdome-planning-poker/dist
# Set working dir
WORKDIR $GOPATH/src/github.com/stevenweathers/thunderdome-planning-poker/
# Fetch dependencies.
RUN go mod download
# Generate swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.3
RUN swag init -g api/api.go -o swaggerdocs
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