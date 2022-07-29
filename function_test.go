package kgo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFun_lenArrayOrSlice(t *testing.T) {
	var res int

	res = lenArrayOrSlice(naturalArr, 3)
	assert.Greater(t, res, 0)

	res = lenArrayOrSlice(naturalArr, 1)
	assert.Greater(t, res, 0)

	res = lenArrayOrSlice(naturalArr, 2)
	assert.Equal(t, res, -1)

	res = lenArrayOrSlice(naturalArr, 9)
	assert.Greater(t, res, 0)
}

func TestFun_isNil(t *testing.T) {
	var s []int
	var i interface{}
	var res bool

	res = isNil(s)
	assert.True(t, res)

	res = isNil(&s)
	assert.False(t, res)

	res = isNil(i)
	assert.True(t, res)
}

func TestFun_GetVariateType(t *testing.T) {
	var res string

	res = GetVariateType(1)
	assert.Equal(t, "int", res)

	res = GetVariateType(intAstronomicalUnit)
	assert.Equal(t, "int64", res)

	res = GetVariateType(flPi1)
	assert.Equal(t, "float32", res)

	res = GetVariateType(floAvogadro)
	assert.Equal(t, "float64", res)

	res = GetVariateType(strHello)
	assert.Equal(t, "string", res)

	res = GetVariateType(true)
	assert.Equal(t, "bool", res)

	res = GetVariateType(rune('你'))
	assert.Equal(t, "int32", res)

	res = GetVariateType('你')
	assert.Equal(t, "int32", res)

	res = GetVariateType([]byte("你好"))
	assert.Equal(t, "[]uint8", res)
}

func BenchmarkFun_GetVariateType(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetVariateType(intAstronomicalUnit)
	}
}

func TestFun_GetVariatePointerAddr(t *testing.T) {
	var tests = []struct {
		input    interface{}
		expected float64
	}{
		{intSpeedLight, 0},
		{strHello, 0},
		{crowd, 0},
	}
	for _, test := range tests {
		actual := GetVariatePointerAddr(test.input)
		assert.Greater(t, actual, int64(test.expected))
	}
}

func BenchmarkFun_GetVariatePointerAddr(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetVariatePointerAddr(intSpeedLight)
	}
}

func TestFun_IsPointer(t *testing.T) {
	var chk bool

	//非指针
	chk = IsPointer(itfObj, false)
	assert.False(t, chk)

	//指针
	chk = IsPointer(orgS1, false)
	assert.True(t, chk)

	//非nil指针
	chk = IsPointer(orgS1, true)
	assert.True(t, chk)

	//空指针
	chk = IsPointer(itfObj, true)
	assert.False(t, chk)
}

func BenchmarkFun_IsPointer(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsPointer(orgS1, true)
	}
}
