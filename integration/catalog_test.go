//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCatalog_Basic(t *testing.T) {
	client := Setup(t)

	project := MustCreateProject(t, client)

	_, warehouseDefault := MustCreateWarehouse(t, client, defaultProjectID)
	_, err := client.CatalogV1(context.Background(), defaultProjectID, warehouseDefault)
	require.NoError(t, err)

	_, warehouseProject := MustCreateWarehouse(t, client, project)
	_, err = client.CatalogV1(context.Background(), project, warehouseProject)
	require.NoError(t, err)
}
