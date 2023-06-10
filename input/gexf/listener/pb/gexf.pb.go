// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: gexf.proto

package pb

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GexfRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileContent []byte `protobuf:"bytes,1,opt,name=file_content,json=fileContent,proto3" json:"file_content,omitempty"`
	NetworkName string `protobuf:"bytes,2,opt,name=network_name,json=networkName,proto3" json:"network_name,omitempty"`
}

func (x *GexfRequest) Reset() {
	*x = GexfRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gexf_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GexfRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GexfRequest) ProtoMessage() {}

func (x *GexfRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gexf_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GexfRequest.ProtoReflect.Descriptor instead.
func (*GexfRequest) Descriptor() ([]byte, []int) {
	return file_gexf_proto_rawDescGZIP(), []int{0}
}

func (x *GexfRequest) GetFileContent() []byte {
	if x != nil {
		return x.FileContent
	}
	return nil
}

func (x *GexfRequest) GetNetworkName() string {
	if x != nil {
		return x.NetworkName
	}
	return ""
}

type GexfResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *GexfResponse) Reset() {
	*x = GexfResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gexf_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GexfResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GexfResponse) ProtoMessage() {}

func (x *GexfResponse) ProtoReflect() protoreflect.Message {
	mi := &file_gexf_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GexfResponse.ProtoReflect.Descriptor instead.
func (*GexfResponse) Descriptor() ([]byte, []int) {
	return file_gexf_proto_rawDescGZIP(), []int{1}
}

func (x *GexfResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_gexf_proto protoreflect.FileDescriptor

var file_gexf_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x67, 0x65, 0x78, 0x66, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1a, 0x67, 0x72,
	0x61, 0x70, 0x68, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x43,
	0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x22, 0x53, 0x0a, 0x0b, 0x47, 0x65, 0x78, 0x66,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x66, 0x69, 0x6c, 0x65, 0x5f,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x66,
	0x69, 0x6c, 0x65, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x6e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x28, 0x0a,
	0x0c, 0x47, 0x65, 0x78, 0x66, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07,
	0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x32, 0x6f, 0x0a, 0x0b, 0x47, 0x65, 0x78, 0x66, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x60, 0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73,
	0x73, 0x47, 0x65, 0x78, 0x66, 0x12, 0x27, 0x2e, 0x67, 0x72, 0x61, 0x70, 0x68, 0x41, 0x6e, 0x61,
	0x6c, 0x79, 0x7a, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x2e, 0x47, 0x65, 0x78, 0x66, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28,
	0x2e, 0x67, 0x72, 0x61, 0x70, 0x68, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x72, 0x44, 0x61,
	0x74, 0x61, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x47, 0x65, 0x78, 0x66,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x18, 0x5a, 0x16, 0x69, 0x6e, 0x70, 0x75,
	0x74, 0x2f, 0x67, 0x65, 0x78, 0x66, 0x2f, 0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x65, 0x72, 0x2f,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gexf_proto_rawDescOnce sync.Once
	file_gexf_proto_rawDescData = file_gexf_proto_rawDesc
)

func file_gexf_proto_rawDescGZIP() []byte {
	file_gexf_proto_rawDescOnce.Do(func() {
		file_gexf_proto_rawDescData = protoimpl.X.CompressGZIP(file_gexf_proto_rawDescData)
	})
	return file_gexf_proto_rawDescData
}

var file_gexf_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_gexf_proto_goTypes = []interface{}{
	(*GexfRequest)(nil),  // 0: graphAnalyzerDataCollector.GexfRequest
	(*GexfResponse)(nil), // 1: graphAnalyzerDataCollector.GexfResponse
}
var file_gexf_proto_depIdxs = []int32{
	0, // 0: graphAnalyzerDataCollector.GexfService.ProcessGexf:input_type -> graphAnalyzerDataCollector.GexfRequest
	1, // 1: graphAnalyzerDataCollector.GexfService.ProcessGexf:output_type -> graphAnalyzerDataCollector.GexfResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_gexf_proto_init() }
func file_gexf_proto_init() {
	if File_gexf_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gexf_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GexfRequest); i {
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
		file_gexf_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GexfResponse); i {
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
			RawDescriptor: file_gexf_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_gexf_proto_goTypes,
		DependencyIndexes: file_gexf_proto_depIdxs,
		MessageInfos:      file_gexf_proto_msgTypes,
	}.Build()
	File_gexf_proto = out.File
	file_gexf_proto_rawDesc = nil
	file_gexf_proto_goTypes = nil
	file_gexf_proto_depIdxs = nil
}