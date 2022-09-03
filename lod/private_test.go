// Copyright 2022 Vitalii Noha vitalii.noga@gmail.com. All rights reserved.

package lod

import (
	"os"

	"gopkg.in/check.v1"
)

func (s *lodSuite) TestLoadEmptyModule(c *check.C) {
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	os.Chdir("sb/empty")
	items, err := loadModule("empty_test.sb", "sb")
	c.Assert(err, check.IsNil)
	c.Assert(items, check.NotNil)
}
