package admin

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	adminRequest "github.com/leoleoasd/EduOJBackend/app/request/admin"
	"github.com/leoleoasd/EduOJBackend/app/response"
	adminResponse "github.com/leoleoasd/EduOJBackend/app/response/admin"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"time"
)

func PostUser(c echo.Context) error {
	req := new(adminRequest.PostUserRequest)
	err, ok := utils.BindAndValidate(req, &c, utils.GetUsernameValidator("Username"))
	if !ok {
		return err
	}
	hashed := utils.HashPassword(req.Password)
	count := 0
	utils.PanicIfDBError(base.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count), "could not query user count")
	if count != 0 {
		return c.JSON(http.StatusBadRequest, response.ErrorResp("POST_USER_DUPLICATE_EMAIL", nil))
	}
	utils.PanicIfDBError(base.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count), "could not query user count")
	if count != 0 {
		return c.JSON(http.StatusBadRequest, response.ErrorResp("POST_USER_DUPLICATE_USERNAME", nil))
	}
	user := models.User{
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    req.Email,
		Password: hashed,
	}
	utils.PanicIfDBError(base.DB.Create(&user), "could not create user")
	token := models.Token{
		Token: utils.RandStr(32),
		User:  user,
	}
	utils.PanicIfDBError(base.DB.Create(&token), "could not create token for user")
	return c.JSON(http.StatusCreated, adminResponse.PostUserResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			models.User `json:"user"`
		}{
			user,
		},
	})
}

func PutUser(c echo.Context) error {
	req := new(adminRequest.PutUserRequest)
	err, ok := utils.BindAndValidate(req, &c, utils.GetUsernameValidator("Username"))
	if !ok {
		return err
	}
	user, err, ok := queryUser(c)
	if !ok {
		return err
	}
	count := 0
	utils.PanicIfDBError(base.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count), "could not query user count")
	if count > 1 || (count == 1 && user.Email != req.Email) {
		return c.JSON(http.StatusBadRequest, response.ErrorResp("PUT_USER_DUPLICATE_EMAIL", nil))
	}
	utils.PanicIfDBError(base.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count), "could not query user count")
	if count > 1 || (count == 1 && user.Username != req.Username) {
		return c.JSON(http.StatusBadRequest, response.ErrorResp("PUT_USER_DUPLICATE_USERNAME", nil))
	}
	hashed := utils.HashPassword(req.Password)
	user.Username = req.Username
	user.Nickname = req.Nickname
	user.Email = req.Email
	user.Password = hashed
	user.UpdatedAt = time.Now()
	utils.PanicIfDBError(base.DB.Save(&user), "could not update user")
	return c.JSON(http.StatusOK, adminResponse.PutUserResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			models.User `json:"user"`
		}{
			user,
		},
	})
}

func DeleteUser(c echo.Context) error {
	user, err, ok := queryUser(c)
	if !ok {
		return err
	}
	utils.PanicIfDBError(base.DB.Delete(&user), "could not delete user")
	return c.JSON(http.StatusOK, response.Response{
		Message: "SUCCESS",
		Error:   nil,
		Data:    nil,
	})
}

func GetUser(c echo.Context) error {

	user, err, ok := queryUser(c)
	if !ok {
		return err
	}
	return c.JSON(http.StatusOK, adminResponse.GetUserResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			models.User `json:"user"`
		}{
			user,
		},
	})
}

func GetUsers(c echo.Context) error {
	req := new(adminRequest.GetUsersRequest)
	err, ok := utils.BindAndValidate(req, &c, utils.GetUsernameValidator("Username"))
	if !ok {
		return err
	}
	var users []models.User
	var total int

	//if req.OrderBy != "" {
	//	err = base.DB.Where("username like ? and nickname like ?", "%"+req.Username+"%", "%"+req.Nickname+"%").Order(req.OrderBy).Find(&users).Error
	//} else {
	//	err = base.DB.Where("username like ? and nickname like ?", "%"+req.Username+"%", "%"+req.Nickname+"%").Find(&users).Error
	//}
	query := base.DB.Model(&models.User{})
	if req.OrderBy != "" {
		order := strings.SplitN(req.OrderBy, ".", 2)
		if len(order) != 2 {
			return c.JSON(http.StatusBadRequest, response.ErrorResp("INVALID_ORDER", nil))
		}
		if !utils.Contain(order[0], []string{"username", "id", "nickname"}) {
			return c.JSON(http.StatusBadRequest, response.ErrorResp("INVALID_ORDER", nil))
		}
		if !utils.Contain(order[1], []string{"ASC", "DESC"}) {
			return c.JSON(http.StatusBadRequest, response.ErrorResp("INVALID_ORDER", nil))
		}
		query = query.Order(strings.Join(order, " "))
	}
	if req.Search != "" {
		query = query.Where("id like %?% or username like %?% or email like %?% or nickname like %?%", req.Search, req.Search, req.Search, req.Search)
	}
	if req.Limit == 0 {
		req.Limit = 20 // Default limit
	}
	err = query.Limit(req.Limit).Offset(req.Offset).Find(&users).Error
	if err != nil {
		panic(errors.Wrap(err, "could not query users"))
	}
	err = query.Count(&total).Error
	if err != nil {
		panic(errors.Wrap(err, "could not query count of users"))
	}

	prevURL := c.Request().URL
	if req.Offset-req.Limit >= 0 {
		prevURL.Query().Set("offset", fmt.Sprint(req.Offset-req.Limit))
	} else {
		prevURL.Query().Set("offset", "0")
	}
	prevURL.Query().Set("limit", fmt.Sprint(req.Limit))

	nextURL := c.Request().URL
	nextURL.Query().Set("offset", fmt.Sprint(req.Offset+req.Limit))
	prevURL.Query().Set("limit", fmt.Sprint(req.Limit))

	return c.JSON(http.StatusOK, adminResponse.GetUsersResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			Users  []models.User `json:"users"`
			Total  int           `json:"total"`
			Count  int           `json:"count"`
			Offset int           `json:"offset"`
			Prev   string        `json:"prev"`
			Next   string        `json:"next"`
		}{
			users,
			total,
			len(users),
			req.Offset,
			prevURL.String(),
			nextURL.String(),
		},
	})
}

func queryUser(c echo.Context) (models.User, error, bool) {
	id := c.Param("id")
	user := models.User{}
	err := base.DB.Where("id = ?", id).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		err = base.DB.Where("username = ?", id).First(&user).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return user, c.JSON(http.StatusNotFound, response.ErrorResp("QUERY_USER_WRONG_ID", nil)), false
			} else {
				panic(errors.Wrap(err, "could not query username"))
			}
		}
	} else if err != nil {
		panic(errors.Wrap(err, "could not query id"))
	}
	return user, nil, true
}

//TODO:add tests
