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

func (pPoint Point) String() string {
	if pPoint.IsAtInfinity() {
		return fmt.Sprint("(\u221E, \u221E)")
	}

	return fmt.Sprintf("(%s, %s)", pPoint.x.String(), pPoint.y.String())
}

// Returns true if a pPoint is equal q
func (pPoint *Point) IsEqual(q *Point) bool {
	return pPoint.x.Cmp(q.x) == 0 && pPoint.y.Cmp(q.y) == 0
}

// Returns the negative of a pointer pPoint
func (pPoint Point) Neg() Point {
	if !pPoint.IsAtInfinity() {
		pPoint = Point{new(big.Int).Set(pPoint.x), new(big.Int).Set(pPoint.y)}
	}

	return Point{new(big.Int).Set(pPoint.x), new(big.Int).Neg(pPoint.y)}
}

// Returns true if a pPoint is a point at the infinity
func (pPoint *Point) IsAtInfinity() bool {
	return pPoint.x == nil && pPoint.y == nil
}

// Add a pPoint pPoint to a point q
func (pPoint *Point) Add(q *Point, a *big.Int, b *big.Int, p *big.Int) *Point {
	if !pPoint.IsOnCurve(a, b, p) {
		msg := fmt.Sprintf("Point '%v' not in curve: a = %s, b = %s, p = %s", pPoint, a, b, p.String())
		panic(msg)
	}

	if !q.IsOnCurve(a, b, p) {
		msg := fmt.Sprintf("Point '%v' not in curve: a = %s, b = %s, p = %s", p, a.String(), b.String(), p.String())
		panic(msg)
	}

	if pPoint.IsAtInfinity() && q.IsAtInfinity() {
		return &Point{}
	}

	if pPoint.IsAtInfinity() {
		return q
	}

	if q.IsAtInfinity() {
		return pPoint
	}

	qNeg := q.Neg()
	if pPoint.IsEqual(&qNeg) {
		return &Point{}
	}

	lambda := new(big.Int)
	if pPoint.IsEqual(q) {
		dividend := new(big.Int).Mul(pPoint.x, pPoint.x)
		dividend.Mul(dividend, big.NewInt(3))
		dividend.Add(dividend, a)
		dividend.Mod(dividend, p)

		divisor := new(big.Int).Mul(big.NewInt(2), pPoint.y)
		ok := divisor.ModInverse(divisor, p)

		if ok == nil {
			return &Point{}
		}

		lambda.Mul(dividend, divisor)
		lambda.Mod(lambda, p)
	} else {
		deltaY := new(big.Int).Sub(q.y, pPoint.y)
		deltaY.Mod(deltaY, p)

		deltaX := new(big.Int).Sub(q.x, pPoint.x)
		ok := deltaX.ModInverse(deltaX, p)

		if ok == nil {
			return &Point{}
		}

		lambda = lambda.Mul(deltaY, deltaX)
		lambda = lambda.Mod(lambda, p)
	}

	x := new(big.Int).Mul(lambda, lambda)
	x.Sub(x, pPoint.x)
	x.Sub(x, q.x)
	x.Mod(x, p)

	y := new(big.Int).Sub(pPoint.x, x)
	y.Mul(y, lambda)
	y.Sub(y, pPoint.y)
	y.Mod(y, p)

	return &Point{x, y}
}

// Multiply a point pPoint by n
func (pPoint *Point) Mul(n int, a, b, p *big.Int) *Point {
	if n == 1 {
		return pPoint
	}

	r := pPoint.Add(pPoint, a, b, p)
	for i := 0; i < n; i++ {
		r = r.Add(pPoint, a, b, p)
	}

	return r
}

// Get the order of a point
func (pPoint *Point) getOrder(a, b, p *big.Int) int {
	counter := 2
	r := pPoint.Add(pPoint, a, b, p)

	for ; !r.IsAtInfinity(); counter++ {
		r = r.Add(pPoint, a, b, p)
	}

	return counter
}

// Returns true if point is on the curve
func (pPoint *Point) IsOnCurve(a, b, p *big.Int) bool {
	if pPoint.IsAtInfinity() {
		return true
	}

	ySquared := new(big.Int).Exp(pPoint.x, big.NewInt(3), nil)
	ySquared.Add(ySquared, new(big.Int).Mul(a, pPoint.x))
	ySquared.Add(ySquared, b)
	ySquared.Mod(ySquared, p)

	ySquaredRight := new(big.Int).Mul(pPoint.y, pPoint.y)
	ySquaredRight.Mod(ySquaredRight, p)

	return ySquaredRight.Cmp(ySquared) == 0
}

// Returns true if a group can be based the set E(a, b)
func isCurveValid(a, b, p *big.Int) bool {
	firstTerm := new(big.Int).Exp(a, big.NewInt(3), nil)
	secondTerm := new(big.Int).Mul(b, b)

	firstTerm.Mul(firstTerm, big.NewInt(4))
	secondTerm.Mul(secondTerm, big.NewInt(27))

	result := new(big.Int).Add(firstTerm, secondTerm)
	result.Mod(result, p)

	return !(result.Int64() == 0)
}

// Returns y^2
func getYSquared(x, a, b, p *big.Int) *big.Int {
	if !isCurveValid(a, b, p) {
		panic("Not valid curve, can't compute y squared!")
	}

	ySquared := new(big.Int).Exp(x, big.NewInt(3), nil)

	ySquared.Add(ySquared, new(big.Int).Mul(a, x))
	ySquared.Add(ySquared, b)

	return ySquared
}

// Given coordinates, returns a point
func getPoint(x, y int64) *Point {
	return &Point{
		big.NewInt(x),
		big.NewInt(y),
	}
}

// Returns the biggest order of all points
func getBiggestOrder(a, b, p *big.Int) int {
	_, orders := getAllPoints(a, b, p)

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

	pInt := p.Int64()
	for xInt := int64(0); xInt < pInt; xInt++ {
		x := big.NewInt(xInt)

		y := getYSquared(x, a, b, p)
		y = y.ModSqrt(y, p)

		if y != nil {
			pPoint := Point{x, y}
			order := pPoint.getOrder(a, b, p)

			points = append(points, pPoint)
			orders = append(orders, order)
			if y.Uint64() != 0 {
				pPoint = Point{x, new(big.Int).Sub(p, y)}
				order = pPoint.getOrder(a, b, p)

				points = append(points, pPoint)
				orders = append(orders, order)
			}
		}
	}

	return points, orders
}
