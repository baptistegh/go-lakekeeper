## lkctl project access

Get project access

### Synopsis

Get project access. By default, current user's access is returned

```
lkctl project access PROJECT-ID [flags]
```

### Examples

```
  # Get default project access
  lkctl project access

  # Get specific project access
  lkctl project access 01986184-3cb1-7526-a98c-72fecfe97731

  # Get project access for a specific user
  lkctl project access 01986184-3cb1-7526-a98c-72fecfe97731 --user oidc~0198618c-5be8-7a82-a0b9-1076c9dd12f0

  # Get project access for a specific role
  lkctl project access 01986184-3cb1-7526-a98c-72fecfe97731 --role oidc~0198618c-5be8-7a82-a0b9-1076c9dd12f0

```

### Options

```
  -h, --help            help for access
  -o, --output string   Output format. One of: json|text (default "text")
      --role string     Filter by role
      --user string     Filter by user
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

