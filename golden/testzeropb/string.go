// Copyright 2019 Daniel Harrison. All Rights Reserved.

package testzeropb

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

func (m TestMessage) String() string {
	var buf strings.Builder
	m.Textual(&buf)
	return buf.String()
}

// Textual returns the contents in the gogo string format.
func (m TestMessage) Textual(w io.Writer) {
	if _, ok := m.offsets.Get(5); ok {
		fmt.Fprintf(w, `uint64=%d `, m.Uint64())
	}
	if _, ok := m.offsets.Get(15); ok {
		fmt.Fprintf(w, `byte_array=%s `, strconv.Quote(string(m.ByteArray())))
	}
	if _, ok := m.offsets.Get(16); ok {
		fmt.Fprintf(w, `enum=%d `, m.Enum())
	}
	var sub TestMessage
	if ok, err := m.Message(&sub); err != nil {
		panic(err)
	} else if ok {
		io.WriteString(w, `message=<`)
		sub.Textual(w)
		io.WriteString(w, `> `)
	}
	for it := m.Messages(); ; {
		if ok, err := it.Next(&sub); err != nil {
			panic(err)
		} else if !ok {
			break
		}
		io.WriteString(w, `messages=<`)
		sub.Textual(w)
		io.WriteString(w, `> `)
	}
}
