package svm

import (
	"errors"
	"os"
	"strconv"
)

const defaultTradeoff = 0.001

var (
	defaultRBFParams  = [][]float64{{1e-5}, {1e-4}, {1e-3}, {1e-2}, {1e-1}, {1e0}, {1e1}, {1e2}}
	defaultPolyPowers = []float64{2}
	defaultPolySums   = []float64{0, 1}
)

// These environment variables specify
// various parameters for the SVM trainer.
const (
	// Set this to "1" to get verbose logs.
	VerboseEnvVar = "SVM_VERBOSE"

	// You may set this to "linear", "rbf", or
	// "polynomial".
	KernelEnvVar = "SVM_KERNEL"

	// The numerical constant used in the
	// RBF kernel.
	RBFParamEnvVar = "SVM_RBF_PARAM"

	// The degree parameter for polynomial kernels.
	PolyDegreeEnvVar = "SVM_POLY_DEGREE"

	// The summed term (before applying the exponential)
	// for polynomial kernels.
	PolySumEnvVar = "SVM_POLY_SUM"

	// The tradeoff between margin size and hinge loss.
	// The higher the tradeoff value, the greater the
	// margin size, but at the expense of correct
	// classifications.
	TradeoffEnvVar = "SVM_TRADEOFF"
)

// TrainerParams specifies parameters for the
// SVM trainer.
type TrainerParams struct {
	Verbose  bool
	Kernels  []*Kernel
	Tradeoff float64
}

// EnvTrainerParams generates TrainerParams
// by reading environment variables.
// If an environment variable is incorrectly
// formatted, this returns an error.
// When a variable is missing, a default value
// or set of values will be used.
func EnvTrainerParams() (*TrainerParams, error) {
	var res TrainerParams
	var err error

	if res.Tradeoff, err = envTradeoff(); err != nil {
		return nil, err
	}
	res.Verbose = (os.Getenv(VerboseEnvVar) == "1")

	kernTypes, err := envKernelTypes()
	if err != nil {
		return nil, err
	}

	for _, kernType := range kernTypes {
		params, err := envKernelParams(kernType)
		if err != nil {
			return nil, err
		}
		for _, param := range params {
			kernel := &Kernel{
				Type:   kernType,
				Params: param,
			}
			res.Kernels = append(res.Kernels, kernel)
		}
	}

	return &res, nil
}

func envTradeoff() (float64, error) {
	if val := os.Getenv(TradeoffEnvVar); val != "" {
		return strconv.ParseFloat(val, 64)
	} else {
		return defaultTradeoff, nil
	}
}

func envKernelTypes() ([]KernelType, error) {
	if val := os.Getenv(KernelEnvVar); val != "" {
		res, ok := map[string]KernelType{
			"linear":     LinearKernel,
			"polynomial": PolynomialKernel,
			"rbf":        RadialBasisKernel,
		}[val]
		if !ok {
			return nil, errors.New("unknown kernel: " + val)
		} else {
			return []KernelType{res}, nil
		}
	} else {
		return []KernelType{LinearKernel, PolynomialKernel, RadialBasisKernel}, nil
	}
}

func envKernelParams(t KernelType) ([][]float64, error) {
	switch t {
	case LinearKernel:
		return [][]float64{}, nil
	case RadialBasisKernel:
		if val := os.Getenv(RBFParamEnvVar); val != "" {
			res, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, errors.New("invalid RBF param: " + val)
			}
			return [][]float64{{res}}, nil
		} else {
			return defaultRBFParams, nil
		}
	case PolynomialKernel:
		powers := defaultPolyPowers
		sums := defaultPolySums
		if val := os.Getenv(PolySumEnvVar); val != "" {
			sum, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, errors.New("invalid poly sum: " + val)
			}
			sums = []float64{sum}
		}
		if val := os.Getenv(PolyDegreeEnvVar); val != "" {
			degree, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, errors.New("invalid poly degree: " + val)
			}
			powers = []float64{degree}
		}
		res := make([][]float64, 0, len(powers)*len(sums))
		for _, power := range powers {
			for _, sum := range sums {
				res = append(res, []float64{sum, power})
			}
		}
		return res, nil
	default:
		panic("unknown kernel: " + strconv.Itoa(int(t)))
	}
}
