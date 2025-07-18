# Lakekeeper Go Client

[![Go Report Card](https://goreportcard.com/badge/github.com/baptistegh/go-lakekeeper)](https://goreportcard.com/report/github.com/baptistegh/go-lakekeeper)
[![GoDoc](https://godoc.org/github.com/baptistegh/go-lakekeeper?status.svg)](https://godoc.org/github.com/baptistegh/go-lakekeeper)
[![codecov](https://codecov.io/gh/baptistegh/go-lakekeeper/graph/badge.svg?token=2WF3AB10RA)](https://codecov.io/gh/baptistegh/go-lakekeeper)

Go Client for [Lakekeeper API](https://docs.lakekeeper.io).

It provides a convenient way to interact with Lakekeeper services from your Go applications.

## Installation

To install the client library, use `go get`:

```sh
go get github.com/baptistegh/go-lakekeeper
```

## Usage

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
)

func main() {
    // Create the OAuth configuration
    oauthConfig := &clientcredentials.Config{
        ClientID:     "lakekeeper-client-id",
        ClientSecret: "lakekeeper-client-secret",
        TokenURL:     "<oidc_provider>/oauth/token",
        Scopes:       []string{"lakekeeper"},
    }

    as := core.OAuthTokenSource{TokenSource: oauthConfig.TokenSource()}
    
    // Create the client and enable the initial bootstrap
    client, err := lakekeeper.NewAuthSourceClient(
        &as,
        baseURL,
        lakekeeper.WithInitialBootstrapV1Enabled(true, true, core.Ptr(managementv1.ApplicationUserType))
    )
    if err != nil {
        log.Fatalf("error creating lakekeeper client, %v", err)
    }

    // You can now use the client to interact with the API
}
```

#### Kubernetes Service Account

```go
import (
    "log"

    "github.com/baptistegh/go-lakekeeper/pkg/core"
    lakekeeper "github.com/baptistegh/go-lakekeeper/pkg/client"
)

func main() {
    as := core.K8sServiceAccountAuthSource{}

    client, err := lakekeeper.NewAuthSourceClient(&as, baseURL)
    if err != nil {
        log.Fatalf("error creating lakekeeper client, %v", err)
    }

    // You can now use the client to interact with the API
}
```

### Accessing API Services

The client is organized into services that correspond to different parts of the Lakekeeper API.

#### Get Server Information

You can get information about the Lakekeeper server instance.

```go
serverInfo, _, err := client.ServerV1().Info()
if err != nil {
    log.Fatalf("Failed to get server info: %v", err)
}

log.Printf("Connected to Lakekeeper version: %s\n", serverInfo.Version)
```

#### Working with Projects

```go
// Get default project
project, _, err := client.ProjectV1().GetDefault()
if err != nil {
    log.Fatalf("Failed to get project: %v", err)
}

fmt.Printf("Default Project ID: %s, Name: %s\n", project.ID, project.Name)
```

#### Working with Project-Scoped Resources (e.g., Roles)

Services for resources like Roles and Warehouses are scoped to a specific project.
You first create a service for that project ID.

```go
// Get a specific role within a project
role, _, err := client.RoleV1(project.ID).Get("a-role-id")
if err != nil {
    log.Fatalf("Failed to get role: %v", err)
}
fmt.Printf("Role Name: %s\n", role.Name)
```

#### Create resources (e.g., Warehouse)

```go
// Set the storage settings (eg. MinIO)
storage, _ := profile.NewS3StorageSettings("bucket-name", "local-01",
    profile.WithEndpoint("http://minio:9000/"),
    profile.WithPathStyleAccess()
)

creds, _ := credential.NewS3CredentialAccessKey("access-key-id", "secret-access-key")

opts := v1.CreateWarehouseOptions{
    Name:              acctest.RandString(8),
    StorageProfile:    storage.AsProfile(),
    StorageCredential: creds.AsCredential(),
    DeleteProfile:     profile.NewTabularDeleteProfileHard().AsProfile(),
}

// Create the warehouse within a project
warehouse, _, err := client.WarehouseV1(project.ID).Create(&opts)
if err != nil {
    log.Fatalf("Failed to create warehouse: %v", err)
}

fmt.Printf("Warehouse with ID %s created!\n", warehouse.ID)
```
