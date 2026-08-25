# lotka-volt

lotka-volt is a Go Lotka-Volterra predator-prey calculator. It computes the
positive equilibrium `(V*, P*)`, integrates the nonlinear system with RK4 plus
a conservation projection that keeps the Hamiltonian-like invariant
`H = delta*ln(V) - gamma*V + beta*ln(P) - alpha*P` within a small drift bound,
and reports axis behavior when either species starts at zero. The service is
available through HTTP JSON endpoints and CLI subcommands with no web page.

## Usage

Run the HTTP server:

```bash
go run . serve -addr :8080
```

Evaluate from the command line:

```bash
go run . eq -alpha 1.1 -beta 0.4 -gamma 0.4 -delta 0.1
go run . orbit -alpha 1.1 -beta 0.4 -gamma 0.4 -delta 0.1 -V0 1.2 -P0 2 -t 60 -n 6000
```

Run the cycle example:

```bash
go run . example -file example/cycle.json
```

The example starts off equilibrium and remains on a closed, approximately
conservative orbit.

## HTTP API

```text
POST /api/eq     {"alpha":1.1,"beta":0.4,"gamma":0.4,"delta":0.1}
POST /api/orbit  {"alpha":1.1,"beta":0.4,"gamma":0.4,"delta":0.1,"V0":1.2,"P0":2,"t_end":60}
GET  /health
```

Non-positive rates or negative initial populations return an error body with
HTTP 400.

## Model

The system uses prey growth `alpha`, predation `beta`, conversion `gamma`, and
predator mortality `delta`:

- `dV/dt = alpha*V - beta*V*P`
- `dP/dt = gamma*V*P - delta*P`
- `V* = delta/gamma`, `P* = alpha/beta`

## Code Layout

```text
internal/lv       equilibrium, H invariant, RK4 integration, scenarios
internal/dynamics trajectory stats, phase summaries, CSV
internal/table    orbit tables and CSV helpers
internal/pop      per-capita rates and population metrics
internal/verify   conservation/equilibrium/axis checks
internal/server   HTTP handlers and JSON responses
internal/cli      subcommand parsing and terminal output
example/          offline scenario JSON files
```

## Build and Test

```bash
export GOTOOLCHAIN=local CGO_ENABLED=0
go build ./...
go test ./...
```

The Dockerfile builds the server binary and starts it on port 8080.
