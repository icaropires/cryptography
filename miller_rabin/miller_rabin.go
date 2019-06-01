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
		fmt.Println("dddddddddddd", d)
	}
	fmt.Println(r, d)
	for j := uint64(0); j < k; j++ {
		key := make([]byte, 32)
		random_bytes := stream_salsa20(8, []byte{4, 2, 4, 4, 1, byte(j), 4, 2}, key)
		a := binary.BigEndian.Uint64(random_bytes)%(n-2) + 2
		fmt.Println("a", a)
		fmt.Println("d", d)

		x := new(big.Int).Exp(big.NewInt(int64(a)), big.NewInt(int64(d)), nil)
		fmt.Println("x", x)

		if x == big.NewInt(1) || x == big.NewInt(int64(n)-1) {
			continue
		}
		flag := false
		for i := 0; i < r-1; i++ {
			x = new(big.Int).Exp(x, big.NewInt(int64(2)), big.NewInt(int64(n)))
			fmt.Println("===============>", x)
			if x == big.NewInt(int64(n)-1) {
				flag = true
			}
		}
		if !flag {
			return false
		}
	}
	return true
}

func main() {
	fmt.Println(miller_rabin(252601, 1000))
	fmt.Println("Probabilidade: ", math.Pow(0.25, 100))

}
