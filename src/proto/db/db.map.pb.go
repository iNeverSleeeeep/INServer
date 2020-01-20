// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: db.map.proto

package db

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
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

type DBDynamicMap struct {
	UUID                 string   `protobuf:"bytes,1,opt,name=UUID,proto3" json:"UUID,omitempty" db:"UUID"`
	SerializedData       []byte   `protobuf:"bytes,2,opt,name=SerializedData,proto3" json:"SerializedData,omitempty" db:"SerializedData"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DBDynamicMap) Reset()         { *m = DBDynamicMap{} }
func (m *DBDynamicMap) String() string { return proto.CompactTextString(m) }
func (*DBDynamicMap) ProtoMessage()    {}
func (*DBDynamicMap) Descriptor() ([]byte, []int) {
	return fileDescriptor_08ce7afbaaa358f7, []int{0}
}
func (m *DBDynamicMap) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DBDynamicMap) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DBDynamicMap.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DBDynamicMap) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DBDynamicMap.Merge(m, src)
}
func (m *DBDynamicMap) XXX_Size() int {
	return m.Size()
}
func (m *DBDynamicMap) XXX_DiscardUnknown() {
	xxx_messageInfo_DBDynamicMap.DiscardUnknown(m)
}

var xxx_messageInfo_DBDynamicMap proto.InternalMessageInfo

func (m *DBDynamicMap) GetUUID() string {
	if m != nil {
		return m.UUID
	}
	return ""
}

func (m *DBDynamicMap) GetSerializedData() []byte {
	if m != nil {
		return m.SerializedData
	}
	return nil
}

type DBStaticMap struct {
	ZoneID               int32    `protobuf:"varint,1,opt,name=ZoneID,proto3" json:"ZoneID,omitempty" db:"ZoneID"`
	MapID                int32    `protobuf:"varint,2,opt,name=MapID,proto3" json:"MapID,omitempty" db:"MapID"`
	UUID                 string   `protobuf:"bytes,3,opt,name=UUID,proto3" json:"UUID,omitempty" db:"UUID"`
	SerializedData       []byte   `protobuf:"bytes,4,opt,name=SerializedData,proto3" json:"SerializedData,omitempty" db:"SerializedData"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DBStaticMap) Reset()         { *m = DBStaticMap{} }
func (m *DBStaticMap) String() string { return proto.CompactTextString(m) }
func (*DBStaticMap) ProtoMessage()    {}
func (*DBStaticMap) Descriptor() ([]byte, []int) {
	return fileDescriptor_08ce7afbaaa358f7, []int{1}
}
func (m *DBStaticMap) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DBStaticMap) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DBStaticMap.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DBStaticMap) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DBStaticMap.Merge(m, src)
}
func (m *DBStaticMap) XXX_Size() int {
	return m.Size()
}
func (m *DBStaticMap) XXX_DiscardUnknown() {
	xxx_messageInfo_DBStaticMap.DiscardUnknown(m)
}

var xxx_messageInfo_DBStaticMap proto.InternalMessageInfo

func (m *DBStaticMap) GetZoneID() int32 {
	if m != nil {
		return m.ZoneID
	}
	return 0
}

func (m *DBStaticMap) GetMapID() int32 {
	if m != nil {
		return m.MapID
	}
	return 0
}

func (m *DBStaticMap) GetUUID() string {
	if m != nil {
		return m.UUID
	}
	return ""
}

func (m *DBStaticMap) GetSerializedData() []byte {
	if m != nil {
		return m.SerializedData
	}
	return nil
}

func init() {
	proto.RegisterType((*DBDynamicMap)(nil), "DBDynamicMap")
	proto.RegisterType((*DBStaticMap)(nil), "DBStaticMap")
}

func init() { proto.RegisterFile("db.map.proto", fileDescriptor_08ce7afbaaa358f7) }

