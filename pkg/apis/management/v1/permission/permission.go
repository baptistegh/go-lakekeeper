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

package permission

import (
	"github.com/baptistegh/go-lakekeeper/pkg/core"
)

type (
	PermissionServiceInterface interface {
		ServerPermission() ServerPermissionServiceInterface
		ProjectPermission() ProjectPermissionServiceInterface
		RolePermission() RolePermissionServiceInterface
		WarehousePermission() WarehousePermissionServiceInterface
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

func (s *PermissionService) WarehousePermission() WarehousePermissionServiceInterface {
	return NewWarehousePermissionService(s.client)
}
