// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: posts.proto

package pb

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type PostsAddRequest struct {
	Text string `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
}

func (m *PostsAddRequest) Reset()         { *m = PostsAddRequest{} }
func (m *PostsAddRequest) String() string { return proto.CompactTextString(m) }
func (*PostsAddRequest) ProtoMessage()    {}
func (*PostsAddRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b14bd1586479c33d, []int{0}
}
func (m *PostsAddRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PostsAddRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PostsAddRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PostsAddRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PostsAddRequest.Merge(m, src)
}
func (m *PostsAddRequest) XXX_Size() int {
	return m.Size()
}
func (m *PostsAddRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PostsAddRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PostsAddRequest proto.InternalMessageInfo

func (m *PostsAddRequest) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

type PostsAddResponse struct {
	PostUrl string `protobuf:"bytes,1,opt,name=postUrl,proto3" json:"postUrl,omitempty"`
}

func (m *PostsAddResponse) Reset()         { *m = PostsAddResponse{} }
func (m *PostsAddResponse) String() string { return proto.CompactTextString(m) }
func (*PostsAddResponse) ProtoMessage()    {}
func (*PostsAddResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b14bd1586479c33d, []int{1}
}
func (m *PostsAddResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PostsAddResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PostsAddResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PostsAddResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PostsAddResponse.Merge(m, src)
}
func (m *PostsAddResponse) XXX_Size() int {
	return m.Size()
}
func (m *PostsAddResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PostsAddResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PostsAddResponse proto.InternalMessageInfo

func (m *PostsAddResponse) GetPostUrl() string {
	if m != nil {
		return m.PostUrl
	}
	return ""
}

func init() {
	proto.RegisterType((*PostsAddRequest)(nil), "meme.PostsAddRequest")
	proto.RegisterType((*PostsAddResponse)(nil), "meme.PostsAddResponse")
}

func init() { proto.RegisterFile("posts.proto", fileDescriptor_b14bd1586479c33d) }

var fileDescriptor_b14bd1586479c33d = []byte{
	// 166 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2e, 0xc8, 0x2f, 0x2e,
	0x29, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xc9, 0x4d, 0xcd, 0x4d, 0x55, 0x52, 0xe5,
	0xe2, 0x0f, 0x00, 0x09, 0x3a, 0xa6, 0xa4, 0x04, 0xa5, 0x16, 0x96, 0xa6, 0x16, 0x97, 0x08, 0x09,
	0x71, 0xb1, 0x94, 0xa4, 0x56, 0x94, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x81, 0xd9, 0x4a,
	0x3a, 0x5c, 0x02, 0x08, 0x65, 0xc5, 0x05, 0xf9, 0x79, 0xc5, 0xa9, 0x42, 0x12, 0x5c, 0xec, 0x20,
	0xf3, 0x42, 0x8b, 0x72, 0xa0, 0x4a, 0x61, 0x5c, 0x23, 0x5b, 0x2e, 0x56, 0xb0, 0x6a, 0x21, 0x13,
	0x2e, 0x66, 0xc7, 0x94, 0x14, 0x21, 0x51, 0x3d, 0x90, 0x5d, 0x7a, 0x68, 0x16, 0x49, 0x89, 0xa1,
	0x0b, 0x43, 0x0c, 0x76, 0x92, 0x39, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f,
	0xe4, 0x18, 0x27, 0x3c, 0x96, 0x63, 0xb8, 0xf0, 0x58, 0x8e, 0xe1, 0xc6, 0x63, 0x39, 0x86, 0x28,
	0xa6, 0x82, 0xa4, 0x24, 0x36, 0xb0, 0xf3, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0xfa, 0x1e,
	0x3d, 0xb7, 0xcd, 0x00, 0x00, 0x00,
}

func (m *PostsAddRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PostsAddRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PostsAddRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Text) > 0 {
		i -= len(m.Text)
		copy(dAtA[i:], m.Text)
		i = encodeVarintPosts(dAtA, i, uint64(len(m.Text)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *PostsAddResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PostsAddResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PostsAddResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.PostUrl) > 0 {
		i -= len(m.PostUrl)
		copy(dAtA[i:], m.PostUrl)
		i = encodeVarintPosts(dAtA, i, uint64(len(m.PostUrl)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintPosts(dAtA []byte, offset int, v uint64) int {
	offset -= sovPosts(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *PostsAddRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Text)
	if l > 0 {
		n += 1 + l + sovPosts(uint64(l))
	}
	return n
}

func (m *PostsAddResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.PostUrl)
	if l > 0 {
		n += 1 + l + sovPosts(uint64(l))
	}
	return n
}

func sovPosts(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPosts(x uint64) (n int) {
	return sovPosts(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *PostsAddRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPosts
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PostsAddRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PostsAddRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Text", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPosts
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPosts
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPosts
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Text = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPosts(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPosts
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPosts
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *PostsAddResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPosts
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PostsAddResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PostsAddResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PostUrl", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPosts
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthPosts
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPosts
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PostUrl = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPosts(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthPosts
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthPosts
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipPosts(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPosts
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowPosts
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowPosts
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthPosts
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupPosts
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthPosts
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthPosts        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPosts          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupPosts = fmt.Errorf("proto: unexpected end of group")
)