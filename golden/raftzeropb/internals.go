// Copyright 2019 Daniel Harrison. All Rights Reserved.

package raftzeropb

// Offsets allows tests to examine the internal offsets.
func (m *Entry) Offsets() ([]uint16, *map[int]uint64) {
	return m.offsets.a[:], &m.offsets.m
}

// Offsets allows tests to examine the internal offsets.
func (m *Message) Offsets() ([]uint16, *map[int]uint64) {
	return m.offsets.a[:], &m.offsets.m
}

// Offsets allows tests to examine the internal offsets.
func (m *SnapshotMetadata) Offsets() ([]uint16, *map[int]uint64) {
	return m.offsets.a[:], &m.offsets.m
}

// Offsets allows tests to examine the internal offsets.
func (m *Snapshot) Offsets() ([]uint16, *map[int]uint64) {
	return m.offsets.a[:], &m.offsets.m
}
