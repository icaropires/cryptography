package main

import (
	//"encoding/base64"
	"encoding/binary"

	//"fmt"
	"io/ioutil"
)

const ROUNDS = 20

func rotl(x, n uint32) uint32 {
	x &= 0xffffffff
	return (x << n) | (x>>(32-n))&0xffffffff
}

func quarter_round(a, b, c, d *uint32) {
	*b ^= rotl(*a+*d, 7)
	*c ^= rotl(*b+*a, 9)
	*d ^= rotl(*c+*b, 13)
	*a ^= rotl(*d+*c, 18)
}

func salsa20_block(in []uint32) []uint32 {
	x := make([]uint32, 16)

	for i := range x {
		x[i] = in[i]
	}

	for i := 0; i < ROUNDS; i += 2 {
		quarter_round(&x[0], &x[4], &x[8], &x[12])
		quarter_round(&x[5], &x[9], &x[13], &x[1])
		quarter_round(&x[10], &x[14], &x[2], &x[6])
		quarter_round(&x[15], &x[3], &x[7], &x[11])

		quarter_round(&x[0], &x[1], &x[2], &x[3])
		quarter_round(&x[5], &x[6], &x[7], &x[4])
		quarter_round(&x[10], &x[11], &x[8], &x[9])
		quarter_round(&x[15], &x[12], &x[13], &x[14])
	}

	out := make([]uint32, 16)
	for i := range x {
		out[i] = x[i] + in[i]
	}

	return out
}

func get_block(nonce, pos, key []uint32) []byte {
	state := make([]uint32, 0, 16)
	cons := make([]uint32, 4)

	// cons = expand 32-byte key
	cons[0], cons[1], cons[2], cons[3] = 1634760805, 857760878, 2036477234, 1797285236

	// Fill salsa state
	state = append(state, cons[0], key[0], key[1], key[2])
	state = append(state, key[3], cons[1], nonce[0], nonce[1])
	state = append(state, pos[0], pos[1], cons[2], key[4])
	state = append(state, key[5], key[6], key[7], cons[3])

	state = salsa20_block(state)

	salsa20_block_bytes := make([]byte, 16*4)
	for i, word := range state {
		binary.LittleEndian.PutUint32(salsa20_block_bytes[i*4:i*4+4], word)
	}

	return salsa20_block_bytes
}

func stream_salsa20(length int, nonce, key []byte) []byte {
	nonce_ints := make([]uint32, 0, len(nonce)/4)
	key_ints := make([]uint32, 0, len(key)/4)

	for i := 0; i < len(nonce)/4; i++ {
		nonce_ints = append(nonce_ints, binary.LittleEndian.Uint32(nonce[i*4:i*4+4]))
	}

	for i := 0; i < len(key)/4; i++ {
		key_ints = append(key_ints, binary.LittleEndian.Uint32(key[i*4:i*4+4]))
	}

	output := make([]byte, 0, length)
	for i := 0; i < (length+63)/64; i++ {
		output = append(output, get_block(nonce_ints, []uint32{uint32(i & 0xffffffff), uint32(i >> 32)}, key_ints)...)
	}

	return output
}

func main() {
	key := make([]byte, 32)
	out := stream_salsa20(512, []byte{4, 2, 4, 2, 4, 2, 4, 2}, key)

	ioutil.WriteFile("out", out, 0644)
	//fmt.Println(base64.StdEncoding.EncodeToString(out))
}
