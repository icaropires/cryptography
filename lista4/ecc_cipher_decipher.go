// Question 3

package main

import (
	"fmt"
	"math/big"
	"os"
)

func main() {
	if len(os.Args) < 8 {
		fmt.Println("Uso incorreto! Exemplo de uso: ./bin [a] [b] [p] [Gx] [Gy] [Px] [Py]")
		return
	}

	aStr, bStr, pStr, gxStr, gyStr, pxStr, pyStr := os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6], os.Args[7]

	curve := &Curve{}
	curve.a, _ = new(big.Int).SetString(aStr, 10)
	curve.b, _ = new(big.Int).SetString(bStr, 10)
	curve.p, _ = new(big.Int).SetString(pStr, 10)

	g := &Point{}
	g.x, _ = new(big.Int).SetString(gxStr, 10)
	g.y, _ = new(big.Int).SetString(gyStr, 10)

	pPoint := &Point{}
	pPoint.x, _ = new(big.Int).SetString(pxStr, 10)
	pPoint.y, _ = new(big.Int).SetString(pyStr, 10)

	privateKey, publicKey := GenKeys(g, curve)

	fmt.Println("Initial Plain Point: ", pPoint)
	fmt.Println("G = ", g)
	fmt.Println("Public key: ", publicKey)

	c1, c2 := Cipher(pPoint, publicKey, g, curve)
	fmt.Println("Cipher point: ", c1, "e", c2)

	plain := Decipher(c1, c2, privateKey, curve)
	fmt.Println("Plain Point: ", plain)
}
