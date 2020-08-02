package middleware_test

import (
	"fmt"
	"github.com/kami-zh/go-capturer"
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
	e := echo.New()
	classA := testClass{}
	classB := testClass{}
	assert.Nil(t, base.DB.AutoMigrate(&testClass{}).Error)
	assert.True(t, base.DB.HasTable(&testClass{}))
	assert.Nil(t, base.DB.Create(&classA).Error)
	assert.Nil(t, base.DB.Create(&classB).Error)
	assert.Nil(t, base.DB.First(&classA).Error)
	assert.Nil(t, base.DB.First(&classB).Error)
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
	assert.Nil(t, base.DB.Create(&testHasPermUserWithoutPerms).Error)
	assert.Nil(t, base.DB.Create(&testHasPermUserWithClassAPerm).Error)
	assert.Nil(t, base.DB.Create(&testHasPermUserWithAllClassAPerms).Error)
	assert.Nil(t, base.DB.Create(&testHasPermUserWithPerm).Error)
	assert.Nil(t, base.DB.Create(&testHasPermAdministrator).Error)
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

	for _, user := range users {
		t.Run(user.Username, func(t *testing.T) {
			group := e.Group("/"+user.Username, setUser(user))
			for _, permTest := range permTests {
				httpResp := (*http.Response)(nil)
				resp := response.Response{}
				if permTest.targetType == nil {
					group.POST("/"+permTest.path, testController, middleware.HasPermission(permTest.permName))
					httpResp = MakeResp(MakeReq(t, "POST", "/"+user.Username+"/"+permTest.path, nil), e)
				} else {
					group.POST("/"+permTest.path+"/:id", testController, middleware.HasPermission(permTest.permName, *permTest.targetType))
					httpResp = MakeResp(MakeReq(t, "POST", fmt.Sprintf("/%s/%s/%d", user.Username, permTest.path, permTest.targetID), nil), e)
				}
				MustJsonDecode(httpResp, &resp)
				if expectedRet[user.Username][permTest.name] {
					assert.Equal(t, http.StatusOK, httpResp.StatusCode)
					JsonEQ(t, responseWithUser(user), resp)
				} else {
					assert.Equal(t, http.StatusForbidden, httpResp.StatusCode)
					assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resp)
				}
			}
		})
	}

	t.Run("testNoUser", func(t *testing.T) {
		e.Group("/noUser").POST("/test_perm_global", testController, middleware.HasPermission("testPerm"))
		resp := response.Response{}
		httpResp := (*http.Response)(nil)
		capturer.CaptureOutput(func() {
			httpResp = MakeResp(MakeReq(t, "POST", "/noUser/test_perm_global", nil), e)
		})
		MustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusInternalServerError, httpResp.StatusCode)
		assert.Equal(t, response.MakeInternalErrorResp(), resp)
	})

	t.Run("testAdministrator", func(t *testing.T) {
		adminGroup := e.Group("/testHasPermAdministrator", setUser(testHasPermAdministrator))
		adminGroup.POST("/testMultipleTarget", testController, middleware.HasPermission("all", "targetA", "targetB"))
		httpResp := (*http.Response)(nil)
		capturer.CaptureOutput(func() {
			httpResp = MakeResp(MakeReq(t, "POST", "/testHasPermAdministrator/testMultipleTarget", nil), e)
		})
		resp := response.Response{}
		MustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusInternalServerError, httpResp.StatusCode)
		assert.Equal(t, response.MakeInternalErrorResp(), resp)
	})

	t.Run("testAdministrator", func(t *testing.T) {
		e.POST("/testNonUser", testController, func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Set("user", "nonUser")
				return next(c)
			}
		}, middleware.HasPermission("all"))
		httpResp := (*http.Response)(nil)
		capturer.CaptureOutput(func() {
			httpResp = MakeResp(MakeReq(t, "POST", "/testNonUser", nil), e)
		})
		resp := response.Response{}
		MustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusInternalServerError, httpResp.StatusCode)
		assert.Equal(t, response.MakeInternalErrorResp(), resp)
	})

	assert.Nil(t, base.DB.DropTable(&testClass{}).Error)
}
