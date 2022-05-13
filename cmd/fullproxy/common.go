package main

type stringSlice struct {
	contents []string
}

func (s *stringSlice) String() string {
	return "{}"
}

func (s *stringSlice) Set(ss string) error {
	s.contents = append(s.contents, ss)
	return nil
}
