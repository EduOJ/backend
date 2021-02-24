package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/app/response/resource"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"net/http"
	"sync"
)

var inviteCodeLock sync.Mutex

func GenerateInviteCode() (code string) {
	inviteCodeLock.Lock()
	defer inviteCodeLock.Unlock()
	var classes []models.Class
	utils.PanicIfDBError(base.DB.Select("invite_code").Find(&classes), "could not find classes for generating invite codes")
	crashed := true
	for crashed {
		// 5: Fixed invite code length
		code = utils.RandStr(5)
		crashed = false
		for _, c := range classes {
			if c.InviteCode == code {
				crashed = true
				continue
			}
		}
	}
	return
}

func CreateClass(c echo.Context) error {
	user := c.Get("user").(models.User)
	req := request.CreateClassRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}
	class := models.Class{
		Name:        req.Name,
		CourseName:  req.CourseName,
		Description: req.Description,
		InviteCode:  GenerateInviteCode(),
		Managers: []models.User{
			user,
		},
		Students: []models.User{},
	}
	utils.PanicIfDBError(base.DB.Create(&class), "could not create class")
	user.GrantRole("class_creator", class)
	return c.JSON(http.StatusCreated, response.CreateClassResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.ClassDetail `json:"class"`
		}{
			resource.GetClassDetail(&class),
		},
	})
}

func GetClass(c echo.Context) error {
	class := models.Class{}
	if err := base.DB.Preload("Managers").Preload("Students").First(&class, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		}
		panic(errors.Wrap(err, "could not get class while getting class"))
	}
	user := c.Get("user").(models.User)
	if user.Can("read_class_secrets", class) || user.Can("read_class_secrets") {
		return c.JSON(http.StatusOK, response.GetClassResponseForAdmin{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*resource.ClassDetail `json:"class"`
			}{
				resource.GetClassDetail(&class),
			},
		})
	} else if user.In(class.Students) {
		return c.JSON(http.StatusOK, response.GetClassResponse{
			Message: "SUCCESS",
			Error:   nil,
			Data: struct {
				*resource.Class `json:"class"`
			}{
				resource.GetClass(&class),
			},
		})
	} else {
		return c.JSON(http.StatusForbidden, response.ErrorResp("PERMISSION_DENIED", nil))
	}
}

func GetClassesIManage(c echo.Context) error {
	user := c.Get("user").(models.User)
	var classes []models.Class
	if err := base.DB.Model(&user).Preload("Managers").Preload("Students").Association("ClassesManaging").Find(&classes); err != nil {
		panic(errors.Wrap(err, "could not find class managing"))
	}
	return c.JSON(http.StatusOK, response.GetClassesIManageResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			Classes []resource.Class `json:"classes"`
		}{
			Classes: resource.GetClassSlice(classes),
		},
	})
}

func GetClassesITake(c echo.Context) error {
	user := c.Get("user").(models.User)
	var classes []models.Class
	if err := base.DB.Model(&user).Preload("Managers").Preload("Students").Association("ClassesTaking").Find(&classes); err != nil {
		panic(errors.Wrap(err, "could not find class taking"))
	}
	return c.JSON(http.StatusOK, response.GetClassesITakeResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			Classes []resource.Class `json:"classes"`
		}{
			Classes: resource.GetClassSlice(classes),
		},
	})
}

func UpdateClass(c echo.Context) error {
	req := request.UpdateClassRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}
	class := models.Class{}
	if err := base.DB.Preload("Managers").Preload("Students").First(&class, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		} else {
			panic(errors.Wrap(err, "could not find class for updating"))
		}
	}
	class.Name = req.Name
	class.CourseName = req.CourseName
	class.Description = req.Description
	utils.PanicIfDBError(base.DB.Save(&class), "could not update class")
	return c.JSON(http.StatusOK, response.UpdateClassResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.ClassDetail `json:"class"`
		}{
			resource.GetClassDetail(&class),
		},
	})
}

func RefreshInviteCode(c echo.Context) error {
	class := models.Class{}
	if err := base.DB.Preload("Managers").Preload("Students").First(&class, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		} else {
			panic(errors.Wrap(err, "could not find class for refreshing invite code"))
		}
	}
	class.InviteCode = GenerateInviteCode()
	utils.PanicIfDBError(base.DB.Save(&class), "could not update class for refreshing invite code")
	return c.JSON(http.StatusOK, response.RefreshInviteCodeResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.ClassDetail `json:"class"`
		}{
			resource.GetClassDetail(&class),
		},
	})
}

func AddStudents(c echo.Context) error {
	req := request.AddStudentsRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}
	class := models.Class{}
	if err := base.DB.Preload("Managers").Preload("Students").First(&class, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		} else {
			panic(errors.Wrap(err, "could not find class for adding students"))
		}
	}
	databaseIds := make([]uint, len(class.Students))
	for i, s := range class.Students {
		databaseIds[i] = s.ID
	}
	ids := utils.IdUniqueInA(req.UserIds, databaseIds)
	if err := class.AddStudents(ids); err != nil {
		panic(errors.Wrap(err, "could not add students"))
	}
	return c.JSON(http.StatusOK, response.AddStudentsResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.ClassDetail `json:"class"`
		}{
			resource.GetClassDetail(&class),
		},
	})
}

func DeleteStudents(c echo.Context) error {
	req := request.DeleteStudentsRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}
	class := models.Class{}
	if err := base.DB.Preload("Managers").Preload("Students").First(&class, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		} else {
			panic(errors.Wrap(err, "could not find class for deleting students"))
		}
	}
	if err := class.DeleteStudents(req.UserIds); err != nil {
		panic(errors.Wrap(err, "could not delete students"))
	}
	return c.JSON(http.StatusOK, response.DeleteStudentsResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.ClassDetail `json:"class"`
		}{
			resource.GetClassDetail(&class),
		},
	})
}

func JoinClass(c echo.Context) error {
	req := request.JoinClassRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}
	class := models.Class{}
	if err := base.DB.Preload("Managers").Preload("Students").First(&class, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		} else {
			panic(errors.Wrap(err, "could not find class for deleting students"))
		}
	}
	if class.InviteCode != req.InviteCode {
		return c.JSON(http.StatusForbidden, response.ErrorResp("WRONG_INVITE_CODE", nil))
	}
	user := c.Get("user").(models.User)
	if user.In(class.Students) {
		return c.JSON(http.StatusBadRequest, response.ErrorResp("ALREADY_IN_CLASS", nil))
	}
	if err := base.DB.Model(&class).Association("Students").Append(&user); err != nil {
		panic(errors.Wrap(err, "could not add student for joining class"))
	}
	return c.JSON(http.StatusOK, response.JoinClassResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.Class `json:"class"`
		}{
			resource.GetClass(&class),
		},
	})
}

func DeleteClass(c echo.Context) error {
	class := models.Class{}
	if err := base.DB.Delete(&class, c.Param("id")).Error; err != nil {
		panic(errors.Wrap(err, "could not find class for deleting"))
	}
	return c.JSON(http.StatusOK, response.Response{
		Message: "SUCCESS",
		Error:   nil,
		Data:    nil,
	})
}
