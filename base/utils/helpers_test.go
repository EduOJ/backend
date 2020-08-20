package utils

import (
	"github.com/jinzhu/gorm"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/database"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestFindUser(t *testing.T) {
	t.Cleanup(database.SetupDatabaseForTest())
	t.Run("findUserNonExist", func(t *testing.T) {
		user, err := FindUser("test_find_user_non_exist")
		assert.Nil(t, user)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
	t.Run("findUserSuccessWithId", func(t *testing.T) {
		user := models.User{
			Username: "test_find_user_id_username",
			Nickname: "test_find_user_id_nickname",
			Email:    "test_find_user_id@mail.com",
			Password: "test_find_user_id_password",
		}
		assert.Nil(t, base.DB.Create(&user).Error)
		foundUser, err := FindUser(strconv.Itoa(int(user.ID)))
		if foundUser != nil {
			assert.Equal(t, user, *foundUser)
		}
		assert.Nil(t, err)
	})
	t.Run("findUserSuccessWithUsername", func(t *testing.T) {
		user := models.User{
			Username: "test_find_user_name_username",
			Nickname: "test_find_user_name_nickname",
			Email:    "test_find_user_name@mail.com",
			Password: "test_find_user_name_password",
		}
		assert.Nil(t, base.DB.Create(&user).Error)
		foundUser, err := FindUser(user.Username)
		if foundUser != nil {
			assert.Equal(t, user, *foundUser)
		}
		assert.Nil(t, err)
	})
}
