package middleware_test

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/middleware"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type testClass struct {
	ID uint `gorm:"primary_key" json:"id"`
}

func (c testClass) TypeName() string {
	return "test_class"
}

func (c testClass) GetID() uint {
	return c.ID
}

func setUser(user models.User) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("user", user)
			return next(c)
		}
	}
}

func responseWithUser(user models.User) response.Response {
	return response.Response{
		Message: "SUCCESS",
		Error:   nil,
		Data:    user,
	}
}

func TestHasPermission(t *testing.T) {
	t.Parallel()
	e := echo.New()
	classA := testClass{ID: 1}
	classB := testClass{ID: 2}
	dummy := "test_class"
	adminRole := models.Role{
		Name:   "admin",
		Target: &dummy,
	}
	permRole := models.Role{
		Name:   "testRole",
		Target: &dummy,
	}
	globalAdminRole := models.Role{
		Name: "globalAdmin",
	}
	globalPermRole := models.Role{
		Name: "globalTestRole",
	}
	base.DB.Create(&adminRole)
	base.DB.Create(&permRole)
	base.DB.Create(&globalAdminRole)
	base.DB.Create(&globalPermRole)
	adminRole.AddPermission("all")
	permRole.AddPermission("testPerm")
	globalAdminRole.AddPermission("all")
	globalPermRole.AddPermission("testPerm")

	testHasPermUserWithoutPerms := models.User{
		Username: "testHasPermUserWithoutPerms",
		Nickname: "uwp",
		Email:    "uwop@e.com",
		Password: "",
	}
	testHasPermUserWithClassAPerm := models.User{
		Username: "testHasPermUserWithClassAPerm",
		Nickname: "uwcap",
		Email:    "uwcap@e.com",
		Password: "",
	}
	testHasPermUserWithAllClassAPerms := models.User{
		Username: "testHasPermUserWithAllClassAPerms",
		Nickname: "uwacap",
		Email:    "uwacap@e.com",
		Password: "",
	}
	testHasPermUserWithPerm := models.User{
		Username: "testHasPermUserWithPerm",
		Nickname: "uwp",
		Email:    "uwp@e.com",
		Password: "",
	}
	testHasPermAdministrator := models.User{
		Username: "testHasPermAdministrator",
		Nickname: "a",
		Email:    "a@e.com",
		Password: "",
	}

	base.DB.Create(&testHasPermUserWithClassAPerm)
	base.DB.Create(&testHasPermUserWithAllClassAPerms)
	base.DB.Create(&testHasPermUserWithPerm)
	base.DB.Create(&testHasPermAdministrator)
	testHasPermUserWithClassAPerm.GrantRole(permRole, classA)
	testHasPermUserWithAllClassAPerms.GrantRole(adminRole, classA)
	testHasPermUserWithPerm.GrantRole(globalPermRole)
	testHasPermAdministrator.GrantRole(globalAdminRole)

	users := []models.User{
		testHasPermUserWithoutPerms,
		testHasPermUserWithClassAPerm,
		testHasPermUserWithAllClassAPerms,
		testHasPermUserWithPerm,
		testHasPermAdministrator,
	}

	permTests := []struct {
		name       string
		path       string
		permName   string
		targetType *string
		targetID   uint
	}{
		{
			name:     "perm_global",
			path:     "test_perm_global",
			permName: "testPerm",
		},
		{
			name:       "perm_a",
			path:       "test_perm",
			permName:   "testPerm",
			targetType: &dummy,
			targetID:   classA.ID,
		},
		{
			name:       "perm_b",
			path:       "test_perm",
			permName:   "testPerm",
			targetType: &dummy,
			targetID:   classB.ID,
		},
		{
			name:     "all_global",
			path:     "test_all_global",
			permName: "nonExitingPerm",
		},
		{
			name:       "all_a",
			path:       "test_all",
			permName:   "non_exiting",
			targetType: &dummy,
			targetID:   classA.ID,
		},
		{
			name:       "all_b",
			path:       "test_all",
			permName:   "nonExitingPerm",
			targetType: &dummy,
			targetID:   classB.ID,
		},
	}

	expectedRet := map[string]map[string]bool{
		"testHasPermUserWithoutPerms": {},
		"testHasPermUserWithClassAPerm": {
			"perm_a": true,
		},
		"testHasPermUserWithAllClassAPerms": {
			"perm_a": true,
			"all_a":  true,
		},
		"testHasPermUserWithPerm": {
			"perm_global": true,
		},
		"testHasPermAdministrator": {
			"perm_global": true,
			"all_global":  true,
		},
	}
	userGroups := make([]*echo.Group, len(users))
	for i, user := range users {
		userGroups[i] = e.Group("/"+user.Username, setUser(user))
		for _, permTest := range permTests {
			if permTest.targetType == nil {
				userGroups[i].POST("/"+permTest.path, testController, middleware.HasPermission(permTest.permName))
			} else {
				userGroups[i].POST("/"+permTest.path+"/:id", testController, middleware.HasPermission(permTest.permName, *permTest.targetType))
			}
		}
	}
	e.POST("/noUser/test_perm_global", testController, middleware.HasPermission("testPerm"))
	e.POST("/testHasPermAdministrator/testMultipleTarget",
		testController,
		setUser(testHasPermAdministrator),
		middleware.HasPermission("all", "targetA", "targetB"))
	e.POST("/testNonUser", testController, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("user", "nonUser")
			return next(c)
		}
	}, middleware.HasPermission("all"))

	for _, user := range users {
		t.Run(user.Username, func(t *testing.T) {
			t.Parallel()
			for _, permTest := range permTests {
				t.Run(permTest.name, func(t *testing.T) {
					t.Parallel()
					httpResp := (*http.Response)(nil)
					resp := response.Response{}
					if permTest.targetType == nil {
						httpResp = makeResp(makeReq(t, "POST", "/"+user.Username+"/"+permTest.path, nil), e)
					} else {
						httpResp = makeResp(makeReq(t, "POST", fmt.Sprintf("/%s/%s/%d", user.Username, permTest.path, permTest.targetID), nil), e)
					}
					mustJsonDecode(httpResp, &resp)
					expectedResult := expectedRet[user.Username][permTest.name]
					if expectedResult {
						assert.Equal(t, http.StatusOK, httpResp.StatusCode)
						jsonEQ(t, responseWithUser(user), resp)
					} else {
						assert.Equal(t, http.StatusForbidden, httpResp.StatusCode)
						assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resp)
					}
				})
			}
		})
	}

	t.Run("testNoUser", func(t *testing.T) {
		t.Parallel()
		resp := response.Response{}
		httpResp := makeResp(makeReq(t, "POST", "/noUser/test_perm_global", nil), e)
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusInternalServerError, httpResp.StatusCode)
		assert.Equal(t, response.MakeInternalErrorResp(), resp)
	})

	t.Run("testAdministrator", func(t *testing.T) {
		t.Parallel()
		httpResp := makeResp(makeReq(t, "POST", "/testHasPermAdministrator/testMultipleTarget", nil), e)
		resp := response.Response{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusInternalServerError, httpResp.StatusCode)
		assert.Equal(t, response.MakeInternalErrorResp(), resp)
	})

	t.Run("testAdministrator", func(t *testing.T) {
		t.Parallel()
		httpResp := makeResp(makeReq(t, "POST", "/testNonUser", nil), e)
		resp := response.Response{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusInternalServerError, httpResp.StatusCode)
		assert.Equal(t, response.MakeInternalErrorResp(), resp)
	})

	t.Run("testIllegalRouteParam", func(t *testing.T) {
		t.Parallel()
		httpResp := makeResp(makeReq(t, "POST", "/testHasPermUserWithoutPerms/test_perm/aaa", nil), e)
		resp := response.Response{}
		mustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusForbidden, httpResp.StatusCode)
		assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resp)
	})
}
