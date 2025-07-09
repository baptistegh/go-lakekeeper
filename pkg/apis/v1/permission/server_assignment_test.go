package permission

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/baptistegh/go-lakekeeper/pkg/core"
	"github.com/baptistegh/go-lakekeeper/pkg/testutil" // Import testutil
	"github.com/hashicorp/go-retryablehttp"
)

func TestServerPermissions_MarshalJSON(t *testing.T) {
	expected := []string{
		`{"role":"a6e5a780-258e-4bee-9bd8-f8ae3f675415","type":"admin"}`,
		`{"role":"9cc096bf-db1f-43f3-bea6-f0819df32db0","type":"operator"}`,
		`{"type":"admin","user":"f5c2329c-8679-44d0-8ea3-167ee14fa94e"}`,
		`{"type":"operator","user":"a0d21f3d-2cbb-4066-8b77-5ec5a21680be"}`,
	}

	given := []ServerAssignment{
		{
			Assignment: AdminServerAssignment,
			Assignee: UserOrRole{
				Type:  RoleType,
				Value: "a6e5a780-258e-4bee-9bd8-f8ae3f675415",
			},
		},
		{
			Assignment: OperatorServerAssignment,
			Assignee: UserOrRole{
				Type:  RoleType,
				Value: "9cc096bf-db1f-43f3-bea6-f0819df32db0",
			},
		},
		{
			Assignment: AdminServerAssignment,
			Assignee: UserOrRole{
				Type:  UserType,
				Value: "f5c2329c-8679-44d0-8ea3-167ee14fa94e",
			},
		},
		{
			Assignment: OperatorServerAssignment,
			Assignee: UserOrRole{
				Type:  UserType,
				Value: "a0d21f3d-2cbb-4066-8b77-5ec5a21680be",
			},
		},
	}

	for k, v := range expected {
		b, err := json.Marshal(given[k])
		if err != nil {
			t.Fatalf("%v", err)
		}
		if string(b) != v {
			t.Fatalf("exepected %s got %s", v, string(b))
		}
	}
}

func TestServerPermissionService_GetServerAccess(t *testing.T) {
	mockResponse := &GetServerAccessResponse{
		AllowedActions: []ServerAction{CreateProjectServerAction, UpdateUsersServerAction},
	}
	mockResponseBody, _ := json.Marshal(mockResponse)

	mockClient := &testutil.MockClient{
		DoFunc: func(req *retryablehttp.Request, v any) (*http.Response, *core.ApiError) {
			if req.Method != http.MethodGet {
				t.Errorf("Expected GET method, got %s", req.Method)
			}
			if req.URL.Path != "/permissions/server/access" { // Path is relative in NewRequest
				t.Errorf("Expected path /permissions/server/access, got %s", req.URL.Path)
			}

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBuffer(mockResponseBody)),
				Header:     make(http.Header),
			}
			if v != nil {
				_ = json.Unmarshal(mockResponseBody, v)
			}
			return resp, nil
		},
		NewRequestFunc: func(method, path string, opt any, options []core.RequestOptionFunc) (*retryablehttp.Request, error) {
			if method != http.MethodGet {
				return nil, fmt.Errorf("expected method %s, got %s", http.MethodGet, method)
			}
			expectedPath := "/permissions/server/access"
			if path != expectedPath {
				return nil, fmt.Errorf("expected path %s, got %s", expectedPath, path)
			}
			// No specific checks for opt or options for now, as they are handled by the DoFunc or not critical for this test.
			req, err := retryablehttp.NewRequest(method, "http://example.com"+path, opt)
			if err != nil {
				return nil, err
			}
			return req, nil
		},
	}

	service := NewServerPermissionService(mockClient)
	ctx := context.Background()
	principal := &UserOrRole{Type: UserType, Value: "test-user"}

	resp, httpResp, apiErr := service.GetServerAccess(ctx, principal)

	if apiErr != nil {
		t.Fatalf("Expected no API error, got %v", apiErr)
	}
	if httpResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %d", httpResp.StatusCode)
	}
	if resp == nil {
		t.Fatal("Expected a response, got nil")
	}
	if len(resp.AllowedActions) != 2 {
		t.Fatalf("Expected 2 allowed actions, got %d", len(resp.AllowedActions))
	}
	if resp.AllowedActions[0] != CreateProjectServerAction || resp.AllowedActions[1] != UpdateUsersServerAction {
		t.Fatalf("Allowed actions mismatch")
	}
}

