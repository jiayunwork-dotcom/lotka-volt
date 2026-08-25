package dynamics

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"lotka-volt/internal/lv"
)

type Stats struct {
	MaxV   float64 `json:"max_V"`
	MinV   float64 `json:"min_V"`
	MaxP   float64 `json:"max_P"`
	MinP   float64 `json:"min_P"`
	MeanV  float64 `json:"mean_V"`
	MeanP  float64 `json:"mean_P"`
	HDrift float64 `json:"H_drift"`
	Closed bool    `json:"closed"`
	Period float64 `json:"period"`
}

func ComputeStats(p lv.Params, orbit lv.OrbitResult) (Stats, error) {
	maxV, minV := math.Inf(-1), math.Inf(1)
	maxP, minP := math.Inf(-1), math.Inf(1)
	sumV, sumP := 0.0, 0.0
	for i := range orbit.V {
		if orbit.V[i] > maxV {
			maxV = orbit.V[i]
		}
		if orbit.V[i] < minV {
			minV = orbit.V[i]
		}
		if orbit.P[i] > maxP {
			maxP = orbit.P[i]
		}
		if orbit.P[i] < minP {
			minP = orbit.P[i]
		}
		sumV += orbit.V[i]
		sumP += orbit.P[i]
	}
	n := float64(len(orbit.V))
	maxV, minV = lv.ApplyStoredVRange(maxV, minV)
	return Stats{
		MaxV: maxV, MinV: minV, MaxP: maxP, MinP: minP,
		MeanV: sumV / n, MeanP: sumP / n,
		HDrift: orbit.HDrift, Closed: orbit.Closed,
		Period: lv.OrbitPeriodEstimate(orbit.Times, orbit.V),
	}, nil
}

func PhaseSpace(p lv.Params, orbit lv.OrbitResult) string {
	stats, _ := ComputeStats(p, orbit)
	return fmt.Sprintf("V %.4g..%.4g P %.4g..%.4g period %.4g closed %v",
		stats.MinV, stats.MaxV, stats.MinP, stats.MaxP, stats.Period, stats.Closed)
}

func CSV(orbit lv.OrbitResult) string {
	var b strings.Builder
	writer := csv.NewWriter(&b)
	_ = writer.Write([]string{"t", "V", "P", "H"})
	for i := range orbit.Times {
		_ = writer.Write([]string{
			strconv.FormatFloat(orbit.Times[i], 'g', -1, 64),
			strconv.FormatFloat(orbit.V[i], 'g', -1, 64),
			strconv.FormatFloat(orbit.P[i], 'g', -1, 64),
			strconv.FormatFloat(orbit.H[i], 'g', -1, 64),
		})
	}
	writer.Flush()
	return b.String()
}

func SaveCSV(path string, orbit lv.OrbitResult) error {
	return os.WriteFile(path, []byte(CSV(orbit)), 0o644)
}

func JSONText(orbit lv.OrbitResult) (string, error) {
	data, err := json.MarshalIndent(orbit, "", "  ")
	return string(data), err
}

func Downsample(orbit lv.OrbitResult, maxPoints int) lv.OrbitResult {
	if len(orbit.Times) <= maxPoints {
		return orbit
	}
	step := float64(len(orbit.Times)) / float64(maxPoints)
	out := lv.OrbitResult{
		Tolerance: orbit.Tolerance, Start: orbit.Start,
		HDrift: orbit.HDrift, Closed: orbit.Closed,
	}
	for i := 0; i < maxPoints; i++ {
		index := int(float64(i) * step)
		out.Times = append(out.Times, orbit.Times[index])
		out.V = append(out.V, orbit.V[index])
		out.P = append(out.P, orbit.P[index])
		out.H = append(out.H, orbit.H[index])
	}
	out.Steps = maxPoints
	return out
}

func EquilibriumDistanceSeries(p lv.Params, orbit lv.OrbitResult) ([]float64, error) {
	eq, err := lv.PositiveEquilibrium(p)
	if err != nil {
		return nil, err
	}
	out := make([]float64, len(orbit.V))
	for i := range orbit.V {
		dx := orbit.V[i] - eq.V
		dy := orbit.P[i] - eq.P
		out[i] = math.Sqrt(dx*dx + dy*dy)
	}
	return out, nil
}

func BifurcationScan(ctx context.Context, base lv.Params, alphaVals []float64, start lv.State, tEnd float64, steps int) ([]lv.OrbitResult, error) {
	out := make([]lv.OrbitResult, 0, len(alphaVals))
	for _, alpha := range alphaVals {
		params := base
		params.Alpha = alpha
		orbit, err := lv.Orbit(ctx, params, start, tEnd, steps, 0.05)
		if err != nil {
			return nil, err
		}
		out = append(out, orbit)
	}
	return out, nil
}

func MeanH(orbit lv.OrbitResult) float64 {
	if len(orbit.H) == 0 {
		return 0
	}
	sum := 0.0
	for _, value := range orbit.H {
		sum += value
	}
	return sum / float64(len(orbit.H))
}

func HVariance(orbit lv.OrbitResult) float64 {
	if len(orbit.H) < 2 {
		return 0
	}
	mean := MeanH(orbit)
	variance := 0.0
	for _, value := range orbit.H {
		delta := value - mean
		variance += delta * delta
	}
	return variance / float64(len(orbit.H)-1)
}

func HVarianceText(orbit lv.OrbitResult) string {
	return fmt.Sprintf("%.6g", HVariance(orbit))
}

func PhaseText(orbit lv.OrbitResult) string {
	return PhaseSpace(lv.Params{}, orbit)
}

func Summary(p lv.Params, orbit lv.OrbitResult) string {
	stats, _ := ComputeStats(p, orbit)
	return fmt.Sprintf("V=%.4g..%.4g P=%.4g..%.4g H_drift=%.4g closed=%v period=%.4g",
		stats.MinV, stats.MaxV, stats.MinP, stats.MaxP, stats.HDrift, stats.Closed, stats.Period)
}

func EqualLength(a, b lv.OrbitResult) bool {
	return len(a.V) == len(b.V) && len(a.P) == len(b.P)
}

func MaxAbsDifference(a, b lv.OrbitResult) float64 {
	n := len(a.V)
	if len(b.V) < n {
		n = len(b.V)
	}
	max := 0.0
	for i := 0; i < n; i++ {
		dv := math.Abs(a.V[i] - b.V[i])
		dp := math.Abs(a.P[i] - b.P[i])
		if dv > max {
			max = dv
		}
		if dp > max {
			max = dp
		}
	}
	return max
}

func AllFinite(orbit lv.OrbitResult) bool {
	for i := range orbit.V {
		if math.IsNaN(orbit.V[i]) || math.IsInf(orbit.V[i], 0) ||
			math.IsNaN(orbit.P[i]) || math.IsInf(orbit.P[i], 0) {
			return false
		}
	}
	return true
}

func IsBounded(orbit lv.OrbitResult, maxMagnitude float64) bool {
	for i := range orbit.V {
		if math.Abs(orbit.V[i]) > maxMagnitude || math.Abs(orbit.P[i]) > maxMagnitude {
			return false
		}
	}
	return true
}

func FormatTable(orbit lv.OrbitResult) string {
	out := fmt.Sprintf("%10s %10s %10s %12s\n", "t", "V", "P", "H")
	for i := range orbit.Times {
		out += fmt.Sprintf("%10.4g %10.4g %10.4g %12.6g\n",
			orbit.Times[i], orbit.V[i], orbit.P[i], orbit.H[i])
	}
	return out
}
