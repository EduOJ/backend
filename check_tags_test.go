// +build !test

package main_test

import "testing"

func TestIsTest(t *testing.T) {
	t.Log("You should run tests with `go test -tags test` !")
	t.FailNow()
}
