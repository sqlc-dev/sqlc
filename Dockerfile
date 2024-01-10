# STEP 1: Build sqlc
FROM golang:1.21.6 AS builder

COPY . /workspace
WORKDIR /workspace

ARG github_ref
ARG github_sha
ARG version
ENV GITHUB_REF=$github_ref
ENV GITHUB_SHA=$github_sha
ENV VERSION=$version
RUN go run scripts/release.go -docker

# STEP 2: Build a tiny image
FROM gcr.io/distroless/base-debian12

COPY --from=builder /workspace/sqlc /workspace/sqlc
ENTRYPOINT ["/workspace/sqlc"]
