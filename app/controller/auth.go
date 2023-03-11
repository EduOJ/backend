package controller

import (
	"bytes"
	"net/http"
	"time"

	"github.com/EduOJ/backend/base/log"

	"github.com/EduOJ/backend/app/request"
	"github.com/EduOJ/backend/app/response"
	"github.com/EduOJ/backend/app/response/resource"
	"github.com/EduOJ/backend/base"
	"github.com/EduOJ/backend/base/event"
	"github.com/EduOJ/backend/base/utils"
	"github.com/EduOJ/backend/database/models"
	"github.com/EduOJ/backend/event/register"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// @summary      Login into an account using email/username and password.
// @description  Login into an account using email/username and password. A token will be returned, together with the
// @description  user's personal data.
// @router       /auth/login [POST]
// @produce      json
// @tags         Auth
// @param        request  body      request.LoginRequest  true  "The login request."
// @success      200      {object}  response.LoginResponse
// @failure      500      {object}  response.Response
// @failure      400      {object}  response.Response{data=[]response.ValidationError}  "Validation error"
// @failure      404      {object}  response.Response                                   "Wrong username, with message `WRONG_USERNAME`"
// @failure      403      {object}  response.Response                                   "Wrong password, with message `WRONG_PASSWORD`"
func Login(c echo.Context) error {
	req := request.LoginRequest{}
	if err, ok := utils.BindAndValidate(&req, c); !ok {
		return err
	}
	user := models.User{}
	err := base.DB.Where("email = ? or username = ?", req.UsernameOrEmail, req.UsernameOrEmail).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("WRONG_USERNAME", nil))
		} else {
			panic(errors.Wrap(err, "could not query username or email"))
		}
	}
	if !utils.VerifyPassword(req.Password, user.Password) {
		return c.JSON(http.StatusForbidden, response.ErrorResp("WRONG_PASSWORD", nil))
	}
	token := models.Token{
		Token:      utils.RandStr(32),
		User:       user,
		RememberMe: req.RememberMe,
	}
	utils.PanicIfDBError(base.DB.Create(&token), "could not create token for users")
	if !user.RoleLoaded {
		user.LoadRoles()
	}
	return c.JSON(http.StatusOK, response.LoginResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			User  resource.UserForAdmin `json:"user"`
			Token string                `json:"token"`
		}{
			User:  *resource.GetUserForAdmin(&user),
			Token: token.Token,
		},
	})
}

// @summary      Register an account, and login into that account.
// @description  Register an account, and login into that account. A token will be returned, together with the
// @description  user's personal data.
// @router       /auth/register [POST]
// @produce      json
// @tags         Auth
// @param        request  body      request.RegisterRequest  true  "The register request."
// @success      201      {object}  response.RegisterResponse
// @failure      500      {object}  response.Response
// @failure      400      {object}  response.Response{data=[]response.ValidationError}  "Validation error"
// @failure      409      {object}  response.Response                                   "Email registered, with message `CONFLICT_EMAIL`"
// @failure      409      {object}  response.Response                                   "Username registered, with message `WRONG_PASSWORD`"
func Register(c echo.Context) error {
	req := request.RegisterRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}
	hashed := utils.HashPassword(req.Password)
	count := int64(0)
	utils.PanicIfDBError(base.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count), "could not query user count")
	if count != 0 {
		return c.JSON(http.StatusConflict, response.ErrorResp("CONFLICT_EMAIL", nil))
	}
	utils.PanicIfDBError(base.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count), "could not query user count")
	if count != 0 {
		return c.JSON(http.StatusConflict, response.ErrorResp("CONFLICT_USERNAME", nil))
	}
	user := models.User{
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    req.Email,
		Password: hashed,
	}
	utils.PanicIfDBError(base.DB.Create(&user), "could not create user")
	if _, err := event.FireEvent("register", &user); err != nil {
		panic(err)
	}
	token := models.Token{
		Token: utils.RandStr(32),
		User:  user,
	}
	utils.PanicIfDBError(base.DB.Create(&token), "could not create token for user")
	return c.JSON(http.StatusCreated, response.RegisterResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			User  resource.UserForAdmin `json:"user"`
			Token string                `json:"token"`
		}{
			User:  *resource.GetUserForAdmin(&user),
			Token: token.Token,
		},
	})
}

// @summary      EmailRegistered returns if an email is registered.
// @Description  EmailRegistered returns if an email is registered. It is mainly used for client side validation.
// @router       /auth/email_registered [GET]
// @produce      json
// @tags         Auth
// @param        email  query     string             true  "The email registered request."
// @success      200    {object}  response.Response  "Email unregistered, with message `SUCCESS`"
// @failure      500    {object}  response.Response
// @failure      400    {object}  response.Response{data=[]response.ValidationError}  "Validation error"
// @failure      409    {object}  response.Response                                   "Email registered, with message `CONFLICT_EMAIL`"
func EmailRegistered(c echo.Context) error {
	req := request.EmailRegisteredRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}
	var count int64
	utils.PanicIfDBError(base.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count), "could not query user count")
	if count != 0 {
		return c.JSON(http.StatusConflict, response.ErrorResp("EMAIL_REGISTERED", nil))
	}
	return c.JSON(http.StatusOK, response.Response{
		Message: "SUCCESS",
		Error:   nil,
		Data:    nil,
	})
}

