// Copyright 2022 Vitalii Noha vitalii.noga@gmail.com. All rights reserved.

package lod

import (
	"testing"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type lodSuite struct {
	mod module
}

var _ = check.Suite(&lodSuite{
	mod: module{
		kind: kindName,
		items: Items{
			appName: map[string]string{},
		},
	}})

const (
	kindName string = "sb"
	appName  string = "test"
)
