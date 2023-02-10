# syntax = docker/dockerfile:1.4

# get modules, if they don't change the cache can be used for faster builds
FROM golang:1.19@sha256:09009d430dbf3fb80a4731546eefb1603f1fdb36fe6354eab138d8275a85e89f AS base
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /src
COPY go.* .
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# build th application
FROM base AS build
# temp mount all files instead of loading into image with COPY
# temp mount module cache
# temp mount go build cache
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-w -s" -o /app/main ./cmd/{{.Base.appName}}/*.go

# Import the binary from build stage

FROM gcr.io/distroless/static:nonroot@sha256:5c9728ecd8683960ffff84837b21ac9f8f1fabbe437221a486fb93ccc4e05d88 as prd
COPY --link --from=build /app/main /
# this is the numeric version of user nonroot:nonroot to check runAsNonRoot in kubernetes
USER 65532:65532
ENTRYPOINT ["/main"]
