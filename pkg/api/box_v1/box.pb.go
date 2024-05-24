// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.1
// source: box.proto

package box_v1

import (
	abstract "Homework-1/pkg/api/abstract"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Box struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Cost    float64 `protobuf:"fixed64,2,opt,name=cost,proto3" json:"cost,omitempty"`
	IsCheck bool    `protobuf:"varint,3,opt,name=isCheck,proto3" json:"isCheck,omitempty"`
	Weight  float64 `protobuf:"fixed64,4,opt,name=weight,proto3" json:"weight,omitempty"`
}

func (x *Box) Reset() {
	*x = Box{}
	if protoimpl.UnsafeEnabled {
		mi := &file_box_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Box) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Box) ProtoMessage() {}

func (x *Box) ProtoReflect() protoreflect.Message {
	mi := &file_box_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Box.ProtoReflect.Descriptor instead.
func (*Box) Descriptor() ([]byte, []int) {
	return file_box_proto_rawDescGZIP(), []int{0}
}

func (x *Box) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Box) GetCost() float64 {
	if x != nil {
		return x.Cost
	}
	return 0
}

func (x *Box) GetIsCheck() bool {
	if x != nil {
		return x.IsCheck
	}
	return false
}

func (x *Box) GetWeight() float64 {
	if x != nil {
		return x.Weight
	}
	return 0
}

type BoxAllInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID        int64                  `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Box       *Box                   `protobuf:"bytes,2,opt,name=box,proto3" json:"box,omitempty"`
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=updatedAt,proto3" json:"updatedAt,omitempty"`
}

func (x *BoxAllInfo) Reset() {
	*x = BoxAllInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_box_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BoxAllInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BoxAllInfo) ProtoMessage() {}

func (x *BoxAllInfo) ProtoReflect() protoreflect.Message {
	mi := &file_box_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BoxAllInfo.ProtoReflect.Descriptor instead.
func (*BoxAllInfo) Descriptor() ([]byte, []int) {
	return file_box_proto_rawDescGZIP(), []int{1}
}

func (x *BoxAllInfo) GetID() int64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *BoxAllInfo) GetBox() *Box {
	if x != nil {
		return x.Box
	}
	return nil
}

func (x *BoxAllInfo) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *BoxAllInfo) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

type BoxCreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Box *Box `protobuf:"bytes,1,opt,name=box,proto3" json:"box,omitempty"`
}

func (x *BoxCreateRequest) Reset() {
	*x = BoxCreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_box_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BoxCreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BoxCreateRequest) ProtoMessage() {}

func (x *BoxCreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_box_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BoxCreateRequest.ProtoReflect.Descriptor instead.
func (*BoxCreateRequest) Descriptor() ([]byte, []int) {
	return file_box_proto_rawDescGZIP(), []int{2}
}

func (x *BoxCreateRequest) GetBox() *Box {
	if x != nil {
		return x.Box
	}
	return nil
}

type BoxIDRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BoxID int64 `protobuf:"varint,1,opt,name=boxID,proto3" json:"boxID,omitempty"`
}

func (x *BoxIDRequest) Reset() {
	*x = BoxIDRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_box_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BoxIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BoxIDRequest) ProtoMessage() {}

