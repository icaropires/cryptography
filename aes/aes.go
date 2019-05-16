package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	BLOCK_SIZE_BYTES        = 16
	BYTE_SIZE_BITS          = 8
	WORD_SIZE_BYTES         = 32
	KEY_SIZE_BYTES          = 4
	EXPANDED_KEY_SIZE_WORDS = 44
	STATE_SIZE_ROWS         = 4
	NUMBER_OF_ROUNDS        = 10
)

var rcon = []uint32{0x01000000, 0x02000000, 0x04000000, 0x08000000, 0x10000000, 0x20000000, 0x40000000, 0x80000000, 0x1b000000, 0x36000000}

// FIPS-197 Figure 7. S-box substitution values in hexadecimal format.
var sbox0 = [][]byte{
	[]byte{0x63, 0x7c, 0x77, 0x7b, 0xf2, 0x6b, 0x6f, 0xc5, 0x30, 0x01, 0x67, 0x2b, 0xfe, 0xd7, 0xab, 0x76},
	[]byte{0xca, 0x82, 0xc9, 0x7d, 0xfa, 0x59, 0x47, 0xf0, 0xad, 0xd4, 0xa2, 0xaf, 0x9c, 0xa4, 0x72, 0xc0},
	[]byte{0xb7, 0xfd, 0x93, 0x26, 0x36, 0x3f, 0xf7, 0xcc, 0x34, 0xa5, 0xe5, 0xf1, 0x71, 0xd8, 0x31, 0x15},
	[]byte{0x04, 0xc7, 0x23, 0xc3, 0x18, 0x96, 0x05, 0x9a, 0x07, 0x12, 0x80, 0xe2, 0xeb, 0x27, 0xb2, 0x75},
	[]byte{0x09, 0x83, 0x2c, 0x1a, 0x1b, 0x6e, 0x5a, 0xa0, 0x52, 0x3b, 0xd6, 0xb3, 0x29, 0xe3, 0x2f, 0x84},
	[]byte{0x53, 0xd1, 0x00, 0xed, 0x20, 0xfc, 0xb1, 0x5b, 0x6a, 0xcb, 0xbe, 0x39, 0x4a, 0x4c, 0x58, 0xcf},
	[]byte{0xd0, 0xef, 0xaa, 0xfb, 0x43, 0x4d, 0x33, 0x85, 0x45, 0xf9, 0x02, 0x7f, 0x50, 0x3c, 0x9f, 0xa8},
	[]byte{0x51, 0xa3, 0x40, 0x8f, 0x92, 0x9d, 0x38, 0xf5, 0xbc, 0xb6, 0xda, 0x21, 0x10, 0xff, 0xf3, 0xd2},
	[]byte{0xcd, 0x0c, 0x13, 0xec, 0x5f, 0x97, 0x44, 0x17, 0xc4, 0xa7, 0x7e, 0x3d, 0x64, 0x5d, 0x19, 0x73},
	[]byte{0x60, 0x81, 0x4f, 0xdc, 0x22, 0x2a, 0x90, 0x88, 0x46, 0xee, 0xb8, 0x14, 0xde, 0x5e, 0x0b, 0xdb},
	[]byte{0xe0, 0x32, 0x3a, 0x0a, 0x49, 0x06, 0x24, 0x5c, 0xc2, 0xd3, 0xac, 0x62, 0x91, 0x95, 0xe4, 0x79},
	[]byte{0xe7, 0xc8, 0x37, 0x6d, 0x8d, 0xd5, 0x4e, 0xa9, 0x6c, 0x56, 0xf4, 0xea, 0x65, 0x7a, 0xae, 0x08},
	[]byte{0xba, 0x78, 0x25, 0x2e, 0x1c, 0xa6, 0xb4, 0xc6, 0xe8, 0xdd, 0x74, 0x1f, 0x4b, 0xbd, 0x8b, 0x8a},
	[]byte{0x70, 0x3e, 0xb5, 0x66, 0x48, 0x03, 0xf6, 0x0e, 0x61, 0x35, 0x57, 0xb9, 0x86, 0xc1, 0x1d, 0x9e},
	[]byte{0xe1, 0xf8, 0x98, 0x11, 0x69, 0xd9, 0x8e, 0x94, 0x9b, 0x1e, 0x87, 0xe9, 0xce, 0x55, 0x28, 0xdf},
	[]byte{0x8c, 0xa1, 0x89, 0x0d, 0xbf, 0xe6, 0x42, 0x68, 0x41, 0x99, 0x2d, 0x0f, 0xb0, 0x54, 0xbb, 0x16},
}

