// Copyright 2019 Daniel Harrison. All Rights Reserved.

package zeropb_test

import (
	"testing"

	"github.com/danhhz/zeropb/golden/raftgogopb"
	"github.com/danhhz/zeropb/golden/raftpb"
	"github.com/danhhz/zeropb/golden/raftzeropb"
	"github.com/golang/protobuf/proto"
)

const testByteArrayLen = 100

func testEntry(t testing.TB) []byte {
	term, index, typ := uint64(1), uint64(2), raftpb.EntryType_EntryNormal
	e := &raftpb.Entry{Term: &term, Index: &index, Type: &typ}
	e.Data = make([]byte, testByteArrayLen)
	for i := range e.Data {
		e.Data[i] = byte(testByteArrayLen + i)
	}
	buf, err := proto.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	return buf
}

func testMessage(t testing.TB, numEntries int) []byte {
	to, from, term, logTerm, index := uint64(1), uint64(2), uint64(3), uint64(4), uint64(5)
	m := &raftpb.Message{
		To: &to, From: &from, Term: &term, LogTerm: &logTerm, Index: &index,
		Context: make([]byte, testByteArrayLen),
		Snapshot: &raftpb.Snapshot{Metadata: &raftpb.SnapshotMetadata{
			Index: &index, Term: &term,
		}},
		Entries: make([]*raftpb.Entry, numEntries),
	}
	for i := range m.Context {
		m.Context[i] = byte(i)
	}
	for i := range m.Entries {
		m.Entries[i] = &raftpb.Entry{
			Term: &term, Index: &index, Data: make([]byte, testByteArrayLen),
		}
		for j := range m.Entries[i].Data {
			m.Entries[i].Data[j] = byte(i*testByteArrayLen + j)
		}
	}
	buf, err := proto.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	return buf
}

func byteSum(b []byte) uint64 {
	var x uint64
	for i := range b {
		x += uint64(b[i])
	}
	return x
}

