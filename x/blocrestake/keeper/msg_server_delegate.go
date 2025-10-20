package keeper

import (
	"context"
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	math "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/lyfeloopinc/lyfebloc-network/x/blocrestake/types"
)

func (s msgServer) Delegate(ctx context.Context, msg *types.MsgDelegate) (*types.MsgDelegateResponse, error) {
	delegator, err := sdk.AccAddressFromBech32(msg.Delegator)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidAddress, fmt.Sprintf("invalid delegator address: %s", err))
	}

	valAddr, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidAddress, fmt.Sprintf("invalid validator address: %s", err))
	}

	if msg.Amount == 0 {
		return nil, errorsmod.Wrap(types.ErrInvalidAmount, "amount must be positive")
	}

	val, err := s.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		if errors.Is(err, stakingtypes.ErrNoValidatorFound) {
			return nil, types.ErrValidatorNotFound
		}
		return nil, errorsmod.Wrap(err, "failed to fetch validator")
	}

	amount := math.NewIntFromUint64(msg.Amount)
	if !amount.IsPositive() {
		return nil, errorsmod.Wrap(types.ErrInvalidAmount, "amount must be positive")
	}

	bondDenom, err := s.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to fetch bond denom")
	}

	coin := sdk.NewCoin(bondDenom, amount)
	if err := s.bankKeeper.SendCoinsFromAccountToModule(ctx, delegator, stakingtypes.BondedPoolName, sdk.NewCoins(coin)); err != nil {
		return nil, errorsmod.Wrap(err, "bank transfer failed")
	}

	if _, err := s.stakingKeeper.Delegate(ctx, delegator, amount, stakingtypes.Unbonded, val, true); err != nil {
		return nil, errorsmod.Wrap(err, "staking delegate failed")
	}

	return &types.MsgDelegateResponse{}, nil
}

func (s msgServer) Undelegate(ctx context.Context, msg *types.MsgUndelegate) (*types.MsgUndelegateResponse, error) {
	delegator, err := sdk.AccAddressFromBech32(msg.Delegator)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidAddress, fmt.Sprintf("invalid delegator address: %s", err))
	}

	valAddr, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidAddress, fmt.Sprintf("invalid validator address: %s", err))
	}

	if msg.Amount == 0 {
		return nil, errorsmod.Wrap(types.ErrInvalidAmount, "amount must be positive")
	}

	amount := math.NewIntFromUint64(msg.Amount)
	if !amount.IsPositive() {
		return nil, errorsmod.Wrap(types.ErrInvalidAmount, "amount must be positive")
	}

	shares := math.LegacyNewDecFromInt(amount)

	if _, _, err := s.stakingKeeper.Undelegate(ctx, delegator, valAddr, shares); err != nil {
		return nil, errorsmod.Wrap(err, "undelegate failed")
	}

	// undelegated tokens move to the not bonded pool automatically; nothing else to do here
	return &types.MsgUndelegateResponse{}, nil
}
