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
	"encoding/json"
	"testing"
)

func TestProjectAssignment_MarshalJSON(t *testing.T) {
	expected := []string{
		`{"role":"a6e5a780-258e-4bee-9bd8-f8ae3f675415","type":"project_admin"}`,
		`{"role":"9cc096bf-db1f-43f3-bea6-f0819df32db0","type":"select"}`,
		`{"type":"security_admin","user":"f5c2329c-8679-44d0-8ea3-167ee14fa94e"}`,
		`{"type":"role_creator","user":"a0d21f3d-2cbb-4066-8b77-5ec5a21680be"}`,
	}

	given := []ProjectAssignment{
		{
			Assignment: AdminProjectAssignment,
			Assignee: UserOrRole{
				Type:  RoleType,
				Value: "a6e5a780-258e-4bee-9bd8-f8ae3f675415",
			},
		},
		{
			Assignment: SelectProjectAssignment,
			Assignee: UserOrRole{
				Type:  RoleType,
				Value: "9cc096bf-db1f-43f3-bea6-f0819df32db0",
			},
		},
		{
			Assignment: SecurityAdminProjectAssignment,
			Assignee: UserOrRole{
				Type:  UserType,
				Value: "f5c2329c-8679-44d0-8ea3-167ee14fa94e",
			},
		},
		{
			Assignment: RoleCreatorProjectAssignment,
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
			t.Fatalf("exepcted %s got %s", v, string(b))
		}
	}
}

func TestProjectAssignment_UnmarshalJSON(t *testing.T) {
	expected := []ProjectAssignment{
		{
			Assignment: DataAdminProjectAssignment,
			Assignee: UserOrRole{
				Type:  RoleType,
				Value: "a6e5a780-258e-4bee-9bd8-f8ae3f675415",
			},
		},
		{
			Assignment: DescribeProjectAssignment,
			Assignee: UserOrRole{
				Type:  RoleType,
				Value: "9cc096bf-db1f-43f3-bea6-f0819df32db0",
			},
		},
		{
			Assignment: CreateProjectAssignment,
			Assignee: UserOrRole{
				Type:  UserType,
				Value: "f5c2329c-8679-44d0-8ea3-167ee14fa94e",
			},
		},
		{
			Assignment: ModifyProjectAssignment,
			Assignee: UserOrRole{
				Type:  UserType,
				Value: "a0d21f3d-2cbb-4066-8b77-5ec5a21680be",
			},
		},
	}

	given := []string{
		`{"role":"a6e5a780-258e-4bee-9bd8-f8ae3f675415","type":"data_admin"}`,
		`{"role":"9cc096bf-db1f-43f3-bea6-f0819df32db0","type":"describe"}`,
		`{"type":"create","user":"f5c2329c-8679-44d0-8ea3-167ee14fa94e"}`,
		`{"type":"modify","user":"a0d21f3d-2cbb-4066-8b77-5ec5a21680be"}`,
	}

	for k, v := range expected {
		var aux ProjectAssignment
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
