package utils

import (
	"context"
	"fmt"
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/minio/minio-go/v7"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"net/url"
	"sync"
	"time"
)

var bucketsCreated sync.Map

func CreateBucket(name string) error {
	_, ok := bucketsCreated.Load(name)
	if ok {
		return nil
	}
	found, err := base.Storage.BucketExists(context.Background(), name)
	if err != nil {
		return errors.Wrap(err, "could not query if bucket exists")
	}
	if !found {
		err = base.Storage.MakeBucket(context.Background(), name, minio.MakeBucketOptions{
			Region: viper.GetString("storage.region"),
		})
		if err != nil {
			return errors.Wrap(err, "could not create bucket")
		}
	}
	bucketsCreated.Store(name, true)
	return nil
}

func GetPresignedURL(bucket string, path string, fileName string) (string, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf(`inline; filename="%s"`, fileName))
	presignedURL, err := base.Storage.PresignedGetObject(context.Background(), bucket, path, time.Second*60, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), err
}
