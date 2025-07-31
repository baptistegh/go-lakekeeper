## lkctl role

Manage roles

```
lkctl role [flags]
```

### Options

```
  -h, --help             help for role
  -p, --project string   Select a project (default "00000000-0000-0000-0000-000000000000")
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

* [lkctl](lkctl.md)	 - A CLI to interact with Lakekeeper's management - and Iceberg catalog APIs powered by go-iceberg.
* [lkctl role access](lkctl_role_access.md)	 - Get role access
* [lkctl role assignments](lkctl_role_assignments.md)	 - Get role assignments
* [lkctl role create](lkctl_role_create.md)	 - Create a new role
* [lkctl role delete](lkctl_role_delete.md)	 - Delete a role by id
* [lkctl role get](lkctl_role_get.md)	 - Get a role by id
* [lkctl role grant](lkctl_role_grant.md)	 - add role assignments
* [lkctl role list](lkctl_role_list.md)	 - List available roles
* [lkctl role update](lkctl_role_update.md)	 - Update role