func BenchmarkDecodeSimpleAccessNone(b *testing.B) {
	buf := testEntry(b)

	var pb raftpb.Entry
	b.Run(`pb`, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if err := proto.Unmarshal(buf, &pb); err != nil {
				b.Fatal(err)
			}
		}
	})

	var gogopb raftgogopb.Entry
	b.Run(`gogopb`, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if err := proto.Unmarshal(buf, &gogopb); err != nil {
				b.Fatal(err)
			}
		}
	})

	var zeropb raftzeropb.Entry
	b.Run(`zeropb`, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if err := zeropb.Decode(buf); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkDecodeSimpleAccessAll(b *testing.B) {
	buf := testEntry(b)

	var pb raftpb.Entry
	b.Run(`pb`, func(b *testing.B) {
		var x uint64
		for i := 0; i < b.N; i++ {
			if err := proto.Unmarshal(buf, &pb); err != nil {
				b.Fatal(err)
			}
			x += pb.GetTerm() + pb.GetIndex() + uint64(pb.GetType()) + byteSum(pb.GetData())
		}
	})

	var gogopb raftgogopb.Entry
	b.Run(`gogopb`, func(b *testing.B) {
		var x uint64
		for i := 0; i < b.N; i++ {
			if err := proto.Unmarshal(buf, &gogopb); err != nil {
				b.Fatal(err)
			}
			x += gogopb.Term + gogopb.Index + uint64(gogopb.Type) + byteSum(pb.Data)
		}
	})

	var zeropb raftzeropb.Entry
	b.Run(`zeropb`, func(b *testing.B) {
		var x uint64
		for i := 0; i < b.N; i++ {
			if err := zeropb.Decode(buf); err != nil {
				b.Fatal(err)
			}
			x += zeropb.Term() + zeropb.Index() + uint64(zeropb.Type()) + byteSum(zeropb.Data())
		}
	})
}

func BenchmarkDecodeSimpleAccessRepeatedly(b *testing.B) {
	const numAccessRepetitions = 3
	buf := testEntry(b)

	var pb raftpb.Entry
	b.Run(`pb`, func(b *testing.B) {
		var x uint64
		for i := 0; i < b.N; i++ {
			if err := proto.Unmarshal(buf, &pb); err != nil {
				b.Fatal(err)
			}
			for i := 0; i < numAccessRepetitions; i++ {
				x += pb.GetTerm() + pb.GetIndex() + uint64(pb.GetType()) + byteSum(pb.GetData())
			}
		}
	})

	var gogopb raftgogopb.Entry
	b.Run(`gogopb`, func(b *testing.B) {
		var x uint64
		for i := 0; i < b.N; i++ {
			if err := proto.Unmarshal(buf, &gogopb); err != nil {
				b.Fatal(err)
			}
			for i := 0; i < numAccessRepetitions; i++ {
				x += gogopb.Term + gogopb.Index + uint64(gogopb.Type) + byteSum(pb.Data)
			}
		}
	})

	var zeropb raftzeropb.Entry
	b.Run(`zeropb`, func(b *testing.B) {
		var x uint64
		for i := 0; i < b.N; i++ {
			if err := zeropb.Decode(buf); err != nil {
				b.Fatal(err)
			}
			for i := 0; i < numAccessRepetitions; i++ {
				x += zeropb.Term() + zeropb.Index() + uint64(zeropb.Type()) + byteSum(zeropb.Data())
			}
		}
	})
}

func BenchmarkDecodeComplexAccessOne(b *testing.B) {
	buf := testMessage(b, 3)

	var pb raftpb.Message
	b.Run(`pb`, func(b *testing.B) {
		var x uint64
		for i := 0; i < b.N; i++ {
			if err := proto.Unmarshal(buf, &pb); err != nil {
				b.Fatal(err)
			}
			x += pb.GetTo()
		}
	})

	var gogopb raftgogopb.Message
	b.Run(`gogopb`, func(b *testing.B) {
		var x uint64
		for i := 0; i < b.N; i++ {
			if err := proto.Unmarshal(buf, &gogopb); err != nil {
				b.Fatal(err)
			}
			x += gogopb.To
		}
	})

	var zeropb raftzeropb.Message
	b.Run(`zeropb`, func(b *testing.B) {
		var x uint64
		for i := 0; i < b.N; i++ {
			if err := zeropb.Decode(buf); err != nil {
				b.Fatal(err)
			}
			x += zeropb.To()
		}
	})
}

func BenchmarkDecodeComplexAccessRepeatedMessage(b *testing.B) {
	buf := testMessage(b, 3)

	var pb raftpb.Message
	b.Run(`pb`, func(b *testing.B) {
		var x uint64
		for i := 0; i < b.N; i++ {
			if err := proto.Unmarshal(buf, &pb); err != nil {
				b.Fatal(err)
			}
			for _, e := range pb.GetEntries() {
				if e == nil {
					continue
				}
				x += e.GetTerm() + e.GetIndex() + uint64(e.GetType()) + byteSum(e.GetData())
			}
		}
	})

	var gogopb raftgogopb.Message
	b.Run(`gogopb`, func(b *testing.B) {
		var x uint64
		for i := 0; i < b.N; i++ {
			if err := proto.Unmarshal(buf, &gogopb); err != nil {
				b.Fatal(err)
			}
			for _, e := range gogopb.Entries {
				x += e.Term + e.Index + uint64(e.Type) + byteSum(e.Data)
			}
		}
	})

	var zeropb raftzeropb.Message
	var zeropbE raftzeropb.Entry
	b.Run(`zeropb`, func(b *testing.B) {
		var x uint64
		for i := 0; i < b.N; i++ {
			if err := zeropb.Decode(buf); err != nil {
				b.Fatal(err)
			}
			it := zeropb.Entries()
			for {
				if ok, err := it.Next(&zeropbE); !ok {
					break
				} else if err != nil {
					b.Fatal(err)
				}
				x += zeropbE.Term() + zeropbE.Index() + uint64(zeropbE.Type()) + byteSum(zeropbE.Data())
			}
		}
	})
}
