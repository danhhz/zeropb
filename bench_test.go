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

func testEntry() raftgogopb.Entry {
	e := raftgogopb.Entry{Term: 1, Index: 2, Type: raftgogopb.EntryNormal, Data: make([]byte, testByteArrayLen)}
	for i := range e.Data {
		e.Data[i] = byte(testByteArrayLen + i)
	}
	return e
}

func testEntryEncoded(t testing.TB) []byte {
	e := testEntry()
	buf, err := proto.Marshal(&e)
	if err != nil {
		t.Fatal(err)
	}
	return buf
}

func testMessage(numEntries int) raftgogopb.Message {
	m := raftgogopb.Message{
		To: 1, From: 2, Term: 3, LogTerm: 4, Index: 5,
		Context: make([]byte, testByteArrayLen),
		Snapshot: raftgogopb.Snapshot{Metadata: raftgogopb.SnapshotMetadata{
			Term: 3, Index: 5,
		}},
		Entries: make([]raftgogopb.Entry, numEntries),
	}
	for i := range m.Context {
		m.Context[i] = byte(i)
	}
	for i := range m.Entries {
		m.Entries[i] = raftgogopb.Entry{
			Term: 3, Index: 5, Data: make([]byte, testByteArrayLen),
		}
		for j := range m.Entries[i].Data {
			m.Entries[i].Data[j] = byte(i*testByteArrayLen + j)
		}
	}
	return m
}

func testMessageEncoded(t testing.TB, numEntries int) []byte {
	m := testMessage(numEntries)
	buf, err := proto.Marshal(&m)
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
	buf := testEntryEncoded(b)

	var pb raftpb.Entry
	b.Run(`pb`, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if err := proto.Unmarshal(buf, &pb); err != nil {
				b.Fatal(err)
			}
		}
		b.SetBytes(int64(len(buf)))
	})

	var gogopb raftgogopb.Entry
	b.Run(`gogopb`, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if err := proto.Unmarshal(buf, &gogopb); err != nil {
				b.Fatal(err)
			}
		}
		b.SetBytes(int64(len(buf)))
	})

	var zeropb raftzeropb.Entry
	b.Run(`zeropb`, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if err := zeropb.Decode(buf); err != nil {
				b.Fatal(err)
			}
		}
		b.SetBytes(int64(len(buf)))
	})
}

func BenchmarkDecodeSimpleAccessAll(b *testing.B) {
	buf := testEntryEncoded(b)

	var pb raftpb.Entry
	b.Run(`pb`, func(b *testing.B) {
		var x uint64
		for i := 0; i < b.N; i++ {
			if err := proto.Unmarshal(buf, &pb); err != nil {
				b.Fatal(err)
			}
			x += pb.GetTerm() + pb.GetIndex() + uint64(pb.GetType()) + byteSum(pb.GetData())
		}
		b.SetBytes(int64(len(buf)))
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
		b.SetBytes(int64(len(buf)))
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
		b.SetBytes(int64(len(buf)))
	})
}

func BenchmarkDecodeSimpleAccessRepeatedly(b *testing.B) {
	const numAccessRepetitions = 3
	buf := testEntryEncoded(b)

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
		b.SetBytes(int64(len(buf)))
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
		b.SetBytes(int64(len(buf)))
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
		b.SetBytes(int64(len(buf)))
	})
}

func BenchmarkDecodeComplexAccessOne(b *testing.B) {
	buf := testMessageEncoded(b, 3)

	var pb raftpb.Message
	b.Run(`pb`, func(b *testing.B) {
		var x uint64
		for i := 0; i < b.N; i++ {
			if err := proto.Unmarshal(buf, &pb); err != nil {
				b.Fatal(err)
			}
			x += pb.GetTo()
		}
		b.SetBytes(int64(len(buf)))
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
		b.SetBytes(int64(len(buf)))
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
		b.SetBytes(int64(len(buf)))
	})
}

func BenchmarkDecodeComplexAccessRepeatedMessage(b *testing.B) {
	buf := testMessageEncoded(b, 3)

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
		b.SetBytes(int64(len(buf)))
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
		b.SetBytes(int64(len(buf)))
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
		b.SetBytes(int64(len(buf)))
	})
}

