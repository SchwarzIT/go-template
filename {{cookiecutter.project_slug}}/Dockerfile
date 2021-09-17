# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY out/bin/{{cookiecutter.app_name}} .
USER nonroot:nonroot

ENTRYPOINT ["/{{cookiecutter.app_name}}"]