var fileDescriptor_08ce7afbaaa358f7 = []byte{
	// 247 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0x49, 0xd2, 0xcb,
	0x4d, 0x2c, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x97, 0xe2, 0x4a, 0xcf, 0x4f, 0xcf, 0x87, 0xb0,
	0x95, 0x8a, 0xb8, 0x78, 0x5c, 0x9c, 0x5c, 0x2a, 0xf3, 0x12, 0x73, 0x33, 0x93, 0x7d, 0x13, 0x0b,
	0x84, 0x14, 0xb9, 0x58, 0x42, 0x43, 0x3d, 0x5d, 0x24, 0x18, 0x15, 0x18, 0x35, 0x38, 0x9d, 0x78,
	0x3f, 0xdd, 0x93, 0xe7, 0x4c, 0x49, 0xb2, 0x52, 0x02, 0x89, 0x29, 0x05, 0x81, 0xa5, 0x84, 0xec,
	0xb9, 0xf8, 0x82, 0x53, 0x8b, 0x32, 0x13, 0x73, 0x32, 0xab, 0x52, 0x53, 0x5c, 0x12, 0x4b, 0x12,
	0x25, 0x98, 0x14, 0x18, 0x35, 0x78, 0x9c, 0xc4, 0x3f, 0xdd, 0x93, 0x17, 0x06, 0x29, 0x46, 0x95,
	0x55, 0x0a, 0x42, 0x53, 0xae, 0x74, 0x80, 0x91, 0x8b, 0xdb, 0xc5, 0x29, 0xb8, 0x24, 0xb1, 0x04,
	0x62, 0xa7, 0x3a, 0x17, 0x5b, 0x54, 0x7e, 0x5e, 0x2a, 0xd4, 0x56, 0x56, 0x27, 0xfe, 0x4f, 0xf7,
	0xe4, 0xb9, 0x41, 0x06, 0x41, 0x44, 0x95, 0x82, 0xa0, 0xd2, 0x42, 0x2a, 0x5c, 0xac, 0xbe, 0x89,
	0x05, 0x9e, 0x2e, 0x60, 0x0b, 0x59, 0x9d, 0xf8, 0x3e, 0xdd, 0x93, 0xe7, 0x02, 0xa9, 0x03, 0x0b,
	0x2a, 0x05, 0x41, 0x24, 0xe1, 0x5e, 0x60, 0x26, 0xc5, 0x0b, 0x2c, 0x24, 0x79, 0xc1, 0x49, 0xfd,
	0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f, 0xe4, 0x18, 0x1f, 0x3c, 0x92, 0x63, 0x9c, 0xf1, 0x58, 0x8e,
	0x21, 0x4a, 0xd4, 0xd3, 0x2f, 0x38, 0xb5, 0xa8, 0x2c, 0xb5, 0x48, 0xbf, 0xb8, 0x28, 0x59, 0x1f,
	0x1c, 0xb8, 0xfa, 0x29, 0x49, 0x49, 0x6c, 0x60, 0x96, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0xf6,
	0xe1, 0x69, 0x33, 0x82, 0x01, 0x00, 0x00,
}

func (m *DBDynamicMap) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DBDynamicMap) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DBDynamicMap) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.SerializedData) > 0 {
		i -= len(m.SerializedData)
		copy(dAtA[i:], m.SerializedData)
		i = encodeVarintDbMap(dAtA, i, uint64(len(m.SerializedData)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.UUID) > 0 {
		i -= len(m.UUID)
		copy(dAtA[i:], m.UUID)
		i = encodeVarintDbMap(dAtA, i, uint64(len(m.UUID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *DBStaticMap) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DBStaticMap) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DBStaticMap) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.SerializedData) > 0 {
		i -= len(m.SerializedData)
		copy(dAtA[i:], m.SerializedData)
		i = encodeVarintDbMap(dAtA, i, uint64(len(m.SerializedData)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.UUID) > 0 {
		i -= len(m.UUID)
		copy(dAtA[i:], m.UUID)
		i = encodeVarintDbMap(dAtA, i, uint64(len(m.UUID)))
		i--
		dAtA[i] = 0x1a
	}
	if m.MapID != 0 {
		i = encodeVarintDbMap(dAtA, i, uint64(m.MapID))
		i--
		dAtA[i] = 0x10
	}
	if m.ZoneID != 0 {
		i = encodeVarintDbMap(dAtA, i, uint64(m.ZoneID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintDbMap(dAtA []byte, offset int, v uint64) int {
	offset -= sovDbMap(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *DBDynamicMap) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.UUID)
	if l > 0 {
		n += 1 + l + sovDbMap(uint64(l))
	}
	l = len(m.SerializedData)
	if l > 0 {
		n += 1 + l + sovDbMap(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *DBStaticMap) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ZoneID != 0 {
		n += 1 + sovDbMap(uint64(m.ZoneID))
	}
	if m.MapID != 0 {
		n += 1 + sovDbMap(uint64(m.MapID))
	}
	l = len(m.UUID)
	if l > 0 {
		n += 1 + l + sovDbMap(uint64(l))
	}
	l = len(m.SerializedData)
	if l > 0 {
		n += 1 + l + sovDbMap(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovDbMap(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozDbMap(x uint64) (n int) {
	return sovDbMap(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *DBDynamicMap) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDbMap
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
			return fmt.Errorf("proto: DBDynamicMap: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DBDynamicMap: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UUID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDbMap
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
				return ErrInvalidLengthDbMap
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDbMap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UUID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SerializedData", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDbMap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthDbMap
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthDbMap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SerializedData = append(m.SerializedData[:0], dAtA[iNdEx:postIndex]...)
			if m.SerializedData == nil {
				m.SerializedData = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDbMap(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDbMap
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthDbMap
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
func (m *DBStaticMap) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDbMap
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
			return fmt.Errorf("proto: DBStaticMap: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DBStaticMap: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ZoneID", wireType)
			}
			m.ZoneID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDbMap
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
				return fmt.Errorf("proto: wrong wireType = %d for field MapID", wireType)
			}
			m.MapID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDbMap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MapID |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UUID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDbMap
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
				return ErrInvalidLengthDbMap
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDbMap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UUID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SerializedData", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDbMap
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthDbMap
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthDbMap
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SerializedData = append(m.SerializedData[:0], dAtA[iNdEx:postIndex]...)
			if m.SerializedData == nil {
				m.SerializedData = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDbMap(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthDbMap
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthDbMap
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
func skipDbMap(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDbMap
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
					return 0, ErrIntOverflowDbMap
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
					return 0, ErrIntOverflowDbMap
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
				return 0, ErrInvalidLengthDbMap
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupDbMap
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthDbMap
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthDbMap        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDbMap          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupDbMap = fmt.Errorf("proto: unexpected end of group")
)
