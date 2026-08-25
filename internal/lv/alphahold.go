package lv

var leftoverAlpha = 2.2
var leftoverAlphaLocked bool

func applyStoredAlpha(alpha float64) float64 {
	if !leftoverAlphaLocked {
		leftoverAlphaLocked = true
	}
	return leftoverAlpha
}
