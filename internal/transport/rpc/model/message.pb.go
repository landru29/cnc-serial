// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.12.4
// source: model/message.proto

package model

import (
	empty "github.com/golang/protobuf/ptypes/empty"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Command struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Data          string                 `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Command) Reset() {
	*x = Command{}
	mi := &file_model_message_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Command) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Command) ProtoMessage() {}

func (x *Command) ProtoReflect() protoreflect.Message {
	mi := &file_model_message_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Command.ProtoReflect.Descriptor instead.
func (*Command) Descriptor() ([]byte, []int) {
	return file_model_message_proto_rawDescGZIP(), []int{0}
}

func (x *Command) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

type Status struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Data          string                 `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Status) Reset() {
	*x = Status{}
	mi := &file_model_message_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Status) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Status) ProtoMessage() {}

func (x *Status) ProtoReflect() protoreflect.Message {
	mi := &file_model_message_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Status.ProtoReflect.Descriptor instead.
func (*Status) Descriptor() ([]byte, []int) {
	return file_model_message_proto_rawDescGZIP(), []int{1}
}

func (x *Status) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

var File_model_message_proto protoreflect.FileDescriptor

const file_model_message_proto_rawDesc = "" +
	"\n" +
	"\x13model/message.proto\x12\x05model\x1a\x1bgoogle/protobuf/empty.proto\"\x1d\n" +
	"\aCommand\x12\x12\n" +
	"\x04data\x18\x01 \x01(\tR\x04data\"\x1c\n" +
	"\x06Status\x12\x12\n" +
	"\x04data\x18\x01 \x01(\tR\x04data2~\n" +
	"\rCommandSender\x127\n" +
	"\vSendCommand\x12\x0e.model.Command\x1a\x16.google.protobuf.Empty\"\x00\x124\n" +
	"\tGetStatus\x12\x16.google.protobuf.Empty\x1a\r.model.Status\"\x00B=Z;github.com/landru29/cnc-serial/internal/transport/rpc/modelb\x06proto3"

var (
	file_model_message_proto_rawDescOnce sync.Once
	file_model_message_proto_rawDescData []byte
)

func file_model_message_proto_rawDescGZIP() []byte {
	file_model_message_proto_rawDescOnce.Do(func() {
		file_model_message_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_model_message_proto_rawDesc), len(file_model_message_proto_rawDesc)))
	})
	return file_model_message_proto_rawDescData
}

var file_model_message_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_model_message_proto_goTypes = []any{
	(*Command)(nil),     // 0: model.Command
	(*Status)(nil),      // 1: model.Status
	(*empty.Empty)(nil), // 2: google.protobuf.Empty
}
var file_model_message_proto_depIdxs = []int32{
	0, // 0: model.CommandSender.SendCommand:input_type -> model.Command
	2, // 1: model.CommandSender.GetStatus:input_type -> google.protobuf.Empty
	2, // 2: model.CommandSender.SendCommand:output_type -> google.protobuf.Empty
	1, // 3: model.CommandSender.GetStatus:output_type -> model.Status
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_model_message_proto_init() }
func file_model_message_proto_init() {
	if File_model_message_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_model_message_proto_rawDesc), len(file_model_message_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_model_message_proto_goTypes,
		DependencyIndexes: file_model_message_proto_depIdxs,
		MessageInfos:      file_model_message_proto_msgTypes,
	}.Build()
	File_model_message_proto = out.File
	file_model_message_proto_goTypes = nil
	file_model_message_proto_depIdxs = nil
}
