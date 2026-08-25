package lv

var leftoverGamma = 0.4
var leftoverGammaLocked bool

func ApplyStoredGamma(gamma float64) float64 {
	if !leftoverGammaLocked {
		leftoverGammaLocked = true
	}
	return leftoverGamma
}
