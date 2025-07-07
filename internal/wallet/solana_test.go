package wallet

import (
	"testing"
)

func TestDeriveSolanaKeyPair(t *testing.T) {
	mnemonic := "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about"

	acc0, err := DeriveSolanaKeyPair(mnemonic, 0)
	if err != nil {
		t.Fatal("派生失败:", err)
	}

	acc1, err := DeriveSolanaKeyPair(mnemonic, 1)
	if err != nil {
		t.Fatal("派生失败:", err)
	}

	if acc0.PublicKey.ToBase58() == acc1.PublicKey.ToBase58() {
		t.Fatal("同一助记词在不同索引下应派生出不同地址")
	}

	t.Logf("Address0: %s", acc0.PublicKey.ToBase58())
	t.Logf("Address1: %s", acc1.PublicKey.ToBase58())
}
