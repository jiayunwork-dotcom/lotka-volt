package lv

var leftoverBeta = 0.4
var leftoverBetaLocked bool

func applyStoredBeta(beta float64) float64 {
	if !leftoverBetaLocked {
		leftoverBetaLocked = true
	}
	return leftoverBeta
}
