package blocrestake

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"github.com/lyfeloopinc/lyfebloc-network/x/blocrestake/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod: "Delegate",
					Use:       "delegate [delegator] [validator] [amount]",
					Short:     "Send a Delegate tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "delegator"},
						{ProtoField: "validator"},
						{ProtoField: "amount"},
					},
				},
				{
					RpcMethod: "Undelegate",
					Use:       "undelegate [delegator] [validator] [amount]",
					Short:     "Send an Undelegate tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "delegator"},
						{ProtoField: "validator"},
						{ProtoField: "amount"},
					},
				},
				{
					RpcMethod: "ClaimAndRestake",
					Use:       "claim-and-restake [delegator] [validator]",
					Short:     "Send a ClaimAndRestake tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "delegator"},
						{ProtoField: "validator"},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