var sbox1 = [][]byte{
	[]byte{0x52, 0x09, 0x6a, 0xd5, 0x30, 0x36, 0xa5, 0x38, 0xbf, 0x40, 0xa3, 0x9e, 0x81, 0xf3, 0xd7, 0xfb},
	[]byte{0x7c, 0xe3, 0x39, 0x82, 0x9b, 0x2f, 0xff, 0x87, 0x34, 0x8e, 0x43, 0x44, 0xc4, 0xde, 0xe9, 0xcb},
	[]byte{0x54, 0x7b, 0x94, 0x32, 0xa6, 0xc2, 0x23, 0x3d, 0xee, 0x4c, 0x95, 0x0b, 0x42, 0xfa, 0xc3, 0x4e},
	[]byte{0x08, 0x2e, 0xa1, 0x66, 0x28, 0xd9, 0x24, 0xb2, 0x76, 0x5b, 0xa2, 0x49, 0x6d, 0x8b, 0xd1, 0x25},
	[]byte{0x72, 0xf8, 0xf6, 0x64, 0x86, 0x68, 0x98, 0x16, 0xd4, 0xa4, 0x5c, 0xcc, 0x5d, 0x65, 0xb6, 0x92},
	[]byte{0x6c, 0x70, 0x48, 0x50, 0xfd, 0xed, 0xb9, 0xda, 0x5e, 0x15, 0x46, 0x57, 0xa7, 0x8d, 0x9d, 0x84},
	[]byte{0x90, 0xd8, 0xab, 0x00, 0x8c, 0xbc, 0xd3, 0x0a, 0xf7, 0xe4, 0x58, 0x05, 0xb8, 0xb3, 0x45, 0x06},
	[]byte{0xd0, 0x2c, 0x1e, 0x8f, 0xca, 0x3f, 0x0f, 0x02, 0xc1, 0xaf, 0xbd, 0x03, 0x01, 0x13, 0x8a, 0x6b},
	[]byte{0x3a, 0x91, 0x11, 0x41, 0x4f, 0x67, 0xdc, 0xea, 0x97, 0xf2, 0xcf, 0xce, 0xf0, 0xb4, 0xe6, 0x73},
	[]byte{0x96, 0xac, 0x74, 0x22, 0xe7, 0xad, 0x35, 0x85, 0xe2, 0xf9, 0x37, 0xe8, 0x1c, 0x75, 0xdf, 0x6e},
	[]byte{0x47, 0xf1, 0x1a, 0x71, 0x1d, 0x29, 0xc5, 0x89, 0x6f, 0xb7, 0x62, 0x0e, 0xaa, 0x18, 0xbe, 0x1b},
	[]byte{0xfc, 0x56, 0x3e, 0x4b, 0xc6, 0xd2, 0x79, 0x20, 0x9a, 0xdb, 0xc0, 0xfe, 0x78, 0xcd, 0x5a, 0xf4},
	[]byte{0x1f, 0xdd, 0xa8, 0x33, 0x88, 0x07, 0xc7, 0x31, 0xb1, 0x12, 0x10, 0x59, 0x27, 0x80, 0xec, 0x5f},
	[]byte{0x60, 0x51, 0x7f, 0xa9, 0x19, 0xb5, 0x4a, 0x0d, 0x2d, 0xe5, 0x7a, 0x9f, 0x93, 0xc9, 0x9c, 0xef},
	[]byte{0xa0, 0xe0, 0x3b, 0x4d, 0xae, 0x2a, 0xf5, 0xb0, 0xc8, 0xeb, 0xbb, 0x3c, 0x83, 0x53, 0x99, 0x61},
	[]byte{0x17, 0x2b, 0x04, 0x7e, 0xba, 0x77, 0xd6, 0x26, 0xe1, 0x69, 0x14, 0x63, 0x55, 0x21, 0x0c, 0x7d},
}

