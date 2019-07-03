// Question 3

package main

import (
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 8 {
		fmt.Println("Uso incorreto! Exemplo de uso: ./bin [a] [b] [p] [Gx] [Gy] [Px] [Py]")
		return
	}

	aStr, bStr, pStr, gxStr, gyStr, pxStr, pyStr := os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6], os.Args[7]

	a, _ := new(big.Int).SetString(aStr, 10)
	b, _ := new(big.Int).SetString(bStr, 10)
	p, _ := new(big.Int).SetString(pStr, 10)
	gx, _ := new(big.Int).SetString(gxStr, 10)
	gy, _ := new(big.Int).SetString(gyStr, 10)
	px, _ := new(big.Int).SetString(pxStr, 10)
	py, _ := new(big.Int).SetString(pyStr, 10)

	biggest := getBiggestOrder(a, b, p)

	rand.Seed(time.Now().UnixNano())
	privateKey := rand.Intn(biggest)

	if privateKey == 0 {
		privateKey++
	}

	g := Point{gx, gy}

	publicKey := g.Mul(privateKey, a, p)

	k := rand.Intn(biggest)
	if k == 0 {
		k++
	}

	pPoint := Point{px, py}
	aux := publicKey.Mul(k, a, p)

	fmt.Println("Initial Plain Point: ", pPoint)
	fmt.Println("G = ", g)
	fmt.Println("Public key: ", publicKey)

	c1 := g.Mul(k, a, p)
	c2 := pPoint.Add(aux, a, p)

	fmt.Println("Cipher point: ", c1, "e", c2)

	aux = c1.Mul(privateKey, a, p)
	aux.y = new(big.Int).Neg(aux.y)

	plain := c2.Add(aux, a, p)
	fmt.Println("Plain Point: ", plain)
}
