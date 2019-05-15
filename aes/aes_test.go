package main

import (
	"crypto/aes"
	"fmt"
	"testing"
)

func TestCopyToState(t *testing.T) {
	block := []byte{
		0xea, 0x04, 0x65, 0x85, 0x83, 0x45, 0x5d, 0x96,
		0x5c, 0x33, 0x98, 0xb0, 0xf0, 0x2d, 0xad, 0xc5,
	}

	state_right_bytes := [][]byte{
		[]byte{0xea, 0x83, 0x5c, 0xf0},
		[]byte{0x04, 0x45, 0x33, 0x2d},
		[]byte{0x65, 0x5d, 0x98, 0xad},
		[]byte{0x85, 0x96, 0xb0, 0xc5},
	}
	state_candidate_bytes := copyToState(block)

	state_right := fmt.Sprintf("%v", state_right_bytes)
	state_candidate := fmt.Sprintf("%v", state_candidate_bytes)

	if state_right != state_candidate {
		t.Errorf("Wrong state %v != %v", state_candidate, state_right)
	}
}

func TestSubByte(t *testing.T) {
	var sbox_value_right byte = 0x5d
	sbox_value_candidate := subByte(0x8d)

	if sbox_value_right != sbox_value_candidate {
		t.Errorf("Wrong subByte %v != %v", sbox_value_candidate, sbox_value_right)
	}
}

func TestSubWord(t *testing.T) {
	subword_right := uint32(0x5da515d2)
	subword_candidate := subWord(0x8d292f7f)

	if subword_right != subword_candidate {
		t.Errorf("Wrong subWord %v != %v", subword_candidate, subword_right)
	}
}

func TestSubBytes(t *testing.T) {
	state := [][]byte{
		[]byte{0xea, 0x83, 0x5c, 0xf0},
		[]byte{0x04, 0x45, 0x33, 0x2d},
		[]byte{0x65, 0x5d, 0x98, 0xad},
		[]byte{0x85, 0x96, 0xb0, 0xc5},
	}

	state_right_bytes := [][]byte{
		[]byte{0x87, 0xec, 0x4a, 0x8c},
		[]byte{0xf2, 0x6e, 0xc3, 0xd8},
		[]byte{0x4d, 0x4c, 0x46, 0x95},
		[]byte{0x97, 0x90, 0xe7, 0xa6},
	}
	state_candidate_bytes := subBytes(state)

	state_right := fmt.Sprintf("%v", state_right_bytes)
	state_candidate := fmt.Sprintf("%v", state_candidate_bytes)

	if state_right != state_candidate {
		t.Errorf("Wrong bytes substitution! %v != %v", state_candidate, state_right)
	}
}

func TestShiftRows(t *testing.T) {
	state := [][]byte{
		[]byte{0xd4, 0xe0, 0xb8, 0x1e},
		[]byte{0x27, 0xbf, 0xb4, 0x41},
		[]byte{0x11, 0x98, 0x5d, 0x52},
		[]byte{0xae, 0xf1, 0xe5, 0x30},
	}

	state_right_bytes := [][]byte{
		[]byte{0xd4, 0xe0, 0xb8, 0x1e},
		[]byte{0xbf, 0xb4, 0x41, 0x27},
		[]byte{0x5d, 0x52, 0x11, 0x98},
		[]byte{0x30, 0xae, 0xf1, 0xe5},
	}

	state_candidate_bytes := shiftRows(state)

	state_right := fmt.Sprintf("%v", state_right_bytes)
	state_candidate := fmt.Sprintf("%v", state_candidate_bytes)

	if state_right != state_candidate {
		t.Errorf("Wrong bytes shifting! %v != %v", state_candidate, state_right)
	}
}

