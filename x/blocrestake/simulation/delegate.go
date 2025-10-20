package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/lyfeloopinc/lyfebloc-network/x/blocrestake/keeper"
	"github.com/lyfeloopinc/lyfebloc-network/x/blocrestake/types"
)

func SimulateMsgDelegate(
	k keeper.Keeper,
	txGen client.TxConfig,
) simtypes.Operation {
	_ = k
	_ = txGen

	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		_ = app
		_ = ctx
		_ = chainID

		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgDelegate{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handle the Delegate simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "Delegate simulation not implemented"), nil, nil
	}
}
