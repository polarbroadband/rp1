// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: dns.proto

package dns

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Zone struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Commit  string               `protobuf:"bytes,1,opt,name=Commit,proto3" json:"Commit,omitempty"`
	Records map[string]*Category `protobuf:"bytes,2,rep,name=Records,proto3" json:"Records,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3" yaml:"records"`
}

func (x *Zone) Reset() {
	*x = Zone{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dns_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Zone) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Zone) ProtoMessage() {}

func (x *Zone) ProtoReflect() protoreflect.Message {
	mi := &file_dns_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Zone.ProtoReflect.Descriptor instead.
func (*Zone) Descriptor() ([]byte, []int) {
	return file_dns_proto_rawDescGZIP(), []int{0}
}

func (x *Zone) GetCommit() string {
	if x != nil {
		return x.Commit
	}
	return ""
}

func (x *Zone) GetRecords() map[string]*Category {
	if x != nil {
		return x.Records
	}
	return nil
}

type Category struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type map[string]*Record `protobuf:"bytes,1,rep,name=Type,proto3" json:"Type,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3" yaml:",inline"`
}

func (x *Category) Reset() {
	*x = Category{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dns_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Category) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Category) ProtoMessage() {}

func (x *Category) ProtoReflect() protoreflect.Message {
	mi := &file_dns_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Category.ProtoReflect.Descriptor instead.
func (*Category) Descriptor() ([]byte, []int) {
	return file_dns_proto_rawDescGZIP(), []int{1}
}

func (x *Category) GetType() map[string]*Record {
	if x != nil {
		return x.Type
	}
	return nil
}

type Record struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Addr []string `protobuf:"bytes,1,rep,name=Addr,proto3" json:"Addr,omitempty" yaml:"addr"`
	TTL  int64    `protobuf:"varint,2,opt,name=TTL,proto3" json:"TTL,omitempty" yaml:"ttl"`
}

func (x *Record) Reset() {
	*x = Record{}
	if protoimpl.UnsafeEnabled {
		mi := &file_dns_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Record) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Record) ProtoMessage() {}

func (x *Record) ProtoReflect() protoreflect.Message {
	mi := &file_dns_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Record.ProtoReflect.Descriptor instead.
func (*Record) Descriptor() ([]byte, []int) {
	return file_dns_proto_rawDescGZIP(), []int{2}
}

func (x *Record) GetAddr() []string {
	if x != nil {
		return x.Addr
	}
	return nil
}

func (x *Record) GetTTL() int64 {
	if x != nil {
		return x.TTL
	}
	return 0
}

var File_dns_proto protoreflect.FileDescriptor

var file_dns_proto_rawDesc = []byte{
	0x0a, 0x09, 0x64, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x64, 0x6e, 0x73,
	0x22, 0x9b, 0x01, 0x0a, 0x04, 0x5a, 0x6f, 0x6e, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x43, 0x6f, 0x6d,
	0x6d, 0x69, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x43, 0x6f, 0x6d, 0x6d, 0x69,
	0x74, 0x12, 0x30, 0x0a, 0x07, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x16, 0x2e, 0x64, 0x6e, 0x73, 0x2e, 0x5a, 0x6f, 0x6e, 0x65, 0x2e, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x52, 0x65, 0x63, 0x6f,
	0x72, 0x64, 0x73, 0x1a, 0x49, 0x0a, 0x0c, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x23, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x64, 0x6e, 0x73, 0x2e, 0x43, 0x61, 0x74, 0x65, 0x67,
	0x6f, 0x72, 0x79, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x7d,
	0x0a, 0x08, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x12, 0x2b, 0x0a, 0x04, 0x54, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x64, 0x6e, 0x73, 0x2e, 0x43,
	0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65, 0x1a, 0x44, 0x0a, 0x09, 0x54, 0x79, 0x70, 0x65, 0x45,
	0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x21, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x64, 0x6e, 0x73, 0x2e, 0x52, 0x65, 0x63, 0x6f,
	0x72, 0x64, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x2e, 0x0a,
	0x06, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x41, 0x64, 0x64, 0x72, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x41, 0x64, 0x64, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x54,
	0x54, 0x4c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x54, 0x54, 0x4c, 0x42, 0x06, 0x5a,
	0x04, 0x2f, 0x64, 0x6e, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_dns_proto_rawDescOnce sync.Once
	file_dns_proto_rawDescData = file_dns_proto_rawDesc
)

func file_dns_proto_rawDescGZIP() []byte {
	file_dns_proto_rawDescOnce.Do(func() {
		file_dns_proto_rawDescData = protoimpl.X.CompressGZIP(file_dns_proto_rawDescData)
	})
	return file_dns_proto_rawDescData
}

var file_dns_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_dns_proto_goTypes = []interface{}{
	(*Zone)(nil),     // 0: dns.Zone
	(*Category)(nil), // 1: dns.Category
	(*Record)(nil),   // 2: dns.Record
	nil,              // 3: dns.Zone.RecordsEntry
	nil,              // 4: dns.Category.TypeEntry
}
var file_dns_proto_depIdxs = []int32{
	3, // 0: dns.Zone.Records:type_name -> dns.Zone.RecordsEntry
	4, // 1: dns.Category.Type:type_name -> dns.Category.TypeEntry
	1, // 2: dns.Zone.RecordsEntry.value:type_name -> dns.Category
	2, // 3: dns.Category.TypeEntry.value:type_name -> dns.Record
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_dns_proto_init() }
func file_dns_proto_init() {
	if File_dns_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_dns_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Zone); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_dns_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Category); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_dns_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Record); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_dns_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_dns_proto_goTypes,
		DependencyIndexes: file_dns_proto_depIdxs,
		MessageInfos:      file_dns_proto_msgTypes,
	}.Build()
	File_dns_proto = out.File
	file_dns_proto_rawDesc = nil
	file_dns_proto_goTypes = nil
	file_dns_proto_depIdxs = nil
}
