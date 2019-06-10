// Copyright 2019 Daniel Harrison. All Rights Reserved.

package zeropb

import (
	"encoding/binary"
	"math"
	"reflect"
	"unsafe"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

// GetBool gets an encoded bool field with the given field id.
func GetBool(buf []byte, offsets []uint16, fieldID int) bool {
	offset := int(offsets[fieldID])
	if offset == 0 {
		return false
	}
	x, _ := binary.Uvarint(buf[offset:])
	return x > 0
}

// GetInt32 gets an encoded int32 field with the given field id.
func GetInt32(buf []byte, offsets []uint16, fieldID int) int32 {
	offset := int(offsets[fieldID])
	if offset == 0 {
		return 0
	}
	x, _ := binary.Uvarint(buf[offset:])
	return int32(x)
}

// GetInt64 gets an encoded int64 field with the given field id.
func GetInt64(buf []byte, offsets []uint16, fieldID int) int64 {
	offset := int(offsets[fieldID])
	if offset == 0 {
		return 0
	}
	x, _ := binary.Uvarint(buf[offset:])
	return int64(x)
}

// GetUint32 gets an encoded uint32 field with the given field id.
func GetUint32(buf []byte, offsets []uint16, fieldID int) uint32 {
	offset := int(offsets[fieldID])
	if offset == 0 {
		return 0
	}
	x, _ := binary.Uvarint(buf[offset:])
	return uint32(x)
}

// GetUint64 gets an encoded uint64 field with the given field id.
func GetUint64(buf []byte, offsets []uint16, fieldID int) uint64 {
	offset := int(offsets[fieldID])
	if offset == 0 {
		return 0
	}
	x, _ := binary.Uvarint(buf[offset:])
	return x
}

// GetZigZagInt32 gets a zig-zag encoded int32 field with the given field id.
func GetZigZagInt32(buf []byte, offsets []uint16, fieldID int) int32 {
	offset := int(offsets[fieldID])
	if offset == 0 {
		return 0
	}
	x, _ := binary.Varint(buf[offset:])
	return int32(x)
}

// GetZigZagInt64 gets a zig-zag encoded int64 field with the given field id.
func GetZigZagInt64(buf []byte, offsets []uint16, fieldID int) int64 {
	offset := int(offsets[fieldID])
	if offset == 0 {
		return 0
	}
	x, _ := binary.Varint(buf[offset:])
	return x
}

// GetFixedUint32 gets a little-endian encoded uint32 field with the given field
// id.
func GetFixedUint32(buf []byte, offsets []uint16, fieldID int) uint32 {
	offset := int(offsets[fieldID])
	if offset == 0 {
		return 0
	}
	return binary.LittleEndian.Uint32(buf[offset : offset+4])
}

// GetFixedUint64 gets a little-endian encoded uint64 field with the given field
// id.
func GetFixedUint64(buf []byte, offsets []uint16, fieldID int) uint64 {
	offset := int(offsets[fieldID])
	if offset == 0 {
		return 0
	}
	return binary.LittleEndian.Uint64(buf[offset : offset+8])
}

// GetFixedInt32 gets a little-endian encoded int32 field with the given field
// id.
func GetFixedInt32(buf []byte, offsets []uint16, fieldID int) int32 {
	offset := int(offsets[fieldID])
	if offset == 0 {
		return 0
	}
	return int32(binary.LittleEndian.Uint32(buf[offset : offset+4]))
}

// GetFixedInt64 gets a little-endian encoded int64 field with the given field
// id.
func GetFixedInt64(buf []byte, offsets []uint16, fieldID int) int64 {
	offset := int(offsets[fieldID])
	if offset == 0 {
		return 0
	}
	return int64(binary.LittleEndian.Uint64(buf[offset : offset+8]))
}

// GetFloat32 gets a little-endian, IEEE 754 encoded float32 field with the
// given field id.
func GetFloat32(buf []byte, offsets []uint16, fieldID int) float32 {
	offset := int(offsets[fieldID])
	if offset == 0 {
		return 0
	}
	bits := binary.LittleEndian.Uint32(buf[offset : offset+4])
	return math.Float32frombits(bits)
}

// GetFloat64 gets a little-endian, IEEE 754 encoded float64 field with the
// given field id.
func GetFloat64(buf []byte, offsets []uint16, fieldID int) float64 {
	offset := int(offsets[fieldID])
	if offset == 0 {
		return 0
	}
	bits := binary.LittleEndian.Uint64(buf[offset : offset+8])
	return math.Float64frombits(bits)
}

// GetString gets an encoded string field with the given field id.
func GetString(buf []byte, offsets []uint16, fieldID int) string {
	offset := int(offsets[fieldID])
	if offset == 0 {
		return ``
	}
	len, lenSize := binary.Uvarint(buf[offset:])
	offset += lenSize
	data := buf[offset : offset+int(len)]
	return *(*string)(unsafe.Pointer(&data))
}

// GetBytes gets an encoded []byte field with the given field id.
func GetBytes(buf []byte, offsets []uint16, fieldID int) []byte {
	offset := int(offsets[fieldID])
	if offset == 0 {
		return nil
	}
	len, lenSize := binary.Uvarint(buf[offset:])
	offset += lenSize
	return buf[offset : offset+int(len)]
}

// SetBool set an encoded bool field with the given field id, overwriting an
// old value if present.
func SetBool(buf *[]byte, offsets []uint16, fieldID int, x bool) {
	var scratch [1]byte
	if x {
		scratch[0] = 1
	}
	var valLen, val []byte = nil, scratch[:]
	setSingle(buf, offsets, fieldID, proto.WireVarint, valLen, val)
}

// SetInt32 set an encoded int32 field with the given field id, overwriting an
// old value if present.
func SetInt32(buf *[]byte, offsets []uint16, fieldID int, x int32) {
	// WIP huh? var scratch [binary.MaxVarintLen32]byte
	var scratch [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(scratch[:], uint64(x))
	var valLen, val []byte = nil, scratch[:n]
	setSingle(buf, offsets, fieldID, proto.WireVarint, valLen, val)
}

// SetInt64 set an encoded int64 field with the given field id, overwriting an
// old value if present.
func SetInt64(buf *[]byte, offsets []uint16, fieldID int, x int64) {
	var scratch [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(scratch[:], uint64(x))
	var valLen, val []byte = nil, scratch[:n]
	setSingle(buf, offsets, fieldID, proto.WireVarint, valLen, val)
}

// SetUint32 set an encoded uint64 field with the given field id, overwriting an
// old value if present.
func SetUint32(buf *[]byte, offsets []uint16, fieldID int, x uint32) {
	var scratch [binary.MaxVarintLen32]byte
	n := binary.PutUvarint(scratch[:], uint64(x))
	var valLen, val []byte = nil, scratch[:n]
	setSingle(buf, offsets, fieldID, proto.WireVarint, valLen, val)
}

// SetUint64 set an encoded uint64 field with the given field id, overwriting an
// old value if present.
func SetUint64(buf *[]byte, offsets []uint16, fieldID int, x uint64) {
	var scratch [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(scratch[:], x)
	var valLen, val []byte = nil, scratch[:n]
	setSingle(buf, offsets, fieldID, proto.WireVarint, valLen, val)
}

// SetZigZagInt32 set a zig-zag encoded int64 field with the given field id,
// overwriting an old value if present.
func SetZigZagInt32(buf *[]byte, offsets []uint16, fieldID int, x int32) {
	var scratch [binary.MaxVarintLen32]byte
	n := binary.PutVarint(scratch[:], int64(x))
	var valLen, val []byte = nil, scratch[:n]
	setSingle(buf, offsets, fieldID, proto.WireVarint, valLen, val)
}

// SetZigZagInt64 set a zig-zag encoded int64 field with the given field id,
// overwriting an old value if present.
func SetZigZagInt64(buf *[]byte, offsets []uint16, fieldID int, x int64) {
	var scratch [binary.MaxVarintLen64]byte
	n := binary.PutVarint(scratch[:], x)
	var valLen, val []byte = nil, scratch[:n]
	setSingle(buf, offsets, fieldID, proto.WireVarint, valLen, val)
}

// SetFixedUint32 set a little-endian encoded uint32 field with the given field
// id, overwriting an old value if present.
func SetFixedUint32(buf *[]byte, offsets []uint16, fieldID int, x uint32) {
	var scratch [4]byte
	binary.LittleEndian.PutUint32(scratch[:], x)
	var valLen, val []byte = nil, scratch[:]
	setSingle(buf, offsets, fieldID, proto.WireFixed32, valLen, val)
}

// SetFixedUint64 set a little-endian encoded uint64 field with the given field
// id, overwriting an old value if present.
func SetFixedUint64(buf *[]byte, offsets []uint16, fieldID int, x uint64) {
	var scratch [8]byte
	binary.LittleEndian.PutUint64(scratch[:], x)
	var valLen, val []byte = nil, scratch[:]
	setSingle(buf, offsets, fieldID, proto.WireFixed64, valLen, val)
}

// SetFixedInt32 set a little-endian encoded int32 field with the given field
// id, overwriting an old value if present.
func SetFixedInt32(buf *[]byte, offsets []uint16, fieldID int, x int32) {
	var scratch [4]byte
	binary.LittleEndian.PutUint32(scratch[:], uint32(x))
	var valLen, val []byte = nil, scratch[:]
	setSingle(buf, offsets, fieldID, proto.WireFixed32, valLen, val)
}

// SetFixedInt64 set a little-endian encoded int64 field with the given field
// id, overwriting an old value if present.
func SetFixedInt64(buf *[]byte, offsets []uint16, fieldID int, x int64) {
	var scratch [8]byte
	binary.LittleEndian.PutUint64(scratch[:], uint64(x))
	var valLen, val []byte = nil, scratch[:]
	setSingle(buf, offsets, fieldID, proto.WireFixed64, valLen, val)
}

// SetFloat32 set a little-endian, IEEE 754 encoded float32 field with the given
// field id, overwriting an old value if present.
func SetFloat32(buf *[]byte, offsets []uint16, fieldID int, x float32) {
	var scratch [4]byte
	bits := math.Float32bits(x)
	binary.LittleEndian.PutUint32(scratch[:], bits)
	var valLen, val []byte = nil, scratch[:]
	setSingle(buf, offsets, fieldID, proto.WireFixed32, valLen, val)
}

// SetFloat64 set a little-endian, IEEE 754 encoded float64 field with the given
// field id, overwriting an old value if present.
func SetFloat64(buf *[]byte, offsets []uint16, fieldID int, x float64) {
	var scratch [8]byte
	bits := math.Float64bits(x)
	binary.LittleEndian.PutUint64(scratch[:], bits)
	var valLen, val []byte = nil, scratch[:]
	setSingle(buf, offsets, fieldID, proto.WireFixed64, valLen, val)
}

// SetString set an encoded string field with the given field id, overwriting an
// old value if present.
func SetString(buf *[]byte, offsets []uint16, fieldID int, x string) {
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
func SetBytes(buf *[]byte, offsets []uint16, fieldID int, x []byte) {
	var scratch [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(scratch[:], uint64(len(x)))
	var valLen, val []byte = scratch[:n], x
	setSingle(buf, offsets, fieldID, proto.WireBytes, valLen, val)
}

// AppendBytes appends one value to an encoded []byte field with the given field
// id.
func AppendBytes(buf *[]byte, offsets []uint16, fieldID int, x []byte) {
	var scratch [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(scratch[:], uint64(len(x)))
	var valLen, val []byte = scratch[:n], x
	appendSingle(buf, offsets, fieldID, proto.WireBytes, valLen, val)
}

func setSingle(buf *[]byte, offsets []uint16, fieldID, typ int, valLen, val []byte) {
	// TODO(dan): We could save space if we didn't encode fields set to the
	// default value.

	var scratch [binary.MaxVarintLen64]byte
	tag := uint64(fieldID<<3 | typ)
	tagLen := binary.PutUvarint(scratch[:], tag)

	if offset := int(offsets[fieldID]); offset != 0 {
		var oldLen int
		switch typ {
		case proto.WireVarint:
			_, oldLen = binary.Uvarint((*buf)[offset:])
		case proto.WireBytes:
			oldValLen, oldValLenLen := binary.Uvarint((*buf)[offset:])
			oldLen = int(oldValLen) + oldValLenLen
		default:
			panic(typ)
		}

		// Performance optimization: if the new value is exactly the same length as
		// the old one, simply overwrite it.
		if int(oldLen) == len(valLen)+len(val) {
			offset += copy((*buf)[offset:], valLen)
			offset += copy((*buf)[offset:], val)
			return
		}

		fieldBegin, fieldEnd := offset-tagLen, offset+oldLen
		// TODO(dan): sanity check that we got the same tag bytes here?
		copy((*buf)[fieldBegin:], (*buf)[fieldEnd:])
		removedBytes := fieldEnd - fieldBegin
		(*buf) = (*buf)[:len(*buf)-removedBytes]

		// Adjust every offset in the map that's larger than fieldEnd down by
		// removedBytes.
		for i, offset := range offsets {
			if int(offset) > fieldEnd {
				offsets[i] -= uint16(removedBytes)
			}
		}
	}
	*buf = append(*buf, scratch[:tagLen]...)

	offset := len(*buf)
	offsets[fieldID] = uint16(offset)
	*buf = append(*buf, valLen...)
	*buf = append(*buf, val...)
}

func appendSingle(buf *[]byte, offsets []uint16, fieldID, typ int, valLen, val []byte) {
	var scratch [binary.MaxVarintLen64]byte

	tag := uint64(fieldID<<3 | typ)
	tagLen := binary.PutUvarint(scratch[:], tag)
	*buf = append(*buf, scratch[:tagLen]...)

	offset := len(*buf)
	// The offset for repeated fields always points at the first one.
	if existing := int(offsets[fieldID]); existing != 0 {
		offsets[fieldID] = uint16(offset)
	}
	*buf = append(*buf, valLen...)
	*buf = append(*buf, val...)
}

// GetRepeatedNonPacked returns a slice of the internal buffer, starting at the
// first instance of a non-packed repeated field and with everything after it.
// It's positioned past the tag of the field, so the very first thing will be a
// varint length. This is exactly the input expected by FindNextField.
func GetRepeatedNonPacked(buf []byte, offsets []uint16, fieldID int) []byte {
	offset := int(offsets[fieldID])
	if offset == 0 {
		return nil
	}
	return buf[offset:]
}

// Decode parses and validates a proto message, filling in a map of field id ->
// offset as it goes.
func Decode(buf []byte, offsets []uint16) error {
	if len(buf) > math.MaxUint16 {
		return errors.Errorf(`cannot decode messages longer than %d bytes`, math.MaxUint16)
	}
	for i := range offsets {
		offsets[i] = 0
	}
	for idx := 0; idx < len(buf); {
		tag, size := binary.Uvarint(buf[idx:])
		idx += size
		field := int(tag >> 3)
		typ := tag & 0x7
		// fmt.Printf("tag=%x field=%x typ=%x\n", tag, field, typ)
		switch typ {
		case proto.WireVarint:
			offsets[field] = uint16(idx)
			_, size := binary.Uvarint(buf[idx:])
			idx += size
		case proto.WireFixed32:
			offsets[field] = uint16(idx)
			idx += 4
		case proto.WireFixed64:
			offsets[field] = uint16(idx)
			idx += 8
		case proto.WireBytes:
			// Set the offset if this is the first one we've found but don't overwrite
			// an earlier offset. The repeated message iterator will start from this
			// offset and re-parse the rest of the message to find the rest of them.
			if existing := offsets[field]; existing == 0 {
				offsets[field] = uint16(idx)
			}
			// TODO(dan): The proto spec specifically allows a single message
			// to be split, but I don't understand why it would happen. It's
			// hard to support so leave it unsupported for now.
			len, lenSize := binary.Uvarint(buf[idx:])
			idx += lenSize + int(len)
		default:
			return errors.Errorf(`unsupported type: %d`, typ)
		}
	}
	return nil
}

// FindNextField parses a bytes field (with no tag) that must be present at the
// very start of buf and finds the next requested field, returning it ready to
// be handed back to FindNextField.
func FindNextField(buf []byte, field int) ([]byte, []byte) {
	// TODO(dan): Remove the duplication between this and Decode.
	if buf == nil {
		return nil, nil
	}
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
		switch typ {
		case proto.WireVarint:
			_, size := binary.Uvarint(buf[idx:])
			idx += size
		case proto.WireFixed32:
			idx += 4
		case proto.WireFixed64:
			idx += 8
		case proto.WireBytes:
			len, lenSize := binary.Uvarint(buf[idx:])
			idx += lenSize + int(len)
		default:
			panic(errors.Errorf(`unsupported type: %d`, typ))
		}
	}
	return nil, msg
}