func getEmptyBlock() []byte {
	return make([]byte, BLOCK_SIZE_BYTES)
}

func copyToState(block []byte) [][]byte {
	state := make([][]byte, STATE_SIZE_ROWS)

	nb := (len(block) * BYTE_SIZE_BITS) / 32
	for i := 0; i < STATE_SIZE_ROWS; i++ {
		state[i] = make([]byte, nb)
	}

	for r := 0; r < STATE_SIZE_ROWS; r++ {
		for c := 0; c < nb; c++ {
			state[r][c] = block[r+STATE_SIZE_ROWS*c]
		}
	}

	return state
}

func copyFromState(state [][]byte) []byte {
	block := make([]byte, BLOCK_SIZE_BYTES)

	nb := (len(block) * BYTE_SIZE_BITS) / 32
	for r := 0; r < STATE_SIZE_ROWS; r++ {
		for c := 0; c < nb; c++ {
			block[r+STATE_SIZE_ROWS*c] = state[r][c]
		}
	}

	return block
}

func rotWord(word uint32) uint32 {
	return word<<BYTE_SIZE_BITS | word>>(WORD_SIZE_BYTES-BYTE_SIZE_BITS)
}

func subByte(b byte) byte {
	x, y := (b&0xf0)>>4, b&0xf

	return sbox0[x][y]
}

func subBytes(state [][]byte) [][]byte {
	nb := len(state[0])
	for r := 0; r < STATE_SIZE_ROWS; r++ {
		for c := 0; c < nb; c++ {
			state[r][c] = subByte(state[r][c])
		}
	}

	return state
}

func invSubByte(b byte) byte {
	x, y := (b&0xf0)>>4, b&0xf

	return sbox1[x][y]
}

func invSubBytes(state [][]byte) [][]byte {
	nb := len(state[0])
	for r := 0; r < STATE_SIZE_ROWS; r++ {
		for c := 0; c < nb; c++ {
			state[r][c] = invSubByte(state[r][c])
		}
	}

	return state
}

func subWord(word uint32) uint32 {
	var result_word uint32 = 0

	for shift_hex := uint32(0); shift_hex < WORD_SIZE_BYTES; shift_hex += 8 {
		b := byte(word & (0xff << shift_hex) >> shift_hex)
		result_word |= uint32(subByte(b)) << shift_hex
	}

	return result_word
}

func shiftRows(state [][]byte) [][]byte {
	for r := 0; r < STATE_SIZE_ROWS; r++ {
		state[r] = append(state[r][r:], state[r][:r]...)
	}

	return state
}

func invShiftRows(state [][]byte) [][]byte {
	nb := len(state[0])

	for r := 0; r < STATE_SIZE_ROWS; r++ {
		state[r] = append(state[r][nb-r:], state[r][:nb-r]...)
	}

	return state
}

