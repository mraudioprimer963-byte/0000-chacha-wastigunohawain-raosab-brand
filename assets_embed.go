package main

import "embed"

//go:embed assets/*.dll
var embeddedAssets embed.FS

var (
	dll1Bytes []byte
	dll2Bytes []byte
)

func loadEmbeddedDLLs() {
	dll1Bytes, _ = embeddedAssets.ReadFile("assets/Cvte.Platform.Basic.dll")
	dll2Bytes, _ = embeddedAssets.ReadFile("assets/SWCoreSharp.SWAuthorization.SWAuthClients.dll")
}