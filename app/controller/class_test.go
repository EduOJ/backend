package controller_test

import (
	"fmt"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/app/response/resource"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"testing"
)

func checkInviteCode(t *testing.T, code string) {
	assert.Regexp(t, regexp.MustCompile(fmt.Sprintf("^[a-zA-Z]{%d}$", viper.GetInt("invite_code_length"))), code)
	var count int64
	assert.NoError(t, base.DB.Model(models.Class{}).Where("invite_code = ?", code).Count(&count).Error)
	assert.Equal(t, int64(1), count)
}

func TestCreateClass(t *testing.T) {
	user := createUserForTest(t, "test_create_class", 1)
	user.GrantRole("admin")
	httpResp := makeResp(makeReq(t, "POST", base.Echo.Reverse("class.createClass"), request.CreateClassRequest{
		Name:        "test_create_class_1_name",
		CourseName:  "test_create_class_1_course_name",
		Description: "test_create_class_1_description",
	}, applyUser(user)))
	assert.Equal(t, http.StatusCreated, httpResp.StatusCode)

	databaseClass := models.Class{}
	assert.NoError(t, base.DB.Preload("Managers").Preload("Students").First(&databaseClass, "name = ? ", "test_create_class_1_name").Error)
	checkInviteCode(t, databaseClass.InviteCode)
	assert.True(t, user.HasRole("class_creator", databaseClass))
	user.LoadRoles()
	databaseClass.Managers[0].LoadRoles()
	user.UpdatedAt = databaseClass.Managers[0].UpdatedAt
	expectedClass := models.Class{
		ID:          databaseClass.ID,
		Name:        "test_create_class_1_name",
		CourseName:  "test_create_class_1_course_name",
		Description: "test_create_class_1_description",
		InviteCode:  databaseClass.InviteCode,
		Managers: []models.User{
			user,
		},
		Students:  []models.User{},
		CreatedAt: databaseClass.CreatedAt,
		UpdatedAt: databaseClass.UpdatedAt,
		DeletedAt: gorm.DeletedAt{},
	}
	assert.Equal(t, expectedClass, databaseClass)
	resp := response.CreateClassResponse{}
	mustJsonDecode(httpResp, &resp)
	assert.Equal(t, response.CreateClassResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.ClassDetail `json:"class"`
		}{
			resource.GetClassDetail(&expectedClass),
		},
	}, resp)
}
