package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/big"
	"math/rand"
	"os"
	"time"
)

const (
	M = 123
	r = 30
	p = 4177
)

func readFile(filepath string) []byte {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic("Não foi possível ler do arquivo")
	}
	file = bytes.Trim(file, "\n")

	//for i := 0; i < len(file)%16; i++ {
	//	file = append(file, ' ')
	//}

	fmt.Println(file)
	return file
}

func mapping(message byte) *Point {
	j := int64(0)
	for {
		y := big.NewInt(int64(message)*r + j)

		x := new(big.Int).Set(y)

		rr := y.ModSqrt(new(big.Int).Sub(new(big.Int).Exp(y, big.NewInt(3), nil), big.NewInt(4)), big.NewInt(4177))

		if rr != nil {
			return &Point{x, y}
		}
		j += 1
	}
}

func main() {
	if len(os.Args) < 9 {
		fmt.Println("Uso incorreto! Exemplo de uso: ./bin [a] [b] [p] [Gx] [Gy] [Px] [Py] [filename]")
		return
	}

	aStr, bStr, pStr, gxStr, gyStr, filename := os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6]

	curve := &Curve{}
	curve.a, _ = new(big.Int).SetString(aStr, 10)
	curve.b, _ = new(big.Int).SetString(bStr, 10)

	gx, _ := new(big.Int).SetString(gxStr, 10)
	gy, _ := new(big.Int).SetString(gyStr, 10)

	biggest := getBiggestOrder(curve)

	message := readFile(filename)

	for _, char := range message {
		fmt.Println(mapping(char))
	}

	r := &Point{}
	s := &Point{}
	k := 0

	for s.IsAtInfinity() {
		for r.IsAtInfinity() || k == 0 {
			rand.Seed(time.Now().UnixNano())
			k := rand.Intn(biggest)
			pPoint := g.Mul(k, curve)
			r := new(big.Int).Mod(pPoint.x, biggest)
		}

		kBig := big.NewInt(k)
		t := new(big.Int).ModInv(kBig, biggest)

		e := hash(filename)

	}

}
