package main

import (
	"crypto/ed25519"
	_ "embed"
	"encoding/base32"
	"strings"
)

//go:embed assets/pubkey.bin
var pubKeyBytes []byte

// VerifyActivationKey takes the User ID exactly as displayed (with dashes)
// and the pasted Activation key. Returns true iff the signature verifies.
func VerifyActivationKey(userID, activationKey string) bool {
	if len(pubKeyBytes) != ed25519.PublicKeySize {
		return false
	}
	pub := ed25519.PublicKey(pubKeyBytes)

	// normalize key: uppercase, strip dashes/spaces
	clean := strings.ToUpper(activationKey)
	clean = strings.NewReplacer("-", "", " ", "", "\t", "", "\r", "", "\n", "").Replace(clean)

	sig, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(clean)
	if err != nil || len(sig) != ed25519.SignatureSize {
		return false
	}
	return ed25519.Verify(pub, []byte(userID), sig)
}