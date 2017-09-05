package dlframework

import (
	"math"
	"sort"
)

type Features []*Feature

// Len is the number of elements in the collection.
func (p Features) Len() int {
	return len(p)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (p Features) Less(i, j int) bool {
	pi := p[i].Probability
	pj := p[j].Probability
	return !(pi < pj || math.IsNaN(float64(pi)) && !math.IsNaN(float64(pj)))
}

// Swap swaps the elements with indexes i and j.
func (p Features) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Features) Sort() {
	sort.Sort(p)
}

func (p Features) Take(n int) Features {
	if p.Len() <= n {
		return p
	}
	return Features(p[:n])
}