func TestAddRoundKey(t *testing.T) {
  block := []byte{
    0x47, 0x37, 0x94, 0xED,
    0x40, 0xD4, 0xE4, 0xA5,
    0xA3, 0x70, 0x3A, 0xA6,
    0x4C, 0x9F, 0x42, 0xBC,
  }
  state := copyToState(block)

  key := []byte{
    0xAC, 0x77, 0x66, 0xF3,
    0x19, 0xFA, 0xDC, 0x21,
    0x28, 0xD1, 0x29, 0x41,
    0x57, 0x5C, 0x00, 0x6A,
  }

  result_candidate_int := addRoundKey(state, key)

  right_result_int := [][]byte{
    []byte{0xEB, 0x59, 0x8B, 0x1B},
    []byte{0x40, 0x2E, 0xA1, 0xC3},
    []byte{0xF2, 0x38, 0x13, 0x42},
    []byte{0x1E, 0x84, 0xE7, 0xD6},
  }

	result_candidate := fmt.Sprintf("%v", result_candidate_int)
	right_result := fmt.Sprintf("%v", right_result_int)
  if right_result != result_candidate {
		t.Errorf("Wrong addRoundKey %v != %v", result_candidate, right_result)
  }
}

func TestKeyExpansionEncrypt(t *testing.T) {
	// Check fips 197, Appendix A, page 27
	key := []byte{
		0x2b, 0x7e, 0x15, 0x16, 0x28, 0xae, 0xd2, 0xa6,
		0xab, 0xf7, 0x15, 0x88, 0x09, 0xcf, 0x4f, 0x3c,
	}

	key_candidate_int := expandKeyEncrypt(make([]uint32, EXPANDED_KEY_SIZE_WORDS), key)
	key_right_int := []uint32{
		0x2b7e1516, 0x28aed2a6, 0xabf71588, 0x9cf4f3c,
		0xa0fafe17, 0x88542cb1, 0x23a33939, 0x2a6c7605,
		0xf2c295f2, 0x7a96b943, 0x5935807a, 0x7359f67f,
		0x3d80477d, 0x4716fe3e, 0x1e237e44, 0x6d7a883b,
		0xef44a541, 0xa8525b7f, 0xb671253b, 0xdb0bad00,
		0xd4d1c6f8, 0x7c839d87, 0xcaf2b8bc, 0x11f915bc,
		0x6d88a37a, 0x110b3efd, 0xdbf98641, 0xca0093fd,
		0x4e54f70e, 0x5f5fc9f3, 0x84a64fb2, 0x4ea6dc4f,
		0xead27321, 0xb58dbad2, 0x312bf560, 0x7f8d292f,
		0xac7766f3, 0x19fadc21, 0x28d12941, 0x575c006e,
		0xd014f9a8, 0xc9ee2589, 0xe13f0cc8, 0xb6630ca6,
	}
	key_candidate := fmt.Sprintf("%v", key_candidate_int)
	key_right := fmt.Sprintf("%v", key_right_int)

	if key_candidate != key_right {
		t.Errorf("Wrong key %v != %v", key_candidate, key_right)
	}
}

func TestKeyExpansionDecrypt(t *testing.T) {
	key := getEmptyBlock()
	key_candidate_int := []uint32(expandKeyDecrypt(key))

	key_right_int := []uint32{3411799988, 300061246, 3478251811, 2383974255, 226657621, 2297109739, 2998141196, 3999567141, 3924997738, 2238290366, 979312103, 1557171241, 2336942831, 1822109652, 3208021081, 1720685006, 2265703722, 3889284411, 3551374221, 3652785559, 4028541365, 1624907793, 880570038, 169131546, 2456059330, 2428799396, 1420185255, 1047056556, 2886511982, 44102758, 3294807811, 1791925771, 3334615909, 2930530568, 3334615909, 2930530568, 1752066669, 1752066669, 1752066669, 1752066669, 0, 0, 0, 0}

	key_candidate := fmt.Sprintf("%v", key_candidate_int)
	key_right := fmt.Sprintf("%v", key_right_int)

	if key_candidate != key_right {
		t.Errorf("Wrong key %v != %v", key_candidate, key_right)
	}
}

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
