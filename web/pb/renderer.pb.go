// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0-devel
// 	protoc        v3.12.4
// source: renderer.proto

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

type UniversalRenderer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Renderer:
	//	*UniversalRenderer_ProfileRenderer
	//	*UniversalRenderer_FeedRenderer
	//	*UniversalRenderer_PostRenderer
	//	*UniversalRenderer_HeaderRenderer
	//	*UniversalRenderer_LoginPageRenderer
	//	*UniversalRenderer_SandboxRenderer
	Renderer isUniversalRenderer_Renderer `protobuf_oneof:"renderer"`
}

func (x *UniversalRenderer) Reset() {
	*x = UniversalRenderer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_renderer_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UniversalRenderer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UniversalRenderer) ProtoMessage() {}

func (x *UniversalRenderer) ProtoReflect() protoreflect.Message {
	mi := &file_renderer_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UniversalRenderer.ProtoReflect.Descriptor instead.
func (*UniversalRenderer) Descriptor() ([]byte, []int) {
	return file_renderer_proto_rawDescGZIP(), []int{0}
}

func (m *UniversalRenderer) GetRenderer() isUniversalRenderer_Renderer {
	if m != nil {
		return m.Renderer
	}
	return nil
}

func (x *UniversalRenderer) GetProfileRenderer() *ProfileRenderer {
	if x, ok := x.GetRenderer().(*UniversalRenderer_ProfileRenderer); ok {
		return x.ProfileRenderer
	}
	return nil
}

func (x *UniversalRenderer) GetFeedRenderer() *FeedRenderer {
	if x, ok := x.GetRenderer().(*UniversalRenderer_FeedRenderer); ok {
		return x.FeedRenderer
	}
	return nil
}

func (x *UniversalRenderer) GetPostRenderer() *PostRenderer {
	if x, ok := x.GetRenderer().(*UniversalRenderer_PostRenderer); ok {
		return x.PostRenderer
	}
	return nil
}

func (x *UniversalRenderer) GetHeaderRenderer() *HeaderRenderer {
	if x, ok := x.GetRenderer().(*UniversalRenderer_HeaderRenderer); ok {
		return x.HeaderRenderer
	}
	return nil
}

func (x *UniversalRenderer) GetLoginPageRenderer() *LoginPageRenderer {
	if x, ok := x.GetRenderer().(*UniversalRenderer_LoginPageRenderer); ok {
		return x.LoginPageRenderer
	}
	return nil
}

func (x *UniversalRenderer) GetSandboxRenderer() *SandboxRenderer {
	if x, ok := x.GetRenderer().(*UniversalRenderer_SandboxRenderer); ok {
		return x.SandboxRenderer
	}
	return nil
}

type isUniversalRenderer_Renderer interface {
	isUniversalRenderer_Renderer()
}

type UniversalRenderer_ProfileRenderer struct {
	ProfileRenderer *ProfileRenderer `protobuf:"bytes,1,opt,name=profileRenderer,proto3,oneof"`
}

type UniversalRenderer_FeedRenderer struct {
	FeedRenderer *FeedRenderer `protobuf:"bytes,2,opt,name=feedRenderer,proto3,oneof"`
}

type UniversalRenderer_PostRenderer struct {
	PostRenderer *PostRenderer `protobuf:"bytes,3,opt,name=postRenderer,proto3,oneof"`
}

type UniversalRenderer_HeaderRenderer struct {
	HeaderRenderer *HeaderRenderer `protobuf:"bytes,4,opt,name=headerRenderer,proto3,oneof"`
}

type UniversalRenderer_LoginPageRenderer struct {
	LoginPageRenderer *LoginPageRenderer `protobuf:"bytes,5,opt,name=loginPageRenderer,proto3,oneof"`
}

type UniversalRenderer_SandboxRenderer struct {
	SandboxRenderer *SandboxRenderer `protobuf:"bytes,6,opt,name=sandboxRenderer,proto3,oneof"`
}

func (*UniversalRenderer_ProfileRenderer) isUniversalRenderer_Renderer() {}

func (*UniversalRenderer_FeedRenderer) isUniversalRenderer_Renderer() {}

func (*UniversalRenderer_PostRenderer) isUniversalRenderer_Renderer() {}

func (*UniversalRenderer_HeaderRenderer) isUniversalRenderer_Renderer() {}

func (*UniversalRenderer_LoginPageRenderer) isUniversalRenderer_Renderer() {}

func (*UniversalRenderer_SandboxRenderer) isUniversalRenderer_Renderer() {}

type ResolveRouteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *ResolveRouteRequest) Reset() {
	*x = ResolveRouteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_renderer_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResolveRouteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResolveRouteRequest) ProtoMessage() {}

