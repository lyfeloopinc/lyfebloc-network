package types

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// StakingKeeper defines the subset of staking keeper functionality required by restaking.
type StakingKeeper interface {
	GetValidator(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error)
	Delegate(ctx context.Context, delAddr sdk.AccAddress, bondAmt sdkmath.Int, status stakingtypes.BondStatus, validator stakingtypes.Validator, subtractAccount bool) (sdkmath.LegacyDec, error)
	BondDenom(ctx context.Context) (string, error)
}
