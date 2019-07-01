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

	a, _ := new(big.Int).SetString(aStr, 10)
	b, _ := new(big.Int).SetString(bStr, 10)
	p, _ := new(big.Int).SetString(pStr, 10)

	points, orders := getAllPoints(a, b, p)

	biggest := getBiggestOrder(orders)

	fmt.Printf("Ponto(s) com maior(es) ordem(ns): ")
	for i, e := range orders {
		if e == biggest {
			if i == len(points)-1 {
				fmt.Printf("%s", points[i])
			} else {
				fmt.Printf("%s ", points[i])
			}
		}
	}
	fmt.Println()

	fmt.Println("Total de pontos: ", len(points), "+ infinito")
	for i, point := range points {
		fmt.Printf("Point #%d: %s, order = %d\n", i+1, point, orders[i])
	}
}
