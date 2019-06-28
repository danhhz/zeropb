// Copyright 2019 Daniel Harrison. All Rights Reserved.

package testzeropb

// Offsets allows tests to examine the internal offsets.
func (m *TestMessage) Offsets() ([]uint16, *map[int]uint64) {
	return m.offsets.a[:], &m.offsets.m
}
