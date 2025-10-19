package types

import (
	"context"

	"cosmossdk.io/core/address"
    stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	
    sdk "github.com/cosmos/cosmos-sdk/types"
)


    
    

    
        // StakingKeeper defines the expected interface for the Staking module.
        type StakingKeeper interface {
        	ConsensusAddressCodec() address.Codec
        	ValidatorByConsAddr(context.Context, sdk.ConsAddress) (stakingtypes.ValidatorI, error)
        	// Methods imported from account should be defined here
        }

    


// AuthKeeper defines the expected interface for the Auth module.
type AuthKeeper interface {
    AddressCodec() address.Codec
    GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
    // Methods imported from account should be defined here
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
    SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
    // Methods imported from bank should be defined here
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}
