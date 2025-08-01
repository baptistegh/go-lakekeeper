# Go Client for Lakekeeper

The client is organized into services that correspond to different parts of the Lakekeeper API.

The two main parts are [Management](#management-api) and [Catalog](#catalog-api-iceberg-rest-catalog).

The Catalog part is handled by the Iceberg Go implementation : [go-iceberg](https://github.com/apache/iceberg-go).

### Installation

To install the client library, use `go get`:

```sh
go get github.com/baptistegh/go-lakekeeper@latest
```

This library requires Go 1.24 or later.

### Client Initialization

First, import the client package.
Then, create a new client using your authentication configurations and the base URL of your Lakekeeper instance.

#### Client Credentials (OIDC)

```go
import (
    "log"

    "golang.org/x/oauth2/clientcredentials"

    "github.com/baptistegh/go-lakekeeper/pkg/core"
    lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
    managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
)

func main() {
    // Create the OAuth configuration
    oauthConfig := &clientcredentials.Config{
        ClientID:     "your-client-id",
        ClientSecret: "your-client-secret",
        TokenURL:     "https://your-idp/oauth2-token-endpoint",
        Scopes:       []string{"lakekeeper"},
    }

    as := core.OAuthTokenSource{TokenSource: oauthConfig.TokenSource()}
    
    // Create the client and enable the initial bootstrap
    client, err := lakekeeper.NewAuthSourceClient(
        context.Background(),
        &as,
        baseURL,
        lakekeeper.WithInitialBootstrapV1Enabled(true, true, core.Ptr(managementv1.ApplicationUserType))
    )
    if err != nil {
        log.Fatalf("error creating lakekeeper client, %v", err)
    }
}
```

#### Kubernetes Service Account

```go
// This gets the service account token 
// usually stored in /var/run/secrets/kubernetes.io/serviceaccount/token
client, err := lakekeeper.NewAuthSourceClient(ctx, &core.K8sServiceAccountAuthSource{}, baseURL)
if err != nil {
    log.Fatalf("error creating lakekeeper client, %v", err)
}
```
