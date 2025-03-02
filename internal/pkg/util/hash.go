package util

import "hash/fnv"

// HashString 计算字符串的哈希值
func HashString(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
