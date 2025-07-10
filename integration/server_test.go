//go:build integration
// +build integration

package integration

import (
	"testing"

	v1 "github.com/baptistegh/go-lakekeeper/pkg/apis/v1"
	"github.com/stretchr/testify/assert"
)

func TestServer_Info(t *testing.T) {
	client := Setup(t)

	info, _, err := client.ServerV1().Info()

	assert.NoError(t, err)
	assert.Equal(t, true, info.Bootstrapped)
	assert.NotEmpty(t, true, info.ServerID)
	assert.NotEmpty(t, true, info.DefaultProjectID)
	assert.Equal(t, info.AuthzBackend, "openfga")
}

func TestServer_Bootstrap(t *testing.T) {
	client := Setup(t)

	resp, err := client.ServerV1().Bootstrap(&v1.BootstrapServerOptions{
		AcceptTermsOfUse: true,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	info, _, err := client.ServerV1().Info()

	assert.NoError(t, err)
	assert.Equal(t, true, info.Bootstrapped)
}
