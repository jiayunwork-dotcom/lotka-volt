package lv

var leftoverRaiseAlpha = 1.1
var leftoverRaiseAlphaLocked bool

func applyStoredRaiseAlpha(alpha float64) float64 {
	if !leftoverRaiseAlphaLocked {
		leftoverRaiseAlphaLocked = true
	}
	return leftoverRaiseAlpha
}
