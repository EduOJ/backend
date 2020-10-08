package utils

import (
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
