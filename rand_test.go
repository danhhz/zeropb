// Copyright 2019 Daniel Harrison. All Rights Reserved.

package zeropb_test

import (
	"log"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/danhhz/zeropb/golden/testgogopb"
	"github.com/danhhz/zeropb/golden/testzeropb"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
)

const alphanum = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`

func randTestMessage(rng *rand.Rand) *testgogopb.TestMessage {
	var fn func(depth int) *testgogopb.TestMessage
	fn = func(depth int) *testgogopb.TestMessage {
		m := &testgogopb.TestMessage{
			Bool:     rng.Intn(2) == 0,
			Int32:    int32(rng.Uint32()),
			Int64:    int64(rng.Uint64()),
			Uint32:   rng.Uint32(),
			Uint64:   rng.Uint64(),
			Sint32:   int32(rng.Uint32()),
			Sint64:   int64(rng.Uint64()),
			Fixed32:  rng.Uint32(),
			Fixed64:  rng.Uint64(),
			Sfixed32: int32(rng.Uint32()),
			Sfixed64: int64(rng.Uint64()),
			Float:    rng.Float32(),
			Double:   rng.Float64(),
			Enum:     testgogopb.TestEnum(rng.Intn(2)),
		}
		if l := rng.Intn(10); l > 0 {
			var buf strings.Builder
			for i := 0; i < l; i++ {
				buf.WriteByte(alphanum[rng.Intn(len(alphanum))])
			}
			m.String_ = buf.String()
		}
		if l := rng.Intn(10); l > 0 {
			m.ByteArray = make([]byte, l)
			for i := range m.ByteArray {
				m.ByteArray[i] = byte(rng.Intn(256))
			}
		}
		if rng.Intn(10) == 0 && depth < 3 {
			m.Message = fn(depth + 1)
		}
		if l := rng.Intn(5); l > 0 && depth < 3 {
			m.Messages = make([]*testgogopb.TestMessage, l)
			for i := range m.Messages {
				m.Messages[i] = fn(depth + 1)
			}
		}
		return m
	}
	return fn(0)
}

func gogoToZeroViaBytes(gogo *testgogopb.TestMessage) testzeropb.TestMessage {
	buf, err := proto.Marshal(gogo)
	if err != nil {
		panic(err)
	}
	var zero testzeropb.TestMessage
	if err := zero.Decode(buf); err != nil {
		panic(err)
	}
	return zero
}

func zeroToGogoViaBytes(zero testzeropb.TestMessage) *testgogopb.TestMessage {
	var gogo testgogopb.TestMessage
	buf := zero.Encode()
	if err := proto.Unmarshal(buf, &gogo); err != nil {
		panic(err)
	}
	return &gogo
}

func gogoToZeroViaCopy(gogo *testgogopb.TestMessage) testzeropb.TestMessage {
	var zero testzeropb.TestMessage
	zero.SetBool(gogo.Bool)
	zero.SetInt32(gogo.Int32)
	zero.SetInt64(gogo.Int64)
	zero.SetUint32(gogo.Uint32)
	zero.SetUint64(gogo.Uint64)
	zero.SetSint32(gogo.Sint32)
	zero.SetSint64(gogo.Sint64)
	zero.SetFixed32(gogo.Fixed32)
	zero.SetFixed64(gogo.Fixed64)
	zero.SetSfixed32(gogo.Sfixed32)
	zero.SetSfixed64(gogo.Sfixed64)
	zero.SetFloat(gogo.Float)
	zero.SetDouble(gogo.Double)
	zero.SetString(gogo.String_)
	zero.SetEnum(uint32(gogo.Enum))
	if x := gogo.ByteArray; len(x) > 0 {
		zero.SetByteArray(gogo.ByteArray)
	}
	if gogo.Message != nil {
		sub := gogoToZeroViaCopy(gogo.Message)
		zero.SetMessage(sub)
	}
	for _, m := range gogo.Messages {
		sub := gogoToZeroViaCopy(m)
		zero.AppendToMessages(sub)
	}
	return zero
}

func zeroToGogoViaCopy(zero testzeropb.TestMessage) *testgogopb.TestMessage {
	gogo := &testgogopb.TestMessage{
		Bool:      zero.Bool(),
		Int32:     zero.Int32(),
		Int64:     zero.Int64(),
		Uint32:    zero.Uint32(),
		Uint64:    zero.Uint64(),
		Sint32:    zero.Sint32(),
		Sint64:    zero.Sint64(),
		Fixed32:   zero.Fixed32(),
		Fixed64:   zero.Fixed64(),
		Sfixed32:  zero.Sfixed32(),
		Sfixed64:  zero.Sfixed64(),
		Double:    zero.Double(),
		Float:     zero.Float(),
		String_:   zero.String(),
		ByteArray: zero.ByteArray(),
		Enum:      testgogopb.TestEnum(zero.Enum()),
	}
	var sub testzeropb.TestMessage
	if ok, err := zero.Message(&sub); err != nil {
		panic(err)
	} else if ok {
		gogo.Message = zeroToGogoViaCopy(sub)
	}
	for it := zero.Messages(); ; {
		if ok, err := it.Next(&sub); err != nil {
			panic(err)
		} else if !ok {
			break
		}
		gogo.Messages = append(gogo.Messages, zeroToGogoViaCopy(sub))
	}
	return gogo
}

func TestRandDecode(t *testing.T) {
	seed := time.Now().UnixNano()
	log.Printf("seed=%d", seed)
	rng := rand.New(rand.NewSource(seed))

	gogo := randTestMessage(rng)
	zero := gogoToZeroViaBytes(gogo)
	roundtripped := zeroToGogoViaCopy(zero)
	require.Equal(t, gogo, roundtripped)
}

func TestRandEncode(t *testing.T) {
	seed := time.Now().UnixNano()
	log.Printf("seed=%d", seed)
	rng := rand.New(rand.NewSource(seed))

	gogo := randTestMessage(rng)
	zero := gogoToZeroViaCopy(gogo)
	roundtripped := zeroToGogoViaBytes(zero)
	require.Equal(t, gogo, roundtripped)
}
