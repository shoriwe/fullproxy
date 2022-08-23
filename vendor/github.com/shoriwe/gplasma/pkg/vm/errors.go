package vm

import "fmt"

var (
	NotOperable   = fmt.Errorf("not operable")
	NotIndexable  = fmt.Errorf("not indexable")
	NotComparable = fmt.Errorf("not comparable")
)
