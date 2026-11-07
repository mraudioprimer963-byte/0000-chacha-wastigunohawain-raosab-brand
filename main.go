package main

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	procMessageBoxW  = user32.NewProc("MessageBoxW")
)

func msgBox(text, title string) {
	t, _ := syscall.UTF16PtrFromString(text)
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	procMessageBoxW.Call(0, uintptr(unsafe.Pointer(t)), uintptr(unsafe.Pointer(titlePtr)), 0)
}

func main() {
	exe, _ := os.Executable()
	exeDir := filepath.Dir(exe)

	if !CheckPlacement(exeDir) {
		msgBox("Place me at the installation directory!\n(Make sure the Main folder contains the target DLL)", "Easiactivation")
		return
	}

	// Just a test to see if it opens on Windows 7!
	msgBox("The EXE successfully opened on Windows 7!", "Success")
	
	// ... (UI logic will go here, but let's make sure the app opens first)
	_ = strings.TrimSpace
}