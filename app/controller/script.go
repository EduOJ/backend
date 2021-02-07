package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"net/http"
)

func GetScript(c echo.Context) error {
	script := models.Script{}
	name := c.Param("name")
	err := base.DB.First(&script, "name = ?", name).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		} else {
			panic(errors.Wrap(err, "could not query script"))
		}
	}
	url, err := utils.GetPresignedURL("scripts", script.Name, script.Filename)
	if err != nil {
		panic(errors.Wrap(err, "could not get presigned url of script"))
	}
	return c.Redirect(http.StatusFound, url)
}
