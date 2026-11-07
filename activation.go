package main

import (
	"os"
	"path/filepath"
)

// Activate writes the embedded DLLs straight into the Main\ folder next to the EXE.
// No temp directory. Overwrites existing files.
func Activate(exeDir string) error {
	mainDir := filepath.Join(exeDir, "Main")
	if err := os.MkdirAll(mainDir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(mainDir, "Cvte.Platform.Basic.dll"), dll1Bytes, 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(mainDir, "SWCoreSharp.SWAuthorization.SWAuthClients.dll"), dll2Bytes, 0644); err != nil {
		return err
	}
	return nil
}

// CheckPlacement verifies the EXE sits next to a Main\ folder that already
// contains at least one of the two target DLLs (matching the old tool's rule).
func CheckPlacement(exeDir string) bool {
	mainDir := filepath.Join(exeDir, "Main")
	if _, err := os.Stat(mainDir); err != nil {
		return false
	}
	for _, name := range []string{
		"Cvte.Platform.Basic.dll",
		"SWCoreSharp.SWAuthorization.SWAuthClients.dll",
	} {
		if _, err := os.Stat(filepath.Join(mainDir, name)); err == nil {
			return true
		}
	}
	return false
}