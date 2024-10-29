package internal

import (
	"crypto/rand"
	"strings"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var charLen = byte(len(chars))

func RandString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)

	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	for i := 0; i < n; i++ {
		sb.WriteByte(chars[b[i]%charLen])
	}

	return sb.String()
}
