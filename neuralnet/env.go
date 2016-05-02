package neuralnet

import (
	"math"
	"os"
	"strconv"
)

const DefaultMaxIterations = 6400

// DefaultHiddenLayerScale specifies how much
// larger the hidden layer is than the output
// layer, by default.
const DefaultHiddenLayerScale = 2.0

// VerboseEnvVar is an environment variable
// which can be set to "1" to make the
// neuralnet print out status reports.
var VerboseEnvVar = "NEURALNET_VERBOSE"

// VerboseStepsEnvVar is an environment
// variable which can be set to "1" to make
// neuralnet print out status reports after
// each iteration of gradient descent.
var VerboseStepsEnvVar = "NEURALNET_VERBOSE_STEPS"

// StepSizeEnvVar is an environment variable
// which can be used to specify the step size
// for use in gradient descent.
var StepSizeEnvVar = "NEURALNET_STEP_SIZE"

// MaxItersEnvVar is an environment variable
// specifying the maximum number of iterations
// of gradient descent to perform.
var MaxItersEnvVar = "NEURALNET_MAX_ITERS"

// HiddenSizeEnvVar is an environment variable
// specifying the number of hidden neurons.
var HiddenSizeEnvVar = "NEURALNET_HIDDEN_SIZE"

func verboseFlag() bool {
	return os.Getenv(VerboseEnvVar) == "1"
}

func verboseStepsFlag() bool {
	return os.Getenv(VerboseStepsEnvVar) == "1"
}

func stepSizes() []float64 {
	if stepSize := os.Getenv(StepSizeEnvVar); stepSize == "" {
		var res []float64
		for power := -20; power < 10; power++ {
			res = append(res, math.Pow(2, float64(power)))
		}
		return res
	} else {
		val, err := strconv.ParseFloat(stepSize, 64)
		if err != nil {
			panic(err)
		}
		return []float64{val}
	}
}

func maxIterations() int {
	if max := os.Getenv(MaxItersEnvVar); max == "" {
		return DefaultMaxIterations
	} else {
		val, err := strconv.Atoi(max)
		if err != nil {
			panic(err)
		}
		return val
	}
}

func hiddenSize(outputCount int) int {
	if size := os.Getenv(HiddenSizeEnvVar); size == "" {
		return int(float64(outputCount)*DefaultHiddenLayerScale + 0.5)
	} else {
		val, err := strconv.Atoi(size)
		if err != nil {
			panic(err)
		}
		return val
	}
}
