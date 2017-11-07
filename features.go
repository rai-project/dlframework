package dlframework

import (
	"math"
	"sort"

	"github.com/pkg/errors"
	"gonum.org/v1/gonum/stat"
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

func (p Features) ProbabilitiesFloat32() []float32 {
	pProbs := make([]float32, p.Len())
	for ii := 0; ii < p.Len(); ii++ {
		pProbs[ii] = p[ii].Probability
	}
	return pProbs
}

func (p Features) ProbabilitiesFloat64() []float64 {
	pProbs := make([]float64, p.Len())
	for ii := 0; ii < p.Len(); ii++ {
		pProbs[ii] = float64(p[ii].Probability)
	}
	return pProbs
}

func (p Features) KullbackLeiblerDivergence(q Features) (float64, error) {
	if p.Len() != q.Len() {
		return 0, errors.Errorf("length mismatch %d != %d", p.Len(), q.Len())
	}

	pProbs := p.ProbabilitiesFloat64()
	qProbs := q.ProbabilitiesFloat64()

	return stat.KullbackLeibler(pProbs, qProbs), nil
}

func (p Features) Correlation(q Features) (float64, error) {
	if p.Len() != q.Len() {
		return 0, errors.Errorf("length mismatch %d != %d", p.Len(), q.Len())
	}

	pProbs := p.ProbabilitiesFloat64()
	qProbs := q.ProbabilitiesFloat64()

	return stat.Correlation(pProbs, qProbs, nil), nil
}
