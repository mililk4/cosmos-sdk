package slashing

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) onValidatorBonded(ctx sdk.Context, address sdk.ConsAddress) {
	// Create a new signing info if necessary
	_, found := k.getValidatorSigningInfo(ctx, address)
	if !found {
		signingInfo := ValidatorSigningInfo{
			StartHeight:         ctx.BlockHeight(),
			IndexOffset:         0,
			JailedUntil:         time.Unix(0, 0),
			SignedBlocksCounter: 0,
		}
		k.setValidatorSigningInfo(ctx, address, signingInfo)
	}

	// Create a new slashing period when a validator is bonded
	slashingPeriod := ValidatorSlashingPeriod{
		ValidatorAddr: address,
		StartHeight:   ctx.BlockHeight(),
		EndHeight:     0,
		SlashedSoFar:  sdk.ZeroDec(),
	}
	k.addOrUpdateValidatorSlashingPeriod(ctx, slashingPeriod)
}

// Mark the slashing period as having ended when a validator begins unbonding
func (k Keeper) onValidatorBeginUnbonding(ctx sdk.Context, address sdk.ConsAddress) {
	slashingPeriod := k.getValidatorSlashingPeriodForHeight(ctx, address, ctx.BlockHeight())
	slashingPeriod.EndHeight = ctx.BlockHeight()
	k.addOrUpdateValidatorSlashingPeriod(ctx, slashingPeriod)
}

//_________________________________________________________________________________________

// Wrapper struct
type Hooks struct {
	k Keeper
}

var _ sdk.StakingHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// Implements sdk.ValidatorHooks
func (h Hooks) OnValidatorBonded(ctx sdk.Context, address sdk.ConsAddress) {
	h.k.onValidatorBonded(ctx, address)
}

// Implements sdk.ValidatorHooks
func (h Hooks) OnValidatorBeginUnbonding(ctx sdk.Context, address sdk.ConsAddress) {
	h.k.onValidatorBeginUnbonding(ctx, address)
}

// nolint - unused hooks
func (h Hooks) OnValidatorCreated(_ sdk.Context, _ sdk.ValAddress)                           {}
func (h Hooks) OnValidatorCommissionChange(_ sdk.Context, _ sdk.ValAddress)                  {}
func (h Hooks) OnValidatorRemoved(_ sdk.Context, _ sdk.ValAddress)                           {}
func (h Hooks) OnDelegationCreated(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress)        {}
func (h Hooks) OnDelegationSharesModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) {}
func (h Hooks) OnDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress)        {}
