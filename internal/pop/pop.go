package pop

import (
	"fmt"
	"math"

	"lotka-volt/internal/lv"
)

func GrowthRate(p lv.Params, s lv.State) (lv.State, error) {
	return lv.DerivationRate(p, s)
}

func PreyPerCapita(p lv.Params, s lv.State) float64 {
	if s.V == 0 {
		return 0
	}
	return p.Alpha - p.Beta*s.P
}

func PredatorPerCapita(p lv.Params, s lv.State) float64 {
	if s.P == 0 {
		return 0
	}
	return p.Gamma*s.V - p.Delta
}

func DoublingTime(rate float64) float64 {
	if rate <= 0 {
		return math.Inf(1)
	}
	return math.Ln2 / rate
}

func HalfLife(rate float64) float64 {
	if rate >= 0 {
		return math.Inf(1)
	}
	return math.Ln2 / (-rate)
}

func EquilibriumState(p lv.Params) (lv.State, error) {
	eq, err := lv.PositiveEquilibrium(p)
	if err != nil {
		return lv.State{}, err
	}
	return lv.State{V: eq.V, P: eq.P}, nil
}

func DistanceToEquilibrium(p lv.Params, s lv.State) (float64, error) {
	eq, err := lv.PositiveEquilibrium(p)
	if err != nil {
		return 0, err
	}
	dx := s.V - eq.V
	dy := s.P - eq.P
	return math.Sqrt(dx*dx + dy*dy), nil
}

func PreyDominant(s lv.State) bool {
	return s.V > s.P
}

func PredatorDominant(s lv.State) bool {
	return s.P > s.V
}

func TotalPopulation(s lv.State) float64 {
	return s.V + s.P
}

func PredatorFraction(s lv.State) float64 {
	total := TotalPopulation(s)
	if total == 0 {
		return 0
	}
	return s.P / total
}

func PreyFraction(s lv.State) float64 {
	total := TotalPopulation(s)
	if total == 0 {
		return 0
	}
	return s.V / total
}

func Ratio(s lv.State) float64 {
	if s.P == 0 {
		return math.Inf(1)
	}
	return s.V / s.P
}

func StateText(s lv.State) string {
	return fmt.Sprintf("V=%.4g P=%.4g total=%.4g", s.V, s.P, TotalPopulation(s))
}

func PerCapitaText(p lv.Params, s lv.State) string {
	return fmt.Sprintf("dV/V=%.4g dP/P=%.4g", PreyPerCapita(p, s), PredatorPerCapita(p, s))
}

func IsIncreasing(s lv.State) bool {
	return s.V > 0 && s.P > 0
}

func IsAtOrigin(s lv.State) bool {
	return s.V == 0 && s.P == 0
}

func IsPreyOnly(s lv.State) bool {
	return s.V > 0 && s.P == 0
}

func IsPredatorOnly(s lv.State) bool {
	return s.V == 0 && s.P > 0
}

func PreyGrowthContribution(p lv.Params, s lv.State) float64 {
	return p.Alpha * s.V
}

func PredationLoss(p lv.Params, s lv.State) float64 {
	return p.Beta * s.V * s.P
}

func PredatorGrowthContribution(p lv.Params, s lv.State) float64 {
	return p.Delta * s.V * s.P
}

func PredatorDeathLoss(p lv.Params, s lv.State) float64 {
	return p.Gamma * s.P
}

func BalanceText(p lv.Params, s lv.State) string {
	return fmt.Sprintf("prey_gain=%.4g prey_loss=%.4g pred_gain=%.4g pred_loss=%.4g",
		PreyGrowthContribution(p, s), PredationLoss(p, s),
		PredatorGrowthContribution(p, s), PredatorDeathLoss(p, s))
}

func PhaseLabel(s lv.State) string {
	switch {
	case IsAtOrigin(s):
		return "origin"
	case IsPreyOnly(s):
		return "prey-axis"
	case IsPredatorOnly(s):
		return "predator-axis"
	default:
		return "interior"
	}
}

func TrajectoryLength(states []lv.State) float64 {
	length := 0.0
	for i := 1; i < len(states); i++ {
		dx := states[i].V - states[i-1].V
		dy := states[i].P - states[i-1].P
		length += math.Sqrt(dx*dx + dy*dy)
	}
	return length
}

func MaxVelocity(p lv.Params, states []lv.State) float64 {
	max := 0.0
	for _, s := range states {
		rate, _ := lv.DerivationRate(p, s)
		magnitude := math.Sqrt(rate.V*rate.V + rate.P*rate.P)
		if magnitude > max {
			max = magnitude
		}
	}
	return max
}

func MeanPopulation(states []lv.State) (float64, float64) {
	sumV, sumP := 0.0, 0.0
	for _, s := range states {
		sumV += s.V
		sumP += s.P
	}
	n := float64(len(states))
	return sumV / n, sumP / n
}

func PopulationStddev(states []lv.State) (float64, float64) {
	if len(states) < 2 {
		return 0, 0
	}
	meanV, meanP := MeanPopulation(states)
	var vv, pp float64
	for _, s := range states {
		dv := s.V - meanV
		dp := s.P - meanP
		vv += dv * dv
		pp += dp * dp
	}
	n := float64(len(states) - 1)
	return math.Sqrt(vv / n), math.Sqrt(pp / n)
}

func PhaseArea(states []lv.State) float64 {
	area := 0.0
	n := len(states)
	if n < 3 {
		return 0
	}
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		area += states[i].V*states[j].P - states[j].V*states[i].P
	}
	return math.Abs(area) / 2
}

func IsCyclic(states []lv.State) bool {
	if len(states) < 4 {
		return false
	}
	first := states[0]
	last := states[len(states)-1]
	return math.Abs(first.V-last.V) < 1e-6 && math.Abs(first.P-last.P) < 1e-6
}

func Summarize(p lv.Params, states []lv.State) string {
	meanV, meanP := MeanPopulation(states)
	return fmt.Sprintf("n=%d mean V=%.4g mean P=%.4g area=%.4g",
		len(states), meanV, meanP, PhaseArea(states))
}

func DescribeTransition(a, b lv.State) string {
	return fmt.Sprintf("%s -> %s", a, b)
}

func Delta(a, b lv.State) lv.State {
	return lv.State{V: b.V - a.V, P: b.P - a.P}
}

func Magnitude(s lv.State) float64 {
	return math.Sqrt(s.V*s.V + s.P*s.P)
}

func Scale(s lv.State, factor float64) lv.State {
	return lv.State{V: s.V * factor, P: s.P * factor}
}

func Normalize(s lv.State) lv.State {
	mag := Magnitude(s)
	if mag == 0 {
		return lv.State{}
	}
	return Scale(s, 1/mag)
}

func Add(a, b lv.State) lv.State {
	return lv.State{V: a.V + b.V, P: a.P + b.P}
}

func IsBounded(states []lv.State, max float64) bool {
	for _, s := range states {
		if s.V > max || s.P > max {
			return false
		}
	}
	return true
}

func AboveCarrying(s lv.State, threshold float64) bool {
	return TotalPopulation(s) > threshold
}

func BelowThreshold(s lv.State, threshold float64) bool {
	return TotalPopulation(s) < threshold
}

func ExtinctionRisk(s lv.State, threshold float64) bool {
	return s.V < threshold || s.P < threshold
}
