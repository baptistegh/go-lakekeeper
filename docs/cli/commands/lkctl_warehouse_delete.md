## lkctl warehouse delete

delete a warehouse by id

```
lkctl warehouse delete WAREHOUSEID [flags]
```

### Examples

```
  # delete a warehouse by id
  lkctl warehouse delete 019861a0-6d4e-7bf3-96c6-9aef2d4a2749
  
  # force delete a warehouse
  lkctl warehouse rm 019861a0-6d4e-7bf3-96c6-9aef2d4a2749 --force
```

### Options

```
      --force   Force delete the warehouse
  -h, --help    help for delete
```

### Options inherited from parent commands

```
      --auth-url string        OAuth2 token endpoint; set this or LAKEKEEPER_AUTH_URL environment variable (default "http://localhost:30080/realms/iceberg/protocol/openid-connect/token")
      --bootstrap              If set to true, the CLI will try to bootstrap the server with the current user first; set this or LAKEKEEPER_BOOTSTRAP environment variable
      --client-id string       OAuth2 client_id; set this or LAKEKEEPER_CLIENT_ID environment variable (default "lakekeeper-admin")
      --client-secret string   OAuth2 client_secret; set this or LAKEKEEPER_CLIENT_SECRET environment variable (default "KNjaj1saNq5yRidVEMdf1vI09Hm0pQaL")
      --debug                  Enable debug mode
  -p, --project string         Select a project (default "00000000-0000-0000-0000-000000000000")
      --scopes strings         OAuth2 scopes; set this or LAKEKEEPER_SCOPE environment variable (default [lakekeeper])
      --server string          Lakekeeper base URL; set this or LAKEKEEPER_SERVER environment variable (default "http://localhost:8181")
```

### SEE ALSO

* [lkctl warehouse](lkctl_warehouse.md)	 - Manage warehouses

