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
	if err := r.checkFile(); err != nil {
		return nil, err
	}
	mod := module{}
	mod.name = r.name
	mod.items = Items{}
	item := ""
	line := ""
	token1 := ""
	token2 := ""
	var err error
	var ioErr error
	for {
		if line == "" {
			line, ioErr = r.buf.ReadString('\n')
			if ioErr != nil && ioErr != io.EOF {
				return nil, ioErr
			}
		}
		line = r.removeComment(line)
		line = strings.Trim(line, trimChars)
		if line != "" {
			// process the line
			token1, token2 = r.splitLine(line)
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
				mod.items[item] = append(mod.items[item], Item{{token1, token2}}...)
			}
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

func (r *reader) readAsync(res chan<- moduleAsync) {
	m, e := r.read()
	res <- moduleAsync{m, e}
}

func (r *reader) readItem(mod *module, name string) (next string, err error) {
	next = ""
	token1 := ""
	token2 := ""
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
		next = r.removeComment(next)
		// process the line
		next = strings.Trim(next, trimChars)
		if next != "" {
			token1, token2 = r.splitLine(next)
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
				mod.items[name] = append(mod.items[name], Item{{token1, token2}}...)
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
				mod.items[name] = append(mod.items[name], Item{{token1, token2}}...)
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

func (r *reader) checkFile() error {
	line := ""
	token1 := ""
	token2 := ""
	var ioErr error
	for {
		line, ioErr = r.buf.ReadString('\n')
		if ioErr != nil && ioErr != io.EOF {
			return ioErr
		}
		line = r.removeComment(line)
		line = strings.Trim(line, trimChars)
		if line != "" {
			token1, token2 = r.splitLine(line)
			// check the current file
			if token1 != AppCode || token2 != "" {
				return fmt.Errorf(FirstTokenInvalid)
			} else {
				return nil
			}
		}
		// check the EOF
		if ioErr != nil {
			break
		}
		line = ""
	}
	return fmt.Errorf(FirstTokenIsMissing)
}

func (r *reader) splitLine(line string) (t1 string, t2 string) {
	pos := strings.Index(line, ItemSeparator)
	if pos > 0 {
		t1 = strings.Trim(line[0:pos], trimChars)
		t2 = strings.Trim(line[pos:], trimChars)
	} else {
		t1 = line
		t2 = ""
	}
	return
}

func (r *reader) removeComment(line string) string {
	index := strings.Index(line, CommentOptCode)
	if index > 0 {
		return line[:index]
	} else {
		return line
	}
}
