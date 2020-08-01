package middleware_test

import (
	"github.com/kami-zh/go-capturer"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/middleware"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
)

type TestClass struct {
	ID uint `gorm:"primary_key" json:"id"`
}

func (c TestClass) TypeName() string {
	return "test_class"
}

func (c TestClass) GetID() uint {
	return c.ID
}

func AddUser(user models.User) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("user", user)
			return next(c)
		}
	}
}

func ResponseWithUser(user models.User) response.Response {
	return response.Response{
		Message: "SUCCESS",
		Error:   nil,
		Data:    user,
	}
}

func TestPermission(t *testing.T) {
	oldEcho := base.Echo
	base.Echo = echo.New()
	t.Cleanup(func() {
		base.Echo = oldEcho
	})
	classA := TestClass{}
	classB := TestClass{}
	assert.Nil(t, base.DB.AutoMigrate(&TestClass{}).Error)
	assert.True(t, base.DB.HasTable(&TestClass{}))
	assert.Nil(t, base.DB.Create(&classA).Error)
	assert.Nil(t, base.DB.Create(&classB).Error)
	assert.Nil(t, base.DB.First(&classA).Error)
	assert.Nil(t, base.DB.First(&classB).Error)
	classAID := strconv.Itoa(int(classA.ID))
	classBID := strconv.Itoa(int(classB.ID))
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

	userWithoutPerms := models.User{
		Username: "userWithoutPerms",
		Nickname: "uwp",
		Email:    "uwop@e.com",
		Password: "",
	}
	userWithClassAPerm := models.User{
		Username: "userWithClassAPerm",
		Nickname: "uwcap",
		Email:    "uwcap@e.com",
		Password: "",
	}
	userWithAllClassAPerms := models.User{
		Username: "userWithAllClassAPerms",
		Nickname: "uwacap",
		Email:    "uwacap@e.com",
		Password: "",
	}
	userWithPerm := models.User{
		Username: "userWithPerm",
		Nickname: "uwp",
		Email:    "uwp@e.com",
		Password: "",
	}
	administrator := models.User{
		Username: "administrator",
		Nickname: "a",
		Email:    "a@e.com",
		Password: "",
	}
	assert.Nil(t, base.DB.Create(&userWithoutPerms).Error)
	assert.Nil(t, base.DB.Create(&userWithClassAPerm).Error)
	assert.Nil(t, base.DB.Create(&userWithAllClassAPerms).Error)
	assert.Nil(t, base.DB.Create(&userWithPerm).Error)
	assert.Nil(t, base.DB.Create(&administrator).Error)
	userWithClassAPerm.GrantRole(permRole, classA)
	userWithAllClassAPerms.GrantRole(adminRole, classA)
	userWithPerm.GrantRole(globalPermRole)
	administrator.GrantRole(globalAdminRole)
	groups := []*echo.Group{
		base.Echo.Group("/noUser"),
		base.Echo.Group("/userWithoutPerms", AddUser(userWithoutPerms)),
		base.Echo.Group("/userWithClassAPerm", AddUser(userWithClassAPerm)),
		base.Echo.Group("/userWithAllClassAPerms", AddUser(userWithAllClassAPerms)),
		base.Echo.Group("/userWithPerm", AddUser(userWithPerm)),
		base.Echo.Group("/administrator", AddUser(administrator)),
	}
	for _, group := range groups {
		group.POST("/test_perm_global", testController, middleware.Permission("testPerm"))
		group.POST("/test_perm/:id", testController, middleware.Permission("testPerm", "test_class"))
		group.POST("/test_all_global", testController, middleware.Permission("all"))
		group.POST("/test_all/:id", testController, middleware.Permission("all", "test_class"))
	}
	t.Run("testNoUser", func(t *testing.T) {
		resp := response.Response{}
		httpResp := (*http.Response)(nil)
		capturer.CaptureOutput(func() {
			httpResp = MakeResp(MakeReq(t, "POST", "/noUser/test_perm_global", nil))
		})
		MustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusInternalServerError, httpResp.StatusCode)
		assert.Equal(t, response.MakeInternalErrorResp(), resp)
	})
	t.Run("testUserWithoutPerms", func(t *testing.T) {
		resps := make([]response.Response, 4)
		httpResps := []*http.Response{
			MakeResp(MakeReq(t, "POST", "/userWithoutPerms/test_perm_global", nil)),
			MakeResp(MakeReq(t, "POST", "/userWithoutPerms/test_perm/"+classAID, nil)),
			MakeResp(MakeReq(t, "POST", "/userWithoutPerms/test_all_global", nil)),
			MakeResp(MakeReq(t, "POST", "/userWithoutPerms/test_all/"+classAID, nil)),
		}
		for index, httpResp := range httpResps {
			MustJsonDecode(httpResp, &resps[index])
		}
		assert.Equal(t, http.StatusForbidden, httpResps[0].StatusCode)
		assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resps[0])
		assert.Equal(t, http.StatusForbidden, httpResps[1].StatusCode)
		assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resps[1])
		assert.Equal(t, http.StatusForbidden, httpResps[2].StatusCode)
		assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resps[2])
		assert.Equal(t, http.StatusForbidden, httpResps[3].StatusCode)
		assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resps[3])
	})
	t.Run("testUserWithClassAPerm", func(t *testing.T) {
		resps := make([]response.Response, 6)
		httpResps := []*http.Response{
			MakeResp(MakeReq(t, "POST", "/userWithClassAPerm/test_perm_global", nil)),
			MakeResp(MakeReq(t, "POST", "/userWithClassAPerm/test_perm/"+classAID, nil)),
			MakeResp(MakeReq(t, "POST", "/userWithClassAPerm/test_perm/"+classBID, nil)),
			MakeResp(MakeReq(t, "POST", "/userWithClassAPerm/test_all_global", nil)),
			MakeResp(MakeReq(t, "POST", "/userWithClassAPerm/test_all/"+classAID, nil)),
			MakeResp(MakeReq(t, "POST", "/userWithClassAPerm/test_all/"+classBID, nil)),
		}
		for index, httpResp := range httpResps {
			MustJsonDecode(httpResp, &resps[index])
		}
		assert.Equal(t, http.StatusForbidden, httpResps[0].StatusCode)
		assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resps[0])
		assert.Equal(t, http.StatusOK, httpResps[1].StatusCode)
		JsonEQ(t, ResponseWithUser(userWithClassAPerm), resps[1])
		assert.Equal(t, http.StatusForbidden, httpResps[2].StatusCode)
		assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resps[2])
		assert.Equal(t, http.StatusForbidden, httpResps[3].StatusCode)
		assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resps[3])
		assert.Equal(t, http.StatusForbidden, httpResps[4].StatusCode)
		assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resps[4])
		assert.Equal(t, http.StatusForbidden, httpResps[5].StatusCode)
		assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resps[5])
	})
	t.Run("testUserWithAllClassAPerms", func(t *testing.T) {
		resps := make([]response.Response, 6)
		httpResps := []*http.Response{
			MakeResp(MakeReq(t, "POST", "/userWithAllClassAPerms/test_perm_global", nil)),
			MakeResp(MakeReq(t, "POST", "/userWithAllClassAPerms/test_perm/"+classAID, nil)),
			MakeResp(MakeReq(t, "POST", "/userWithAllClassAPerms/test_perm/"+classBID, nil)),
			MakeResp(MakeReq(t, "POST", "/userWithAllClassAPerms/test_all_global", nil)),
			MakeResp(MakeReq(t, "POST", "/userWithAllClassAPerms/test_all/"+classAID, nil)),
			MakeResp(MakeReq(t, "POST", "/userWithAllClassAPerms/test_all/"+classBID, nil)),
		}
		for index, httpResp := range httpResps {
			MustJsonDecode(httpResp, &resps[index])
		}
		assert.Equal(t, http.StatusForbidden, httpResps[0].StatusCode)
		assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resps[0])
		assert.Equal(t, http.StatusOK, httpResps[1].StatusCode)
		JsonEQ(t, ResponseWithUser(userWithAllClassAPerms), resps[1])
		assert.Equal(t, http.StatusForbidden, httpResps[2].StatusCode)
		assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resps[2])
		assert.Equal(t, http.StatusForbidden, httpResps[3].StatusCode)
		assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resps[3])
		assert.Equal(t, http.StatusOK, httpResps[4].StatusCode)
		JsonEQ(t, ResponseWithUser(userWithAllClassAPerms), resps[4])
		assert.Equal(t, http.StatusForbidden, httpResps[5].StatusCode)
		assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resps[5])
	})
	t.Run("testUserWithPerm", func(t *testing.T) {
		resps := make([]response.Response, 6)
		httpResps := []*http.Response{
			MakeResp(MakeReq(t, "POST", "/userWithPerm/test_perm_global", nil)),
			MakeResp(MakeReq(t, "POST", "/userWithPerm/test_perm/"+classAID, nil)),
			MakeResp(MakeReq(t, "POST", "/userWithPerm/test_perm/"+classBID, nil)),
			MakeResp(MakeReq(t, "POST", "/userWithPerm/test_all_global", nil)),
			MakeResp(MakeReq(t, "POST", "/userWithPerm/test_all/"+classAID, nil)),
			MakeResp(MakeReq(t, "POST", "/userWithPerm/test_all/"+classBID, nil)),
		}
		for index, httpResp := range httpResps {
			MustJsonDecode(httpResp, &resps[index])
		}
		assert.Equal(t, http.StatusOK, httpResps[0].StatusCode)
		JsonEQ(t, ResponseWithUser(userWithPerm), resps[0])
		assert.Equal(t, http.StatusForbidden, httpResps[1].StatusCode)
		assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resps[1])
		assert.Equal(t, http.StatusForbidden, httpResps[2].StatusCode)
		assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resps[2])
		assert.Equal(t, http.StatusForbidden, httpResps[3].StatusCode)
		assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resps[3])
		assert.Equal(t, http.StatusForbidden, httpResps[4].StatusCode)
		assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resps[4])
		assert.Equal(t, http.StatusForbidden, httpResps[5].StatusCode)
		assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resps[5])
	})
	t.Run("testAdministrator", func(t *testing.T) {
		resps := make([]response.Response, 6)
		httpResps := []*http.Response{
			MakeResp(MakeReq(t, "POST", "/administrator/test_perm_global", nil)),
			MakeResp(MakeReq(t, "POST", "/administrator/test_perm/"+classAID, nil)),
			MakeResp(MakeReq(t, "POST", "/administrator/test_perm/"+classBID, nil)),
			MakeResp(MakeReq(t, "POST", "/administrator/test_all_global", nil)),
			MakeResp(MakeReq(t, "POST", "/administrator/test_all/"+classAID, nil)),
			MakeResp(MakeReq(t, "POST", "/administrator/test_all/"+classBID, nil)),
		}
		for index, httpResp := range httpResps {
			MustJsonDecode(httpResp, &resps[index])
		}
		assert.Equal(t, http.StatusOK, httpResps[0].StatusCode)
		JsonEQ(t, ResponseWithUser(administrator), resps[0])
		assert.Equal(t, http.StatusForbidden, httpResps[1].StatusCode)
		assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resps[1])
		assert.Equal(t, http.StatusForbidden, httpResps[2].StatusCode)
		assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resps[2])
		assert.Equal(t, http.StatusOK, httpResps[3].StatusCode)
		JsonEQ(t, ResponseWithUser(administrator), resps[3])
		assert.Equal(t, http.StatusForbidden, httpResps[4].StatusCode)
		assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resps[4])
		assert.Equal(t, http.StatusForbidden, httpResps[5].StatusCode)
		assert.Equal(t, response.ErrorResp("PERMISSION_DENIED", nil), resps[5])
	})

	t.Run("testAdministrator", func(t *testing.T) {
		groups[5].POST("/testMultipleTarget", testController, middleware.Permission("all", "targetA", "targetB"))
		httpResp := (*http.Response)(nil)
		capturer.CaptureOutput(func() {
			httpResp = MakeResp(MakeReq(t, "POST", "/administrator/testMultipleTarget", nil))
		})
		resp := response.Response{}
		MustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusInternalServerError, httpResp.StatusCode)
		assert.Equal(t, response.MakeInternalErrorResp(), resp)
	})

	t.Run("testAdministrator", func(t *testing.T) {
		base.Echo.POST("/testNonUser", testController, func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Set("user", "nonUser")
				return next(c)
			}
		}, middleware.Permission("all"))
		httpResp := (*http.Response)(nil)
		capturer.CaptureOutput(func() {
			httpResp = MakeResp(MakeReq(t, "POST", "/testNonUser", nil))
		})
		resp := response.Response{}
		MustJsonDecode(httpResp, &resp)
		assert.Equal(t, http.StatusInternalServerError, httpResp.StatusCode)
		assert.Equal(t, response.MakeInternalErrorResp(), resp)
	})

	assert.Nil(t, base.DB.DropTable(&TestClass{}).Error)
}
