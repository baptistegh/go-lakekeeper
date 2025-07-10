package integration

import (
	"testing"

	v1 "github.com/baptistegh/go-lakekeeper/pkg/apis/v1"
	"github.com/stretchr/testify/assert"
)

func TestServer_Info(t *testing.T) {
	t.Parallel()

	client := Setup(t)

	info, _, err := client.ServerV1().Info()

	assert.NoError(t, err)
	assert.Equal(t, false, info.Bootstrapped)
}

func TestServer_Bootstrap(t *testing.T) {
	t.Parallel()

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
