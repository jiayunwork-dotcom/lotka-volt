package lv

import "fmt"

type Params struct {
	Alpha float64 `json:"alpha"`
	Beta  float64 `json:"beta"`
	Gamma float64 `json:"gamma"`
	Delta float64 `json:"delta"`
}

type State struct {
	V float64 `json:"V"`
	P float64 `json:"P"`
}

type Equilibrium struct {
	V float64 `json:"V"`
	P float64 `json:"P"`
}

func (p Params) String() string {
	return fmt.Sprintf("alpha=%.4g beta=%.4g gamma=%.4g delta=%.4g", p.Alpha, p.Beta, p.Gamma, p.Delta)
}

func (s State) String() string {
	return fmt.Sprintf("V=%.6g P=%.6g", s.V, s.P)
}

func (s State) IsFinite() bool {
	return !bad(s.V) && !bad(s.P)
}

func bad(value float64) bool {
	return value != value || value > 1e300 || value < -1e300
}

func (p Params) Validate() error {
	if p.Alpha <= 0 {
		return fmt.Errorf("alpha must be positive, got %g", p.Alpha)
	}
	if p.Beta <= 0 {
		return fmt.Errorf("beta must be positive, got %g", p.Beta)
	}
	if p.Gamma <= 0 {
		return fmt.Errorf("gamma must be positive, got %g", p.Gamma)
	}
	if p.Delta <= 0 {
		return fmt.Errorf("delta must be positive, got %g", p.Delta)
	}
	return nil
}

func ValidateState(s State) error {
	if s.V < 0 {
		return fmt.Errorf("prey V must be non-negative, got %g", s.V)
	}
	if s.P < 0 {
		return fmt.Errorf("predator P must be non-negative, got %g", s.P)
	}
	if !s.IsFinite() {
		return fmt.Errorf("state must be finite")
	}
	return nil
}

func IsZero(s State) bool {
	return s.V == 0 && s.P == 0
}

func IsOnPreyAxis(s State) bool {
	return s.P == 0 && s.V >= 0
}

func IsOnPredatorAxis(s State) bool {
	return s.V == 0 && s.P >= 0
}

func SameState(a, b State, tolerance float64) bool {
	return abs(a.V-b.V) <= tolerance && abs(a.P-b.P) <= tolerance
}

func abs(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}
