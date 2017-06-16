// Code generated by protoc-gen-gogo.
// source: GLMClassifier.proto
// DO NOT EDIT!

package CoreML

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Ignoring public import of StringToInt64Map from DataStructures.proto

// Ignoring public import of Int64ToStringMap from DataStructures.proto

// Ignoring public import of StringToDoubleMap from DataStructures.proto

// Ignoring public import of Int64ToDoubleMap from DataStructures.proto

// Ignoring public import of StringVector from DataStructures.proto

// Ignoring public import of Int64Vector from DataStructures.proto

// Ignoring public import of DoubleVector from DataStructures.proto

type GLMClassifier_PostEvaluationTransform int32

const (
	GLMClassifier_Logit  GLMClassifier_PostEvaluationTransform = 0
	GLMClassifier_Probit GLMClassifier_PostEvaluationTransform = 1
)

var GLMClassifier_PostEvaluationTransform_name = map[int32]string{
	0: "Logit",
	1: "Probit",
}
var GLMClassifier_PostEvaluationTransform_value = map[string]int32{
	"Logit":  0,
	"Probit": 1,
}

func (x GLMClassifier_PostEvaluationTransform) String() string {
	return proto.EnumName(GLMClassifier_PostEvaluationTransform_name, int32(x))
}
func (GLMClassifier_PostEvaluationTransform) EnumDescriptor() ([]byte, []int) {
	return fileDescriptorGLMClassifier, []int{0, 0}
}

type GLMClassifier_ClassEncoding int32

const (
	GLMClassifier_ReferenceClass GLMClassifier_ClassEncoding = 0
	GLMClassifier_OneVsRest      GLMClassifier_ClassEncoding = 1
)

var GLMClassifier_ClassEncoding_name = map[int32]string{
	0: "ReferenceClass",
	1: "OneVsRest",
}
var GLMClassifier_ClassEncoding_value = map[string]int32{
	"ReferenceClass": 0,
	"OneVsRest":      1,
}

func (x GLMClassifier_ClassEncoding) String() string {
	return proto.EnumName(GLMClassifier_ClassEncoding_name, int32(x))
}
func (GLMClassifier_ClassEncoding) EnumDescriptor() ([]byte, []int) {
	return fileDescriptorGLMClassifier, []int{0, 1}
}

// *
// A generalized linear model classifier.
type GLMClassifier struct {
	Weights                 []*GLMClassifier_DoubleArray          `protobuf:"bytes,1,rep,name=weights" json:"weights,omitempty"`
	Offset                  []float64                             `protobuf:"fixed64,2,rep,packed,name=offset" json:"offset,omitempty"`
	PostEvaluationTransform GLMClassifier_PostEvaluationTransform `protobuf:"varint,3,opt,name=postEvaluationTransform,proto3,enum=CoreML.GLMClassifier_PostEvaluationTransform" json:"postEvaluationTransform,omitempty"`
	ClassEncoding           GLMClassifier_ClassEncoding           `protobuf:"varint,4,opt,name=classEncoding,proto3,enum=CoreML.GLMClassifier_ClassEncoding" json:"classEncoding,omitempty"`
	// *
	// Required class label mapping.
	//
	// Types that are valid to be assigned to ClassLabels:
	//	*GLMClassifier_StringClassLabels
	//	*GLMClassifier_Int64ClassLabels
	ClassLabels isGLMClassifier_ClassLabels `protobuf_oneof:"ClassLabels"`
}

func (m *GLMClassifier) Reset()                    { *m = GLMClassifier{} }
func (m *GLMClassifier) String() string            { return proto.CompactTextString(m) }
func (*GLMClassifier) ProtoMessage()               {}
func (*GLMClassifier) Descriptor() ([]byte, []int) { return fileDescriptorGLMClassifier, []int{0} }

type isGLMClassifier_ClassLabels interface {
	isGLMClassifier_ClassLabels()
	MarshalTo([]byte) (int, error)
	Size() int
}

