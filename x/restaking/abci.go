package restaking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/lyfeloopinc/lyfebloc-network/x/restaking/keeper"
)

// EndBlocker scans withdrawal events and applies auto-restake logic.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	events := ctx.EventManager().Events()
	for _, ev := range events {
		if ev.Type != distributiontypes.EventTypeWithdrawRewards {
			continue
		}

		var (
			delAddr sdk.AccAddress
			valAddr sdk.ValAddress
			rewards sdk.Coins
		)

		for _, attr := range ev.Attributes {
			switch string(attr.Key) {
			case distributiontypes.AttributeKeyDelegator:
				addr, err := sdk.AccAddressFromBech32(string(attr.Value))
				if err == nil {
					delAddr = addr
				}
			case distributiontypes.AttributeKeyValidator:
				addr, err := sdk.ValAddressFromBech32(string(attr.Value))
				if err == nil {
					valAddr = addr
				}
			case sdk.AttributeKeyAmount:
				coins, err := sdk.ParseCoinsNormalized(string(attr.Value))
				if err == nil {
					rewards = coins
				}
			}
		}

		if delAddr == nil || valAddr == nil || rewards.Empty() {
			continue
		}

		portion := k.AutoRestakeRewards(ctx, rewards)
		if portion == nil || portion.IsZero() {
			continue
		}

		if err := k.RestakeDelegate(ctx, delAddr, valAddr, portion); err != nil {
			ctx.Logger().Error("auto-restake delegate failed", "err", err, "delegator", delAddr.String(), "validator", valAddr.String())
			continue
		}

		ctx.Logger().Info("auto-restake executed", "delegator", delAddr.String(), "validator", valAddr.String(), "amount", portion.String(), "height", ctx.BlockHeight())
	}
}
