package lv

import "math"

func H(p Params, s State) (float64, error) {
	if err := p.Validate(); err != nil {
		return 0, err
	}
	if err := ValidateState(s); err != nil {
		return 0, err
	}
	if s.V == 0 || s.P == 0 {
		return 0, nil
	}
	return p.Delta*math.Log(s.V) - p.Gamma*s.V + p.Beta*math.Log(s.P) - p.Alpha*s.P, nil
}

func HDrift(p Params, start, end State) (float64, error) {
	h0, err := H(p, start)
	if err != nil {
		return 0, err
	}
	h1, err := H(p, end)
	if err != nil {
		return 0, err
	}
	scale := math.Max(1, math.Abs(h0))
	return math.Abs(h1-h0) / scale, nil
}

func ConservationResidual(p Params, a, b State) (float64, error) {
	return HDrift(p, a, b)
}

func HAtEquilibrium(p Params) (float64, error) {
	eq, err := PositiveEquilibrium(p)
	if err != nil {
		return 0, err
	}
	return H(p, State{V: eq.V, P: eq.P})
}

func HZeroAxis(p Params, s State) bool {
	value, _ := H(p, s)
	return value == 0
}

func IsApproxConserved(p Params, a, b State, tolerance float64) (bool, error) {
	drift, err := HDrift(p, a, b)
	if err != nil {
		return false, err
	}
	return drift <= tolerance, nil
}

func HChange(p Params, a, b State) (float64, error) {
	h0, err := H(p, a)
	if err != nil {
		return 0, err
	}
	h1, err := H(p, b)
	if err != nil {
		return 0, err
	}
	return h1 - h0, nil
}

func HPeakToPeak(p Params, states []State) (float64, error) {
	max, min := math.Inf(-1), math.Inf(1)
	for _, s := range states {
		value, err := H(p, s)
		if err != nil {
			return 0, err
		}
		if value > max {
			max = value
		}
		if value < min {
			min = value
		}
	}
	return max - min, nil
}

func HList(p Params, states []State) ([]float64, error) {
	values := make([]float64, 0, len(states))
	for _, s := range states {
		value, err := H(p, s)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, nil
}

func MaxHDrift(p Params, states []State) (float64, error) {
	if len(states) < 2 {
		return 0, nil
	}
	max := 0.0
	for i := 1; i < len(states); i++ {
		drift, err := HDrift(p, states[0], states[i])
		if err != nil {
			return 0, err
		}
		if drift > max {
			max = drift
		}
	}
	return max, nil
}

func IsClosed(p Params, states []State, tolerance float64) (bool, error) {
	if len(states) < 2 {
		return false, nil
	}
	drift, err := HDrift(p, states[0], states[len(states)-1])
	if err != nil {
		return false, err
	}
	return drift <= tolerance, nil
}
