// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0-devel
// 	protoc        v3.12.4
// source: posts.proto

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

type ToggleLikeRequest_Action int32

const (
	ToggleLikeRequest_LIKE   ToggleLikeRequest_Action = 0
	ToggleLikeRequest_UNLIKE ToggleLikeRequest_Action = 1
)

// Enum value maps for ToggleLikeRequest_Action.
var (
	ToggleLikeRequest_Action_name = map[int32]string{
		0: "LIKE",
		1: "UNLIKE",
	}
	ToggleLikeRequest_Action_value = map[string]int32{
		"LIKE":   0,
		"UNLIKE": 1,
	}
)

func (x ToggleLikeRequest_Action) Enum() *ToggleLikeRequest_Action {
	p := new(ToggleLikeRequest_Action)
	*p = x
	return p
}

func (x ToggleLikeRequest_Action) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ToggleLikeRequest_Action) Descriptor() protoreflect.EnumDescriptor {
	return file_posts_proto_enumTypes[0].Descriptor()
}

func (ToggleLikeRequest_Action) Type() protoreflect.EnumType {
	return &file_posts_proto_enumTypes[0]
}

func (x ToggleLikeRequest_Action) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ToggleLikeRequest_Action.Descriptor instead.
func (ToggleLikeRequest_Action) EnumDescriptor() ([]byte, []int) {
	return file_posts_proto_rawDescGZIP(), []int{2, 0}
}

type PostsAddRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Text    string `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	PhotoId string `protobuf:"bytes,2,opt,name=photoId,proto3" json:"photoId,omitempty"`
}

func (x *PostsAddRequest) Reset() {
	*x = PostsAddRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_posts_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostsAddRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostsAddRequest) ProtoMessage() {}

func (x *PostsAddRequest) ProtoReflect() protoreflect.Message {
	mi := &file_posts_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostsAddRequest.ProtoReflect.Descriptor instead.
func (*PostsAddRequest) Descriptor() ([]byte, []int) {
	return file_posts_proto_rawDescGZIP(), []int{0}
}

func (x *PostsAddRequest) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *PostsAddRequest) GetPhotoId() string {
	if x != nil {
		return x.PhotoId
	}
	return ""
}

type PostsAddResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PostUrl string `protobuf:"bytes,1,opt,name=postUrl,proto3" json:"postUrl,omitempty"`
}

func (x *PostsAddResponse) Reset() {
	*x = PostsAddResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_posts_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostsAddResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostsAddResponse) ProtoMessage() {}

func (x *PostsAddResponse) ProtoReflect() protoreflect.Message {
	mi := &file_posts_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostsAddResponse.ProtoReflect.Descriptor instead.
func (*PostsAddResponse) Descriptor() ([]byte, []int) {
	return file_posts_proto_rawDescGZIP(), []int{1}
}

func (x *PostsAddResponse) GetPostUrl() string {
	if x != nil {
		return x.PostUrl
	}
	return ""
}

type ToggleLikeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Action ToggleLikeRequest_Action `protobuf:"varint,1,opt,name=action,proto3,enum=meme.ToggleLikeRequest_Action" json:"action,omitempty"`
	PostId string                   `protobuf:"bytes,2,opt,name=postId,proto3" json:"postId,omitempty"`
}

func (x *ToggleLikeRequest) Reset() {
	*x = ToggleLikeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_posts_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ToggleLikeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ToggleLikeRequest) ProtoMessage() {}

func (x *ToggleLikeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_posts_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ToggleLikeRequest.ProtoReflect.Descriptor instead.
func (*ToggleLikeRequest) Descriptor() ([]byte, []int) {
	return file_posts_proto_rawDescGZIP(), []int{2}
}

func (x *ToggleLikeRequest) GetAction() ToggleLikeRequest_Action {
	if x != nil {
		return x.Action
	}
	return ToggleLikeRequest_LIKE
}

func (x *ToggleLikeRequest) GetPostId() string {
	if x != nil {
		return x.PostId
	}
	return ""
}

type ToggleLikeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LikesCount int32 `protobuf:"varint,1,opt,name=likesCount,proto3" json:"likesCount,omitempty"`
}

func (x *ToggleLikeResponse) Reset() {
	*x = ToggleLikeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_posts_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ToggleLikeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ToggleLikeResponse) ProtoMessage() {}

func (x *ToggleLikeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_posts_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ToggleLikeResponse.ProtoReflect.Descriptor instead.
func (*ToggleLikeResponse) Descriptor() ([]byte, []int) {
	return file_posts_proto_rawDescGZIP(), []int{3}
}

func (x *ToggleLikeResponse) GetLikesCount() int32 {
	if x != nil {
		return x.LikesCount
	}
	return 0
}

type AddCommentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Text   string `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	PostId string `protobuf:"bytes,2,opt,name=postId,proto3" json:"postId,omitempty"`
}

