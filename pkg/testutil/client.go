// Copyright 2025 Baptiste Gouhoury <baptiste.gouhoury@scalend.fr>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testutil

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/baptistegh/go-lakekeeper/pkg/client"
	"github.com/hashicorp/go-retryablehttp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setup sets up a test HTTP server along with a Client that is
// configured to talk to that test server.  Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func ServerMux(t *testing.T) (*http.ServeMux, *client.Client) {
	t.Helper()
	// mux is the HTTP request multiplexer used with the test server.
	mux := http.NewServeMux()

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	c, err := client.NewClient(t.Context(), "", server.URL)
	if err != nil {
		t.Fatal(err)
	}

	return mux, c
}

func TestURL(t *testing.T, r *http.Request, want string) {
	if got := r.RequestURI; got != want {
		t.Errorf("Request url: %+v, want %s", got, want)
	}
}

func TestMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %s, want %s", got, want)
	}
}

func TestHeader(t *testing.T, r *http.Request, key, want string) {
	if got := r.Header.Get(key); got != want {
		t.Errorf("Request header: %s, want %s", got, want)
	}
}

// Tests that a given form attribute has a value in a form request. Useful
// for testing file upload API requests.
func TestFormBody(t *testing.T, r *http.Request, key, want string) {
	if got := r.FormValue(key); got != want {
		t.Errorf("Request body for key %s got: %s, want %s", key, got, want)
	}
}

// testBodyJSON tests that the JSON request body is what we expect. The want
// argument is typically either a struct, a map[string]string, or a
// map[string]any, though other types are handled as well.
//
// Calls t.Fatal if decoding the request body fails, failing the test
// immediately.
//
// When the request body is not equal to "want", the error is reported but the
// test is allowed to continue. You can use the return value to end the test on
// error: returns true if the decoded body is identical to want, false
// otherwise.
func TestBodyJSON[T any](t *testing.T, r *http.Request, want T) bool {
	var got T

	if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
		t.Fatalf("Failed to decode JSON from request body: %v", err)
	}

	return assert.Equal(t, want, got)
}

// testParam checks whether the given request contains the expected parameter and whether the parameter has the expected value.
func TestParam(t *testing.T, r *http.Request, key, value string) {
	require.True(t, r.URL.Query().Has(key), "Request does not contain the %q parameter", key)
	assert.Len(t, r.URL.Query()[key], 1, "Request contains multiple %q parameters when only one is expected", key)
	require.Equal(t, value, r.URL.Query().Get(key))
}

func MustWriteHTTPResponse(t *testing.T, w io.Writer, fixturePath string) {
	t.Helper()
	f, err := os.Open(fixturePath)
	if err != nil {
		t.Fatalf("error opening fixture file: %v", err)
	}
	defer f.Close()

	if _, err = io.Copy(w, f); err != nil {
		t.Fatalf("error writing response: %v", err)
	}
}

// mustWriteJSONResponse writes a JSON response to w.
// It uses t.Fatal to stop the test and report an error if encoding the response fails.
// This helper is useful when implementing handlers in unit tests.
func MustWriteJSONResponse(t *testing.T, w io.Writer, response any) {
	t.Helper()
	if err := json.NewEncoder(w).Encode(response); err != nil {
		t.Fatalf("Failed to write response: %v", err)
	}
}

// mustWriteErrorResponse writes an error response to w in a format that CheckResponse can parse.
// It uses t.Fatal to stop the test and report an error if encoding the response fails.
// This is useful when testing error conditions.
func MustWriteErrorResponse(t *testing.T, w io.Writer, err error) {
	t.Helper()
	MustWriteJSONResponse(t, w, map[string]any{
		"error": err.Error(),
	})
}

func ErrorOption(*retryablehttp.Request) error {
	return errors.New("RequestOptionFunc returns an error")
}

func LoadFixture(t *testing.T, filePath string) []byte {
	t.Helper()
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}

	return content
}
