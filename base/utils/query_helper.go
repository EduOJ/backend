package utils

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

type HttpError struct {
	Code    int
	Message string
	Err     error
}

func (e HttpError) Error() string {
	return fmt.Sprintf("[%d]%s", e.Code, e.Message)
}

func (e HttpError) Response(c echo.Context) error {
	return c.JSON(e.Code, response.ErrorResp(e.Message, e.Err))
}

func Paginator(query *gorm.DB, limit, offset int, requestURL *url.URL, output interface{}) (total int, prevUrl, nextUrl *string, err error) {
	if limit == 0 {
		limit = 20 // Default limit
	}
	err = query.Count(&total).Error
	if err != nil {
		err = errors.Wrap(err, "could not query count of objects")
		return
	}
	err = query.Limit(limit).Offset(offset).Find(output).Error
	if err != nil {
		err = errors.Wrap(err, "could not query objects")
		return
	}
	count := reflect.ValueOf(output).Elem().Len()

	if offset-limit >= 0 {
		prevURL := requestURL
		q, err := url.ParseQuery(prevURL.RawQuery)
		if err != nil {
			panic(errors.Wrap(err, "could not parse query for url"))
		}
		q.Set("offset", fmt.Sprint(offset-limit))
		q.Set("limit", fmt.Sprint(limit))
		prevURL.RawQuery = q.Encode()
		temp := prevURL.String()
		prevUrl = &temp
	} else {
		prevUrl = nil
	}
	if offset+count < total {
		nextURL := requestURL
		q, err := url.ParseQuery(nextURL.RawQuery)
		if err != nil {
			panic(errors.Wrap(err, "could not parse query for url"))
		}
		q.Set("offset", fmt.Sprint(offset+limit))
		q.Set("limit", fmt.Sprint(limit))
		nextURL.RawQuery = q.Encode()
		temp := nextURL.String()
		nextUrl = &temp
	} else {
		nextUrl = nil
	}
	return
}

func Sorter(query *gorm.DB, orderBy string, columns ...string) (*gorm.DB, error) {
	if orderBy != "" {
		order := strings.SplitN(orderBy, ".", 2)
		if len(order) != 2 {
			return nil, HttpError{
				Code:    http.StatusBadRequest,
				Message: "INVALID_ORDER",
			}
		}
		if !Contain(order[0], columns) {
			return nil, HttpError{
				Code:    http.StatusBadRequest,
				Message: "INVALID_ORDER",
			}
		}
		if !Contain(order[1], []string{"ASC", "DESC"}) {
			return nil, HttpError{
				Code:    http.StatusBadRequest,
				Message: "INVALID_ORDER",
			}
		}
		query = query.Order(strings.Join(order, " "))
	}
	return query, nil
}

func FindUser(id string) (*models.User, error) {
	user := models.User{}
	err := base.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		err = base.DB.Where("username = ?", id).First(&user).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, err
			} else {
				panic(errors.Wrap(err, "could not query user"))
			}
		}
	}
	return &user, nil
}
