// Code generated by protoc-gen-go. DO NOT EDIT.
// source: blog.proto

package blog

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Blog struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	AuthorId             string   `protobuf:"bytes,2,opt,name=author_id,json=authorId,proto3" json:"author_id,omitempty"`
	Title                string   `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	Body                 string   `protobuf:"bytes,4,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Blog) Reset()         { *m = Blog{} }
func (m *Blog) String() string { return proto.CompactTextString(m) }
func (*Blog) ProtoMessage()    {}
func (*Blog) Descriptor() ([]byte, []int) {
	return fileDescriptor_6745b25902462fb1, []int{0}
}

func (m *Blog) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Blog.Unmarshal(m, b)
}
func (m *Blog) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Blog.Marshal(b, m, deterministic)
}
func (m *Blog) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Blog.Merge(m, src)
}
func (m *Blog) XXX_Size() int {
	return xxx_messageInfo_Blog.Size(m)
}
func (m *Blog) XXX_DiscardUnknown() {
	xxx_messageInfo_Blog.DiscardUnknown(m)
}

var xxx_messageInfo_Blog proto.InternalMessageInfo

func (m *Blog) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Blog) GetAuthorId() string {
	if m != nil {
		return m.AuthorId
	}
	return ""
}

func (m *Blog) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Blog) GetBody() string {
	if m != nil {
		return m.Body
	}
	return ""
}

func init() {
	proto.RegisterType((*Blog)(nil), "Blog")
}

func init() {
	proto.RegisterFile("blog.proto", fileDescriptor_6745b25902462fb1)
}

var fileDescriptor_6745b25902462fb1 = []byte{
	// 125 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0xca, 0xc9, 0x4f,
	0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0x8a, 0xe5, 0x62, 0x71, 0x02, 0xf2, 0x84, 0xf8, 0xb8,
	0x98, 0x32, 0x53, 0x24, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0x80, 0x2c, 0x21, 0x69, 0x2e, 0xce,
	0xc4, 0xd2, 0x92, 0x8c, 0xfc, 0xa2, 0x78, 0xa0, 0x30, 0x13, 0x58, 0x98, 0x03, 0x22, 0xe0, 0x99,
	0x22, 0x24, 0xc2, 0xc5, 0x5a, 0x92, 0x59, 0x92, 0x93, 0x2a, 0xc1, 0x0c, 0x96, 0x80, 0x70, 0x84,
	0x84, 0xb8, 0x58, 0x92, 0xf2, 0x53, 0x2a, 0x25, 0x58, 0xc0, 0x82, 0x60, 0xb6, 0x11, 0x2f, 0x17,
	0x37, 0xc8, 0xf8, 0xe0, 0xd4, 0xa2, 0xb2, 0xcc, 0xe4, 0xd4, 0x24, 0x36, 0xb0, 0xa5, 0xc6, 0x80,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x27, 0xd8, 0x52, 0xf8, 0x82, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// BlogServiceClient is the client API for BlogService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BlogServiceClient interface {
}

type blogServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBlogServiceClient(cc grpc.ClientConnInterface) BlogServiceClient {
	return &blogServiceClient{cc}
}

// BlogServiceServer is the server API for BlogService service.
type BlogServiceServer interface {
}

// UnimplementedBlogServiceServer can be embedded to have forward compatible implementations.
type UnimplementedBlogServiceServer struct {
}

func RegisterBlogServiceServer(s *grpc.Server, srv BlogServiceServer) {
	s.RegisterService(&_BlogService_serviceDesc, srv)
}

var _BlogService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "BlogService",
	HandlerType: (*BlogServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
	Metadata:    "blog.proto",
}
