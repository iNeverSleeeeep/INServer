// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: world-database.proto

package msg

import (
	data "INServer/src/proto/data"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type LoadStaticMapReq struct {
	ZoneID               int32    `protobuf:"varint,1,opt,name=ZoneID,proto3" json:"ZoneID,omitempty"`
	StaticMapID          int32    `protobuf:"varint,2,opt,name=StaticMapID,proto3" json:"StaticMapID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoadStaticMapReq) Reset()         { *m = LoadStaticMapReq{} }
func (m *LoadStaticMapReq) String() string { return proto.CompactTextString(m) }
func (*LoadStaticMapReq) ProtoMessage()    {}
func (*LoadStaticMapReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_b07d4da67d0ecb77, []int{0}
}
func (m *LoadStaticMapReq) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LoadStaticMapReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LoadStaticMapReq.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LoadStaticMapReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoadStaticMapReq.Merge(m, src)
}
func (m *LoadStaticMapReq) XXX_Size() int {
	return m.Size()
}
func (m *LoadStaticMapReq) XXX_DiscardUnknown() {
	xxx_messageInfo_LoadStaticMapReq.DiscardUnknown(m)
}

var xxx_messageInfo_LoadStaticMapReq proto.InternalMessageInfo

func (m *LoadStaticMapReq) GetZoneID() int32 {
	if m != nil {
		return m.ZoneID
	}
	return 0
}

func (m *LoadStaticMapReq) GetStaticMapID() int32 {
	if m != nil {
		return m.StaticMapID
	}
	return 0
}

type LoadStaticMapResp struct {
	Map                  *data.MapData `protobuf:"bytes,1,opt,name=Map,proto3" json:"Map,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *LoadStaticMapResp) Reset()         { *m = LoadStaticMapResp{} }
func (m *LoadStaticMapResp) String() string { return proto.CompactTextString(m) }
func (*LoadStaticMapResp) ProtoMessage()    {}
func (*LoadStaticMapResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_b07d4da67d0ecb77, []int{1}
}
func (m *LoadStaticMapResp) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LoadStaticMapResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LoadStaticMapResp.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LoadStaticMapResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoadStaticMapResp.Merge(m, src)
}
func (m *LoadStaticMapResp) XXX_Size() int {
	return m.Size()
}
func (m *LoadStaticMapResp) XXX_DiscardUnknown() {
	xxx_messageInfo_LoadStaticMapResp.DiscardUnknown(m)
}

var xxx_messageInfo_LoadStaticMapResp proto.InternalMessageInfo

func (m *LoadStaticMapResp) GetMap() *data.MapData {
	if m != nil {
		return m.Map
	}
	return nil
}

func init() {
	proto.RegisterType((*LoadStaticMapReq)(nil), "LoadStaticMapReq")
	proto.RegisterType((*LoadStaticMapResp)(nil), "LoadStaticMapResp")
}

func init() { proto.RegisterFile("world-database.proto", fileDescriptor_b07d4da67d0ecb77) }

var fileDescriptor_b07d4da67d0ecb77 = []byte{
	// 182 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x29, 0xcf, 0x2f, 0xca,
	0x49, 0xd1, 0x4d, 0x49, 0x2c, 0x49, 0x4c, 0x4a, 0x2c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x97, 0xe2, 0xcc, 0x4d, 0x2c, 0x80, 0x30, 0x95, 0x7c, 0xb8, 0x04, 0x7c, 0xf2, 0x13, 0x53, 0x82,
	0x4b, 0x12, 0x4b, 0x32, 0x93, 0x7d, 0x13, 0x0b, 0x82, 0x52, 0x0b, 0x85, 0xc4, 0xb8, 0xd8, 0xa2,
	0xf2, 0xf3, 0x52, 0x3d, 0x5d, 0x24, 0x18, 0x15, 0x18, 0x35, 0x58, 0x83, 0xa0, 0x3c, 0x21, 0x05,
	0x2e, 0x6e, 0xb8, 0x3a, 0x4f, 0x17, 0x09, 0x26, 0xb0, 0x24, 0xb2, 0x90, 0x92, 0x3e, 0x97, 0x20,
	0x9a, 0x69, 0xc5, 0x05, 0x42, 0x52, 0x5c, 0xcc, 0xbe, 0x89, 0x05, 0x60, 0xb3, 0xb8, 0x8d, 0x38,
	0xf4, 0x7c, 0x13, 0x0b, 0x5c, 0x12, 0x4b, 0x12, 0x83, 0x40, 0x82, 0x4e, 0x1a, 0x27, 0x1e, 0xc9,
	0x31, 0x5e, 0x78, 0x24, 0xc7, 0xf8, 0xe0, 0x91, 0x1c, 0xe3, 0x8c, 0xc7, 0x72, 0x0c, 0x51, 0x62,
	0x9e, 0x7e, 0xc1, 0xa9, 0x45, 0x65, 0xa9, 0x45, 0xfa, 0xc5, 0x45, 0xc9, 0xfa, 0x60, 0x47, 0xea,
	0xe7, 0x16, 0xa7, 0x27, 0xb1, 0x81, 0x99, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc3, 0x2d,
	0x65, 0x6a, 0xd2, 0x00, 0x00, 0x00,
}

func (m *LoadStaticMapReq) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LoadStaticMapReq) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LoadStaticMapReq) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.StaticMapID != 0 {
		i = encodeVarintWorldDatabase(dAtA, i, uint64(m.StaticMapID))
		i--
		dAtA[i] = 0x10
	}
	if m.ZoneID != 0 {
		i = encodeVarintWorldDatabase(dAtA, i, uint64(m.ZoneID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *LoadStaticMapResp) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LoadStaticMapResp) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LoadStaticMapResp) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Map != nil {
		{
			size, err := m.Map.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintWorldDatabase(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintWorldDatabase(dAtA []byte, offset int, v uint64) int {
	offset -= sovWorldDatabase(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *LoadStaticMapReq) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ZoneID != 0 {
		n += 1 + sovWorldDatabase(uint64(m.ZoneID))
	}
	if m.StaticMapID != 0 {
		n += 1 + sovWorldDatabase(uint64(m.StaticMapID))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *LoadStaticMapResp) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Map != nil {
		l = m.Map.Size()
		n += 1 + l + sovWorldDatabase(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovWorldDatabase(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozWorldDatabase(x uint64) (n int) {
	return sovWorldDatabase(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *LoadStaticMapReq) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowWorldDatabase
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
			return fmt.Errorf("proto: LoadStaticMapReq: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LoadStaticMapReq: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ZoneID", wireType)
			}
			m.ZoneID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWorldDatabase
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ZoneID |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StaticMapID", wireType)
			}
			m.StaticMapID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWorldDatabase
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StaticMapID |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipWorldDatabase(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthWorldDatabase
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthWorldDatabase
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *LoadStaticMapResp) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowWorldDatabase
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
			return fmt.Errorf("proto: LoadStaticMapResp: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LoadStaticMapResp: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Map", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowWorldDatabase
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthWorldDatabase
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthWorldDatabase
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Map == nil {
				m.Map = &data.MapData{}
			}
			if err := m.Map.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipWorldDatabase(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthWorldDatabase
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthWorldDatabase
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipWorldDatabase(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowWorldDatabase
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
					return 0, ErrIntOverflowWorldDatabase
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
					return 0, ErrIntOverflowWorldDatabase
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
				return 0, ErrInvalidLengthWorldDatabase
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupWorldDatabase
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthWorldDatabase
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthWorldDatabase        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowWorldDatabase          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupWorldDatabase = fmt.Errorf("proto: unexpected end of group")
)