type GLMClassifier_StringClassLabels struct {
	StringClassLabels *StringVector `protobuf:"bytes,100,opt,name=stringClassLabels,oneof"`
}
type GLMClassifier_Int64ClassLabels struct {
	Int64ClassLabels *Int64Vector `protobuf:"bytes,101,opt,name=int64ClassLabels,oneof"`
}

func (*GLMClassifier_StringClassLabels) isGLMClassifier_ClassLabels() {}
func (*GLMClassifier_Int64ClassLabels) isGLMClassifier_ClassLabels()  {}

func (m *GLMClassifier) GetClassLabels() isGLMClassifier_ClassLabels {
	if m != nil {
		return m.ClassLabels
	}
	return nil
}

func (m *GLMClassifier) GetWeights() []*GLMClassifier_DoubleArray {
	if m != nil {
		return m.Weights
	}
	return nil
}

func (m *GLMClassifier) GetOffset() []float64 {
	if m != nil {
		return m.Offset
	}
	return nil
}

func (m *GLMClassifier) GetPostEvaluationTransform() GLMClassifier_PostEvaluationTransform {
	if m != nil {
		return m.PostEvaluationTransform
	}
	return GLMClassifier_Logit
}

func (m *GLMClassifier) GetClassEncoding() GLMClassifier_ClassEncoding {
	if m != nil {
		return m.ClassEncoding
	}
	return GLMClassifier_ReferenceClass
}

func (m *GLMClassifier) GetStringClassLabels() *StringVector {
	if x, ok := m.GetClassLabels().(*GLMClassifier_StringClassLabels); ok {
		return x.StringClassLabels
	}
	return nil
}

func (m *GLMClassifier) GetInt64ClassLabels() *Int64Vector {
	if x, ok := m.GetClassLabels().(*GLMClassifier_Int64ClassLabels); ok {
		return x.Int64ClassLabels
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*GLMClassifier) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _GLMClassifier_OneofMarshaler, _GLMClassifier_OneofUnmarshaler, _GLMClassifier_OneofSizer, []interface{}{
		(*GLMClassifier_StringClassLabels)(nil),
		(*GLMClassifier_Int64ClassLabels)(nil),
	}
}

func _GLMClassifier_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*GLMClassifier)
	// ClassLabels
	switch x := m.ClassLabels.(type) {
	case *GLMClassifier_StringClassLabels:
		_ = b.EncodeVarint(100<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.StringClassLabels); err != nil {
			return err
		}
	case *GLMClassifier_Int64ClassLabels:
		_ = b.EncodeVarint(101<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Int64ClassLabels); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("GLMClassifier.ClassLabels has unexpected type %T", x)
	}
	return nil
}

func _GLMClassifier_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*GLMClassifier)
	switch tag {
	case 100: // ClassLabels.stringClassLabels
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(StringVector)
		err := b.DecodeMessage(msg)
		m.ClassLabels = &GLMClassifier_StringClassLabels{msg}
		return true, err
	case 101: // ClassLabels.int64ClassLabels
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Int64Vector)
		err := b.DecodeMessage(msg)
		m.ClassLabels = &GLMClassifier_Int64ClassLabels{msg}
		return true, err
	default:
		return false, nil
	}
}

