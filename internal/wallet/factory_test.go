package wallet

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"log"
	"math/big"
	"os"
	"testing"
)

func TestGetAllEthereumWalletFromMnemonic(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("No .env file found")
	}

	mnemonic := os.Getenv("TEST_MNEMONIC")
	infura_api_key := os.Getenv("INFURA_API_KEY")
	fmt.Println(mnemonic)
	fmt.Println(infura_api_key)

	// 连接Sepolia RPC
	client, err := ethclient.Dial(fmt.Sprintf("https://sepolia.infura.io/v3/%s", infura_api_key))
	if err != nil {
		log.Fatal(err)
	}

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	// 批量派生地址并查询余额
	for i := 0; i < 20; i++ {
		// 标准路径 m/44'/60'/0'/0/{i}
		path := hdwallet.MustParseDerivationPath(fmt.Sprintf("m/44'/60'/0'/0/%d", i))
		account, err := wallet.Derive(path, false)
		if err != nil {
			log.Printf("派生地址失败 index=%d: %v", i, err)
			continue
		}

		balance, err := client.BalanceAt(context.Background(), account.Address, nil)
		if err != nil {
			log.Printf("查询余额失败 地址=%s: %v", account.Address.Hex(), err)
			continue
		}

		if balance.Cmp(big.NewInt(0)) > 0 {
			fmt.Printf("账户 %d 地址: %s 余额: %s wei\n", i, account.Address.Hex(), balance.String())
		}
	}
}
