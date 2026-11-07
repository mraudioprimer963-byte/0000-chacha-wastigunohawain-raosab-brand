package main

import "embed"

//go:embed assets/*.dll
var embeddedAssets embed.FS

// We will check these at runtime to see what the user actually put in the assets folder
var (
	dll1Bytes []byte // Cvte.Platform.Basic.dll
	dll2Bytes []byte // SWCoreSharp.SWAuthorization.SWAuthClients.dll
)

func loadEmbeddedDLLs() {
	var err error
	// Try loading DLL1 (won't crash if missing)
	dll1Bytes, err = embeddedAssets.ReadFile("assets/Cvte.Platform.Basic.dll")
	if err != nil {
		dll1Bytes = nil
	}

	// Try loading DLL2 (won't crash if missing)
	dll2Bytes, err = embeddedAssets.ReadFile("assets/SWCoreSharp.SWAuthorization.SWAuthClients.dll")
	if err != nil {
		dll2Bytes = nil
	}
}