func (x *AddCommentRequest) Reset() {
	*x = AddCommentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_posts_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddCommentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddCommentRequest) ProtoMessage() {}

func (x *AddCommentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_posts_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddCommentRequest.ProtoReflect.Descriptor instead.
func (*AddCommentRequest) Descriptor() ([]byte, []int) {
	return file_posts_proto_rawDescGZIP(), []int{4}
}

func (x *AddCommentRequest) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *AddCommentRequest) GetPostId() string {
	if x != nil {
		return x.PostId
	}
	return ""
}

type AddCommentResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AddCommentResponse) Reset() {
	*x = AddCommentResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_posts_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddCommentResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddCommentResponse) ProtoMessage() {}

func (x *AddCommentResponse) ProtoReflect() protoreflect.Message {
	mi := &file_posts_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddCommentResponse.ProtoReflect.Descriptor instead.
func (*AddCommentResponse) Descriptor() ([]byte, []int) {
	return file_posts_proto_rawDescGZIP(), []int{5}
}

type CommentComposerRenderer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PostId      string `protobuf:"bytes,1,opt,name=postId,proto3" json:"postId,omitempty"`
	Placeholder string `protobuf:"bytes,2,opt,name=placeholder,proto3" json:"placeholder,omitempty"`
}

func (x *CommentComposerRenderer) Reset() {
	*x = CommentComposerRenderer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_posts_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommentComposerRenderer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommentComposerRenderer) ProtoMessage() {}

func (x *CommentComposerRenderer) ProtoReflect() protoreflect.Message {
	mi := &file_posts_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommentComposerRenderer.ProtoReflect.Descriptor instead.
func (*CommentComposerRenderer) Descriptor() ([]byte, []int) {
	return file_posts_proto_rawDescGZIP(), []int{6}
}

func (x *CommentComposerRenderer) GetPostId() string {
	if x != nil {
		return x.PostId
	}
	return ""
}

func (x *CommentComposerRenderer) GetPlaceholder() string {
	if x != nil {
		return x.Placeholder
	}
	return ""
}

var File_posts_proto protoreflect.FileDescriptor

var file_posts_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x70, 0x6f, 0x73, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x6d,
	0x65, 0x6d, 0x65, 0x22, 0x3f, 0x0a, 0x0f, 0x50, 0x6f, 0x73, 0x74, 0x73, 0x41, 0x64, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x68,
	0x6f, 0x74, 0x6f, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x68, 0x6f,
	0x74, 0x6f, 0x49, 0x64, 0x22, 0x2c, 0x0a, 0x10, 0x50, 0x6f, 0x73, 0x74, 0x73, 0x41, 0x64, 0x64,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x6f, 0x73, 0x74,
	0x55, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x6f, 0x73, 0x74, 0x55,
	0x72, 0x6c, 0x22, 0x83, 0x01, 0x0a, 0x11, 0x54, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x4c, 0x69, 0x6b,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x36, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x6d, 0x65, 0x6d, 0x65, 0x2e,
	0x54, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x4c, 0x69, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x16, 0x0a, 0x06, 0x70, 0x6f, 0x73, 0x74, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x70, 0x6f, 0x73, 0x74, 0x49, 0x64, 0x22, 0x1e, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x08, 0x0a, 0x04, 0x4c, 0x49, 0x4b, 0x45, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06,
	0x55, 0x4e, 0x4c, 0x49, 0x4b, 0x45, 0x10, 0x01, 0x22, 0x34, 0x0a, 0x12, 0x54, 0x6f, 0x67, 0x67,
	0x6c, 0x65, 0x4c, 0x69, 0x6b, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1e,
	0x0a, 0x0a, 0x6c, 0x69, 0x6b, 0x65, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0a, 0x6c, 0x69, 0x6b, 0x65, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x3f,
	0x0a, 0x11, 0x41, 0x64, 0x64, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x6f, 0x73, 0x74, 0x49,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x6f, 0x73, 0x74, 0x49, 0x64, 0x22,
	0x14, 0x0a, 0x12, 0x41, 0x64, 0x64, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x53, 0x0a, 0x17, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74,
	0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x73, 0x65, 0x72, 0x52, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x65, 0x72,
	0x12, 0x16, 0x0a, 0x06, 0x70, 0x6f, 0x73, 0x74, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x70, 0x6f, 0x73, 0x74, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x70, 0x6c, 0x61, 0x63,
	0x65, 0x68, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70,
	0x6c, 0x61, 0x63, 0x65, 0x68, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x32, 0xbf, 0x01, 0x0a, 0x05, 0x50,
	0x6f, 0x73, 0x74, 0x73, 0x12, 0x34, 0x0a, 0x03, 0x41, 0x64, 0x64, 0x12, 0x15, 0x2e, 0x6d, 0x65,
	0x6d, 0x65, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x73, 0x41, 0x64, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x16, 0x2e, 0x6d, 0x65, 0x6d, 0x65, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x73, 0x41,
	0x64, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3f, 0x0a, 0x0a, 0x54, 0x6f,
	0x67, 0x67, 0x6c, 0x65, 0x4c, 0x69, 0x6b, 0x65, 0x12, 0x17, 0x2e, 0x6d, 0x65, 0x6d, 0x65, 0x2e,
	0x54, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x4c, 0x69, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x18, 0x2e, 0x6d, 0x65, 0x6d, 0x65, 0x2e, 0x54, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x4c,
	0x69, 0x6b, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3f, 0x0a, 0x0a, 0x41,
	0x64, 0x64, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x17, 0x2e, 0x6d, 0x65, 0x6d, 0x65,
	0x2e, 0x41, 0x64, 0x64, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x18, 0x2e, 0x6d, 0x65, 0x6d, 0x65, 0x2e, 0x41, 0x64, 0x64, 0x43, 0x6f, 0x6d,
	0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x05, 0x5a, 0x03,
	0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_posts_proto_rawDescOnce sync.Once
	file_posts_proto_rawDescData = file_posts_proto_rawDesc
)

