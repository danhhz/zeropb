// Copyright 2019 Daniel Harrison. All Rights Reserved.

package testzeropb

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/danhhz/zeropb"
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
	offsets := zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m)
	if _, ok := offsets.Get(1); ok {
		fmt.Fprintf(w, `bool=%t `, m.Bool())
	}
	if _, ok := offsets.Get(2); ok {
		fmt.Fprintf(w, `int32=%d `, m.Int32())
	}
	if _, ok := offsets.Get(3); ok {
		fmt.Fprintf(w, `int64=%d `, m.Int64())
	}
	if _, ok := offsets.Get(4); ok {
		fmt.Fprintf(w, `uint32=%d `, m.Uint32())
	}
	if _, ok := offsets.Get(5); ok {
		fmt.Fprintf(w, `uint64=%d `, m.Uint64())
	}
	if _, ok := offsets.Get(6); ok {
		fmt.Fprintf(w, `sint32=%d `, m.Sint32())
	}
	if _, ok := offsets.Get(7); ok {
		fmt.Fprintf(w, `sint64=%d `, m.Sint64())
	}
	if _, ok := offsets.Get(8); ok {
		fmt.Fprintf(w, `fixed32=%d `, m.Fixed32())
	}
	if _, ok := offsets.Get(9); ok {
		fmt.Fprintf(w, `fixed64=%d `, m.Fixed64())
	}
	if _, ok := offsets.Get(10); ok {
		fmt.Fprintf(w, `sfixed32=%d `, m.Sfixed32())
	}
	if _, ok := offsets.Get(11); ok {
		fmt.Fprintf(w, `sfixed64=%d `, m.Sfixed64())
	}
	if _, ok := offsets.Get(12); ok {
		fmt.Fprintf(w, `double=%v `, m.Double())
	}
	if _, ok := offsets.Get(13); ok {
		fmt.Fprintf(w, `float=%v `, m.Float())
	}
	if _, ok := offsets.Get(14); ok {
		fmt.Fprintf(w, `string=%s `, m.String())
	}
	if _, ok := offsets.Get(15); ok {
		fmt.Fprintf(w, `byte_array=%s `, strconv.Quote(string(m.ByteArray())))
	}
	if _, ok := offsets.Get(16); ok {
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
