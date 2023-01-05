package goincv

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/spf13/cast"
)

func InterfaceConvertValue(d interface{}) *Value {
	return &Value{
		data: d,
	}
}

func InitValue(initValue interface{}, shape Shape) *Value {
	val := InterfaceConvertValue(initValue)
	vShape := val.GetShapeBySlice()
	vs := Values{}

	dims := 1
	for i := range shape {
		dims = dims * shape[i]
	}
	if len(vShape) == 1 {
		for i := 0; i < dims; i++ {
			vs = append(vs, val)
		}
	} else {
		nvas := val.Flatten()
		for i := 0; i < dims; i++ {
			vs = append(vs, nvas...)
		}
		vs = vs[:dims]
	}

	val = vs.EncodeValue()
	val.Reshape(shape)
	return val
}

type Value struct {
	data interface{}
}

type Values []*Value

func (v Values) EncodeValue() *Value {
	ret := []interface{}{}
	for i := range v {
		ret = append(ret, v[i].data)
	}
	return &Value{
		data: ret,
	}
}

func (v Value) Raw() interface{} {
	return v.data
}

func (v *Value) GetShapeBySlice() Shape {
	var dShape Shape
	xv := reflect.ValueOf(v.data)
	if xv.Type().Kind() != reflect.Slice {
		return Shape{}
	}

	dShape = append(dShape, xv.Len())
	for {
		xv = xv.Index(0)
		if xv.Type().Kind() == reflect.Interface {
			xv = reflect.ValueOf(InterfaceArrayFix(xv.Interface()))
		}
		if xv.Type().Kind() != reflect.Slice {
			break
		}
		dShape = append(dShape, xv.Len())
	}
	return dShape

}

func (v *Value) Array() Values {
	dv := reflect.ValueOf(v.data)
	var ret Values
	for i := 0; i < dv.Len(); i++ {
		ret = append(ret, &Value{
			data: dv.Index(i).Interface(),
		})
	}
	return ret
}

func (v *Value) Flatten() Values {
	shape := v.GetShapeBySlice()
	if len(shape) == 0 {
		return Values{v}
	}
	var ret Values
	dv := reflect.ValueOf(v.data)
	for i := 0; i < dv.Len(); i++ {
		xv := dv.Index(i)
		vt := &Value{
			data: xv.Interface(),
		}
		if len(shape) > 1 {
			ret = append(ret, vt.Flatten()...)
		} else {
			ret = append(ret, vt)
		}
	}
	return ret
}

func (v *Value) Reshape(reshape Shape) (*Value, error) {
	var err error
	data := v.Flatten()
	size := len(data)
	if len(reshape) == 1 {
		if size == reshape[0] {
			return data.EncodeValue(), nil
		} else {
			err = errors.New(fmt.Sprintf("a len == %d , b len == %d , a!=b", size, reshape[0]))
			return nil, err
		}
	}

	if size%reshape[0] != 0 {
		err = errors.New(fmt.Sprint("size%reshape[0]  != 0    ", size, reshape[0]))
		return nil, err
	}

	bachtSize := size / reshape[0]
	var retFloat Values
	for k := 0; k < reshape[0]; k++ {
		itemValue := data[k*bachtSize : (k+1)*bachtSize]

		retItem, err := itemValue.EncodeValue().Reshape(reshape[1:])
		if err != nil {
			return nil, err
		}
		retFloat = append(retFloat, retItem)
	}
	return retFloat.EncodeValue(), nil

}

//转换代码

func (v *Value) Float32() float32 {
	if v, ok := v.data.(float32); ok {
		return v
	}
	return cast.ToFloat32(v.data)
}

func (v *Value) To1DFloat32() (ret []float32) {
	if v, ok := v.data.([]float32); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].Float32())
	}
	return ret
}

func (v *Value) To2DFloat32() (ret [][]float32) {
	if v, ok := v.data.([][]float32); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To1DFloat32())
	}
	return ret
}
func (v *Value) To3DFloat32() (ret [][][]float32) {
	if v, ok := v.data.([][][]float32); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To2DFloat32())
	}
	return ret
}
func (v *Value) To4DFloat32() (ret [][][][]float32) {
	if v, ok := v.data.([][][][]float32); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To3DFloat32())
	}
	return ret
}
func (v *Value) To5DFloat32() (ret [][][][][]float32) {
	if v, ok := v.data.([][][][][]float32); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To4DFloat32())
	}
	return ret
}

func (v *Value) Float64() float64 {
	if v, ok := v.data.(float64); ok {
		return v
	}
	return cast.ToFloat64(v.data)
}
func (v *Value) To1DFloat64() (ret []float64) {
	if v, ok := v.data.([]float64); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].Float64())
	}
	return ret
}

