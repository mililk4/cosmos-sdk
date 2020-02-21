package keeper_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// tests GetDelegation, GetDelegatorDelegations, SetDelegation, RemoveDelegation, GetDelegatorDelegations
//func TestDelegation(t *testing.T) {
//	_, app, ctx := getBaseSimappWithCustomKeeper()
//
//	addrs := simapp.AddTestAddrs(app, ctx, 3, sdk.NewInt(10000))
//
//	//construct the validators
//	amts := []sdk.Int{sdk.NewInt(9), sdk.NewInt(8), sdk.NewInt(7)}
//	var validators [3]types.Validator
//	for i, amt := range amts {
//		validators[i] = types.NewValidator(sdk.ValAddress(addrs[i]), PKs[i], types.Description{})
//		validators[i], _ = validators[i].AddTokensFromDel(amt)
//	}
//
//	validators[0] = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validators[0], true)
//	validators[1] = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validators[1], true)
//	validators[2] = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validators[2], true)
//
//	// first add a validators[0] to delegate too
//	bond1to1 := types.NewDelegation(addrs[0], sdk.ValAddress(addrs[0]), sdk.NewDec(9))
//
//	// check the empty keeper first
//	_, found := app.StakingKeeper.GetDelegation(ctx, addrs[0], sdk.ValAddress(addrs[0]))
//	require.False(t, found)
//
//	// set and retrieve a record
//	app.StakingKeeper.SetDelegation(ctx, bond1to1)
//	resBond, found := app.StakingKeeper.GetDelegation(ctx, addrs[0], sdk.ValAddress(addrs[0]))
//	require.True(t, found)
//	require.True(t, bond1to1.Equal(resBond))
//
//	// modify a records, save, and retrieve
//	bond1to1.Shares = sdk.NewDec(99)
//	app.StakingKeeper.SetDelegation(ctx, bond1to1)
//	resBond, found = app.StakingKeeper.GetDelegation(ctx, addrs[0], sdk.ValAddress(addrs[0]))
//	require.True(t, found)
//	require.True(t, bond1to1.Equal(resBond))
//
//	// add some more records
//	bond1to2 := types.NewDelegation(addrs[0], sdk.ValAddress(addrs[1]), sdk.NewDec(9))
//	bond1to3 := types.NewDelegation(addrs[0], sdk.ValAddress(addrs[2]), sdk.NewDec(9))
//	bond2to1 := types.NewDelegation(addrs[1], sdk.ValAddress(addrs[0]), sdk.NewDec(9))
//	bond2to2 := types.NewDelegation(addrs[1], sdk.ValAddress(addrs[1]), sdk.NewDec(9))
//	bond2to3 := types.NewDelegation(addrs[1], sdk.ValAddress(addrs[2]), sdk.NewDec(9))
//	app.StakingKeeper.SetDelegation(ctx, bond1to2)
//	app.StakingKeeper.SetDelegation(ctx, bond1to3)
//	app.StakingKeeper.SetDelegation(ctx, bond2to1)
//	app.StakingKeeper.SetDelegation(ctx, bond2to2)
//	app.StakingKeeper.SetDelegation(ctx, bond2to3)
//
//	// test all bond retrieve capabilities
//	resBonds := app.StakingKeeper.GetDelegatorDelegations(ctx, addrs[0], 5)
//	require.Equal(t, 3, len(resBonds))
//	require.True(t, bond1to1.Equal(resBonds[0]))
//	require.True(t, bond1to2.Equal(resBonds[1]))
//	require.True(t, bond1to3.Equal(resBonds[2]))
//	resBonds = app.StakingKeeper.GetAllDelegatorDelegations(ctx, addrs[0])
//	require.Equal(t, 3, len(resBonds))
//	resBonds = app.StakingKeeper.GetDelegatorDelegations(ctx, addrs[0], 2)
//	require.Equal(t, 2, len(resBonds))
//	resBonds = app.StakingKeeper.GetDelegatorDelegations(ctx, addrs[1], 5)
//	require.Equal(t, 3, len(resBonds))
//	require.True(t, bond2to1.Equal(resBonds[0]))
//	require.True(t, bond2to2.Equal(resBonds[1]))
//	require.True(t, bond2to3.Equal(resBonds[2]))
//	allBonds := app.StakingKeeper.GetAllDelegations(ctx)
//	require.Equal(t, 6, len(allBonds))
//	require.True(t, bond1to1.Equal(allBonds[0]))
//	require.True(t, bond1to2.Equal(allBonds[1]))
//	require.True(t, bond1to3.Equal(allBonds[2]))
//	require.True(t, bond2to1.Equal(allBonds[3]))
//	require.True(t, bond2to2.Equal(allBonds[4]))
//	require.True(t, bond2to3.Equal(allBonds[5]))
//
//	resVals := app.StakingKeeper.GetDelegatorValidators(ctx, addrs[0], 3)
//	require.Equal(t, 3, len(resVals))
//	resVals = app.StakingKeeper.GetDelegatorValidators(ctx, addrs[1], 4)
//	require.Equal(t, 3, len(resVals))
//
//	for i := 0; i < 3; i++ {
//		resVal, err := app.StakingKeeper.GetDelegatorValidator(ctx, addrs[0], sdk.ValAddress(addrs[i]))
//		require.Nil(t, err)
//		require.Equal(t, addrs[i], resVal.GetOperator())
//
//		resVal, err = app.StakingKeeper.GetDelegatorValidator(ctx, addrs[1], sdk.ValAddress(addrs[i]))
//		require.Nil(t, err)
//		require.Equal(t, addrs[i], resVal.GetOperator())
//
//		resDels := app.StakingKeeper.GetValidatorDelegations(ctx, sdk.ValAddress(addrs[i]))
//		require.Len(t, resDels, 2)
//	}
//
//	// delete a record
//	app.StakingKeeper.RemoveDelegation(ctx, bond2to3)
//	_, found = app.StakingKeeper.GetDelegation(ctx, addrs[1], sdk.ValAddress(addrs[2]))
//	require.False(t, found)
//	resBonds = app.StakingKeeper.GetDelegatorDelegations(ctx, addrs[1], 5)
//	require.Equal(t, 2, len(resBonds))
//	require.True(t, bond2to1.Equal(resBonds[0]))
//	require.True(t, bond2to2.Equal(resBonds[1]))
//
//	resBonds = app.StakingKeeper.GetAllDelegatorDelegations(ctx, addrs[1])
//	require.Equal(t, 2, len(resBonds))
//
//	// delete all the records from delegator 2
//	app.StakingKeeper.RemoveDelegation(ctx, bond2to1)
//	app.StakingKeeper.RemoveDelegation(ctx, bond2to2)
//	_, found = app.StakingKeeper.GetDelegation(ctx, addrs[1], sdk.ValAddress(addrs[0]))
//	require.False(t, found)
//	_, found = app.StakingKeeper.GetDelegation(ctx, addrs[1], sdk.ValAddress(addrs[1]))
//	require.False(t, found)
//	resBonds = app.StakingKeeper.GetDelegatorDelegations(ctx, addrDels[1], 5)
//	require.Equal(t, 0, len(resBonds))
//}
//
// tests Get/Set/Remove UnbondingDelegation
func TestUnbondingDelegation(t *testing.T) {
	_, app, ctx := getBaseSimappWithCustomKeeper()

	ubd := types.NewUnbondingDelegation(
		addrDels[0],
		addrVals[0],
		0,
		time.Unix(0, 0),
		sdk.NewInt(5),
	)

	// set and retrieve a record
	app.StakingKeeper.SetUnbondingDelegation(ctx, ubd)
	resUnbond, found := app.StakingKeeper.GetUnbondingDelegation(ctx, addrDels[0], addrVals[0])
	require.True(t, found)
	require.True(t, ubd.Equal(resUnbond))

	// modify a records, save, and retrieve
	ubd.Entries[0].Balance = sdk.NewInt(21)
	app.StakingKeeper.SetUnbondingDelegation(ctx, ubd)

	resUnbonds := app.StakingKeeper.GetUnbondingDelegations(ctx, addrDels[0], 5)
	require.Equal(t, 1, len(resUnbonds))

	resUnbonds = app.StakingKeeper.GetAllUnbondingDelegations(ctx, addrDels[0])
	require.Equal(t, 1, len(resUnbonds))

	resUnbond, found = app.StakingKeeper.GetUnbondingDelegation(ctx, addrDels[0], addrVals[0])
	require.True(t, found)
	require.True(t, ubd.Equal(resUnbond))

	// delete a record
	app.StakingKeeper.RemoveUnbondingDelegation(ctx, ubd)
	_, found = app.StakingKeeper.GetUnbondingDelegation(ctx, addrDels[0], addrVals[0])
	require.False(t, found)

	resUnbonds = app.StakingKeeper.GetUnbondingDelegations(ctx, addrDels[0], 5)
	require.Equal(t, 0, len(resUnbonds))

	resUnbonds = app.StakingKeeper.GetAllUnbondingDelegations(ctx, addrDels[0])
	require.Equal(t, 0, len(resUnbonds))
}

