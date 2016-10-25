package adaptivelimiter

import "math"

// The magnitude of delta will be greater than min, or be 0
func MinSmoother(min float64) SmoothingFunction {
	return func(delta float64) float64 {
		if math.Abs(delta) <= min {
			return 0
		}
		return delta
	}
}

// The magnitude of delta will be less than max, or be max
func MaxSmoother(max float64) SmoothingFunction {
	return func(delta float64) float64 {
		if math.Abs(delta) > max {
			return math.Copysign(max, delta)
		}
		return delta
	}
}

func ConstantGainSmoother(gain float64) SmoothingFunction {
	return func(delta float64) float64 {
		return gain * delta
	}
}

func GaussianSmoother(mean, stddev float64) SmoothingFunction {
	gaussian := func(x float64) float64 {
		if x < 0 {
			x = -x
		}
		return 1 / (2 * math.Sqrt(x) * math.Pi) * math.Exp(-(x-mean)*(x-mean)/(2*stddev*stddev))
	}
	return func(delta float64) float64 {
		return gaussian(mean) * delta * gaussian(delta)
	}
}

func LogSmoother() SmoothingFunction {
	return func(delta float64) float64 {
		if delta < 0 {
			return -math.Log(-delta + 1)
		}
		return math.Log(delta + 1)
	}
}
