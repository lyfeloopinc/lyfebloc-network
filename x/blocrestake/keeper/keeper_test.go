package keeper_test

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"cosmossdk.io/core/address"
	"cosmossdk.io/math"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	storetypes "cosmossdk.io/store/types"

	_ "github.com/lyfeloopinc/lyfebloc-network/app"
	"github.com/lyfeloopinc/lyfebloc-network/x/blocrestake/keeper"
	module "github.com/lyfeloopinc/lyfebloc-network/x/blocrestake/module"
	"github.com/lyfeloopinc/lyfebloc-network/x/blocrestake/types"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestMsgServerClaimAndRestake(t *testing.T) {
	t.Parallel()

	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModule{})
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(storeKey)
	sdkCtx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_test")).Ctx

	bankKeeper := newMockBankKeeper()
	stakingKeeper := newMockStakingKeeper("ulbt")
	distributionKeeper := newMockDistributionKeeper(bankKeeper)

	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())

	k := keeper.NewKeeper(
		storeService,
		encCfg.Codec,
		addressCodec,
		authority,
		nil,
		bankKeeper,
		stakingKeeper,
		distributionKeeper,
	)
	require.NoError(t, k.Params.Set(sdkCtx, types.DefaultParams()))

	msgServer := keeper.NewMsgServerImpl(k)

	delegator := sdk.AccAddress(bytes.Repeat([]byte{0x1}, 20))
	validator := sdk.ValAddress(bytes.Repeat([]byte{0x2}, 20))

	stakingKeeper.addValidator(stakingtypes.Validator{OperatorAddress: validator.String()})

	bondDenom, err := stakingKeeper.BondDenom(sdkCtx)
	require.NoError(t, err)

	rewardCoin := sdk.NewCoin(bondDenom, math.NewInt(100_000))
	require.NoError(t, bankKeeper.MintCoins(sdkCtx, distributiontypes.ModuleName, sdk.NewCoins(rewardCoin)))
	distributionKeeper.setRewards(delegator, validator, sdk.NewCoins(rewardCoin))

	beforeBal := bankKeeper.GetBalance(sdkCtx, delegator, bondDenom)
	require.True(t, beforeBal.Amount.IsZero(), "expected zero balance prior to claim")

	msg := &types.MsgClaimAndRestake{
		Delegator: delegator.String(),
		Validator: validator.String(),
	}

	_, err = msgServer.ClaimAndRestake(sdkCtx, msg)
	require.NoError(t, err)

	delegatedAmt := stakingKeeper.delegatedAmount(delegator)
	require.True(t, delegatedAmt.Equal(rewardCoin.Amount), "expected full reward to be re-delegated")

	afterBal := bankKeeper.GetBalance(sdkCtx, delegator, bondDenom)
	require.True(t, afterBal.Amount.IsZero(), "delegator balance should be zero after restake")

	events := sdk.UnwrapSDKContext(sdkCtx).EventManager().Events()
	found := false
	for _, evt := range events {
		if evt.Type == types.EventTypeClaimAndRestake {
			found = true
			break
		}
	}
	require.True(t, found, "expected claim_and_restake event to be emitted")
}

// -----------------------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------------------

type mockBankKeeper struct {
	accounts       map[string]sdk.Coins
	moduleAccounts map[string]sdk.Coins
}

func newMockBankKeeper() *mockBankKeeper {
	return &mockBankKeeper{
		accounts:       make(map[string]sdk.Coins),
		moduleAccounts: make(map[string]sdk.Coins),
	}
}

func (m *mockBankKeeper) ensureAccount(key string) sdk.Coins {
	coins, ok := m.accounts[key]
	if !ok {
		coins = sdk.NewCoins()
		m.accounts[key] = coins
	}
	return coins
}

func (m *mockBankKeeper) ensureModule(module string) sdk.Coins {
	coins, ok := m.moduleAccounts[module]
	if !ok {
		coins = sdk.NewCoins()
		m.moduleAccounts[module] = coins
	}
	return coins
}

func (m *mockBankKeeper) SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error {
	key := senderAddr.String()
	bal := m.ensureAccount(key)
	for _, coin := range amt {
		if bal.AmountOf(coin.Denom).LT(coin.Amount) {
			return sdkerrors.ErrInsufficientFunds
		}
	}
	m.accounts[key] = bal.Sub(amt...)

	moduleBal := m.ensureModule(recipientModule)
	m.moduleAccounts[recipientModule] = moduleBal.Add(amt...)
	return nil
}

