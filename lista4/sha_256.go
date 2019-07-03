package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	h0 = uint32(0x6a09e667)
	h1 = uint32(0xbb67ae85)
	h2 = uint32(0x3c6ef372)
	h3 = uint32(0xa54ff53a)
	h4 = uint32(0x510e527f)
	h5 = uint32(0x9b05688c)
	h6 = uint32(0x1f83d9ab)
	h7 = uint32(0x5be0cd19)
)

const (
	BYTE_SIZE_BITS   = 8
	BLOCK_SIZE_BITS  = 512
	BLOCK_SIZE_BYTES = 64
)

var k = []uint32{
	0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,
	0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3, 0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,
	0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
	0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967,
	0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13, 0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85,
	0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
	0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3,
	0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208, 0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2}

func process(blocks [][]uint32) string {
	for _, block := range blocks {
		for i := 16; i < BLOCK_SIZE_BYTES; i++ {
			s0 := rotate(block[i-15], 7) ^ rotate(block[i-15], 18) ^ (block[i-15] >> 3)
			s1 := rotate(block[i-2], 17) ^ rotate(block[i-2], 19) ^ (block[i-2] >> 10)
			block[i] = block[i-16] + s0 + block[i-7] + s1
		}

		a := (h0)
		b := (h1)
		c := (h2)
		d := (h3)
		e := (h4)
		f := (h5)
		g := (h6)
		h := (h7)

		for i := 0; i < BLOCK_SIZE_BYTES; i++ {
			s0 := rotate(uint32(a), 2) ^ rotate(uint32(a), 13) ^ rotate(uint32(a), 22)
			maj := uint32((a & b) ^ (a & c) ^ (b & c))
			temp2 := uint32(s0 + maj)
			s1 := rotate(uint32(e), 6) ^ rotate(uint32(e), 11) ^ rotate(uint32(e), 25)
			ch := uint32((e & f) ^ ((^e) & g))
			temp1 := uint32(uint32(h) + s1 + ch + k[i] + block[i])

			h = g
			g = f
			f = e
			e = uint32(uint32(d) + temp1)
			d = c
			c = b
			b = a
			a = uint32(temp1 + temp2)
		}
		h0 = (h0 + a)
		h1 = (h1 + b)
		h2 = (h2 + c)
		h3 = (h3 + d)
		h4 = (h4 + e)
		h5 = (h5 + f)
		h6 = (h6 + g)
		h7 = (h7 + h)

	}
	hashParts := []uint32{h0, h1, h2, h3, h4, h5, h6, h7}
	var hash string
	for _, h := range hashParts {
		hash += fmt.Sprintf("%08x", h)
	}
	return hash
}

func rotate(word, size uint32) uint32 {
	return word>>size | word<<(32-size)
}

func textFromFile(filepath string) [][]uint32 {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic("Não foi possível ler do arquivo")
	}

	message := file
	paddedMessage := message

	paddingLength := BLOCK_SIZE_BYTES - len(message)%BLOCK_SIZE_BYTES - BYTE_SIZE_BITS

	for i := 0; i < paddingLength; i++ {
		if i == 0 {
			paddedMessage = append(message, 128)
		} else {
			paddedMessage = append(paddedMessage, 0)
		}
	}

	messageLength := make([]byte, 8)
	binary.BigEndian.PutUint64(messageLength, 8*uint64(len(message)))
	finalMessage := append(paddedMessage, messageLength...)

	var blocks [][]uint32
	for i := 0; i < len(finalMessage)/64; i++ {
		var words []uint32
		for j := 0; j < 64; j++ {
			if j < 16 {
				words = append(words, binary.BigEndian.Uint32(finalMessage[(i*64+j*4):(i*64+j*4+4)]))
			} else {
				words = append(words, 0)
			}
		}
		blocks = append(blocks, words)
	}

	return blocks
}

func main() {

	filename := os.Args[1]
	blocks := textFromFile(filename)

	hash := process(blocks)

	fmt.Printf("SHA256(%s)= %s\n", filename, hash)
}
