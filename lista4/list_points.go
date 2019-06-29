// Questão 2

package main

import (
	"fmt"
	"math/big"
	"os"
)

type Point struct {
	x *big.Int
	y *big.Int
}

// Returns y^2
func curve(a, b, x *big.Int) *big.Int {
	aux := new(big.Int).Exp(x, new(big.Int).SetUint64(3), nil)

	aux.Add(aux, new(big.Int).Mul(a, x))
	aux.Add(aux, b)

	return aux
}

func ListCurvePoints(a, b, p *big.Int) []Point {
	points := make([]Point, 0)
	for xInt := int64(0); xInt < p.Int64(); xInt++ {
		x := new(big.Int).SetInt64(xInt)

		ySquare := curve(a, b, x)
		y := new(big.Int).ModSqrt(ySquare, p)

		if y != nil {
			points = append(points, Point{x, y})
			if y.Uint64() != 0 {
				points = append(points, Point{x, new(big.Int).Sub(p, y)})
			}
		}
	}

	return points
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Uso incorreto! Exemplo de uso: ./bin [a] [b] [p]")
		return
	}

	aStr, bStr, pStr := os.Args[1], os.Args[2], os.Args[3]

	a, _ := new(big.Int).SetString(aStr, 10)
	b, _ := new(big.Int).SetString(bStr, 10)
	p, _ := new(big.Int).SetString(pStr, 10)

	points := ListCurvePoints(a, b, p)

	fmt.Println("Total de pontos: ", len(points))
	for i, point := range points {
		fmt.Printf("Point #%d: (%s, %s)\n", i+1, point.x.String(), point.y.String())
	}
}
