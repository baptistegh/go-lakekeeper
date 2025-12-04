package v1_test

import (
	"net/http"
	"testing"

	managementv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1"
	permissionv1 "github.com/baptistegh/go-lakekeeper/pkg/apis/management/v1/permission"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/baptistegh/go-lakekeeper/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerService_Info(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	mux.HandleFunc("/management/v1/info", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.MustWriteHTTPResponse(t, w, "testdata/server_info.json")
	})

	info, resp, err := client.ServerV1().Info(t.Context())
	require.NoError(t, err)
	assert.NotNil(t, resp)

	want := &managementv1.ServerInfo{
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

	opts := &managementv1.BootstrapServerOptions{AcceptTermsOfUse: true}

	mux.HandleFunc("/management/v1/bootstrap", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		if !testutil.TestBodyJSON(t, r, opts) {
			t.Fatalf("error wrong body")
		}
		w.WriteHeader(http.StatusNoContent)
	})

	r, err := client.ServerV1().Bootstrap(t.Context(), opts)
	require.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, http.StatusNoContent, r.StatusCode)
}

func TestServerService_GetAllowedActions(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	opt := &managementv1.GetServerAllowedActionsOptions{
		PrincipalUser: core.Ptr("oidc~testuser"),
		PrincipalRole: core.Ptr("testrole"),
	}

	mux.HandleFunc("/management/v1/server/actions", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.MustWriteHTTPResponse(t, w, "./testdata/server_get_actions.json")
		testutil.TestParam(t, r, "principalUser", "oidc~testuser")
		testutil.TestParam(t, r, "principalRole", "testrole")
	})

	access, resp, err := client.ServerV1().GetAllowedActions(t.Context(), opt)
	require.NoError(t, err)
	assert.NotNil(t, resp)

	want := &managementv1.GetServerAllowedActionsResponse{
		AllowedActions: []permissionv1.ServerAction{
			permissionv1.CreateProject,
			permissionv1.UpdateUsers,
			permissionv1.DeleteUsers,
			permissionv1.ListUsers,
			permissionv1.ProvisionUsers,
		},
	}

	assert.Equal(t, want, access)
}
