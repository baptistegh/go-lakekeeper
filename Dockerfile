# Copyright 2025 Baptiste Gouhoury <baptiste.gouhoury@scalend.fr>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.24-alpine AS builder

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
