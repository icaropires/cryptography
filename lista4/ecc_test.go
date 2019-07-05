package main

import (
	"fmt"
	"math/big"
	"testing"
)

func TestIsAtInfinity(t *testing.T) {
	pPoint := &Point{}

	if !pPoint.IsAtInfinity() {
		t.Errorf("Is at Infinity!")
	}
}

func TestAddInverse(t *testing.T) {
	curve := &Curve{big.NewInt(1), big.NewInt(1), big.NewInt(23), big.NewInt(28)}

	pPoint := getPoint(3, 13)
	q := getPoint(3, -13)

	r := pPoint.Add(q, curve)
	if !r.IsAtInfinity() {
		t.Errorf("P + (-P) must be at Infinity")
	}
}

func TestAddPInfinity(t *testing.T) {
	curve := &Curve{big.NewInt(1), big.NewInt(1), big.NewInt(23), big.NewInt(28)}

	pPoint := &Point{}
	q := getPoint(3, -13)

	r := pPoint.Add(q, curve)
	if !q.IsEqual(r) {
		t.Errorf("Q + Infinity must be equal Q")
	}
}

func TestAddQInfinity(t *testing.T) {
	curve := &Curve{big.NewInt(1), big.NewInt(1), big.NewInt(23), big.NewInt(28)}

	pPoint := getPoint(3, -13)
	q := &Point{}

	r := pPoint.Add(q, curve)
	if !r.IsEqual(pPoint) {
		t.Errorf("P + At Infinity must be equal P")
	}
}

func TestAddAllInfinity(t *testing.T) {
	curve := &Curve{big.NewInt(1), big.NewInt(1), big.NewInt(23), big.NewInt(28)}

	pPoint := &Point{}
	q := &Point{}

	r := pPoint.Add(q, curve)
	if !r.IsAtInfinity() {
		t.Errorf("At Infinity + At Infinity must be equal Point at Infinity")
	}
}

func TestAddPxEqualQx(t *testing.T) {
	curve := &Curve{big.NewInt(1), big.NewInt(1), big.NewInt(23), big.NewInt(28)}

	pPoint := getPoint(11, 3)
	q := getPoint(11, 20)

	r := pPoint.Add(q, curve)
	if !r.IsAtInfinity() {
		t.Errorf("P + Q with Px = Qx must be equal Point at Infinity")
	}
}

func TestAddEqualYZero(t *testing.T) {
	curve := &Curve{big.NewInt(1), big.NewInt(1), big.NewInt(23), big.NewInt(28)}

	pPoint := getPoint(4, 0)
	q := getPoint(4, 0)

	r := pPoint.Add(q, curve)
	if !r.IsAtInfinity() {
		t.Errorf("P + P with Py = 0 must be equal Point at Infinity")
	}
}

func TestMul(t *testing.T) {
	fmt.Println("=========================================================")
	fmt.Println("=========================================================")
	fmt.Println("=========================================================")
	fmt.Println("=========================================================")
	curve := &Curve{big.NewInt(0), big.NewInt(-4), big.NewInt(257), big.NewInt(258)}
	n := big.NewInt(101)
	pPoint := getPoint(2, 2)
	pPointRight := getPoint(197, 167)

	r := pPoint.Mul(n, curve)

	if !r.IsEqual(pPointRight) {
		t.Errorf("Wrong multiplication! %v != %v", pPointRight, r)
	}
}
