# Authentication

You can authenticate to Lakekeeper using the CLI with both `flags` or `env` variables.

Eg. with flags:

```sh
lkctl info \
    --server http://localhost:8181 \
    --auth-url http://localhost:30080/realms/iceberg/protocol/openid-connect/token \
    --client-id spark \
    --client-secret 2OR3eRvYfSZzzZ16MlPd95jhLnOaLM \
    --scope lakekeeper
```

Eg. with environment variables:

```sh
export LAKEKEEPER_SERVER=http://localhost:8181
export LAKEKEEPER_AUTH_URL=http://localhost:30080/realms/iceberg/protocol/openid-connect/token
export LAKEKEEPER_CLIENT_ID=spark
export LAKEKEEPER_CLIENT_SECRET=2OR3eRvYfSZzzZ16MlPd95jhLnOaLM
export LAKEKEEPER_SCOPE=lakekeeper

lkctl info
```

You can also set these variables in a `.env` file.