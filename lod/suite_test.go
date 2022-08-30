package lod

import (
	"io/ioutil"
	"testing"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type SModuleSuite struct {
}

var _ = check.Suite(&SModuleSuite{})

func copyFile(src, dst string) error {
	input, err := ioutil.ReadFile(src)
	if err == nil {
		return ioutil.WriteFile(dst, input, 0644)
	} else {
		return err
	}
}
