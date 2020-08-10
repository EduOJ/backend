package utils

import (
	"fmt"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/pkg/errors"
	"reflect"
	"regexp"
	"sync"
	"time"
)

var initAuth sync.Once
var SessionTimeout time.Duration
var RememberMeTimeout time.Duration
var SessionCount int

func initAuthConfig() {
	sessionTimeoutInt := config.MustGet("auth.session_timeout", 1200).Value().(int)
	SessionTimeout = time.Second * time.Duration(sessionTimeoutInt)
	RememberMeTimeoutInt := config.MustGet("auth.remember_me_timeout", 604800).Value().(int)
	RememberMeTimeout = time.Second * time.Duration(RememberMeTimeoutInt)
	SessionCount = config.MustGet("auth.session_count", 10).Value().(int)
}

func IsTokenExpired(token models.Token) bool {
	initAuth.Do(initAuthConfig)
	if token.RememberMe {
		return token.UpdatedAt.Add(RememberMeTimeout).Before(time.Now())
	} else {
		return token.UpdatedAt.Add(SessionTimeout).Before(time.Now())
	}
}

//TODO: Use this function in timed tasks
func CleanUpExpiredTokens() error {
	initAuth.Do(initAuthConfig)
	var users []models.User
	err := base.DB.Model(models.User{}).Find(&users).Error
	if err != nil {
		return errors.Wrap(err, "could not find users")
	}
	for _, user := range users {
		var tokens []models.Token
		var tokenIds []uint
		storedTokenCount := 0
		err = base.DB.Preload("User").Where("user_id = ?", user.ID).Order("updated_at desc").Find(&tokens).Error
		if err != nil {
			return errors.Wrap(err, "could not find tokens")
		}
		for _, token := range tokens {
			if storedTokenCount < SessionCount && !IsTokenExpired(token) {
				storedTokenCount++
				continue
			}
			tokenIds = append(tokenIds, token.ID)
		}
		err = base.DB.Delete(models.Token{}, "id in (?)", tokenIds).Error
		if err != nil {
			return errors.Wrap(err, "could not delete tokens")
		}
	}
	return nil
}

//The username custom validating function generator
func IsValidUsername(structField string) ValidateFunc {
	//The username custom validating function
	return func(req interface{}) (ret bool, tag string) {
		//get username from request using reflect
		username := reflect.ValueOf(req).Elem().FieldByName(structField).String()
		//validate username
		ret, err := regexp.MatchString(fmt.Sprintf("[\\w\\d_]{%d}", len(username)), username)
		if err != nil {
			panic(errors.Wrap(err, "could not finish username verification"))
		}
		if !ret {
			tag = "Illegal username (only alphas numbers and underscores are allowed to use in username)"
		}
		return
	}
}

func GetUsernameValidation(structField string) CustomValidation {
	return CustomValidation{
		requestField: "username",                   //the field in the request
		validate:     IsValidUsername(structField), //the field in the structure (usually "Username")
	}
}