func (v *Value) To2DFloat64() (ret [][]float64) {
	if v, ok := v.data.([][]float64); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To1DFloat64())
	}
	return ret
}
func (v *Value) To3DFloat64() (ret [][][]float64) {
	if v, ok := v.data.([][][]float64); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To2DFloat64())
	}
	return ret
}
func (v *Value) To4DFloat64() (ret [][][][]float64) {
	if v, ok := v.data.([][][][]float64); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To3DFloat64())
	}
	return ret
}
func (v *Value) To5DFloat64() (ret [][][][][]float64) {
	if v, ok := v.data.([][][][][]float64); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To4DFloat64())
	}
	return ret
}

func (v *Value) Int64() int64 {
	if v, ok := v.data.(int64); ok {
		return v
	}
	return cast.ToInt64(v.data)
}

func (v *Value) To1DInt64() (ret []int64) {
	if v, ok := v.data.([]int64); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].Int64())
	}
	return ret
}

func (v *Value) To2DInt64() (ret [][]int64) {
	if v, ok := v.data.([][]int64); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To1DInt64())
	}
	return ret
}
func (v *Value) To3DInt64() (ret [][][]int64) {
	if v, ok := v.data.([][][]int64); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To2DInt64())
	}
	return ret
}
func (v *Value) To4DInt64() (ret [][][][]int64) {
	if v, ok := v.data.([][][][]int64); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To3DInt64())
	}
	return ret
}
func (v *Value) To5DInt64() (ret [][][][][]int64) {
	if v, ok := v.data.([][][][][]int64); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To4DInt64())
	}
	return ret
}

func (v *Value) Int32() int32 {
	if v, ok := v.data.(int32); ok {
		return v
	}
	return cast.ToInt32(v.data)
}

func (v *Value) To1DInt32() (ret []int32) {
	if v, ok := v.data.([]int32); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].Int32())
	}
	return ret
}

func (v *Value) To2DInt32() (ret [][]int32) {
	if v, ok := v.data.([][]int32); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To1DInt32())
	}
	return ret
}
func (v *Value) To3DInt32() (ret [][][]int32) {
	if v, ok := v.data.([][][]int32); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To2DInt32())
	}
	return ret
}
func (v *Value) To4DInt32() (ret [][][][]int32) {
	if v, ok := v.data.([][][][]int32); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To3DInt32())
	}
	return ret
}
func (v *Value) To5DInt32() (ret [][][][][]int32) {
	if v, ok := v.data.([][][][][]int32); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To4DInt32())
	}
	return ret
}

func (v *Value) Int16() int16 {
	if v, ok := v.data.(int16); ok {
		return v
	}
	return cast.ToInt16(v.data)
}

func (v *Value) To1DInt16() (ret []int16) {
	if v, ok := v.data.([]int16); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].Int16())
	}
	return ret
}

func (v *Value) To2DInt16() (ret [][]int16) {
	if v, ok := v.data.([][]int16); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To1DInt16())
	}
	return ret
}
func (v *Value) To3DInt16() (ret [][][]int16) {
	if v, ok := v.data.([][][]int16); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To2DInt16())
	}
	return ret
}
func (v *Value) To4DInt16() (ret [][][][]int16) {
	if v, ok := v.data.([][][][]int16); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To3DInt16())
	}
	return ret
}
func (v *Value) To5DInt16() (ret [][][][][]int16) {
	if v, ok := v.data.([][][][][]int16); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To4DInt16())
	}
	return ret
}
func (v *Value) Int8() int8 {
	if v, ok := v.data.(int8); ok {
		return v
	}
	return cast.ToInt8(v.data)
}

func (v *Value) To1DInt8() (ret []int8) {
	if v, ok := v.data.([]int8); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].Int8())
	}
	return ret
}

func (v *Value) To2DInt8() (ret [][]int8) {
	if v, ok := v.data.([][]int8); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To1DInt8())
	}
	return ret
}
func (v *Value) To3DInt8() (ret [][][]int8) {
	if v, ok := v.data.([][][]int8); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To2DInt8())
	}
	return ret
}
func (v *Value) To4DInt8() (ret [][][][]int8) {
	if v, ok := v.data.([][][][]int8); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To3DInt8())
	}
	return ret
}
func (v *Value) To5DInt8() (ret [][][][][]int8) {
	if v, ok := v.data.([][][][][]int8); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To4DInt8())
	}
	return ret
}
func (v *Value) Int() int {
	if v, ok := v.data.(int); ok {
		return v
	}
	return cast.ToInt(v.data)
}

func (v *Value) To1DInt() (ret []int) {
	if v, ok := v.data.([]int); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].Int())
	}
	return ret
}

func (v *Value) To2DInt() (ret [][]int) {
	if v, ok := v.data.([][]int); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To1DInt())
	}
	return ret
}
func (v *Value) To3DInt() (ret [][][]int) {
	if v, ok := v.data.([][][]int); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To2DInt())
	}
	return ret
}
func (v *Value) To4DInt() (ret [][][][]int) {
	if v, ok := v.data.([][][][]int); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To3DInt())
	}
	return ret
}
func (v *Value) To5DInt() (ret [][][][][]int) {
	if v, ok := v.data.([][][][][]int); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To4DInt())
	}
	return ret
}

