package my_protocol

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// for tests
func newAES256WithRandomKey() cipher.Block {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic(err)
	}
	cip, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	return cip
}

func padPKCS7(b []byte) []byte {
	padSize := aes.BlockSize - (len(b) % aes.BlockSize)
	pad := bytes.Repeat([]byte{byte(padSize)}, padSize)
	return append(b, pad...)
}

func unpadPKCS7(b []byte) []byte {
	padSize := int(b[len(b)-1])
	return b[:len(b)-padSize]
}
