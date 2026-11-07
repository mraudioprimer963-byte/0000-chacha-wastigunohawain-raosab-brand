package main

import "github.com/lxn/walk"

func main() {
	// Windows GUI app — no console window
	walk.AppendToWalkLog = false
	runUI()
}