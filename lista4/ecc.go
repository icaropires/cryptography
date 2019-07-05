package main

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
)

type Point struct {
	x *big.Int
	y *big.Int
}

type Curve struct {
	a *big.Int
	b *big.Int
	p *big.Int
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
func (pPoint Point) Neg() *Point {
	if !pPoint.IsAtInfinity() {
		pPoint = Point{new(big.Int).Set(pPoint.x), new(big.Int).Set(pPoint.y)}
	}

	return &Point{new(big.Int).Set(pPoint.x), new(big.Int).Neg(pPoint.y)}
}

// Returns true if a pPoint is a point at the infinity
func (pPoint *Point) IsAtInfinity() bool {
	return pPoint.x == nil && pPoint.y == nil
}

// Add a pPoint pPoint to a point q
func (pPoint *Point) Add(q *Point, curve *Curve) *Point {
	if !pPoint.IsOnCurve(curve) {
		msg := fmt.Sprintf("Point '%v' not in curve: %v", pPoint, curve)
		panic(msg)
	}

	if !q.IsOnCurve(curve) {
		msg := fmt.Sprintf("Point '%v' not in curve: %v", q, curve)
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
	if pPoint.IsEqual(qNeg) {
		return &Point{}
	}

	lambda := new(big.Int)
	if pPoint.IsEqual(q) {
		dividend := new(big.Int).Mul(pPoint.x, pPoint.x)
		dividend.Mul(dividend, big.NewInt(3))
		dividend.Add(dividend, curve.a)
		dividend.Mod(dividend, curve.p)

		divisor := new(big.Int).Mul(big.NewInt(2), pPoint.y)
		ok := divisor.ModInverse(divisor, curve.p)

		if ok == nil {
			return &Point{}
		}

		lambda.Mul(dividend, divisor)
		lambda.Mod(lambda, curve.p)
	} else {
		deltaY := new(big.Int).Sub(q.y, pPoint.y)
		deltaY.Mod(deltaY, curve.p)

		deltaX := new(big.Int).Sub(q.x, pPoint.x)
		ok := deltaX.ModInverse(deltaX, curve.p)

		if ok == nil {
			return &Point{}
		}

		lambda = lambda.Mul(deltaY, deltaX)
		lambda = lambda.Mod(lambda, curve.p)
	}

	x := new(big.Int).Mul(lambda, lambda)
	x.Sub(x, pPoint.x)
	x.Sub(x, q.x)
	x.Mod(x, curve.p)

	y := new(big.Int).Sub(pPoint.x, x)
	y.Mul(y, lambda)
	y.Sub(y, pPoint.y)
	y.Mod(y, curve.p)

	return &Point{x, y}
}

// Multiply a point pPoint by n
func (pPoint *Point) Mul(n *big.Int, curve *Curve) *Point {
	if n.Uint64() == 1 {
		return pPoint
	}

	r := pPoint.Add(pPoint, curve)
	for i := int64(0); i < n.Int64()-2; i++ {
		r = r.Add(pPoint, curve)
	}

	return r
}

// Get the order of a point
func (pPoint *Point) getOrder(curve *Curve) int {
	counter := 2
	r := pPoint.Add(pPoint, curve)

	for ; !r.IsAtInfinity(); counter++ {
		r = r.Add(pPoint, curve)
	}

	return counter
}

// Returns true if point is on the curve
func (pPoint *Point) IsOnCurve(curve *Curve) bool {
	if pPoint.IsAtInfinity() {
		return true
	}

	ySquared := new(big.Int).Exp(pPoint.x, big.NewInt(3), nil)
	ySquared.Add(ySquared, new(big.Int).Mul(curve.a, pPoint.x))
	ySquared.Add(ySquared, curve.b)
	ySquared.Mod(ySquared, curve.p)

	ySquaredRight := new(big.Int).Mul(pPoint.y, pPoint.y)
	ySquaredRight.Mod(ySquaredRight, curve.p)

	return ySquaredRight.Cmp(ySquared) == 0
}

func (curve *Curve) String() string {
	return fmt.Sprintf("Curve(%s, %s, %s)", curve.a.String(), curve.b.String(), curve.p.String())
}

// Returns true if a group can be based the set E(a, b)
func isCurveValid(curve *Curve) bool {
	firstTerm := new(big.Int).Exp(curve.a, big.NewInt(3), nil)
	secondTerm := new(big.Int).Mul(curve.b, curve.b)

	firstTerm.Mul(firstTerm, big.NewInt(4))
	secondTerm.Mul(secondTerm, big.NewInt(27))

	result := new(big.Int).Add(firstTerm, secondTerm)
	result.Mod(result, curve.p)

	return !(result.Int64() == 0)
}

// Returns y^2
func getYSquared(x *big.Int, curve *Curve) *big.Int {
	if !isCurveValid(curve) {
		panic("Not valid curve, can't compute y squared!")
	}

	ySquared := new(big.Int).Exp(x, big.NewInt(3), nil)

	ySquared.Add(ySquared, new(big.Int).Mul(curve.a, x))
	ySquared.Add(ySquared, curve.b)

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
func getBiggestOrder(curve *Curve) int64 {
	_, orders := getAllPoints(curve)

	biggest := float64(orders[0])
	for _, e := range orders {
		biggest = math.Max(float64(biggest), float64(e))
	}

	return int64(biggest)
}

// Get all points of a curve with its order
func getAllPoints(curve *Curve) ([]Point, []int) {
	points := make([]Point, 0)
	orders := make([]int, 0)

	pInt := curve.p.Int64()
	for xInt := int64(0); xInt < pInt; xInt++ {
		x := big.NewInt(xInt)

		y := getYSquared(x, curve)
		y = y.ModSqrt(y, curve.p)

		if y != nil {
			pPoint := Point{x, y}
			order := pPoint.getOrder(curve)

			points = append(points, pPoint)
			orders = append(orders, order)
			if y.Uint64() != 0 {
				pPoint = Point{x, new(big.Int).Sub(curve.p, y)}
				order = pPoint.getOrder(curve)

				points = append(points, pPoint)
				orders = append(orders, order)
			}
		}
	}

	return points, orders
}

// N is the upper value for privateKey
func GenKeys(g *Point, curve *Curve, n *big.Int) (privateKey *big.Int, publicKey *Point) {
	privateKey = getRandom(n)
	publicKey = g.Mul(privateKey, curve)

	return
}

func getRandom(n *big.Int) *big.Int {
	zero := big.NewInt(0)
	random := zero

	for random.Cmp(zero) == 0 {
		random, _ = rand.Int(rand.Reader, n)
	}

	return random
}

func Cipher(pPoint, publicKey, g *Point, curve *Curve) (c1, c2 *Point) {
	biggest := getBiggestOrder(curve)
	k := getRandom(big.NewInt(biggest))

	c1 = g.Mul(k, curve)

	aux := publicKey.Mul(k, curve)
	c2 = pPoint.Add(aux, curve)

	return
}

func Decipher(c1, c2 *Point, privateKey *big.Int, curve *Curve) *Point {
	aux := c1.Mul(privateKey, curve)
	aux = aux.Neg()

	plain := c2.Add(aux, curve)

	return plain
}

func MapCharToPoint(message byte) *Point {
	const (
		p = 4177
		r = 30
	)

	j := int64(0)
	for {
		y := big.NewInt(int64(message)*r + j)

		x := new(big.Int).Set(y)

		rr := y.ModSqrt(new(big.Int).Sub(new(big.Int).Exp(y, big.NewInt(3), nil), big.NewInt(4)), big.NewInt(p))

		if rr != nil {
			return &Point{x, y}
		}
		j += 1
	}
}
