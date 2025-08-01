## lkctl

A CLI to interact with Lakekeeper's management - and Iceberg catalog APIs powered by go-iceberg.

```
lkctl [flags]
```

### Options

```
      --auth-url string        OAuth2 token endpoint; set this or LAKEKEEPER_AUTH_URL environment variable (default "http://localhost:30080/realms/iceberg/protocol/openid-connect/token")
      --bootstrap              If set to true, the CLI will try to bootstrap the server with the current user first; set this or LAKEKEEPER_BOOTSTRAP environment variable
      --client-id string       OAuth2 client_id; set this or LAKEKEEPER_CLIENT_ID environment variable (default "lakekeeper-admin")
      --client-secret string   OAuth2 client_secret; set this or LAKEKEEPER_CLIENT_SECRET environment variable (default "KNjaj1saNq5yRidVEMdf1vI09Hm0pQaL")
      --debug                  Enable debug mode
  -h, --help                   help for lkctl
      --scopes strings         OAuth2 scopes; set this or LAKEKEEPER_SCOPE environment variable (default [lakekeeper])
      --server string          Lakekeeper base URL; set this or LAKEKEEPER_SERVER environment variable (default "http://localhost:8181")
```

### SEE ALSO

* [lkctl catalog](lkctl_catalog.md)	 - Interacts with catalogs (not implemented)
* [lkctl project](lkctl_project.md)	 - Manage projects
* [lkctl role](lkctl_role.md)	 - Manage roles
* [lkctl server](lkctl_server.md)	 - Manage server
* [lkctl user](lkctl_user.md)	 - Manage users
* [lkctl version](lkctl_version.md)	 - Print version information
* [lkctl warehouse](lkctl_warehouse.md)	 - Manage warehouses
* [lkctl whoami](lkctl_whoami.md)	 - Print the current user