func TestServerPermissionService_GetServerAssignments(t *testing.T) {
	mockResponse := &GetServerAssignmentsResponse{
		Assignments: []ServerAssignment{
			{Assignee: UserOrRole{Type: UserType, Value: "user1"}, Assignment: AdminServerAssignment},
			{Assignee: UserOrRole{Type: RoleType, Value: "role1"}, Assignment: OperatorServerAssignment},
		},
	}
	mockResponseBody, _ := json.Marshal(mockResponse)

	mockClient := &testutil.MockClient{
		DoFunc: func(req *retryablehttp.Request, v any) (*http.Response, *core.ApiError) {
			if req.Method != http.MethodGet {
				t.Errorf("Expected GET method, got %s", req.Method)
			}
			if req.URL.Path != "/permissions/server/assignments" { // Path is relative in NewRequest
				t.Errorf("Expected path /permissions/server/assignments, got %s", req.URL.Path)
			}
			if req.URL.RawQuery != "relations=admin&relations=operator" {
				t.Errorf("Expected query relations=admin&relations=operator, got %s", req.URL.RawQuery)
			}

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBuffer(mockResponseBody)),
				Header:     make(http.Header),
			}
			if v != nil {
				_ = json.Unmarshal(mockResponseBody, v)
			}
			return resp, nil
		},
		NewRequestFunc: func(method, path string, opt any, options []core.RequestOptionFunc) (*retryablehttp.Request, error) {
			if method != http.MethodGet {
				return nil, fmt.Errorf("expected method %s, got %s", http.MethodGet, method)
			}
			expectedPath := "/permissions/server/assignments"
			if path != expectedPath {
				return nil, fmt.Errorf("expected path %s, got %s", expectedPath, path)
			}
			// No specific checks for opt or options for now, as they are handled by the DoFunc or not critical for this test.
			req, err := retryablehttp.NewRequest(method, "http://example.com"+path, opt)
			if err != nil {
				return nil, err
			}
			return req, nil
		},
	}

	service := NewServerPermissionService(mockClient)
	ctx := context.Background()
	relations := []ServerRelation{AdminServerRelation, OperatorServerRelation}

	resp, httpResp, apiErr := service.GetServerAssignments(ctx, relations)

	if apiErr != nil {
		t.Fatalf("Expected no API error, got %v", apiErr)
	}
	if httpResp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %d", httpResp.StatusCode)
	}
	if resp == nil {
		t.Fatal("Expected a response, got nil")
	}
	if len(resp.Assignments) != 2 {
		t.Fatalf("Expected 2 assignments, got %d", len(resp.Assignments))
	}
	if resp.Assignments[0].Assignee.Value != "user1" || resp.Assignments[1].Assignee.Value != "role1" {
		t.Fatalf("Assignments mismatch")
	}
}

