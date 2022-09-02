// Copyright 2022 Vitalii Noha vitalii.noga@gmail.com. All rights reserved.

package lod

func GetModuleFileName(name string) string {
	return getModuleFileName(name)
}

func IsModuleExists(name string) bool {
	return isModuleExists(name)
}

func IsItemExists(kind, item string) (bool, string) {
	return isItemExists(kind, item)
}
