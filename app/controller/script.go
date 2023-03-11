package controller

import (
	"net/http"

	"github.com/EduOJ/backend/app/response"
	"github.com/EduOJ/backend/base"
	"github.com/EduOJ/backend/base/utils"
	"github.com/EduOJ/backend/database/models"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gorm.io/gorm"
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