func _GLMClassifier_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*GLMClassifier)
	// ClassLabels
	switch x := m.ClassLabels.(type) {
	case *GLMClassifier_StringClassLabels:
		s := proto.Size(x.StringClassLabels)
		n += proto.SizeVarint(100<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *GLMClassifier_Int64ClassLabels:
		s := proto.Size(x.Int64ClassLabels)
		n += proto.SizeVarint(101<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type GLMClassifier_DoubleArray struct {
	Value []float64 `protobuf:"fixed64,1,rep,packed,name=value" json:"value,omitempty"`
}

func (m *GLMClassifier_DoubleArray) Reset()         { *m = GLMClassifier_DoubleArray{} }
func (m *GLMClassifier_DoubleArray) String() string { return proto.CompactTextString(m) }
func (*GLMClassifier_DoubleArray) ProtoMessage()    {}
func (*GLMClassifier_DoubleArray) Descriptor() ([]byte, []int) {
	return fileDescriptorGLMClassifier, []int{0, 0}
}

func (m *GLMClassifier_DoubleArray) GetValue() []float64 {
	if m != nil {
		return m.Value
	}
	return nil
}

func init() {
	proto.RegisterType((*GLMClassifier)(nil), "CoreML.GLMClassifier")
	proto.RegisterType((*GLMClassifier_DoubleArray)(nil), "CoreML.GLMClassifier.DoubleArray")
	proto.RegisterEnum("CoreML.GLMClassifier_PostEvaluationTransform", GLMClassifier_PostEvaluationTransform_name, GLMClassifier_PostEvaluationTransform_value)
	proto.RegisterEnum("CoreML.GLMClassifier_ClassEncoding", GLMClassifier_ClassEncoding_name, GLMClassifier_ClassEncoding_value)
}
func (m *GLMClassifier) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GLMClassifier) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Weights) > 0 {
		for _, msg := range m.Weights {
			dAtA[i] = 0xa
			i++
			i = encodeVarintGLMClassifier(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	if len(m.Offset) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintGLMClassifier(dAtA, i, uint64(len(m.Offset)*8))
		for _, num := range m.Offset {
			f1 := math.Float64bits(float64(num))
			dAtA[i] = uint8(f1)
			i++
			dAtA[i] = uint8(f1 >> 8)
			i++
			dAtA[i] = uint8(f1 >> 16)
			i++
			dAtA[i] = uint8(f1 >> 24)
			i++
			dAtA[i] = uint8(f1 >> 32)
			i++
			dAtA[i] = uint8(f1 >> 40)
			i++
			dAtA[i] = uint8(f1 >> 48)
			i++
			dAtA[i] = uint8(f1 >> 56)
			i++
		}
	}
	if m.PostEvaluationTransform != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintGLMClassifier(dAtA, i, uint64(m.PostEvaluationTransform))
	}
	if m.ClassEncoding != 0 {
		dAtA[i] = 0x20
		i++
		i = encodeVarintGLMClassifier(dAtA, i, uint64(m.ClassEncoding))
	}
	if m.ClassLabels != nil {
		nn2, err := m.ClassLabels.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += nn2
	}
	return i, nil
}

func (m *GLMClassifier_StringClassLabels) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.StringClassLabels != nil {
		dAtA[i] = 0xa2
		i++
		dAtA[i] = 0x6
		i++
		i = encodeVarintGLMClassifier(dAtA, i, uint64(m.StringClassLabels.Size()))
		n3, err := m.StringClassLabels.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n3
	}
	return i, nil
}
func (m *GLMClassifier_Int64ClassLabels) MarshalTo(dAtA []byte) (int, error) {
	i := 0
	if m.Int64ClassLabels != nil {
		dAtA[i] = 0xaa
		i++
		dAtA[i] = 0x6
		i++
		i = encodeVarintGLMClassifier(dAtA, i, uint64(m.Int64ClassLabels.Size()))
		n4, err := m.Int64ClassLabels.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n4
	}
	return i, nil
}
func (m *GLMClassifier_DoubleArray) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GLMClassifier_DoubleArray) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Value) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintGLMClassifier(dAtA, i, uint64(len(m.Value)*8))
		for _, num := range m.Value {
			f5 := math.Float64bits(float64(num))
			dAtA[i] = uint8(f5)
			i++
			dAtA[i] = uint8(f5 >> 8)
			i++
			dAtA[i] = uint8(f5 >> 16)
			i++
			dAtA[i] = uint8(f5 >> 24)
			i++
			dAtA[i] = uint8(f5 >> 32)
			i++
			dAtA[i] = uint8(f5 >> 40)
			i++
			dAtA[i] = uint8(f5 >> 48)
			i++
			dAtA[i] = uint8(f5 >> 56)
			i++
		}
	}
	return i, nil
}

