package main

import (
	"fmt"
)

const BLOCK_SIZE = 16

func getEmptyBlock() []byte {
	return make([]byte, BLOCK_SIZE)
}

func Encrypt(plaintext []byte) []byte {
	return []byte("")
}

func Decrypt(cyphertext []byte) []byte {
	return []byte("")
}

func main() {
	fmt.Println("I'm the super AES!!!")
}
