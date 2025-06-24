package wallet

import (
	"fmt"
	"github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
)

// GenerateMnemonic creates a new BIP39 mnemonic phrase
func GenerateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}

// DeriveEVMAddress derives an EVM-compatible address from the mnemonic and index
func DeriveEVMAddress(mnemonic string, index int) (string, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return "", err
	}
	path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", index))
	account, err := wallet.Derive(path, false)
	if err != nil {
		return "", err
	}
	return account.Address.Hex(), nil
}
