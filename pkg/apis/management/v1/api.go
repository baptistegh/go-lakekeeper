package v1

import "github.com/baptistegh/go-lakekeeper/pkg/core"

const (
	ApiManagementVersionPath = "/management/v1"

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
