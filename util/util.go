package util

import "math/rand"

// 生成随机不重复ID
func GID() uint32 {
	return rand.Uint32()
}