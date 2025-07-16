package permission

import (
	"github.com/baptistegh/go-lakekeeper/pkg/core"
)

type (
	PermissionServiceInterface interface {
		ServerPermission() ServerPermissionServiceInterface
	}

	// PermissionService handles communication with permission endpoints of the Lakekeeper API.
	//
	// Lakekeeper API docs:
	// https://docs.lakekeeper.io/docs/nightly/api/management/#tag/permissions
	PermissionService struct {
		client core.Client
	}
)

func NewPermissionService(client core.Client) PermissionServiceInterface {
	return &PermissionService{
		client: client,
	}
}

func (s *PermissionService) ServerPermission() ServerPermissionServiceInterface {
	return NewServerPermissionService(s.client)
}
