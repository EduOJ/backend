package utils

import (
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/pkg/errors"
	"sync"
)

var bucketsCreated sync.Map

func CreateBucket(name string) error {
	_, ok := bucketsCreated.Load(name)
	if ok {
		return nil
	}
	found, err := base.Storage.BucketExists(name)
	if err != nil {
		return errors.Wrap(err, "could not query if bucket exists")
	}
	if !found {
		err = base.Storage.MakeBucket(name, config.MustGet("storage.region", "us-east-1").String())
		if err != nil {
			return errors.Wrap(err, "could not create bucket")
		}
	}
	bucketsCreated.Store(name, true)
	return nil
}
