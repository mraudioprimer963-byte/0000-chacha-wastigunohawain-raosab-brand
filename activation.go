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

func Activate(exeDir string) error {
	mainDir := filepath.Join(exeDir, "Main")
	if err := os.MkdirAll(mainDir, 0755); err != nil { return err }
	wroteSomething := false
	if len(dll1Bytes) > 0 {
		if err := os.WriteFile(filepath.Join(mainDir, dll1Name), dll1Bytes, 0644); err != nil { return fmt.Errorf("err dll1: %w", err) }
		wroteSomething = true
	}
	if len(dll2Bytes) > 0 {
		if err := os.WriteFile(filepath.Join(mainDir, dll2Name), dll2Bytes, 0644); err != nil { return fmt.Errorf("err dll2: %w", err) }
		wroteSomething = true
	}
	if !wroteSomething { return fmt.Errorf("no DLLs embedded") }
	return nil
}

func CheckPlacement(exeDir string) bool {
	mainDir := filepath.Join(exeDir, "Main")
	if _, err := os.Stat(mainDir); err != nil { return false }
	if _, err := os.Stat(filepath.Join(mainDir, dll1Name)); err == nil { return true }
	if _, err := os.Stat(filepath.Join(mainDir, dll2Name)); err == nil { return true }
	return false
}