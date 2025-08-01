## lkctl warehouse list

List warehouses

```
lkctl warehouse list [flags]
```

### Examples

```
  # List warehouses
  lkctl warehouse ls
  
  # Filter by inactive status
  lkctl warehouse ls --status inactive
```

### Options

```
  -h, --help             help for list
  -o, --output string    Output format. One of: json|text|wide (default "text")
      --status strings   Filter by status. Can be repeated multiple times to filter by multiple statuses. One of: active|inactice
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

