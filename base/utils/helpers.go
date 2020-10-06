package utils

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/log"
	validator2 "github.com/leoleoasd/EduOJBackend/base/validator"
	"github.com/minio/minio-go"
	"github.com/pkg/errors"
	"io"
	"mime/multipart"
	"net/http"
)

func PanicIfDBError(db *gorm.DB, message string) {
	if db.Error != nil {
		panic(errors.Wrap(db.Error, message))
	}
}

func BindAndValidate(req interface{}, c echo.Context) (err error, ok bool) {
	if err := c.Bind(req); err != nil {
		panic(err)
	}
	if err := c.Validate(req); err != nil {
		if e, ok := err.(validator.ValidationErrors); ok {
			validationErrors := make([]response.ValidationError, len(e))
			for i, v := range e {
				validationErrors[i] = response.ValidationError{
					Field:       v.Field(),
					Reason:      v.Tag(),
					Translation: v.Translate(validator2.Trans),
				}
			}
			return c.JSON(http.StatusBadRequest, response.ErrorResp("VALIDATION_ERROR", validationErrors)), false
		}
		log.Error(errors.Wrap(err, "validate failed"), c)
		return response.InternalErrorResp(c), false
	}
	return nil, true
}

func MustPutObject(object *multipart.FileHeader, ctx context.Context, bucket string, path string) {
	src, err := object.Open()
	if err != nil {
		panic(err)
	}
	defer func() {
		err := src.Close()
		if err != nil {
			panic(errors.Wrap(err, "could not close file reader"))
		}
	}()
	_, err = base.Storage.PutObjectWithContext(ctx, bucket, path, src, object.Size, minio.PutObjectOptions{})
	if err != nil {
		panic(errors.Wrap(err, "could write file to s3 storage."))
	}
}

func MustGetObject(bucket string, path string) *minio.Object {
	object, err := base.Storage.GetObject(bucket, path, minio.GetObjectOptions{})
	if err != nil {
		panic(err)
	}
	_, err = object.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}
	return object
}
