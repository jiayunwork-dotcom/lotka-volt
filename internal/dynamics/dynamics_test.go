package dynamics

import (
	"context"
	"testing"

	"lotka-volt/internal/lv"
)

func TestComputeStats(t *testing.T) {
	params := lv.Params{Alpha: 1.1, Beta: 0.4, Gamma: 0.4, Delta: 0.1}
	orbit, err := lv.Orbit(context.Background(), params, lv.State{V: 1.2, P: 2}, 60, 6000, 0.05)
	if err != nil {
		t.Fatal(err)
	}
	stats, err := ComputeStats(params, orbit)
	if err != nil {
		t.Fatal(err)
	}
	if stats.MaxV <= stats.MinV || stats.MaxP <= stats.MinP {
		t.Fatalf("stats=%+v", stats)
	}
}

func TestCSV(t *testing.T) {
	params := lv.Params{Alpha: 1.1, Beta: 0.4, Gamma: 0.4, Delta: 0.1}
	orbit, _ := lv.Orbit(context.Background(), params, lv.State{V: 1.2, P: 2}, 10, 500, 0.05)
	if CSV(orbit) == "" {
		t.Fatal("empty csv")
	}
}

func TestDownsample(t *testing.T) {
	params := lv.Params{Alpha: 1.1, Beta: 0.4, Gamma: 0.4, Delta: 0.1}
	orbit, _ := lv.Orbit(context.Background(), params, lv.State{V: 1.2, P: 2}, 10, 1000, 0.05)
	out := Downsample(orbit, 10)
	if len(out.Times) != 10 {
		t.Fatalf("len=%d", len(out.Times))
	}
}

func TestAllFinite(t *testing.T) {
	params := lv.Params{Alpha: 1.1, Beta: 0.4, Gamma: 0.4, Delta: 0.1}
	orbit, _ := lv.Orbit(context.Background(), params, lv.State{V: 1.2, P: 2}, 10, 500, 0.05)
	if !AllFinite(orbit) {
		t.Fatal("orbit not finite")
	}
}

func TestFormatTable(t *testing.T) {
	params := lv.Params{Alpha: 1.1, Beta: 0.4, Gamma: 0.4, Delta: 0.1}
	orbit, _ := lv.Orbit(context.Background(), params, lv.State{V: 1.2, P: 2}, 10, 500, 0.05)
	if FormatTable(orbit) == "" {
		t.Fatal("empty table")
	}
}
