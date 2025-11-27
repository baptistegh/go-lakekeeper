package permission

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWarehouseAssignment_MarshalJSON(t *testing.T) {
	expected := []string{
		`{"role":"a6e5a780-258e-4bee-9bd8-f8ae3f675415","type":"ownership"}`,
		`{"role":"9cc096bf-db1f-43f3-bea6-f0819df32db0","type":"pass_grants"}`,
		`{"role":"9cc096bf-db1f-43f3-bea6-f0819df32db0","type":"manage_grants"}`,
		`{"role":"9cc096bf-db1f-43f3-bea6-f0819df32db0","type":"describe"}`,
		`{"type":"select","user":"9cc096bf-db1f-43f3-bea6-f0819df32db0"}`,
		`{"type":"create","user":"9cc096bf-db1f-43f3-bea6-f0819df32db0"}`,
		`{"type":"modify","user":"9cc096bf-db1f-43f3-bea6-f0819df32db0"}`,
	}

	given := []WarehouseAssignment{
		{
			Assignment: OwnershipWarehouseAssignment,
			Assignee: UserOrRole{
				Type:  RoleType,
				Value: "a6e5a780-258e-4bee-9bd8-f8ae3f675415",
			},
		},
		{
			Assignment: PassGrantsAdminWarehouseAssignment,
			Assignee: UserOrRole{
				Type:  RoleType,
				Value: "9cc096bf-db1f-43f3-bea6-f0819df32db0",
			},
		},
		{
			Assignment: ManageGrantsAdminWarehouseAssignment,
			Assignee: UserOrRole{
				Type:  RoleType,
				Value: "9cc096bf-db1f-43f3-bea6-f0819df32db0",
			},
		},
		{
			Assignment: DescribeWarehouseAssignment,
			Assignee: UserOrRole{
				Type:  RoleType,
				Value: "9cc096bf-db1f-43f3-bea6-f0819df32db0",
			},
		},
		{
			Assignment: SelectWarehouseAssignment,
			Assignee: UserOrRole{
				Type:  UserType,
				Value: "9cc096bf-db1f-43f3-bea6-f0819df32db0",
			},
		},
		{
			Assignment: CreateWarehouseAssignment,
			Assignee: UserOrRole{
				Type:  UserType,
				Value: "9cc096bf-db1f-43f3-bea6-f0819df32db0",
			},
		},
		{
			Assignment: ModifyWarehouseAssignment,
			Assignee: UserOrRole{
				Type:  UserType,
				Value: "9cc096bf-db1f-43f3-bea6-f0819df32db0",
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

func TestWarehouseAssignment_Getters(t *testing.T) {
	t.Run("user assignee", func(t *testing.T) {
		wa := WarehouseAssignment{
			Assignee: UserOrRole{
				Type:  UserType,
				Value: "user-id-123",
			},
			Assignment: OwnershipWarehouseAssignment,
		}

		assert.Equal(t, "ownership", wa.GetAssignment())
		assert.Equal(t, "user-id-123", wa.GetPrincipalID())
		assert.Equal(t, UserType, wa.GetPrincipalType())
	})

	t.Run("role assignee", func(t *testing.T) {
		wa := WarehouseAssignment{
			Assignee: UserOrRole{
				Type:  RoleType,
				Value: "role-id-456",
			},
			Assignment: SelectWarehouseAssignment,
		}

		assert.Equal(t, "select", wa.GetAssignment())
		assert.Equal(t, "role-id-456", wa.GetPrincipalID())
		assert.Equal(t, RoleType, wa.GetPrincipalType())
	})
}

func TestWarehouseAssignment_UnmarshalJSON(t *testing.T) {
	expected := []WarehouseAssignment{
		{
			Assignment: OwnershipWarehouseAssignment,
			Assignee: UserOrRole{
				Type:  RoleType,
				Value: "a6e5a780-258e-4bee-9bd8-f8ae3f675415",
			},
		},
		{
			Assignment: DescribeWarehouseAssignment,
			Assignee: UserOrRole{
				Type:  RoleType,
				Value: "9cc096bf-db1f-43f3-bea6-f0819df32db0",
			},
		},
		{
			Assignment: CreateWarehouseAssignment,
			Assignee: UserOrRole{
				Type:  UserType,
				Value: "f5c2329c-8679-44d0-8ea3-167ee14fa94e",
			},
		},
		{
			Assignment: ModifyWarehouseAssignment,
			Assignee: UserOrRole{
				Type:  UserType,
				Value: "a0d21f3d-2cbb-4066-8b77-5ec5a21680be",
			},
		},
		{
			Assignment: SelectWarehouseAssignment,
			Assignee: UserOrRole{
				Type:  UserType,
				Value: "a0d21f3d-2cbb-4066-8b77-5ec5a21680be",
			},
		},
		{
			Assignment: ManageGrantsAdminWarehouseAssignment,
			Assignee: UserOrRole{
				Type:  UserType,
				Value: "a0d21f3d-2cbb-4066-8b77-5ec5a21680be",
			},
		},
		{
			Assignment: PassGrantsAdminWarehouseAssignment,
			Assignee: UserOrRole{
				Type:  UserType,
				Value: "a0d21f3d-2cbb-4066-8b77-5ec5a21680be",
			},
		},
	}

	given := []string{
		`{"role":"a6e5a780-258e-4bee-9bd8-f8ae3f675415","type":"ownership"}`,
		`{"role":"9cc096bf-db1f-43f3-bea6-f0819df32db0","type":"describe"}`,
		`{"type":"create","user":"f5c2329c-8679-44d0-8ea3-167ee14fa94e"}`,
		`{"type":"modify","user":"a0d21f3d-2cbb-4066-8b77-5ec5a21680be"}`,
		`{"type":"select","user":"a0d21f3d-2cbb-4066-8b77-5ec5a21680be"}`,
		`{"type":"manage_grants","user":"a0d21f3d-2cbb-4066-8b77-5ec5a21680be"}`,
		`{"type":"pass_grants","user":"a0d21f3d-2cbb-4066-8b77-5ec5a21680be"}`,
	}

	for k, v := range expected {
		var aux WarehouseAssignment
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
