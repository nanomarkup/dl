// Copyright 2023 Vitalii Noha vitalii.noga@gmail.com. All rights reserved.

package dl

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func newReader(file *os.File) *reader {
	return &reader{file.Name(), bufio.NewReader(file)}
}

func (r *reader) read() (*module, error) {
	mod := module{}
	mod.name = r.name
	mod.items = Items{}
	item := ""
	line := ""
	token1 := ""
	token2 := ""
	pos := 0
	index := 1
	cindex := 0
	trimChars := " \t\n\r"
	var err error
	var ioErr error
	for {
		if line == "" {
			line, ioErr = r.buf.ReadString('\n')
			if ioErr != nil && ioErr != io.EOF {
				return nil, ioErr
			}
		}
		// remove comment
		cindex = strings.Index(line, CommentOptCode)
		if cindex > 0 {
			line = line[:cindex]
		}
		// process the line
		line = strings.Trim(line, trimChars)
		if line != "" {
			pos = strings.Index(line, ItemSeparator)
			if pos > 0 {
				token1 = strings.Trim(line[0:pos], trimChars)
				token2 = strings.Trim(line[pos:], trimChars)
			} else {
				token1 = line
				token2 = ""
			}
			if index == 1 {
				// check and initialize a kind of module
				// if token1 != kind || token2 != "" {
				// 	return nil, fmt.Errorf(FirstTokenInvalidF, kind)
				// }
				mod.kind = token1
			} else {
				if token1[len(token1)-1:] == ItemOptCode {
					// new item found, process it
					if token2 != "" {
						return nil, fmt.Errorf(LineSyntaxInvalidF, line)
					}
					item = token1[:len(token1)-1]
					line, err = r.readItem(&mod, item)
					if err != nil {
						if err == io.EOF {
							ioErr = err
						} else {
							return nil, err
						}
					}
					if line != "" && ioErr == nil {
						continue
					}
				} else {
					// add new dependency item
					if mod.items[item] == nil {
						mod.items[item] = Item{}
					}
					mod.items[item][token1] = token2
				}
			}
			index++
		}
		// check the EOF
		if ioErr != nil {
			break
		}
		line = ""
	}
	// add empty item if it exists
	if mod.items[item] == nil {
		mod.items[item] = Item{}
	}
	return &mod, nil
}

func (r *reader) readItem(mod *module, name string) (next string, err error) {
	next = ""
	token1 := ""
	token2 := ""
	pos := 0
	cindex := 0
	trimChars := " \t\n\r"
	groupItem := ""
	groupName := ""
	var ioErr error
	for {
		if next == "" {
			next, ioErr = r.buf.ReadString('\n')
			if ioErr != nil && ioErr != io.EOF {
				return "", ioErr
			}
		}
		// remove comment
		cindex = strings.Index(next, CommentOptCode)
		if cindex > 0 {
			next = next[:cindex]
		}
		// process the line
		next = strings.Trim(next, trimChars)
		if next != "" {
			pos = strings.Index(next, ItemSeparator)
			if pos > 0 {
				token1 = strings.Trim(next[0:pos], trimChars)
				token2 = strings.Trim(next[pos:], trimChars)
			} else {
				token1 = next
				token2 = ""
			}
			if token1[len(token1)-1:] == ItemOptCode {
				// do not process the next item
				return next, nil
			} else if token1 == InitEndOptCode || token2 == InitEndOptCode {
				// the current group ended, check syntax and do nothing
				if token2 != "" {
					return "", fmt.Errorf(LineSyntaxInvalidF, next)
				}
				return "", nil
			} else if token2 != "" && token2[len(token2)-1:] == InitBegOptCode {
				// add a new group item and process it
				groupItem = strings.Trim(token2[:len(token2)-1], trimChars)
				groupName = getNextGroupName()
				token2 = fmt.Sprintf("[%s]%s", groupName, groupItem)
				// remove the pointer from the item
				if strings.HasPrefix(groupItem, "*") {
					groupItem = fmt.Sprintf("[%s]%s", groupName, groupItem[1:])
				} else {
					groupItem = token2
				}
				// add new dependency to the existing item
				if mod.items[name] == nil {
					mod.items[name] = Item{}
				}
				mod.items[name][token1] = token2
				// read the group item
				next, err = r.readItem(mod, groupItem)
				if err != nil {
					return "", err
				}
				if next != "" {
					continue
				}
			} else {
				// add new dependency item
				if mod.items[name] == nil {
					mod.items[name] = Item{}
				}
				mod.items[name][token1] = token2
			}
		}
		next = ""
		// check the EOF
		if ioErr != nil {
			break
		}
	}
	return
}
