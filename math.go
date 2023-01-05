package goincv

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"

	"github.com/spf13/cast"
)

type Shape []int

func ShapeForm32(s []int32) Shape {
	ret := Shape{}
	for i := range s {
		ret = append(ret, int(s[i]))
	}
	return ret
}

func (s Shape) To64() []int64 {
	ret := []int64{}
	for i := range s {
		ret = append(ret, int64(s[i]))
	}
	return ret
}

func (s Shape) To32() []int32 {
	ret := []int32{}
	for i := range s {
		ret = append(ret, int32(s[i]))
	}
	return ret
}

func InterfaceArrayFix(d interface{}) interface{} {
	_dt, _ := json.Marshal(d)
	json.Unmarshal(_dt, &d)
	return d
}

func GetShapeBySlice(d interface{}) Shape {
	var dShape Shape
	xv := reflect.ValueOf(d)
	if xv.Type().Kind() != reflect.Slice {
		return Shape{}
	}
	dShape = append(dShape, xv.Len())
	for {
		xv = xv.Index(0)

		if xv.Type().Kind() == reflect.Interface {
			xv = reflect.ValueOf(InterfaceArrayFix(xv.Interface()))
			// log.Println("xv.Interface():", xv.Interface())
		}

		if xv.Type().Kind() != reflect.Slice {
			break
		}
		dShape = append(dShape, xv.Len())

	}
	return dShape
}

func FlattenFloat64(d interface{}) []float64 {
	if fd, ok := d.([]float64); ok {
		return fd
	}
	if fd, ok := d.([]float32); ok {
		return F32ToF64(fd)
	}
	var ret []float64
	dv := reflect.ValueOf(d)
	for i := 0; i < dv.Len(); i++ {
		xv := dv.Index(i)
		ret = append(ret, FlattenFloat64(xv.Interface())...)
	}
	return ret
}

func FlattenFloat32(d interface{}) []float32 {
	if fd, ok := d.([]float32); ok {
		return fd
	}
	if fd, ok := d.([]float64); ok {
		return F64ToF32(fd)
	}
	var ret []float32
	dv := reflect.ValueOf(d)
	for i := 0; i < dv.Len(); i++ {
		xv := dv.Index(i)
		ret = append(ret, FlattenFloat32(xv.Interface())...)
	}
	return ret
}

func FlattenUint8(d interface{}) []uint8 {
	if fd, ok := d.([]uint8); ok {
		return fd
	}
	if fd, ok := d.(uint8); ok {
		return []uint8{fd}
	}
	var ret []uint8
	dv := reflect.ValueOf(d)
	for i := 0; i < dv.Len(); i++ {
		xv := dv.Index(i)
		ret = append(ret, FlattenUint8(xv.Interface())...)
	}
	return ret
}

type DataType string

const (
	DataTypeFloat32 DataType = "DataTypeFloat32"
	DataTypeFloat64 DataType = "DataTypeFloat64"
	DataTypeUInt8   DataType = "DataTypeUInt8"
)

// func Reshape(input interface{}, reshape Shape, types DataType) (ret interface{}, err error) {
// 	dShape := GetShapeBySlice(input)

// 	if len(dShape) == 1 {
// 		if len(reshape) == 1 {
// 			if dShape[0] == reshape[0] {
// 				return input, nil
// 			} else {
// 				log.Println(input)
// 				err = errors.New(fmt.Sprintf("a len == %d , b len == %d , a!=b", dShape[0], reshape[0]))
// 				return
// 			}
// 		}

// 		inputValue := reflect.ValueOf(input)
// 		size := inputValue.Len()
// 		if size%reshape[0] != 0 {
// 			err = errors.New("size%reshape[0]  != 0")
// 			return
// 		}

// 		bachtSize := inputValue.Len() / reshape[0]

