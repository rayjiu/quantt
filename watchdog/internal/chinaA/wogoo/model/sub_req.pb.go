// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.3
// source: sub_req.proto

package model

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

// 股票分时请求
type SubReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CommonReq  *CommonReq       `protobuf:"bytes,1,opt,name=commonReq,proto3" json:"commonReq,omitempty"`   // 基础请求
	SubMarkets []*SubMarketInfo `protobuf:"bytes,2,rep,name=subMarkets,proto3" json:"subMarkets,omitempty"` // 订阅市场信息
	SubStocks  []*SubStockInfo  `protobuf:"bytes,3,rep,name=subStocks,proto3" json:"subStocks,omitempty"`   // 订阅股票信息
}

func (x *SubReq) Reset() {
	*x = SubReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sub_req_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubReq) ProtoMessage() {}

func (x *SubReq) ProtoReflect() protoreflect.Message {
	mi := &file_sub_req_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubReq.ProtoReflect.Descriptor instead.
func (*SubReq) Descriptor() ([]byte, []int) {
	return file_sub_req_proto_rawDescGZIP(), []int{0}
}

func (x *SubReq) GetCommonReq() *CommonReq {
	if x != nil {
		return x.CommonReq
	}
	return nil
}

func (x *SubReq) GetSubMarkets() []*SubMarketInfo {
	if x != nil {
		return x.SubMarkets
	}
	return nil
}

func (x *SubReq) GetSubStocks() []*SubStockInfo {
	if x != nil {
		return x.SubStocks
	}
	return nil
}

var File_sub_req_proto protoreflect.FileDescriptor

var file_sub_req_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x75, 0x62, 0x5f, 0x72, 0x65, 0x71, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x1a, 0x10, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x72,
	0x65, 0x71, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x73, 0x75, 0x62, 0x5f, 0x6d, 0x61,
	0x72, 0x6b, 0x65, 0x74, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x14, 0x73, 0x75, 0x62, 0x5f, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa1, 0x01, 0x0a, 0x06, 0x53, 0x75, 0x62, 0x52, 0x65, 0x71,
	0x12, 0x2e, 0x0a, 0x09, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x43, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x52, 0x65, 0x71, 0x52, 0x09, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x52, 0x65, 0x71,
	0x12, 0x34, 0x0a, 0x0a, 0x73, 0x75, 0x62, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x53, 0x75, 0x62,
	0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0a, 0x73, 0x75, 0x62, 0x4d,
	0x61, 0x72, 0x6b, 0x65, 0x74, 0x73, 0x12, 0x31, 0x0a, 0x09, 0x73, 0x75, 0x62, 0x53, 0x74, 0x6f,
	0x63, 0x6b, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x2e, 0x53, 0x75, 0x62, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x09,
	0x73, 0x75, 0x62, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x73, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_sub_req_proto_rawDescOnce sync.Once
	file_sub_req_proto_rawDescData = file_sub_req_proto_rawDesc
)

func file_sub_req_proto_rawDescGZIP() []byte {
	file_sub_req_proto_rawDescOnce.Do(func() {
		file_sub_req_proto_rawDescData = protoimpl.X.CompressGZIP(file_sub_req_proto_rawDescData)
	})
	return file_sub_req_proto_rawDescData
}

var file_sub_req_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_sub_req_proto_goTypes = []interface{}{
	(*SubReq)(nil),        // 0: model.SubReq
	(*CommonReq)(nil),     // 1: model.CommonReq
	(*SubMarketInfo)(nil), // 2: model.SubMarketInfo
	(*SubStockInfo)(nil),  // 3: model.SubStockInfo
}
var file_sub_req_proto_depIdxs = []int32{
	1, // 0: model.SubReq.commonReq:type_name -> model.CommonReq
	2, // 1: model.SubReq.subMarkets:type_name -> model.SubMarketInfo
	3, // 2: model.SubReq.subStocks:type_name -> model.SubStockInfo
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_sub_req_proto_init() }
func file_sub_req_proto_init() {
	if File_sub_req_proto != nil {
		return
	}
	file_common_req_proto_init()
	file_sub_market_info_proto_init()
	file_sub_stock_info_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_sub_req_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubReq); i {
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
			RawDescriptor: file_sub_req_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_sub_req_proto_goTypes,
		DependencyIndexes: file_sub_req_proto_depIdxs,
		MessageInfos:      file_sub_req_proto_msgTypes,
	}.Build()
	File_sub_req_proto = out.File
	file_sub_req_proto_rawDesc = nil
	file_sub_req_proto_goTypes = nil
	file_sub_req_proto_depIdxs = nil
}