// @summary      Request a password reset.
// @description  Request a password reset by email or username. Will check for if the user's email is
// @description  verified, then send an email with a token to reset the password. The token will be valid
// @description  for 30 minitues.
// @router       /auth/password_reset [POST]
// @produce      json
// @tags         Auth
// @param        request  body      request.RequestResetPasswordRequest                 true  "username or email"
// @success      200      {object}  response.RequestResetPasswordResponse               "email sent"
// @failure      400      {object}  response.Response{data=[]response.ValidationError}  "Validation error"
// @success      404      {object}  response.Response                                   "user not found, with message `NOT_FOUND`"
// @failure      406      {object}  response.Response                                   "Email not verified, with message `EMAIL_NOT_VERIFIED`"
// @security     ApiKeyAuth
func RequestResetPassword(c echo.Context) error {
	req := request.RequestResetPasswordRequest{}
	if err, ok := utils.BindAndValidate(&req, c); !ok {
		return err
	}
	user := models.User{}
	err := base.DB.Where("email = ? or username = ?", req.UsernameOrEmail, req.UsernameOrEmail).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		} else {
			panic(errors.Wrap(err, "could not query username or email"))
		}
	}
	if !user.EmailVerified {
		return c.JSON(http.StatusNotAcceptable, response.ErrorResp("EMAIL_NOT_VERIFIED", nil))
	}
	verification := models.EmailVerificationToken{
		User:  &user,
		Email: user.Email,
		Token: utils.RandStr(5),
		Used:  false,
	}
	if err := base.DB.Save(&verification).Error; err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	if err := base.Template.Execute(buf, map[string]string{
		"Code":     verification.Token,
		"Nickname": user.Nickname,
	}); err != nil {
		panic(err)
	}
	go func() {
		if err := utils.SendMail(user.Email, "Your email verification code for reset password", buf.String()); err != nil {
			panic(err)
		}
	}()
	return c.JSON(http.StatusOK, response.RequestResetPasswordResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data:    nil,
	})
}

// @summary      Resend a verification email.
// @description  Resend a verification email. Will check for if the user's email is already
// @description  verified, then send an email with a token to verify the email. The token will be valid
// @description  for 30 minitues.
// @router       /user/resend_email_verification [POST]
// @produce      json
// @tags         Auth
// @success      200  {object}  response.ResendEmailVerificationResponse  "email sent"
// @failure      406  {object}  response.Response                         "Email verified, with message `EMAIL_VERIFIED`"
// @security     ApiKeyAuth
func ResendEmailVerification(c echo.Context) error {
	user, ok := c.Get("user").(models.User)
	if !ok {
		panic("could not convert my user into type models.User")
	}
	if user.EmailVerified {
		return c.JSON(http.StatusNotAcceptable, response.ErrorResp("EMAIL_VERIFIED", nil))
	}
	verification := models.EmailVerificationToken{
		User:  &user,
		Email: user.Email,
		Token: utils.RandStr(5),
		Used:  false,
	}
	if err := base.DB.Save(&verification).Error; err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	if err := base.Template.Execute(buf, map[string]string{
		"Code":     verification.Token,
		"Nickname": user.Nickname,
	}); err != nil {
		panic(err)
	}
	go func() {
		if err := utils.SendMail(user.Email, "Your email verification code for email verification", buf.String()); err != nil {
			log.Fatal(err)
		}
	}()
	return c.JSON(http.StatusOK, response.ResendEmailVerificationResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data:    nil,
	})
}

// @summary      Do a password reset.
// @description  Do a password reset by email or username. Will check the if the given code is valid, then reset
// @description  the password, logging out all sessions.
// @router       /auth/password_reset [PUT]
// @produce      json
// @tags         Auth
// @param        request  body      request.DoResetPasswordRequest                      true  "username or email"
// @success      200      {object}  response.EmailVerificationResponse                  "email sent"
// @failure      400      {object}  response.Response{data=[]response.ValidationError}  "Validation error"
// @success      403      {object}  response.Response                                   "invalid token, with message `WRONG_CODE`"
// @success      408      {object}  response.Response                                   "the verification code is expired, with message `CODE_EXPIRED`"
// @success      408      {object}  response.Response                                   "the verification code is used, with message `CODE_USED`"
// @success      404      {object}  response.Response                                   "user not found, with message `NOT_FOUND`"
// @security     ApiKeyAuth
func DoResetPassword(c echo.Context) error {
	req := request.DoResetPasswordRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}
	user := models.User{}
	err = base.DB.Where("email = ? or username = ?", req.UsernameOrEmail, req.UsernameOrEmail).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, response.ErrorResp("NOT_FOUND", nil))
		} else {
			panic(errors.Wrap(err, "could not query username or email"))
		}
	}
	var code models.EmailVerificationToken
	err = base.DB.Where("user_id = ? and token = ? and email = ?", user.ID, req.Token, user.Email).First(&code).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusUnauthorized, response.ErrorResp("WRONG_CODE", nil))
		} else {
			panic(err)
		}
	}
	if code.CreatedAt.Before(time.Now().Add(-30 * time.Minute)) {
		return c.JSON(http.StatusRequestTimeout, response.ErrorResp("CODE_EXPIRED", nil))
	}
	if code.Used {
		return c.JSON(http.StatusRequestTimeout, response.ErrorResp("CODE_USED", nil))
	}
	code.Used = true
	utils.PanicIfDBError(base.DB.Save(&code), "could not save verification code")
	user.Password = utils.HashPassword(req.Password)
	utils.PanicIfDBError(base.DB.Save(&user), "could not save user")
	base.DB.Where("user_id = ?", user.ID).Delete(models.Token{}) // logout existing user
	return c.JSON(http.StatusOK, response.EmailVerificationResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data:    nil,
	})
}