// 		retFloat := []interface{}{}
// 		for k := 0; k < reshape[0]; k++ {
// 			itemValue := inputValue.Slice(k*bachtSize, (k+1)*bachtSize)
// 			item := itemValue.Interface()
// 			retItem, err := Reshape(item, reshape[1:], types)
// 			if err != nil {
// 				return nil, err
// 			}
// 			retFloat = append(retFloat, retItem)
// 		}
// 		return retFloat, nil
// 	} else {
// 		if types == DataTypeFloat64 {
// 			return Reshape(FlattenFloat64(input), reshape, types)
// 		}
// 		if types == DataTypeFloat32 {
// 			return Reshape(FlattenFloat32(input), reshape, types)
// 		}
// 		if types == DataTypeUInt8 {
// 			return Reshape(FlattenUint8(input), reshape, types)
// 		}
// 		return nil, errors.New("不支持的类型")
// 	}

// }

func Reshape(input interface{}, reshape Shape, types DataType) (ret interface{}, err error) {

	size := 0
	if types == DataTypeFloat64 {
		tmp := FlattenFloat64(input)
		input = tmp
		size = len(tmp)
	} else if types == DataTypeFloat32 {
		tmp := FlattenFloat32(input)
		input = tmp
		size = len(tmp)
	} else if types == DataTypeUInt8 {
		tmp := FlattenUint8(input)
		input = tmp
		size = len(tmp)
	} else {
		return nil, errors.New("不支持的类型")
	}

	if len(reshape) == 1 {
		if size == reshape[0] {
			// log.Println("dShape[0]:", dShape[0])
			return input, nil
		} else {
			// log.Println(input)
			err = errors.New(fmt.Sprintf("a len == %d , b len == %d , a!=b", size, reshape[0]))
			return
		}
	}

	inputValue := reflect.ValueOf(input)
	if size%reshape[0] != 0 {
		err = errors.New("size%reshape[0]  != 0")
		return
	}

	bachtSize := inputValue.Len() / reshape[0]

	retFloat := []interface{}{}
	for k := 0; k < reshape[0]; k++ {
		itemValue := inputValue.Slice(k*bachtSize, (k+1)*bachtSize)
		item := itemValue.Interface()
		retItem, err := Reshape(item, reshape[1:], types)
		if err != nil {
			return nil, err
		}
		retFloat = append(retFloat, retItem)
	}
	return retFloat, nil

}

func Cosine(a []float64, b []float64) (cosine float64, err error) {
	count := 0
	length_a := len(a)
	length_b := len(b)
	if length_a > length_b {
		count = length_a
	} else {
		count = length_b
	}
	sumA := 0.0
	s1 := 0.0
	s2 := 0.0
	for k := 0; k < count; k++ {
		if k >= length_a {
			s2 += math.Pow(b[k], 2)
			continue
		}
		if k >= length_b {
			s1 += math.Pow(a[k], 2)
			continue
		}
		sumA += a[k] * b[k]
		s1 += math.Pow(a[k], 2)
		s2 += math.Pow(b[k], 2)
	}
	if s1 == 0 || s2 == 0 {
		return 0.0, errors.New("Vectors should not be null (all zeros)")
	}
	return sumA / (math.Sqrt(s1) * math.Sqrt(s2)), nil
}

func To2DF32(input interface{}) (ret [][]float32) {

	xv := reflect.ValueOf(input)
	if xv.Type().Kind() != reflect.Slice {
		return nil
	}
	l1 := xv.Len()
	for i := 0; i < l1; i++ {
		xv2 := xv.Index(i)

		items := FlattenFloat32(xv2.Interface())

		ret = append(ret, items)
	}
	return
}

func ReshapeTo2DF32(input interface{}, reshape Shape) (ret [][]float32) {
	tmp, err := Reshape(input, reshape, DataTypeFloat32)
	if err != nil {
		log.Println("ReshapeTo2DF32 Reshape error:", err)
		return nil
	}
	return To2DF32(tmp)
}

func To3DF32(input interface{}) (ret [][][]float32) {

	xv := reflect.ValueOf(input)
	if xv.Type().Kind() != reflect.Slice {
		return nil
	}

	for i := 0; i < xv.Len(); i++ {
		xv2 := xv.Index(i)
		ret = append(ret, To2DF32(xv2.Interface()))
	}
	return
}

