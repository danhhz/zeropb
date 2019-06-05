// Copyright 2019 Daniel Harrison. All Rights Reserved.

package zeropb

import (
	"encoding/binary"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

// GetUint32 gets an encoded uint32 field with the given field id.
func GetUint32(buf []byte, offsets *FastIntMap, fieldID int) uint32 {
	offset, ok := offsets.Get(fieldID)
	if !ok {
		return 0
	}
	x, _ := binary.Uvarint(buf[offset:])
	return uint32(x)
}

// GetUint64 gets an encoded uint64 field with the given field id.
func GetUint64(buf []byte, offsets *FastIntMap, fieldID int) uint64 {
	offset, ok := offsets.Get(fieldID)
	if !ok {
		return 0
	}
	x, _ := binary.Uvarint(buf[offset:])
	return uint64(x)
}

// GetBytes gets an encoded []byte field with the given field id.
func GetBytes(buf []byte, offsets *FastIntMap, fieldID int) []byte {
	offset, ok := offsets.Get(fieldID)
	if !ok {
		return nil
	}
	len, lenSize := binary.Uvarint(buf[offset:])
	offset += lenSize
	return buf[offset : offset+int(len)]
}

// SetUint32 set an encoded uint64 field with the given field id, overwriting an
// old value if present.
func SetUint32(buf *[]byte, offsets *FastIntMap, fieldID int, x uint32) {
	var scratch [binary.MaxVarintLen32]byte
	n := binary.PutUvarint(scratch[:], uint64(x))
	var valLen, val []byte = nil, scratch[:n]
	setSingle(buf, offsets, fieldID, proto.WireVarint, valLen, val)
}

// SetUint64 set an encoded uint64 field with the given field id, overwriting an
// old value if present.
func SetUint64(buf *[]byte, offsets *FastIntMap, fieldID int, x uint64) {
	var scratch [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(scratch[:], x)
	var valLen, val []byte = nil, scratch[:n]
	setSingle(buf, offsets, fieldID, proto.WireVarint, valLen, val)
}

// SetBytes set an encoded []byte field with the given field id, overwriting an
// old value if present.
func SetBytes(buf *[]byte, offsets *FastIntMap, fieldID int, x []byte) {
	var scratch [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(scratch[:], uint64(len(x)))
	var valLen, val []byte = scratch[:n], x
	setSingle(buf, offsets, fieldID, proto.WireBytes, valLen, val)
}

// AppendBytes appends one value to an encoded []byte field with the given field
// id.
func AppendBytes(buf *[]byte, offsets *FastIntMap, fieldID int, x []byte) {
	var scratch [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(scratch[:], uint64(len(x)))
	var valLen, val []byte = scratch[:n], x
	appendSingle(buf, offsets, fieldID, proto.WireBytes, valLen, val)
}

func setSingle(buf *[]byte, offsets *FastIntMap, fieldID, typ int, valLen, val []byte) {
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
		// removedBytes. In general, setting a FastIntMap while iterating it is bad
		// news, but we happen to know that setting a smaller (still positive) value
		// for an entry already in the map will not cause it to flap from small to
		// large. This is still a little scary though.
		offsets.ForEach(func(fieldID int, offset int) {
			if offset > fieldEnd {
				offsets.Set(fieldID, offset-removedBytes)
			}
		})
	}
	*buf = append(*buf, scratch[:tagLen]...)

	offset := len(*buf)
	offsets.Set(fieldID, offset)
	*buf = append(*buf, valLen...)
	*buf = append(*buf, val...)
}

func appendSingle(buf *[]byte, offsets *FastIntMap, fieldID, typ int, valLen, val []byte) {
	var scratch [binary.MaxVarintLen64]byte

	tag := uint64(fieldID<<3 | typ)
	tagLen := binary.PutUvarint(scratch[:], tag)
	*buf = append(*buf, scratch[:tagLen]...)

	offset := len(*buf)
	// The offset for repeated fields always points at the first one.
	if _, ok := offsets.Get(fieldID); !ok {
		offsets.Set(fieldID, offset)
	}
	*buf = append(*buf, valLen...)
	*buf = append(*buf, val...)
}

// GetRepeatedMessage WIP
func GetRepeatedMessage(buf []byte, offsets *FastIntMap, fieldID int) []byte {
	offset, ok := offsets.Get(fieldID)
	if !ok {
		return nil
	}
	return buf[offset:]
}

// Decode parses and validates a proto message, filling in a map of field id ->
// offset as it goes.
func Decode(buf []byte, offsets *FastIntMap) error {
	offsets.Clear()
	for idx := 0; idx < len(buf); {
		tag, size := binary.Uvarint(buf[idx:])
		idx += size
		field := int(tag >> 3)
		typ := tag & 0x7
		// fmt.Printf("tag=%x field=%x typ=%x\n", tag, field, typ)
		switch typ {
		case proto.WireVarint:
			offsets.Set(field, idx)
			_, size := binary.Uvarint(buf[idx:])
			idx += size
		case proto.WireBytes:
			// Set the offset if this is the first one we've found but don't overwrite
			// an earlier offset. The repeated message iterator will start from this
			// offset and re-parse the rest of the message to find the rest of them.
			if _, ok := offsets.Get(field); !ok {
				offsets.Set(field, idx)
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
		case 0:
			_, size := binary.Uvarint(buf[idx:])
			idx += size
		case 2:
			len, lenSize := binary.Uvarint(buf[idx:])
			idx += lenSize + int(len)
		default:
			panic(errors.Errorf(`unsupported type: %d`, typ))
		}
	}
	return nil, msg
}
