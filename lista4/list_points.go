// Questão 2

package main

import (
	"fmt"
	"math"
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

		if divisor == nil {
			fmt.Println("Zero division!")
			return Point{}
		}

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

func ListCurvePoints(a, b, p *big.Int) ([]Point, []int) {
	points := make([]Point, 0)
	orders := make([]int, 0)

	for xInt := int64(0); xInt < p.Int64(); xInt++ {
		x := big.NewInt(xInt)

		ySquare := curve(a, b, x)
		y := new(big.Int).ModSqrt(ySquare, p)

		if y != nil {
			p_point := Point{x, y}
			order := p_point.getOrder(p, a)

			points = append(points, p_point)
			orders = append(orders, order)
			if y.Uint64() != 0 {
				p_point = Point{x, new(big.Int).Sub(p, y)}
				order = p_point.getOrder(p, a)

				points = append(points, p_point)
				orders = append(orders, order)
			}
		}
	}

	return points, orders
}

func getPoint(x, y int64) Point {
	return Point{
		big.NewInt(x),
		big.NewInt(y),
	}
}

func (p_point *Point) getOrder(p *big.Int, a *big.Int) int {
	if p_point.x.Cmp(p_point.x) == 0 && p_point.y.Cmp(p_point.y) == 0 && p_point.y.Cmp(big.NewInt(0)) == 0 { // Zero division!!!
		return 2
	}
	r := p_point.Add(p_point, p, a)

	counter := 2

	for {
		counter += 1

		if r.x.Cmp(p_point.x) == 0 { // Zero division!!!
			break
		}

		if r.x.Cmp(p_point.x) == 0 && r.y.Cmp(p_point.y) == 0 && r.y.Cmp(big.NewInt(0)) == 0 { // Zero division!!!
			break
		}

		r = r.Add(p_point, p, a)
	}

	return counter
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

	points, orders := ListCurvePoints(a, b, p)

	// Get biggest
	biggest := float64(orders[0])
	for _, e := range orders {
		biggest = math.Max(float64(biggest), float64(e))
	}

	fmt.Printf("Ponto(s) com maior(es) ordem(ns): ")
	for i, e := range orders {
		if float64(e) == biggest {
			if i == len(points)-1 {
				fmt.Printf("%s", points[i])
			} else {
				fmt.Printf("%s ", points[i])
			}
		}
	}
	fmt.Println()

	fmt.Println("Total de pontos: ", len(points))
	for i, point := range points {
		fmt.Printf("Point #%d: %s, order = %d\n", i+1, point, orders[i])
	}
}