func ReshapeTo3DF32(input interface{}, reshape Shape) (ret [][][]float32) {
	tmp, err := Reshape(input, reshape, DataTypeFloat32)
	if err != nil {
		log.Println("ReshapeTo3DF32 Reshape error:", err)
		return nil
	}
	return To3DF32(tmp)
}

func To4DF32(input interface{}) (ret [][][][]float32) {
	xv := reflect.ValueOf(input)
	if xv.Type().Kind() != reflect.Slice {
		return nil
	}

	for i := 0; i < xv.Len(); i++ {
		xv2 := xv.Index(i)
		ret = append(ret, To3DF32(xv2.Interface()))
	}
	return
}

func To2DU8(input interface{}) (ret [][]uint8) {

	xv := reflect.ValueOf(input)
	if xv.Type().Kind() != reflect.Slice {
		return nil
	}
	l1 := xv.Len()
	for i := 0; i < l1; i++ {
		xv2 := xv.Index(i)

		items := FlattenUint8(xv2.Interface())

		ret = append(ret, items)
	}
	return
}

func To3DU8(input interface{}) (ret [][][]uint8) {

	xv := reflect.ValueOf(input)
	if xv.Type().Kind() != reflect.Slice {
		return nil
	}

	for i := 0; i < xv.Len(); i++ {
		xv2 := xv.Index(i)
		ret = append(ret, To2DU8(xv2.Interface()))
	}
	return
}

// ToSliceE casts an interface to a []interface{} type.
func ToSlice(oval interface{}) []interface{} {
	var s []interface{}

	switch v := oval.(type) {
	case []interface{}:
		return append(s, v...)
	case []float32:
		for i := range v {
			s = append(s, v[i])
		}
	case []float64:
		for i := range v {
			s = append(s, v[i])
		}
	case []int64:
		for i := range v {
			s = append(s, v[i])
		}
	case []int32:
	case []int:
	case []uint:
	case []uint8:
	case []uint16:
	case []uint32:
	case []uint64:
		for i := range v {
			s = append(s, v[i])
		}
	default:
	}
	// log.Println("v:", oval)
	return s
}

func Softmax(arr interface{}) int {
	list := ToSlice(arr)
	maxIndex := -1
	var maxNum float64 = -99999999
	for i := range list {
		n := cast.ToFloat64(fmt.Sprint(list[i]))
		// log.Println(n)
		if n > maxNum {
			maxNum = n
			maxIndex = i
		}

	}
	return maxIndex
}
func Mean(v []float64) float64 {
	var res float64 = 0
	var n int = len(v)
	for i := 0; i < n; i++ {
		res += v[i]
	}
	return res / float64(n)
}

func Variance(v []float64) float64 {
	var res float64 = 0
	var m = Mean(v)
	var n int = len(v)
	for i := 0; i < n; i++ {
		res += (v[i] - m) * (v[i] - m)
	}
	return res / float64(n-1)
}
func Std(v []float64) float64 {
	return math.Sqrt(Variance(v))
}

type Stack []byte

func (s *Stack) Push(v byte) {
	*s = append(*s, v)
}
func (s *Stack) Pop() byte {
	if len(*s) == 0 {
		return 0
	}
	len := len(*s)
	v := (*s)[len-1]
	*s = (*s)[:len-1]
	return v
}

func (s *Stack) Look() byte {
	if len(*s) == 0 {
		return 0
	}
	l := len(*s)
	return (*s)[l-1]
}

func (s *Stack) Len() int {

	return len(*s)
}

func F32ToF64(f32 []float32) (f64 []float64) {
	for i := range f32 {
		f64 = append(f64, float64(f32[i]))
	}
	return
}

func F64ToF32(f64 []float64) (f32 []float32) {
	for i := range f64 {
		f32 = append(f32, float32(f64[i]))
	}
	return
}