func expandKey(key []byte) []uint32 {
	words := make([]uint32, EXPANDED_KEY_SIZE_WORDS)

	for i := uint32(0); i < KEY_SIZE_BYTES; i++ {
		count := uint32(0)
		for shift_hex := uint32(8); shift_hex <= WORD_SIZE_BYTES; shift_hex += 8 {
			words[i] |= uint32(key[KEY_SIZE_BYTES*i+count]) << (WORD_SIZE_BYTES - shift_hex)
			count++
		}
	}

	for i := uint32(KEY_SIZE_BYTES); i < EXPANDED_KEY_SIZE_WORDS; i++ {
		temp := words[i-1]

		if (i % KEY_SIZE_BYTES) == 0 {
			temp = subWord(rotWord(temp)) ^ rcon[(i/KEY_SIZE_BYTES)-1]
		}

		words[i] = words[i-KEY_SIZE_BYTES] ^ temp
	}

	return words
}

func xtime(value byte) byte {
	if value&0x80 != 0 {
		return (((value << 1) ^ 0x1B) & 0xFF)
	}
	return (value << 1)
}

func mixColumns(state [][]byte) [][]byte {
	newState := make([][]byte, STATE_SIZE_ROWS)
	for i := 0; i < STATE_SIZE_ROWS; i++ {
		newState[i] = make([]byte, STATE_SIZE_ROWS)
	}
	for i := uint32(0); i < STATE_SIZE_ROWS; i++ {
		newState[0][i] = xtime(state[0][i]) ^ (state[1][i] ^ xtime(state[1][i])) ^ state[2][i] ^ state[3][i]
		newState[1][i] = state[0][i] ^ xtime(state[1][i]) ^ (state[2][i] ^ xtime(state[2][i])) ^ state[3][i]
		newState[2][i] = state[0][i] ^ state[1][i] ^ xtime(state[2][i]) ^ (state[3][i] ^ xtime(state[3][i]))
		newState[3][i] = (state[0][i] ^ xtime(state[0][i])) ^ state[1][i] ^ state[2][i] ^ xtime(state[3][i])
	}
	return newState
}

func mult(value byte, coef uint32) byte {
	switch coef {
	case 9:
		return xtime(xtime(xtime(value))) ^ value
	case 11:
		return xtime(xtime(xtime(value))^value) ^ value
	case 13:
		return xtime(xtime(xtime(value)^value)) ^ value
	case 14:
		return xtime(xtime(xtime(value)^value) ^ value)
	default:
		return 0x0
	}
}
func invMixColumns(state [][]byte) [][]byte {
	newState := make([][]byte, STATE_SIZE_ROWS)
	for i := 0; i < STATE_SIZE_ROWS; i++ {
		newState[i] = make([]byte, STATE_SIZE_ROWS)
	}
	for i := uint32(0); i < STATE_SIZE_ROWS; i++ {
		newState[0][i] = mult(state[0][i], 14) ^ mult(state[1][i], 11) ^ mult(state[2][i], 13) ^ mult(state[3][i], 9)
		newState[1][i] = mult(state[0][i], 9) ^ mult(state[1][i], 14) ^ mult(state[2][i], 11) ^ mult(state[3][i], 13)
		newState[2][i] = mult(state[0][i], 13) ^ mult(state[1][i], 9) ^ mult(state[2][i], 14) ^ mult(state[3][i], 11)
		newState[3][i] = mult(state[0][i], 11) ^ mult(state[1][i], 13) ^ mult(state[2][i], 9) ^ mult(state[3][i], 14)
	}
	return newState
}

func addRoundKey(state [][]byte, key []byte) [][]byte {
	for r := uint32(0); r < STATE_SIZE_ROWS; r++ {
		for c := uint32(0); c < STATE_SIZE_ROWS; c++ {
			state[r][c] ^= key[r+STATE_SIZE_ROWS*c]
		}
	}

	return state
}

func wordToByte(keys []uint32, start int) []byte {
	var key_bytes []byte

	for i := start; i < start+4; i++ {
		key_bytes_aux := make([]byte, 4)
		binary.BigEndian.PutUint32(key_bytes_aux, keys[i])

		key_bytes = append(key_bytes, key_bytes_aux...)
	}

	return key_bytes
}

