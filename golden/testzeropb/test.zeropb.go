// Code generated by protoc-gen-zeropb. DO NOT EDIT.

package testzeropb

import "github.com/danhhz/zeropb"

type TestMessage struct {
  buf []byte
  offsets struct {
    a [35]uint16
    m map[int]uint64
  }
}

var repeatedFields_TestMessage = zeropb.RepeatedFields{
  18: struct{}{},
  19: struct{}{},
  20: struct{}{},
  21: struct{}{},
  22: struct{}{},
  23: struct{}{},
  24: struct{}{},
  25: struct{}{},
  26: struct{}{},
  27: struct{}{},
  28: struct{}{},
  29: struct{}{},
  30: struct{}{},
  31: struct{}{},
  32: struct{}{},
  33: struct{}{},
  34: struct{}{},
}

func (m *TestMessage) Encode() []byte {
  return m.buf
}

func (m *TestMessage) Decode(buf []byte) error {
  m.buf = buf
  return zeropb.Decode(m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), repeatedFields_TestMessage)
}

func (m *TestMessage) Reset(buf []byte) {
  if len(buf) > 0 {
    panic(`buf must be empty`)
  }
  m.buf = buf
  zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m).Clear()
}

func (m *TestMessage) Bool() bool {
  return zeropb.GetBool(m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 1)
}

func (m *TestMessage) SetBool(x bool) {
  zeropb.SetBool(&m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 1, x)
}

func (m *TestMessage) Int32() int32 {
  return zeropb.GetInt32(m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 2)
}

func (m *TestMessage) SetInt32(x int32) {
  zeropb.SetInt32(&m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 2, x)
}

func (m *TestMessage) Int64() int64 {
  return zeropb.GetInt64(m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 3)
}

func (m *TestMessage) SetInt64(x int64) {
  zeropb.SetInt64(&m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 3, x)
}

func (m *TestMessage) Uint32() uint32 {
  return zeropb.GetUint32(m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 4)
}

func (m *TestMessage) SetUint32(x uint32) {
  zeropb.SetUint32(&m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 4, x)
}

func (m *TestMessage) Uint64() uint64 {
  return zeropb.GetUint64(m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 5)
}

func (m *TestMessage) SetUint64(x uint64) {
  zeropb.SetUint64(&m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 5, x)
}

func (m *TestMessage) Sint32() int32 {
  return zeropb.GetZigZagInt32(m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 6)
}

func (m *TestMessage) SetSint32(x int32) {
  zeropb.SetZigZagInt32(&m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 6, x)
}

func (m *TestMessage) Sint64() int64 {
  return zeropb.GetZigZagInt64(m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 7)
}

func (m *TestMessage) SetSint64(x int64) {
  zeropb.SetZigZagInt64(&m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 7, x)
}

func (m *TestMessage) Fixed32() uint32 {
  return zeropb.GetFixedUint32(m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 8)
}

func (m *TestMessage) SetFixed32(x uint32) {
  zeropb.SetFixedUint32(&m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 8, x)
}

func (m *TestMessage) Fixed64() uint64 {
  return zeropb.GetFixedUint64(m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 9)
}

func (m *TestMessage) SetFixed64(x uint64) {
  zeropb.SetFixedUint64(&m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 9, x)
}

func (m *TestMessage) Sfixed32() int32 {
  return zeropb.GetFixedInt32(m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 10)
}

func (m *TestMessage) SetSfixed32(x int32) {
  zeropb.SetFixedInt32(&m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 10, x)
}

func (m *TestMessage) Sfixed64() int64 {
  return zeropb.GetFixedInt64(m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 11)
}

func (m *TestMessage) SetSfixed64(x int64) {
  zeropb.SetFixedInt64(&m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 11, x)
}

func (m *TestMessage) Double() float64 {
  return zeropb.GetFloat64(m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 12)
}

func (m *TestMessage) SetDouble(x float64) {
  zeropb.SetFloat64(&m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 12, x)
}

func (m *TestMessage) Float() float32 {
  return zeropb.GetFloat32(m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 13)
}

func (m *TestMessage) SetFloat(x float32) {
  zeropb.SetFloat32(&m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 13, x)
}

func (m *TestMessage) String() string {
  return zeropb.GetString(m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 14)
}

func (m *TestMessage) SetString(x string) {
  zeropb.SetString(&m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 14, x)
}

func (m *TestMessage) ByteArray() []byte {
  return zeropb.GetBytes(m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 15)
}

func (m *TestMessage) SetByteArray(x []byte) {
  zeropb.SetBytes(&m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 15, x)
}

func (m *TestMessage) Enum() uint32 {
  return zeropb.GetUint32(m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 16)
}

func (m *TestMessage) SetEnum(x uint32) {
  zeropb.SetUint32(&m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 16, x)
}

func (m *TestMessage) Message(x *TestMessage) (bool, error) {
  buf := zeropb.GetBytes(m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 17)
  if buf == nil {
    return false, nil
  }
  return true, x.Decode(buf)
}

func (m *TestMessage) SetMessage(x TestMessage) {
  buf := x.Encode()
  zeropb.SetBytes(&m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 17, buf)
}

type TestMessageTestMessageIterator []byte

func (i *TestMessageTestMessageIterator) Next(m *TestMessage) (bool, error) {
  var buf []byte
  *i, buf = zeropb.FindNextField((*i), 34)
  if buf == nil {
    return false, nil
  }
  return true, m.Decode(buf)
}

func (m *TestMessage) Messages() TestMessageTestMessageIterator {
  return TestMessageTestMessageIterator(zeropb.GetRepeatedNonPacked(m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 34))
}

func (m *TestMessage) AppendToMessages(x TestMessage) {
  buf := x.Encode()
  zeropb.AppendBytes(&m.buf, zeropb.WrapOffsets(m.offsets.a[:], &m.offsets.m), 34, buf)
}

