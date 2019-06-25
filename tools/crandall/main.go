// Command crandall searches for Crandall primes close to word boundaries.
package main

import (
	"flag"
	"fmt"

	"github.com/mmcloughlin/ec3/prime"
)

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
				p := prime.NewCrandall(n, c)
				P := p.Int()
				if P.ProbablyPrime(*trials) {
					fmt.Printf("%12s\tn=%d\tc=%d\tbits=%d\tp=%d\n", p, n, c, P.BitLen(), P)
				}
			}
		}
	}
}
