// Copyright 2022 Vitalii Noha vitalii.noga@gmail.com. All rights reserved.

package dl

import (
	"fmt"
)

var attrs = struct {
	kindFmt string
	itemFmt string
	depFmt  string
}{
	"%s\n",
	"%s" + ItemOptCode + "\n",
	"%s" + ItemSeparator + "%s\n",
}

func (m *module) Kind() string {
	return m.kind
}

func (m *module) Items() Items {
	return m.items
}

func (m *module) AddItem(item string) error {
	if _, found := m.items[item]; found {
		return fmt.Errorf(ItemExistsF, item)
	}
	m.items[item] = Item{}
	return nil
}

func (m *module) AddDependency(item, dependency, resolver string, update bool) error {
	curr, found := m.items[item]
	if !found {
		return fmt.Errorf(ItemIsMissingF, item)
	}
	if _, found := curr[dependency]; found && !update {
		return fmt.Errorf(DepItemExistsF, dependency, item)
	}
	curr[dependency] = resolver
	return nil
}

func (m *module) DeleteItem(item string) error {
	delete(m.items, item)
	return nil
}

func (m *module) DeleteDependency(item, dependency string) error {
	if curr, found := m.items[item]; found {
		delete(curr, dependency)
	}
	return nil
}

func (m *module) Dependency(item, dep string) string {
	if deps := m.items[item]; deps != nil {
		if res, found := deps[dep]; found {
			return fmt.Sprintf(attrs.depFmt, dep, res)
		}
	}
	return ""
}
