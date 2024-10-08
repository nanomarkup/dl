// Copyright 2022 Vitalii Noha vitalii.noga@gmail.com. All rights reserved.

package dl

import (
	"os"

	"gopkg.in/check.v1"
)

func (s *lodSuite) TestLoadEmptyModule(c *check.C) {
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	os.Chdir("test/empty")

	file, err := os.Open(getModuleFileName("empty_test.sb", "sb"))
	c.Assert(err, check.IsNil)
	defer file.Close()
	items, err := newReader(file).read()
	c.Assert(err, check.IsNil)
	c.Assert(items, check.NotNil)
}

func (s *lodSuite) TestLoadInitializedStruct(c *check.C) {
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	os.Chdir("test/init")

	file, err := os.Open(getModuleFileName("init_test.sb", "sb"))
	c.Assert(err, check.IsNil)
	defer file.Close()
	items, err := newReader(file).read()
	c.Assert(err, check.IsNil)
	c.Assert(items, check.NotNil)
}
