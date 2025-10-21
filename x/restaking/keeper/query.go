package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	restakingv1 "github.com/lyfeloopinc/lyfebloc-network/lyfeblocnetwork/restaking/v1"
)

type queryServer struct {
	keeper Keeper
}

var _ restakingv1.QueryServer = queryServer{}

// NewQueryServer returns an implementation of the restaking gRPC query service.
func NewQueryServer(k Keeper) restakingv1.QueryServer {
	return queryServer{keeper: k}
}

func (q queryServer) Params(ctx context.Context, _ *restakingv1.QueryParamsRequest) (*restakingv1.QueryParamsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	ratio := q.keeper.GetAutoRestakeRatio(sdkCtx)
	return &restakingv1.QueryParamsResponse{AutoRestakeRatio: ratio.String()}, nil
}
