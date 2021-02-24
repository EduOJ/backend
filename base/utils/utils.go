package utils

import (
	"math/rand"
	"reflect"
	"unsafe"
)

// Random string generator by https://stackoverflow.com/a/22892986/8031146

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandStr generates a string with given length.
// Only contains a-z and A-Z.
func RandStr(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func Contain(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

// IdUniqueInA returns IDs that in a but not in b
// Reference: https://stackoverflow.com/questions/52120488/what-is-the-most-efficient-way-to-get-the-intersection-and-exclusions-from-two-a
func IdUniqueInA(a []uint, b []uint) (uniqueA []uint) {
	m := make(map[uint]uint8)
	for _, v := range a {
		m[v] |= 1 << 0
	}
	for _, v := range b {
		m[v] |= 1 << 1
	}
	for v, sel := range m {
		if sel == 1<<0 {
			uniqueA = append(uniqueA, v)
		}
	}
	return
}
