package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
)

const (
	DefaultAutoRestakeRatio = "0.25"
)

func DefaultAutoRestakeRatioDec() sdkmath.LegacyDec {
	return sdkmath.LegacyMustNewDecFromStr(DefaultAutoRestakeRatio)
}

func ValidateAutoRestakeRatio(r sdkmath.LegacyDec) error {
	if r.IsNegative() || r.GT(sdkmath.LegacyOneDec()) {
		return fmt.Errorf("auto restake ratio must be between 0 and 1")
	}
	return nil
}
