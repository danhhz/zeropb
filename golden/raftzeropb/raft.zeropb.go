// Code generated by protoc-gen-zeropb. DO NOT EDIT.

package raftzeropb

import "github.com/danhhz/zeropb"

type Entry struct {
  buf []byte
  offsets [5]uint16
}

func (m *Entry) Encode() []byte {
  return m.buf
}

func (m *Entry) Decode(buf []byte) error {
  m.buf = buf
  return zeropb.Decode(m.buf, m.offsets[:])
}

func (m *Entry) Reset(buf []byte) {
  if len(buf) > 0 {
    panic(`buf must be empty`)
  }
  m.buf = buf
  for i := range m.offsets {
    m.offsets[i] = 0
  }
}

func (m *Entry) Term() uint64 {
  return zeropb.GetUint64(m.buf, m.offsets[:], 2)
}

func (m *Entry) SetTerm(x uint64) {
  zeropb.SetUint64(&m.buf, m.offsets[:], 2, x)
}

func (m *Entry) Index() uint64 {
  return zeropb.GetUint64(m.buf, m.offsets[:], 3)
}

func (m *Entry) SetIndex(x uint64) {
  zeropb.SetUint64(&m.buf, m.offsets[:], 3, x)
}

func (m *Entry) Type() uint32 {
  return zeropb.GetUint32(m.buf, m.offsets[:], 1)
}

func (m *Entry) SetType(x uint32) {
  zeropb.SetUint32(&m.buf, m.offsets[:], 1, x)
}

func (m *Entry) Data() []byte {
  return zeropb.GetBytes(m.buf, m.offsets[:], 4)
}

func (m *Entry) SetData(x []byte) {
  zeropb.SetBytes(&m.buf, m.offsets[:], 4, x)
}

type SnapshotMetadata struct {
  buf []byte
  offsets [4]uint16
}

func (m *SnapshotMetadata) Encode() []byte {
  return m.buf
}

func (m *SnapshotMetadata) Decode(buf []byte) error {
  m.buf = buf
  return zeropb.Decode(m.buf, m.offsets[:])
}

func (m *SnapshotMetadata) Reset(buf []byte) {
  if len(buf) > 0 {
    panic(`buf must be empty`)
  }
  m.buf = buf
  for i := range m.offsets {
    m.offsets[i] = 0
  }
}

func (m *SnapshotMetadata) ConfState(x *ConfState) (bool, error) {
  buf := zeropb.GetBytes(m.buf, m.offsets[:], 1)
  if buf == nil {
    return false, nil
  }
  return true, x.Decode(buf)
}

func (m *SnapshotMetadata) SetConfState(x ConfState) {
  buf := x.Encode()
  zeropb.SetBytes(&m.buf, m.offsets[:], 1, buf)
}

func (m *SnapshotMetadata) Index() uint64 {
  return zeropb.GetUint64(m.buf, m.offsets[:], 2)
}

func (m *SnapshotMetadata) SetIndex(x uint64) {
  zeropb.SetUint64(&m.buf, m.offsets[:], 2, x)
}

func (m *SnapshotMetadata) Term() uint64 {
  return zeropb.GetUint64(m.buf, m.offsets[:], 3)
}

func (m *SnapshotMetadata) SetTerm(x uint64) {
  zeropb.SetUint64(&m.buf, m.offsets[:], 3, x)
}

type Snapshot struct {
  buf []byte
  offsets [3]uint16
}

func (m *Snapshot) Encode() []byte {
  return m.buf
}

func (m *Snapshot) Decode(buf []byte) error {
  m.buf = buf
  return zeropb.Decode(m.buf, m.offsets[:])
}

func (m *Snapshot) Reset(buf []byte) {
  if len(buf) > 0 {
    panic(`buf must be empty`)
  }
  m.buf = buf
  for i := range m.offsets {
    m.offsets[i] = 0
  }
}

func (m *Snapshot) Data() []byte {
  return zeropb.GetBytes(m.buf, m.offsets[:], 1)
}

