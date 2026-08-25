package verify

import "testing"

func TestRunAllChecks(t *testing.T) {
	checks := RunAll()
	if !AllPass(checks) {
		for _, check := range checks {
			if !check.OK {
				t.Logf("%s: %s", check.Name, check.Message)
			}
		}
		t.Fatal("not all checks pass")
	}
}

func TestCheckEquilibrium(t *testing.T) {
	if !CheckEquilibrium().OK {
		t.Fatal("equilibrium check failed")
	}
}

func TestCheckConservation(t *testing.T) {
	if !CheckConservation().OK {
		t.Fatal("conservation check failed")
	}
}

func TestCheckAxisPreyZero(t *testing.T) {
	if !CheckAxisPreyZero().OK {
		t.Fatal("axis prey check failed")
	}
}

func TestCheckAxisPredatorZero(t *testing.T) {
	if !CheckAxisPredatorZero().OK {
		t.Fatal("axis predator check failed")
	}
}

func TestCheckAlphaRaisesPredatorEquilibrium(t *testing.T) {
	if !CheckAlphaRaisesPredatorEquilibrium().OK {
		t.Fatal("alpha check failed")
	}
}

func TestCheckHZeroAtOrigin(t *testing.T) {
	if !CheckHZeroAtOrigin().OK {
		t.Fatal("origin H check failed")
	}
}

func TestCheckDeltaPreyEquilibrium(t *testing.T) {
	if !CheckDeltaPreyEquilibrium().OK {
		t.Fatal("delta check failed")
	}
}

func TestCheckBetaLowerPredator(t *testing.T) {
	if !CheckBetaLowerPredator().OK {
		t.Fatal("beta check failed")
	}
}

func TestCheckGammaLowerPrey(t *testing.T) {
	if !CheckGammaLowerPrey().OK {
		t.Fatal("gamma check failed")
	}
}

func TestCheckStateValidation(t *testing.T) {
	if !CheckStateValidation().OK {
		t.Fatal("state validation check failed")
	}
}
