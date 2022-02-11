package tools

import (
	"io/ioutil"
	"strings"
)

func ReadStringArrayFromFile(filePath string) (ss []string, err error) {
	var content []byte
	content, err = ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	ss = strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")
	return
}
