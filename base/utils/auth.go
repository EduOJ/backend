package utils

import (
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/pkg/errors"
	"time"
)

var initialized bool
var SessionTimeout time.Duration
var RememberMeTimeout time.Duration
var SessionCount int

func InitAuthConfig() {
	sessionTimeoutInt := config.MustGet("auth.session_timeout", 1200).Value().(int)
	SessionTimeout = time.Second * time.Duration(sessionTimeoutInt)
	RememberMeTimeoutInt := config.MustGet("auth.remember_me_timeout", 604800).Value().(int)
	RememberMeTimeout = time.Second * time.Duration(RememberMeTimeoutInt)
	SessionCount = config.MustGet("auth.session_count", 10).Value().(int)
	initialized = true
}

func IsTokenExpired(token models.Token) bool {
	if !initialized {
		InitAuthConfig()
	}
	var timeout time.Duration
	if token.RememberMe {
		timeout = RememberMeTimeout
	} else {
		timeout = SessionTimeout
	}
	return token.UpdatedAt.Add(timeout).Before(time.Now())
}

//TODO: Use this function in timed tasks
func CleanUpExpiredTokens() error {
	InitAuthConfig()
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
