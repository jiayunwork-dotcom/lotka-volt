package verify

import (
	"context"
	"fmt"
	"math"

	"lotka-volt/internal/lv"
)

type Check struct {
	Name    string `json:"name"`
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func testParams() lv.Params {
	return lv.Params{Alpha: 1.1, Beta: 0.4, Gamma: 0.4, Delta: 0.1}
}

func CheckEquilibrium() Check {
	params := testParams()
	eq, err := lv.PositiveEquilibrium(params)
	if err != nil {
		return Check{Name: "equilibrium", OK: false, Message: err.Error()}
	}
	ok := math.Abs(eq.V-0.25) < 1e-12 && math.Abs(eq.P-2.75) < 1e-12
	return Check{Name: "equilibrium", OK: ok, Message: fmt.Sprintf("V*=%.6g P*=%.6g", eq.V, eq.P)}
}

func CheckEquilibriumStationary() Check {
	ok, err := lv.IsEquilibriumStationary(testParams())
	if err != nil {
		return Check{Name: "equilibrium-stationary", OK: false, Message: err.Error()}
	}
	return Check{Name: "equilibrium-stationary", OK: ok, Message: fmt.Sprintf("stationary=%v", ok)}
}

func CheckConservation() Check {
	params := testParams()
	start := lv.State{V: 1.2, P: 2.0}
	orbit, err := lv.Orbit(context.Background(), params, start, 60, 6000, 0.05)
	if err != nil {
		return Check{Name: "conservation", OK: false, Message: err.Error()}
	}
	ok := orbit.HDrift < 0.05
	return Check{Name: "conservation", OK: ok, Message: fmt.Sprintf("H_drift=%.6g", orbit.HDrift)}
}

func CheckClosed() Check {
	params := testParams()
	orbit, err := lv.Orbit(context.Background(), params, lv.State{V: 1.2, P: 2.0}, 60, 6000, 0.05)
	if err != nil {
		return Check{Name: "closed", OK: false, Message: err.Error()}
	}
	ok := orbit.Closed
	return Check{Name: "closed", OK: ok, Message: fmt.Sprintf("closed=%v", orbit.Closed)}
}

func CheckAxisPreyZero() Check {
	params := testParams()
	states, err := lv.IntegrateAxis(params, lv.State{V: 0, P: 5}, 10, 1000)
	if err != nil {
		return Check{Name: "axis-prey-zero", OK: false, Message: err.Error()}
	}
	last := states[len(states)-1]
	ok := last.V == 0 && last.P > 0 && last.P < 5
	return Check{Name: "axis-prey-zero", OK: ok, Message: fmt.Sprintf("V=%.6g P=%.6g", last.V, last.P)}
}

func CheckAxisPredatorZero() Check {
	params := testParams()
	states, err := lv.IntegrateAxis(params, lv.State{V: 3, P: 0}, 10, 1000)
	if err != nil {
		return Check{Name: "axis-predator-zero", OK: false, Message: err.Error()}
	}
	last := states[len(states)-1]
	ok := last.P == 0 && last.V > 0
	return Check{Name: "axis-predator-zero", OK: ok, Message: fmt.Sprintf("V=%.6g P=%.6g", last.V, last.P)}
}

func CheckAlphaRaisesPredatorEquilibrium() Check {
	low := testParams()
	high := testParams()
	high.Alpha = low.Alpha * 2
	eqLow, _ := lv.PositiveEquilibrium(low)
	eqHigh, _ := lv.PositiveEquilibrium(high)
	ok := eqHigh.P > eqLow.P
	return Check{Name: "alpha-predator", OK: ok, Message: fmt.Sprintf("%.6g -> %.6g", eqLow.P, eqHigh.P)}
}

func CheckValidationRejects() Check {
	err := lv.Params{Alpha: 0, Beta: 1, Gamma: 1, Delta: 1}.Validate()
	ok := err != nil
	return Check{Name: "validation", OK: ok, Message: fmt.Sprintf("err=%v", err)}
}

func RunAll() []Check {
	return []Check{
		CheckEquilibrium(),
		CheckEquilibriumStationary(),
		CheckConservation(),
		CheckClosed(),
		CheckAxisPreyZero(),
		CheckAxisPredatorZero(),
		CheckAlphaRaisesPredatorEquilibrium(),
		CheckValidationRejects(),
	}
}

func AllPass(checks []Check) bool {
	for _, check := range checks {
		if !check.OK {
			return false
		}
	}
	return true
}

func FormatChecks(checks []Check) string {
	out := ""
	for _, check := range checks {
		state := "PASS"
		if !check.OK {
			state = "FAIL"
		}
		out += fmt.Sprintf("%-22s %s %s\n", check.Name, state, check.Message)
	}
	return out
}

func CheckHZeroAtOrigin() Check {
	value, err := lv.H(testParams(), lv.State{V: 0, P: 0})
	if err != nil {
		return Check{Name: "h-origin", OK: false, Message: err.Error()}
	}
	ok := value == 0
	return Check{Name: "h-origin", OK: ok, Message: fmt.Sprintf("H=%.6g", value)}
}

func CheckEquilibriumHStable() Check {
	params := testParams()
	eq, _ := lv.PositiveEquilibrium(params)
	start := lv.State{V: eq.V + 0.1, P: eq.P + 0.1}
	orbit, err := lv.Orbit(context.Background(), params, start, 30, 3000, 0.05)
	if err != nil {
		return Check{Name: "eq-h-stable", OK: false, Message: err.Error()}
	}
	ok := orbit.HDrift < 0.05
	return Check{Name: "eq-h-stable", OK: ok, Message: fmt.Sprintf("drift=%.6g", orbit.HDrift)}
}

func CheckDeltaPreyEquilibrium() Check {
	low := testParams()
	high := testParams()
	high.Delta = low.Delta * 2
	eqLow, _ := lv.PositiveEquilibrium(low)
	eqHigh, _ := lv.PositiveEquilibrium(high)
	ok := eqHigh.V > eqLow.V
	return Check{Name: "delta-prey", OK: ok, Message: fmt.Sprintf("%.6g -> %.6g", eqLow.V, eqHigh.V)}
}

func CheckBetaLowerPredator() Check {
	low := testParams()
	high := testParams()
	high.Beta = low.Beta * 2
	eqLow, _ := lv.PositiveEquilibrium(low)
	eqHigh, _ := lv.PositiveEquilibrium(high)
	ok := eqHigh.P < eqLow.P
	return Check{Name: "beta-predator", OK: ok, Message: fmt.Sprintf("%.6g -> %.6g", eqLow.P, eqHigh.P)}
}

func CheckGammaLowerPrey() Check {
	low := testParams()
	high := testParams()
	high.Gamma = low.Gamma * 2
	eqLow, _ := lv.PositiveEquilibrium(low)
	eqHigh, _ := lv.PositiveEquilibrium(high)
	ok := eqHigh.V < eqLow.V
	return Check{Name: "gamma-prey", OK: ok, Message: fmt.Sprintf("%.6g -> %.6g", eqLow.V, eqHigh.V)}
}

func CheckStateValidation() Check {
	err := lv.ValidateState(lv.State{V: -1, P: 1})
	ok := err != nil
	return Check{Name: "state-validation", OK: ok, Message: fmt.Sprintf("err=%v", err)}
}

func CheckHPeakToPeakSmall() Check {
	params := testParams()
	orbit, _ := lv.Orbit(context.Background(), params, lv.State{V: 1.2, P: 2.0}, 60, 6000, 0.05)
	states := make([]lv.State, len(orbit.V))
	for i := range orbit.V {
		states[i] = lv.State{V: orbit.V[i], P: orbit.P[i]}
	}
	spread, err := lv.HPeakToPeak(params, states)
	if err != nil {
		return Check{Name: "h-peak", OK: false, Message: err.Error()}
	}
	ok := spread < 0.1
	return Check{Name: "h-peak", OK: ok, Message: fmt.Sprintf("spread=%.6g", spread)}
}
