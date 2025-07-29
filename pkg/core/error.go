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
	ApiError struct {
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

func (e *ApiError) Error() string {
	if e.Response == nil {
		errMsg := "unexpected error response"
		if len(e.Message) > 0 {
			errMsg = fmt.Sprintf("%s, %s", errMsg, e.Message)
		}
		if e.Cause != nil {
			errMsg = fmt.Sprintf("%s, %v", errMsg, e.Cause)
		}
		return errMsg
	}
	return fmt.Sprintf("api error, code=%d message=%s type=%s", e.Response.Code, e.Response.Message, e.Response.Type)
}

func (e *ApiError) Type() string {
	if e.Response == nil {
		return "Unknown"
	}
	return e.Response.Type
}

func (e *ApiError) IsAuthError() bool {
	return e.StatusCode == http.StatusUnauthorized || e.StatusCode == http.StatusForbidden
}

func (e *ApiError) WithCause(err error) *ApiError {
	e.Cause = err
	return e
}

func (e *ApiError) WithMessage(format string, a ...any) *ApiError {
	e.Message = fmt.Sprintf(format, a...)
	return e
}

func ApiErrorFromResponse(response *http.Response) *ApiError {
	var apiErr ApiError

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

	// Try to unmarshal into ApiError
	if err := json.Unmarshal(bodyBytes, &apiErr); err != nil {
		apiErr.Message = string(bodyBytes) // fallback: use raw body as message
	}

	// Try to unmarshal into ApiError
	if err := json.Unmarshal(bodyBytes, &apiErr); err != nil {
		apiErr.Message = string(bodyBytes) // fallback: use raw body as message
	}

	apiErr.Status = response.Status
	apiErr.StatusCode = response.StatusCode

	return &apiErr
}

func ApiErrorFromMessage(format string, a ...any) *ApiError {
	return &ApiError{
		Message: fmt.Sprintf(format, a...),
	}
}

func ApiErrorFromError(err error) *ApiError {
	if err == nil {
		return nil
	}
	return &ApiError{
		Cause: err,
	}
}