func TestUnbondDelegation(t *testing.T) {
	_, app, ctx := getBaseSimappWithCustomKeeper()

	codec := simapp.NewAppCodec()
	app.StakingKeeper = keeper.NewKeeper(
		codec.Staking,
		app.GetKey(staking.StoreKey),
		app.BankKeeper,
		app.SupplyKeeper,
		app.GetSubspace(staking.ModuleName),
	)

	startTokens := sdk.TokensFromConsensusPower(10)
	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)

	require.NoError(t,
		app.BankKeeper.SetBalances(
			ctx,
			notBondedPool.GetAddress(),
			sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), startTokens)),
		),
	)
	app.SupplyKeeper.SetModuleAccount(ctx, notBondedPool)

	// create a validator and a delegator to that validator
	// note this validator starts not-bonded
	validator := types.NewValidator(addrVals[0], PKs[0], types.Description{})

	validator, issuedShares := validator.AddTokensFromDel(startTokens)
	require.Equal(t, startTokens, issuedShares.RoundInt())

	validator = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator, true)

	delegation := types.NewDelegation(addrDels[0], addrVals[0], issuedShares)
	app.StakingKeeper.SetDelegation(ctx, delegation)

	bondTokens := sdk.TokensFromConsensusPower(6)
	amount, err := app.StakingKeeper.Unbond(ctx, addrDels[0], addrVals[0], bondTokens.ToDec())
	require.NoError(t, err)
	require.Equal(t, bondTokens, amount) // shares to be added to an unbonding delegation

	delegation, found := app.StakingKeeper.GetDelegation(ctx, addrDels[0], addrVals[0])
	require.True(t, found)
	validator, found = app.StakingKeeper.GetValidator(ctx, addrVals[0])
	require.True(t, found)

	remainingTokens := startTokens.Sub(bondTokens)
	require.Equal(t, remainingTokens, delegation.Shares.RoundInt())
	require.Equal(t, remainingTokens, validator.BondedTokens())
}

