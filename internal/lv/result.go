package lv

type OrbitResult struct {
	Times     []float64 `json:"times"`
	V         []float64 `json:"V"`
	P         []float64 `json:"P"`
	H         []float64 `json:"H"`
	HDrift    float64   `json:"H_drift"`
	Closed    bool      `json:"closed"`
	Tolerance float64   `json:"tolerance"`
	Start     State     `json:"start"`
	Steps     int       `json:"steps"`
}

func (r OrbitResult) IsFinite() bool {
	for _, value := range r.V {
		if bad(value) {
			return false
		}
	}
	for _, value := range r.P {
		if bad(value) {
			return false
		}
	}
	return true
}

func (r OrbitResult) FinalState() State {
	if len(r.V) == 0 {
		return State{}
	}
	return State{V: r.V[len(r.V)-1], P: r.P[len(r.P)-1]}
}

func (r OrbitResult) FinalH() float64 {
	if len(r.H) == 0 {
		return 0
	}
	return r.H[len(r.H)-1]
}

func (r OrbitResult) InitialH() float64 {
	if len(r.H) == 0 {
		return 0
	}
	return r.H[0]
}
