package pop

import (
	"testing"

	"lotka-volt/internal/lv"
)

func TestPerCapitaRates(t *testing.T) {
	params := lv.Params{Alpha: 1.1, Beta: 0.4, Gamma: 0.4, Delta: 0.1}
	s := lv.State{V: 1.2, P: 2.0}
	if PreyPerCapita(params, s) <= 0 {
		t.Fatal("prey per capita non-positive")
	}
	if PredatorPerCapita(params, s) <= 0 {
		t.Fatal("predator per capita non-positive")
	}
}

func TestDoublingAndHalfLife(t *testing.T) {
	if DoublingTime(0.1) <= 0 {
		t.Fatal("doubling time invalid")
	}
	if HalfLife(-0.1) <= 0 {
		t.Fatal("half life invalid")
	}
}

func TestPopulationHelpers(t *testing.T) {
	s := lv.State{V: 2, P: 3}
	if TotalPopulation(s) != 5 {
		t.Fatal("total wrong")
	}
	if PredatorFraction(s) != 0.6 {
		t.Fatal("fraction wrong")
	}
}

func TestPhaseLabel(t *testing.T) {
	if PhaseLabel(lv.State{V: 0, P: 0}) != "origin" {
		t.Fatal("origin label wrong")
	}
	if PhaseLabel(lv.State{V: 1, P: 0}) != "prey-axis" {
		t.Fatal("prey axis label wrong")
	}
	if PhaseLabel(lv.State{V: 1, P: 1}) != "interior" {
		t.Fatal("interior label wrong")
	}
}

func TestTrajectoryLength(t *testing.T) {
	states := []lv.State{{V: 0, P: 0}, {V: 3, P: 4}}
	if TrajectoryLength(states) != 5 {
		t.Fatalf("length=%g", TrajectoryLength(states))
	}
}

func TestMeanPopulation(t *testing.T) {
	states := []lv.State{{V: 1, P: 2}, {V: 3, P: 4}}
	meanV, meanP := MeanPopulation(states)
	if meanV != 2 || meanP != 3 {
		t.Fatalf("mean %g %g", meanV, meanP)
	}
}

func TestPhaseArea(t *testing.T) {
	states := []lv.State{{V: 0, P: 0}, {V: 1, P: 0}, {V: 1, P: 1}, {V: 0, P: 1}}
	if PhaseArea(states) != 1 {
		t.Fatalf("area=%g", PhaseArea(states))
	}
}
