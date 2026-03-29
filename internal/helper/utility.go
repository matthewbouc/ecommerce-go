package helper

import (
	"crypto/rand"
	"strconv"
)

func RandomNumbers(length int) (int, error) {
	const digits = "0123456789"
	const base = byte(len(digits))       // 10
	const maxUnbiased = 256 - (256 % 10) // 250 — largest multiple of 10 that fits in a byte

	result := make([]byte, length)
	i := 0
	for i < length {
		var b [1]byte
		if _, err := rand.Read(b[:]); err != nil {
			return 0, err
		}
		if b[0] < maxUnbiased {
			result[i] = digits[b[0]%base]
			i++
		}
	}

	return strconv.Atoi(string(result))
}
