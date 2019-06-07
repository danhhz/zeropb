// Copyright 2019 Daniel Harrison. All Rights Reserved.

package generate

import (
	"strings"
	"unicode"
)

// Generate returns the codegen part of the requested zeropb messages.
func Generate(req *CodeGeneratorRequest) (*CodeGeneratorResponse, error) {
	filesToGenerate := map[string]struct{}{}
	for _, f := range req.FileToGenerate {
		filesToGenerate[f] = struct{}{}
	}
	files := make([]*CodeGeneratorResponse_File, 0, len(filesToGenerate))
	for _, fileReq := range req.ProtoFile {
		if fileReq == nil {
			continue
		}
		if _, ok := filesToGenerate[fileReq.GetName()]; !ok {
			continue
		}
		fileRes, err := generateFile(fileReq)
		if err != nil {
			return nil, err
		}
		files = append(files, fileRes)
	}
	res := &CodeGeneratorResponse{File: files}
	return res, nil
}

func generateFile(req *FileDescriptorProto) (*CodeGeneratorResponse_File, error) {
	var buf strings.Builder
	st := NewStringTree(&buf, `  `)
	var indent StringTreeIndent
	st.Write(indent, "// Code generated by protoc-gen-zeropb. DO NOT EDIT.\n\n")
	// TODO(dan): Use the go package if it's set.
	st.Write(indent, "package ", req.GetPackage(), "\n\n")
	st.Write(indent, `import "github.com/danhhz/zeropb"`, "\n\n")

	for _, m := range req.MessageType {
		if m == nil {
			continue
		}
		generateMessage(st, indent, *m)
	}

	name := strings.ReplaceAll(req.GetName(), `.proto`, `.zeropb.go`)
	content := buf.String()
	res := &CodeGeneratorResponse_File{Name: &name, Content: &content}
	return res, nil
}

func generateMessage(st *StringTree, indent StringTreeIndent, m DescriptorProto) {
	messageGoType := m.GetName()

	st.Write(indent, "type ", messageGoType, " struct {\n")
	st.Write(indent.Next(), "buf []byte\n")
	st.Write(indent.Next(), "offsets zeropb.FastIntMap\n")
	st.Write(indent, "}\n\n")

	st.Write(indent, "func (m *", messageGoType, ") Encode() []byte {\n")
	st.Write(indent.Next(), "return m.buf\n")
	st.Write(indent, "}\n\n")

	st.Write(indent, "func (m *", messageGoType, ") Decode(buf []byte) error {\n")
	st.Write(indent.Next(), "m.buf = buf\n")
	st.Write(indent.Next(), "return zeropb.Decode(m.buf, &m.offsets)\n")
	st.Write(indent, "}\n\n")

	st.Write(indent, "func (m *", messageGoType, ") Reset(buf []byte) {\n")
	st.Write(indent.Next(), "if len(buf) > 0 {\n")
	st.Write(indent.Next().Next(), "panic(`buf must be empty`)\n")
	st.Write(indent.Next(), "}\n")
	st.Write(indent.Next(), "m.buf = buf\n")
	st.Write(indent.Next(), "m.offsets.Clear()\n")
	st.Write(indent, "}\n\n")

	for _, f := range m.Field {
		if f == nil {
			continue
		}
		switch f.GetLabel() {
		case FieldDescriptorProto_LABEL_OPTIONAL, FieldDescriptorProto_LABEL_REQUIRED:
			switch f.GetType() {
			case FieldDescriptorProto_TYPE_MESSAGE:
				generateMessageField(st, indent, messageGoType, *f)
			default:
				generateSimpleField(st, indent, messageGoType, *f)
			}
		case FieldDescriptorProto_LABEL_REPEATED:
			switch f.GetType() {
			case FieldDescriptorProto_TYPE_MESSAGE:
				generateRepeatedMessageField(st, indent, messageGoType, *f)
			default:
				// WIP implement repeated simple fields
				continue
			}
		}
	}
}

func generateSimpleField(
	st *StringTree, indent StringTreeIndent, messageGoType string, f FieldDescriptorProto,
) {
	fieldGoName := toGoCase(f.GetName())
	fieldGoType, fieldFnName := fieldToGoTypeSimple(f.GetType())

	st.Write(indent, "func (m *", messageGoType, ") ", fieldGoName, "() ", fieldGoType, " {\n")
	st.Write(indent.Next(), "return zeropb.Get", fieldFnName, "(m.buf, &m.offsets, ", f.GetNumber(), ")\n")
	st.Write(indent, "}\n\n")

	st.Write(indent, "func (m *", messageGoType, ") Set", fieldGoName, "(x ", fieldGoType, ") {\n")
	st.Write(indent.Next(), "zeropb.Set", fieldFnName, "(&m.buf, &m.offsets, ", f.GetNumber(), ", x)\n")
	st.Write(indent, "}\n\n")
}

