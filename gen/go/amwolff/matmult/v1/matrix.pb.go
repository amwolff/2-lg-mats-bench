// Code generated by protoc-gen-go. DO NOT EDIT.
// source: amwolff/matmult/v1/matrix.proto

package matmultv1

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

// Matrix is a building block of MatrixProductAPI request/response.
type Matrix struct {
	Columns              []*Matrix_Column `protobuf:"bytes,1,rep,name=columns,proto3" json:"columns,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Matrix) Reset()         { *m = Matrix{} }
func (m *Matrix) String() string { return proto.CompactTextString(m) }
func (*Matrix) ProtoMessage()    {}
func (*Matrix) Descriptor() ([]byte, []int) {
	return fileDescriptor_151b85c330b84a14, []int{0}
}

func (m *Matrix) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Matrix.Unmarshal(m, b)
}
func (m *Matrix) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Matrix.Marshal(b, m, deterministic)
}
func (m *Matrix) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Matrix.Merge(m, src)
}
func (m *Matrix) XXX_Size() int {
	return xxx_messageInfo_Matrix.Size(m)
}
func (m *Matrix) XXX_DiscardUnknown() {
	xxx_messageInfo_Matrix.DiscardUnknown(m)
}

var xxx_messageInfo_Matrix proto.InternalMessageInfo

func (m *Matrix) GetColumns() []*Matrix_Column {
	if m != nil {
		return m.Columns
	}
	return nil
}

// Eigen uses column-major order to store its matrices:
// https://eigen.tuxfamily.org/dox/group__TopicStorageOrders.html
// I think it's fair to adjust to Eigen.
type Matrix_Column struct {
	Coefficients         []float64 `protobuf:"fixed64,1,rep,packed,name=coefficients,proto3" json:"coefficients,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Matrix_Column) Reset()         { *m = Matrix_Column{} }
func (m *Matrix_Column) String() string { return proto.CompactTextString(m) }
func (*Matrix_Column) ProtoMessage()    {}
func (*Matrix_Column) Descriptor() ([]byte, []int) {
	return fileDescriptor_151b85c330b84a14, []int{0, 0}
}

func (m *Matrix_Column) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Matrix_Column.Unmarshal(m, b)
}
func (m *Matrix_Column) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Matrix_Column.Marshal(b, m, deterministic)
}
func (m *Matrix_Column) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Matrix_Column.Merge(m, src)
}
func (m *Matrix_Column) XXX_Size() int {
	return xxx_messageInfo_Matrix_Column.Size(m)
}
func (m *Matrix_Column) XXX_DiscardUnknown() {
	xxx_messageInfo_Matrix_Column.DiscardUnknown(m)
}

var xxx_messageInfo_Matrix_Column proto.InternalMessageInfo

func (m *Matrix_Column) GetCoefficients() []float64 {
	if m != nil {
		return m.Coefficients
	}
	return nil
}

func init() {
	proto.RegisterType((*Matrix)(nil), "amwolff.matmult.v1.Matrix")
	proto.RegisterType((*Matrix_Column)(nil), "amwolff.matmult.v1.Matrix.Column")
}

func init() { proto.RegisterFile("amwolff/matmult/v1/matrix.proto", fileDescriptor_151b85c330b84a14) }

var fileDescriptor_151b85c330b84a14 = []byte{
	// 193 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4f, 0xcc, 0x2d, 0xcf,
	0xcf, 0x49, 0x4b, 0xd3, 0xcf, 0x4d, 0x2c, 0xc9, 0x2d, 0xcd, 0x29, 0xd1, 0x2f, 0x33, 0x04, 0x31,
	0x8b, 0x32, 0x2b, 0xf4, 0x0a, 0x8a, 0xf2, 0x4b, 0xf2, 0x85, 0x84, 0xa0, 0x0a, 0xf4, 0xa0, 0x0a,
	0xf4, 0xca, 0x0c, 0x95, 0xca, 0xb9, 0xd8, 0x7c, 0xc1, 0x6a, 0x84, 0xac, 0xb9, 0xd8, 0x93, 0xf3,
	0x73, 0x4a, 0x73, 0xf3, 0x8a, 0x25, 0x18, 0x15, 0x98, 0x35, 0xb8, 0x8d, 0x14, 0xf5, 0x30, 0xd5,
	0xeb, 0x41, 0x14, 0xeb, 0x39, 0x83, 0x55, 0x06, 0xc1, 0x74, 0x48, 0x19, 0x70, 0xb1, 0x41, 0x84,
	0x84, 0xd4, 0xb8, 0x78, 0x92, 0xf3, 0x53, 0xd3, 0xd2, 0x32, 0x93, 0x33, 0x53, 0xf3, 0x4a, 0x20,
	0x66, 0x31, 0x3a, 0x31, 0x09, 0x30, 0x06, 0xa1, 0x88, 0x3b, 0xa5, 0x70, 0x89, 0x25, 0xe7, 0xe7,
	0x62, 0xb1, 0xc2, 0x89, 0x1b, 0x62, 0x47, 0x00, 0xc8, 0xcd, 0x1e, 0x8c, 0x01, 0x8c, 0x51, 0x9c,
	0x50, 0xc9, 0x32, 0xc3, 0x45, 0x4c, 0xcc, 0x8e, 0xbe, 0x11, 0xab, 0x98, 0x84, 0x1c, 0xa1, 0xda,
	0x7c, 0xa1, 0xda, 0xc2, 0x0c, 0x4f, 0xc1, 0x05, 0x63, 0xa0, 0x82, 0x31, 0x61, 0x86, 0x49, 0x6c,
	0x60, 0x9f, 0x1b, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xc6, 0x8a, 0x31, 0x8e, 0x1c, 0x01, 0x00,
	0x00,
}