// +build windows

package kgo

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOS_IsWindows(t *testing.T) {
	res := KOS.IsWindows()
	assert.True(t, res)
}

func BenchmarkOS_IsWindows(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		KOS.IsWindows()
	}
}