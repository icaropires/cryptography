package main

import "fmt"

const ROUNDS = 20

func rotl(a, b byte) byte {
	return (a << b) | (a >> (8 - b))
}

func quarter_round(a, b, c, d *byte) {
	*b ^= rotl(*a+*d, 7)
	*c ^= rotl(*b+*a, 9)
	*d ^= rotl(*c+*b, 13)
	*a ^= rotl(*d+*c, 18)
}

func salsa20_block(out, in []byte) {
	x := make([]byte, 16)

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
		quarter_round(&x[10], &x[11], &x[8], &x[0])
		quarter_round(&x[15], &x[12], &x[13], &x[14])
	}

	for i := 0; i < 16; i++ {
		out[i] = x[i] + in[i]
	}
}

func main() {
	out := make([]byte, 16)
	in := []byte("abobora com pato")

	salsa20_block(out, in)

	fmt.Println(out)
}
