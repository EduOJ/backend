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
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"net/http"
)

func generateInviteCode() (code string) {
	var classes []models.Class
	utils.PanicIfDBError(base.DB.Select("invite_code").Find(&classes), "could not find classes for generating invite codes")
	crashed := true
	for crashed {
		code = utils.RandStr(viper.GetInt("invite_code_length"))
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
		InviteCode:  generateInviteCode(),
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
	if err := base.DB.Preload("managers").Preload("students").First(&class, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		}
		panic(errors.Wrap(err, "could not get class while getting class"))
	}
	user := c.Get("user").(models.User)
	if user.Can("read_class", class) || user.Can("manage_class") {
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
		return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
	}
}

func GetClassesIManage(c echo.Context) error {
	user := c.Get("user").(models.User)
	var classes []models.Class
	if err := base.DB.Model(&user).Association("classes_managing").Find(&classes); err != nil {
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
	if err := base.DB.Model(&user).Association("classes_taking").Find(&classes); err != nil {
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
	if err := base.DB.Preload("managers").Preload("students").First(&class, c.Param("id")).Error; err != nil {
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
	return c.JSON(http.StatusCreated, response.UpdateClassResponse{
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
	req := request.RefreshInviteCodeRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}
	class := models.Class{}
	if err := base.DB.Preload("managers").Preload("students").First(&class, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		} else {
			panic(errors.Wrap(err, "could not find class for refreshing invite code"))
		}
	}
	class.Name = generateInviteCode()
	utils.PanicIfDBError(base.DB.Save(&class), "could not update class for refreshing invite code")
	return c.JSON(http.StatusCreated, response.RefreshInviteCodeResponse{
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
	if err := base.DB.First(&class, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		} else {
			panic(errors.Wrap(err, "could not find class for adding students"))
		}
	}
	if err := class.AddStudents(req.UserIds); err != nil {
		panic(errors.Wrap(err, "could not add students"))
	}
	return c.JSON(http.StatusCreated, response.AddStudentsResponse{
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
	if err := base.DB.First(&class, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		} else {
			panic(errors.Wrap(err, "could not find class for deleting students"))
		}
	}
	if err := class.DeleteStudents(req.UserIds); err != nil {
		panic(errors.Wrap(err, "could not delete students"))
	}
	return c.JSON(http.StatusCreated, response.DeleteStudentsResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.ClassDetail `json:"class"`
		}{
			resource.GetClassDetail(&class),
		},
	})
}

func DeleteClass(c echo.Context) error {
	req := request.UpdateClassRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}
	class := models.Class{}
	if err := base.DB.Delete(&class, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		} else {
			panic(errors.Wrap(err, "could not find class for deleting"))
		}
	}
	return c.JSON(http.StatusCreated, response.Response{
		Message: "SUCCESS",
		Error:   nil,
		Data:    nil,
	})
}
