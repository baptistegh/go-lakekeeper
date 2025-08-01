## lkctl user delete

Delete a user by id

```
lkctl user delete USERID [flags]
```

### Examples

```
  # Delete a user
  lkctl user delete oidc~d223d88c-85b6-4859-b5c5-27f3825e47f6
```

### Options

```
  -h, --help   help for delete
```

### Options inherited from parent commands

```
      --auth-url string        OAuth2 token endpoint; set this or LAKEKEEPER_AUTH_URL environment variable (default "http://localhost:30080/realms/iceberg/protocol/openid-connect/token")
      --bootstrap              If set to true, the CLI will try to bootstrap the server with the current user first; set this or LAKEKEEPER_BOOTSTRAP environment variable
      --client-id string       OAuth2 client_id; set this or LAKEKEEPER_CLIENT_ID environment variable (default "lakekeeper-admin")
      --client-secret string   OAuth2 client_secret; set this or LAKEKEEPER_CLIENT_SECRET environment variable (default "KNjaj1saNq5yRidVEMdf1vI09Hm0pQaL")
      --debug                  Enable debug mode
      --scopes strings         OAuth2 scopes; set this or LAKEKEEPER_SCOPE environment variable (default [lakekeeper])
      --server string          Lakekeeper base URL; set this or LAKEKEEPER_SERVER environment variable (default "http://localhost:8181")
```

### SEE ALSO

* [lkctl user](lkctl_user.md)	 - Manage users

