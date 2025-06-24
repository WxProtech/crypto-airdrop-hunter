package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"

	"github.com/WxProtech/crypto-airdrop-hunter/internal/wallet"
)

func main() {
	app := &cli.App{
		Name:  "wallet-cli",
		Usage: "Generate or import mnemonics and derive EVM addresses",
		Commands: []*cli.Command{
			{
				Name:  "generate",
				Usage: "Generate a new mnemonic",
				Action: func(c *cli.Context) error {
					mnemonic, err := wallet.GenerateMnemonic()
					if err != nil {
						return err
					}
					fmt.Println("Generated Mnemonic:", mnemonic)
					return nil
				},
			},
			{
				Name:  "derive",
				Usage: "Derive EVM address from mnemonic",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "mnemonic", Required: true},
					&cli.IntFlag{Name: "index", Value: 0},
				},
				Action: func(c *cli.Context) error {
					addr, err := wallet.DeriveEthereumAddress(c.String("mnemonic"), c.Int("index"))
					if err != nil {
						return err
					}
					fmt.Println("Derived Address:", addr)
					return nil
				},
			},
			{
				Name:  "aes",
				Usage: "use aes algorithm to encrypt or decrypt mnemonic",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "mnemonic",
						Usage:    "Mnemonic to encrypt or decrypt",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "password",
						Usage:    "password",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "mode",
						Usage:    "Operation mode: encrypt or decrypt",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					mnemonic := c.String("mnemonic")
					mode := c.String("mode")
					password := c.String("password")

					if mode != "encrypt" && mode != "decrypt" {
						return fmt.Errorf("invalid mode: %s, must be 'encrypt' or 'decrypt'", mode)
					}

					if mode == "encrypt" {
						encrypted, err := wallet.EncryptMnemonic(mnemonic, password)
						if err != nil {
							return err
						}
						fmt.Println("Encrypted:", encrypted)
					} else {
						decrypted, err := wallet.DecryptMnemonic(mnemonic, password)
						if err != nil {
							return err
						}
						fmt.Println("Decrypted:", decrypted)
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
