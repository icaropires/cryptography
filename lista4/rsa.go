package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
)

func phi(p, q *big.Int) *big.Int {
	one := big.NewInt(1)

	pAux := new(big.Int).Sub(p, one)
	qAux := new(big.Int).Sub(q, one)

	return new(big.Int).Mul(pAux, qAux)
}

func isValidKey(key, phiN *big.Int) bool {
	gcd := new(big.Int).GCD(nil, nil, phiN, key)

	return gcd.Uint64() == 1 && (phiN.Cmp(key) == 1)
}

func genPrivateKey(key, phiN *big.Int) *big.Int {
	return new(big.Int).ModInverse(key, phiN)
}

func Encrypt(n, key *big.Int, text string) string {
	encryptAux := func(keyInternal *big.Int, chunk byte) string {
		chunkInt := new(big.Int)
		chunkInt.SetUint64(uint64(chunk))

		return new(big.Int).Exp(chunkInt, keyInternal, n).String()
	}

	data := ""
	separator := "\n"
	for i := 0; i < len(text); i++ {
		data += (encryptAux(key, text[i]) + separator)
	}

	return data
}

func Decrypt(n, key *big.Int, text string) string {
	decryptAux := func(keyInternal, chunk *big.Int) uint64 {
		return new(big.Int).Exp(chunk, keyInternal, n).Uint64()
	}

	data := ""
	separator := "\n"
	chunks := strings.Split(text, separator)
	for _, chunk := range chunks {
		chunkInt := new(big.Int)
		chunkInt.SetString(chunk, 10)
		data += string(decryptAux(key, chunkInt))
	}

	return data
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Incorrect Usage. Usage: ./bin [p] [q] [key]")
	}

	pStr, qStr, keyStr := os.Args[1], os.Args[2], os.Args[3]

	p := new(big.Int)
	q := new(big.Int)
	key := new(big.Int)

	p.SetString(pStr, 10)
	q.SetString(qStr, 10)
	key.SetString(keyStr, 10)

	phiN := phi(p, q)
	for !isValidKey(key, phiN) {
		fmt.Println("Chave inválida! Insira uma válida!")
		fmt.Scanln(&key)
	}

	var operation string
	for operation != "d" && operation != "e" {
		fmt.Println("Insira a operação. 'd' ou 'D' para decifrar e 'e' ou 'E' para cifrar")
		fmt.Scanln(&operation)

		operation = strings.ToLower(operation)
	}

	fmt.Println("Qual arquivo de entrada?")
	var inFile string
	fmt.Scanln(&inFile)

	fmt.Println("Onde deseja salvar o resultado?")
	var outFile string
	fmt.Scanln(&outFile)

	textByte, _ := ioutil.ReadFile(inFile)
	text := string(textByte)

	n := new(big.Int).Mul(p, q)
	var data string

	if operation == "d" {
		privKey := genPrivateKey(key, phiN)
		data = Decrypt(n, privKey, text)
	} else {
		data = Encrypt(n, key, text)
	}

	ioutil.WriteFile(outFile, []byte(data), 0644)
}
