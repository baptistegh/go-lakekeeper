

FROM golang:1.25-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# print platform for debug
RUN echo "Building for $GOOS/$GOARCH"

RUN go build -o lkctl ./cmd

FROM docker.io/bitnami/minideb:bookworm

LABEL maintainer="baptiste.gouhoury@scalend.fr"

ENV HOME="/" 

RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y --no-install-recommends ca-certificates jq

RUN apt-get clean \
    && rm -rf /var/lib/apt/lists /var/cache/apt/archives

COPY --chmod=555 --from=builder /app/lkctl /usr/local/bin/lkctl

USER 1001
ENTRYPOINT ["/usr/local/bin/lkctl"]
CMD [ "--help" ]
