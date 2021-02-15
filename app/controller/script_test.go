package controller_test

import (
	"bytes"
	"context"
	"encoding/xml"
	"github.com/leoleoasd/EduOJBackend/app/response"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/database/models"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestGetScript(t *testing.T) {
	script := models.Script{
		Name:     "test_get_script",
		Filename: "test_get_script.zip",
	}
	assert.NoError(t, base.DB.Create(&script).Error)
	file := newFileContent("test_get_script", "test_get_script.zip", b64Encode("test_get_script_zip_content"))
	_, err := base.Storage.PutObject(context.Background(), "scripts", script.Name, file.reader, file.size, minio.PutObjectOptions{})
	assert.NoError(t, err)
	_, err = file.reader.Seek(0, io.SeekStart)
	assert.NoError(t, err)

	noFileScript := models.Script{
		Name:     "test_no_file_script",
		Filename: "test_no_file_script",
	}
	assert.NoError(t, base.DB.Create(&noFileScript).Error)

	t.Run("Success", func(t *testing.T) {
		req := makeReq(t, "GET", base.Echo.Reverse("judger.getScript", script.Name), "", judgerAuthorize)
		resp := makeResp(req)
		assert.Equal(t, http.StatusFound, resp.StatusCode)
		content := getPresignedURLContent(t, resp.Header.Get("Location"))
		assert.Equal(t, "test_get_script_zip_content", content)
	})

	t.Run("MissingName", func(t *testing.T) {
		req := makeReq(t, "GET", base.Echo.Reverse("judger.getScript"), "", judgerAuthorize)
		resp := makeResp(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		jsonEQ(t, response.ErrorResp("NOT_FOUND", nil), resp)
	})

	t.Run("MissingFile", func(t *testing.T) {
		req := makeReq(t, "GET", base.Echo.Reverse("judger.getScript", noFileScript.Name), "", judgerAuthorize)
		resp := makeResp(req)
		assert.Equal(t, http.StatusFound, resp.StatusCode)
		assert.NotEmpty(t, resp.Header.Get("Location"))
		resp, err := http.Get(resp.Header.Get("Location"))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		bodyBuf := bytes.Buffer{}
		_, err = bodyBuf.ReadFrom(resp.Body)
		assert.NoError(t, err)
		var xmlresp struct {
			Name     xml.Name `xml:"Error"`
			Code     string   `xml:"Code"`
			Resource string   `xml:"Resource"`
		}
		assert.NoError(t, xml.Unmarshal(bodyBuf.Bytes(), &xmlresp))
		assert.Equal(t, struct {
			Name     xml.Name `xml:"Error"`
			Code     string   `xml:"Code"`
			Resource string   `xml:"Resource"`
		}{
			Code:     "NoSuchKey",
			Resource: "test_no_file_script",
		}, xmlresp)
	})
}