// @summary      Verify a user's email.
// @description  Verify a user's email. Will check for if the user's email is
// @description  verified, then send an email with a token to verify The token will be valid
// @description  for 30 minitues.
// @router       /user/email_verification [POST]
// @produce      json
// @tags         Auth
// @param        request  body      request.VerifyEmailRequest                      true "token"
// @success      200      {object}  response.EmailVerificationResponse                  "email sent"
// @failure      400      {object}  response.Response{data=[]response.ValidationError}  "Validation error"
// @success      403      {object}  response.Response                                   "invalid token, with message `WRONG_CODE`"
// @success      408      {object}  response.Response                                   "the verification code is expired, with message `CODE_EXPIRED`"
// @success      408      {object}  response.Response                                   "the verification code is used, with message `CODE_USED`"
// @success      404      {object}  response.Response                                   "user not found, with message `NOT_FOUND`"
// @security     ApiKeyAuth
func VerifyEmail(c echo.Context) error {
	req := request.VerifyEmailRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}
	user := c.Get("user").(models.User)
	if user.EmailVerified {
		return c.JSON(http.StatusNotAcceptable, response.ErrorResp("EMAIL_VERIFIED", nil))
	}
	var code models.EmailVerificationToken
	err = base.DB.Where("user_id = ? and token = ? and email = ?", user.ID, req.Token, user.Email).First(&code).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusUnauthorized, response.ErrorResp("WRONG_CODE", nil))
		} else {
			panic(err)
		}
	}
	if code.CreatedAt.Before(time.Now().Add(-30 * time.Minute)) {
		return c.JSON(http.StatusRequestTimeout, response.ErrorResp("CODE_EXPIRED", nil))
	}
	if code.Used {
		return c.JSON(http.StatusRequestTimeout, response.ErrorResp("CODE_USED", nil))
	}
	code.Used = true
	utils.PanicIfDBError(base.DB.Save(&code), "could not save verification code")
	user.EmailVerified = true
	utils.PanicIfDBError(base.DB.Save(&user), "could not save user")
	return c.JSON(http.StatusOK, response.EmailVerificationResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data:    nil,
	})
}

// @summary      Update current user's email if not verified.
// @description  Change current user's email only if the email is not verified.
// @description  The new email can not be the same as other users'.
// @router       /user/update_email [PUT]
// @produce      json
// @tags         Auth
// @param        request  body      request.UpdateEmailRequest  true  "New email"
// @success      200      {object}  response.UpdateEmailResponse
// @failure      406      {object}  response.Response  "Email verified, with message `EMAIL_VERIFIED`"
// @failure      409      {object}  response.Response  "New email confilct, with message `CONFLICT_EMAIL`"
// @security     ApiKeyAuth
func UpdateEmail(c echo.Context) error {
	user, ok := c.Get("user").(models.User)
	if !ok {
		panic("could not convert my user into type models.User")
	}
	req := request.UpdateEmailRequest{}
	err, ok := utils.BindAndValidate(&req, c)
	if !ok {
		return err
	}
	if user.EmailVerified {
		return c.JSON(http.StatusNotAcceptable, response.ErrorResp("EMAIL_VERIFIED", nil))
	}
	count := int64(0)
	utils.PanicIfDBError(base.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count), "could not query user count")
	if count > 1 || (count == 1 && user.Email != req.Email) {
		return c.JSON(http.StatusConflict, response.ErrorResp("CONFLICT_EMAIL", nil))
	}
	user.EmailVerified = false
	base.DB.Delete(&models.EmailVerificationToken{}, "user_id = ?", user.ID)
	user.Email = req.Email
	utils.PanicIfDBError(base.DB.Omit(clause.Associations).Save(&user), "could not update email")
	register.SendVerificationEmail(&user)
	return c.JSON(http.StatusOK, response.UpdateEmailResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			*resource.UserForAdmin `json:"user"`
		}{
			resource.GetUserForAdmin(&user),
		},
	})
}
