package main

import (
	"crypto/aes"
	"testing"
)

func TestEncrypt(t *testing.T) {
	key := getEmptyBlock()
	block, _ := aes.NewCipher(key)

	src_right, dst_right := []byte("abobora com pato"), getEmptyBlock()
	block.Encrypt(dst_right, src_right)

	dst_candidate := Encrypt(src_right)

	if string(dst_candidate) != string(dst_right) {
		t.Errorf("Wrong cyphertext %v != %v", dst_candidate, dst_right)
	}
}

func TestDecrypt(t *testing.T) {
	key := getEmptyBlock()
	block, _ := aes.NewCipher(key)

	src_right, dst_right := []byte{71, 118, 98, 53, 113, 103, 151, 136, 112, 28, 251, 73, 140, 147, 165, 44}, getEmptyBlock()
	block.Decrypt(dst_right, src_right)

	dst_candidate := Decrypt(src_right)

	if string(dst_candidate) != string(dst_right) {
		t.Errorf("Wrong plaintext %v != %v", dst_candidate, dst_right)
	}
}
