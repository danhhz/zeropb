// Copyright 2019 Daniel Harrison. All Rights Reserved.

package zeropb_test

import (
	fmt "fmt"
	"testing"

	"github.com/danhhz/zeropb/golden/raftpb"
	"github.com/danhhz/zeropb/golden/raftzeropb"
	"github.com/golang/protobuf/proto"
)

func TestEntry(t *testing.T) {
	term, index, typ := uint64(1), uint64(2), raftpb.EntryType_EntryNormal
	e := &raftpb.Entry{Term: &term, Index: &index, Type: &typ, Data: []byte{5, 6}}
	buf, err := proto.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}

	var ez raftzeropb.Entry
	if err := ez.Decode(buf); err != nil {
		t.Fatal(err)
	}
	fmt.Println(ez.Term(), ez.Index(), ez.Type(), ez.Data())
}

func TestMessage(t *testing.T) {
	term1, term2 := uint64(1), uint64(2)
	m := &raftpb.Message{Entries: []*raftpb.Entry{
		{Term: &term1},
		{Term: &term2},
	}}
	buf, err := proto.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}

	var mz raftzeropb.Message
	if err := mz.Decode(buf); err != nil {
		t.Fatal(err)
	}
	var ez raftzeropb.Entry
	for it := mz.Entries(); ; {
		if ok, err := it.Next(&ez); err != nil {
			t.Fatal(err)
		} else if !ok {
			break
		}
		fmt.Println(ez.Term(), ez.Index(), ez.Type(), ez.Data())
	}
}
