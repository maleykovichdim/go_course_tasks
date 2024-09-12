// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: index_service.proto

package indexpb

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

type Document struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Url   string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	Title string `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	Body  string `protobuf:"bytes,4,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *Document) Reset() {
	*x = Document{}
	if protoimpl.UnsafeEnabled {
		mi := &file_index_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Document) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Document) ProtoMessage() {}

func (x *Document) ProtoReflect() protoreflect.Message {
	mi := &file_index_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Document.ProtoReflect.Descriptor instead.
func (*Document) Descriptor() ([]byte, []int) {
	return file_index_service_proto_rawDescGZIP(), []int{0}
}

func (x *Document) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Document) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *Document) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Document) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

type Service struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Counter uint32                            `protobuf:"varint,1,opt,name=counter,proto3" json:"counter,omitempty"`
	Links   []*Document                       `protobuf:"bytes,2,rep,name=links,proto3" json:"links,omitempty"`
	Index   map[string]*Service_MapFieldEntry `protobuf:"bytes,3,rep,name=index,proto3" json:"index,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Service) Reset() {
	*x = Service{}
	if protoimpl.UnsafeEnabled {
		mi := &file_index_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Service) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Service) ProtoMessage() {}

func (x *Service) ProtoReflect() protoreflect.Message {
	mi := &file_index_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Service.ProtoReflect.Descriptor instead.
func (*Service) Descriptor() ([]byte, []int) {
	return file_index_service_proto_rawDescGZIP(), []int{1}
}

func (x *Service) GetCounter() uint32 {
	if x != nil {
		return x.Counter
	}
	return 0
}

func (x *Service) GetLinks() []*Document {
	if x != nil {
		return x.Links
	}
	return nil
}

func (x *Service) GetIndex() map[string]*Service_MapFieldEntry {
	if x != nil {
		return x.Index
	}
	return nil
}

type Service_MapFieldEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Index []uint32 `protobuf:"varint,1,rep,packed,name=index,proto3" json:"index,omitempty"`
}

func (x *Service_MapFieldEntry) Reset() {
	*x = Service_MapFieldEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_index_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Service_MapFieldEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Service_MapFieldEntry) ProtoMessage() {}

func (x *Service_MapFieldEntry) ProtoReflect() protoreflect.Message {
	mi := &file_index_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Service_MapFieldEntry.ProtoReflect.Descriptor instead.
func (*Service_MapFieldEntry) Descriptor() ([]byte, []int) {
	return file_index_service_proto_rawDescGZIP(), []int{1, 0}
}

func (x *Service_MapFieldEntry) GetIndex() []uint32 {
	if x != nil {
		return x.Index
	}
	return nil
}

var File_index_service_proto protoreflect.FileDescriptor

var file_index_service_proto_rawDesc = []byte{
	0x0a, 0x13, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x70, 0x62, 0x22, 0x56,
	0x0a, 0x08, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72,
	0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x14, 0x0a, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x22, 0x80, 0x02, 0x0a, 0x07, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x12, 0x27, 0x0a, 0x05,
	0x6c, 0x69, 0x6e, 0x6b, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x70, 0x62, 0x2e, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x05,
	0x6c, 0x69, 0x6e, 0x6b, 0x73, 0x12, 0x31, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x70, 0x62, 0x2e, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x1a, 0x25, 0x0a, 0x0d, 0x4d, 0x61, 0x70, 0x46,
	0x69, 0x65, 0x6c, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64,
	0x65, 0x78, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x1a,
	0x58, 0x0a, 0x0a, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x34, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e,
	0x2e, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x70, 0x62, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x4d, 0x61, 0x70, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x18, 0x5a, 0x16, 0x2e, 0x2e, 0x2f,
	0x70, 0x6b, 0x67, 0x2f, 0x63, 0x72, 0x61, 0x77, 0x6c, 0x65, 0x72, 0x2f, 0x69, 0x6e, 0x64, 0x65,
	0x78, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_index_service_proto_rawDescOnce sync.Once
	file_index_service_proto_rawDescData = file_index_service_proto_rawDesc
)

func file_index_service_proto_rawDescGZIP() []byte {
	file_index_service_proto_rawDescOnce.Do(func() {
		file_index_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_index_service_proto_rawDescData)
	})
	return file_index_service_proto_rawDescData
}

var file_index_service_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_index_service_proto_goTypes = []any{
	(*Document)(nil),              // 0: indexpb.Document
	(*Service)(nil),               // 1: indexpb.Service
	(*Service_MapFieldEntry)(nil), // 2: indexpb.Service.MapFieldEntry
	nil,                           // 3: indexpb.Service.IndexEntry
}
var file_index_service_proto_depIdxs = []int32{
	0, // 0: indexpb.Service.links:type_name -> indexpb.Document
	3, // 1: indexpb.Service.index:type_name -> indexpb.Service.IndexEntry
	2, // 2: indexpb.Service.IndexEntry.value:type_name -> indexpb.Service.MapFieldEntry
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_index_service_proto_init() }
func file_index_service_proto_init() {
	if File_index_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_index_service_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Document); i {
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
		file_index_service_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Service); i {
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
		file_index_service_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*Service_MapFieldEntry); i {
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
			RawDescriptor: file_index_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_index_service_proto_goTypes,
		DependencyIndexes: file_index_service_proto_depIdxs,
		MessageInfos:      file_index_service_proto_msgTypes,
	}.Build()
	File_index_service_proto = out.File
	file_index_service_proto_rawDesc = nil
	file_index_service_proto_goTypes = nil
	file_index_service_proto_depIdxs = nil
}