func textFromFile(filepath string) [][]byte {
	file, err := ioutil.ReadFile(filepath)

	if err != nil {
		panic("Não foi possível ler do arquivo")
	}

	var blocks [][]byte
	for i := 0; i <= len(file)-16; i += 16 {
		blocks = append(blocks, file[i:i+16])
	}

	return blocks
}

func writeBlocksToFile(filename string, blocks [][]byte) {
	var result_block []byte

	for _, block := range blocks {
		result_block = append(result_block, block...)
	}

	erro := ioutil.WriteFile(filename, result_block, 0644)

	if erro != nil {
		panic("Não foi possível escrever no arquivo")
	}
}

func Encrypt(block []byte, key []byte) []byte {
	var nb = (len(block) * BYTE_SIZE_BITS) / 32

	keys := expandKey(key)
	state := copyToState(block)

	state = addRoundKey(state, wordToByte(keys, 0))

	for round := 1; round < NUMBER_OF_ROUNDS; round++ {
		state = subBytes(state)
		state = shiftRows(state)
		state = mixColumns(state)
		state = addRoundKey(state, wordToByte(keys, round*nb))
	}

	state = subBytes(state)
	state = shiftRows(state)
	state = addRoundKey(state, wordToByte(keys, NUMBER_OF_ROUNDS*nb))

	return copyFromState(state)
}

func Decrypt(cyphertext []byte, key []byte) []byte {
	var nb = (len(cyphertext) * BYTE_SIZE_BITS) / 32

	keys := expandKey(key)
	state := copyToState(cyphertext)

	state = addRoundKey(state, wordToByte(keys, NUMBER_OF_ROUNDS*nb))

	for round := NUMBER_OF_ROUNDS - 1; round > 0; round-- {
		state = invShiftRows(state)
		state = invSubBytes(state)
		state = addRoundKey(state, wordToByte(keys, round*nb))
		state = invMixColumns(state)
	}

	state = invShiftRows(state)
	state = invSubBytes(state)
	state = addRoundKey(state, wordToByte(keys, 0))

	return copyFromState(state)

}

func ofb(block []byte, iv []byte, key []byte) []byte {
	ciphertext := make([]byte, BLOCK_SIZE_BYTES)
	result := make([]byte, BLOCK_SIZE_BYTES)
	if len(iv) > 0 {
		for i := uint32(0); i < BLOCK_SIZE_BYTES; i++ {
			result = Encrypt(iv, key)
		}
	} else {
		for i := uint32(0); i < BLOCK_SIZE_BYTES; i++ {
			result = Encrypt(result, key)
		}
	}
	for i := uint32(0); i < BLOCK_SIZE_BYTES; i++ {
		ciphertext[i] = result[i] ^ block[i]
	}
	return ciphertext
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("\nUsage example: go run aes.go [d|e] [filename] [key]\n")
		panic("Wrong usage!")
	}

	op, filename, key_string := os.Args[1], os.Args[2], os.Args[3]

	if len(key_string) < 16 {
		key_string += strings.Repeat(" ", 16-len(key_string))
	}

	key := []byte(key_string)

	if len(key) > 16 {
		key = key[:16]
	}

	blocks := textFromFile(filename)
	var cipher_blocks [][]byte

	if op == "e" {
		for i, block := range blocks {
			if i == 0 {
				cipher_blocks = append(cipher_blocks, ofb(block, key, key))
			} else {
				cipher_blocks = append(cipher_blocks, ofb(block, []byte{}, key))
			}
		}

		writeBlocksToFile(filename+".cipher", cipher_blocks)
	} else if op == "d" {
		for i, block := range blocks {
			if i == 0 {
				cipher_blocks = append(cipher_blocks, ofb(block, key, key))
			} else {
				cipher_blocks = append(cipher_blocks, ofb(block, []byte{}, key))
			}
		}
		writeBlocksToFile(filename+".deciphered", cipher_blocks)
	} else {
		fmt.Println("\nInvalid operation!")
		panic("Invalid operation")
	}
}
