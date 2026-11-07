package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

var (
	user32          = syscall.NewLazyDLL("user32.dll")
	procMessageBoxW = user32.NewProc("MessageBoxW")
)

func runUI() {
	loadEmbeddedDLLs()
	exe, _ := os.Executable()
	exeDir := filepath.Dir(exe)

	if !CheckPlacement(exeDir) {
		msg("Place me at the installation directory!\n(Make sure the Main folder contains the target DLL)", "Error")
		return
	}

	uid := DeriveUserID(GetBiosSerial())
	
	inputKey := prompt("User ID: " + uid + "\n\nEnter Activation Key:")
	if inputKey == "" { return }

	if !VerifyActivationKey(uid, inputKey) {
		msg("Invalid activation key.", "Error")
		return
	}

	if err := Activate(exeDir); err != nil {
		msg("Activation failed: "+err.Error(), "Error")
		return
	}
	msg("Activation completed successfully.", "Success")
}

func msg(text, title string) {
	t, _ := syscall.UTF16PtrFromString(text)
	ti, _ := syscall.UTF16PtrFromString(title)
	procMessageBoxW.Call(0, uintptr(unsafe.Pointer(t)), uintptr(unsafe.Pointer(ti)), 0)
}

func prompt(promptText string) string {
	tmpFile := filepath.Join(os.Getenv("TEMP"), "easiact_input.txt")
	os.Remove(tmpFile)
	
	script := fmt.Sprintf(`InputBox("%s", "Easiactivation")`, strings.ReplaceAll(promptText, `"`, `""`))
	vbsFile := filepath.Join(os.Getenv("TEMP"), "easiact.vbs")
	os.WriteFile(vbsFile, []byte(script), 0644)
	
	cmdStr := fmt.Sprintf(`cscript //nologo "%s" > "%s"`, vbsFile, tmpFile)
	c, _ := syscall.UTF16PtrFromString(cmdStr)
	
	cmd := syscall.NewLazyDLL("kernel32.dll").NewProc("WinExec")
	cmd.Call(uintptr(unsafe.Pointer(c)), 0)
	
	for {
		if _, err := os.Stat(tmpFile); err == nil { break }
	}
	out, _ := os.ReadFile(tmpFile)
	os.Remove(tmpFile)
	os.Remove(vbsFile)
	return strings.TrimSpace(string(out))
}