func (x *ResolveRouteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_renderer_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResolveRouteRequest.ProtoReflect.Descriptor instead.
func (*ResolveRouteRequest) Descriptor() ([]byte, []int) {
	return file_renderer_proto_rawDescGZIP(), []int{1}
}

func (x *ResolveRouteRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type SandboxRenderer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SandboxRenderer) Reset() {
	*x = SandboxRenderer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_renderer_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SandboxRenderer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SandboxRenderer) ProtoMessage() {}

func (x *SandboxRenderer) ProtoReflect() protoreflect.Message {
	mi := &file_renderer_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SandboxRenderer.ProtoReflect.Descriptor instead.
func (*SandboxRenderer) Descriptor() ([]byte, []int) {
	return file_renderer_proto_rawDescGZIP(), []int{2}
}

var File_renderer_proto protoreflect.FileDescriptor

var file_renderer_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x04, 0x6d, 0x65, 0x6d, 0x65, 0x1a, 0x0a, 0x61, 0x70, 0x69, 0x32, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x0b, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x0b, 0x70, 0x6f, 0x73, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa2, 0x03, 0x0a,
	0x11, 0x55, 0x6e, 0x69, 0x76, 0x65, 0x72, 0x73, 0x61, 0x6c, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x65, 0x72, 0x12, 0x41, 0x0a, 0x0f, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x6d, 0x65,
	0x6d, 0x65, 0x2e, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x65, 0x72, 0x48, 0x00, 0x52, 0x0f, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x65, 0x72, 0x12, 0x38, 0x0a, 0x0c, 0x66, 0x65, 0x65, 0x64, 0x52, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6d, 0x65,
	0x6d, 0x65, 0x2e, 0x46, 0x65, 0x65, 0x64, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x48,
	0x00, 0x52, 0x0c, 0x66, 0x65, 0x65, 0x64, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x12,
	0x38, 0x0a, 0x0c, 0x70, 0x6f, 0x73, 0x74, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6d, 0x65, 0x6d, 0x65, 0x2e, 0x50, 0x6f, 0x73,
	0x74, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x48, 0x00, 0x52, 0x0c, 0x70, 0x6f, 0x73,
	0x74, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x12, 0x3e, 0x0a, 0x0e, 0x68, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x14, 0x2e, 0x6d, 0x65, 0x6d, 0x65, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x48, 0x00, 0x52, 0x0e, 0x68, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x12, 0x47, 0x0a, 0x11, 0x6c, 0x6f, 0x67,
	0x69, 0x6e, 0x50, 0x61, 0x67, 0x65, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x6d, 0x65, 0x6d, 0x65, 0x2e, 0x4c, 0x6f, 0x67, 0x69,
	0x6e, 0x50, 0x61, 0x67, 0x65, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x48, 0x00, 0x52,
	0x11, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x50, 0x61, 0x67, 0x65, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x65, 0x72, 0x12, 0x41, 0x0a, 0x0f, 0x73, 0x61, 0x6e, 0x64, 0x62, 0x6f, 0x78, 0x52, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x6d, 0x65,
	0x6d, 0x65, 0x2e, 0x53, 0x61, 0x6e, 0x64, 0x62, 0x6f, 0x78, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x65, 0x72, 0x48, 0x00, 0x52, 0x0f, 0x73, 0x61, 0x6e, 0x64, 0x62, 0x6f, 0x78, 0x52, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x65, 0x72, 0x42, 0x0a, 0x0a, 0x08, 0x72, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65,
	0x72, 0x22, 0x27, 0x0a, 0x13, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76, 0x65, 0x52, 0x6f, 0x75, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x11, 0x0a, 0x0f, 0x53, 0x61,
	0x6e, 0x64, 0x62, 0x6f, 0x78, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x32, 0x4b, 0x0a,
	0x05, 0x55, 0x74, 0x69, 0x6c, 0x73, 0x12, 0x42, 0x0a, 0x0c, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x76,
	0x65, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x12, 0x19, 0x2e, 0x6d, 0x65, 0x6d, 0x65, 0x2e, 0x52, 0x65,
	0x73, 0x6f, 0x6c, 0x76, 0x65, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x17, 0x2e, 0x6d, 0x65, 0x6d, 0x65, 0x2e, 0x55, 0x6e, 0x69, 0x76, 0x65, 0x72, 0x73,
	0x61, 0x6c, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x72, 0x42, 0x05, 0x5a, 0x03, 0x2f, 0x70,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_renderer_proto_rawDescOnce sync.Once
	file_renderer_proto_rawDescData = file_renderer_proto_rawDesc
)

func file_renderer_proto_rawDescGZIP() []byte {
	file_renderer_proto_rawDescOnce.Do(func() {
		file_renderer_proto_rawDescData = protoimpl.X.CompressGZIP(file_renderer_proto_rawDescData)
	})
	return file_renderer_proto_rawDescData
}

var file_renderer_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_renderer_proto_goTypes = []interface{}{
	(*UniversalRenderer)(nil),   // 0: meme.UniversalRenderer
	(*ResolveRouteRequest)(nil), // 1: meme.ResolveRouteRequest
	(*SandboxRenderer)(nil),     // 2: meme.SandboxRenderer
	(*ProfileRenderer)(nil),     // 3: meme.ProfileRenderer
	(*FeedRenderer)(nil),        // 4: meme.FeedRenderer
	(*PostRenderer)(nil),        // 5: meme.PostRenderer
	(*HeaderRenderer)(nil),      // 6: meme.HeaderRenderer
	(*LoginPageRenderer)(nil),   // 7: meme.LoginPageRenderer
}
var file_renderer_proto_depIdxs = []int32{
	3, // 0: meme.UniversalRenderer.profileRenderer:type_name -> meme.ProfileRenderer
	4, // 1: meme.UniversalRenderer.feedRenderer:type_name -> meme.FeedRenderer
	5, // 2: meme.UniversalRenderer.postRenderer:type_name -> meme.PostRenderer
	6, // 3: meme.UniversalRenderer.headerRenderer:type_name -> meme.HeaderRenderer
	7, // 4: meme.UniversalRenderer.loginPageRenderer:type_name -> meme.LoginPageRenderer
	2, // 5: meme.UniversalRenderer.sandboxRenderer:type_name -> meme.SandboxRenderer
	1, // 6: meme.Utils.ResolveRoute:input_type -> meme.ResolveRouteRequest
	0, // 7: meme.Utils.ResolveRoute:output_type -> meme.UniversalRenderer
	7, // [7:8] is the sub-list for method output_type
	6, // [6:7] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_renderer_proto_init() }
func file_renderer_proto_init() {
	if File_renderer_proto != nil {
		return
	}
	file_api2_proto_init()
	file_login_proto_init()
	file_posts_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_renderer_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UniversalRenderer); i {
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
		file_renderer_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResolveRouteRequest); i {
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
		file_renderer_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SandboxRenderer); i {
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
	file_renderer_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*UniversalRenderer_ProfileRenderer)(nil),
		(*UniversalRenderer_FeedRenderer)(nil),
		(*UniversalRenderer_PostRenderer)(nil),
		(*UniversalRenderer_HeaderRenderer)(nil),
		(*UniversalRenderer_LoginPageRenderer)(nil),
		(*UniversalRenderer_SandboxRenderer)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_renderer_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_renderer_proto_goTypes,
		DependencyIndexes: file_renderer_proto_depIdxs,
		MessageInfos:      file_renderer_proto_msgTypes,
	}.Build()
	File_renderer_proto = out.File
	file_renderer_proto_rawDesc = nil
	file_renderer_proto_goTypes = nil
	file_renderer_proto_depIdxs = nil
}
