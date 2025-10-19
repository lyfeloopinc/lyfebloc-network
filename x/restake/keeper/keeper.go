package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"

	"github.com/lyfeloopinc/lyfebloc-network/x/blocrestake/types"
)

type (
	Keeper struct {
		storeKey   sdk.StoreKey
		cdc        codec.BinaryCodec
		bankKeeper keeper.Keeper
		stkKeeper  *stakingkeeper.Keeper
		ibcKeeper  ibcexported.ChannelKeeper
	}
)

// NewKeeper creates a new blocrestake Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	key sdk.StoreKey,
	bk keeper.Keeper,
	sk *stakingkeeper.Keeper,
	ibc ibcexported.ChannelKeeper,
) Keeper {
	return Keeper{
		storeKey:   key,
		cdc:        cdc,
		bankKeeper: bk,
		stkKeeper:  sk,
		ibcKeeper:  ibc,
	}
}

// DelegateTokens lets a user delegate tokens to a validator
func (k Keeper) DelegateTokens(ctx sdk.Context, delegator sdk.AccAddress, validator sdk.ValAddress, amount sdk.Coin) error {
	val, found := k.stkKeeper.GetValidator(ctx, validator)
	if !found {
		return types.ErrValidatorNotFound
	}

	_, err := k.stkKeeper.Delegate(ctx, delegator, amount.Amount, stakingtypes.Unbonded, val, true)
	return err
}

// UndelegateTokens undelegates tokens from a validator
func (k Keeper) UndelegateTokens(ctx sdk.Context, delegator sdk.AccAddress, validator sdk.ValAddress, amount sdk.Coin) error {
	_, err := k.stkKeeper.Undelegate(ctx, delegator, validator, amount.Amount)
	return err
}

// ClaimAndRestake claims rewards and re-delegates them automatically
func (k Keeper) ClaimAndRestake(ctx sdk.Context, delegator sdk.AccAddress, validator sdk.ValAddress) error {
	rewards := sdk.NewCoin("ulbt", sdk.NewInt(1000)) // Placeholder — hook this to distribution module later
	err := k.DelegateTokens(ctx, delegator, validator, rewards)
	return err
}
