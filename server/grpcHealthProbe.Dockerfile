ARG GO_VERSION=1.14

FROM golang:${GO_VERSION}-alpine AS builder

ENV CGO_ENABLED=0

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
# Git is required for fetching the dependencies.
RUN apk add ca-certificates git

RUN go get github.com/grpc-ecosystem/grpc-health-probe

# Final stage: the running container.
FROM alpine AS final

# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/

# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /go/bin/grpc-health-probe /usr/bin/local/grpc-health-probe

# Perform any further action as an unprivileged user.
USER nobody:nobody

# This will sleep indefinitely until SIGTERM or SIGINT occur https://stackoverflow.com/a/35770783
CMD exec /bin/sh -c "trap : TERM INT; (while true; do sleep 1000; done) & wait"