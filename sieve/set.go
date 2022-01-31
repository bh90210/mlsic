package sieve

import (
	"fmt"
	"sort"

	"github.com/fatih/set"
)

func hh() {
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
