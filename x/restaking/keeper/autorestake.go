package keeper

import (
    "context"

    sdk "github.com/cosmos/cosmos-sdk/types"
    stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// RestakeDelegate delegates the provided portion back to the validator for the delegator.
func (k Keeper) RestakeDelegate(ctx sdk.Context, delegator sdk.AccAddress, validator sdk.ValAddress, portion sdk.Coins) error {
	if portion.IsZero() {
		return nil
	}

	bondDenom, err := k.stakingKeeper.BondDenom(contextWithSDK(ctx))
	if err != nil {
		return err
	}

	amount := portion.AmountOf(bondDenom)
	if amount.IsZero() {
		return nil
	}

	val, err := k.stakingKeeper.GetValidator(contextWithSDK(ctx), validator)
	if err != nil {
		return err
	}

	_, err = k.stakingKeeper.Delegate(
		contextWithSDK(ctx),
		delegator,
		amount,
		stakingtypes.Unbonded,
		val,
		true,
	)
	return err
}

func contextWithSDK(ctx sdk.Context) context.Context {
	return sdk.WrapSDKContext(ctx)
}
