package lv

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type Scenario struct {
	Name   string  `json:"name"`
	Params Params  `json:"params"`
	Start  State   `json:"start"`
	TEnd   float64 `json:"t_end"`
	Steps  int     `json:"steps"`
}

func LoadScenario(path string) (Scenario, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Scenario{}, err
	}
	var scenario Scenario
	if err := json.Unmarshal(data, &scenario); err != nil {
		return Scenario{}, fmt.Errorf("parse %s: %w", path, err)
	}
	if err := scenario.Params.Validate(); err != nil {
		return Scenario{}, err
	}
	if err := ValidateState(scenario.Start); err != nil {
		return Scenario{}, err
	}
	if scenario.Steps <= 0 {
		return Scenario{}, fmt.Errorf("steps must be positive")
	}
	return scenario, nil
}

func RunScenario(ctx context.Context, scenario Scenario) (OrbitResult, error) {
	return Orbit(ctx, scenario.Params, scenario.Start, scenario.TEnd, scenario.Steps, 0.05)
}

func SaveScenario(path string, scenario Scenario) error {
	data, err := json.MarshalIndent(scenario, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func CycleExample() Scenario {
	return Scenario{
		Name:   "cycle",
		Params: Params{Alpha: 1.1, Beta: 0.4, Gamma: 0.4, Delta: 0.1},
		Start:  State{V: 1.2, P: 2.0},
		TEnd:   60,
		Steps:  6000,
	}
}

func Examples() []Scenario {
	return []Scenario{
		CycleExample(),
		{
			Name: "equilibrium", Params: Params{Alpha: 1.1, Beta: 0.4, Gamma: 0.4, Delta: 0.1},
			Start: State{V: 0.25, P: 2.75}, TEnd: 20, Steps: 2000,
		},
		{
			Name: "prey-axis", Params: Params{Alpha: 1.1, Beta: 0.4, Gamma: 0.4, Delta: 0.1},
			Start: State{V: 0, P: 5}, TEnd: 10, Steps: 1000,
		},
	}
}

func ScenarioPaths() []string {
	return []string{"example/cycle.json"}
}

func Describe(scenario Scenario) string {
	return fmt.Sprintf("%s: %s start=%s t=%.4g steps=%d",
		scenario.Name, scenario.Params, scenario.Start, scenario.TEnd, scenario.Steps)
}

func ExampleText() string {
	return "alpha=1.1 beta=0.4 gamma=0.4 delta=0.1, V0=1.2 P0=2.0"
}
