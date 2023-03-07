// Copyright 2022 Vitalii Noha vitalii.noga@gmail.com. All rights reserved.

package dl

import (
	"fmt"
)

var attrs = struct {
	typeFmt string
	itemFmt string
	depFmt  string
}{
	"%s\n",
	"%s" + ItemOptCode + "\n",
	"%s" + ItemSeparator + "%s\n",
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
	for i, row := range curr {
		if row[0] == dependency {
			if update {
				curr[i][1] = resolver
				return nil
			} else {
				return fmt.Errorf(DepItemExistsF, dependency, item)
			}
		}
	}
	m.items[item] = append(m.items[item], [][]string{{dependency, resolver}}...)
	return nil
}

func (m *module) DeleteItem(item string) error {
	delete(m.items, item)
	return nil
}

func (m *module) DeleteDependency(item, dependency string) error {
	if curr, found := m.items[item]; found {
		for i, row := range curr {
			if row[0] == dependency {
				curr = append(curr[:i], curr[i+1:]...)
				m.items[item] = curr
				break
			}
		}
	}
	return nil
}

func (m *module) Dependency(item, dep string) string {
	if deps := m.items[item]; deps != nil {
		for _, row := range deps {
			if row[0] == dep {
				if len(row) > 1 {
					return fmt.Sprintf(attrs.depFmt, dep, row[1])
				} else {
					return fmt.Sprintf(attrs.depFmt, dep, "")
				}
			}
		}
	}
	return ""
}
