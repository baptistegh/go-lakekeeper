

FROM bitnami/minideb:bookworm

ARG TARGETPLATFORM
ARG TARGETOS
ARG TARGETARCH

ENV HOME="/"

RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y --no-install-recommends ca-certificates jq

RUN apt-get clean \
    && rm -rf /var/lib/apt/lists /var/cache/apt/archives

COPY --chmod=555 $TARGETPLATFORM/lkctl-${TARGETOS}-${TARGETARCH} /usr/local/bin/lkctl

USER 1001
ENTRYPOINT ["/usr/local/bin/lkctl"]
CMD [ "--help" ]
