// Package sieve TODO: document
//
// references
// https://research.gold.ac.uk/id/eprint/15753/1/11.2-Dimitris-Exarchos-&-Daniel-Jones.pdf
// http://kunstmusik.github.io/score/sieves.html
// https://github.com/deckarep/golang-set
package sieve

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/fatih/set"
	"github.com/go-audio/audio"
)

type Sieve struct {
}

func NewSieve() *Sieve {
	return &Sieve{}
}

func (s *Sieve) Fill(buf *audio.PCMBuffer) error {
	if s == nil {
		return nil
	}

	numChans := 1
	if f := buf.Format; f != nil {
		numChans = f.NumChannels
	}

	frameCount := buf.NumFrames()
	for i := 0; i < frameCount; i++ {
		for j := 0; j < numChans; j++ {
			buf.F64[i*numChans+j] = rand.Float64()
		}
	}

	return nil
}

func Intervallic(s []int) []int {
	var i []int
loop:
	for k := range s {
		switch {
		case k+1 >= len(s):
			break loop

		case s[k] < s[k+1]:
			i = append(i, s[k+1]-s[k])

		case s[k] > s[k+1]:
			i = append(i, s[k]-s[k+1])
		}
	}
	return i
}

func HH() {
	s := set.New(set.ThreadSafe)
	s.Add(1, 3, 4, 5)

	t := set.New(set.ThreadSafe)
	t.Add(1, 2, 3)

	// check if they are the same
	if !s.IsEqual(t) {
		fmt.Println("s is not equal to t")
	}

	// if s contains all elements of t
	if s.IsSubset(t) {
		fmt.Println("t is a subset of s")
	}

	// ... or if s is a superset of t
	if t.IsSuperset(s) {
		fmt.Println("s is a superset of t")
	}

	c := set.Union(s, t)
	sieve := make([]int, len(c.List()))
	for i, v := range c.List() {
		sieve[i] = v.(int)
	}
	sort.Ints(sieve)
	fmt.Println(sieve)
	// fmt.Println("union", c.List())

	// contains items which is in both a and b
	// [berlin]
	c = set.Intersection(s, t)
	fmt.Println(c)

	// contains items which are in a but not in b
	// [ankara, san francisco]
	c = set.Difference(s, t)
	fmt.Println(c)

	// contains items which are in one of either, but not in both.
	// [frankfurt, ankara, san francisco]
	c = set.SymmetricDifference(s, t)
	fmt.Println(c)
}
