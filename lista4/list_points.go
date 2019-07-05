// Questão 2

package main

import (
	"fmt"
	"math/big"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Uso incorreto! Exemplo de uso: ./bin [a] [b] [p]")
		return
	}

	aStr, bStr, pStr := os.Args[1], os.Args[2], os.Args[3]

	curve := &Curve{}
	curve.a, _ = new(big.Int).SetString(aStr, 10)
	curve.b, _ = new(big.Int).SetString(bStr, 10)
	curve.p, _ = new(big.Int).SetString(pStr, 10)

	points, orders := getAllPoints(curve)
	biggest := getBiggestOrder(curve)

	fmt.Printf("Ponto(s) com maior(es) ordem(ns): ")
	for i, e := range orders {
		if e == int(biggest) {
			if i == len(points)-1 {
				fmt.Printf("%s", points[i])
			} else {
				fmt.Printf("%s ", points[i])
			}
		}
	}
	fmt.Println()

	fmt.Println("Total de pontos: ", len(points), "+ ", "{\u221E}")
	for i, point := range points {
		fmt.Printf("Point #%2d: %s, order = %2d\n", i+1, point, orders[i])
	}
}
