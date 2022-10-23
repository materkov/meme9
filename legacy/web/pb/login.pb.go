// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0-devel
// 	protoc        v3.12.4
// source: login.proto

package pb

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

type LoginPageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Renderer *LoginPageRenderer `protobuf:"bytes,1,opt,name=renderer,proto3" json:"renderer,omitempty"`
}

func (x *LoginPageResponse) Reset() {
	*x = LoginPageResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_login_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginPageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginPageResponse) ProtoMessage() {}

func (x *LoginPageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_login_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginPageResponse.ProtoReflect.Descriptor instead.
func (*LoginPageResponse) Descriptor() ([]byte, []int) {
	return file_login_proto_rawDescGZIP(), []int{0}
}

func (x *LoginPageResponse) GetRenderer() *LoginPageRenderer {
	if x != nil {
		return x.Renderer
	}
	return nil
}

type LoginPageRenderer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AuthUrl string `protobuf:"bytes,1,opt,name=authUrl,proto3" json:"authUrl,omitempty"`
	Text    string `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
}

func (x *LoginPageRenderer) Reset() {
	*x = LoginPageRenderer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_login_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginPageRenderer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginPageRenderer) ProtoMessage() {}

func (x *LoginPageRenderer) ProtoReflect() protoreflect.Message {
	mi := &file_login_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginPageRenderer.ProtoReflect.Descriptor instead.
func (*LoginPageRenderer) Descriptor() ([]byte, []int) {
	return file_login_proto_rawDescGZIP(), []int{1}
}

func (x *LoginPageRenderer) GetAuthUrl() string {
	if x != nil {
		return x.AuthUrl
	}
	return ""
}

func (x *LoginPageRenderer) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

var File_login_proto protoreflect.FileDescriptor

var file_login_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x6d,
	0x65, 0x6d, 0x65, 0x22, 0x48, 0x0a, 0x11, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x50, 0x61, 0x67, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x33, 0x0a, 0x08, 0x72, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x6d, 0x65, 0x6d,
	0x65, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x50, 0x61, 0x67, 0x65, 0x52, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x65, 0x72, 0x52, 0x08, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x22, 0x41, 0x0a,
	0x11, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x50, 0x61, 0x67, 0x65, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x65, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x75, 0x74, 0x68, 0x55, 0x72, 0x6c, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x75, 0x74, 0x68, 0x55, 0x72, 0x6c, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74,
	0x42, 0x05, 0x5a, 0x03, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_login_proto_rawDescOnce sync.Once
	file_login_proto_rawDescData = file_login_proto_rawDesc
)

func file_login_proto_rawDescGZIP() []byte {
	file_login_proto_rawDescOnce.Do(func() {
		file_login_proto_rawDescData = protoimpl.X.CompressGZIP(file_login_proto_rawDescData)
	})
	return file_login_proto_rawDescData
}

var file_login_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_login_proto_goTypes = []interface{}{
	(*LoginPageResponse)(nil), // 0: meme.LoginPageResponse
	(*LoginPageRenderer)(nil), // 1: meme.LoginPageRenderer
}
var file_login_proto_depIdxs = []int32{
	1, // 0: meme.LoginPageResponse.renderer:type_name -> meme.LoginPageRenderer
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_login_proto_init() }
func file_login_proto_init() {
	if File_login_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_login_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginPageResponse); i {
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
		file_login_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginPageRenderer); i {
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
			RawDescriptor: file_login_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_login_proto_goTypes,
		DependencyIndexes: file_login_proto_depIdxs,
		MessageInfos:      file_login_proto_msgTypes,
	}.Build()
	File_login_proto = out.File
	file_login_proto_rawDesc = nil
	file_login_proto_goTypes = nil
	file_login_proto_depIdxs = nil
}