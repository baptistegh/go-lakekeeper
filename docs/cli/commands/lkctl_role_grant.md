## lkctl role grant

add role assignments

```
lkctl role grant ROLEID [flags]
```

### Options

```
      --assignments strings   Assignments to use; can be repeated multiple times to add multiple assignments
  -h, --help                  help for grant
      --roles strings         Grant access to roles; can be repeated multiple times to add multiple roles
      --users strings         Grant access to users; can be repeated multiple times to add multiple users
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

* [lkctl role](lkctl_role.md)	 - Manage roles

