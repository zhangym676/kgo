package gohelper

import (
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)

// Int2Str 将整数转换为字符串
func (kc *LkkConvert) Int2Str(val interface{}) string {
	switch val.(type) {
	// Integers
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val)
	// Type is not integers, return empty string
	default:
		return ""
	}
}

// Float2Str 将浮点数转换为字符串,length为小数位数
func (kc *LkkConvert) Float2Str(val interface{}, length int) string {
	switch val.(type) {
	// Floats
	case float32:
		return strconv.FormatFloat(float64(val.(float32)), 'f', length, 32)
	case float64:
		return strconv.FormatFloat(val.(float64), 'f', length, 64)
	// Type is not floats, return empty string
	default:
		return ""
	}
}

// Bool2Str 将布尔值转换为字符串
func (kc *LkkConvert) Bool2Str(val bool) string {
	if val {
		return "true"
	}
	return "false"
}

// StrictStr2Int 严格将字符串转换为有符号整型,bitSize为类型位数,strict为是否严格检查
func (kc *LkkConvert) StrictStr2Int(val string, bitSize int, strict bool) int64 {
	res, err := strconv.ParseInt(val, 0, bitSize)
	if err != nil {
		if strict {
			panic(err)
		}
	}
	return res
}

// Str2Int 将字符串转换为int
func (kc *LkkConvert) Str2Int(val string) int {
	res, _ := strconv.Atoi(val)
	return res
}

// Str2Int8 将字符串转换为int8
func (kc *LkkConvert) Str2Int8(val string) int8 {
	return int8(kc.StrictStr2Int(val, 8, false))
}

// Str2Int16 将字符串转换为int16
func (kc *LkkConvert) Str2Int16(val string) int16 {
	return int16(kc.StrictStr2Int(val, 16, false))
}

// Str2Int32 将字符串转换为int32
func (kc *LkkConvert) Str2Int32(val string) int32 {
	return int32(kc.StrictStr2Int(val, 32, false))
}

// Str2Int64 将字符串转换为int64
func (kc *LkkConvert) Str2Int64(val string) int64 {
	return kc.StrictStr2Int(val, 64, false)
}

// StrictStr2Uint 严格将字符串转换为无符号整型,bitSize为类型位数,strict为是否严格检查
func (kc *LkkConvert) StrictStr2Uint(val string, bitSize int, strict bool) uint64 {
	res, err := strconv.ParseUint(val, 0, bitSize)
	if err != nil {
		if strict {
			panic(err)
		}
	}
	return res
}

// Str2Uint 将字符串转换为uint
func (kc *LkkConvert) Str2Uint(val string) uint {
	return uint(kc.StrictStr2Uint(val, 0, false))
}

// Str2Uint8 将字符串转换为uint8
func (kc *LkkConvert) Str2Uint8(val string) uint8 {
	return uint8(kc.StrictStr2Uint(val, 8, false))
}

// Str2Uint16 将字符串转换为uint16
func (kc *LkkConvert) Str2Uint16(val string) uint16 {
	return uint16(kc.StrictStr2Uint(val, 16, false))
}

// Str2Uint32 将字符串转换为uint32
func (kc *LkkConvert) Str2Uint32(val string) uint32 {
	return uint32(kc.StrictStr2Uint(val, 32, false))
}

// Str2Uint64 将字符串转换为uint64
func (kc *LkkConvert) Str2Uint64(val string) uint64 {
	return uint64(kc.StrictStr2Uint(val, 64, false))
}

// StrictStr2Float 严格将字符串转换为浮点型,bitSize为类型位数,strict为是否严格检查
func (kc *LkkConvert) StrictStr2Float(val string, bitSize int, strict bool) float64 {
	res, err := strconv.ParseFloat(val, bitSize)
	if err != nil {
		if strict {
			panic(err)
		}
	}
	return res
}

// Str2Float32 将字符串转换为float32
func (kc *LkkConvert) Str2Float32(val string) float32 {
	return float32(kc.StrictStr2Float(val, 32, false))
}

// Str2Float64 将字符串转换为float64
func (kc *LkkConvert) Str2Float64(val string) float64 {
	return float64(kc.StrictStr2Float(val, 64, false))
}

// Str2Bool 将字符串转换为布尔值
func (kc *LkkConvert) Str2Bool(val string) bool {
	if val == "true" || val == "True" || val == "TRUE" {
		return true
	}
	return false
}

// Int2Bool 将整数转换为布尔值
func (kc *LkkConvert) Int2Bool(val interface{}) bool {
	switch val.(type) {
	case int:
		return (val.(int) > 0)
	case int8:
		return (val.(int8) > 0)
	case int16:
		return (val.(int16) > 0)
	case int32:
		return (val.(int32) > 0)
	case int64:
		return (val.(int64) > 0)
	case uint:
		return (val.(uint) > 0)
	case uint8:
		return (val.(uint8) > 0)
	case uint16:
		return (val.(uint16) > 0)
	case uint32:
		return (val.(uint32) > 0)
	case uint64:
		return (val.(uint64) > 0)
	default:
		return false
	}
}

// Str2ByteSlice 将字符串转换为字节切片;该方法零拷贝,但不安全,仅当临时需将长字符串转换且不长时间保存时可以使用.
func (kc *LkkConvert) Str2ByteSlice(val string) []byte {
	pSliceHeader := &reflect.SliceHeader{}
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&val))
	pSliceHeader.Data = strHeader.Data
	pSliceHeader.Len = strHeader.Len
	pSliceHeader.Cap = strHeader.Len
	return *(*[]byte)(unsafe.Pointer(pSliceHeader))
}

// BytesSlice2Str 将字节切片转换为字符串,零拷贝
func (kc *LkkConvert) BytesSlice2Str(val []byte) string {
	return *(*string)(unsafe.Pointer(&val))
}