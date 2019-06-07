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

_**Warning**_: There is a very specific set of tradeoffs made in zeropb that
make it unsuitable for many uses. In particular, with lazy decoding you "pay for
what you use" and with zero allocations you get improved performance in hot
paths. However, lazy decoding means that bad usage patterns can be quite slow
and zero allocations means awkward ergonomics (though it will get a bit better
than it currently is). If you're not specifically looking to make these
tradeoffs, you'd probably be better served by another library.


## Compromises

In contrast to [FlatBuffers] and [Cap’n Proto], the protobuf wire format was not
specifically designed for allocation-less use. This means that zeropb ends up as
somewhat more awkward to use than other protobuf variants for Go, such as the
official Google one or gogoprotobuf.

It also means that it is not always possible to decode without heap allocations.
This list will become _considerably_ less restrictive over time, but the current
set of requirements for allocation-less use are as follows:

- All field ids are less than 32.
- The encoded/decoded messages are not too large, where too large is complicated
  to define but can be thought of as a handful of small fields and one big or
  repeated field.

  Specifically, if largest offset of the first appearance of a field in an
  encoded message is 15 or greater, there will be one allocation per message.
  That is, an encoded message is a concatenation of field id/value pairs. Field
  ids may be repeated. Find the offset of the first appearance of each field id
  in the encoded message. The largest of these offsets must be less than 15.


## The Road to v1.0.0

- [x] Support message encoding
- [x] Support all protobuf field types.
  - [ ] Also support them for repeated fields.
- [ ] Generate and use Go types for each protobuf enum. They are currently
  treated at `uint32`s.
- [ ] Improve ergonomics of constructing trees of messages.
- [ ] More test coverage.
- [ ] Fuzz testing.
- [x] More benchmark coverage.
- [ ] Verify that the accessors are all being inlined.
- [ ] Instead of blindly using FastIntMap for each message, tailor the offsets
  to the actual field ids in the message. This would allow us to be
  allocation-free for are larger set of messages with (the common case of) small
  and dense field ids.
