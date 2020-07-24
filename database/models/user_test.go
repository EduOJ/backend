package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestClass struct{
	HasRole
	typeName string
	id uint
}

func (c TestClass)  TypeName() string{
	return c.typeName
}

func (c TestClass) ID() uint{
	return c.id
}

func TestCan(t *testing.T) {
	classString := "class"

	class_A := TestClass{
		typeName: "class",
		id: 65,
	}
	class_B := TestClass{
		typeName: "class",
		id: 66,
	}
	teacherAddHomework := Permission{
		ID: 100,
		RoleID: 201,
		Name: "add_homework",
	}

	teacherCheckHomework := Permission{
		ID: 101,
		RoleID: 200,
		Name: "check_homework",
	}

	assistantCheckHomework := Permission{
		ID: 101,
		RoleID: 201,
		Name: "check_homework",
	}

	adminAll := Permission{
		ID: 102,
		RoleID: 202,
		Name: "all",
	}

	teacher :=Role{
		ID:200,
		Name: "teacher",
		Target: &classString,
		Permissions: []Permission{
			teacherAddHomework,
			teacherCheckHomework,
		},
	}

	assistant :=Role{
		ID:201,
		Name: "Assistant",
		Target: &classString,
		Permissions: []Permission{
			assistantCheckHomework,
		},
	}

	admin :=Role{
		ID:202,
		Name: "Admin",
		Target: &classString,
		Permissions: []Permission{
			adminAll,
		},
	}

	testUser0TeacherClassA := UserHasRole{
		ID: 300,
		UserID: 400,
		RoleID: 200,
		Role: teacher,
		TargetID: 65,
	}

	testUser0AssistantClassB := UserHasRole{
		ID: 301,
		UserID: 400,
		RoleID: 201,
		Role: assistant,
		TargetID: 66,
	}

	testUser1AdminClassA := UserHasRole{
		ID: 302,
		UserID: 401,
		RoleID: 202,
		Role: admin,
		TargetID: 65,
	}

	testUser1TeacherClassB := UserHasRole{
		ID: 303,
		UserID: 401,
		RoleID: 200,
		Role: teacher,
		TargetID: 66,
	}

	testUser0 := User{
		ID:       400,
		Username: "test_user_0",
		Nickname: "tu0",
		Email:    "tu0@e.com",
		Password: "",

		Roles: []UserHasRole{
			testUser0TeacherClassA,
			testUser0AssistantClassB,
		},
		RoleLoaded: true,
	}

	testUser1 := User{
		ID:       401,
		Username: "test_user_1",
		Nickname: "tu1",
		Email:    "tu1@e.com",
		Password: "",

		Roles: []UserHasRole{
			testUser1AdminClassA,
			testUser1TeacherClassB,
		},
		RoleLoaded: true,
	}

	t.Run("normalPermission", func(t *testing.T) {
		assert.Equal(t, true, testUser0.Can("add_homework", class_A))
		assert.Equal(t, false, testUser0.Can("add_homework", class_B))
		assert.Equal(t, true, testUser0.Can("check_homework", class_A))
		assert.Equal(t, true, testUser0.Can("check_homework", class_B))
		assert.Equal(t, false, testUser0.Can("check_homework"))
		assert.Equal(t, true, testUser1.Can("add_homework", class_B))
		assert.Equal(t, true, testUser1.Can("check_homework", class_B))
	})
	t.Run("allPermission", func(t *testing.T) {
		assert.Equal(t, true, testUser1.Can("all", class_A))
		assert.Equal(t, false, testUser1.Can("all", class_B))
		assert.Equal(t, true, testUser1.Can("add_homework", class_A))
		assert.Equal(t, true, testUser1.Can("check_homework", class_A))
		assert.Equal(t, true, testUser1.Can("remove_homework", class_A))
	})

}