func encodeFixed64GLMClassifier(dAtA []byte, offset int, v uint64) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	dAtA[offset+4] = uint8(v >> 32)
	dAtA[offset+5] = uint8(v >> 40)
	dAtA[offset+6] = uint8(v >> 48)
	dAtA[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32GLMClassifier(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintGLMClassifier(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *GLMClassifier) Size() (n int) {
	var l int
	_ = l
	if len(m.Weights) > 0 {
		for _, e := range m.Weights {
			l = e.Size()
			n += 1 + l + sovGLMClassifier(uint64(l))
		}
	}
	if len(m.Offset) > 0 {
		n += 1 + sovGLMClassifier(uint64(len(m.Offset)*8)) + len(m.Offset)*8
	}
	if m.PostEvaluationTransform != 0 {
		n += 1 + sovGLMClassifier(uint64(m.PostEvaluationTransform))
	}
	if m.ClassEncoding != 0 {
		n += 1 + sovGLMClassifier(uint64(m.ClassEncoding))
	}
	if m.ClassLabels != nil {
		n += m.ClassLabels.Size()
	}
	return n
}

func (m *GLMClassifier_StringClassLabels) Size() (n int) {
	var l int
	_ = l
	if m.StringClassLabels != nil {
		l = m.StringClassLabels.Size()
		n += 2 + l + sovGLMClassifier(uint64(l))
	}
	return n
}
func (m *GLMClassifier_Int64ClassLabels) Size() (n int) {
	var l int
	_ = l
	if m.Int64ClassLabels != nil {
		l = m.Int64ClassLabels.Size()
		n += 2 + l + sovGLMClassifier(uint64(l))
	}
	return n
}
func (m *GLMClassifier_DoubleArray) Size() (n int) {
	var l int
	_ = l
	if len(m.Value) > 0 {
		n += 1 + sovGLMClassifier(uint64(len(m.Value)*8)) + len(m.Value)*8
	}
	return n
}

func sovGLMClassifier(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozGLMClassifier(x uint64) (n int) {
	return sovGLMClassifier(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GLMClassifier) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGLMClassifier
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: GLMClassifier: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GLMClassifier: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Weights", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGLMClassifier
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGLMClassifier
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Weights = append(m.Weights, &GLMClassifier_DoubleArray{})
			if err := m.Weights[len(m.Weights)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType == 1 {
				var v uint64
				if (iNdEx + 8) > l {
					return io.ErrUnexpectedEOF
				}
				iNdEx += 8
				v = uint64(dAtA[iNdEx-8])
				v |= uint64(dAtA[iNdEx-7]) << 8
				v |= uint64(dAtA[iNdEx-6]) << 16
				v |= uint64(dAtA[iNdEx-5]) << 24
				v |= uint64(dAtA[iNdEx-4]) << 32
				v |= uint64(dAtA[iNdEx-3]) << 40
				v |= uint64(dAtA[iNdEx-2]) << 48
				v |= uint64(dAtA[iNdEx-1]) << 56
				v2 := float64(math.Float64frombits(v))
				m.Offset = append(m.Offset, v2)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGLMClassifier
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthGLMClassifier
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v uint64
					if (iNdEx + 8) > l {
						return io.ErrUnexpectedEOF
					}
					iNdEx += 8
					v = uint64(dAtA[iNdEx-8])
					v |= uint64(dAtA[iNdEx-7]) << 8
					v |= uint64(dAtA[iNdEx-6]) << 16
					v |= uint64(dAtA[iNdEx-5]) << 24
					v |= uint64(dAtA[iNdEx-4]) << 32
					v |= uint64(dAtA[iNdEx-3]) << 40
					v |= uint64(dAtA[iNdEx-2]) << 48
					v |= uint64(dAtA[iNdEx-1]) << 56
					v2 := float64(math.Float64frombits(v))
					m.Offset = append(m.Offset, v2)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Offset", wireType)
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PostEvaluationTransform", wireType)
			}
			m.PostEvaluationTransform = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGLMClassifier
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PostEvaluationTransform |= (GLMClassifier_PostEvaluationTransform(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClassEncoding", wireType)
			}
			m.ClassEncoding = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGLMClassifier
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ClassEncoding |= (GLMClassifier_ClassEncoding(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 100:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StringClassLabels", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGLMClassifier
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGLMClassifier
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &StringVector{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.ClassLabels = &GLMClassifier_StringClassLabels{v}
			iNdEx = postIndex
		case 101:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Int64ClassLabels", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGLMClassifier
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGLMClassifier
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &Int64Vector{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.ClassLabels = &GLMClassifier_Int64ClassLabels{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGLMClassifier(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthGLMClassifier
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
func (m *GLMClassifier_DoubleArray) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGLMClassifier
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: DoubleArray: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DoubleArray: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType == 1 {
				var v uint64
				if (iNdEx + 8) > l {
					return io.ErrUnexpectedEOF
				}
				iNdEx += 8
				v = uint64(dAtA[iNdEx-8])
				v |= uint64(dAtA[iNdEx-7]) << 8
				v |= uint64(dAtA[iNdEx-6]) << 16
				v |= uint64(dAtA[iNdEx-5]) << 24
				v |= uint64(dAtA[iNdEx-4]) << 32
				v |= uint64(dAtA[iNdEx-3]) << 40
				v |= uint64(dAtA[iNdEx-2]) << 48
				v |= uint64(dAtA[iNdEx-1]) << 56
				v2 := float64(math.Float64frombits(v))
				m.Value = append(m.Value, v2)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowGLMClassifier
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthGLMClassifier
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				for iNdEx < postIndex {
					var v uint64
					if (iNdEx + 8) > l {
						return io.ErrUnexpectedEOF
					}
					iNdEx += 8
					v = uint64(dAtA[iNdEx-8])
					v |= uint64(dAtA[iNdEx-7]) << 8
					v |= uint64(dAtA[iNdEx-6]) << 16
					v |= uint64(dAtA[iNdEx-5]) << 24
					v |= uint64(dAtA[iNdEx-4]) << 32
					v |= uint64(dAtA[iNdEx-3]) << 40
					v |= uint64(dAtA[iNdEx-2]) << 48
					v |= uint64(dAtA[iNdEx-1]) << 56
					v2 := float64(math.Float64frombits(v))
					m.Value = append(m.Value, v2)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipGLMClassifier(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthGLMClassifier
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
func skipGLMClassifier(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGLMClassifier
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
					return 0, ErrIntOverflowGLMClassifier
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGLMClassifier
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
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthGLMClassifier
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowGLMClassifier
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipGLMClassifier(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthGLMClassifier = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGLMClassifier   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("GLMClassifier.proto", fileDescriptorGLMClassifier) }

var fileDescriptorGLMClassifier = []byte{
	// 381 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0x41, 0x6f, 0xd3, 0x30,
	0x14, 0xc7, 0x63, 0xb2, 0x05, 0xed, 0x45, 0x99, 0x82, 0x57, 0xb1, 0x68, 0x87, 0x28, 0x74, 0x97,
	0x5c, 0xa8, 0x50, 0x40, 0x5c, 0x38, 0x6d, 0xed, 0xc4, 0x26, 0x65, 0xa2, 0xf2, 0xd0, 0xee, 0x4e,
	0xfa, 0x12, 0x2c, 0x05, 0xbb, 0xb2, 0x1d, 0x10, 0xdf, 0x82, 0xef, 0xc3, 0x17, 0xe0, 0xc8, 0x47,
	0x40, 0xe5, 0x8b, 0xa0, 0xa6, 0x14, 0x25, 0x6a, 0x7b, 0x7c, 0x7a, 0xff, 0xdf, 0xcf, 0x7e, 0xcf,
	0x86, 0xb3, 0xf7, 0xf9, 0xfd, 0xb4, 0xe1, 0xc6, 0x88, 0x4a, 0xa0, 0x9e, 0x2c, 0xb5, 0xb2, 0x8a,
	0x7a, 0x53, 0xa5, 0xf1, 0x3e, 0xbf, 0x18, 0xcd, 0xb8, 0xe5, 0x0f, 0x56, 0xb7, 0xa5, 0x6d, 0x35,
	0x9a, 0x4d, 0x77, 0xfc, 0xe3, 0x08, 0x82, 0x01, 0x45, 0xdf, 0xc1, 0xd3, 0xaf, 0x28, 0xea, 0x4f,
	0xd6, 0x44, 0x24, 0x71, 0x53, 0x3f, 0x7b, 0x31, 0xd9, 0x18, 0x26, 0x43, 0xfb, 0x4c, 0xb5, 0x45,
	0x83, 0x57, 0x5a, 0xf3, 0x6f, 0x6c, 0x4b, 0xd0, 0xe7, 0xe0, 0xa9, 0xaa, 0x32, 0x68, 0xa3, 0x27,
	0x89, 0x9b, 0x12, 0xf6, 0xaf, 0xa2, 0x35, 0x9c, 0x2f, 0x95, 0xb1, 0x37, 0x5f, 0x78, 0xd3, 0x72,
	0x2b, 0x94, 0xfc, 0xa8, 0xb9, 0x34, 0x95, 0xd2, 0x9f, 0x23, 0x37, 0x21, 0xe9, 0x69, 0xf6, 0x72,
	0xff, 0x21, 0xf3, 0xfd, 0x10, 0x3b, 0x64, 0xa3, 0x77, 0x10, 0x94, 0x6b, 0xfc, 0x46, 0x96, 0x6a,
	0x21, 0x64, 0x1d, 0x1d, 0x75, 0xfa, 0xcb, 0xfd, 0xfa, 0x69, 0x3f, 0xca, 0x86, 0x24, 0x9d, 0xc1,
	0x33, 0x63, 0xb5, 0x90, 0x75, 0x97, 0xca, 0x79, 0x81, 0x8d, 0x89, 0x16, 0x09, 0x49, 0xfd, 0x6c,
	0xb4, 0xd5, 0x3d, 0x74, 0x81, 0x47, 0x2c, 0xad, 0xd2, 0xb7, 0x0e, 0xdb, 0x05, 0xe8, 0x15, 0x84,
	0x42, 0xda, 0xb7, 0x6f, 0xfa, 0x12, 0xec, 0x24, 0x67, 0x5b, 0xc9, 0xdd, 0xba, 0xff, 0xdf, 0xb1,
	0x13, 0xbf, 0xb8, 0x04, 0xbf, 0xb7, 0x6c, 0x3a, 0x82, 0xe3, 0xf5, 0xe0, 0xd8, 0x3d, 0x0f, 0x61,
	0x9b, 0x62, 0xfc, 0x0a, 0xce, 0x0f, 0x2c, 0x8b, 0x9e, 0xc0, 0x71, 0xae, 0x6a, 0x61, 0x43, 0x87,
	0x02, 0x78, 0x73, 0xad, 0x0a, 0x61, 0x43, 0x32, 0xce, 0x20, 0x18, 0xcc, 0x4f, 0x29, 0x9c, 0x32,
	0xac, 0x50, 0xa3, 0x2c, 0xb1, 0xeb, 0x84, 0x0e, 0x0d, 0xe0, 0xe4, 0x83, 0xc4, 0x47, 0xc3, 0xd0,
	0xd8, 0x90, 0x5c, 0x07, 0xe0, 0xf7, 0x6e, 0x76, 0x4d, 0x7f, 0xae, 0x62, 0xf2, 0x6b, 0x15, 0x93,
	0xdf, 0xab, 0x98, 0x7c, 0xff, 0x13, 0x3b, 0xb7, 0xee, 0xdc, 0x29, 0xbc, 0xee, 0x6b, 0xbd, 0xfe,
	0x1b, 0x00, 0x00, 0xff, 0xff, 0x72, 0x80, 0x01, 0xef, 0x8f, 0x02, 0x00, 0x00,
}
