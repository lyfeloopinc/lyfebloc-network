package module

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc"

	"cosmossdk.io/core/appmodule"

	restakingv1 "github.com/lyfeloopinc/lyfebloc-network/lyfeblocnetwork/restaking/v1"
	restaking "github.com/lyfeloopinc/lyfebloc-network/x/restaking"
	"github.com/lyfeloopinc/lyfebloc-network/x/restaking/keeper"
	"github.com/lyfeloopinc/lyfebloc-network/x/restaking/types"
)

type AppModule struct {
	keeper keeper.Keeper
}

var _ appmodule.AppModule = AppModule{}
var _ appmodule.HasServices = AppModule{}
var _ appmodule.HasEndBlocker = AppModule{}

func NewAppModule(k keeper.Keeper) AppModule {
	return AppModule{keeper: k}
}

func (AppModule) IsAppModule() {}

func (AppModule) Name() string { return types.ModuleName }

func (am AppModule) RegisterServices(registrar grpc.ServiceRegistrar) error {
	restakingv1.RegisterQueryServer(registrar, keeper.NewQueryServer(am.keeper))
	return nil
}

func (am AppModule) EndBlock(ctx context.Context) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	restaking.EndBlocker(sdkCtx, am.keeper)
	return nil
}
