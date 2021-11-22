package utils

import (
	"bytes"
	"crypto/rand"
	"math/big"
)

const (
	charTables = "abcdefghijklmnopqrstuvwxyz1234567890"
)

func CreateRandomString(len int) string {
	b := bytes.NewBufferString(charTables)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	var result string
	for i := 0; i < len; i++ {
		randomInt, err := rand.Int(rand.Reader, bigInt)
		if err != nil {
			continue
		}
		result += string(charTables[randomInt.Int64()])
	}
	return result
}