- [ ] Preserve unknown fields.
- [ ] Support nested messages.
- [ ] Support references proto messages from other files/packages.
- [ ] The protobuf spec allows for a non-repeated message to be encoded across
  multiple field id/value pairs, which must then be merged. Decode returns an
  error for this (I don't understand when this would even happen).

Additional current restrictions which may be addressed:
- Groups (deprecated) are not supported.
- Extensions are not supported.
- The generated structs do not implement the `proto.Message` interface.
- Repeated fields can only be iterated, not indexed.


## Usage

Install the codegen plugin binary, possibly with go get:

```sh
> go get -u github.com/danhhz/zeropb/cmd/protoc-gen-zeropb
```

Then run protoc to get the generated code. (If you have `$GOPATH/bin` in your
`PATH`, you can omit the --plugin argument.)

```sh
> protoc --plugin=$GOPATH/bin/protoc-gen-zeropb raft.proto --zeropb_out=.
```

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

(See the full current generated [raft.zeropb.go] code for [raft.proto].)

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
    DecodeSimpleAccessNone/pb-8                     177ns ± 0%
    DecodeSimpleAccessNone/gogopb-8                 102ns ± 0%
    DecodeSimpleAccessNone/zeropb-8                69.8ns ± 1%

    name                                         speed
    DecodeSimpleAccessNone/pb-8                   610MB/s ± 1%
    DecodeSimpleAccessNone/gogopb-8              1.05GB/s ± 2%
    DecodeSimpleAccessNone/zeropb-8              1.55GB/s ± 1%

    name                                         allocs/op
    DecodeSimpleAccessNone/pb-8                      4.00 ± 0%
    DecodeSimpleAccessNone/gogopb-8                  1.00 ± 0%
    DecodeSimpleAccessNone/zeropb-8                  0.00

An apples-to-apples comparison of speed of decoding a message and all of its
fields.

    name                                         time/op
    DecodeSimpleAccessAll/pb-8                      228ns ± 0%
    DecodeSimpleAccessAll/gogopb-8                  150ns ± 2%
    DecodeSimpleAccessAll/zeropb-8                  162ns ± 0%

    name                                         speed
    DecodeSimpleAccessAll/pb-8                    473MB/s ± 0%
    DecodeSimpleAccessAll/gogopb-8                717MB/s ± 2%
    DecodeSimpleAccessAll/zeropb-8                662MB/s ± 2%

    name                                         allocs/op
    DecodeSimpleAccessAll/pb-8                       4.00 ± 0%
    DecodeSimpleAccessAll/gogopb-8                   1.00 ± 0%
    DecodeSimpleAccessAll/zeropb-8                   0.00

Fields are decoded lazily and not cached, so repeatedly using the same fields is
slower than other libraries.

    name                                         time/op
    DecodeSimpleAccessRepeatedly/pb-8               347ns ± 1%
    DecodeSimpleAccessRepeatedly/gogopb-8           259ns ± 0%
    DecodeSimpleAccessRepeatedly/zeropb-8           352ns ± 0%

    name                                         speed
    DecodeSimpleAccessRepeatedly/pb-8             311MB/s ± 0%
    DecodeSimpleAccessRepeatedly/gogopb-8         416MB/s ± 0%
    DecodeSimpleAccessRepeatedly/zeropb-8         306MB/s ± 0%

    name                                         allocs/op
    DecodeSimpleAccessRepeatedly/pb-8                4.00 ± 0%
    DecodeSimpleAccessRepeatedly/gogopb-8            1.00 ± 0%
    DecodeSimpleAccessRepeatedly/zeropb-8            0.00

The speedup of lazy field decoding is even more pronounced when using one field
out of a more complex message. In this case, there are byte and message fields
that are not used but waste cpu and cause allocations because of eager decoding.
Also note that zeropb's message decode cost is roughly proportional to the
number of top-level fields set in the message and independent of any fields in
sub-messages.

    name                                         time/op
    DecodeComplexAccessOne/pb-8                    1.66µs ± 3%
    DecodeComplexAccessOne/gogopb-8                 829ns ± 2%
    DecodeComplexAccessOne/zeropb-8                 485ns ± 2%

    name                                         speed
    DecodeComplexAccessOne/pb-8                   278MB/s ± 3%
    DecodeComplexAccessOne/gogopb-8               533MB/s ±16%
    DecodeComplexAccessOne/zeropb-8               948MB/s ± 2%

    name                                         allocs/op
    DecodeComplexAccessOne/pb-8                      33.0 ± 0%
    DecodeComplexAccessOne/gogopb-8                  7.00 ± 0%
    DecodeComplexAccessOne/zeropb-8                  0.00

A measurement of zeropb's unfortunate re-parsing of the message for repeated
fields.

    name                                         time/op
    DecodeComplexAccessRepeatedMessage/pb-8        1.83µs ± 0%
    DecodeComplexAccessRepeatedMessage/gogopb-8    1.00µs ± 1%
    DecodeComplexAccessRepeatedMessage/zeropb-8    1.19µs ± 1%

    name                                         speed
    DecodeComplexAccessRepeatedMessage/pb-8       251MB/s ± 0%
    DecodeComplexAccessRepeatedMessage/gogopb-8   459MB/s ± 1%
    DecodeComplexAccessRepeatedMessage/zeropb-8   386MB/s ± 1%

    name                                         allocs/op
    DecodeComplexAccessRepeatedMessage/pb-8          33.0 ± 0%
    DecodeComplexAccessRepeatedMessage/gogopb-8      7.00 ± 0%
    DecodeComplexAccessRepeatedMessage/zeropb-8      0.00

An apples-to-apples comparison of setting every field in a message and encoding
it. Encode in zeropb is a no-op, the encoded message is maintained with each
call to a field setter, so if we pulled setting the fields out of this
benchmark, zeropb would be infinitely fast :-D!

    name                                         time/op
    EncodeSimpleSetAll/pb-8                         212ns ± 0%
    EncodeSimpleSetAll/gogopb-8                     175ns ± 1%
    EncodeSimpleSetAll/zeropb-8                     135ns ± 0%

    name                                         speed
    EncodeSimpleSetAll/pb-8                       508MB/s ± 0%
    EncodeSimpleSetAll/gogopb-8                   616MB/s ± 1%
    EncodeSimpleSetAll/zeropb-8                   800MB/s ± 0%

    name                                         allocs/op
    EncodeSimpleSetAll/pb-8                          2.00 ± 0%
    EncodeSimpleSetAll/gogopb-8                      2.00 ± 0%
    EncodeSimpleSetAll/zeropb-8                      0.00

Similar to repeatedly reading a field from a decoded message, repeatedly setting
a field is slower than other libraries. This is because we maintain the encoded
message with each call to a field setter, so they're doing much more work than
setting a field on a go struct.

    name                                         time/op
    EncodeSimpleSetRepeatedly/pb-8                  264ns ± 2%
    EncodeSimpleSetRepeatedly/gogopb-8              188ns ± 1%
    EncodeSimpleSetRepeatedly/zeropb-8              366ns ± 0%

    name                                         speed
    EncodeSimpleSetRepeatedly/pb-8                409MB/s ± 2%
    EncodeSimpleSetRepeatedly/gogopb-8            571MB/s ± 1%
    EncodeSimpleSetRepeatedly/zeropb-8            294MB/s ± 0%

    name                                         allocs/op
    EncodeSimpleSetRepeatedly/pb-8                   4.00 ± 0%
    EncodeSimpleSetRepeatedly/gogopb-8               2.00 ± 0%
    EncodeSimpleSetRepeatedly/zeropb-8               0.00

Encoding a complex message is currently slower than other libraries. This is
expected to get faster, but I haven't yet done speed of light measurements to
see exactly how much faster.

This benchmark currently causes allocations because the resulting encoded
message is large enough that it triggers one of the conditions described in the
Compromises section, but this is not intrinsic. These allocations will go away
when we specialize FastIntMap.

    name                                         time/op
    EncodeComplex/pb-8                             1.06µs ± 0%
    EncodeComplex/gogopb-8                          754ns ± 5%
    EncodeComplex/zeropb-8                         1.34µs ± 1%

    name                                         speed
    EncodeComplex/pb-8                            420MB/s ± 0%
    EncodeComplex/gogopb-8                        610MB/s ± 5%
    EncodeComplex/zeropb-8                        325MB/s ± 1%

    name                                         allocs/op
    EncodeComplex/pb-8                               7.00 ± 0%
    EncodeComplex/gogopb-8                           3.00 ± 0%
    EncodeComplex/zeropb-8                           2.00 ± 0%


[Protocol Buffer]: https://developers.google.com/protocol-buffers/
[FlatBuffers]: https://google.github.io/flatbuffers/
[Cap’n Proto]: https://capnproto.org/
[raft.zeropb.go]: https://github.com/danhhz/zeropb/blob/master/golden/raftzeropb/raft.zeropb.go
[raft.proto]: https://github.com/danhhz/zeropb/blob/master/golden/raft.proto