func (m *Snapshot) SetData(x []byte) {
  zeropb.SetBytes(&m.buf, m.offsets[:], 1, x)
}

func (m *Snapshot) Metadata(x *SnapshotMetadata) (bool, error) {
  buf := zeropb.GetBytes(m.buf, m.offsets[:], 2)
  if buf == nil {
    return false, nil
  }
  return true, x.Decode(buf)
}

func (m *Snapshot) SetMetadata(x SnapshotMetadata) {
  buf := x.Encode()
  zeropb.SetBytes(&m.buf, m.offsets[:], 2, buf)
}

type Message struct {
  buf []byte
  offsets [13]uint16
}

func (m *Message) Encode() []byte {
  return m.buf
}

func (m *Message) Decode(buf []byte) error {
  m.buf = buf
  return zeropb.Decode(m.buf, m.offsets[:])
}

func (m *Message) Reset(buf []byte) {
  if len(buf) > 0 {
    panic(`buf must be empty`)
  }
  m.buf = buf
  for i := range m.offsets {
    m.offsets[i] = 0
  }
}

func (m *Message) Type() uint32 {
  return zeropb.GetUint32(m.buf, m.offsets[:], 1)
}

func (m *Message) SetType(x uint32) {
  zeropb.SetUint32(&m.buf, m.offsets[:], 1, x)
}

func (m *Message) To() uint64 {
  return zeropb.GetUint64(m.buf, m.offsets[:], 2)
}

func (m *Message) SetTo(x uint64) {
  zeropb.SetUint64(&m.buf, m.offsets[:], 2, x)
}

func (m *Message) From() uint64 {
  return zeropb.GetUint64(m.buf, m.offsets[:], 3)
}

func (m *Message) SetFrom(x uint64) {
  zeropb.SetUint64(&m.buf, m.offsets[:], 3, x)
}

func (m *Message) Term() uint64 {
  return zeropb.GetUint64(m.buf, m.offsets[:], 4)
}

func (m *Message) SetTerm(x uint64) {
  zeropb.SetUint64(&m.buf, m.offsets[:], 4, x)
}

func (m *Message) LogTerm() uint64 {
  return zeropb.GetUint64(m.buf, m.offsets[:], 5)
}

func (m *Message) SetLogTerm(x uint64) {
  zeropb.SetUint64(&m.buf, m.offsets[:], 5, x)
}

func (m *Message) Index() uint64 {
  return zeropb.GetUint64(m.buf, m.offsets[:], 6)
}

func (m *Message) SetIndex(x uint64) {
  zeropb.SetUint64(&m.buf, m.offsets[:], 6, x)
}

type MessageEntryIterator []byte

func (i *MessageEntryIterator) Next(m *Entry) (bool, error) {
  var buf []byte
  *i, buf = zeropb.FindNextField((*i), 7)
  if buf == nil {
    return false, nil
  }
  return true, m.Decode(buf)
}

func (m *Message) Entries() MessageEntryIterator {
  return MessageEntryIterator(zeropb.GetRepeatedNonPacked(m.buf, m.offsets[:], 7))
}

func (m *Message) AppendToEntries(x Entry) {
  buf := x.Encode()
  zeropb.AppendBytes(&m.buf, m.offsets[:], 7, buf)
}

func (m *Message) Commit() uint64 {
  return zeropb.GetUint64(m.buf, m.offsets[:], 8)
}

func (m *Message) SetCommit(x uint64) {
  zeropb.SetUint64(&m.buf, m.offsets[:], 8, x)
}

func (m *Message) Snapshot(x *Snapshot) (bool, error) {
  buf := zeropb.GetBytes(m.buf, m.offsets[:], 9)
  if buf == nil {
    return false, nil
  }
  return true, x.Decode(buf)
}

func (m *Message) SetSnapshot(x Snapshot) {
  buf := x.Encode()
  zeropb.SetBytes(&m.buf, m.offsets[:], 9, buf)
}

func (m *Message) Reject() bool {
  return zeropb.GetBool(m.buf, m.offsets[:], 10)
}

