package lv

import (
	"context"
	"fmt"
	"math"
)

func Derivative(p Params, s State) (State, error) {
	return DerivationRate(p, s)
}

func StepRK4(p Params, s State, dt float64) (State, error) {
	if dt <= 0 {
		return State{}, fmt.Errorf("dt must be positive")
	}
	k1, err := Derivative(p, s)
	if err != nil {
		return State{}, err
	}
	k2, err := Derivative(p, State{V: s.V + 0.5*dt*k1.V, P: s.P + 0.5*dt*k1.P})
	if err != nil {
		return State{}, err
	}
	k3, err := Derivative(p, State{V: s.V + 0.5*dt*k2.V, P: s.P + 0.5*dt*k2.P})
	if err != nil {
		return State{}, err
	}
	k4, err := Derivative(p, State{V: s.V + dt*k3.V, P: s.P + dt*k3.P})
	if err != nil {
		return State{}, err
	}
	return State{
		V: s.V + dt/6*(k1.V+2*k2.V+2*k3.V+k4.V),
		P: s.P + dt/6*(k1.P+2*k2.P+2*k3.P+k4.P),
	}, nil
}

func Integrate(ctx context.Context, p Params, start State, tEnd float64, steps int) ([]State, []float64, error) {
	if err := p.Validate(); err != nil {
		return nil, nil, err
	}
	if err := ValidateState(start); err != nil {
		return nil, nil, err
	}
	if tEnd < 0 {
		return nil, nil, fmt.Errorf("tEnd must be non-negative")
	}
	if steps <= 0 {
		return nil, nil, fmt.Errorf("steps must be positive")
	}
	dt := tEnd / float64(steps)
	states := make([]State, 0, steps+1)
	times := make([]float64, 0, steps+1)
	states = append(states, start)
	times = append(times, 0)
	s := start
	for i := 1; i <= steps; i++ {
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		default:
		}
		next, err := StepRK4(p, s, dt)
		if err != nil {
			return nil, nil, err
		}
		next, err = ConservationCorrection(p, s, next, 1e-10)
		if err != nil {
			return nil, nil, err
		}
		if next.V < 0 {
			next.V = 0
		}
		if next.P < 0 {
			next.P = 0
		}
		states = append(states, next)
		times = append(times, float64(i)*dt)
		s = next
	}
	return states, times, nil
}

func ConservationCorrection(p Params, start, next State, tolerance float64) (State, error) {
	if next.V == 0 || next.P == 0 {
		return next, nil
	}
	target, err := H(p, start)
	if err != nil {
		return State{}, err
	}
	k := 1.0
	for i := 0; i < 60; i++ {
		current, err := H(p, State{V: next.V * k, P: next.P * k})
		if err != nil {
			return State{}, err
		}
		residual := current - target
		if math.Abs(residual) <= tolerance*math.Max(1, math.Abs(target)) {
			return State{V: next.V * k, P: next.P * k}, nil
		}
		derivative := (p.Delta+p.Beta)/k - (p.Gamma*next.V + p.Alpha*next.P)
		if math.Abs(derivative) < 1e-15 {
			break
		}
		k -= residual / derivative
		if k <= 0 {
			k = 0.5
		}
	}
	return State{V: next.V * k, P: next.P * k}, nil
}

func Orbit(ctx context.Context, p Params, start State, tEnd float64, steps int, tolerance float64) (OrbitResult, error) {
	states, times, err := Integrate(ctx, p, start, tEnd, steps)
	if err != nil {
		return OrbitResult{}, err
	}
	hs, err := HList(p, states)
	if err != nil {
		return OrbitResult{}, err
	}
	drift, err := MaxHDrift(p, states)
	if err != nil {
		return OrbitResult{}, err
	}
	closed, err := IsClosed(p, states, tolerance)
	if err != nil {
		return OrbitResult{}, err
	}
	return OrbitResult{
		Times: times, V: statesToV(states), P: statesToP(states), H: hs,
		HDrift: drift, Closed: closed, Tolerance: tolerance,
		Start: start, Steps: steps,
	}, nil
}

func statesToV(states []State) []float64 {
	out := make([]float64, len(states))
	for i, s := range states {
		out[i] = s.V
	}
	return out
}

func statesToP(states []State) []float64 {
	out := make([]float64, len(states))
	for i, s := range states {
		out[i] = s.P
	}
	return out
}

func IsEquilibriumStationary(p Params) (bool, error) {
	eq, err := PositiveEquilibrium(p)
	if err != nil {
		return false, err
	}
	return IsStationary(p, State{V: eq.V, P: eq.P})
}

func AxisPreyZero(p Params, start State, tEnd float64, steps int) ([]State, error) {
	if start.V != 0 {
		return nil, fmt.Errorf("prey must be zero")
	}
	return IntegrateAxis(p, start, tEnd, steps)
}

func IntegrateAxis(p Params, start State, tEnd float64, steps int) ([]State, error) {
	states := make([]State, 0, steps+1)
	states = append(states, start)
	s := start
	dt := tEnd / float64(steps)
	for i := 0; i < steps; i++ {
		rate, err := Derivative(p, s)
		if err != nil {
			return nil, err
		}
		next := State{V: s.V + dt*rate.V, P: s.P + dt*rate.P}
		if next.V < 0 {
			next.V = 0
		}
		if next.P < 0 {
			next.P = 0
		}
		states = append(states, next)
		s = next
	}
	return states, nil
}

func AxisPredatorDecayExact(p Params, start State, t float64) State {
	return State{V: 0, P: start.P * math.Exp(-p.Gamma*t)}
}

func OrbitAmplitude(states []State) float64 {
	maxV, minV := 0.0, math.Inf(1)
	for _, s := range states {
		if s.V > maxV {
			maxV = s.V
		}
		if s.V < minV {
			minV = s.V
		}
	}
	return maxV - minV
}

func OrbitPeriodEstimate(times, vs []float64) float64 {
	if len(times) < 4 {
		return 0
	}
	firstMax := -1
	secondMax := -1
	for i, v := range vs {
		if i == 0 || i == len(vs)-1 {
			continue
		}
		if v > vs[i-1] && v > vs[i+1] {
			if firstMax < 0 {
				firstMax = i
			} else if secondMax < 0 {
				secondMax = i
				break
			}
		}
	}
	if firstMax < 0 || secondMax < 0 {
		return 0
	}
	return times[secondMax] - times[firstMax]
}
