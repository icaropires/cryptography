package main

import (
	"crypto/aes"
	"fmt"
	"testing"
)

func TestSubWordByte(t *testing.T) {
	var sbox_value_right byte = 0x5d
	sbox_value_candidate := subWordByte(0x8d)

	if sbox_value_right != sbox_value_candidate {
		t.Errorf("Wrong subWordByte %v != %v", sbox_value_right, sbox_value_candidate)
	}
}

func TestSubWord(t *testing.T) {
	subword_right := uint32(0x5da515d2)
	subword_candidate := subWord(0x8d292f7f)

	if subword_right != subword_candidate {
		t.Errorf("Wrong subWord %v != %v", subword_right, subword_candidate)
	}
}

func TestKeyExpansionEncrypt(t *testing.T) {
	key := getEmptyBlock()
	key_candidate_int := []uint32(expandKeyEncrypt(key))

	key_right_int := []uint32{0, 0, 0, 0, 1667457890, 1667457890, 1667457890, 1667457890, 3382220955, 2868640761, 3382220955, 2868640761, 1345623952, 4207897705, 861402354, 2578190091, 2077886190, 2165664391, 2990710389, 737055102, 2284531327, 155075832, 3145521805, 2425506803, 2236309996, 2356487444, 923402137, 2811999338, 2266461473, 190992437, 1013690284, 2616204230, 855898382, 945924411, 67765911, 2683968849, 3805861041, 3669589386, 3736304413, 1095329356, 3411799988, 300061246, 3478251811, 2383974255}

	key_candidate := fmt.Sprintf("%v", key_candidate_int)
	key_right := fmt.Sprintf("%v", key_right_int)

	if key_candidate != key_right {
		t.Errorf("Wrong key %v != %v", key_candidate, key_right)
	}
}

func TestKeyExpansionDecrypt(t *testing.T) {
	key := getEmptyBlock()
	key_candidate_int := []uint32(expandKeyEncrypt(key))

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