func (m *Message) SetReject(x bool) {
  zeropb.SetBool(&m.buf, m.offsets[:], 10, x)
}

func (m *Message) RejectHint() uint64 {
  return zeropb.GetUint64(m.buf, m.offsets[:], 11)
}

func (m *Message) SetRejectHint(x uint64) {
  zeropb.SetUint64(&m.buf, m.offsets[:], 11, x)
}

func (m *Message) Context() []byte {
  return zeropb.GetBytes(m.buf, m.offsets[:], 12)
}

func (m *Message) SetContext(x []byte) {
  zeropb.SetBytes(&m.buf, m.offsets[:], 12, x)
}

type HardState struct {
  buf []byte
  offsets [4]uint16
}

func (m *HardState) Encode() []byte {
  return m.buf
}

func (m *HardState) Decode(buf []byte) error {
  m.buf = buf
  return zeropb.Decode(m.buf, m.offsets[:])
}

func (m *HardState) Reset(buf []byte) {
  if len(buf) > 0 {
    panic(`buf must be empty`)
  }
  m.buf = buf
  for i := range m.offsets {
    m.offsets[i] = 0
  }
}

func (m *HardState) Term() uint64 {
  return zeropb.GetUint64(m.buf, m.offsets[:], 1)
}

func (m *HardState) SetTerm(x uint64) {
  zeropb.SetUint64(&m.buf, m.offsets[:], 1, x)
}

func (m *HardState) Vote() uint64 {
  return zeropb.GetUint64(m.buf, m.offsets[:], 2)
}

func (m *HardState) SetVote(x uint64) {
  zeropb.SetUint64(&m.buf, m.offsets[:], 2, x)
}

func (m *HardState) Commit() uint64 {
  return zeropb.GetUint64(m.buf, m.offsets[:], 3)
}

func (m *HardState) SetCommit(x uint64) {
  zeropb.SetUint64(&m.buf, m.offsets[:], 3, x)
}

type ConfState struct {
  buf []byte
  offsets [3]uint16
}

func (m *ConfState) Encode() []byte {
  return m.buf
}

func (m *ConfState) Decode(buf []byte) error {
  m.buf = buf
  return zeropb.Decode(m.buf, m.offsets[:])
}

func (m *ConfState) Reset(buf []byte) {
  if len(buf) > 0 {
    panic(`buf must be empty`)
  }
  m.buf = buf
  for i := range m.offsets {
    m.offsets[i] = 0
  }
}

type ConfChange struct {
  buf []byte
  offsets [5]uint16
}

func (m *ConfChange) Encode() []byte {
  return m.buf
}

func (m *ConfChange) Decode(buf []byte) error {
  m.buf = buf
  return zeropb.Decode(m.buf, m.offsets[:])
}

func (m *ConfChange) Reset(buf []byte) {
  if len(buf) > 0 {
    panic(`buf must be empty`)
  }
  m.buf = buf
  for i := range m.offsets {
    m.offsets[i] = 0
  }
}

func (m *ConfChange) ID() uint64 {
  return zeropb.GetUint64(m.buf, m.offsets[:], 1)
}

func (m *ConfChange) SetID(x uint64) {
  zeropb.SetUint64(&m.buf, m.offsets[:], 1, x)
}

func (m *ConfChange) Type() uint32 {
  return zeropb.GetUint32(m.buf, m.offsets[:], 2)
}

func (m *ConfChange) SetType(x uint32) {
  zeropb.SetUint32(&m.buf, m.offsets[:], 2, x)
}

func (m *ConfChange) NodeID() uint64 {
  return zeropb.GetUint64(m.buf, m.offsets[:], 3)
}

func (m *ConfChange) SetNodeID(x uint64) {
  zeropb.SetUint64(&m.buf, m.offsets[:], 3, x)
}

func (m *ConfChange) Context() []byte {
  return zeropb.GetBytes(m.buf, m.offsets[:], 4)
}

func (m *ConfChange) SetContext(x []byte) {
  zeropb.SetBytes(&m.buf, m.offsets[:], 4, x)
}

