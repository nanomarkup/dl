// Copyright 2022 Vitalii Noha vitalii.noga@gmail.com. All rights reserved.

package lod

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"gopkg.in/check.v1"
)

func (s *lodSuite) TestReadAll(c *check.C) {
	m := Manager{}
	m.SetLogger(hclog.New(&hclog.LoggerOptions{
		Name:   "test",
		Level:  hclog.Trace,
		Output: os.Stdout,
	}))
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	os.Chdir("sb/app")
	fmt.Println(os.Getwd())
	r, e := m.ReadAll("sb")
	if e != nil {
		fmt.Println(e.Error())
		c.Error()
		return
	}
	c.Assert(r, check.NotNil)
}
