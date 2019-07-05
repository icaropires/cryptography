package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
)

func readFile(filepath string) (string, string, string) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic("Não foi possível ler do arquivo")
	}
	file = bytes.Trim(file, "\n")
	brPos := strings.IndexByte(string(file), '\n')

	if brPos == -1 {
		return "", "", ""
	}

	line := string(file[:brPos])
	var r, s, z string
	n, err := fmt.Sscanf(line, "signature: %s %s %s\n", &r, &s, &z)
	if n == 0 || err != nil {
		return "", "", ""
	}

	return r, s, z
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Uso incorreto! Exemplo de uso: go run ecc.go sha_256.go verify.go [filename] [publicKeyX] [publicKeyY]")
		return
	}

	filename, pkXStr, pkYStr := os.Args[1], os.Args[2], os.Args[3]

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Assinatura digital inválida. Usou a chave pública correta?:", r)
		}
	}()

	rStr, sStr, zStr := readFile(filename)
	if rStr == "" || sStr == "" || zStr == "" {
		fmt.Println("Assinatura digital inválida: format invalid on file")
		return
	}

	r, _ := new(big.Int).SetString(rStr, 10)
	s, _ := new(big.Int).SetString(sStr, 10)
	z, _ := new(big.Int).SetString(zStr, 10)

	fmt.Println("z =", z)

	// Using curve secp256k1
	curve := &Curve{}
	curve.a, _ = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000000", 16)
	curve.b, _ = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000007", 16)
	curve.p, _ = new(big.Int).SetString("fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f", 16)
	curve.n, _ = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)

	g := &Point{}
	g.x, _ = new(big.Int).SetString("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", 16)
	g.y, _ = new(big.Int).SetString("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8", 16)

	publicKey := &Point{}
	publicKey.x, _ = new(big.Int).SetString(pkXStr, 10)
	publicKey.y, _ = new(big.Int).SetString(pkYStr, 10)

	w := new(big.Int).ModInverse(s, curve.n)
	if w == nil {
		fmt.Printf("Assinatura digital inválida: s='%v' e n='%v' não coprimos\n", s, curve.n)
		return
	}

	u1 := new(big.Int).Mul(z, w)
	u1.Mod(u1, curve.n)

	u2 := new(big.Int).Mul(r, w)
	u2.Mod(u2, curve.n)

	aux1 := g.Mul(u1, curve)
	aux2 := publicKey.Mul(u2, curve)

	pPoint := aux1.Add(aux2, curve)
	if pPoint.IsAtInfinity() {
		fmt.Println("Assinatura digital inválida: pPoint infinity")
		return
	}
	pPoint.x.Mod(pPoint.x, curve.n)

	if pPoint.x.Cmp(r) == 0 {
		fmt.Println("Assinatura digital válida")
	} else {
		fmt.Printf("Assinatura digital inválida: Px='%v' != r='%v'\n", pPoint.x, r)
	}
}
