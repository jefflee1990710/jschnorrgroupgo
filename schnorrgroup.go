package jschnorrgroupgo

import (
	"crypto/rand"
	"fmt"
	"jprimego"
	"math/big"
)

// SchnorrGroup ...
type SchnorrGroup struct {
	p *big.Int
	q *big.Int
}

// New init Schnorr Group with p and q
func (sg *SchnorrGroup) New(n int64) {
	var ONE = new(big.Int).SetInt64(1)
	var TWO = new(big.Int).SetInt64(2)
	next := true
	for next {
		q := jprimego.FastGeneratePrime(n)
		var p big.Int
		p.Mul(TWO, q)
		p.Add(&p, ONE)

		sg.p = &p
		sg.q = q

		if p.ProbablyPrime(40) {
			next = false
		}
	}
}

// GetGenerator get generate base on Schnorr Group requirement
func (sg *SchnorrGroup) GetGenerator() big.Int {
	var ONE = new(big.Int).SetInt64(1)
	var TWO = new(big.Int).SetInt64(2)

	P1 := new(big.Int)
	P1.Sub(sg.p, ONE)

	g := new(big.Int)

	var found = false
	for !found {
		h := createRandomBetween(TWO, P1)
		g.Exp(h, TWO, sg.p)
		if g.Cmp(ONE) != 0 {
			found = true
		}
	}
	return *g
}

// Summary print information of p and q
func (sg *SchnorrGroup) Summary() {
	fmt.Println("p = ", sg.p.Text(10))
	fmt.Println("q = ", sg.q.Text(10))
}

// GetP retrieve p of the group
func (sg *SchnorrGroup) GetP() *big.Int {
	return sg.p
}

// GetQ retrieve p of the group
func (sg *SchnorrGroup) GetQ() *big.Int {
	return sg.q
}

// createRandomBetween create random number between start and end inclusively
func createRandomBetween(start *big.Int, end *big.Int) *big.Int {
	var ONE = new(big.Int).SetInt64(1)

	len := new(big.Int)
	len.Sub(end, start)
	len.Add(len, ONE)

	n, err := rand.Int(rand.Reader, len)
	if err != nil {
		panic(err)
	}
	r := new(big.Int)
	r.Add(n, start)
	return r
}
