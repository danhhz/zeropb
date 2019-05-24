// Copyright 2019 Daniel Harrison. All Rights Reserved.

package generate

import (
	fmt "fmt"
	"io"
	"strconv"
	"strings"
)

// StringTreeIndent represents the number of stops to indent.
type StringTreeIndent int

// Next returns an indent that's stops step further.
func (i StringTreeIndent) Next() StringTreeIndent {
	return StringTreeIndent(int(i) + 1)
}

// StringTree is a wrapper around io.Writer that removes boilerplate.
type StringTree struct {
	w      io.Writer
	indent string
	err    error
}

// NewStringTree returns a StringTree with the request string used for each step
// of indentation.
func NewStringTree(w io.Writer, indent string) *StringTree {
	return &StringTree{w: w, indent: indent}
}

// Write passes along the arguments to the wrapped io.Writer, type sniffing and
// stringifying as necessary.
func (st *StringTree) Write(strs ...interface{}) {
	for _, str := range strs {
		if st.err != nil {
			return
		}
		switch x := str.(type) {
		case StringTreeIndent:
			_, st.err = st.w.Write([]byte(strings.Repeat(st.indent, int(x))))
		case []byte:
			_, st.err = st.w.Write(x)
		case string:
			_, st.err = st.w.Write([]byte(x))
		case fmt.Stringer:
			s := x.String()
			_, st.err = st.w.Write([]byte(s))
		case int32:
			_, st.err = st.w.Write([]byte(strconv.Itoa(int(x))))
		default:
			panic(fmt.Sprintf(`unknown %T: %v`, str, str))
		}
	}
}

// Error returns the first error encountered, if any.
func (st *StringTree) Error() error {
	return st.err
}
