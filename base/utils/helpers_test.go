package utils

import (
	"bytes"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/leoleoasd/EduOJBackend/base/exit"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/leoleoasd/EduOJBackend/database"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"
)

func jsonEQ(t *testing.T, expected, actual interface{}) {
	assert.JSONEq(t, mustJsonEncode(t, expected), mustJsonEncode(t, actual))
}

func mustJsonEncode(t *testing.T, data interface{}) string {
	var err error
	if dataResp, ok := data.(*http.Response); ok {
		data, err = ioutil.ReadAll(dataResp.Body)
		assert.Equal(t, nil, err)
	}
	if dataString, ok := data.(string); ok {
		data = []byte(dataString)
	}
	if dataBytes, ok := data.([]byte); ok {
		err := json.Unmarshal(dataBytes, &data)
		assert.Equal(t, nil, err)
	}
	j, err := json.Marshal(data)
	if err != nil {
		t.Fatal(data, err)
	}
	return string(j)
}

func TestMain(m *testing.M) {
	defer database.SetupDatabaseForTest()()
	defer exit.SetupExitForTest()()
	configFile := bytes.NewBufferString(`debug: false
server:
  port: 8080
  origin:
    - http://127.0.0.1:8000
auth:
  session_timeout: 1200
  remember_me_timeout: 604800
  session_count: 10`)
	err := config.ReadConfig(configFile)
	if err != nil {
		panic(err)
	}
	log.Disable()

	os.Exit(m.Run())
}

func TestFindUser(t *testing.T) {
	t.Run("findUserNonExist", func(t *testing.T) {
		t.Parallel()
		user, err := FindUser("test_find_user_non_exist")
		assert.Nil(t, user)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
	t.Run("findUserSuccessWithId", func(t *testing.T) {
		t.Parallel()
		user := models.User{
			Username: "test_find_user_id_username",
			Nickname: "test_find_user_id_nickname",
			Email:    "test_find_user_id@mail.com",
			Password: "test_find_user_id_password",
		}
		assert.Nil(t, base.DB.Create(&user).Error)
		foundUser, err := FindUser(strconv.Itoa(int(user.ID)))
		jsonEQ(t, user, foundUser)
		assert.Nil(t, err)
	})
	t.Run("findUserSuccessWithUsername", func(t *testing.T) {
		t.Parallel()
		user := models.User{
			Username: "test_find_user_name_username",
			Nickname: "test_find_user_name_nickname",
			Email:    "test_find_user_name@mail.com",
			Password: "test_find_user_name_password",
		}
		assert.Nil(t, base.DB.Create(&user).Error)
		foundUser, err := FindUser(user.Username)
		jsonEQ(t, user, foundUser)
		assert.Nil(t, err)
	})
}
