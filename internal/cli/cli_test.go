package cli

import (
	"context"
	"testing"

	"lotka-volt/internal/lv"
)

func TestValidateParams(t *testing.T) {
	if err := validateParams(1.1, 0.4, 0.4, 0.1); err != nil {
		t.Fatal(err)
	}
	if err := validateParams(0, 0.4, 0.4, 0.1); err == nil {
		t.Fatal("accepted zero alpha")
	}
}

func TestValidateState(t *testing.T) {
	if err := validateState(1.2, 2); err != nil {
		t.Fatal(err)
	}
	if err := validateState(-1, 2); err == nil {
		t.Fatal("accepted negative prey")
	}
}

func TestLoadExample(t *testing.T) {
	scenario, err := loadExample("../../example/cycle.json")
	if err != nil {
		t.Fatal(err)
	}
	if scenario.Params.Alpha != 1.1 {
		t.Fatalf("scenario=%+v", scenario)
	}
}

func TestExamplePaths(t *testing.T) {
	if len(ExamplePaths()) != 1 {
		t.Fatal("expected one example path")
	}
}

func TestEvaluateExample(t *testing.T) {
	scenario, _ := loadExample("../../example/cycle.json")
	orbit, err := lv.RunScenario(context.Background(), scenario)
	if err != nil {
		t.Fatal(err)
	}
	if orbit.HDrift >= 0.05 {
		t.Fatalf("drift=%g", orbit.HDrift)
	}
}
