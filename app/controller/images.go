package controller

import (
	"github.com/gabriel-vasile/mimetype"
	"github.com/labstack/echo/v4"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/leoleoasd/EduOJBackend/base/log"
	"github.com/leoleoasd/EduOJBackend/base/utils"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/minio/minio-go"
	"github.com/pkg/errors"
	"io"
	"net/http"
	path2 "path"
)

func CreateImage(c echo.Context) error {
	user, ok := c.Get("user").(models.User)
	if !ok {
		panic("could not convert my user into type models.User")
	}
	file, err := c.FormFile("file")
	if err != nil {
		panic(err)
	}
	count := 1
	fileIndex := ""
	for count != 0 {
		fileIndex = utils.RandStr(32)
		utils.PanicIfDBError(base.DB.Model(&models.Image{}).Where("file_path = ?", fileIndex).Count(&count), "could not save image")
	}
	fileModel := models.Image{
		Filename: file.Filename,
		FilePath: fileIndex,
		User:     user,
	}
	utils.PanicIfDBError(base.DB.Save(&fileModel), "could not save image")
	src, err := file.Open()
	if err != nil {
		panic(err)
	}
	defer src.Close()
	mime, err := mimetype.DetectReader(src)
	if err != nil {
		log.Error("could not detect MIME of image!")
		log.Error(err)
		return c.JSON(http.StatusForbidden, response.CreateImageResponse{
			Message: "ILLEGAL_TYPE",
			Error:   nil,
			Data: struct {
				FilePath *string `json:"filename"`
			}{},
		})
	}
	if mime.String()[:5] != "image" || mime.Extension() != path2.Ext(file.Filename) {
		return c.JSON(http.StatusForbidden, response.CreateImageResponse{
			Message: "ILLEGAL_TYPE",
			Error:   nil,
			Data: struct {
				FilePath *string `json:"filename"`
			}{},
		})
	}
	_, err = src.Seek(0, io.SeekStart)
	if err != nil {
		panic(errors.Wrap(err, "could not seek to file start"))
	}

	found, err := base.Storage.BucketExists("images")
	if err != nil {
		panic(errors.Wrap(err, "could not query if bucket exists"))
	}
	if !found {
		err = base.Storage.MakeBucket("images", config.MustGet("storage.region", "us-east-1").String())
		if err != nil {
			panic(errors.Wrap(err, "could not query if bucket exists"))
		}
	}
	_, err = base.Storage.PutObjectWithContext(c.Request().Context(), "images", fileIndex, src, file.Size, minio.PutObjectOptions{
		ContentType: mime.String(),
	})
	if err != nil {
		panic(errors.Wrap(err, "could write image to s3 storage."))
	}
	filePath := base.Echo.Reverse("image.getImage", fileIndex)
	return c.JSON(http.StatusCreated, response.CreateImageResponse{
		Message: "SUCCESS",
		Error:   nil,
		Data: struct {
			FilePath *string `json:"filename"`
		}{
			&filePath,
		},
	})
}
