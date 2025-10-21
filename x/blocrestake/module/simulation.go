package blocrestake

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

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
		opWeightMsgDelegate          = "op_weight_msg_blocrestake_delegate"
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
		blocrestakesimulation.SimulateMsgDelegate(am.keeper, simState.TxConfig),
	))

	const (
		opWeightMsgUndelegate          = "op_weight_msg_blocrestake_undelegate"
		defaultWeightMsgUndelegate int = 100
	)

	var weightMsgUndelegate int
	simState.AppParams.GetOrGenerate(opWeightMsgUndelegate, &weightMsgUndelegate, nil,
		func(_ *rand.Rand) {
			weightMsgUndelegate = defaultWeightMsgUndelegate
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUndelegate,
		blocrestakesimulation.SimulateMsgUndelegate(am.keeper, simState.TxConfig),
	))

	const (
		opWeightMsgClaimAndRestake          = "op_weight_msg_blocrestake_claim_and_restake"
		defaultWeightMsgClaimAndRestake int = 100
	)

	var weightMsgClaimAndRestake int
	simState.AppParams.GetOrGenerate(opWeightMsgClaimAndRestake, &weightMsgClaimAndRestake, nil,
		func(_ *rand.Rand) {
			weightMsgClaimAndRestake = defaultWeightMsgClaimAndRestake
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgClaimAndRestake,
		blocrestakesimulation.SimulateMsgClaimAndRestake(am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
