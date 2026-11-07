package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type activationWindow struct {
	*walk.MainWindow
	userIDEdit  *walk.LineEdit
	keyEdit     *walk.LineEdit
	statusLabel *walk.Label
	exeDir      string
}

func runUI() {
	exe, err := os.Executable()
	if err != nil {
		walk.MsgBox(nil, "Easiactivation", "Cannot determine EXE path: "+err.Error(), walk.MsgBoxIconError)
		return
	}
	exeDir := filepath.Dir(exe)

	// Placement check FIRST (same behavior as old tool)
	if !CheckPlacement(exeDir) {
		walk.MsgBox(nil, "Easiactivation", "Place me at the installation directory!", walk.MsgBoxOK|walk.MsgBoxIconError)
		return
	}

	w := &activationWindow{exeDir: exeDir}

	mainWin := MainWindow{
		AssignTo: &w.MainWindow,
		Title:    "Easiactivation",
		Size:     Size{Width: 480, Height: 260},
		Layout:   VBox{MarginsZero: false},
		Children: []Widget{
			Label{Text: "Easiactivation"},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{Text: "User id", MinSize: Size{Width: 120}},
					HSpacer{},
					PushButton{
						Text: "Generate",
						OnClicked: func() {
							serial := GetBiosSerial()
							uid := DeriveUserID(serial)
							w.userIDEdit.SetText(uid)
							w.statusLabel.SetText("User id generated from this PC. Select it and press Ctrl+C to copy.")
						},
					},
				},
			},
			LineEdit{AssignTo: &w.userIDEdit, ReadOnly: false},
			Label{Text: "Activation key"},
			LineEdit{AssignTo: &w.keyEdit},
			PushButton{
				Text: "VERIFY && ACTIVATE",
				OnClicked: func() {
					userID := strings.TrimSpace(w.userIDEdit.Text())
					key := strings.TrimSpace(w.keyEdit.Text())
					if len(digitsOnly(userID)) != 20 {
						w.statusLabel.SetText("User id must contain exactly 20 digits.")
						return
					}
					if !VerifyActivationKey(userID, key) {
						w.statusLabel.SetText("Invalid activation key. Please check and try again.")
						return
					}
					if err := Activate(w.exeDir); err != nil {
						w.statusLabel.SetText("Activation failed: " + err.Error())
						return
					}
					walk.MsgBox(w, "Easiactivation", "Activation completed successfully.", walk.MsgBoxOK|walk.MsgBoxIconInformation)
					w.Close()
				},
			},
			Label{AssignTo: &w.statusLabel, Text: ""},
		},
	}

	if _, err := mainWin.Run(); err != nil {
		walk.MsgBox(nil, "Easiactivation", "UI error: "+err.Error(), walk.MsgBoxIconError)
	}
}

// digitsOnly returns only 0-9 characters from s.
func digitsOnly(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}