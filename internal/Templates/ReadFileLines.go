package Templates

import (
	"io/ioutil"
	"strings"
)

func ReadFileLines(filePath string) ([]string, error) {
	content, readingError := ioutil.ReadFile(filePath)
	if readingError != nil {
		return nil, readingError
	}
	stringContent := string(content)
	return strings.Split(stringContent, "\n"), nil
}
