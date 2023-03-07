// Copyright 2022 Vitalii Noha vitalii.noga@gmail.com. All rights reserved.

package dl

import (
	"bytes"
	"fmt"
	"sort"
)

func (f *Formatter) Item(name string, deps [][]string) string {
	if deps == nil {
		return ""
	}
	var res bytes.Buffer
	res.WriteString(fmt.Sprintf(attrs.itemFmt, name))
	var d, r string
	for _, row := range deps {
		d = row[0]
		if len(row) > 1 {
			r = row[1]
		} else {
			r = ""
		}
		res.WriteString(fmt.Sprintf("\t"+attrs.depFmt, d, r))
	}
	return res.String()
}

func (f *Formatter) String(module Module) string {
	var res bytes.Buffer
	res.WriteString(fmt.Sprintf(attrs.typeFmt, AppCode))
	// sort items
	items := module.Items()
	sorted := make([]string, 0, len(items))
	for item := range items {
		sorted = append(sorted, item)
	}
	sort.Strings(sorted)
	// add items
	for _, item := range sorted {
		res.WriteString("\n" + f.Item(item, items[item]))
	}
	return res.String()
}
