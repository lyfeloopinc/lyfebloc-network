package keeper

import (
	"errors"

	"cosmossdk.io/collections"
	corestore "cosmossdk.io/core/store"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/lyfeloopinc/lyfebloc-network/x/restaking/types"
)

// Keeper manages chain-wide restake parameters and execution.
type Keeper struct {
	storeService  corestore.KVStoreService
	stakingKeeper types.StakingKeeper

	schema             collections.Schema
	autoRestakeRatioIt collections.Item[sdkmath.LegacyDec]
}

func NewKeeper(storeService corestore.KVStoreService, stakingKeeper types.StakingKeeper) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService:       storeService,
		stakingKeeper:      stakingKeeper,
		autoRestakeRatioIt: collections.NewItem(sb, types.AutoRestakeRatioKey, "auto_restake_ratio", sdk.LegacyDecValue),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.schema = schema

	return k
}

// Schema returns the keeper schema.
func (k Keeper) Schema() collections.Schema {
	return k.schema
}

// GetAutoRestakeRatio returns the configured auto restake ratio or the default.
func (k Keeper) GetAutoRestakeRatio(ctx sdk.Context) sdkmath.LegacyDec {
	ratio, err := k.autoRestakeRatioIt.Get(ctx)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return types.DefaultAutoRestakeRatioDec()
		}
		panic(err)
	}
	return ratio
}

// SetAutoRestakeRatio stores a new auto restake ratio after validation.
func (k Keeper) SetAutoRestakeRatio(ctx sdk.Context, ratio sdkmath.LegacyDec) error {
	if err := types.ValidateAutoRestakeRatio(ratio); err != nil {
		return err
	}
	return k.autoRestakeRatioIt.Set(ctx, ratio)
}

// AutoRestakeRewards calculates the portion of rewards to restake based on the global ratio.
func (k Keeper) AutoRestakeRewards(ctx sdk.Context, rewards sdk.Coins) sdk.Coins {
	if rewards.IsZero() {
		return nil
	}

	ratio := k.GetAutoRestakeRatio(ctx)
	if ratio.IsZero() {
		return nil
	}

	restake := sdk.NewCoins()
	for _, coin := range rewards {
		portion := sdkmath.LegacyNewDecFromInt(coin.Amount).Mul(ratio).TruncateInt()
		if portion.IsPositive() {
			restake = restake.Add(sdk.NewCoin(coin.Denom, portion))
		}
	}

	if restake.IsZero() {
		return nil
	}

	return restake
}
