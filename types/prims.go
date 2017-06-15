package types

import "github.com/chewxy/hm"

const (
	Bit               hm.TypeConst = "Bit"
	Boolean           hm.TypeConst = "Boolean"
	Integer           hm.TypeConst = "Integer"
	Integer8          hm.TypeConst = "Integer8"
	Integer16         hm.TypeConst = "Integer16"
	Integer32         hm.TypeConst = "Integer32"
	Integer64         hm.TypeConst = "Integer64"
	UnsignedInteger   hm.TypeConst = "UnsignedInteger"
	UnsignedInteger8  hm.TypeConst = "UnsignedInteger8"
	UnsignedInteger16 hm.TypeConst = "UnsignedInteger16"
	UnsignedInteger32 hm.TypeConst = "UnsignedInteger32"
	UnsignedInteger64 hm.TypeConst = "UnsignedInteger64"
	Float             hm.TypeConst = "Float"
	Float32           hm.TypeConst = "Float32"
	Float64           hm.TypeConst = "Float64"
	String            hm.TypeConst = "String"
	Top               hm.TypeConst = "⊤"
	Bottom            hm.TypeConst = "⊥"
	List              hm.TypeConst = "List"
	Image             hm.TypeConst = "Image"
)

func NewList(t hm.Type) hm.Type {
	return NewConstructor(List, t)
}
