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
