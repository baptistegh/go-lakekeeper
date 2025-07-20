package core

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestK8sAuthSource(t *testing.T) {
	t.Parallel()

	t.Run("Provided ServiceAccountTokenPath", func(t *testing.T) {
		t.Parallel()

		path := "/tmp/service-account-go-lakekeeper-test"
		if err := os.WriteFile(path, []byte("service-account-token"), 0o600); err != nil {
			t.Fatalf("Failed to create service account token file: %v", err)
		}

		t.Cleanup(func() {
			_ = os.Remove(path)
		})

		as := K8sServiceAccountAuthSource{
			ServiceAccountTokenPath: Ptr(path),
		}

		err := as.Init(context.Background())
		assert.NoError(t, err)

		key, value, err := as.Header(context.Background())
		assert.NoError(t, err)

		assert.Equal(t, "Authorization", key)
		assert.Contains(t, "Bearer service-account-token", value)
	})

	t.Run("Default ServiceAccountTokenPath", func(t *testing.T) {
		t.Parallel()

		as := K8sServiceAccountAuthSource{}

		err := as.Init(context.Background())
		assert.Error(t, err, "failed to read service account token")

		assert.Equal(t, Ptr("/var/run/secrets/kubernetes.io/serviceaccount/token"), as.ServiceAccountTokenPath)
	})
}