//func TestUnbondingDelegationsMaxEntries(t *testing.T) {
//	app := simapp.Setup(false)
//	ctx := app.BaseApp.NewContext(false, abci.Header{})
//
//	codec := simapp.NewAppCodec()
//	app.StakingKeeper = keeper.NewKeeper(
//		codec.Staking,
//		app.GetKey(staking.StoreKey),
//		app.BankKeeper,
//		app.SupplyKeeper,
//		app.GetSubspace(staking.ModuleName),
//	)
//
//	startTokens := sdk.TokensFromConsensusPower(10)
//
//	bondDenom := app.StakingKeeper.BondDenom(ctx)
//	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
//
//	err := app.BankKeeper.SetBalances(ctx, notBondedPool.GetAddress(), sdk.NewCoins(sdk.NewCoin(bondDenom, startTokens)))
//	require.NoError(t, err)
//	app.SupplyKeeper.SetModuleAccount(ctx, notBondedPool)
//
//	// create a validator and a delegator to that validator
//	validator := types.NewValidator(addrVals[0], PKs[0], types.Description{})
//
//	validator, issuedShares := validator.AddTokensFromDel(startTokens)
//	require.Equal(t, startTokens, issuedShares.RoundInt())
//
//	validator = keeper.TestingUpdateValidator(app.StakingKeeper, ctx, validator, true)
//	require.True(sdk.IntEq(t, startTokens, validator.BondedTokens()))
//	require.True(t, validator.IsBonded())
//
//	delegation := types.NewDelegation(addrDels[0], addrVals[0], issuedShares)
//	app.StakingKeeper.SetDelegation(ctx, delegation)
//
//	maxEntries := app.StakingKeeper.MaxEntries(ctx)
//
//	oldBonded := app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetBondedPool(ctx).GetAddress(), bondDenom).Amount
//	oldNotBonded := app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetNotBondedPool(ctx).GetAddress(), bondDenom).Amount
//
//	// should all pass
//	var completionTime time.Time
//	for i := uint32(0); i < maxEntries; i++ {
//		var err error
//		completionTime, err = app.StakingKeeper.Undelegate(ctx, addrDels[0], addrVals[0], sdk.NewDec(1))
//		require.NoError(t, err)
//	}
//
//	newBonded := app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetBondedPool(ctx).GetAddress(), bondDenom).Amount
//	newNotBonded := app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetNotBondedPool(ctx).GetAddress(), bondDenom).Amount
//	require.True(sdk.IntEq(t, newBonded, oldBonded.SubRaw(int64(maxEntries))))
//	require.True(sdk.IntEq(t, newNotBonded, oldNotBonded.AddRaw(int64(maxEntries))))
//
//	oldBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetBondedPool(ctx).GetAddress(), bondDenom).Amount
//	oldNotBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetNotBondedPool(ctx).GetAddress(), bondDenom).Amount
//
//	// an additional unbond should fail due to max entries
//	_, err = app.StakingKeeper.Undelegate(ctx, addrDels[0], addrVals[0], sdk.NewDec(1))
//	require.Error(t, err)
//
//	newBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetBondedPool(ctx).GetAddress(), bondDenom).Amount
//	newNotBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetNotBondedPool(ctx).GetAddress(), bondDenom).Amount
//
//	require.True(sdk.IntEq(t, newBonded, oldBonded))
//	require.True(sdk.IntEq(t, newNotBonded, oldNotBonded))
//
//	// mature unbonding delegations
//	ctx = ctx.WithBlockTime(completionTime)
//	err = app.StakingKeeper.CompleteUnbonding(ctx, addrDels[0], addrVals[0])
//	require.NoError(t, err)
//
//	newBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetBondedPool(ctx).GetAddress(), bondDenom).Amount
//	newNotBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetNotBondedPool(ctx).GetAddress(), bondDenom).Amount
//	require.True(sdk.IntEq(t, newBonded, oldBonded))
//	require.True(sdk.IntEq(t, newNotBonded, oldNotBonded.SubRaw(int64(maxEntries))))
//
//	oldNotBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetNotBondedPool(ctx).GetAddress(), bondDenom).Amount
//
//	// unbonding  should work again
//	_, err = app.StakingKeeper.Undelegate(ctx, addrDels[0], addrVals[0], sdk.NewDec(1))
//	require.NoError(t, err)
//
//	newBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetBondedPool(ctx).GetAddress(), bondDenom).Amount
//	newNotBonded = app.BankKeeper.GetBalance(ctx, app.StakingKeeper.GetNotBondedPool(ctx).GetAddress(), bondDenom).Amount
//	require.True(sdk.IntEq(t, newBonded, oldBonded.SubRaw(1)))
//	require.True(sdk.IntEq(t, newNotBonded, oldNotBonded.AddRaw(1)))
//}
//
//// test undelegating self delegation from a validator pushing it below MinSelfDelegation
//// shift it from the bonded to unbonding state and jailed
//func TestUndelegateSelfDelegationBelowMinSelfDelegation(t *testing.T) {
//	ctx, _, bk, keeper, _ := CreateTestInput(t, false, 0)
//	delTokens := sdk.TokensFromConsensusPower(10)
//	delCoins := sdk.NewCoins(sdk.NewCoin(keeper.BondDenom(ctx), delTokens))
//
//	//create a validator with a self-delegation
//	validator := types.NewValidator(addrVals[0], PKs[0], types.Description{})
//
//	validator.MinSelfDelegation = delTokens
//	validator, issuedShares := validator.AddTokensFromDel(delTokens)
//	require.Equal(t, delTokens, issuedShares.RoundInt())
//
//	// add bonded tokens to pool for delegations
//	notBondedPool := keeper.GetNotBondedPool(ctx)
//	oldNotBonded := bk.GetAllBalances(ctx, notBondedPool.GetAddress())
//	err := bk.SetBalances(ctx, notBondedPool.GetAddress(), oldNotBonded.Add(delCoins...))
//	require.NoError(t, err)
//	keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)
//
//	validator = TestingUpdateValidator(keeper, ctx, validator, true)
//	require.True(t, validator.IsBonded())
//
//	selfDelegation := types.NewDelegation(sdk.AccAddress(addrVals[0].Bytes()), addrVals[0], issuedShares)
//	keeper.SetDelegation(ctx, selfDelegation)
//
//	// add bonded tokens to pool for delegations
//	bondedPool := keeper.GetBondedPool(ctx)
//	oldBonded := bk.GetAllBalances(ctx, bondedPool.GetAddress())
//	err = bk.SetBalances(ctx, bondedPool.GetAddress(), oldBonded.Add(delCoins...))
//	require.NoError(t, err)
//	keeper.supplyKeeper.SetModuleAccount(ctx, bondedPool)
//
//	// create a second delegation to this validator
//	keeper.DeleteValidatorByPowerIndex(ctx, validator)
//	validator, issuedShares = validator.AddTokensFromDel(delTokens)
//	require.True(t, validator.IsBonded())
//	require.Equal(t, delTokens, issuedShares.RoundInt())
//
//	// add bonded tokens to pool for delegations
//	oldBonded = bk.GetAllBalances(ctx, bondedPool.GetAddress())
//	err = bk.SetBalances(ctx, bondedPool.GetAddress(), oldBonded.Add(delCoins...))
//	require.NoError(t, err)
//	keeper.supplyKeeper.SetModuleAccount(ctx, bondedPool)
//
//	validator = TestingUpdateValidator(keeper, ctx, validator, true)
//	delegation := types.NewDelegation(addrDels[0], addrVals[0], issuedShares)
//	keeper.SetDelegation(ctx, delegation)
//
//	val0AccAddr := sdk.AccAddress(addrVals[0].Bytes())
//	_, err = keeper.Undelegate(ctx, val0AccAddr, addrVals[0], sdk.TokensFromConsensusPower(6).ToDec())
//	require.NoError(t, err)
//
//	// end block
//	updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
//	require.Equal(t, 1, len(updates))
//
//	validator, found := keeper.GetValidator(ctx, addrVals[0])
//	require.True(t, found)
//	require.Equal(t, sdk.TokensFromConsensusPower(14), validator.Tokens)
//	require.Equal(t, sdk.Unbonding, validator.Status)
//	require.True(t, validator.Jailed)
//}
//
//func TestUndelegateFromUnbondingValidator(t *testing.T) {
//	ctx, _, bk, keeper, _ := CreateTestInput(t, false, 0)
//	delTokens := sdk.TokensFromConsensusPower(10)
//	delCoins := sdk.NewCoins(sdk.NewCoin(keeper.BondDenom(ctx), delTokens))
//
//	//create a validator with a self-delegation
//	validator := types.NewValidator(addrVals[0], PKs[0], types.Description{})
//
//	validator, issuedShares := validator.AddTokensFromDel(delTokens)
//	require.Equal(t, delTokens, issuedShares.RoundInt())
//
//	// add bonded tokens to pool for delegations
//	notBondedPool := keeper.GetNotBondedPool(ctx)
//	oldNotBonded := bk.GetAllBalances(ctx, notBondedPool.GetAddress())
//	err := bk.SetBalances(ctx, notBondedPool.GetAddress(), oldNotBonded.Add(delCoins...))
//	require.NoError(t, err)
//	keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)
//
//	validator = TestingUpdateValidator(keeper, ctx, validator, true)
//	require.True(t, validator.IsBonded())
//
//	selfDelegation := types.NewDelegation(sdk.AccAddress(addrVals[0].Bytes()), addrVals[0], issuedShares)
//	keeper.SetDelegation(ctx, selfDelegation)
//
//	bondedPool := keeper.GetBondedPool(ctx)
//	oldBonded := bk.GetAllBalances(ctx, bondedPool.GetAddress())
//	err = bk.SetBalances(ctx, bondedPool.GetAddress(), oldBonded.Add(delCoins...))
//	require.NoError(t, err)
//	keeper.supplyKeeper.SetModuleAccount(ctx, bondedPool)
//
//	// create a second delegation to this validator
//	keeper.DeleteValidatorByPowerIndex(ctx, validator)
//
//	validator, issuedShares = validator.AddTokensFromDel(delTokens)
//	require.Equal(t, delTokens, issuedShares.RoundInt())
//
//	oldBonded = bk.GetAllBalances(ctx, bondedPool.GetAddress())
//	err = bk.SetBalances(ctx, bondedPool.GetAddress(), oldBonded.Add(delCoins...))
//	require.NoError(t, err)
//	keeper.supplyKeeper.SetModuleAccount(ctx, bondedPool)
//
//	validator = TestingUpdateValidator(keeper, ctx, validator, true)
//	delegation := types.NewDelegation(addrDels[0], addrVals[0], issuedShares)
//	keeper.SetDelegation(ctx, delegation)
//
//	oldBonded = bk.GetAllBalances(ctx, bondedPool.GetAddress())
//	err = bk.SetBalances(ctx, bondedPool.GetAddress(), oldBonded.Add(delCoins...))
//	require.NoError(t, err)
//	keeper.supplyKeeper.SetModuleAccount(ctx, bondedPool)
//
//	header := ctx.BlockHeader()
//	blockHeight := int64(10)
//	header.Height = blockHeight
//	blockTime := time.Unix(333, 0)
//	header.Time = blockTime
//	ctx = ctx.WithBlockHeader(header)
//
//	// unbond the all self-delegation to put validator in unbonding state
//	val0AccAddr := sdk.AccAddress(addrVals[0].Bytes())
//	_, err = keeper.Undelegate(ctx, val0AccAddr, addrVals[0], delTokens.ToDec())
//	require.NoError(t, err)
//
//	// end block
//	updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
//	require.Equal(t, 1, len(updates))
//
//	validator, found := keeper.GetValidator(ctx, addrVals[0])
//	require.True(t, found)
//	require.Equal(t, blockHeight, validator.UnbondingHeight)
//	params := keeper.GetParams(ctx)
//	require.True(t, blockTime.Add(params.UnbondingTime).Equal(validator.UnbondingTime))
//
//	blockHeight2 := int64(20)
//	blockTime2 := time.Unix(444, 0).UTC()
//	ctx = ctx.WithBlockHeight(blockHeight2)
//	ctx = ctx.WithBlockTime(blockTime2)
//
//	// unbond some of the other delegation's shares
//	_, err = keeper.Undelegate(ctx, addrDels[0], addrVals[0], sdk.NewDec(6))
//	require.NoError(t, err)
//
//	// retrieve the unbonding delegation
//	ubd, found := keeper.GetUnbondingDelegation(ctx, addrDels[0], addrVals[0])
//	require.True(t, found)
//	require.Len(t, ubd.Entries, 1)
//	require.True(t, ubd.Entries[0].Balance.Equal(sdk.NewInt(6)))
//	assert.Equal(t, blockHeight2, ubd.Entries[0].CreationHeight)
//	assert.True(t, blockTime2.Add(params.UnbondingTime).Equal(ubd.Entries[0].CompletionTime))
//}
//
//func TestUndelegateFromUnbondedValidator(t *testing.T) {
//	ctx, _, bk, keeper, _ := CreateTestInput(t, false, 1)
//	delTokens := sdk.TokensFromConsensusPower(10)
//	delCoins := sdk.NewCoins(sdk.NewCoin(keeper.BondDenom(ctx), delTokens))
//
//	// add bonded tokens to pool for delegations
//	notBondedPool := keeper.GetNotBondedPool(ctx)
//	oldNotBonded := bk.GetAllBalances(ctx, notBondedPool.GetAddress())
//	err := bk.SetBalances(ctx, notBondedPool.GetAddress(), oldNotBonded.Add(delCoins...))
//	require.NoError(t, err)
//	keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)
//
//	// create a validator with a self-delegation
//	validator := types.NewValidator(addrVals[0], PKs[0], types.Description{})
//
//	valTokens := sdk.TokensFromConsensusPower(10)
//	validator, issuedShares := validator.AddTokensFromDel(valTokens)
//	require.Equal(t, valTokens, issuedShares.RoundInt())
//	validator = TestingUpdateValidator(keeper, ctx, validator, true)
//	require.True(t, validator.IsBonded())
//
//	val0AccAddr := sdk.AccAddress(addrVals[0].Bytes())
//	selfDelegation := types.NewDelegation(val0AccAddr, addrVals[0], issuedShares)
//	keeper.SetDelegation(ctx, selfDelegation)
//
//	bondedPool := keeper.GetBondedPool(ctx)
//	oldBonded := bk.GetAllBalances(ctx, bondedPool.GetAddress())
//	err = bk.SetBalances(ctx, bondedPool.GetAddress(), oldBonded.Add(delCoins...))
//	require.NoError(t, err)
//	keeper.supplyKeeper.SetModuleAccount(ctx, bondedPool)
//
//	// create a second delegation to this validator
//	keeper.DeleteValidatorByPowerIndex(ctx, validator)
//	validator, issuedShares = validator.AddTokensFromDel(delTokens)
//	require.Equal(t, delTokens, issuedShares.RoundInt())
//	validator = TestingUpdateValidator(keeper, ctx, validator, true)
//	require.True(t, validator.IsBonded())
//	delegation := types.NewDelegation(addrDels[0], addrVals[0], issuedShares)
//	keeper.SetDelegation(ctx, delegation)
//
//	ctx = ctx.WithBlockHeight(10)
//	ctx = ctx.WithBlockTime(time.Unix(333, 0))
//
//	// unbond the all self-delegation to put validator in unbonding state
//	_, err = keeper.Undelegate(ctx, val0AccAddr, addrVals[0], valTokens.ToDec())
//	require.NoError(t, err)
//
//	// end block
//	updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
//	require.Equal(t, 1, len(updates))
//
//	validator, found := keeper.GetValidator(ctx, addrVals[0])
//	require.True(t, found)
//	require.Equal(t, ctx.BlockHeight(), validator.UnbondingHeight)
//	params := keeper.GetParams(ctx)
//	require.True(t, ctx.BlockHeader().Time.Add(params.UnbondingTime).Equal(validator.UnbondingTime))
//
//	// unbond the validator
//	ctx = ctx.WithBlockTime(validator.UnbondingTime)
//	keeper.UnbondAllMatureValidatorQueue(ctx)
//
//	// Make sure validator is still in state because there is still an outstanding delegation
//	validator, found = keeper.GetValidator(ctx, addrVals[0])
//	require.True(t, found)
//	require.Equal(t, validator.Status, sdk.Unbonded)
//
//	// unbond some of the other delegation's shares
//	unbondTokens := sdk.TokensFromConsensusPower(6)
//	_, err = keeper.Undelegate(ctx, addrDels[0], addrVals[0], unbondTokens.ToDec())
//	require.NoError(t, err)
//
//	// unbond rest of the other delegation's shares
//	remainingTokens := delTokens.Sub(unbondTokens)
//	_, err = keeper.Undelegate(ctx, addrDels[0], addrVals[0], remainingTokens.ToDec())
//	require.NoError(t, err)
//
//	//  now validator should now be deleted from state
//	validator, found = keeper.GetValidator(ctx, addrVals[0])
//	require.False(t, found, "%v", validator)
//}
//
//func TestUnbondingAllDelegationFromValidator(t *testing.T) {
//	ctx, _, bk, keeper, _ := CreateTestInput(t, false, 0)
//	delTokens := sdk.TokensFromConsensusPower(10)
//	delCoins := sdk.NewCoins(sdk.NewCoin(keeper.BondDenom(ctx), delTokens))
//
//	// add bonded tokens to pool for delegations
//	notBondedPool := keeper.GetNotBondedPool(ctx)
//	oldNotBonded := bk.GetAllBalances(ctx, notBondedPool.GetAddress())
//	err := bk.SetBalances(ctx, notBondedPool.GetAddress(), oldNotBonded.Add(delCoins...))
//	require.NoError(t, err)
//	keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)
//
//	//create a validator with a self-delegation
//	validator := types.NewValidator(addrVals[0], PKs[0], types.Description{})
//
//	valTokens := sdk.TokensFromConsensusPower(10)
//	validator, issuedShares := validator.AddTokensFromDel(valTokens)
//	require.Equal(t, valTokens, issuedShares.RoundInt())
//
//	validator = TestingUpdateValidator(keeper, ctx, validator, true)
//	require.True(t, validator.IsBonded())
//	val0AccAddr := sdk.AccAddress(addrVals[0].Bytes())
//
//	selfDelegation := types.NewDelegation(val0AccAddr, addrVals[0], issuedShares)
//	keeper.SetDelegation(ctx, selfDelegation)
//
//	// create a second delegation to this validator
//	keeper.DeleteValidatorByPowerIndex(ctx, validator)
//	validator, issuedShares = validator.AddTokensFromDel(delTokens)
//	require.Equal(t, delTokens, issuedShares.RoundInt())
//
//	bondedPool := keeper.GetBondedPool(ctx)
//	oldBonded := bk.GetAllBalances(ctx, bondedPool.GetAddress())
//	err = bk.SetBalances(ctx, bondedPool.GetAddress(), oldBonded.Add(delCoins...))
//	require.NoError(t, err)
//	keeper.supplyKeeper.SetModuleAccount(ctx, bondedPool)
//
//	validator = TestingUpdateValidator(keeper, ctx, validator, true)
//	require.True(t, validator.IsBonded())
//
//	delegation := types.NewDelegation(addrDels[0], addrVals[0], issuedShares)
//	keeper.SetDelegation(ctx, delegation)
//
//	ctx = ctx.WithBlockHeight(10)
//	ctx = ctx.WithBlockTime(time.Unix(333, 0))
//
//	// unbond the all self-delegation to put validator in unbonding state
//	_, err = keeper.Undelegate(ctx, val0AccAddr, addrVals[0], valTokens.ToDec())
//	require.NoError(t, err)
//
//	// end block
//	updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
//	require.Equal(t, 1, len(updates))
//
//	// unbond all the remaining delegation
//	_, err = keeper.Undelegate(ctx, addrDels[0], addrVals[0], delTokens.ToDec())
//	require.NoError(t, err)
//
//	// validator should still be in state and still be in unbonding state
//	validator, found := keeper.GetValidator(ctx, addrVals[0])
//	require.True(t, found)
//	require.Equal(t, validator.Status, sdk.Unbonding)
//
//	// unbond the validator
//	ctx = ctx.WithBlockTime(validator.UnbondingTime)
//	keeper.UnbondAllMatureValidatorQueue(ctx)
//
//	// validator should now be deleted from state
//	_, found = keeper.GetValidator(ctx, addrVals[0])
//	require.False(t, found)
//}
//
//// Make sure that that the retrieving the delegations doesn't affect the state
//func TestGetRedelegationsFromSrcValidator(t *testing.T) {
//	ctx, _, _, keeper, _ := CreateTestInput(t, false, 0)
//
//	rd := types.NewRedelegation(addrDels[0], addrVals[0], addrVals[1], 0,
//		time.Unix(0, 0), sdk.NewInt(5),
//		sdk.NewDec(5))
//
//	// set and retrieve a record
//	keeper.SetRedelegation(ctx, rd)
//	resBond, found := keeper.GetRedelegation(ctx, addrDels[0], addrVals[0], addrVals[1])
//	require.True(t, found)
//
//	// get the redelegations one time
//	redelegations := keeper.GetRedelegationsFromSrcValidator(ctx, addrVals[0])
//	require.Equal(t, 1, len(redelegations))
//	require.True(t, redelegations[0].Equal(resBond))
//
//	// get the redelegations a second time, should be exactly the same
//	redelegations = keeper.GetRedelegationsFromSrcValidator(ctx, addrVals[0])
//	require.Equal(t, 1, len(redelegations))
//	require.True(t, redelegations[0].Equal(resBond))
//}
//
//// tests Get/Set/Remove/Has UnbondingDelegation
//func TestRedelegation(t *testing.T) {
//	ctx, _, _, keeper, _ := CreateTestInput(t, false, 0)
//
//	rd := types.NewRedelegation(addrDels[0], addrVals[0], addrVals[1], 0,
//		time.Unix(0, 0), sdk.NewInt(5),
//		sdk.NewDec(5))
//
//	// test shouldn't have and redelegations
//	has := keeper.HasReceivingRedelegation(ctx, addrDels[0], addrVals[1])
//	require.False(t, has)
//
//	// set and retrieve a record
//	keeper.SetRedelegation(ctx, rd)
//	resRed, found := keeper.GetRedelegation(ctx, addrDels[0], addrVals[0], addrVals[1])
//	require.True(t, found)
//
//	redelegations := keeper.GetRedelegationsFromSrcValidator(ctx, addrVals[0])
//	require.Equal(t, 1, len(redelegations))
//	require.True(t, redelegations[0].Equal(resRed))
//
//	redelegations = keeper.GetRedelegations(ctx, addrDels[0], 5)
//	require.Equal(t, 1, len(redelegations))
//	require.True(t, redelegations[0].Equal(resRed))
//
//	redelegations = keeper.GetAllRedelegations(ctx, addrDels[0], nil, nil)
//	require.Equal(t, 1, len(redelegations))
//	require.True(t, redelegations[0].Equal(resRed))
//
//	// check if has the redelegation
//	has = keeper.HasReceivingRedelegation(ctx, addrDels[0], addrVals[1])
//	require.True(t, has)
//
//	// modify a records, save, and retrieve
//	rd.Entries[0].SharesDst = sdk.NewDec(21)
//	keeper.SetRedelegation(ctx, rd)
//
//	resRed, found = keeper.GetRedelegation(ctx, addrDels[0], addrVals[0], addrVals[1])
//	require.True(t, found)
//	require.True(t, rd.Equal(resRed))
//
//	redelegations = keeper.GetRedelegationsFromSrcValidator(ctx, addrVals[0])
//	require.Equal(t, 1, len(redelegations))
//	require.True(t, redelegations[0].Equal(resRed))
//
//	redelegations = keeper.GetRedelegations(ctx, addrDels[0], 5)
//	require.Equal(t, 1, len(redelegations))
//	require.True(t, redelegations[0].Equal(resRed))
//
//	// delete a record
//	keeper.RemoveRedelegation(ctx, rd)
//	_, found = keeper.GetRedelegation(ctx, addrDels[0], addrVals[0], addrVals[1])
//	require.False(t, found)
//
//	redelegations = keeper.GetRedelegations(ctx, addrDels[0], 5)
//	require.Equal(t, 0, len(redelegations))
//
//	redelegations = keeper.GetAllRedelegations(ctx, addrDels[0], nil, nil)
//	require.Equal(t, 0, len(redelegations))
//}
//
//func TestRedelegateToSameValidator(t *testing.T) {
//	ctx, _, bk, keeper, _ := CreateTestInput(t, false, 0)
//	valTokens := sdk.TokensFromConsensusPower(10)
//	startCoins := sdk.NewCoins(sdk.NewCoin(keeper.BondDenom(ctx), valTokens))
//
//	// add bonded tokens to pool for delegations
//	notBondedPool := keeper.GetNotBondedPool(ctx)
//	oldNotBonded := bk.GetAllBalances(ctx, notBondedPool.GetAddress())
//	err := bk.SetBalances(ctx, notBondedPool.GetAddress(), oldNotBonded.Add(startCoins...))
//	require.NoError(t, err)
//	keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)
//
//	// create a validator with a self-delegation
//	validator := types.NewValidator(addrVals[0], PKs[0], types.Description{})
//	validator, issuedShares := validator.AddTokensFromDel(valTokens)
//	require.Equal(t, valTokens, issuedShares.RoundInt())
//	validator = TestingUpdateValidator(keeper, ctx, validator, true)
//	require.True(t, validator.IsBonded())
//
//	val0AccAddr := sdk.AccAddress(addrVals[0].Bytes())
//	selfDelegation := types.NewDelegation(val0AccAddr, addrVals[0], issuedShares)
//	keeper.SetDelegation(ctx, selfDelegation)
//
//	_, err = keeper.BeginRedelegation(ctx, val0AccAddr, addrVals[0], addrVals[0], sdk.NewDec(5))
//	require.Error(t, err)
//}
//
//func TestRedelegationMaxEntries(t *testing.T) {
//	ctx, _, bk, keeper, _ := CreateTestInput(t, false, 0)
//	startTokens := sdk.TokensFromConsensusPower(20)
//	startCoins := sdk.NewCoins(sdk.NewCoin(keeper.BondDenom(ctx), startTokens))
//
//	// add bonded tokens to pool for delegations
//	notBondedPool := keeper.GetNotBondedPool(ctx)
//	oldNotBonded := bk.GetAllBalances(ctx, notBondedPool.GetAddress())
//	err := bk.SetBalances(ctx, notBondedPool.GetAddress(), oldNotBonded.Add(startCoins...))
//	require.NoError(t, err)
//	keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)
//
//	// create a validator with a self-delegation
//	validator := types.NewValidator(addrVals[0], PKs[0], types.Description{})
//	valTokens := sdk.TokensFromConsensusPower(10)
//	validator, issuedShares := validator.AddTokensFromDel(valTokens)
//	require.Equal(t, valTokens, issuedShares.RoundInt())
//	validator = TestingUpdateValidator(keeper, ctx, validator, true)
//	val0AccAddr := sdk.AccAddress(addrVals[0].Bytes())
//	selfDelegation := types.NewDelegation(val0AccAddr, addrVals[0], issuedShares)
//	keeper.SetDelegation(ctx, selfDelegation)
//
//	// create a second validator
//	validator2 := types.NewValidator(addrVals[1], PKs[1], types.Description{})
//	validator2, issuedShares = validator2.AddTokensFromDel(valTokens)
//	require.Equal(t, valTokens, issuedShares.RoundInt())
//
//	validator2 = TestingUpdateValidator(keeper, ctx, validator2, true)
//	require.Equal(t, sdk.Bonded, validator2.Status)
//
//	maxEntries := keeper.MaxEntries(ctx)
//
//	// redelegations should pass
//	var completionTime time.Time
//	for i := uint32(0); i < maxEntries; i++ {
//		var err error
//		completionTime, err = keeper.BeginRedelegation(ctx, val0AccAddr, addrVals[0], addrVals[1], sdk.NewDec(1))
//		require.NoError(t, err)
//	}
//
//	// an additional redelegation should fail due to max entries
//	_, err = keeper.BeginRedelegation(ctx, val0AccAddr, addrVals[0], addrVals[1], sdk.NewDec(1))
//	require.Error(t, err)
//
//	// mature redelegations
//	ctx = ctx.WithBlockTime(completionTime)
//	err = keeper.CompleteRedelegation(ctx, val0AccAddr, addrVals[0], addrVals[1])
//	require.NoError(t, err)
//
//	// redelegation should work again
//	_, err = keeper.BeginRedelegation(ctx, val0AccAddr, addrVals[0], addrVals[1], sdk.NewDec(1))
//	require.NoError(t, err)
//}
//
//func TestRedelegateSelfDelegation(t *testing.T) {
//	ctx, _, bk, keeper, _ := CreateTestInput(t, false, 0)
//	startTokens := sdk.TokensFromConsensusPower(30)
//	startCoins := sdk.NewCoins(sdk.NewCoin(keeper.BondDenom(ctx), startTokens))
//
//	// add bonded tokens to pool for delegations
//	notBondedPool := keeper.GetNotBondedPool(ctx)
//	oldNotBonded := bk.GetAllBalances(ctx, notBondedPool.GetAddress())
//	err := bk.SetBalances(ctx, notBondedPool.GetAddress(), oldNotBonded.Add(startCoins...))
//	require.NoError(t, err)
//	keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)
//
//	//create a validator with a self-delegation
//	validator := types.NewValidator(addrVals[0], PKs[0], types.Description{})
//	valTokens := sdk.TokensFromConsensusPower(10)
//	validator, issuedShares := validator.AddTokensFromDel(valTokens)
//	require.Equal(t, valTokens, issuedShares.RoundInt())
//
//	validator = TestingUpdateValidator(keeper, ctx, validator, true)
//
//	val0AccAddr := sdk.AccAddress(addrVals[0].Bytes())
//	selfDelegation := types.NewDelegation(val0AccAddr, addrVals[0], issuedShares)
//	keeper.SetDelegation(ctx, selfDelegation)
//
//	// create a second validator
//	validator2 := types.NewValidator(addrVals[1], PKs[1], types.Description{})
//	validator2, issuedShares = validator2.AddTokensFromDel(valTokens)
//	require.Equal(t, valTokens, issuedShares.RoundInt())
//	validator2 = TestingUpdateValidator(keeper, ctx, validator2, true)
//	require.Equal(t, sdk.Bonded, validator2.Status)
//
//	// create a second delegation to validator 1
//	delTokens := sdk.TokensFromConsensusPower(10)
//	validator, issuedShares = validator.AddTokensFromDel(delTokens)
//	require.Equal(t, delTokens, issuedShares.RoundInt())
//	validator = TestingUpdateValidator(keeper, ctx, validator, true)
//
//	delegation := types.NewDelegation(addrDels[0], addrVals[0], issuedShares)
//	keeper.SetDelegation(ctx, delegation)
//
//	_, err = keeper.BeginRedelegation(ctx, val0AccAddr, addrVals[0], addrVals[1], delTokens.ToDec())
//	require.NoError(t, err)
//
//	// end block
//	updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
//	require.Equal(t, 2, len(updates))
//
//	validator, found := keeper.GetValidator(ctx, addrVals[0])
//	require.True(t, found)
//	require.Equal(t, valTokens, validator.Tokens)
//	require.Equal(t, sdk.Unbonding, validator.Status)
//}
//
//func TestRedelegateFromUnbondingValidator(t *testing.T) {
//	ctx, _, bk, keeper, _ := CreateTestInput(t, false, 0)
//	startTokens := sdk.TokensFromConsensusPower(30)
//	startCoins := sdk.NewCoins(sdk.NewCoin(keeper.BondDenom(ctx), startTokens))
//
//	// add bonded tokens to pool for delegations
//	notBondedPool := keeper.GetNotBondedPool(ctx)
//	oldNotBonded := bk.GetAllBalances(ctx, notBondedPool.GetAddress())
//	err := bk.SetBalances(ctx, notBondedPool.GetAddress(), oldNotBonded.Add(startCoins...))
//	require.NoError(t, err)
//	keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)
//
//	//create a validator with a self-delegation
//	validator := types.NewValidator(addrVals[0], PKs[0], types.Description{})
//
//	valTokens := sdk.TokensFromConsensusPower(10)
//	validator, issuedShares := validator.AddTokensFromDel(valTokens)
//	require.Equal(t, valTokens, issuedShares.RoundInt())
//	validator = TestingUpdateValidator(keeper, ctx, validator, true)
//	val0AccAddr := sdk.AccAddress(addrVals[0].Bytes())
//	selfDelegation := types.NewDelegation(val0AccAddr, addrVals[0], issuedShares)
//	keeper.SetDelegation(ctx, selfDelegation)
//
//	// create a second delegation to this validator
//	keeper.DeleteValidatorByPowerIndex(ctx, validator)
//	delTokens := sdk.TokensFromConsensusPower(10)
//	validator, issuedShares = validator.AddTokensFromDel(delTokens)
//	require.Equal(t, delTokens, issuedShares.RoundInt())
//	validator = TestingUpdateValidator(keeper, ctx, validator, true)
//	delegation := types.NewDelegation(addrDels[0], addrVals[0], issuedShares)
//	keeper.SetDelegation(ctx, delegation)
//
//	// create a second validator
//	validator2 := types.NewValidator(addrVals[1], PKs[1], types.Description{})
//	validator2, issuedShares = validator2.AddTokensFromDel(valTokens)
//	require.Equal(t, valTokens, issuedShares.RoundInt())
//	validator2 = TestingUpdateValidator(keeper, ctx, validator2, true)
//
//	header := ctx.BlockHeader()
//	blockHeight := int64(10)
//	header.Height = blockHeight
//	blockTime := time.Unix(333, 0)
//	header.Time = blockTime
//	ctx = ctx.WithBlockHeader(header)
//
//	// unbond the all self-delegation to put validator in unbonding state
//	_, err = keeper.Undelegate(ctx, val0AccAddr, addrVals[0], delTokens.ToDec())
//	require.NoError(t, err)
//
//	// end block
//	updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
//	require.Equal(t, 1, len(updates))
//
//	validator, found := keeper.GetValidator(ctx, addrVals[0])
//	require.True(t, found)
//	require.Equal(t, blockHeight, validator.UnbondingHeight)
//	params := keeper.GetParams(ctx)
//	require.True(t, blockTime.Add(params.UnbondingTime).Equal(validator.UnbondingTime))
//
//	//change the context
//	header = ctx.BlockHeader()
//	blockHeight2 := int64(20)
//	header.Height = blockHeight2
//	blockTime2 := time.Unix(444, 0)
//	header.Time = blockTime2
//	ctx = ctx.WithBlockHeader(header)
//
//	// unbond some of the other delegation's shares
//	redelegateTokens := sdk.TokensFromConsensusPower(6)
//	_, err = keeper.BeginRedelegation(ctx, addrDels[0], addrVals[0], addrVals[1], redelegateTokens.ToDec())
//	require.NoError(t, err)
//
//	// retrieve the unbonding delegation
//	ubd, found := keeper.GetRedelegation(ctx, addrDels[0], addrVals[0], addrVals[1])
//	require.True(t, found)
//	require.Len(t, ubd.Entries, 1)
//	assert.Equal(t, blockHeight, ubd.Entries[0].CreationHeight)
//	assert.True(t, blockTime.Add(params.UnbondingTime).Equal(ubd.Entries[0].CompletionTime))
//}
//
//func TestRedelegateFromUnbondedValidator(t *testing.T) {
//	ctx, _, bk, keeper, _ := CreateTestInput(t, false, 0)
//	startTokens := sdk.TokensFromConsensusPower(30)
//	startCoins := sdk.NewCoins(sdk.NewCoin(keeper.BondDenom(ctx), startTokens))
//
//	// add bonded tokens to pool for delegations
//	notBondedPool := keeper.GetNotBondedPool(ctx)
//	oldNotBonded := bk.GetAllBalances(ctx, notBondedPool.GetAddress())
//	err := bk.SetBalances(ctx, notBondedPool.GetAddress(), oldNotBonded.Add(startCoins...))
//	require.NoError(t, err)
//	keeper.supplyKeeper.SetModuleAccount(ctx, notBondedPool)
//
//	//create a validator with a self-delegation
//	validator := types.NewValidator(addrVals[0], PKs[0], types.Description{})
//
//	valTokens := sdk.TokensFromConsensusPower(10)
//	validator, issuedShares := validator.AddTokensFromDel(valTokens)
//	require.Equal(t, valTokens, issuedShares.RoundInt())
//	validator = TestingUpdateValidator(keeper, ctx, validator, true)
//	val0AccAddr := sdk.AccAddress(addrVals[0].Bytes())
//	selfDelegation := types.NewDelegation(val0AccAddr, addrVals[0], issuedShares)
//	keeper.SetDelegation(ctx, selfDelegation)
//
//	// create a second delegation to this validator
//	keeper.DeleteValidatorByPowerIndex(ctx, validator)
//	delTokens := sdk.TokensFromConsensusPower(10)
//	validator, issuedShares = validator.AddTokensFromDel(delTokens)
//	require.Equal(t, delTokens, issuedShares.RoundInt())
//	validator = TestingUpdateValidator(keeper, ctx, validator, true)
//	delegation := types.NewDelegation(addrDels[0], addrVals[0], issuedShares)
//	keeper.SetDelegation(ctx, delegation)
//
//	// create a second validator
//	validator2 := types.NewValidator(addrVals[1], PKs[1], types.Description{})
//	validator2, issuedShares = validator2.AddTokensFromDel(valTokens)
//	require.Equal(t, valTokens, issuedShares.RoundInt())
//	validator2 = TestingUpdateValidator(keeper, ctx, validator2, true)
//	require.Equal(t, sdk.Bonded, validator2.Status)
//
//	ctx = ctx.WithBlockHeight(10)
//	ctx = ctx.WithBlockTime(time.Unix(333, 0))
//
//	// unbond the all self-delegation to put validator in unbonding state
//	_, err = keeper.Undelegate(ctx, val0AccAddr, addrVals[0], delTokens.ToDec())
//	require.NoError(t, err)
//
//	// end block
//	updates := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
//	require.Equal(t, 1, len(updates))
//
//	validator, found := keeper.GetValidator(ctx, addrVals[0])
//	require.True(t, found)
//	require.Equal(t, ctx.BlockHeight(), validator.UnbondingHeight)
//	params := keeper.GetParams(ctx)
//	require.True(t, ctx.BlockHeader().Time.Add(params.UnbondingTime).Equal(validator.UnbondingTime))
//
//	// unbond the validator
//	keeper.unbondingToUnbonded(ctx, validator)
//
//	// redelegate some of the delegation's shares
//	redelegationTokens := sdk.TokensFromConsensusPower(6)
//	_, err = keeper.BeginRedelegation(ctx, addrDels[0], addrVals[0], addrVals[1], redelegationTokens.ToDec())
//	require.NoError(t, err)
//
//	// no red should have been found
//	red, found := keeper.GetRedelegation(ctx, addrDels[0], addrVals[0], addrVals[1])
//	require.False(t, found, "%v", red)
//}
