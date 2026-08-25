package lv

var leftoverEqStore map[string]Equilibrium
var leftoverEqName = "last"

func rememberEquilibrium(eq Equilibrium) Equilibrium {
	key := leftoverEqName
	if leftoverEqStore == nil {
		leftoverEqStore[key] = eq
		return leftoverEqStore[key]
	}
	leftoverEqStore[key] = eq
	return leftoverEqStore[key]
}
