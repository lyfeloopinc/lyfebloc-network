package blocrestake

import (
	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"math/rand"

	blocrestakesimulation "github.com/lyfeloopinc/lyfebloc-network/x/blocrestake/simulation"
	"github.com/lyfeloopinc/lyfebloc-network/x/blocrestake/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	blocrestakeGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&blocrestakeGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgDelegate          = "op_weight_msg_blocrestake"
		defaultWeightMsgDelegate int = 100
	)

	var weightMsgDelegate int
	simState.AppParams.GetOrGenerate(opWeightMsgDelegate, &weightMsgDelegate, nil,
		func(_ *rand.Rand) {
			weightMsgDelegate = defaultWeightMsgDelegate
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDelegate,
		blocrestakesimulation.SimulateMsgDelegate(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}