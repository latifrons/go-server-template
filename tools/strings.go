package tools

import (
	"errors"
	"io/ioutil"
	"strings"
)

func ConfigListToMap(s string) map[string]string {
	v := map[string]string{}
	slots := strings.Split(s, ";")
	for _, slot := range slots {
		slotElems := strings.Split(slot, ":")
		if len(slotElems) != 2 {
			panic("bad config")
		}
		v[slotElems[0]] = slotElems[1]
	}
	return v
}

func ReadStringArrayFromFile(filePath string) (ss []string, err error) {
	var content []byte
	content, err = ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	ss = strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")
	return
}

func ParseBool(s string) (b bool, err error) {
	s = strings.ToLower(s)
	if s == "true" || s == "1" {
		b = true
	} else if s == "false" || s == "0" {
		b = false
	} else {
		err = errors.New("parse failed")
	}
	return
}

func ContainsIgnoreCase(s string, subStr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(subStr))
}
