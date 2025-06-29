package wallet

import (
	"testing"
)

func TestDeriveSolanaKeyPair(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"
	account, err := DeriveSolanaKeyPair(mnemonic, 0)
	if err != nil {
		t.Fatal("派生失败:", err)
	}
	t.Logf("Address: %s", account.PublicKey.ToBase58())
	t.Logf("Private: %x", account.PrivateKey)
}
