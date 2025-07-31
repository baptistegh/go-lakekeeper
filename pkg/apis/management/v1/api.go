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

package v1

import "github.com/baptistegh/go-lakekeeper/pkg/core"

const (
	APIManagementVersionPath = "/management/v1"

	ProjectIDHeader = "x-project-id"
)

type (
	ListOptions struct {
		// Next page token
		PageToken *string `url:"pageToken,omitempty"`
		// Signals an upper bound of the number of results that a client will receive.
		// Default: 100
		PageSize *int64 `url:"pageSize,omitempty"`
	}

	ListResponse struct {
		// Token to fetch the next page
		NextPageToken *string `json:"next-page-token,omitempty"`
	}
)

// WithProject add the correct header in order to select a project
// for the request. The default user project is used otherwise.
func WithProject(id string) core.RequestOptionFunc {
	return core.WithHeader(ProjectIDHeader, id)
}
