package core

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestError_From(t *testing.T) {
	t.Parallel()

	t.Run("FromError", func(t *testing.T) {
		t.Parallel()
		given := APIErrorFromError(errors.New("error message, testing"))

		expected := "unexpected error response, error message, testing"

		require.ErrorContains(t, given, expected)
	})

	t.Run("FromErro Nil", func(t *testing.T) {
		t.Parallel()
		given := APIErrorFromError(nil)

		assert.Nil(t, given)
	})

	t.Run("From Message", func(t *testing.T) {
		t.Parallel()
		given := APIErrorFromMessage("error message %s", "testing")

		expected := "unexpected error response, error message testing"

		require.ErrorContains(t, given, expected)
	})

	t.Run("FromResponse", func(t *testing.T) {
		t.Parallel()
		mux := http.NewServeMux()

		server := httptest.NewServer(mux)
		t.Cleanup(server.Close)

		mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			if r.Method == http.MethodGet {
				if _, err := fmt.Fprintf(w, `{"error":{"code": 32,"message": "testing message","stack": ["stack1"],"type": "error-type"}}`); err != nil {
					t.Fatalf("error writing http response, %v", err)
				}
			}
		})

		resp, err := http.Get(server.URL + "/error")
		require.NoError(t, err, "failed to read service account token")

		given := APIErrorFromResponse(resp)

		assert.Equal(t, 32, given.Response.Code)
		assert.Equal(t, "testing message", given.Response.Message)
		assert.Equal(t, "error-type", given.Type())

		require.ErrorContains(t, given, "api error, code=32 message=testing message type=error-type")
	})
}

func TestError_With(t *testing.T) {
	t.Parallel()

	t.Run("WithCause", func(t *testing.T) {
		t.Parallel()
		apiErr := &APIError{}

		given := apiErr.WithCause(errors.New("testing error"))

		require.ErrorContains(t, given, "unexpected error response, testing error")
	})

	t.Run("WithMessage", func(t *testing.T) {
		t.Parallel()
		apiErr := &APIError{}

		given := apiErr.WithMessage("message is %s", "testing")

		require.ErrorContains(t, given, "unexpected error response, message is testing")
	})
}

func TestError_IsAuthError(t *testing.T) {
	t.Parallel()
	unauthorized := &APIError{
		StatusCode: 401,
	}

	forbidden := &APIError{
		StatusCode: 403,
	}

	assert.True(t, unauthorized.IsAuthError())
	assert.True(t, forbidden.IsAuthError())
}

func TestError_Type(t *testing.T) {
	t.Parallel()

	t.Run("Unknown", func(t *testing.T) {
		t.Parallel()
		given := &APIError{}

		assert.Equal(t, "Unknown", given.Type())
	})
}
