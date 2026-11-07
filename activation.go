package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	dll1Name = "Cvte.Platform.Basic.dll"
	dll2Name = "SWCoreSharp.SWAuthorization.SWAuthClients.dll"
)

// Activate writes ONLY the embedded DLLs directly into the Main\ folder next to the EXE.
func Activate(exeDir string) error {
	mainDir := filepath.Join(exeDir, "Main")
	if err := os.MkdirAll(mainDir, 0755); err != nil {
		return err
	}

	wroteSomething := false

	// Write DLL1 if it was embedded
	if dll1Bytes != nil && len(dll1Bytes) > 0 {
		if err := os.WriteFile(filepath.Join(mainDir, dll1Name), dll1Bytes, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", dll1Name, err)
		}
		wroteSomething = true
	}

	// Write DLL2 if it was embedded
	if dll2Bytes != nil && len(dll2Bytes) > 0 {
		if err := os.WriteFile(filepath.Join(mainDir, dll2Name), dll2Bytes, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", dll2Name, err)
		}
		wroteSomething = true
	}

	if !wroteSomething {
		return fmt.Errorf("no DLLs were embedded in this EXE")
	}

	return nil
}

// CheckPlacement verifies the EXE sits next to a Main\ folder that contains at least ONE of the target DLLs.
func CheckPlacement(exeDir string) bool {
	mainDir := filepath.Join(exeDir, "Main")
	if _, err := os.Stat(mainDir); err != nil {
		return false // Main folder doesn't exist
	}

	// If EITHER DLL exists in the Main folder, placement is valid.
	if _, err := os.Stat(filepath.Join(mainDir, dll1Name)); err == nil {
		return true
	}
	if _, err := os.Stat(filepath.Join(mainDir, dll2Name)); err == nil {
		return true
	}

	return false
}