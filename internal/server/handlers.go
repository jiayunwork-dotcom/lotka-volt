package server

import (
	"net/http"

	"lotka-volt/internal/lv"
)

type eqRequest struct {
	Alpha float64 `json:"alpha"`
	Beta  float64 `json:"beta"`
	Gamma float64 `json:"gamma"`
	Delta float64 `json:"delta"`
}

type orbitRequest struct {
	eqRequest
	V0    float64 `json:"V0"`
	P0    float64 `json:"P0"`
	TEnd  float64 `json:"t_end"`
	Steps int     `json:"steps"`
}

func eqHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}
	var req eqRequest
	if !readJSON(w, r, &req) {
		return
	}
	params := lv.Params{Alpha: req.Alpha, Beta: req.Beta, Gamma: req.Gamma, Delta: req.Delta}
	equilibria, err := lv.Equilibria(params)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"equilibria": equilibria,
		"positive":   equilibria[1],
	})
}

func orbitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w, "POST")
		return
	}
	var req orbitRequest
	if !readJSON(w, r, &req) {
		return
	}
	if req.Steps <= 0 {
		req.Steps = 6000
	}
	params := lv.Params{Alpha: req.Alpha, Beta: req.Beta, Gamma: req.Gamma, Delta: req.Delta}
	result, err := lv.Orbit(r.Context(), params, lv.State{V: req.V0, P: req.P0}, req.TEnd, req.Steps, 0.05)
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w, "GET")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"name": "lotka-volt", "version": "1.0.0"})
}
