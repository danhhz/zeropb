// Copyright 2019 Daniel Harrison. All Rights Reserved.

package testzeropb

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Text returns the contents in the proto text format.
func (m TestMessage) Text() string {
	var buf strings.Builder
	m.WriteText(&buf)
	return buf.String()
}

// WriteText writes the contents in the proto text format.
func (m TestMessage) WriteText(w io.Writer) {
	// TODO(dan): A bunch of these aren't right, but this is currently only used
	// while debugging test failures. Revisit if we start codegen'ing this.
	if offset := m.offsets[1]; offset != 0 {
		fmt.Fprintf(w, `bool=%t `, m.Bool())
	}
	if offset := m.offsets[2]; offset != 0 {
		fmt.Fprintf(w, `int32=%d `, m.Int32())
	}
	if offset := m.offsets[3]; offset != 0 {
		fmt.Fprintf(w, `int64=%d `, m.Int64())
	}
	if offset := m.offsets[4]; offset != 0 {
		fmt.Fprintf(w, `uint32=%d `, m.Uint32())
	}
	if offset := m.offsets[5]; offset != 0 {
		fmt.Fprintf(w, `uint64=%d `, m.Uint64())
	}
	if offset := m.offsets[6]; offset != 0 {
		fmt.Fprintf(w, `sint32=%d `, m.Sint32())
	}
	if offset := m.offsets[7]; offset != 0 {
		fmt.Fprintf(w, `sint64=%d `, m.Sint64())
	}
	if offset := m.offsets[8]; offset != 0 {
		fmt.Fprintf(w, `fixed32=%d `, m.Fixed32())
	}
	if offset := m.offsets[9]; offset != 0 {
		fmt.Fprintf(w, `fixed64=%d `, m.Fixed64())
	}
	if offset := m.offsets[10]; offset != 0 {
		fmt.Fprintf(w, `sfixed32=%d `, m.Sfixed32())
	}
	if offset := m.offsets[11]; offset != 0 {
		fmt.Fprintf(w, `sfixed64=%d `, m.Sfixed64())
	}
	if offset := m.offsets[12]; offset != 0 {
		fmt.Fprintf(w, `double=%v `, m.Double())
	}
	if offset := m.offsets[13]; offset != 0 {
		fmt.Fprintf(w, `float=%v `, m.Float())
	}
	if offset := m.offsets[14]; offset != 0 {
		fmt.Fprintf(w, `string=%s `, m.String())
	}
	if offset := m.offsets[15]; offset != 0 {
		fmt.Fprintf(w, `byte_array=%s `, strconv.Quote(string(m.ByteArray())))
	}
	if offset := m.offsets[16]; offset != 0 {
		fmt.Fprintf(w, `enum=%d `, m.Enum())
	}
	var sub TestMessage
	if ok, err := m.Message(&sub); err != nil {
		panic(err)
	} else if ok {
		io.WriteString(w, `message=<`)
		sub.WriteText(w)
		io.WriteString(w, `> `)
	}
	for it := m.Messages(); ; {
		if ok, err := it.Next(&sub); err != nil {
			panic(err)
		} else if !ok {
			break
		}
		io.WriteString(w, `messages=<`)
		sub.WriteText(w)
		io.WriteString(w, `> `)
	}
}
