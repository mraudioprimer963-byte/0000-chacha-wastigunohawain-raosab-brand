package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base32"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	initFlag := flag.Bool("init", false, "Generate new keypair, print hex, write assets/pubkey.bin")
	flag.Parse()

	if *initFlag {
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			die(err)
		}
		fmt.Println("PRIVATE KEY (hex, 64 chars) — STORE OFFLINE, NEVER COMMIT:")
		fmt.Println(hex.EncodeToString(priv.Seed()))
		fmt.Println()
		fmt.Println("PUBLIC KEY (hex, 64 chars):")
		fmt.Println(hex.EncodeToString(pub))
		fmt.Println()
		if err := os.MkdirAll("assets", 0755); err != nil {
			die(err)
		}
		if err := os.WriteFile("assets/pubkey.bin", []byte(pub), 0644); err != nil {
			die(err)
		}
		fmt.Println("Wrote assets/pubkey.bin (32 raw bytes). Rebuild Easiactivation.exe.")
		return
	}

	// sign mode
	privHex := os.Getenv("EASIAct_PRIVATE_KEY")
	if privHex == "" {
		die(fmt.Errorf("set EASIAct_PRIVATE_KEY env var to the 64-char hex seed"))
	}
	seed, err := hex.DecodeString(strings.TrimSpace(privHex))
	if err != nil || len(seed) != ed25519.SeedSize {
		die(fmt.Errorf("EASIAct_PRIVATE_KEY must be 64 hex chars (32-byte seed)"))
	}
	priv := ed25519.NewKeyFromSeed(seed)

	if flag.NArg() < 1 {
		die(fmt.Errorf("usage: EASIAct_PRIVATE_KEY=<hex> keygen \"1234-5678-9012-3456-7890\""))
	}
	userID := flag.Arg(0)
	sig := ed25519.Sign(priv, []byte(userID))
	encoded := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(sig)
	fmt.Println(encoded)
}

func die(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}