package IOTools

import (
	"io/ioutil"
	"strings"
)

func ReadLines(filePath string) ([]string, error) {
	content, readingError := ioutil.ReadFile(filePath)
	if readingError != nil {
		return nil, readingError
	}
	return strings.Split(string(content), "\n"), nil
}