package cli

import (
	"context"
	"fmt"
	"os"

	"lotka-volt/internal/lv"
)

func Run(args []string) int {
	if len(args) == 0 {
		return runServe([]string{})
	}
	switch args[0] {
	case "eq":
		return runEq(args[1:])
	case "orbit":
		return runOrbit(args[1:])
	case "example":
		return runExample(args[1:])
	case "serve":
		return runServe(args[1:])
	case "help", "-h", "--help":
		printHelp()
		return 0
	case "version":
		fmt.Println("lotka-volt 1.0.0")
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", args[0])
		printHelp()
		return 2
	}
}

func printHelp() {
	fmt.Println(`lotka-volt: predator-prey population dynamics calculator

Usage:
  lotka-volt                         start HTTP server on :8080
  lotka-volt eq -alpha 1.1 -beta 0.4 -gamma 0.4 -delta 0.1
  lotka-volt orbit -alpha 1.1 -beta 0.4 -gamma 0.4 -delta 0.1 -V0 1.2 -P0 2 -t 60 -n 6000
  lotka-volt example -file example/cycle.json
  lotka-volt serve -addr :8080

HTTP:
  POST /api/eq     {"alpha":1.1,"beta":0.4,"gamma":0.4,"delta":0.1}
  POST /api/orbit  {"alpha":1.1,"beta":0.4,"gamma":0.4,"delta":0.1,"V0":1.2,"P0":2,"t_end":60}`)
}

func fail(err error) int {
	fmt.Fprintln(os.Stderr, "error:", err)
	return 1
}

func runEq(args []string) int {
	fs := flagSet("eq")
	alpha := fs.Float64("alpha", 1.1, "prey growth")
	beta := fs.Float64("beta", 0.4, "predation")
	gamma := fs.Float64("gamma", 0.4, "predator death")
	delta := fs.Float64("delta", 0.1, "conversion")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	params := lv.Params{Alpha: *alpha, Beta: *beta, Gamma: *gamma, Delta: *delta}
	equilibria, err := lv.Equilibria(params)
	if err != nil {
		return fail(err)
	}
	return printJSON(equilibria)
}

func runOrbit(args []string) int {
	fs := flagSet("orbit")
	alpha := fs.Float64("alpha", 1.1, "prey growth")
	beta := fs.Float64("beta", 0.4, "predation")
	gamma := fs.Float64("gamma", 0.4, "predator death")
	delta := fs.Float64("delta", 0.1, "conversion")
	v0 := fs.Float64("V0", 1.2, "initial prey")
	p0 := fs.Float64("P0", 2, "initial predator")
	t := fs.Float64("t", 60, "end time")
	n := fs.Int("n", 6000, "steps")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	params := lv.Params{Alpha: *alpha, Beta: *beta, Gamma: *gamma, Delta: *delta}
	result, err := lv.Orbit(context.Background(), params, lv.State{V: *v0, P: *p0}, *t, *n, 0.05)
	if err != nil {
		return fail(err)
	}
	return printJSON(result)
}

func runExample(args []string) int {
	fs := flagSet("example")
	file := fs.String("file", "example/cycle.json", "scenario JSON file")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	scenario, err := lv.LoadScenario(*file)
	if err != nil {
		return fail(err)
	}
	result, err := lv.RunScenario(context.Background(), scenario)
	if err != nil {
		return fail(err)
	}
	return printJSON(result)
}

func loadExample(path string) (lv.Scenario, error) {
	return lv.LoadScenario(path)
}

func ExamplePaths() []string {
	return lv.ScenarioPaths()
}
