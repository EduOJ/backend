package utils

import (
	"bufio"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/EduOJ/backend/app/response"
	"github.com/EduOJ/backend/base"
	"github.com/EduOJ/backend/base/log"
	validator2 "github.com/EduOJ/backend/base/validator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func PanicIfDBError(db *gorm.DB, message string) {
	if db.Error != nil {
		panic(errors.Wrap(db.Error, message))
	}
}

func BindAndValidate(req interface{}, c echo.Context) (err error, ok bool) {
	// TODO: return a HttpError instead of writing to response.
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.ErrorResp("BAD_REQUEST_PARAMETER", nil)), false
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
	// TODO: use reader
	src, err := object.Open()
	if err != nil {
		panic(err)
	}
	defer src.Close()
	_, err = base.Storage.PutObject(ctx, bucket, path, src, object.Size, minio.PutObjectOptions{})
	if err != nil {
		panic(errors.Wrap(err, "could write file to s3 storage."))
	}
}

func MustPutTestCase(sanitize bool, object *multipart.FileHeader, ctx context.Context, bucket string, path string) {
	src, err := object.Open()
	if err != nil {
		panic(err)
	}
	defer src.Close()

	var fileSize int64

	if sanitize {
		scanner := bufio.NewScanner(src)
		tempFile, err := os.CreateTemp("", "tempFile*.txt")
		if err != nil {
			panic(err)
		}
		writer := bufio.NewWriter(tempFile)
		for scanner.Scan() {
			line := strings.ReplaceAll(scanner.Text(), "\r\n", "\n") // replace '\r\n' to '\n'
			_, err := fmt.Fprint(writer, line)
			if !strings.HasSuffix(line, "\n") {
				line += "\n"
			}
			if err != nil {
				panic(err)
			}
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
		if err := writer.Flush(); err != nil {
			panic(err)
		}

		fileInfo, err := tempFile.Stat()
		if err != nil {
			panic(err)
		}
		fileSize = fileInfo.Size()

		src, err = os.Open(tempFile.Name())
		if err != nil {
			panic(err)
		}
		defer src.Close()
	} else {
		fileSize = object.Size
	}

	_, err = base.Storage.PutObject(ctx, bucket, path, src, fileSize, minio.PutObjectOptions{})
	if err != nil {
		panic(errors.Wrap(err, "could write file to s3 storage."))
	}
}

func MustGetObject(c echo.Context, bucket string, path string) *minio.Object {
	object, err := base.Storage.GetObject(c.Request().Context(), bucket, path, minio.GetObjectOptions{})
	if err != nil {
		panic(err)
	}
	return object
}
