package models

import (
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/database"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestClass struct {
	ID uint
}

func (c TestClass) TypeName() string {
	return "test_class"
}

func (c TestClass) GetID() uint {
	return c.ID
}

func TestUser_GrantRole(t *testing.T) {
	t.Cleanup(database.SetupDatabaseForTest())
	u := User{
		Username: "test_user_grant_role",
		Nickname: "test_user_grant_role",
		Email:    "test_user_grant_role",
		Password: "test_user_grant_role",
	}
	base.DB.Create(&u)
	r := Role{
		Name: "ttt",
	}
	base.DB.Create(&r)
	u.GrantRole(r)
	assert.Equal(t, r, u.Roles[0].Role)
	{
		dummy := "test_class"
		r = Role{
			Name:   "ttt123",
			Target: &dummy,
		}
	}
	base.DB.Create(&r)
	u.GrantRole(r, TestClass{ID: 2})
	assert.Equal(t, r, u.Roles[1].Role)
	assert.Equal(t, uint(2), u.Roles[1].TargetID)
}

func TestCan(t *testing.T) {
	t.Cleanup(database.SetupDatabaseForTest())
	base.DB.AutoMigrate(&TestClass{})
	classA := TestClass{}
	classB := TestClass{}
	base.DB.Create(&classA)
	base.DB.Create(&classB)
	dummy := "test_class"
	teacher := Role{
		Name:   "teacher",
		Target: &dummy,
	}
	assistant := Role{
		Name:   "assistant",
		Target: &dummy,
	}
	admin := Role{
		Name:   "admin",
		Target: &dummy,
	}
	globalRole := Role{
		Name: "global_role",
	}
	globalAdmin := Role{
		Name: "global_admin",
	}
	base.DB.Create(&teacher)
	base.DB.Create(&assistant)
	base.DB.Create(&admin)
	base.DB.Create(&globalRole)
	base.DB.Create(&globalAdmin)
	teacher.AddPermission("permission_teacher")
	teacher.AddPermission("permission_both")
	assistant.AddPermission("permission_both")
	admin.AddPermission("all")
	globalRole.AddPermission("global_permission")
	globalAdmin.AddPermission("all")

	testUser0 := User{
		Username: "test_user_0",
		Nickname: "tu0",
		Email:    "tu0@e.com",
		Password: "",
	}
	testUser1 := User{
		Username: "test_user_1",
		Nickname: "tu1",
		Email:    "tu1@e.com",
		Password: "",
	}
	base.DB.Create(&testUser0)
	base.DB.Create(&testUser1)
	testUser0.GrantRole(teacher, classA)
	testUser0.GrantRole(assistant, classB)
	testUser1.GrantRole(admin, classB)
	testUser0.GrantRole(globalRole)
	testUser1.GrantRole(globalAdmin)
	t.Run("scoped", func(t *testing.T) {
		t.Run("normal", func(t *testing.T) {
			assert := assert.New(t)
			assert.True(testUser0.Can("permission_teacher", classA))
			assert.False(testUser0.Can("permission_teacher", classB))
			assert.True(testUser0.Can("permission_both", classA))
			assert.True(testUser0.Can("permission_both", classB))
			assert.False(testUser0.Can("permission_both"))
			assert.True(testUser1.Can("permission_teacher", classB))
			assert.True(testUser1.Can("permission_both", classB))
		})
		t.Run("admin", func(t *testing.T) {
			assert := assert.New(t)
			assert.False(testUser1.Can("all", classA))
			assert.True(testUser1.Can("all", classB))
			assert.True(testUser1.Can("permission_teacher", classB))
			assert.True(testUser1.Can("permission_both", classB))
			assert.True(testUser1.Can("permission_non_existing", classB))
		})
	})
	t.Run("global", func(t *testing.T) {
		assert := assert.New(t)
		assert.True(testUser0.Can("global_permission"))
		assert.False(testUser0.Can("non_existing_permission"))
		assert.True(testUser1.Can("global_permission"))
		assert.True(testUser1.Can("non_existing_permission"))
	})
	assert.Panics(t, func() {
		testUser0.Can("xxx", classA, classB)
	})
}
