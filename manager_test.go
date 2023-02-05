// Copyright 2022 Vitalii Noha vitalii.noga@gmail.com. All rights reserved.

package dl

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"gopkg.in/check.v1"
)

func (s *lodSuite) TestReadUseFileName(c *check.C) {
	m := Manager{}
	m.SetLogger(hclog.New(&hclog.LoggerOptions{
		Name:   "test",
		Level:  hclog.Trace,
		Output: os.Stdout,
	}))
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	os.Chdir("test/app")
	r, e := m.Read("app_test.sb")
	if e != nil {
		fmt.Println(e.Error())
		c.Error()
		return
	}
	c.Assert(r, check.NotNil)
}

func (s *lodSuite) TestReadUseFilePath(c *check.C) {
	m := Manager{}
	m.SetLogger(hclog.New(&hclog.LoggerOptions{
		Name:   "test",
		Level:  hclog.Trace,
		Output: os.Stdout,
	}))
	wd, _ := os.Getwd()
	path := filepath.Join(wd, "test", "app", "app_test.sb")
	r, e := m.Read(path)
	if e != nil {
		fmt.Println(e.Error())
		c.Error()
		return
	}
	c.Assert(r, check.NotNil)
}

func (s *lodSuite) TestReadAll(c *check.C) {
	m := Manager{}
	m.Kind = kindName
	m.SetLogger(hclog.New(&hclog.LoggerOptions{
		Name:   "test",
		Level:  hclog.Trace,
		Output: os.Stdout,
	}))
	wd, _ := os.Getwd()
	defer os.Chdir(wd)
	os.Chdir("test/app")
	r, e := m.ReadAll()
	if e != nil {
		fmt.Println(e.Error())
		c.Error()
		return
	}
	c.Assert(r, check.NotNil)
}
