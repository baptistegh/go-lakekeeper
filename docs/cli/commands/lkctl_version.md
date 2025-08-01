## lkctl version

Print version information

```
lkctl version [flags]
```

### Examples

```
  # Print the full version of client and server to stdout
  lkctl version

  # Print only full version of the client - no connection to server will be made
  lkctl version --client

  # Print the full version of client and server in JSON format
  lkctl version

  # Print only client and server core version strings in YAML format
  lkctl version --short
```

### Options

```
      --client          client version only (no server required)
  -h, --help            help for version
  -o, --output string   Output format. One of: json|text|short (default "text")
      --short           print just the version number
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

