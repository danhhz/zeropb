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

func randTestMessage(rng *rand.Rand, maxDepth int) *testgogopb.TestMessage {
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
		if rng.Intn(10) == 0 && depth < maxDepth {
			m.Message = fn(depth + 1)
		}
		if l := rng.Intn(5); l > 0 && depth < maxDepth {
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

	gogo := randTestMessage(rng, 3)
	zero := gogoToZeroViaBytes(gogo)
	roundtripped := zeroToGogoViaCopy(zero)
	require.Equal(t, gogo, roundtripped)
}

func TestRandEncode(t *testing.T) {
	seed := time.Now().UnixNano()
	log.Printf("seed=%d", seed)
	rng := rand.New(rand.NewSource(seed))

	gogo := randTestMessage(rng, 3)
	zero := gogoToZeroViaCopy(gogo)
	roundtripped := zeroToGogoViaBytes(zero)
	require.Equal(t, gogo, roundtripped)
}

func mutate(rng *rand.Rand, gogo *testgogopb.TestMessage, zero *testzeropb.TestMessage) {
	randMessage := func() (testgogopb.TestMessage, testzeropb.TestMessage) {
		gogoSub := randTestMessage(rng, 0)
		buf, err := proto.Marshal(gogoSub)
		if err != nil {
			panic(err)
		}
		var zeroSub testzeropb.TestMessage
		if err := zeroSub.Decode(buf); err != nil {
			panic(err)
		}
		return *gogoSub, zeroSub
	}
	_ = randMessage

	mutationFns := []func(){
		func() {
			bool := rng.Intn(2) == 0
			gogo.Bool = bool
			zero.SetBool(bool)
		},
		func() {
			int32 := int32(rng.Uint64())
			gogo.Int32 = int32
			zero.SetInt32(int32)
		},
		func() {
			int64 := int64(rng.Uint64())
			gogo.Int64 = int64
			zero.SetInt64(int64)
		},
		func() {
			uint32 := uint32(rng.Uint64())
			gogo.Uint32 = uint32
			zero.SetUint32(uint32)
		},
		func() {
			uint64 := rng.Uint64()
			gogo.Uint64 = uint64
			zero.SetUint64(uint64)
		},
		func() {
			sint32 := int32(rng.Uint64())
			gogo.Sint32 = sint32
			zero.SetSint32(sint32)
		},
		func() {
			sint64 := int64(rng.Uint64())
			gogo.Sint64 = sint64
			zero.SetSint64(sint64)
		},
		func() {
			fixed32 := uint32(rng.Uint64())
			gogo.Fixed32 = fixed32
			zero.SetFixed32(fixed32)
		},
		func() {
			fixed64 := rng.Uint64()
			gogo.Fixed64 = fixed64
			zero.SetFixed64(fixed64)
		},
		func() {
			sfixed32 := int32(rng.Uint64())
			gogo.Sfixed32 = sfixed32
			zero.SetSfixed32(sfixed32)
		},
		func() {
			sfixed64 := int64(rng.Uint64())
			gogo.Sfixed64 = sfixed64
			zero.SetSfixed64(sfixed64)
		},
		func() {
			float := rng.Float32()
			gogo.Float = float
			zero.SetFloat(float)
		},
		func() {
			double := rng.Float64()
			gogo.Double = double
			zero.SetDouble(double)
		},
		func() {
			buf := make([]byte, rng.Intn(10))
			for i := range buf {
				buf[i] = alphanum[rng.Intn(len(alphanum))]
			}
			gogo.String_ = string(buf)
			zero.SetString(string(buf))
		},
		func() {
			bytes := make([]byte, rng.Intn(10)+1)
			rng.Read(bytes)
			gogo.ByteArray = bytes
			zero.SetByteArray(bytes)
		},
		func() {
			enum := testgogopb.TestEnum(rng.Intn(2))
			gogo.Enum = enum
			zero.SetEnum(uint32(enum))
		},
		func() {
			gogoSub, zeroSub := randMessage()
			gogo.Message = &gogoSub
			zero.SetMessage(zeroSub)
		},
		func() {
			gogoSub, zeroSub := randMessage()
			gogo.Messages = append(gogo.Messages, &gogoSub)
			zero.AppendToMessages(zeroSub)
		},
	}
	mutationFns[rng.Intn(len(mutationFns))]()
}

func TestRandMutations(t *testing.T) {
	const iterations = 1000
	seed := time.Now().UnixNano()
	seed = 1560546759782267000
	log.Printf("seed=%d", seed)
	rng := rand.New(rand.NewSource(seed))

	var gogo testgogopb.TestMessage
	var zero testzeropb.TestMessage

	defer func() {
		if err := recover(); err != nil && err != `cannot create a proto this big` {
			t.Fatal(err)
		}
	}()
	for i := 0; i < iterations; i++ {
		mutate(rng, &gogo, &zero)

		viaBytes := zeroToGogoViaBytes(zero)
		require.Equal(t, gogo, *viaBytes)
		viaCopy := zeroToGogoViaCopy(zero)
		require.Equal(t, gogo, *viaCopy)

		zeroViaBytes := gogoToZeroViaBytes(&gogo)
		require.Equal(t, &gogo, zeroToGogoViaCopy(zeroViaBytes))
		zeroViaCopy := gogoToZeroViaCopy(&gogo)
		require.Equal(t, &gogo, zeroToGogoViaBytes(zeroViaCopy))
	}
}
