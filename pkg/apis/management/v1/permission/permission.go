package permission

import (
	"github.com/baptistegh/go-lakekeeper/pkg/core"
)

type (
	PermissionServiceInterface interface {
		ServerPermission() ServerPermissionServiceInterface
		ProjectPermission() ProjectPermissionServiceInterface
		RolePermission() RolePermissionServiceInterface
		WarehousePermission(projectID string) WarehousePermissionServiceInterface
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

func (s *PermissionService) ProjectPermission() ProjectPermissionServiceInterface {
	return NewProjectPermissionService(s.client)
}

func (s *PermissionService) RolePermission() RolePermissionServiceInterface {
	return NewRolePermissionService(s.client)
}

func (s *PermissionService) WarehousePermission(projectID string) WarehousePermissionServiceInterface {
	return NewWarehousePermissionService(s.client, projectID)
}
