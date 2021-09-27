# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY out/bin/{{.appName}} .
USER nonroot:nonroot

ENTRYPOINT ["/{{.appName}}"]
