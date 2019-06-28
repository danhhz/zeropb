// Copyright 2019 Daniel Harrison. All Rights Reserved.

package zeropb

import (
	"encoding/binary"
	"io"
	"math"
	"reflect"
	"unsafe"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

// GetBool gets an encoded bool field with the given field id.
func GetBool(buf []byte, offsets Offsets, fieldID int) bool {
	offset, ok := offsets.Get(fieldID)
	if !ok {
		return false
	}
	x, _ := binary.Uvarint(buf[offset:])
	return x > 0
}

// GetInt32 gets an encoded int32 field with the given field id.
func GetInt32(buf []byte, offsets Offsets, fieldID int) int32 {
	offset, ok := offsets.Get(fieldID)
	if !ok {
		return 0
	}
	x, _ := binary.Uvarint(buf[offset:])
	return int32(x)
}

// GetInt64 gets an encoded int64 field with the given field id.
func GetInt64(buf []byte, offsets Offsets, fieldID int) int64 {
	offset, ok := offsets.Get(fieldID)
	if !ok {
		return 0
	}
	x, _ := binary.Uvarint(buf[offset:])
	return int64(x)
}

// GetUint32 gets an encoded uint32 field with the given field id.
func GetUint32(buf []byte, offsets Offsets, fieldID int) uint32 {
	offset, ok := offsets.Get(fieldID)
	if !ok {
		return 0
	}
	x, _ := binary.Uvarint(buf[offset:])
	return uint32(x)
}

// GetUint64 gets an encoded uint64 field with the given field id.
func GetUint64(buf []byte, offsets Offsets, fieldID int) uint64 {
	offset, ok := offsets.Get(fieldID)
	if !ok {
		return 0
	}
	x, _ := binary.Uvarint(buf[offset:])
	return x
}

// GetZigZagInt32 gets a zig-zag encoded int32 field with the given field id.
func GetZigZagInt32(buf []byte, offsets Offsets, fieldID int) int32 {
	offset, ok := offsets.Get(fieldID)
	if !ok {
		return 0
	}
	x, _ := binary.Varint(buf[offset:])
	return int32(x)
}

// GetZigZagInt64 gets a zig-zag encoded int64 field with the given field id.
func GetZigZagInt64(buf []byte, offsets Offsets, fieldID int) int64 {
	offset, ok := offsets.Get(fieldID)
	if !ok {
		return 0
	}
	x, _ := binary.Varint(buf[offset:])
	return x
}

// GetFixedUint32 gets a little-endian encoded uint32 field with the given field
// id.
func GetFixedUint32(buf []byte, offsets Offsets, fieldID int) uint32 {
	offset, ok := offsets.Get(fieldID)
	if !ok {
		return 0
	}
	return binary.LittleEndian.Uint32(buf[offset : offset+4])
}

// GetFixedUint64 gets a little-endian encoded uint64 field with the given field
// id.
func GetFixedUint64(buf []byte, offsets Offsets, fieldID int) uint64 {
	offset, ok := offsets.Get(fieldID)
	if !ok {
		return 0
	}
	return binary.LittleEndian.Uint64(buf[offset : offset+8])
}

// GetFixedInt32 gets a little-endian encoded int32 field with the given field
// id.
func GetFixedInt32(buf []byte, offsets Offsets, fieldID int) int32 {
	offset, ok := offsets.Get(fieldID)
	if !ok {
		return 0
	}
	return int32(binary.LittleEndian.Uint32(buf[offset : offset+4]))
}

// GetFixedInt64 gets a little-endian encoded int64 field with the given field
// id.
func GetFixedInt64(buf []byte, offsets Offsets, fieldID int) int64 {
	offset, ok := offsets.Get(fieldID)
	if !ok {
		return 0
	}
	return int64(binary.LittleEndian.Uint64(buf[offset : offset+8]))
}

// GetFloat32 gets a little-endian, IEEE 754 encoded float32 field with the
// given field id.
func GetFloat32(buf []byte, offsets Offsets, fieldID int) float32 {
	offset, ok := offsets.Get(fieldID)
	if !ok {
		return 0
	}
	bits := binary.LittleEndian.Uint32(buf[offset : offset+4])
	return math.Float32frombits(bits)
}

// GetFloat64 gets a little-endian, IEEE 754 encoded float64 field with the
// given field id.
func GetFloat64(buf []byte, offsets Offsets, fieldID int) float64 {
	offset, ok := offsets.Get(fieldID)
	if !ok {
		return 0
	}
	bits := binary.LittleEndian.Uint64(buf[offset : offset+8])
	return math.Float64frombits(bits)
}

// GetString gets an encoded string field with the given field id.
func GetString(buf []byte, offsets Offsets, fieldID int) string {
	offset, ok := offsets.Get(fieldID)
	if !ok {
		return ``
	}
	len, lenSize := binary.Uvarint(buf[offset:])
	offset += uint64(lenSize)
	data := buf[offset : offset+len]
	return *(*string)(unsafe.Pointer(&data))
}

// GetBytes gets an encoded []byte field with the given field id.
func GetBytes(buf []byte, offsets Offsets, fieldID int) []byte {
	offset, ok := offsets.Get(fieldID)
	if !ok {
		return nil
	}
	len, lenSize := binary.Uvarint(buf[offset:])
	offset += uint64(lenSize)
	return buf[offset : offset+len]
}

// SetBool set an encoded bool field with the given field id, overwriting an
// old value if present.
func SetBool(buf *[]byte, offsets Offsets, fieldID int, x bool) {
	var scratch [1]byte
	if x {
		scratch[0] = 1
	}
	var valLen, val []byte = nil, scratch[:]
	setSingle(buf, offsets, fieldID, proto.WireVarint, valLen, val)
}

// SetInt32 set an encoded int32 field with the given field id, overwriting an
// old value if present.
func SetInt32(buf *[]byte, offsets Offsets, fieldID int, x int32) {
	// WIP huh? var scratch [binary.MaxVarintLen32]byte
	var scratch [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(scratch[:], uint64(x))
	var valLen, val []byte = nil, scratch[:n]
	setSingle(buf, offsets, fieldID, proto.WireVarint, valLen, val)
}

// SetInt64 set an encoded int64 field with the given field id, overwriting an
// old value if present.
func SetInt64(buf *[]byte, offsets Offsets, fieldID int, x int64) {
	var scratch [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(scratch[:], uint64(x))
	var valLen, val []byte = nil, scratch[:n]
	setSingle(buf, offsets, fieldID, proto.WireVarint, valLen, val)
}

// SetUint32 set an encoded uint64 field with the given field id, overwriting an
// old value if present.
func SetUint32(buf *[]byte, offsets Offsets, fieldID int, x uint32) {
	var scratch [binary.MaxVarintLen32]byte
	n := binary.PutUvarint(scratch[:], uint64(x))
	var valLen, val []byte = nil, scratch[:n]
	setSingle(buf, offsets, fieldID, proto.WireVarint, valLen, val)
}

// SetUint64 set an encoded uint64 field with the given field id, overwriting an
// old value if present.
func SetUint64(buf *[]byte, offsets Offsets, fieldID int, x uint64) {
	var scratch [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(scratch[:], x)
	var valLen, val []byte = nil, scratch[:n]
	setSingle(buf, offsets, fieldID, proto.WireVarint, valLen, val)
}

// SetZigZagInt32 set a zig-zag encoded int64 field with the given field id,
// overwriting an old value if present.
func SetZigZagInt32(buf *[]byte, offsets Offsets, fieldID int, x int32) {
	var scratch [binary.MaxVarintLen32]byte
	n := binary.PutVarint(scratch[:], int64(x))
	var valLen, val []byte = nil, scratch[:n]
	setSingle(buf, offsets, fieldID, proto.WireVarint, valLen, val)
}

// SetZigZagInt64 set a zig-zag encoded int64 field with the given field id,
// overwriting an old value if present.
func SetZigZagInt64(buf *[]byte, offsets Offsets, fieldID int, x int64) {
	var scratch [binary.MaxVarintLen64]byte
	n := binary.PutVarint(scratch[:], x)
	var valLen, val []byte = nil, scratch[:n]
	setSingle(buf, offsets, fieldID, proto.WireVarint, valLen, val)
}

// SetFixedUint32 set a little-endian encoded uint32 field with the given field
// id, overwriting an old value if present.
func SetFixedUint32(buf *[]byte, offsets Offsets, fieldID int, x uint32) {
	var scratch [4]byte
	binary.LittleEndian.PutUint32(scratch[:], x)
	var valLen, val []byte = nil, scratch[:]
	setSingle(buf, offsets, fieldID, proto.WireFixed32, valLen, val)
}

// SetFixedUint64 set a little-endian encoded uint64 field with the given field
// id, overwriting an old value if present.
func SetFixedUint64(buf *[]byte, offsets Offsets, fieldID int, x uint64) {
	var scratch [8]byte
	binary.LittleEndian.PutUint64(scratch[:], x)
	var valLen, val []byte = nil, scratch[:]
	setSingle(buf, offsets, fieldID, proto.WireFixed64, valLen, val)
}

// SetFixedInt32 set a little-endian encoded int32 field with the given field
// id, overwriting an old value if present.
func SetFixedInt32(buf *[]byte, offsets Offsets, fieldID int, x int32) {
	var scratch [4]byte
	binary.LittleEndian.PutUint32(scratch[:], uint32(x))
	var valLen, val []byte = nil, scratch[:]
	setSingle(buf, offsets, fieldID, proto.WireFixed32, valLen, val)
}

// SetFixedInt64 set a little-endian encoded int64 field with the given field
// id, overwriting an old value if present.
func SetFixedInt64(buf *[]byte, offsets Offsets, fieldID int, x int64) {
	var scratch [8]byte
	binary.LittleEndian.PutUint64(scratch[:], uint64(x))
	var valLen, val []byte = nil, scratch[:]
	setSingle(buf, offsets, fieldID, proto.WireFixed64, valLen, val)
}

// SetFloat32 set a little-endian, IEEE 754 encoded float32 field with the given
// field id, overwriting an old value if present.
func SetFloat32(buf *[]byte, offsets Offsets, fieldID int, x float32) {
	var scratch [4]byte
	bits := math.Float32bits(x)
	binary.LittleEndian.PutUint32(scratch[:], bits)
	var valLen, val []byte = nil, scratch[:]
	setSingle(buf, offsets, fieldID, proto.WireFixed32, valLen, val)
}

// SetFloat64 set a little-endian, IEEE 754 encoded float64 field with the given
// field id, overwriting an old value if present.
func SetFloat64(buf *[]byte, offsets Offsets, fieldID int, x float64) {
	var scratch [8]byte
	bits := math.Float64bits(x)
	binary.LittleEndian.PutUint64(scratch[:], bits)
	var valLen, val []byte = nil, scratch[:]
	setSingle(buf, offsets, fieldID, proto.WireFixed64, valLen, val)
}

// SetString set an encoded string field with the given field id, overwriting an
// old value if present.
func SetString(buf *[]byte, offsets Offsets, fieldID int, x string) {
	var scratch [binary.MaxVarintLen64]byte
	hdr := *(*reflect.StringHeader)(unsafe.Pointer(&x))
	data := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: hdr.Data,
		Len:  hdr.Len,
		Cap:  hdr.Len,
	}))
	n := binary.PutUvarint(scratch[:], uint64(len(data)))
	var valLen, val []byte = scratch[:n], data
	setSingle(buf, offsets, fieldID, proto.WireBytes, valLen, val)
}

// SetBytes set an encoded []byte field with the given field id, overwriting an
// old value if present.
func SetBytes(buf *[]byte, offsets Offsets, fieldID int, x []byte) {
	var scratch [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(scratch[:], uint64(len(x)))
	var valLen, val []byte = scratch[:n], x
	setSingle(buf, offsets, fieldID, proto.WireBytes, valLen, val)
}

// AppendBytes appends one value to an encoded []byte field with the given field
// id.
func AppendBytes(buf *[]byte, offsets Offsets, fieldID int, x []byte) {
	var scratch [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(scratch[:], uint64(len(x)))
	var valLen, val []byte = scratch[:n], x
	appendSingle(buf, offsets, fieldID, proto.WireBytes, valLen, val)
}

func setSingle(buf *[]byte, offsets Offsets, fieldID, typ int, valLen, val []byte) {
	// TODO(dan): We could save space if we didn't encode fields set to the
	// default value.

	var scratch [binary.MaxVarintLen64]byte
	tag := uint64(fieldID<<3 | typ)
	tagLen := binary.PutUvarint(scratch[:], tag)

	if offset, ok := offsets.Get(fieldID); ok {
		var oldLen int
		switch typ {
		case proto.WireVarint:
			_, oldLen = binary.Uvarint((*buf)[offset:])
		case proto.WireFixed32:
			oldLen = 4
		case proto.WireFixed64:
			oldLen = 8
		case proto.WireBytes:
			oldValLen, oldValLenLen := binary.Uvarint((*buf)[offset:])
			oldLen = int(oldValLen) + oldValLenLen
		default:
			panic(typ)
		}

		// Performance optimization: if the new value is exactly the same length as
		// the old one, simply overwrite it.
		if int(oldLen) == len(valLen)+len(val) {
			offset += uint64(copy((*buf)[offset:], valLen))
			offset += uint64(copy((*buf)[offset:], val))
			return
		}

		fieldBegin, fieldEnd := offset-uint64(tagLen), offset+uint64(oldLen)
		// TODO(dan): sanity check that we got the same tag bytes here?
		copy((*buf)[fieldBegin:], (*buf)[fieldEnd:])
		removedBytes := fieldEnd - fieldBegin
		(*buf) = (*buf)[:uint64(len(*buf))-removedBytes]

		// Adjust every offset in the map that's larger than fieldEnd down by
		// removedBytes.
		//
		// TODO(dan): Figure out how to move all this internal knowledge into
		// Offsets.
		for fieldID, offset := range offsets.a {
			if offset != math.MaxUint16 && uint64(offset) > fieldEnd {
				// Guaranteed to still be > 0 and < MaxOffsetInArray.
				offsets.a[fieldID] = offset - uint16(removedBytes)
			}
		}
		if offsets.m != nil {
			for fieldID, offset := range *offsets.m {
				if offset > fieldEnd {
					offsets.Set(fieldID, offset-removedBytes)
				}
			}
		}
	}
	*buf = append(*buf, scratch[:tagLen]...)

	offset := len(*buf)
	offsets.Set(fieldID, uint64(offset))
	*buf = append(*buf, valLen...)
	*buf = append(*buf, val...)
}

func appendSingle(buf *[]byte, offsets Offsets, fieldID, typ int, valLen, val []byte) {
	var scratch [binary.MaxVarintLen64]byte

	tag := uint64(fieldID<<3 | typ)
	tagLen := binary.PutUvarint(scratch[:], tag)
	*buf = append(*buf, scratch[:tagLen]...)

	offset := len(*buf)
	// The offset for repeated fields always points at the first one.
	if _, ok := offsets.Get(fieldID); !ok {
		offsets.Set(fieldID, uint64(offset))
	}
	*buf = append(*buf, valLen...)
	*buf = append(*buf, val...)
}

// GetRepeatedNonPacked returns a slice of the internal buffer, starting at the
// first instance of a non-packed repeated field and with everything after it.
// It's positioned past the tag of the field, so the very first thing will be a
// varint length. This is exactly the input expected by FindNextField.
func GetRepeatedNonPacked(buf []byte, offsets Offsets, fieldID int) []byte {
	offset, ok := offsets.Get(fieldID)
	if !ok {
		return nil
	}
	return buf[offset:]
}

// Decode parses and validates a proto message, filling in a map of field id ->
// offset as it goes.
func Decode(buf []byte, offsets Offsets, repeatedFields RepeatedFields) error {
	offsets.Clear()
	for idx := 0; ; {
		if idx == len(buf) {
			return nil
		} else if idx > len(buf) {
			return io.ErrUnexpectedEOF
		}
		tag, size := binary.Uvarint(buf[idx:])
		if size == 0 {
			return io.ErrUnexpectedEOF
		}
		idx += size
		if idx >= len(buf) {
			return io.ErrUnexpectedEOF
		}

		fieldID := int(tag >> 3)
		typ := tag & 0x7
		// fmt.Printf("tag=%x field=%d typ=%d idx=%d\n", tag, fieldID, typ, idx)

		// TODO(dan): It's odd that we're passing in only the repeatedness of the
		// field. Maybe pass in some sort of field descriptor instead?
		if _, ok := repeatedFields[fieldID]; ok {
			// For a repeated field, set the offset if this is the first one we've
			// found but don't overwrite an earlier offset. The repeated field
			// iterator will start from this offset and re-parse the rest of the
			// message to find the rest of them.
			if _, ok := offsets.Get(fieldID); !ok {
				offsets.Set(fieldID, uint64(idx))
			}
		} else {
			// For a non-repeated field, always overwrite the offset. The proto spec
			// says to keep only the last one.
			offsets.Set(fieldID, uint64(idx))
			// TODO(dan): The proto spec specifically allows a single message to be
			// split, but I don't understand why it would happen. It's hard to
			// support so leave it unsupported for now.
		}

		s := fieldDataSize(typ, buf[idx:])
		if s == 0 {
			return io.ErrUnexpectedEOF
		}
		idx += s
	}
}

// FindNextField parses a bytes field (with no tag) that must be present at the
// very start of buf and finds the next requested field, returning it ready to
// be handed back to FindNextField.
func FindNextField(buf []byte, field int) ([]byte, []byte) {
	if buf == nil {
		return nil, nil
	}
	// NB: FindNextField doesn't do the same bounds checking that Decode does
	// because Decode has already validated the input.
	msgLen, lenSize := binary.Uvarint(buf)
	msg := buf[lenSize : lenSize+int(msgLen)]
	for idx := lenSize + int(msgLen); idx < len(buf); {
		tag, size := binary.Uvarint(buf[idx:])
		idx += size
		f := int(tag >> 3)
		typ := tag & 0x7
		if f == field {
			return buf[idx:], msg
		}
		idx += fieldDataSize(typ, buf[idx:])
	}
	return nil, msg
}

func fieldDataSize(typ uint64, buf []byte) int {
	switch typ {
	case proto.WireVarint:
		_, size := binary.Uvarint(buf)
		return size
	case proto.WireFixed32:
		return 4
	case proto.WireFixed64:
		return 8
	case proto.WireBytes:
		len, lenSize := binary.Uvarint(buf)
		if lenSize == 0 {
			return 0
		}
		return lenSize + int(len)
	default:
		panic(errors.Errorf(`unsupported type: %d`, typ))
	}
}

// RepeatedFields indicates which fields are repeated, indexed by field id.
type RepeatedFields map[int]struct{}
