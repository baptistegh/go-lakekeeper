package v1

import "github.com/baptistegh/go-lakekeeper/pkg/core"

const (
	ApiManagementVersionPath = "/management/v1"

	ProjectIDHeader = "x-project-id"
)

// WithProject add the correct header in order to select a project
// for the request. The default user project is used otherwise.
func WithProject(id string) core.RequestOptionFunc {
	return core.WithHeader(ProjectIDHeader, id)
}
