package main

import _ "embed"

//go:embed assets/Cvte.Platform.Basic.dll
var dll1Bytes []byte

//go:embed assets/SWCoreSharp.SWAuthorization.SWAuthClients.dll
var dll2Bytes []byte

//go:embed assets/pubkey.bin
var _embedCheckPubkey []byte // referenced by crypto.go via its own embed; keep file non-empty