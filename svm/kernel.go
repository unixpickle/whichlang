package svm

import (
	"math"
	"strconv"

	"github.com/unixpickle/num-analysis/linalg"
)

type KernelType int

const (
	// LinearKernel generates a linear classifier
	// with no parameters.
	LinearKernel KernelType = iota

	// PolynomialKernel computes inner products as
	// (x*y + k1)^k2, where k1 and k2 are parameters.
	PolynomialKernel

	// RadialBasisKernel computes inner products as
	// exp(-k1*||x-y||^2), where k1 is a parameter.
	RadialBasisKernel
)

// A Kernel computes inner products of vectors
// after transforming them into some space.
type Kernel struct {
	Type   KernelType
	Params []float64
}

func (k *Kernel) Product(v1, v2 linalg.Vector) float64 {
	switch k.Type {
	case LinearKernel:
		return v1.Dot(v2)
	case PolynomialKernel:
		if len(k.Params) != 2 {
			panic("expected two parameters for polynomial kernel")
		}
		return math.Pow(v1.Dot(v2)+k.Params[0], k.Params[1])
	case RadialBasisKernel:
		if len(k.Params) != 1 {
			panic("expected one parameter for radial basis kernel")
		}
		diff := v1.Copy().Scale(-1).Add(v2)
		return math.Exp(-k.Params[0] * diff.Dot(diff))
	default:
		panic("unknown kernel type: " + strconv.Itoa(int(k.Type)))
	}
}
