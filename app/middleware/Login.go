package middleware

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/request"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

var tokenEffectiveTime int

func Login(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(request.LoggedRequest)
		if err := c.Bind(req); err != nil {
			panic(err)
		}
		if err := c.Validate(req); err != nil {
			if e, ok := err.(validator.ValidationErrors); ok {
				validationErrors := make([]response.ValidationError, len(e))
				for i, v := range e {
					validationErrors[i] = response.ValidationError{
						Field:  v.Field(),
						Reason: v.Tag(),
					}
				}
				return c.JSON(http.StatusBadRequest, response.ErrorResp(1, "validation error", validationErrors))
			}
			log.Error(errors.Wrap(err, "validate failed"), c)
			return response.InternalErrorResp(c)
		}
		user, err := utils.GetUserFromToken(req.Token)
		fmt.Println("err:",err,"user:",user)
		if err != nil && err != gorm.ErrRecordNotFound{
			log.Error(errors.Wrap(err, "fail to get user from token"), c)
			return response.InternalErrorResp(c)
		}
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusBadRequest, response.ErrorResp(2, "invalid token", nil))
		}
		tetStr := fmt.Sprintf("-%ds",tokenEffectiveTime)
		tetTime , err := time.ParseDuration(tetStr)
		//TODO: remember to use time.Now().UTC() at other operation time of user
		if time.Now().UTC().Add(tetTime).After(user.UpdatedAt) {
			return c.JSON(http.StatusBadRequest, response.ErrorResp(3, "outdated token", nil))
		}
		user.UpdatedAt = time.Now().UTC()
		c.Set("user",user)
		return next(c)
	}
}

func InitTokenEffectiveTime(tet int){
	tokenEffectiveTime = tet
}