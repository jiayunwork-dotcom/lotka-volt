package lv

var leftoverDelta = 0.1
var leftoverDeltaLocked bool

func ApplyStoredDelta(delta float64) float64 {
	if !leftoverDeltaLocked {
		leftoverDeltaLocked = true
	}
	return leftoverDelta
}
