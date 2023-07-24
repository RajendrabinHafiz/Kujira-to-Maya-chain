package mayachain

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/blang/semver"

	"gitlab.com/mayachain/mayanode/common"
	"gitlab.com/mayachain/mayanode/common/cosmos"
	"gitlab.com/mayachain/mayanode/constants"
)

var IsValidMAYANameV1 = regexp.MustCompile(`^[a-zA-Z0-9+_-]+$`).MatchString

// ManageMAYANameHandler a handler to process MsgNetworkFee messages
type ManageMAYANameHandler struct {
	mgr Manager
}

// NewManageMAYANameHandler create a new instance of network fee handler
func NewManageMAYANameHandler(mgr Manager) ManageMAYANameHandler {
	return ManageMAYANameHandler{mgr: mgr}
}

// Run is the main entry point for network fee logic
func (h ManageMAYANameHandler) Run(ctx cosmos.Context, m cosmos.Msg) (*cosmos.Result, error) {
	msg, ok := m.(*MsgManageMAYAName)
	if !ok {
		return nil, errInvalidMessage
	}
	if err := h.validate(ctx, *msg); err != nil {
		ctx.Logger().Error("MsgManageMAYAName failed validation", "error", err)
		return nil, err
	}
	result, err := h.handle(ctx, *msg)
	if err != nil {
		ctx.Logger().Error("fail to process MsgManageMAYAName", "error", err)
	}
	return result, err
}

func (h ManageMAYANameHandler) validate(ctx cosmos.Context, msg MsgManageMAYAName) error {
	version := h.mgr.GetVersion()
	if version.GTE(semver.MustParse("0.1.0")) {
		return h.validateV1(ctx, msg)
	}
	return errBadVersion
}

func (h ManageMAYANameHandler) validateNameV1(n string) error {
	// validate MAYAName
	if len(n) > 30 {
		return errors.New("MAYAName cannot exceed 30 characters")
	}
	if !IsValidMAYANameV1(n) {
		return errors.New("invalid MAYAName")
	}
	return nil
}

func (h ManageMAYANameHandler) validateV1(ctx cosmos.Context, msg MsgManageMAYAName) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	exists := h.mgr.Keeper().MAYANameExists(ctx, msg.Name)

	if !exists {
		// mayaname doesn't appear to exist, let's validate the name
		if err := h.validateNameV1(msg.Name); err != nil {
			return err
		}
		registrationFee := h.mgr.GetConstants().GetInt64Value(constants.TNSRegisterFee)
		if msg.Coin.Amount.LTE(cosmos.NewUint(uint64(registrationFee))) {
			return fmt.Errorf("not enough funds")
		}
	} else {
		name, err := h.mgr.Keeper().GetMAYAName(ctx, msg.Name)
		if err != nil {
			return err
		}

		// if this mayaname is already owned, check signer has ownership. If
		// expiration is past, allow different user to take ownership
		if !name.Owner.Equals(msg.Signer) && ctx.BlockHeight() <= name.ExpireBlockHeight {
			ctx.Logger().Error("no authorization", "owner", name.Owner)
			return fmt.Errorf("no authorization: owned by %s", name.Owner)
		}

		// ensure user isn't inflating their expire block height artificaially
		if name.ExpireBlockHeight < msg.ExpireBlockHeight {
			return errors.New("cannot artificially inflate expire block height")
		}
	}

	return nil
}

// handle process MsgManageMAYAName
func (h ManageMAYANameHandler) handle(ctx cosmos.Context, msg MsgManageMAYAName) (*cosmos.Result, error) {
	version := h.mgr.GetVersion()
	if version.GTE(semver.MustParse("0.1.0")) {
		return h.handleV1(ctx, msg)
	}
	return nil, errBadVersion
}

// handle process MsgManageMAYAName
func (h ManageMAYANameHandler) handleV1(ctx cosmos.Context, msg MsgManageMAYAName) (*cosmos.Result, error) {
	var err error

	enable, _ := h.mgr.Keeper().GetMimir(ctx, "MAYANames")
	if enable == 0 {
		return nil, fmt.Errorf("MAYANames are currently disabled")
	}

	tn := MAYAName{Name: msg.Name, Owner: msg.Signer, PreferredAsset: common.EmptyAsset}
	exists := h.mgr.Keeper().MAYANameExists(ctx, msg.Name)
	if exists {
		tn, err = h.mgr.Keeper().GetMAYAName(ctx, msg.Name)
		if err != nil {
			return nil, err
		}
	}

	registrationFeePaid := cosmos.ZeroUint()
	fundPaid := cosmos.ZeroUint()

	// check if user is trying to extend expiration
	if !msg.Coin.Amount.IsZero() {
		// check that MAYAName is still valid, can't top up an invalid MAYAName
		if err := h.validateNameV1(msg.Name); err != nil {
			return nil, err
		}
		var addBlocks int64
		// registration fee is for BASEChain addresses only
		if !exists {
			// minus registration fee
			registrationFee := fetchConfigInt64(ctx, h.mgr, constants.TNSRegisterFee)
			msg.Coin.Amount = common.SafeSub(msg.Coin.Amount, cosmos.NewUint(uint64(registrationFee)))
			registrationFeePaid = cosmos.NewUint(uint64(registrationFee))
			addBlocks = h.mgr.GetConstants().GetInt64Value(constants.BlocksPerYear) // registration comes with 1 free year
		}
		feePerBlock := fetchConfigInt64(ctx, h.mgr, constants.TNSFeePerBlock)
		fundPaid = msg.Coin.Amount
		addBlocks += (int64(msg.Coin.Amount.Uint64()) / feePerBlock)
		if tn.ExpireBlockHeight < ctx.BlockHeight() {
			tn.ExpireBlockHeight = ctx.BlockHeight() + addBlocks
		} else {
			tn.ExpireBlockHeight += addBlocks
		}
	}

	// check if we need to reduce the expire time, upon user request
	if msg.ExpireBlockHeight > 0 && msg.ExpireBlockHeight < tn.ExpireBlockHeight {
		tn.ExpireBlockHeight = msg.ExpireBlockHeight
	}

	// check if we need to update the preferred asset
	if !tn.PreferredAsset.Equals(msg.PreferredAsset) && !msg.PreferredAsset.IsEmpty() {
		tn.PreferredAsset = msg.PreferredAsset
	}

	tn.SetAlias(msg.Chain, msg.Address) // update address
	if !msg.Owner.Empty() {
		tn.Owner = msg.Owner // update owner
	}
	h.mgr.Keeper().SetMAYAName(ctx, tn)

	evt := NewEventMAYAName(tn.Name, msg.Chain, msg.Address, registrationFeePaid, fundPaid, tn.ExpireBlockHeight, tn.Owner)
	if err := h.mgr.EventMgr().EmitEvent(ctx, evt); nil != err {
		ctx.Logger().Error("fail to emit MAYAName event", "error", err)
	}

	return &cosmos.Result{}, nil
}
