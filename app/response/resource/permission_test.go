package resource_test

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/leoleoasd/EduOJBackend/app/response/resource"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"testing"
	"time"
)

type roleWithTargetID struct {
	role models.Role
	id   uint
}

func createPermissionForTest(name string, index uint, roleId uint) (permission models.Permission) {
	permission = models.Permission{
		ID:     index,
		RoleID: roleId,
		Name:   fmt.Sprintf("test_%s_permission_%d", name, index),
	}
	return
}

func createRoleForTest(name string, index uint, permissionCount uint) (role models.Role) {
	target := fmt.Sprintf("test_%s_role_%d_target", name, index)
	permissions := make([]models.Permission, permissionCount)
	for i := range permissions {
		permissions[i] = createPermissionForTest(name, uint(i), index)
	}
	role = models.Role{
		ID:          index,
		Name:        fmt.Sprintf("test_%s_role_%d", name, index),
		Target:      &target,
		Permissions: permissions,
	}
	return
}

func createUserForTest(name string, index uint, roles ...roleWithTargetID) (user models.User) {
	user = models.User{
		ID:         index,
		Username:   fmt.Sprintf("test_%s_user_%d", name, index),
		Nickname:   fmt.Sprintf("test_%s_user_%d_nick", name, index),
		Email:      fmt.Sprintf("test_%s_user_%d@e.e", name, index),
		Password:   utils.HashPassword(fmt.Sprintf("test_%s_user_%d_pwd", name, index)),
		Roles:      make([]models.UserHasRole, len(roles)),
		RoleLoaded: true,
		CreatedAt:  time.Date(int(index), 1, 1, 1, 1, 1, 1, time.FixedZone("test_zone", 0)),
		UpdatedAt:  time.Date(int(index), 2, 2, 2, 2, 2, 2, time.FixedZone("test_zone", 0)),
		DeletedAt:  nil,
	}
	for i, role := range roles {
		user.Roles[i] = models.UserHasRole{
			ID:       uint(i),
			UserID:   user.ID,
			RoleID:   role.role.ID,
			Role:     role.role,
			TargetID: role.id,
		}
	}
	return
}

func TestGetPermission(t *testing.T) {
	permission := createPermissionForTest("get_permission", 0, 0)
	actualP := resource.GetPermission(&permission)
	expectedP := resource.Permission{
		ID:   0,
		Name: "test_get_permission_permission_0",
	}
	assert.Equal(t, expectedP, actualP)
}

func TestGetRoleAndGetRoleSlice(t *testing.T) {
	role1 := createRoleForTest("get_role", 1, 1)
	role2 := createRoleForTest("get_role", 2, 2)
	user := createUserForTest("get_role", 0,
		roleWithTargetID{role: role1, id: 1},
		roleWithTargetID{role: role2, id: 2},
	)
	t.Run("testGetRole", func(t *testing.T) {
		actualR := resource.GetRole(&user.Roles[0])
		target := "test_get_role_role_1_target"
		expectedR := resource.Role{
			ID:     0,
			Name:   "test_get_role_role_1",
			Target: &target,
			Permissions: []resource.Permission{
				{ID: 0, Name: "test_get_role_permission_0"},
			},
			TargetID: 1,
		}
		assert.Equal(t, expectedR, actualR)
	})
	t.Run("testGetRoleSlice", func(t *testing.T) {
		actualRS := resource.GetRoleSlice(user.Roles)
		target1 := "test_get_role_role_1_target"
		target2 := "test_get_role_role_2_target"
		expectedRS := []resource.Role{
			{
				ID:     0,
				Name:   "test_get_role_role_1",
				Target: &target1,
				Permissions: []resource.Permission{
					{ID: 0, Name: "test_get_role_permission_0"},
				},
				TargetID: 1,
			},
			{
				ID:     0,
				Name:   "test_get_role_role_2",
				Target: &target2,
				Permissions: []resource.Permission{
					{ID: 0, Name: "test_get_role_permission_0"},
					{ID: 1, Name: "test_get_role_permission_1"},
				},
				TargetID: 2,
			},
		}
		assert.Equal(t, expectedRS, actualRS)
	})
}
