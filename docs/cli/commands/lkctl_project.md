## lkctl project

Manage projects

```
lkctl project [flags]
```

### Options

```
  -h, --help   help for project
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
* [lkctl project access](lkctl_project_access.md)	 - Get project access
* [lkctl project assignments](lkctl_project_assignments.md)	 - Get project assignments
* [lkctl project create](lkctl_project_create.md)	 - Create a new project
* [lkctl project delete](lkctl_project_delete.md)	 - Delete a project
* [lkctl project get](lkctl_project_get.md)	 - Get a project by id
* [lkctl project grant](lkctl_project_grant.md)	 - add project assignments
* [lkctl project list](lkctl_project_list.md)	 - List all the available projects for the current user
* [lkctl project rename](lkctl_project_rename.md)	 - Rename a project

