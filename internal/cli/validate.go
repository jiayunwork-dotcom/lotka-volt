package cli

import (
	"fmt"
	"strconv"

	"lotka-volt/internal/lv"
)

func parseFloat(text, label string) (float64, error) {
	value, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be a number, got %q", label, text)
	}
	return value, nil
}

func validateParams(alpha, beta, gamma, delta float64) error {
	return lv.Params{Alpha: alpha, Beta: beta, Gamma: gamma, Delta: delta}.Validate()
}

func validateState(v, p float64) error {
	return lv.ValidateState(lv.State{V: v, P: p})
}

func requirePositive(value float64, label string) error {
	if value <= 0 {
		return fmt.Errorf("%s must be positive, got %g", label, value)
	}
	return nil
}

func normalize(text string) string {
	if len(text) > 0 && text[0] == '"' {
		if parsed, err := strconv.Unquote(text); err == nil {
			return parsed
		}
	}
	return text
}
