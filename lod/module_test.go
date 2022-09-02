// Copyright 2022 Vitalii Noha vitalii.noga@gmail.com. All rights reserved.

package lod

import (
	"fmt"

	"gopkg.in/check.v1"
)

func (s *lodSuite) TestApp(c *check.C) {
	item, err := s.mod.App(appName)
	c.Assert(err, check.IsNil)
	c.Assert(item, check.NotNil)
}

func (s *lodSuite) TestApps(c *check.C) {
	items, err := s.mod.Apps()
	c.Assert(err, check.IsNil)
	c.Assert(items, check.NotNil)
	c.Assert(len(items), check.Equals, 1)
}

func (s *lodSuite) TestKind(c *check.C) {
	c.Assert(s.mod.Kind(), check.Equals, kindName)
}

func (s *lodSuite) TestItems(c *check.C) {
	items := s.mod.Items()
	c.Assert(items, check.NotNil)
	c.Assert(len(items) > 0, check.Equals, true)
}

func (s *lodSuite) TestDependency(c *check.C) {
	itemName := "dep-item"
	depName := "dependency"
	resName := "resolver"
	s.mod.AddItem(itemName)
	s.mod.AddDependency(itemName, depName, resName, false)
	c.Assert(s.mod.Dependency(itemName, depName), check.Equals, fmt.Sprintf(attrs.depFmt, depName, resName))
}

func (s *lodSuite) TestAddItem(c *check.C) {
	itemName := "add-item"
	s.mod.AddItem(itemName)
	c.Assert(s.mod.Items()[itemName], check.NotNil)
}

func (s *lodSuite) TestDeleteItem(c *check.C) {
	itemName := "delete-item"
	s.mod.AddItem(itemName)
	l := len(s.mod.Items())
	s.mod.DeleteItem(itemName)
	c.Assert(len(s.mod.Items()), check.Equals, l-1)
}

func (s *lodSuite) TestAddDependency(c *check.C) {
	itemName := "new-dep-item"
	depName := "new-dependency"
	resName := "new-resolver"
	s.mod.AddItem(itemName)
	s.mod.AddDependency(itemName, depName, resName, false)
	items := s.mod.Items()
	c.Assert(items, check.NotNil)
	item := items[itemName]
	c.Assert(item, check.NotNil)
	res, found := item[depName]
	c.Assert(found, check.Equals, true)
	c.Assert(res, check.Equals, resName)
	// test with the update flag is false
	s.mod.AddDependency(itemName, depName, "new-resolver-2", false)
	res = item[depName]
	c.Assert(res, check.Equals, resName)
	// test with the update flag is true
	s.mod.AddDependency(itemName, depName, "new-resolver-2", true)
	res = item[depName]
	c.Assert(res, check.Equals, "new-resolver-2")
}

func (s *lodSuite) TestAddDeleteDependency(c *check.C) {
	itemName := "del-dep-item"
	depName := "del-dependency"
	s.mod.AddItem(itemName)
	s.mod.AddDependency(itemName, depName, "new-resolver", false)
	item := s.mod.items[itemName]
	l := len(item)
	s.mod.DeleteDependency(itemName, depName)
	c.Assert(len(item), check.Equals, l-1)
	res, found := item[depName]
	c.Assert(res, check.Equals, "")
	c.Assert(found, check.Equals, false)
}
