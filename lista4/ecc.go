package main

import (
	"fmt"
	"math"
	"math/big"
)

type Point struct {
	x *big.Int
	y *big.Int
}

func (p_point Point) String() string {
	return fmt.Sprintf("(%s, %s)", p_point.x.String(), p_point.y.String())
}

func (p_point *Point) Add(q *Point, a *big.Int, p *big.Int) Point {
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
			fmt.Println("Zero division! Y is zero")
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

func (p_point *Point) getOrder(p *big.Int, a *big.Int) int {
	if p_point.x.Cmp(p_point.x) == 0 && p_point.y.Cmp(p_point.y) == 0 && p_point.y.Cmp(big.NewInt(0)) == 0 { // Zero division!!!
		return 2
	}
	r := p_point.Add(p_point, a, p)

	counter := 2

	for {
		counter += 1

		if r.x.Cmp(p_point.x) == 0 { // Zero division!!!
			break
		}

		if r.x.Cmp(p_point.x) == 0 && r.y.Cmp(p_point.y) == 0 && r.y.Cmp(big.NewInt(0)) == 0 { // Zero division!!!
			break
		}

		r = r.Add(p_point, a, p)
	}

	return counter
}

func (p_point *Point) Mul(n int, a *big.Int, p *big.Int) Point {
	if n == 1 {
		return *p_point
	}

	if p_point.x.Cmp(p_point.x) == 0 && p_point.y.Cmp(p_point.y) == 0 && p_point.y.Cmp(big.NewInt(0)) == 0 { // Zero division!!!
		//fmt.Println("Not Valid multiplication! Y = 0")
		return *p_point
	}
	r := p_point.Add(p_point, a, p)

	for i := 2; i < n; i++ {
		if r.x.Cmp(p_point.x) == 0 { // Zero division!!!
			//fmt.Println("Not Valid multiplication! X = X")
			break
		}

		if r.x.Cmp(p_point.x) == 0 && r.y.Cmp(p_point.y) == 0 && r.y.Cmp(big.NewInt(0)) == 0 { // Zero division!!!
			//fmt.Println("Not Valid multiplication! Next is infinity")
			break
		}

		r = r.Add(p_point, a, p)
	}

	return r
}

// Returns y^2
func curve(a, b, x *big.Int) *big.Int {
	aux := new(big.Int).Exp(x, big.NewInt(3), nil)

	aux.Add(aux, new(big.Int).Mul(a, x))
	aux.Add(aux, b)

	return aux
}

func getPoint(x, y int64) Point {
	return Point{
		big.NewInt(x),
		big.NewInt(y),
	}
}

func getBiggestOrder(orders []int) int {
	biggest := float64(orders[0])
	for _, e := range orders {
		biggest = math.Max(float64(biggest), float64(e))
	}

	return int(biggest)
}

// Get all points of a curve with its order
func getAllPoints(a, b, p *big.Int) ([]Point, []int) {
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
