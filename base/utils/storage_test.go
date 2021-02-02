package utils

import (
	"bytes"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/minio/minio-go"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
)

func getPresignedURLContent(t *testing.T, presignedUrl string) (content string) {
	resp, err := http.Get(presignedUrl)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	length, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	assert.Nil(t, err)
	body := make([]byte, length)
	_, err = resp.Body.Read(body)
	return string(body)
}

func TestMustCreateBucket(t *testing.T) {
	t.Run("testMustCreateBucketExistingBucket", func(t *testing.T) {
		assert.Nil(t, base.Storage.MakeBucket("existing-bucket", config.MustGet("storage.region", "us-east-1").String()))
		found, err := base.Storage.BucketExists("existing-bucket")
		assert.True(t, found)
		assert.Nil(t, err)
		assert.Nil(t, CreateBucket("existing-bucket"))
		found, err = base.Storage.BucketExists("existing-bucket")
		assert.True(t, found)
		assert.Nil(t, err)
	})
	t.Run("testMustCreateBucketNonExistingBucket", func(t *testing.T) {
		found, err := base.Storage.BucketExists("non-existing-bucket")
		assert.False(t, found)
		assert.Nil(t, err)
		assert.Nil(t, CreateBucket("non-existing-bucket"))
		found, err = base.Storage.BucketExists("non-existing-bucket")
		assert.True(t, found)
		assert.Nil(t, err)
	})
}

func TestGetPresignedURL(t *testing.T) {
	b := []byte("test_get_presigned_url")
	reader := bytes.NewReader(b)
	n, err := base.Storage.PutObject("test-bucket", "test_get_presigned_url_object", reader, int64(len(b)), minio.PutObjectOptions{})
	assert.Equal(t, int64(len(b)), n)
	assert.Nil(t, err)
	presignedUrl, err := GetPresignedURL("test-bucket", "test_get_presigned_url_object", "test_get_presigned_url_file_name")
	assert.Nil(t, err)
	resp, err := http.Get(presignedUrl)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, string(b), getPresignedURLContent(t, presignedUrl))
}
