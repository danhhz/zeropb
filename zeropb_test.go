// Copyright 2019 Daniel Harrison. All Rights Reserved.

package zeropb_test

import (
	"testing"

	"github.com/danhhz/zeropb/golden/raftzeropb"
	"github.com/stretchr/testify/require"
)

func TestEncode(t *testing.T) {
	var e raftzeropb.Entry
	require.Equal(t, uint64(0), e.Index())
	require.Equal(t, uint64(0), e.Term())
	require.Equal(t, []byte(nil), e.Data())

	// Set previously unset fields.
	e.SetIndex(1)
	e.SetTerm(2)
	e.SetData([]byte{3, 4})
	require.Equal(t, uint64(1), e.Index())
	require.Equal(t, uint64(2), e.Term())
	require.Equal(t, []byte{3, 4}, e.Data())

	// Overwrite previously set fields with same sized data (fast path).
	e.SetIndex(5)
	require.Equal(t, uint64(5), e.Index())
	e.SetData([]byte{7, 8})
	require.Equal(t, []byte{7, 8}, e.Data())

	// Overwrite previously set data with different sized data.
	e.SetData([]byte{9, 10, 11})
	require.Equal(t, []byte{9, 10, 11}, e.Data())
	e.SetIndex(1000000)
	require.Equal(t, uint64(1000000), e.Index())

	// The previous updates had to remove some data in the middle of the encoding,
	// which then requires updating the offsets map. Double check that everything
	// is still what we expect it to be.
	require.Equal(t, uint64(1000000), e.Index())
	require.Equal(t, uint64(2), e.Term())
	require.Equal(t, []byte{9, 10, 11}, e.Data())
}
