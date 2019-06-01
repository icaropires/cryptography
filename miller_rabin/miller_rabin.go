package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

func miller_rabin(n uint64, k uint64) bool {
	if n == 2 || n == 3 {
		return true
	}

	if n%2 == 0 {
		return false
	}

	r := 0
	d := n - 1
	for d%2 == 0 {
		r += 1
		d /= 2
	}
	for j := uint64(0); j < k; j++ {
		key := make([]byte, 32)
		random_bytes := stream_salsa20(8, []byte{4, 2, 4, 4, 1, byte(j), 4, 2}, key)
		a := binary.BigEndian.Uint64(random_bytes)%(n-2) + 2

		x := new(big.Int).Exp(big.NewInt(int64(a)), big.NewInt(int64(d)), big.NewInt(int64(n)))

		if x.Uint64() == 1 || x.Uint64() == n-1 {
			continue
		}
		entrou := false
		for i := 0; i < r-1; i++ {
			x = new(big.Int).Exp(x, big.NewInt(int64(2)), big.NewInt(int64(n)))
			if x.Uint64() == n-1 {
				entrou = true
				break
			}
		}
		if !entrou {
			return false
		}
	}
	return true
}

func main() {
	var ent int
	var k float64
	fmt.Printf("Número: ")
	fmt.Scanf("%d", &ent)
	fmt.Printf("Iterações: ")
	fmt.Scanf("%f", &k)
	if miller_rabin(uint64(ent), 10) {
		fmt.Println("O número é primo com probabilidade: ", 1.0-math.Pow(0.25, k))
	} else {
		fmt.Println("O número é composto")
	}
}
