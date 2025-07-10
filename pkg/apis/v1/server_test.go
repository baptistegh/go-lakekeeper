package v1_test

import (
	"net/http"
	"testing"

	v1 "github.com/baptistegh/go-lakekeeper/pkg/apis/v1"
	"github.com/baptistegh/go-lakekeeper/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestServerService_Info(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	mux.HandleFunc("/management/v1/info", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.MustWriteHTTPResponse(t, w, "testdata/server_info.json")
	})

	info, resp, err := client.ServerV1().Info()
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	want := &v1.ServerInfo{
		AuthzBackend:                 "openfga",
		Bootstrapped:                 true,
		DefaultProjectID:             "01f2fdfc-81fc-444d-8368-5b6701566e35",
		AWSSystemIdentitiesEnabled:   false,
		AzureSystemIdentitiesEnabled: false,
		GCPSystemIdentitiesEnabled:   false,
		ServerID:                     "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d",
		Version:                      "v0.9.0",
		Queues:                       []string{"string"},
	}

	assert.Equal(t, want, info)
}

func TestServerService_Bootstrap(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	opts := &v1.BootstrapServerOptions{AcceptTermsOfUse: true}

	mux.HandleFunc("/management/v1/bootstrap", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		if !testutil.TestBodyJSON(t, r, opts) {
			t.Fatalf("error wrong body")
		}
		w.WriteHeader(http.StatusNoContent)
	})

	r, err := client.ServerV1().Bootstrap(opts)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)
}
