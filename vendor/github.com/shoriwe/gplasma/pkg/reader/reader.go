package reader

import (
	"bytes"
	"io"
)

type Reader interface {
	Next()
	Redo()
	HasNext() bool
	Index() int
	Char() rune
}

type StringReader struct {
	content []rune
	index   int
	length  int
}

func (s *StringReader) Next() {
	s.index++
}

func (s *StringReader) Redo() {
	s.index--
}

func (s *StringReader) HasNext() bool {
	return s.index < s.length
}

func (s *StringReader) Index() int {
	return s.index
}

func (s *StringReader) Char() rune {
	return s.content[s.index]
}

func NewStringReader(code string) *StringReader {
	runeCode := []rune(code)
	return &StringReader{
		content: runeCode,
		index:   0,
		length:  len(runeCode),
	}
}

func NewStringReaderFromFile(file io.ReadCloser) *StringReader {
	defer file.Close()
	content, readingError := io.ReadAll(file)
	if readingError != nil {
		panic(readingError)
	}
	if bytes.Equal(content[:3], []byte{0xef, 0xbb, 0xbf}) {
		content = content[3:]
	}
	content = bytes.ReplaceAll(content, []byte{'\r', '\n'}, []byte{'\n'})
	return NewStringReader(string(content))
}
