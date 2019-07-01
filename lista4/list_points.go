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

func (p_point Point) String() string {
	return fmt.Sprintf("(%s, %s)", p_point.x.String(), p_point.y.String())
}

func (p_point *Point) Add(q *Point, p *big.Int, a *big.Int) Point {
	if p_point.x == q.x && p_point.y == new(big.Int).Neg(q.y) {
		fmt.Println("Isso aqui vai dar infinito")
		return Point{}
	}

	var lambda *big.Int
	if p_point.x.Cmp(q.x) == 0 && p_point.y.Cmp(q.y) == 0 {
		dividendAux := new(big.Int).Exp(p_point.x, big.NewInt(2), nil)
		dividendAux = new(big.Int).Mul(big.NewInt(3), dividendAux)

		dividend := new(big.Int).Add(dividendAux, a)
		dividend = new(big.Int).Mod(dividend, p)

		divisor := new(big.Int).Mul(big.NewInt(2), p_point.y)
		divisor = new(big.Int).ModInverse(divisor, p)

		lambdaAux := new(big.Int).Mul(dividend, divisor)
		lambda = new(big.Int).Mod(lambdaAux, p)
	} else {
		deltaY := new(big.Int).Sub(q.y, p_point.y)
		deltaX := new(big.Int).Sub(q.x, p_point.x)

		deltaY = new(big.Int).Mod(deltaY, p)
		deltaX = new(big.Int).ModInverse(deltaX, p)

		if deltaX == nil {
			fmt.Println("Zero division!")
			return Point{}
		}

		lambdaAux := new(big.Int).Mul(deltaY, deltaX)
		lambda = new(big.Int).Mod(lambdaAux, p)
	}

	xAux := new(big.Int).Exp(lambda, big.NewInt(2), nil)
	xAux = new(big.Int).Sub(xAux, p_point.x)
	xAux = new(big.Int).Sub(xAux, q.x)
	x := new(big.Int).Mod(xAux, p)

	yAux := new(big.Int).Mul(lambda, new(big.Int).Sub(p_point.x, x))
	yAux = new(big.Int).Sub(yAux, p_point.y)
	y := new(big.Int).Mod(yAux, p)

	return Point{x, y}
}

// Returns y^2
func curve(a, b, x *big.Int) *big.Int {
	aux := new(big.Int).Exp(x, big.NewInt(3), nil)

	aux.Add(aux, new(big.Int).Mul(a, x))
	aux.Add(aux, b)

	return aux
}

func ListCurvePoints(a, b, p *big.Int) []Point {
	points := make([]Point, 0)
	for xInt := int64(0); xInt < p.Int64(); xInt++ {
		x := big.NewInt(xInt)

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

func getPoint(x, y int64) Point {
	return Point{
		big.NewInt(x),
		big.NewInt(y),
	}
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
		fmt.Printf("Point #%d: %s\n", i+1, point)
	}

	p_point := getPoint(3, 10)
	q := getPoint(9, 7)

	fmt.Printf("%s + %s = ", p_point, q)
	r := p_point.Add(&q, p, a)
	fmt.Println(r)
}
