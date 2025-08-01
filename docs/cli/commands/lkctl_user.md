## lkctl user

Manage users

```
lkctl user [flags]
```

### Options

```
  -h, --help   help for user
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
* [lkctl user create](lkctl_user_create.md)	 - Create a new user
* [lkctl user delete](lkctl_user_delete.md)	 - Delete a user by id
* [lkctl user get](lkctl_user_get.md)	 - Get a user by id
* [lkctl user list](lkctl_user_list.md)	 - List users

