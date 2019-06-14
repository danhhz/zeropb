// Copyright 2019 Daniel Harrison. All Rights Reserved.

package zeropb_test

import (
	"encoding/binary"
	"math"
	"testing"

	"github.com/danhhz/zeropb/golden/raftzeropb"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	var e raftzeropb.Entry
	assert.Equal(t, uint64(0), e.Index())
	assert.Equal(t, uint64(0), e.Term())
	assert.Equal(t, []byte(nil), e.Data())

	// Set previously unset fields.
	e.SetIndex(1)
	e.SetTerm(2)
	e.SetData([]byte{3, 4})
	assert.Equal(t, uint64(1), e.Index())
	assert.Equal(t, uint64(2), e.Term())
	assert.Equal(t, []byte{3, 4}, e.Data())

	// Overwrite previously set fields with same sized data (fast path).
	e.SetIndex(5)
	assert.Equal(t, uint64(5), e.Index())
	e.SetData([]byte{7, 8})
	assert.Equal(t, []byte{7, 8}, e.Data())

	// Overwrite previously set data with different sized data.
	e.SetData([]byte{9, 10, 11})
	assert.Equal(t, []byte{9, 10, 11}, e.Data())
	e.SetIndex(1000000)
	assert.Equal(t, uint64(1000000), e.Index())

	// The previous updates had to remove some data in the middle of the encoding,
	// which then requires updating the offsets map. Double check that everything
	// is still what we expect it to be.
	assert.Equal(t, uint64(1000000), e.Index())
	assert.Equal(t, uint64(2), e.Term())
	assert.Equal(t, []byte{9, 10, 11}, e.Data())
}

type encodeBuf []byte

func (b *encodeBuf) appendTag(fieldID, typ int) *encodeBuf {
	tag := uint64(fieldID<<3 | typ)
	return b.appendVarint(tag)
}

func (b *encodeBuf) appendVarint(x uint64) *encodeBuf {
	var scratch [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(scratch[:], x)
	*b = append(*b, scratch[:n]...)
	return b
}

func TestDecode(t *testing.T) {
	var e raftzeropb.Entry

	// Max message length
	assert.EqualError(t, e.Decode(make([]byte, math.MaxUint16+1)),
		`cannot decode messages longer than 65535 bytes`)

	// Incomplete varint tag
	{
		var buf encodeBuf
		buf.appendVarint(1000000)
		buf = buf[:len(buf)-1]
		assert.EqualError(t, e.Decode(buf), `unexpected EOF`)
	}

	// Varint tag with no data
	{
		var buf encodeBuf
		buf.appendTag(1, proto.WireVarint)
		assert.EqualError(t, e.Decode(buf), `unexpected EOF`)
	}

	// Varint tag with incomplete data
	{
		var buf encodeBuf
		buf.appendTag(1, proto.WireVarint)
		buf.appendVarint(1000000)
		buf = buf[:len(buf)-1]
		assert.EqualError(t, e.Decode(buf), `unexpected EOF`)
	}

	// Fixed32 tag with insufficient data
	{
		var buf encodeBuf
		buf.appendTag(1, proto.WireFixed32)
		buf = append(buf, 0)
		assert.EqualError(t, e.Decode(buf), `unexpected EOF`)
	}

	// Fixed64 tag with insufficient data
	{
		var buf encodeBuf
		buf.appendTag(1, proto.WireFixed64)
		buf = append(buf, 0)
		assert.EqualError(t, e.Decode(buf), `unexpected EOF`)
	}

	// Bytes tag with no length
	{
		var buf encodeBuf
		buf.appendTag(1, proto.WireBytes)
		assert.EqualError(t, e.Decode(buf), `unexpected EOF`)
	}

	// Varint tag with incomplete length
	{
		var buf encodeBuf
		buf.appendTag(1, proto.WireBytes)
		buf.appendVarint(1000000)
		buf = buf[:len(buf)-1]
		assert.EqualError(t, e.Decode(buf), `unexpected EOF`)
	}

	// Bytes tag with length longer than remaining data
	{
		var buf encodeBuf
		buf.appendTag(1, proto.WireBytes)
		buf.appendVarint(10)
		buf = append(buf, 1, 2, 3)
		assert.EqualError(t, e.Decode(buf), `unexpected EOF`)
	}
}