func BenchmarkEncodeSimpleSetAll(b *testing.B) {
	e := testEntry()

	pbBuf := proto.NewBuffer(make([]byte, 0, 100+testByteArrayLen))
	b.Run(`pb`, func(b *testing.B) {
		var x int64
		for i := 0; i < b.N; i++ {
			pbBuf.Reset()
			typ := raftpb.EntryType(e.Type)
			pb := raftpb.Entry{Term: &e.Term, Index: &e.Index, Type: &typ, Data: e.Data}
			if err := pbBuf.Marshal(&pb); err != nil {
				b.Fatal(err)
			}
			x += int64(len(pbBuf.Bytes()))
		}
		b.SetBytes(x / int64(b.N))
	})

	gogopbBuf := proto.NewBuffer(make([]byte, 0, 100+testByteArrayLen))
	b.Run(`gogopb`, func(b *testing.B) {
		var x int64
		for i := 0; i < b.N; i++ {
			gogopbBuf.Reset()
			gogopb := raftgogopb.Entry{Term: e.Term, Index: e.Index, Type: e.Type, Data: e.Data}
			if err := gogopbBuf.Marshal(&gogopb); err != nil {
				b.Fatal(err)
			}
			x += int64(len(gogopbBuf.Bytes()))
		}
		b.SetBytes(x / int64(b.N))
	})

	zeropbBuf := make([]byte, 0, 100+testByteArrayLen)
	b.Run(`zeropb`, func(b *testing.B) {
		var x int64
		for i := 0; i < b.N; i++ {
			var zeropb raftzeropb.Entry
			zeropb.Reset(zeropbBuf)
			zeropb.SetIndex(e.Index)
			zeropb.SetTerm(e.Term)
			zeropb.SetType(uint32(e.Type))
			zeropb.SetData(e.Data)
			x += int64(len(zeropb.Encode()))
		}
		b.SetBytes(x / int64(b.N))
	})
}

func BenchmarkEncodeSimpleSetRepeatedly(b *testing.B) {
	const numSetRepetitions = 3
	const bufLen = 100 + testByteArrayLen
	e := testEntry()

	pbBuf := proto.NewBuffer(make([]byte, 0, bufLen))
	b.Run(`pb`, func(b *testing.B) {
		var x int64
		for i := 0; i < b.N; i++ {
			pbBuf.Reset()
			var pb raftpb.Entry
			for j := 0; j < numSetRepetitions; j++ {
				typ := raftpb.EntryType(e.Type)
				pb = raftpb.Entry{Term: &e.Term, Index: &e.Index, Type: &typ, Data: e.Data}
			}
			if err := pbBuf.Marshal(&pb); err != nil {
				b.Fatal(err)
			}
			x += int64(len(pbBuf.Bytes()))
		}
		b.SetBytes(x / int64(b.N))
	})

	gogopbBuf := proto.NewBuffer(make([]byte, 0, bufLen))
	b.Run(`gogopb`, func(b *testing.B) {
		var x int64
		for i := 0; i < b.N; i++ {
			gogopbBuf.Reset()
			var gogopb raftgogopb.Entry
			for j := 0; j < numSetRepetitions; j++ {
				gogopb = raftgogopb.Entry{Term: e.Term, Index: e.Index, Type: e.Type, Data: e.Data}
			}
			if err := gogopbBuf.Marshal(&gogopb); err != nil {
				b.Fatal(err)
			}
			x += int64(len(gogopbBuf.Bytes()))
		}
		b.SetBytes(x / int64(b.N))
	})

	zeropbBuf := make([]byte, 0, bufLen)
	b.Run(`zeropb`, func(b *testing.B) {
		var x int64
		for i := 0; i < b.N; i++ {
			var zeropb raftzeropb.Entry
			zeropb.Reset(zeropbBuf)
			for j := 0; j < numSetRepetitions; j++ {
				zeropb.SetIndex(e.Index)
				zeropb.SetTerm(e.Term)
				zeropb.SetType(uint32(e.Type))
				zeropb.SetData(e.Data)
			}
			x += int64(len(zeropb.Encode()))
		}
		b.SetBytes(x / int64(b.N))
	})
}