func TestServerPermissionService_UpdateServerAssignments(t *testing.T) {
	mockClient := &testutil.MockClient{
		DoFunc: func(req *retryablehttp.Request, v any) (*http.Response, *core.ApiError) {
			if req.Method != http.MethodPost {
				t.Errorf("Expected POST method, got %s", req.Method)
			}
			if req.URL.Path != "/permissions/server/assignments" { // Path is relative in NewRequest
				t.Errorf("Expected path /permissions/server/assignments, got %s", req.URL.Path)
			}

			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("Failed to read request body: %v", err)
			}
			var requestBody UpdateServerAssignmentsRequest
			_ = json.Unmarshal(bodyBytes, &requestBody)

			if len(requestBody.Writes) != 1 || requestBody.Writes[0].Assignee.Value != "new-user" {
				t.Errorf("Request body writes mismatch: %+v", requestBody.Writes)
			}
			if len(requestBody.Deletes) != 1 || requestBody.Deletes[0].Assignee.Value != "old-user" {
				t.Errorf("Request body deletes mismatch: %+v", requestBody.Deletes)
			}

			resp := &http.Response{
				StatusCode: http.StatusNoContent,
				Body:       io.NopCloser(bytes.NewBufferString("")),
				Header:     make(http.Header),
			}
			return resp, nil
		},
		NewRequestFunc: func(method, path string, opt any, options []core.RequestOptionFunc) (*retryablehttp.Request, error) {
			if method != http.MethodPost {
				return nil, fmt.Errorf("expected method %s, got %s", http.MethodPost, method)
			}
			expectedPath := "/permissions/server/assignments"
			if path != expectedPath {
				return nil, fmt.Errorf("expected path %s, got %s", expectedPath, path)
			}
			// No specific checks for opt or options for now, as they are handled by the DoFunc or not critical for this test.
			req, err := retryablehttp.NewRequest(method, "http://example.com"+path, opt)
			if err != nil {
				return nil, err
			}
			return req, nil
		},
	}

	service := NewServerPermissionService(mockClient)
	ctx := context.Background()
	body := &UpdateServerAssignmentsRequest{
		Writes: []ServerAssignment{
			{Assignee: UserOrRole{Type: UserType, Value: "new-user"}, Assignment: AdminServerAssignment},
		},
		Deletes: []ServerAssignment{
			{Assignee: UserOrRole{Type: UserType, Value: "old-user"}, Assignment: OperatorServerAssignment},
		},
	}

	httpResp, apiErr := service.UpdateServerAssignments(ctx, body)

	if apiErr != nil {
		t.Fatalf("Expected no API error, got %v", apiErr)
	}
	if httpResp.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status No Content, got %d", httpResp.StatusCode)
	}
}

func TestServerPermissions_UnmarshalJSON(t *testing.T) {
	expected := []ServerAssignment{
		{
			Assignment: AdminServerAssignment,
			Assignee: UserOrRole{
				Type:  RoleType,
				Value: "a6e5a780-258e-4bee-9bd8-f8ae3f675415",
			},
		},
		{
			Assignment: OperatorServerAssignment,
			Assignee: UserOrRole{
				Type:  RoleType,
				Value: "9cc096bf-db1f-43f3-bea6-f0819df32db0",
			},
		},
		{
			Assignment: AdminServerAssignment,
			Assignee: UserOrRole{
				Type:  UserType,
				Value: "f5c2329c-8679-44d0-8ea3-167ee14fa94e",
			},
		},
		{
			Assignment: OperatorServerAssignment,
			Assignee: UserOrRole{
				Type:  UserType,
				Value: "a0d21f3d-2cbb-4066-8b77-5ec5a21680be",
			},
		},
	}

	given := []string{
		`{"role":"a6e5a780-258e-4bee-9bd8-f8ae3f675415","type":"admin"}`,
		`{"role":"9cc096bf-db1f-43f3-bea6-f0819df32db0","type":"operator"}`,
		`{"type":"admin","user":"f5c2329c-8679-44d0-8ea3-167ee14fa94e"}`,
		`{"type":"operator","user":"a0d21f3d-2cbb-4066-8b77-5ec5a21680be"}`,
	}

	for k, v := range expected {
		var aux ServerAssignment
		err := json.Unmarshal([]byte(given[k]), &aux)
		if err != nil {
			t.Fatalf("%v", err)
		}

		if v.Assignment != aux.Assignment {
			t.Fatalf("expected %s got %s", v.Assignment, aux.Assignment)
		}

		if v.Assignee.Type != aux.Assignee.Type {
			t.Fatalf("expected %s got %s", v.Assignee.Type, aux.Assignee.Type)
		}

		if v.Assignee.Value != aux.Assignee.Value {
			t.Fatalf("expected %s got %s", v.Assignee.Type, aux.Assignee.Value)
		}
	}
}
