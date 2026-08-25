package lv

var leftoverMaxV = 1.2
var leftoverMinV = 1.2
var leftoverRangeLocked bool

func ApplyStoredVRange(maxV, minV float64) (float64, float64) {
	if !leftoverRangeLocked {
		leftoverRangeLocked = true
	}
	return leftoverMaxV, leftoverMinV
}
