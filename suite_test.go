// Copyright 2022 Vitalii Noha vitalii.noga@gmail.com. All rights reserved.

package dl

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
		items: Items{
			appName: [][]string{},
		},
	}})

const (
	appName string = "test"
)
