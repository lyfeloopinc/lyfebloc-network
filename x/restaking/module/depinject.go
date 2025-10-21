package module

import (
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"github.com/cosmos/cosmos-sdk/codec"

	restakingmodpb "github.com/lyfeloopinc/lyfebloc-network/lyfeblocnetwork/restaking/module/v1"
	"github.com/lyfeloopinc/lyfebloc-network/x/restaking/keeper"
	"github.com/lyfeloopinc/lyfebloc-network/x/restaking/types"
)

var _ depinject.OnePerModuleType = AppModule{}

func (AppModule) IsOnePerModuleType() {}

func init() {
	appconfig.Register(
		&restakingmodpb.Module{},
		appconfig.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	StoreService  store.KVStoreService
	Cdc           codec.Codec
	StakingKeeper types.StakingKeeper
}

type ModuleOutputs struct {
	depinject.Out

	RestakingKeeper keeper.Keeper
	Module          appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	k := keeper.NewKeeper(in.StoreService, in.StakingKeeper)
	m := NewAppModule(k)

	return ModuleOutputs{
		RestakingKeeper: k,
		Module:          m,
	}
}
