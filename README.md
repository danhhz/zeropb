zeropb
======

### This is a proof of concept and _not ready for production_ use.

zeropb is a [Protocol Buffer] runtime and code generator for Go. The primary
goal is "zero" heap allocations at runtime, including encoding, decoding, field
access and field mutation.

Protocol Buffers (protobufs) are a language-independent data interchange format,
designed specifically to accommodate schemas that evolve over time. A protobuf
user defines schemas (messages) in a `.proto` file and then uses the `protoc`
compiler to codegen a language-native representation of each message. zeropb
uses the protoc plugin system to do the code generation.


## Compromises

In contrast to [FlatBuffers] and [Cap’n Proto], the protobuf wire format was not
specifically designed for allocation-less use. This means that zeropb ends up as
somewhat more awkward to use than other protobuf variants for Go, such as the
official Google one or gogoprotobuf.

It also means that it is not always possible to decode without heap allocations.
This list will become _considerably_ less restrictive over time, but the current
set of requirements for allocation-less use are as follows:

- All field ids are less than 32.
- For decoding, the largest offset of the first appearance of a field in an
  encoded message is less than 15. That is, an encoded message is a
  concatenation of field id/value pairs. Field ids may be repeated. Find the
  offset of the first appearance of each field id in the encoded message. The
  largest of these offsets must be less than 15.

Additional current restrictions which may be addressed:
- Many field types are not supported: float, double, sint32/sint64,
  fixed32/fixed64, sfixed32/sfixed64, groups (deprecated).
