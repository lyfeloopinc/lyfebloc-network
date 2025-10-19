package restake

import (
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"

	"github.com/lyfeloopinc/lyfebloc-network/x/restake/keeper"
	"github.com/lyfeloopinc/lyfebloc-network/x/restake/types"
)

var _ depinject.OnePerModuleType = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (AppModule) IsOnePerModuleType() {}

func init() {
	appconfig.Register(
	    &types.Module{},
		appconfig.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	Config       *types.Module
	StoreService  store.KVStoreService
	Cdc          codec.Codec
	AddressCodec address.Codec

	AuthKeeper types.AuthKeeper
	BankKeeper types.BankKeeper
    StakingKeeper types.StakingKeeper

    IBCKeeperFn func() *ibckeeper.Keeper `optional:"true"` 
}

type ModuleOutputs struct {
	depinject.Out

	RestakeKeeper keeper.Keeper
	Module appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	// default to governance authority if not provided
	authority := authtypes.NewModuleAddress(types.GovModuleName)
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}
	k := keeper.NewKeeper(
		in.StoreService,
	    in.Cdc,
		in.AddressCodec,
	    authority, 
		in.IBCKeeperFn,
        in.BankKeeper,
        in.StakingKeeper,
	)
	m := NewAppModule(in.Cdc, k, in.AuthKeeper, in.BankKeeper)

	return ModuleOutputs{RestakeKeeper: k, Module: m}
}
