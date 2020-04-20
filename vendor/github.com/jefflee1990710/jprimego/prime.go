package jprimego

import (
	"container/list"
	"crypto/rand"
	"math/big"
)

// FastGeneratePrime find prime number by finding nearest prime of a large integer
func FastGeneratePrime(bitLength int64) *big.Int {
	var TWO = new(big.Int).SetInt64(2)
	var THREE = new(big.Int).SetInt64(3)

	BL := new(big.Int).SetInt64(int64(bitLength))

	if BL.Cmp(TWO) == -1 {
		panic("bitLength cannot smaller then 2")
	} else if BL.Cmp(TWO) == 0 {
		return THREE
	} else {
		max := getRangeByBitLength(BL)
		n := createRandomBetween(TWO, max)
		return findNearestPrimeNumber(n, max)
	}
}

// sign1Gen is a mapping function map {1, 2, 3, 4} -> {-1, 1, -1, 1}
func sign1Gen(k int) int {
	return 1 - 2*(k%2)
}

// sign2Gen is a mapping function map {1, 2, 3, 4} -> {-1, 1, 1, -1}
func sign2Gen(k int) int {
	return 1 - 2*((k%3)%2)
}

// generate4Possible ....
func generate4Possible(k int) *list.List {
	rs := list.New()
	for i := 0; i < 4; i++ {
		sign1 := sign1Gen(i + 1)
		sign2 := sign2Gen(i + 1)
		r := (sign1) * (6*k + (sign2))
		rs.PushBack(int64(r))
	}
	return rs
}

// generateNextPossiblePrime generate candidate for prime testing +-(6k +- 1)
func generateNextPossiblePrime(c *big.Int) chan *big.Int {
	ch := make(chan *big.Int)
	go func() {
		k := 0
		next := true
		for next {
			l := generate4Possible(k)
			for e := l.Front(); e != nil; e = e.Next() {
				// c + p
				P := new(big.Int).SetInt64(e.Value.(int64))
				var r big.Int
				r.Sub(c, P)
				ch <- &r
			}
			k = k + 1
		}
	}()
	return ch
}

// findNearestPrimeNumber find the nearest prime number which must smaller then a limit
func findNearestPrimeNumber(n *big.Int, max *big.Int) *big.Int {
	var SIX = new(big.Int).SetInt64(6)
	var ZERO = new(big.Int).SetInt64(1)

	c := new(big.Int).Mod(n, SIX)
	c.Sub(n, c)

	cnt := 0
	for p := range generateNextPossiblePrime(c) {
		if p.Cmp(ZERO) == 1 && (max == nil || p.Cmp(max) == -1) {
			cnt = 0
			if p.ProbablyPrime(40) {
				return p
			}
		} else {
			cnt = cnt + 1
		}
		if cnt == 4 {
			panic("Prime number not found in range")
		}
	}
	return n
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

// getRangeByBitLength get maximun value of a specific bit length integer
func getRangeByBitLength(bitLength *big.Int) *big.Int {
	var ONE = new(big.Int).SetInt64(1)
	var TWO = new(big.Int).SetInt64(2)

	// s = 2k - 1
	s := new(big.Int).Exp(TWO, bitLength, nil)
	s = s.Sub(s, ONE)
	return s
}
