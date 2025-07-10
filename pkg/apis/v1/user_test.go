package v1_test

import (
	"net/http"
	"testing"

	v1 "github.com/baptistegh/go-lakekeeper/pkg/apis/v1"
	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/baptistegh/go-lakekeeper/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestUserService_Get(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	userID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	mux.HandleFunc("/management/v1/user/"+userID, func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.MustWriteHTTPResponse(t, w, "testdata/get_user.json")
	})

	user, resp, err := client.UserV1().Get(userID)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	want := &v1.User{
		ID:              userID,
		Name:            "test-user",
		Email:           core.Ptr("test@example.com"),
		UserType:        v1.HumanUserType,
		CreatedAt:       "2019-08-24T14:15:22Z",
		UpdatedAt:       core.Ptr("2019-08-24T14:15:22Z"),
		LastUpdatedWith: "create-endpoint",
	}

	assert.Equal(t, want, user)
}

func TestUserService_Whoami(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	mux.HandleFunc("/management/v1/whoami", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodGet)
		testutil.MustWriteHTTPResponse(t, w, "testdata/get_user.json")
	})

	user, resp, err := client.UserV1().Whoami()
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	want := &v1.User{
		ID:              "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d",
		Name:            "test-user",
		Email:           core.Ptr("test@example.com"),
		UserType:        v1.HumanUserType,
		CreatedAt:       "2019-08-24T14:15:22Z",
		UpdatedAt:       core.Ptr("2019-08-24T14:15:22Z"),
		LastUpdatedWith: "create-endpoint",
	}

	assert.Equal(t, want, user)
}

func TestUserService_Provision(t *testing.T) {
	t.Parallel()

	mux, client := testutil.ServerMux(t)

	opts := v1.ProvisionUserOptions{
		ID:             core.Ptr("a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"),
		Email:          core.Ptr("test@example.com"),
		Name:           core.Ptr("test-user"),
		UserType:       core.Ptr(v1.HumanUserType),
		UpdateIfExists: core.Ptr(true),
	}

	mux.HandleFunc("/management/v1/user", func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodPost)
		if !testutil.TestBodyJSON(t, r, &opts) {
			t.Fatalf("error wrong body")
		}
		w.WriteHeader(http.StatusCreated)
		testutil.MustWriteHTTPResponse(t, w, "testdata/get_user.json")
	})

	want := &v1.User{
		ID:              "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d",
		Email:           core.Ptr("test@example.com"),
		Name:            "test-user",
		UserType:        v1.HumanUserType,
		CreatedAt:       "2019-08-24T14:15:22Z",
		UpdatedAt:       core.Ptr("2019-08-24T14:15:22Z"),
		LastUpdatedWith: "create-endpoint",
	}

	user, resp, err := client.UserV1().Provision(&opts)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	assert.Equal(t, want, user)
}

func TestUserService_Delete(t *testing.T) {
	t.Parallel()
	mux, client := testutil.ServerMux(t)

	userID := "a4b2c1d0-e3f4-5a6b-7c8d-9e0f1a2b3c4d"

	mux.HandleFunc("/management/v1/user/"+userID, func(w http.ResponseWriter, r *http.Request) {
		testutil.TestMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.UserV1().Delete(userID)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