func file_posts_proto_rawDescGZIP() []byte {
	file_posts_proto_rawDescOnce.Do(func() {
		file_posts_proto_rawDescData = protoimpl.X.CompressGZIP(file_posts_proto_rawDescData)
	})
	return file_posts_proto_rawDescData
}

var file_posts_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_posts_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_posts_proto_goTypes = []interface{}{
	(ToggleLikeRequest_Action)(0),   // 0: meme.ToggleLikeRequest.Action
	(*PostsAddRequest)(nil),         // 1: meme.PostsAddRequest
	(*PostsAddResponse)(nil),        // 2: meme.PostsAddResponse
	(*ToggleLikeRequest)(nil),       // 3: meme.ToggleLikeRequest
	(*ToggleLikeResponse)(nil),      // 4: meme.ToggleLikeResponse
	(*AddCommentRequest)(nil),       // 5: meme.AddCommentRequest
	(*AddCommentResponse)(nil),      // 6: meme.AddCommentResponse
	(*CommentComposerRenderer)(nil), // 7: meme.CommentComposerRenderer
}
var file_posts_proto_depIdxs = []int32{
	0, // 0: meme.ToggleLikeRequest.action:type_name -> meme.ToggleLikeRequest.Action
	1, // 1: meme.Posts.Add:input_type -> meme.PostsAddRequest
	3, // 2: meme.Posts.ToggleLike:input_type -> meme.ToggleLikeRequest
	5, // 3: meme.Posts.AddComment:input_type -> meme.AddCommentRequest
	2, // 4: meme.Posts.Add:output_type -> meme.PostsAddResponse
	4, // 5: meme.Posts.ToggleLike:output_type -> meme.ToggleLikeResponse
	6, // 6: meme.Posts.AddComment:output_type -> meme.AddCommentResponse
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_posts_proto_init() }
func file_posts_proto_init() {
	if File_posts_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_posts_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostsAddRequest); i {
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
		file_posts_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostsAddResponse); i {
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
		file_posts_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ToggleLikeRequest); i {
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
		file_posts_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ToggleLikeResponse); i {
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
		file_posts_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddCommentRequest); i {
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
		file_posts_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddCommentResponse); i {
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
		file_posts_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommentComposerRenderer); i {
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
			RawDescriptor: file_posts_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_posts_proto_goTypes,
		DependencyIndexes: file_posts_proto_depIdxs,
		EnumInfos:         file_posts_proto_enumTypes,
		MessageInfos:      file_posts_proto_msgTypes,
	}.Build()
	File_posts_proto = out.File
	file_posts_proto_rawDesc = nil
	file_posts_proto_goTypes = nil
	file_posts_proto_depIdxs = nil
}
