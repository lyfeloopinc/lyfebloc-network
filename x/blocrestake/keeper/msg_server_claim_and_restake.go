package keeper

import (
	"context"
	"errors"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/lyfeloopinc/lyfebloc-network/x/blocrestake/types"
)

// ClaimAndRestake claims rewards via the distribution module, then re-delegates them to the same validator.
func (s msgServer) ClaimAndRestake(ctx context.Context, msg *types.MsgClaimAndRestake) (*types.MsgClaimAndRestakeResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	delAddr, err := sdk.AccAddressFromBech32(msg.Delegator)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidAddress, "invalid delegator address")
	}
	valAddr, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidAddress, "invalid validator address")
	}

	// 1. Ensure validator exists
	val, err := s.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		if errors.Is(err, stakingtypes.ErrNoValidatorFound) {
			return nil, types.ErrValidatorNotFound
		}
		return nil, sdkerrors.Wrap(err, "failed to fetch validator")
	}

	// 2. Withdraw delegator’s pending rewards from the distribution module
	rewards, err := s.distributionKeeper.WithdrawDelegationRewards(ctx, delAddr, valAddr)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to withdraw rewards")
	}

	bondDenom, err := s.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to fetch bond denom")
	}

	amount := rewards.AmountOf(bondDenom)
	if amount.IsZero() {
		return nil, sdkerrors.Wrapf(types.ErrInsufficientFunds, "no %s rewards available to restake", bondDenom)
	}

	coin := sdk.NewCoin(bondDenom, amount)
	if err := s.bankKeeper.SendCoinsFromAccountToModule(ctx, delAddr, stakingtypes.BondedPoolName, sdk.NewCoins(coin)); err != nil {
		return nil, sdkerrors.Wrap(err, "bank transfer failed")
	}

	if _, err := s.stakingKeeper.Delegate(ctx, delAddr, amount, stakingtypes.Unbonded, val, true); err != nil {
		return nil, sdkerrors.Wrap(err, "restake delegation failed")
	}

	sdkCtx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeClaimAndRestake,
			sdk.NewAttribute(types.AttributeKeyDelegator, msg.Delegator),
			sdk.NewAttribute(types.AttributeKeyValidator, msg.Validator),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})

	return &types.MsgClaimAndRestakeResponse{}, nil
}
