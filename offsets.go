// Copyright 2019 Daniel Harrison. All Rights Reserved.

package zeropb

import (
	"math"
)

const (
	// offsetUnset indicates that the offset is not set.
	offsetUnset uint16 = 0
	// offsetInMap indicates that the offset is stored in the map.
	offsetInMap uint16 = math.MaxUint16
	// offsetMaxInArray is the maximum offset storable in the array.
	offsetMaxInArray uint16 = offsetInMap - 1
)

// Offsets is a mapping between field id and byte offset into the message. The
// offset always points after the varint tag containing the field id and field
// type.
//
// Offsets smaller than OffsetMaxInArray are stored in the array. Otherwise, the
// map is lazily allocated once and the offset is stored in there. Invariant:
// the fieldID is present in the map if and only iff `o.a[fieldID] ==
// OffsetInMap`. If the field is unset, `o.a[fieldID]` will be OffsetUnset.
//
// - Non-repeated fields point at the last instance. The spec says to ignore
//   everything but the last one.
// - Repeated fields point at the first instance. Iterating one of these fields
//   will partially re-parse the message starting at this first instance.
// - All offsets point to after the varint tag containing the field id and type.
//   Primitive fields (varint, fixed32, fixed64) point directly at the field
//   data. Bytes fields point at the varint length. Packed repeated fields are
//   not yet supported.
type Offsets struct {
	a []uint16
	m *map[int]uint64
}

// WrapOffsets wraps the array and map fields from a generated struct into
// an Offsets helper.
func WrapOffsets(a []uint16, m *map[int]uint64) Offsets {
	return Offsets{a: a, m: m}
}

// Get returns the offset of the given field or false if it's not present.
func (o Offsets) Get(fieldID int) (uint64, bool) {
	if offset := o.a[fieldID]; offset != offsetInMap {
		return uint64(offset), offset != offsetUnset
	}
	// If the array had the sentinel, the map is guaranteed to exist and the
	// value is guaranteed to be set.
	return (*o.m)[fieldID], true
}

// Set stores the given field id offset.
func (o Offsets) Set(fieldID int, offset uint64) {
	if offset <= uint64(offsetMaxInArray) {
		if previous := o.a[fieldID]; previous == offsetInMap {
			delete(*o.m, fieldID)
		}
		o.a[fieldID] = uint16(offset)
	} else {
		// No need to check what was previously in the array, just
		// unconditionally overwrite it.
		o.a[fieldID] = offsetInMap
		if *o.m == nil {
			// TODO(dan): Preallocate this with a capacity of the number of
			// field ids to make sure it's only ever 1 allocation.
			*o.m = make(map[int]uint64, fieldID)
		}
		(*o.m)[fieldID] = offset
	}
}

// Clear removes all entries.
func (o Offsets) Clear() {
	for i := range o.a {
		o.a[i] = offsetUnset
	}
	if *o.m != nil {
		for f := range *o.m {
			delete(*o.m, f)
		}
	}
}