func (v *Value) Uint() uint {
	if v, ok := v.data.(uint); ok {
		return v
	}
	return cast.ToUint(v.data)
}

func (v *Value) To1DUint() (ret []uint) {
	if v, ok := v.data.([]uint); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].Uint())
	}
	return ret
}

func (v *Value) To2DUint() (ret [][]uint) {
	if v, ok := v.data.([][]uint); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To1DUint())
	}
	return ret
}
func (v *Value) To3DUint() (ret [][][]uint) {
	if v, ok := v.data.([][][]uint); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To2DUint())
	}
	return ret
}
func (v *Value) To4DUint() (ret [][][][]uint) {
	if v, ok := v.data.([][][][]uint); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To3DUint())
	}
	return ret
}
func (v *Value) To5DUint() (ret [][][][][]uint) {
	if v, ok := v.data.([][][][][]uint); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To4DUint())
	}
	return ret
}
func (v *Value) Uint64() uint64 {
	if v, ok := v.data.(uint64); ok {
		return v
	}
	return cast.ToUint64(v.data)
}

func (v *Value) To1DUint64() (ret []uint64) {
	if v, ok := v.data.([]uint64); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].Uint64())
	}
	return ret
}

func (v *Value) To2DUint64() (ret [][]uint64) {
	if v, ok := v.data.([][]uint64); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To1DUint64())
	}
	return ret
}
func (v *Value) To3DUint64() (ret [][][]uint64) {
	if v, ok := v.data.([][][]uint64); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To2DUint64())
	}
	return ret
}
func (v *Value) To4DUint64() (ret [][][][]uint64) {
	if v, ok := v.data.([][][][]uint64); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To3DUint64())
	}
	return ret
}
func (v *Value) To5DUint64() (ret [][][][][]uint64) {
	if v, ok := v.data.([][][][][]uint64); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To4DUint64())
	}
	return ret
}

func (v *Value) Uint32() uint32 {
	if v, ok := v.data.(uint32); ok {
		return v
	}
	return cast.ToUint32(v.data)
}

func (v *Value) To1DUint32() (ret []uint32) {
	if v, ok := v.data.([]uint32); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].Uint32())
	}
	return ret
}

func (v *Value) To2DUint32() (ret [][]uint32) {
	if v, ok := v.data.([][]uint32); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To1DUint32())
	}
	return ret
}
func (v *Value) To3DUint32() (ret [][][]uint32) {
	if v, ok := v.data.([][][]uint32); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To2DUint32())
	}
	return ret
}
func (v *Value) To4DUint32() (ret [][][][]uint32) {
	if v, ok := v.data.([][][][]uint32); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To3DUint32())
	}
	return ret
}
func (v *Value) To5DUint32() (ret [][][][][]uint32) {
	if v, ok := v.data.([][][][][]uint32); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To4DUint32())
	}
	return ret
}
func (v *Value) Uint16() uint16 {
	if v, ok := v.data.(uint16); ok {
		return v
	}
	return cast.ToUint16(v.data)
}

func (v *Value) To1DUint16() (ret []uint16) {
	if v, ok := v.data.([]uint16); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].Uint16())
	}
	return ret
}

func (v *Value) To2DUint16() (ret [][]uint16) {
	if v, ok := v.data.([][]uint16); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To1DUint16())
	}
	return ret
}
func (v *Value) To3DUint16() (ret [][][]uint16) {
	if v, ok := v.data.([][][]uint16); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To2DUint16())
	}
	return ret
}
func (v *Value) To4DUint16() (ret [][][][]uint16) {
	if v, ok := v.data.([][][][]uint16); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To3DUint16())
	}
	return ret
}
func (v *Value) To5DUint16() (ret [][][][][]uint16) {
	if v, ok := v.data.([][][][][]uint16); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To4DUint16())
	}
	return ret
}
func (v *Value) Uint8() uint8 {
	if v, ok := v.data.(uint8); ok {
		return v
	}
	return cast.ToUint8(v.data)
}

func (v *Value) To1DUint8() (ret []uint8) {
	if v, ok := v.data.([]uint8); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].Uint8())
	}
	return ret
}

func (v *Value) To2DUint8() (ret [][]uint8) {
	if v, ok := v.data.([][]uint8); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To1DUint8())
	}
	return ret
}
func (v *Value) To3DUint8() (ret [][][]uint8) {
	if v, ok := v.data.([][][]uint8); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To2DUint8())
	}
	return ret
}
func (v *Value) To4DUint8() (ret [][][][]uint8) {
	if v, ok := v.data.([][][][]uint8); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To3DUint8())
	}
	return ret
}
func (v *Value) To5DUint8() (ret [][][][][]uint8) {
	if v, ok := v.data.([][][][][]uint8); ok {
		return v
	}
	vs := v.Array()
	for i := range vs {
		ret = append(ret, vs[i].To4DUint8())
	}
	return ret
}
