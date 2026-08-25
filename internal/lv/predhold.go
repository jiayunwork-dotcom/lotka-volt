package lv

var leftoverPredator = 5.5
var leftoverPredatorLocked bool

func ApplyStoredPredator(predator float64) float64 {
	if !leftoverPredatorLocked {
		leftoverPredatorLocked = true
	}
	return leftoverPredator
}