func (x *BoxIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_box_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BoxIDRequest.ProtoReflect.Descriptor instead.
func (*BoxIDRequest) Descriptor() ([]byte, []int) {
	return file_box_proto_rawDescGZIP(), []int{3}
}

func (x *BoxIDRequest) GetBoxID() int64 {
	if x != nil {
		return x.BoxID
	}
	return 0
}

type BoxListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BoxAllInfo []*BoxAllInfo        `protobuf:"bytes,1,rep,name=boxAllInfo,proto3" json:"boxAllInfo,omitempty"`
	Pagination *abstract.Pagination `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (x *BoxListResponse) Reset() {
	*x = BoxListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_box_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BoxListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BoxListResponse) ProtoMessage() {}

func (x *BoxListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_box_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BoxListResponse.ProtoReflect.Descriptor instead.
func (*BoxListResponse) Descriptor() ([]byte, []int) {
	return file_box_proto_rawDescGZIP(), []int{4}
}

func (x *BoxListResponse) GetBoxAllInfo() []*BoxAllInfo {
	if x != nil {
		return x.BoxAllInfo
	}
	return nil
}

func (x *BoxListResponse) GetPagination() *abstract.Pagination {
	if x != nil {
		return x.Pagination
	}
	return nil
}

var File_box_proto protoreflect.FileDescriptor

var file_box_proto_rawDesc = []byte{
	0x0a, 0x09, 0x62, 0x6f, 0x78, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0e, 0x61, 0x62, 0x73,
	0x74, 0x72, 0x61, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5f, 0x0a, 0x03, 0x42, 0x6f,
	0x78, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x73, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x04, 0x63, 0x6f, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x69, 0x73, 0x43,
	0x68, 0x65, 0x63, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x69, 0x73, 0x43, 0x68,
	0x65, 0x63, 0x6b, 0x12, 0x16, 0x0a, 0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x06, 0x77, 0x65, 0x69, 0x67, 0x68, 0x74, 0x22, 0xa8, 0x01, 0x0a, 0x0a,
	0x42, 0x6f, 0x78, 0x41, 0x6c, 0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x03, 0x62, 0x6f,
	0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x04, 0x2e, 0x42, 0x6f, 0x78, 0x52, 0x03, 0x62,
	0x6f, 0x78, 0x12, 0x38, 0x0a, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x38, 0x0a, 0x09,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x2a, 0x0a, 0x10, 0x42, 0x6f, 0x78, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x03, 0x62, 0x6f,
	0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x04, 0x2e, 0x42, 0x6f, 0x78, 0x52, 0x03, 0x62,
	0x6f, 0x78, 0x22, 0x24, 0x0a, 0x0c, 0x42, 0x6f, 0x78, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x62, 0x6f, 0x78, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x05, 0x62, 0x6f, 0x78, 0x49, 0x44, 0x22, 0x6b, 0x0a, 0x0f, 0x42, 0x6f, 0x78, 0x4c,
	0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x0a, 0x62,
	0x6f, 0x78, 0x41, 0x6c, 0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x0b, 0x2e, 0x42, 0x6f, 0x78, 0x41, 0x6c, 0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0a, 0x62, 0x6f,
	0x78, 0x41, 0x6c, 0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x2b, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x69,
	0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x50,
	0x61, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x69, 0x6e,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x32, 0xad, 0x02, 0x0a, 0x0a, 0x42, 0x6f, 0x78, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x4b, 0x0a, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x42, 0x6f,
	0x78, 0x12, 0x11, 0x2e, 0x42, 0x6f, 0x78, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x19, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x13, 0x3a, 0x01,
	0x2a, 0x22, 0x0e, 0x2f, 0x62, 0x6f, 0x78, 0x5f, 0x76, 0x31, 0x2f, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x12, 0x4c, 0x0a, 0x09, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x42, 0x6f, 0x78, 0x12, 0x0d,
	0x2e, 0x42, 0x6f, 0x78, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x2a, 0x16, 0x2f, 0x62, 0x6f, 0x78, 0x5f, 0x76, 0x31,
	0x2f, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x2f, 0x7b, 0x62, 0x6f, 0x78, 0x49, 0x44, 0x7d, 0x12,
	0x3d, 0x0a, 0x09, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x6f, 0x78, 0x65, 0x73, 0x12, 0x05, 0x2e, 0x50,
	0x61, 0x67, 0x65, 0x1a, 0x10, 0x2e, 0x42, 0x6f, 0x78, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x17, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x11, 0x3a, 0x01, 0x2a,
	0x22, 0x0c, 0x2f, 0x62, 0x6f, 0x78, 0x5f, 0x76, 0x31, 0x2f, 0x6c, 0x69, 0x73, 0x74, 0x12, 0x45,
	0x0a, 0x0a, 0x47, 0x65, 0x74, 0x42, 0x6f, 0x78, 0x42, 0x79, 0x49, 0x44, 0x12, 0x0d, 0x2e, 0x42,
	0x6f, 0x78, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0b, 0x2e, 0x42, 0x6f,
	0x78, 0x41, 0x6c, 0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15,
	0x12, 0x13, 0x2f, 0x62, 0x6f, 0x78, 0x5f, 0x76, 0x31, 0x2f, 0x67, 0x65, 0x74, 0x2f, 0x7b, 0x62,
	0x6f, 0x78, 0x49, 0x44, 0x7d, 0x42, 0x20, 0x5a, 0x1e, 0x63, 0x72, 0x75, 0x64, 0x2d, 0x70, 0x76,
	0x7a, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x62, 0x6f, 0x78, 0x5f, 0x76, 0x31,
	0x3b, 0x62, 0x6f, 0x78, 0x5f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_box_proto_rawDescOnce sync.Once
	file_box_proto_rawDescData = file_box_proto_rawDesc
)

func file_box_proto_rawDescGZIP() []byte {
	file_box_proto_rawDescOnce.Do(func() {
		file_box_proto_rawDescData = protoimpl.X.CompressGZIP(file_box_proto_rawDescData)
	})
	return file_box_proto_rawDescData
}

var file_box_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_box_proto_goTypes = []interface{}{
	(*Box)(nil),                      // 0: Box
	(*BoxAllInfo)(nil),               // 1: BoxAllInfo
	(*BoxCreateRequest)(nil),         // 2: BoxCreateRequest
	(*BoxIDRequest)(nil),             // 3: BoxIDRequest
	(*BoxListResponse)(nil),          // 4: BoxListResponse
	(*timestamppb.Timestamp)(nil),    // 5: google.protobuf.Timestamp
	(*abstract.Pagination)(nil),      // 6: Pagination
	(*abstract.Page)(nil),            // 7: Page
	(*abstract.MessageResponse)(nil), // 8: MessageResponse
}
var file_box_proto_depIdxs = []int32{
	0,  // 0: BoxAllInfo.box:type_name -> Box
	5,  // 1: BoxAllInfo.createdAt:type_name -> google.protobuf.Timestamp
	5,  // 2: BoxAllInfo.updatedAt:type_name -> google.protobuf.Timestamp
	0,  // 3: BoxCreateRequest.box:type_name -> Box
	1,  // 4: BoxListResponse.boxAllInfo:type_name -> BoxAllInfo
	6,  // 5: BoxListResponse.pagination:type_name -> Pagination
	2,  // 6: BoxService.CreateBox:input_type -> BoxCreateRequest
	3,  // 7: BoxService.DeleteBox:input_type -> BoxIDRequest
	7,  // 8: BoxService.ListBoxes:input_type -> Page
	3,  // 9: BoxService.GetBoxByID:input_type -> BoxIDRequest
	8,  // 10: BoxService.CreateBox:output_type -> MessageResponse
	8,  // 11: BoxService.DeleteBox:output_type -> MessageResponse
	4,  // 12: BoxService.ListBoxes:output_type -> BoxListResponse
	1,  // 13: BoxService.GetBoxByID:output_type -> BoxAllInfo
	10, // [10:14] is the sub-list for method output_type
	6,  // [6:10] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_box_proto_init() }
func file_box_proto_init() {
	if File_box_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_box_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Box); i {
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
		file_box_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BoxAllInfo); i {
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
		file_box_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BoxCreateRequest); i {
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
		file_box_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BoxIDRequest); i {
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
		file_box_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BoxListResponse); i {
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
			RawDescriptor: file_box_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_box_proto_goTypes,
		DependencyIndexes: file_box_proto_depIdxs,
		MessageInfos:      file_box_proto_msgTypes,
	}.Build()
	File_box_proto = out.File
	file_box_proto_rawDesc = nil
	file_box_proto_goTypes = nil
	file_box_proto_depIdxs = nil
}
