package utils

import (
	"github.com/leoleoasd/EduOJBackend/base"
	"github.com/leoleoasd/EduOJBackend/base/config"
	"github.com/pkg/errors"
)

func MustCreateBucket(name string) {
	found, err := base.Storage.BucketExists(name)
	if err != nil {
		panic(errors.Wrap(err, "could not query if bucket exists"))
	}
	if !found {
		err = base.Storage.MakeBucket(name, config.MustGet("storage.region", "us-east-1").String())
		if err != nil {
			panic(errors.Wrap(err, "could not create bucket"))
		}
	}
}

func MustCreateBuckets(names ...string) {
	for _, name := range names {
		MustCreateBucket(name)
	}
}
