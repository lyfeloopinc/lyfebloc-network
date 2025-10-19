package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "cosmossdk.io/errors"
	"github.com/lyfeloopinc/lyfebloc-network/x/blocrestake/types"
)

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

func (s msgServer) Delegate(goCtx context.Context, msg *types.MsgDelegate) (*types.MsgDelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	delegator, err := sdk.AccAddressFromBech32(msg.Delegator)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidAddress, err.Error())
	}
	validator, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrInvalidAddress, err.Error())
	}
	coin := sdk.NewCoin("ulbt", sdk.NewInt(int64(msg.Amount)))
	if err := s.DelegateTokens(ctx, delegator, validator, coin); err != nil {
		return nil, err
	}
	return &types.MsgDelegateResponse{}, nil
}

func (s msgServer) Undelegate(goCtx context.Context, msg *types.MsgUndelegate) (*types.MsgUndelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	delegator, _ := sdk.AccAddressFromBech32(msg.Delegator)
	validator, _ := sdk.ValAddressFromBech32(msg.Validator)
	coin := sdk.NewCoin("ulbt", sdk.NewInt(int64(msg.Amount)))
	if err := s.UndelegateTokens(ctx, delegator, validator, coin); err != nil {
		return nil, err
	}
	return &types.MsgUndelegateResponse{}, nil
}

func (s msgServer) ClaimAndRestake(goCtx context.Context, msg *types.MsgClaimAndRestake) (*types.MsgClaimAndRestakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	delegator, _ := sdk.AccAddressFromBech32(msg.Delegator)
	validator, _ := sdk.ValAddressFromBech32(msg.Validator)
	if err := s.ClaimAndRestake(ctx, delegator, validator); err != nil {
		return nil, err
	}
	return &types.MsgClaimAndRestakeResponse{}, nil
}
