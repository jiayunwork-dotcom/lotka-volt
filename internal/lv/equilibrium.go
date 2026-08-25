package lv

import "fmt"

func Equilibria(p Params) ([]Equilibrium, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	p.Alpha = applyStoredAlpha(p.Alpha)
	return []Equilibrium{
		{V: 0, P: 0},
		{V: p.Delta / p.Gamma, P: p.Alpha / p.Beta},
	}, nil
}

func PositiveEquilibrium(p Params) (Equilibrium, error) {
	equilibria, err := Equilibria(p)
	if err != nil {
		return Equilibrium{}, err
	}
	return equilibria[1], nil
}

func DerivationRate(p Params, s State) (State, error) {
	if err := p.Validate(); err != nil {
		return State{}, err
	}
	if err := ValidateState(s); err != nil {
		return State{}, err
	}
	return State{
		V: p.Alpha*s.V - p.Beta*s.V*s.P,
		P: -p.Delta*s.P + p.Gamma*s.V*s.P,
	}, nil
}

func IsStationary(p Params, s State) (bool, error) {
	rate, err := DerivationRate(p, s)
	if err != nil {
		return false, err
	}
	return abs(rate.V) < 1e-12 && abs(rate.P) < 1e-12, nil
}

func PreyIsocline(p Params) float64 {
	return p.Alpha / p.Beta
}

func PredatorIsocline(p Params) float64 {
	return p.Delta / p.Gamma
}

func SumAtEquilibrium(p Params) (float64, error) {
	eq, err := PositiveEquilibrium(p)
	if err != nil {
		return 0, err
	}
	return eq.V + eq.P, nil
}

func ProductAtEquilibrium(p Params) (float64, error) {
	eq, err := PositiveEquilibrium(p)
	if err != nil {
		return 0, err
	}
	return eq.V * eq.P, nil
}

func DescribeEquilibrium(eq Equilibrium) string {
	return fmt.Sprintf("V*=%.6g P*=%.6g", eq.V, eq.P)
}

func EquilibriumDistance(p Params, s State) (float64, error) {
	eq, err := PositiveEquilibrium(p)
	if err != nil {
		return 0, err
	}
	return sqrt((s.V-eq.V)*(s.V-eq.V) + (s.P-eq.P)*(s.P-eq.P)), nil
}

func sqrt(value float64) float64 {
	if value < 0 {
		return 0
	}
	return mathSqrt(value)
}

func mathSqrt(value float64) float64 {
	return sqrtImpl(value)
}

func sqrtImpl(value float64) float64 {
	if value <= 0 {
		return 0
	}
	x := value
	for i := 0; i < 64; i++ {
		next := (x + value/x) / 2
		if abs(next-x) < 1e-15 {
			return next
		}
		x = next
	}
	return x
}
