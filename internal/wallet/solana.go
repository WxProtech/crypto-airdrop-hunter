package wallet

import (
	"encoding/hex"
	"fmt"

	"crypto/sha512"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/pbkdf2"
)

// DeriveSolanaKeyPair 从助记词派生 Solana 地址（基于 ed25519）
func DeriveSolanaKeyPair(mnemonic string, index uint32) (types.Account, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return types.Account{}, fmt.Errorf("invalid mnemonic")
	}

	seed := bip39.NewSeed(mnemonic, "")

	// incorporate the index into the derivation so different indices yield
	// different key pairs. The salt format is not BIP-44 compliant but is
	// sufficient for deterministic unique keys.
	salt := []byte{
		'e', 'd', '2', '5', '5', '1', '9', ' ', 's', 'e', 'e', 'd',
		byte(index >> 24), byte(index >> 16), byte(index >> 8), byte(index),
	}

	derived := pbkdf2.Key(seed[:], salt, 2048, 32, sha512.New)

	// # 生成solana的地址
	account, err := types.AccountFromSeed(derived)
	if err != nil {
		return types.Account{}, err
	}

	return account, nil
}

func PrintSolanaAddress(mnemonic string, index uint32) error {
	acc, err := DeriveSolanaKeyPair(mnemonic, index)
	if err != nil {
		return err
	}
	fmt.Println("Address:", acc.PublicKey.ToBase58())
	fmt.Println("Private:", hex.EncodeToString(acc.PrivateKey))
	return nil
}
