## lkctl project assignments

Get project assignments

```
lkctl project assignments PROJECT-ID [flags]
```

### Examples

```
  # Get default project assignments
  lkctl project assignments

  # Filter by assignment type
  lkctl project assignments 01986184-3cb1-7526-a98c-72fecfe97731 --relations project_admin

  # Filter by multiple assignment types
  lkctl project assignments 01986184-3cb1-7526-a98c-72fecfe97731 --relations project_admin --relations select

```

### Options

```
  -h, --help                help for assignments
  -o, --output string       Output format. One of: json|text (default "text")
      --relations strings   Filter by relations. (Can be repeated multiple times to add multiple relations, also supports comma separated relations)
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

* [lkctl project](lkctl_project.md)	 - Manage projects