func generateMessageField(
	st *StringTree, indent StringTreeIndent, messageGoType string, f FieldDescriptorProto,
) {
	fieldGoName := toGoCase(f.GetName())
	fieldGoType := fieldToGoTypeMessage(f.GetTypeName())

	st.Write(indent, "func (m *", messageGoType, ") ", fieldGoName, "(x *", fieldGoType, ") (bool, error) {\n")
	st.Write(indent.Next(), "buf := zeropb.GetBytes(m.buf, &m.offsets, ", f.GetNumber(), ")\n")
	st.Write(indent.Next(), "if buf == nil {\n")
	st.Write(indent.Next().Next(), "return false, nil\n")
	st.Write(indent.Next(), "}\n")
	st.Write(indent.Next(), "return true, x.Decode(buf)\n")
	st.Write(indent, "}\n\n")

	st.Write(indent, "func (m *", messageGoType, ") Set", fieldGoName, "(x ", fieldGoType, ") {\n")
	st.Write(indent.Next(), "buf := x.Encode()\n")
	st.Write(indent.Next(), "zeropb.SetBytes(&m.buf, &m.offsets, ", f.GetNumber(), ", buf)\n")
	st.Write(indent, "}\n\n")
}

func generateRepeatedMessageField(
	st *StringTree, indent StringTreeIndent, messageGoType string, f FieldDescriptorProto,
) {
	fieldGoName := toGoCase(f.GetName())
	fieldGoType := fieldToGoTypeMessage(f.GetTypeName())
	itGoType := messageGoType + fieldGoType + "Iterator"

	st.Write(indent, "type ", itGoType, " []byte\n\n")

	st.Write(indent, "func (i *", itGoType, ") Next(m *", fieldGoType, ") (bool, error) {\n")
	st.Write(indent.Next(), "var buf []byte\n")
	st.Write(indent.Next(), "*i, buf = zeropb.FindNextField((*i), ", f.GetNumber(), ")\n")
	st.Write(indent.Next(), "if buf == nil {\n")
	st.Write(indent.Next().Next(), "return false, nil\n")
	st.Write(indent.Next(), "}\n")
	st.Write(indent.Next(), "return true, m.Decode(buf)\n")
	st.Write(indent, "}\n\n")

	st.Write(indent, "func (m *", messageGoType, ") ", fieldGoName, "() ", itGoType, " {\n")
	st.Write(indent.Next(), "return ", itGoType, "(zeropb.GetRepeatedNonPacked(m.buf, &m.offsets, ", f.GetNumber(), "))\n")
	st.Write(indent, "}\n\n")

	st.Write(indent, "func (m *", messageGoType, ") AppendTo", fieldGoName, "(x ", fieldGoType, ") {\n")
	st.Write(indent.Next(), "buf := x.Encode()\n")
	st.Write(indent.Next(), "zeropb.AppendBytes(&m.buf, &m.offsets, ", f.GetNumber(), ", buf)\n")
	st.Write(indent, "}\n\n")
}

func toGoCase(s string) string {
	titleNext := true
	mapFn := func(r rune) rune {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			titleNext = true
			return -1
		}
		if titleNext {
			titleNext = false
			return unicode.ToTitle(r)
		}
		return r
	}
	return strings.Map(mapFn, s)
}

func fieldToGoTypeSimple(typ FieldDescriptorProto_Type) (fieldGoType, fieldFnName string) {
	switch typ {
	case FieldDescriptorProto_TYPE_BOOL:
		return `bool`, `Bool`
	case FieldDescriptorProto_TYPE_INT32:
		return `int32`, `Int32`
	case FieldDescriptorProto_TYPE_INT64:
		return `int64`, `Int64`
	case FieldDescriptorProto_TYPE_UINT32:
		return `uint32`, `Uint32`
	case FieldDescriptorProto_TYPE_UINT64:
		return `uint64`, `Uint64`
	case FieldDescriptorProto_TYPE_SINT32:
		return `int32`, `ZigZagInt32`
	case FieldDescriptorProto_TYPE_SINT64:
		return `int64`, `ZigZagInt64`
	case FieldDescriptorProto_TYPE_FIXED32:
		return `uint32`, `FixedUint32`
	case FieldDescriptorProto_TYPE_FIXED64:
		return `uint64`, `FixedUint64`
	case FieldDescriptorProto_TYPE_SFIXED32:
		return `int32`, `FixedInt32`
	case FieldDescriptorProto_TYPE_SFIXED64:
		return `int64`, `FixedInt64`
	case FieldDescriptorProto_TYPE_DOUBLE:
		return `float64`, `Float64`
	case FieldDescriptorProto_TYPE_FLOAT:
		return `float32`, `Float32`
	case FieldDescriptorProto_TYPE_STRING:
		return `string`, `String`
	case FieldDescriptorProto_TYPE_BYTES:
		return `[]byte`, `Bytes`
	case FieldDescriptorProto_TYPE_ENUM:
		return `uint32`, `Uint32`
	}
	panic(typ)
}

func fieldToGoTypeMessage(typName string) string {
	if idx := strings.LastIndex(typName, `.`); idx != -1 {
		typName = typName[idx+1:]
	}
	return typName
}
