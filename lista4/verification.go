package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
)

const (
	M = 123
	r = 30
	p = 4177
)

func readFile(filepath string) (string, string) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic("Não foi possível ler do arquivo")
	}
	file = bytes.Trim(file, "\n")
	brPos := strings.IndexByte(string(file), '\n')

	if brPos == -1 {
		return nil, nil
	}

	line := string(file[:brPos])
	fmt.Println("Line", line)
	var r, s string
	n, err := fmt.Sscanf(line, "signature: %s %s\n", &r, &s)
	if n == 0 || err != nil {
		return nil, nil
	}
	fmt.Println("rs", r, s)

	return r, s
}

func main() {

	aStr, bStr, pStr, gxStr, gyStr, pkXStr, pkYStr, filename := os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6], os.Args[7], os.Args[8]

	rStr, sStr := readFile(filename)

	if rStr == nil || sStr == nil {
		fmt.Println("Assinatura digital inválida")
		return
	}

	r, _ := new(big.Int).SetString(rStr, 10)
	s, _ := new(big.Int).SetString(sStr, 10)

	curve := &Curve{}
	curve.a, _ = new(big.Int).SetString(aStr, 10)
	curve.b, _ = new(big.Int).SetString(bStr, 10)
	curve.p, _ = new(big.Int).SetString(pStr, 10)

	g := &Point{}
	g.x, _ = new(big.Int).SetString(gxStr, 10)
	g.y, _ = new(big.Int).SetString(gyStr, 10)

	publicKey := &Point{}
	publicKey.x, _ = new(big.Int).SetString(pkXStr, 10)
	publicKey.y, _ = new(big.Int).SetString(pkYStr, 10)

	fmt.Println("----------->", r, s)

	biggest := getBiggestOrder(curve)

	n := big.NewInt(int64(biggest))

	e, _ := new(big.Int).SetString(hash(filename), 16)
	numerator := new(big.Int).Mod(e, n)

	denominator := new(big.Int).ModInverse(s, n)

	if denominator == nil {
		fmt.Println("Assinatura digital inválida")
		return
	}

	u1 := new(big.Int).Mul(numerator, denominator)

	numerator = new(big.Int).Mod(r, n)
	u2 := new(big.Int).Mul(numerator, denominator)

	aux := g.Mul(int(u1.Int64()), curve)

	aux2 := publicKey.Mul(int(u2.Int64()), curve)

	pPoint := aux.Add(aux2, curve)

	pPoint.x.Mod(pPoint.x, n)

	if pPoint.x.Cmp(r) == 0 {
		fmt.Println("Assinatura digital válida")
	} else {
		fmt.Println("Assinatura digital inválida")
	}
	fmt.Println(pPoint.x, r)

}
