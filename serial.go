package main

import (
	"crypto/sha256"
	"fmt"
	"os/exec"
	"strings"
)

var placeholders = map[string]bool{
	"": true, "TOBEFILLEDBYOEM": true, "TOBEFILLEDBYOEM.": true, "NONE": true,
	"DEFAULTSTRING": true, "SERIALNUMBER": true, "0": true, "00000000": true,
	"NULL": true, "N/A": true, "NA": true,
}

func GetBiosSerial() string {
	if s := tryPowerShell(); s != "" { return s }
	if s := tryWMIC(); s != "" { return s }
	if s := tryMachineGUID(); s != "" { return s }
	if s := tryComputerName(); s != "" { return s }
	return "FALLBACK-UNKNOWN"
}

func tryPowerShell() string {
	cmd := exec.Command("powershell.exe", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command",
		`try{ $b=(Get-CimInstance Win32_BIOS).SerialNumber; $m=(Get-CimInstance Win32_BaseBoard).SerialNumber; $u=(Get-CimInstance Win32_ComputerSystemProduct).UUID; Write-Output ($b+"|"+$m+"|"+$u) }catch{ Write-Output "" }`)
	out, err := cmd.Output()
	if err != nil { return "" }
	return normalizeSerial(string(out))
}

func tryWMIC() string {
	cmd := exec.Command("cmd", "/C", "wmic bios get serialnumber /value")
	out, err := cmd.Output()
	if err != nil { return "" }
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "SerialNumber=") {
			return normalizeSerial(strings.TrimPrefix(line, "SerialNumber="))
		}
	}
	return ""
}

func tryMachineGUID() string {
	cmd := exec.Command("reg", "query", `HKLM\SOFTWARE\Microsoft\Cryptography`, "/v", "MachineGuid")
	out, err := cmd.Output()
	if err != nil { return "" }
	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 3 && fields[0] == "MachineGuid" {
			return normalizeSerial(fields[len(fields)-1])
		}
	}
	return ""
}

func tryComputerName() string {
	cmd := exec.Command("cmd", "/C", "echo %COMPUTERNAME%")
	out, err := cmd.Output()
	if err != nil { return "" }
	name := strings.TrimSpace(string(out))
	if name == "" { return "" }
	return "FALLBACK-" + strings.ToUpper(name)
}

func normalizeSerial(s string) string {
	s = strings.ToUpper(s)
	replacer := strings.NewReplacer(" ", "", "\t", "", "\r", "", "\n", "", "-", "")
	s = replacer.Replace(s)
	if placeholders[s] { return "" }
	return s
}

func DeriveUserID(serial string) string {
	seed := serial + "|EASIACT-2024-v1"
	sum := sha256.Sum256([]byte(seed))
	var n uint64
	for i := 0; i < 8; i++ { n = (n << 8) | uint64(sum[i]) }
	const mod = uint64(10000000000000000000)
	n = n % mod
	digits := fmt.Sprintf("%019d", n)
	extra := sum[8] % 10
	digits = fmt.Sprintf("%d%s", extra, digits)
	return fmt.Sprintf("%s-%s-%s-%s-%s", digits[0:4], digits[4:8], digits[8:12], digits[12:16], digits[16:20])
}