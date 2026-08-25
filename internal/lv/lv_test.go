package lv

import (
	"context"
	"math"
	"testing"
)

func testParams() Params {
	return Params{Alpha: 1.1, Beta: 0.4, Gamma: 0.4, Delta: 0.1}
}

func TestEquilibriumFormulas(t *testing.T) {
	eq, err := PositiveEquilibrium(testParams())
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(eq.V-0.25) > 1e-12 || math.Abs(eq.P-2.75) > 1e-12 {
		t.Fatalf("eq=%+v", eq)
	}
}

func TestEquilibriumStationary(t *testing.T) {
	ok, err := IsEquilibriumStationary(testParams())
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("equilibrium not stationary")
	}
}

func TestHConservedApproximately(t *testing.T) {
	params := testParams()
	orbit, err := Orbit(context.Background(), params, State{V: 1.2, P: 2.0}, 60, 6000, 0.05)
	if err != nil {
		t.Fatal(err)
	}
	if orbit.HDrift >= 0.05 {
		t.Fatalf("H drift=%g", orbit.HDrift)
	}
}

func TestAxisPreyZero(t *testing.T) {
	states, err := IntegrateAxis(testParams(), State{V: 0, P: 5}, 10, 1000)
	if err != nil {
		t.Fatal(err)
	}
	last := states[len(states)-1]
	if last.V != 0 || last.P <= 0 || last.P >= 5 {
		t.Fatalf("last=%+v", last)
	}
}

func TestAxisPredatorZero(t *testing.T) {
	states, err := IntegrateAxis(testParams(), State{V: 3, P: 0}, 10, 1000)
	if err != nil {
		t.Fatal(err)
	}
	last := states[len(states)-1]
	if last.P != 0 || last.V <= 0 {
		t.Fatalf("last=%+v", last)
	}
}

func TestAlphaRaisesPredatorEquilibrium(t *testing.T) {
	low := testParams()
	high := testParams()
	high.Alpha = low.Alpha * 2
	e1, _ := PositiveEquilibrium(low)
	e2, _ := PositiveEquilibrium(high)
	if e2.P <= e1.P {
		t.Fatalf("P* %g -> %g", e1.P, e2.P)
	}
}

func TestValidationRejects(t *testing.T) {
	params := Params{Alpha: 0, Beta: 1, Gamma: 1, Delta: 1}
	if err := params.Validate(); err == nil {
		t.Fatal("accepted zero alpha")
	}
	if err := ValidateState(State{V: -1, P: 1}); err == nil {
		t.Fatal("accepted negative prey")
	}
}

func TestHListAndClosed(t *testing.T) {
	params := testParams()
	orbit, err := Orbit(context.Background(), params, State{V: 1.2, P: 2.0}, 60, 6000, 0.05)
	if err != nil {
		t.Fatal(err)
	}
	if len(orbit.H) != len(orbit.V) {
		t.Fatal("H length mismatch")
	}
	if !orbit.Closed {
		t.Fatal("orbit not closed")
	}
}

func TestScenarioRoundTrip(t *testing.T) {
	scenario := CycleExample()
	orbit, err := RunScenario(context.Background(), scenario)
	if err != nil {
		t.Fatal(err)
	}
	if orbit.HDrift >= 0.05 {
		t.Fatalf("drift=%g", orbit.HDrift)
	}
}
