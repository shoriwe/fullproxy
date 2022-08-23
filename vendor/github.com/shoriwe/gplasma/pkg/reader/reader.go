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
	Line() int
	Char() rune
}

type StringReader struct {
	content []rune
	index   int
	line    int
	length  int
}

func (s *StringReader) Line() int {
	return s.line
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
	character := s.content[s.index]
	if character == '\n' {
		s.line++
	}
	return character
}

func NewStringReader(code string) Reader {
	runeCode := []rune(code)
	return &StringReader{
		content: runeCode,
		line:    1,
		index:   0,
		length:  len(runeCode),
	}
}

func NewStringReaderFromFile(file io.ReadCloser) Reader {
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
