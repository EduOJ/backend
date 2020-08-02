package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"net/http"
	"strconv"
)

type hasRole struct {
	ID   uint
	Name string
}

func (h *hasRole) GetID() uint {
	return h.ID
}
func (h *hasRole) TypeName() string {
	return h.Name
}

func HasPermission(perm string, targets ...string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_u := c.Get("user")
			if _u == nil {
				log.Fatalf("%s don't have a login check middleware!!", c.Path())
				return response.InternalErrorResp(c)
			}
			if u, ok := _u.(models.User); ok {
				can := false
				if len(targets) == 0 {
					can = u.Can(perm)
				} else if len(targets) == 1 {
					id, err := strconv.ParseUint(c.Param("id"), 10, strconv.IntSize)
					if err != nil {
						log.Errorf("illegal id: %s", c.Param("id"))
					}
					can = u.Can(perm, &hasRole{
						ID:   uint(id),
						Name: targets[0],
					})
				} else {
					log.Fatalf("%s registered permission middleware with more than one permission!", c.Path())
					return response.InternalErrorResp(c)
				}
				if can {
					return next(c)
				} else {
					return c.JSON(http.StatusForbidden, response.ErrorResp("PERMISSION_DENIED", nil))
				}
			}
			log.Fatalf("%s's user is not a models.User!", c.Path())
			return response.InternalErrorResp(c)
		}
	}
}
