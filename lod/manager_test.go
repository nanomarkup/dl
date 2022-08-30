package lod

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"gopkg.in/check.v1"
)

func (s *SModuleSuite) TestLoading(c *check.C) {
	m := Manager{}
	m.SetLogger(hclog.New(&hclog.LoggerOptions{
		Name:   "test",
		Level:  hclog.Level(1),
		Output: os.Stdout,
	}))
	r, e := m.ReadAll("sb")
	if e != nil {
		fmt.Println(e.Error())
		c.Error()
		return
	}
	c.Assert(r, check.NotNil)
}
