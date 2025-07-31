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

package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type (
	APIError struct {
		Status     string         `json:"-"`
		StatusCode int            `json:"-"`
		Message    string         `json:"-"`
		Response   *ErrorResponse `json:"error"`
		Cause      error          `json:"-"`
	}

	ErrorResponse struct {
		Code    int      `json:"code"`
		Message string   `json:"message"`
		Stack   []string `json:"stack"`
		Type    string   `json:"type"`
	}
)

func (e *APIError) Error() string {
	if e.Response == nil {
		errMsg := "unexpected error response"
		if e.Message != "" {
			errMsg = fmt.Sprintf("%s, %s", errMsg, e.Message)
		}
		if e.Cause != nil {
			errMsg = fmt.Sprintf("%s, %v", errMsg, e.Cause)
		}
		return errMsg
	}
	return fmt.Sprintf("api error, code=%d message=%s type=%s", e.Response.Code, e.Response.Message, e.Response.Type)
}

func (e *APIError) Type() string {
	if e.Response == nil {
		return "Unknown"
	}
	return e.Response.Type
}

func (e *APIError) IsAuthError() bool {
	return e.StatusCode == http.StatusUnauthorized || e.StatusCode == http.StatusForbidden
}

func (e *APIError) WithCause(err error) *APIError {
	e.Cause = err
	return e
}

func (e *APIError) WithMessage(format string, a ...any) *APIError {
	e.Message = fmt.Sprintf(format, a...)
	return e
}

func APIErrorFromResponse(response *http.Response) *APIError {
	var apiErr APIError

	// Read the body once
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		apiErr.Message = "failed to read response body"
		apiErr.Status = response.Status
		apiErr.StatusCode = response.StatusCode
		return &apiErr
	}

	// Restore the body for potential further use
	response.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Try to unmarshal into APIError
	if err := json.Unmarshal(bodyBytes, &apiErr); err != nil {
		apiErr.Message = string(bodyBytes) // fallback: use raw body as message
	}

	// Try to unmarshal into APIError
	if err := json.Unmarshal(bodyBytes, &apiErr); err != nil {
		apiErr.Message = string(bodyBytes) // fallback: use raw body as message
	}

	apiErr.Status = response.Status
	apiErr.StatusCode = response.StatusCode

	return &apiErr
}

func APIErrorFromMessage(format string, a ...any) *APIError {
	return &APIError{
		Message: fmt.Sprintf(format, a...),
	}
}

func APIErrorFromError(err error) *APIError {
	if err == nil {
		return nil
	}
	return &APIError{
		Cause: err,
	}
}
