// Command crandall searches for Crandall primes close to word boundaries.
package main

import (
	"flag"
	"fmt"
	"math/big"
)

// Crandall is a prime of the form 2ⁿ - c.
type Crandall struct {
	N int
	C int
}

// NewCrandall constructs a Crandall prime.
func NewCrandall(n, c int) Crandall {
	return Crandall{N: n, C: c}
}

func (p Crandall) String() string {
	return fmt.Sprintf("2^%d%+d", p.N, -p.C)
}

// Int returns the prime as a big integer.
func (p Crandall) Int() *big.Int {
	one := big.NewInt(1)
	c := big.NewInt(int64(p.C))
	e := new(big.Int).Lsh(one, uint(p.N))
	return new(big.Int).Sub(e, c)
}

// Command line flags.
var (
	trials = flag.Int("trials", 32, "how many trials of the Miller-Rabin test")
	s      = flag.Int("s", 64, "word size in bits")
	minc   = flag.Int("min-c", -128, "minimum value of c to try")
	maxc   = flag.Int("max-c", 128, "maximum value of c to try")
)

func main() {
	flag.Parse()

	for b := 128; b <= 1024; b += *s {
		for s := 0; s <= 2; s++ {
			n := b - s
			for c := *minc; c < *maxc; c++ {
				p := NewCrandall(n, c)
				P := p.Int()
				prime := P.ProbablyPrime(*trials)
				if prime {
					fmt.Printf("%12s\tn=%d\tc=%d\tbits=%d\tp=%d\n", p, n, c, P.BitLen(), P)
				}
			}
		}
	}
}
