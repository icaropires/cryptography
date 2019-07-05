package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
)

func readFile(filepath string) (string, string) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic("Não foi possível ler do arquivo")
	}
	file = bytes.Trim(file, "\n")
	brPos := strings.IndexByte(string(file), '\n')

	if brPos == -1 {
		return "", ""
	}

	line := string(file[:brPos])
	var r, s string
	n, err := fmt.Sscanf(line, "signature: %s %s\n", &r, &s)
	if n == 0 || err != nil {
		return "", ""
	}

	return r, s
}

func main() {
	if len(os.Args) < 8 {
		fmt.Println("Uso incorreto! Exemplo de uso: go run ecc.go sha_256.go veirfy.go [a] [b] [p] [Gx] [Gy] [publicKeyX] [publicKeyY] [filename]")
		return
	}

	aStr, bStr, pStr, gxStr, gyStr, pkXStr, pkYStr, filename := os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6], os.Args[7], os.Args[8]

	rStr, sStr := readFile(filename)
	if rStr == "" || sStr == "" {
		fmt.Println("Assinatura digital inválida: format invalid on file")
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

	biggest := getBiggestOrder(curve)
	n := big.NewInt(int64(biggest))

	w := new(big.Int).ModInverse(s, n)
	if w == nil {
		fmt.Printf("Assinatura digital inválida: s='%v' e n='%v' não coprimos\n", s, n)
		return
	}

	e, _ := new(big.Int).SetString(hash(filename), 16)
	z := e.Rsh(e, uint(e.BitLen()-n.BitLen())) // Fips 180

	u1 := new(big.Int).Mul(z, w)
	u1.Mod(u1, n)

	u2 := new(big.Int).Mul(r, w)
	u2.Mod(u2, n)

	aux1 := g.Mul(int(u1.Int64()), curve)
	aux2 := publicKey.Mul(int(u2.Int64()), curve)

	pPoint := aux1.Add(aux2, curve)
	if pPoint.IsAtInfinity() {
		fmt.Println("Assinatura digital inválida: pPoint infinity")
	}
	pPoint.x.Mod(pPoint.x, n)

	if pPoint.x.Cmp(r) == 0 {
		fmt.Println("Assinatura digital válida")
	} else {
		fmt.Printf("Assinatura digital inválida: Px='%v' != r='%v'\n", pPoint.x, r)
	}
}
