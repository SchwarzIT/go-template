# syntax = docker/dockerfile:1.4

# get modules, if they don't change the cache can be used for faster builds
FROM golang:1.25.10@sha256:69a7b8cc06cfb629da20d23be8f8f359a1f823481ab1a14b4e1ece153ea1d1c0 AS base
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

FROM gcr.io/distroless/static:nonroot@sha256:47b2d72ff90843eb8a768b5c2f89b40741843b639d065b9b937b07cd59b479c6 as prd
COPY --link --from=build /app/main /
# this is the numeric version of user nonroot:nonroot to check runAsNonRoot in kubernetes
USER 65532:65532
ENTRYPOINT ["/main"]
