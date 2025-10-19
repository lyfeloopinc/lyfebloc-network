package keeper

import (
	"context"

    "github.com/lyfeloopinc/lyfebloc-network/x/blocrestake/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)


func (k msgServer) Delegate(ctx context.Context,  msg *types.MsgDelegate) (*types.MsgDelegateResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(err, "invalid authority address")
	}

    // TODO: Handle the message

	return &types.MsgDelegateResponse{}, nil
}