func (m *mockBankKeeper) SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	moduleBal := m.ensureModule(senderModule)
	for _, coin := range amt {
		if moduleBal.AmountOf(coin.Denom).LT(coin.Amount) {
			return sdkerrors.ErrInsufficientFunds
		}
	}
	m.moduleAccounts[senderModule] = moduleBal.Sub(amt...)

	key := recipientAddr.String()
	bal := m.ensureAccount(key)
	m.accounts[key] = bal.Add(amt...)
	return nil
}

func (m *mockBankKeeper) MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error {
	moduleBal := m.ensureModule(moduleName)
	m.moduleAccounts[moduleName] = moduleBal.Add(amt...)
	return nil
}

func (m *mockBankKeeper) GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	bal := m.accounts[addr.String()]
	return sdk.NewCoin(denom, bal.AmountOf(denom))
}

// -----------------------------------------------------------------------------

type mockStakingKeeper struct {
	validators   map[string]stakingtypes.Validator
	delegations  map[string]math.Int
	bondDenomStr string
}

func newMockStakingKeeper(bondDenom string) *mockStakingKeeper {
	return &mockStakingKeeper{
		validators:   make(map[string]stakingtypes.Validator),
		delegations:  make(map[string]math.Int),
		bondDenomStr: bondDenom,
	}
}

func (m *mockStakingKeeper) addValidator(val stakingtypes.Validator) {
	m.validators[val.OperatorAddress] = val
}

func (m *mockStakingKeeper) GetValidator(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
	val, ok := m.validators[addr.String()]
	if !ok {
		return stakingtypes.Validator{}, stakingtypes.ErrNoValidatorFound
	}
	return val, nil
}

func (m *mockStakingKeeper) Delegate(ctx context.Context, delAddr sdk.AccAddress, amt math.Int, status stakingtypes.BondStatus, validator stakingtypes.Validator, subtractAccount bool) (math.LegacyDec, error) {
	key := delAddr.String()
	current, ok := m.delegations[key]
	if !ok {
		current = math.ZeroInt()
	}
	m.delegations[key] = current.Add(amt)
	return math.LegacyNewDecFromInt(amt), nil
}

func (m *mockStakingKeeper) Undelegate(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares math.LegacyDec) (time.Time, math.Int, error) {
	return time.Time{}, math.ZeroInt(), nil
}

func (m *mockStakingKeeper) BondDenom(ctx context.Context) (string, error) {
	return m.bondDenomStr, nil
}

func (m *mockStakingKeeper) delegatedAmount(addr sdk.AccAddress) math.Int {
	return m.delegations[addr.String()]
}

// -----------------------------------------------------------------------------

type mockDistributionKeeper struct {
	bank    *mockBankKeeper
	rewards map[string]sdk.Coins
}

func newMockDistributionKeeper(bank *mockBankKeeper) *mockDistributionKeeper {
	return &mockDistributionKeeper{
		bank:    bank,
		rewards: make(map[string]sdk.Coins),
	}
}

func (m *mockDistributionKeeper) rewardKey(del sdk.AccAddress, val sdk.ValAddress) string {
	return del.String() + "|" + val.String()
}

func (m *mockDistributionKeeper) setRewards(del sdk.AccAddress, val sdk.ValAddress, coins sdk.Coins) {
	m.rewards[m.rewardKey(del, val)] = coins
}

func (m *mockDistributionKeeper) WithdrawDelegationRewards(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error) {
	key := m.rewardKey(delAddr, valAddr)
	coins, ok := m.rewards[key]
	if !ok {
		return nil, distributiontypes.ErrNoValidatorDistInfo
	}
	if err := m.bank.SendCoinsFromModuleToAccount(ctx, distributiontypes.ModuleName, delAddr, coins); err != nil {
		return nil, err
	}
	delete(m.rewards, key)
	return coins, nil
}

type fixture struct {
	ctx          sdk.Context
	keeper       keeper.Keeper
	addressCodec address.Codec
}

func initFixture(t *testing.T) *fixture {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModule{})
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(storeKey)
	sdkCtx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_test")).Ctx

	bank := newMockBankKeeper()
	staking := newMockStakingKeeper("ulbt")
	distr := newMockDistributionKeeper(bank)

	authority := authtypes.NewModuleAddress(govtypes.ModuleName)
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())

	k := keeper.NewKeeper(
		storeService,
		encCfg.Codec,
		addressCodec,
		authority,
		nil,
		bank,
		staking,
		distr,
	)

	if err := k.Params.Set(sdkCtx, types.DefaultParams()); err != nil {
		t.Fatalf("failed to set params: %v", err)
	}

	return &fixture{
		ctx:          sdkCtx,
		keeper:       k,
		addressCodec: addressCodec,
	}
}
