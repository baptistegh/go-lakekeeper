#!/usr/bin/env bash
set -euo pipefail

set -a            
source .env
set +a

CONTAINER_ENGINE="${CONTAINER_ENGINE:-docker}"
CONTAINER_COMPOSE_ENGINE="${CONTAINER_COMPOSE_ENGINE:-docker compose}"

$CONTAINER_COMPOSE_ENGINE down --volumes
LAKEKEEPER_VERSION=${LAKEKEEPER_VERSION:-latest-main} $CONTAINER_COMPOSE_ENGINE up -d

./scripts/await-healthy.sh

sleep 10

go test -v -tags integration ./integration/...
