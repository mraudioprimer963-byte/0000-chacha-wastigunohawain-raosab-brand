package main

import (
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

var (
	user32          = syscall.NewLazyDLL("user32.dll")
	procMessageBoxW = user32.NewProc("MessageBoxW")
)

func msgBox(text, title string) {
	t, _ := syscall.UTF16PtrFromString(text)
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	procMessageBoxW.Call(0, uintptr(unsafe.Pointer(t)), uintptr(unsafe.Pointer(titlePtr)), 0)
}

func main() {
	exe, err := os.Executable()
	if err != nil {
		msgBox("Cannot determine EXE path", "Error")
		return
	}
	exeDir := filepath.Dir(exe)

	if !CheckPlacement(exeDir) {
		msgBox("Place me at the installation directory!\n(Make sure the Main folder contains the target DLL)", "Easiactivation")
		return
	}

	msgBox("SUCCESS! The EXE works perfectly on Windows 7.", "Easiactivation")
}