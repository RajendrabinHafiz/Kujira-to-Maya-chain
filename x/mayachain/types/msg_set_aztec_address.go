package types

import (
	"gitlab.com/mayachain/mayanode/common"
	"gitlab.com/mayachain/mayanode/common/cosmos"
)

// NewMsgSetAztecAddress is a constructor function for NewMsgAddNodeKeys
func NewMsgSetAztecAddress(aztecAddress common.Address, signer cosmos.AccAddress) *MsgSetAztecAddress {
	return &MsgSetAztecAddress{
		AztecAddress: aztecAddress,
		Signer:       signer,
	}
}

// Route should return the router key of the module
func (m *MsgSetAztecAddress) Route() string { return RouterKey }

// Type should return the action
func (m MsgSetAztecAddress) Type() string { return "set_aztec_address" }

// ValidateBasic runs stateless checks on the message
func (m *MsgSetAztecAddress) ValidateBasic() error {
	if m.Signer.Empty() {
		return cosmos.ErrInvalidAddress(m.Signer.String())
	}

	if !m.AztecAddress.IsChain(common.AZTECChain) {
		return cosmos.ErrInvalidAddress(m.AztecAddress.String())
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (m *MsgSetAztecAddress) GetSignBytes() []byte {
	return cosmos.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgSetAztecAddress) GetSigners() []cosmos.AccAddress {
	return []cosmos.AccAddress{m.Signer}
}