func BenchmarkEncodeComplex(b *testing.B) {
	const numEntries = 3
	const bufLen = 3 * (100 + testByteArrayLen)
	m := testMessage(numEntries)

	pbBuf := proto.NewBuffer(make([]byte, 0, bufLen))
	b.Run(`pb`, func(b *testing.B) {
		var x int64
		for i := 0; i < b.N; i++ {
			pbBuf.Reset()

			pb := &raftpb.Message{
				To: &m.To, From: &m.From, Term: &m.Term, LogTerm: &m.LogTerm, Index: &m.Index,
				Context: m.Context,
				Snapshot: &raftpb.Snapshot{Metadata: &raftpb.SnapshotMetadata{
					Index: &m.Snapshot.Metadata.Index, Term: &m.Snapshot.Metadata.Term,
				}},
				Entries: make([]*raftpb.Entry, len(m.Entries)),
			}
			for i := range m.Entries {
				pb.Entries[i] = &raftpb.Entry{
					Term: &m.Entries[i].Term, Index: &m.Entries[i].Index, Data: m.Entries[i].Data,
				}
			}

			if err := pbBuf.Marshal(pb); err != nil {
				b.Fatal(err)
			}
			x += int64(len(pbBuf.Bytes()))
		}
		b.SetBytes(x / int64(b.N))
	})

	gogopbBuf := proto.NewBuffer(make([]byte, 0, bufLen))
	b.Run(`gogopb`, func(b *testing.B) {
		var x int64
		for i := 0; i < b.N; i++ {
			gogopbBuf.Reset()

			gogopb := raftgogopb.Message{
				To: m.To, From: m.From, Term: m.Term, LogTerm: m.LogTerm, Index: m.Index,
				Context: m.Context,
				Snapshot: raftgogopb.Snapshot{Metadata: raftgogopb.SnapshotMetadata{
					Index: m.Snapshot.Metadata.Index, Term: m.Snapshot.Metadata.Term,
				}},
				Entries: make([]raftgogopb.Entry, len(m.Entries)),
			}
			for i := range m.Entries {
				gogopb.Entries[i] = raftgogopb.Entry{
					Term: m.Entries[i].Term, Index: m.Entries[i].Index, Data: m.Entries[i].Data,
				}
			}

			if err := gogopbBuf.Marshal(&gogopb); err != nil {
				b.Fatal(err)
			}
			x += int64(len(gogopbBuf.Bytes()))
		}
		b.SetBytes(x / int64(b.N))
	})

	// TODO(dan): Eeek. This is not ergonomic.
	zeropbBufs := make([][]byte, 4)
	for i := range zeropbBufs {
		zeropbBufs[i] = make([]byte, 0, bufLen)
	}
	b.Run(`zeropb`, func(b *testing.B) {
		var x int64

		for i := 0; i < b.N; i++ {
			var zeroM raftzeropb.Message
			var zeroE raftzeropb.Entry
			var zeroS raftzeropb.Snapshot
			var zeroSM raftzeropb.SnapshotMetadata

			zeroM.Reset(zeropbBufs[0])
			zeroM.SetTo(m.To)
			zeroM.SetFrom(m.From)
			zeroM.SetTerm(m.Term)
			zeroM.SetLogTerm(m.LogTerm)
			zeroM.SetIndex(m.Index)
			zeroM.SetContext(m.Context)

			zeroSM.Reset(zeropbBufs[1])
			zeroSM.SetIndex(m.Snapshot.Metadata.Index)
			zeroSM.SetTerm(m.Snapshot.Metadata.Term)

			zeroS.Reset(zeropbBufs[2])
			zeroS.SetMetadata(zeroSM)

			for i := range m.Entries {
				zeroE.Reset(zeropbBufs[3])
				zeroE.SetTerm(m.Entries[i].Term)
				zeroE.SetIndex(m.Entries[i].Index)
				zeroE.SetData(m.Entries[i].Data)
				zeroM.AppendToEntries(zeroE)
			}

			x += int64(len(zeroM.Encode()))
		}
		b.SetBytes(x / int64(b.N))
	})
}
