package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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
	if len(os.Args) < 6 {
		fmt.Println("Uso incorreto! Exemplo de uso: ./bin [a] [b] [p] [Gx] [Gy] [filename]")
		return
	}

	aStr, bStr, pStr, gxStr, gyStr, filename := os.Args[1], os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6]

	curve := &Curve{}
	curve.a, _ = new(big.Int).SetString(aStr, 10)
	curve.b, _ = new(big.Int).SetString(bStr, 10)
	curve.p, _ = new(big.Int).SetString(pStr, 10)

	g := &Point{}
	g.x, _ = new(big.Int).SetString(gxStr, 10)
	g.y, _ = new(big.Int).SetString(gyStr, 10)

	biggest := getBiggestOrder(curve)

	n := big.NewInt(int64(biggest))

	message := readFile(filename)

	for _, char := range message {
		mapping(char)
	}

	r := big.NewInt(0)
	s := big.NewInt(0)
	privateKey, publicKey := GenKeys(g, curve)
	k := 0

	for s.Uint64() == 0 {
		for r.Uint64() == 0 || k == 0 {
			rand.Seed(time.Now().UnixNano())
			k = rand.Intn(int(n.Int64()) - 1)
			pPoint := g.Mul(k, curve)
			r = new(big.Int).Mod(pPoint.x, n)
		}

		kBig := big.NewInt(int64(k))

		e, _ := new(big.Int).SetString(hash(filename), 16)

		z := big.NewInt(e.Int64() >> uint(e.BitLen()-n.BitLen()))

		numerator := new(big.Int).Add(z, new(big.Int).Mul(big.NewInt(int64(privateKey)), r))
		denominator := new(big.Int).ModInverse(kBig, n)

		if denominator != nil {
			denominator.Mod(denominator, n)
			numerator.Mod(numerator, n)
			s = new(big.Int).Mul(numerator, denominator)
			s.Mod(s, n)
		} else {
			s = big.NewInt(0)
			k = 0
		}
	}

	fmt.Println(r, s, n, publicKey)
	InsertStringToFile(filename, "signature: "+r.String()+" "+s.String()+"\n", 0)

}
