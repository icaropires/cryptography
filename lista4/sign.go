package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
)

func File2lines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LinesFromReader(f)
}

func LinesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

/**
 * Insert sting to n-th line of file.
 * If you want to insert a line, append newline '\n' to the end of the string.
 */
func InsertStringToFile(path, str string, index int) error {
	lines, err := File2lines(path)
	if err != nil {
		return err
	}

	fileContent := ""
	for i, line := range lines {
		if i == index {
			fileContent += str
		}
		fileContent += line
		fileContent += "\n"
	}

	return ioutil.WriteFile(path, []byte(fileContent), 0644)
}

func readFile(filepath string) []byte {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic("Não foi possível ler do arquivo")
	}
	file = bytes.Trim(file, "\n")

	return file
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso incorreto! Exemplo de uso: go run ecc.go sha_256.go sign.go [filename]")
		return
	}

	filename := os.Args[1]

	// Using curve secp256k1
	curve := &Curve{}
	curve.a, _ = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000000", 16)
	curve.b, _ = new(big.Int).SetString("0000000000000000000000000000000000000000000000000000000000000007", 16)
	curve.p, _ = new(big.Int).SetString("fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f", 16)
	curve.n, _ = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)

	g := &Point{}
	g.x, _ = new(big.Int).SetString("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", 16)
	g.y, _ = new(big.Int).SetString("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8", 16)

	r := big.NewInt(0)
	s := big.NewInt(0)
	k := big.NewInt(0)

	privateKey, publicKey := GenKeys(g, curve)
	z, _ := new(big.Int).SetString(hash(filename), 16)
	//z := e.Rsh(e, uint(e.BitLen()-curve.n.BitLen())) // FIPS 180

	for s.Uint64() == 0 {
		for r.Uint64() == 0 || k.Uint64() == 0 {
			k = getRandom(curve.n)

			pPoint := g.Mul(k, curve)
			r = pPoint.x.Mod(pPoint.x, curve.n)
		}

		numerator := new(big.Int).Mul(privateKey, r)
		numerator.Add(numerator, z)

		denominator := new(big.Int).ModInverse(k, curve.n)

		if denominator != nil {
			s = new(big.Int).Mul(numerator, denominator)
			s.Mod(s, curve.n)
		} else {
			s = big.NewInt(0)
			k = big.NewInt(0)
		}
	}

	InsertStringToFile(filename, "signature: "+r.String()+" "+s.String()+"\n", 0)
	fmt.Printf("Sua chave pública, para ser usada na verificação, é (X Y):\n%v %v\n", publicKey.x, publicKey.y)
}