- Unknown fields are not preserved.
- Protobuf enums are treated as `uint32`s.
- The protobuf spec allows for a non-repeated message to be encoded across
  multiple field id/value pairs, which must then be merged. Decode returns an
  error for this (I don't understand when this would even happen).
- The generated structs do not implement the `proto.Message` interface.
- Extensions are not supported.
- Nested messages are not supported.
- Referencing a message defined in another file probably doesn't work.

Additional restrictions that are unlikely to ever be addressed:
- Repeated fields can only be iterated, not indexed.


## Usage

The following protobuf messages (a subset of the ones defined by the etcd/raft
library) will generate two Go structs, `Entry` and `Message`.

```protobuf
message Entry {
  uint64 term = 2;
  uint64 index = 3;
  EntryType type = 1;
  bytes data = 4;
}

message Snapshot {
  bytes data = 1;
  SnapshotMetadata metadata = 2;
}

message Message {
  MessageType type = 1;
  uint64 to = 2;
  uint64 from = 3;
  uint64 term = 4;
  uint64 logTerm = 5;
  uint64 index = 6;
  repeated Entry entries = 7;
  uint64 commit = 8;
  Snapshot snapshot = 9;
  bool reject = 10;
  uint64 rejectHint = 11;
  bytes context = 12;
}
```

The generated structs are thin wrappers around the encoded bytes with a (usually
stack-allocated) map of the offset of the first occurrence of each field.

Notably, the protobuf fields are not fields on the generated struct. Fields
values are lazily decoded and returned by methods on the struct that are named
after the field. This allows for cheaply extracting one field of an encoded
message with many fields.

```golang
struct Entry {
  buf []byte
  offsets zeropb.FastIntMap
}

func (e *Entry) Decode(buf []byte) error { ... }
func (e *Entry) Term() uint64 { ... }
func (e *Entry) Index() uint64 { ... }
func (e *Entry) Type() uint32 { ... }
func (e *Entry) Data() []byte { ... }
```

The `Decode` message accepts an encoded message and fills this map of offsets
and checks the validity of the encoded bytes. As mentioned above, it does not
decode the field values. The encoded message must not be modified for the
lifetime of the decoded struct.

```golang
var buf []byte = getEncodedMessage()
var entry Entry
if err := entry.Decode(buf); err != nil {
  return err
}
fmt.Println(entry.Term())
```

To decode a field containing some other protobuf message, an instance of the
generated struct for the message is given to be filled. Note that in contrast to
other Go protobuf variants, field values are decoded lazily and so there is the
possibility of an invalid embedded message, which means this field access may
return an error.

```golang
var message Message = getDecodedMessage()
var snapshot Snapshot
if err := message.Snapshot(&snapshot); err != nil {
  return err
}
fmt.Println(snapshot.Data())
```

The field accessor for a repeated field returns an iterator. To keep zeropb
allocation free, this iterator re-parses part of the encoded message each time
it iterates. This means repeatedly iterating a repeated field is more costly
than in other Go protobuf variants. In the worst case (a repeated field of
length one in a message with many field id/value pairs), it is more costly for
the first iteration.

```golang
var message Message = getDecodedMessage()
var entry Entry
for entryIt := message.Entries();; {
  if err := entryIt.Next(&entry); err != nil {
    return err
  }
  fmt.Println(entry.Term())
}
```


## Benchmarks

Measuring pure decode speed by decoding a message and accessing no fields. Note
that this is not an apples-to-apples comparison because zeropb decodes the
fields lazily (on access).

    name                                         time/op
    DecodeSimpleAccessNone/pb-8                   176ns ± 1%
    DecodeSimpleAccessNone/gogopb-8               101ns ± 1%
    DecodeSimpleAccessNone/zeropb-8              62.8ns ± 1%

    name                                         alloc/op
    DecodeSimpleAccessNone/pb-8                    136B ± 0%
    DecodeSimpleAccessNone/gogopb-8                112B ± 0%
    DecodeSimpleAccessNone/zeropb-8               0.00B

    name                                         allocs/op
    DecodeSimpleAccessNone/pb-8                    4.00 ± 0%
    DecodeSimpleAccessNone/gogopb-8                1.00 ± 0%
    DecodeSimpleAccessNone/zeropb-8                0.00

An apples-to-apples comparison of speed of decoding a message and all of its
fields.

    name                                         time/op
    DecodeSimpleAccessAll/pb-8                    226ns ± 1%
    DecodeSimpleAccessAll/gogopb-8                152ns ± 1%
    DecodeSimpleAccessAll/zeropb-8                157ns ± 0%

    name                                         alloc/op
    DecodeSimpleAccessAll/pb-8                     136B ± 0%
    DecodeSimpleAccessAll/gogopb-8                 112B ± 0%
    DecodeSimpleAccessAll/zeropb-8                0.00B

    name                                         allocs/op
    DecodeSimpleAccessAll/pb-8                     4.00 ± 0%
    DecodeSimpleAccessAll/gogopb-8                 1.00 ± 0%
    DecodeSimpleAccessAll/zeropb-8                 0.00

Fields are decoded lazily and not cached, so repeatedly using the same fields is
slower than other libraries.

    name                                         time/op
    DecodeSimpleAccessRepeatedly/pb-8             332ns ± 0%
    DecodeSimpleAccessRepeatedly/gogopb-8         253ns ± 1%
    DecodeSimpleAccessRepeatedly/zeropb-8         349ns ± 0%

    name                                         alloc/op
    DecodeSimpleAccessRepeatedly/pb-8              136B ± 0%
    DecodeSimpleAccessRepeatedly/gogopb-8          112B ± 0%
    DecodeSimpleAccessRepeatedly/zeropb-8         0.00B

    name                                         allocs/op
    DecodeSimpleAccessRepeatedly/pb-8              4.00 ± 0%
    DecodeSimpleAccessRepeatedly/gogopb-8          1.00 ± 0%
    DecodeSimpleAccessRepeatedly/zeropb-8          0.00

The speedup of lazy field decoding is even more pronounced when using one field
out of a more complex message. In this case, there are byte and message fields
that are not used but waste cpu and even cause allocations in eager decoding.
Also note that zeropb's message decode cost is roughly proportional to the
number of top-level fields set in the message.

    name                                         time/op
    DecodeComplexAccessOne/pb-8                  1.39µs ± 0%
    DecodeComplexAccessOne/gogopb-8               771ns ± 1%
    DecodeComplexAccessOne/zeropb-8               379ns ± 0%

    name                                         alloc/op
    DecodeComplexAccessOne/pb-8                    976B ± 0%
    DecodeComplexAccessOne/gogopb-8                960B ± 0%
    DecodeComplexAccessOne/zeropb-8               0.00B

    name                                         allocs/op
    DecodeComplexAccessOne/pb-8                    25.0 ± 0%
    DecodeComplexAccessOne/gogopb-8                7.00 ± 0%
    DecodeComplexAccessOne/zeropb-8                0.00

A measurement of zeropb's unfortunate re-parsing of the message for repeated
fields.

    name                                         time/op
    DecodeComplexAccessRepeatedMessage/pb-8      1.55µs ± 0%
    DecodeComplexAccessRepeatedMessage/gogopb-8   931ns ± 1%
    DecodeComplexAccessRepeatedMessage/zeropb-8   894ns ± 1%

    name                                         alloc/op
    DecodeComplexAccessRepeatedMessage/pb-8        976B ± 0%
    DecodeComplexAccessRepeatedMessage/gogopb-8    960B ± 0%
    DecodeComplexAccessRepeatedMessage/zeropb-8   0.00B

    name                                         allocs/op
    DecodeComplexAccessRepeatedMessage/pb-8        25.0 ± 0%
    DecodeComplexAccessRepeatedMessage/gogopb-8    7.00 ± 0%
    DecodeComplexAccessRepeatedMessage/zeropb-8    0.00


## The Road to v1.0.0

- [ ] Support message encoding
- [ ] Instead of blindly using FastIntMap for each message, tailor the offsets
  to the actual field ids in the message. This would allow us to be
  allocation-free for are larger set of messages with (the common case of) small
  and dense field ids.
- [ ] Support all protobuf field types.
- [ ] Generate and use Go types for each protobuf enum.
- [ ] More test coverage.
- [ ] Fuzz testing.
- [ ] More benchmark coverage.
- [ ] Verify that the accessors are all being inlined.


[Protocol Buffer]: https://developers.google.com/protocol-buffers/
[FlatBuffers]: https://google.github.io/flatbuffers/
[Cap’n Proto]: https://capnproto.org/
