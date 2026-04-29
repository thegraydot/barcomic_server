package barcomic

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

var ErrNoTool = errors.New(
	"keystroke: no injection tool found; install xdotool (X11) or ydotool (Wayland)",
)
var ErrUnsupportedOS = fmt.Errorf("keystroke: unsupported OS %q", runtime.GOOS)
var ErrInvalidInput = errors.New("keystroke: input contains non-digit characters")

// validateInput hard-gates before any exec call
// Barcodes must contain only ASCII digits 0-9
func validateInput(s string) error {
	if len(s) == 0 {
		return ErrInvalidInput
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return ErrInvalidInput
		}
	}
	return nil
}

// TypeBarcode types the barcode string into the currently focused window
// and presses Enter. Returns an error if the platform tool is unavailable
// or if input validation fails.
func TypeBarcode(barcode string) error {
	if err := validateInput(barcode); err != nil {
		return err
	}
	switch runtime.GOOS {
	case "windows":
		args := buildWindowsArgs(barcode)
		return exec.Command(args[0], args[1:]...).Run()
	case "darwin":
		args := buildDarwinArgs(barcode)
		return exec.Command(args[0], args[1:]...).Run()
	case "linux":
		return typeLinux(barcode)
	default:
		return ErrUnsupportedOS
	}
}

func buildWindowsArgs(barcode string) []string {
	script := fmt.Sprintf(
		`Add-Type -AssemblyName System.Windows.Forms; `+
			`[System.Windows.Forms.SendKeys]::SendWait('%s'); `+
			`[System.Windows.Forms.SendKeys]::SendWait('{ENTER}')`,
		barcode,
	)
	return []string{"powershell.exe", "-NoProfile", "-NonInteractive", "-Command", script}
}

func buildDarwinArgs(barcode string) []string {
	script := fmt.Sprintf(
		`tell application "System Events"`+"\n"+
			`  keystroke "%s"`+"\n"+
			`  key code 36`+"\n"+
			`end tell`,
		barcode,
	)
	return []string{"osascript", "-e", script}
}

func buildLinuxX11TypeArgs(barcode string) []string {
	// --clearmodifiers: ensures Shift/Ctrl/etc are not accidentally held
	// --: prevents barcodes starting with '-' being treated as flags
	return []string{"xdotool", "type", "--clearmodifiers", "--", barcode}
}

func buildLinuxX11ReturnArgs() []string {
	return []string{"xdotool", "key", "Return"}
}

func buildLinuxWaylandTypeArgs(barcode string) []string {
	return []string{"ydotool", "type", "--", barcode}
}

func buildLinuxWaylandReturnArgs() []string {
	// evdev keycode 28 = Enter; 28:1 = keydown, 28:0 = keyup
	return []string{"ydotool", "key", "28:1", "28:0"}
}

func runX11(barcode string) error {
	typeArgs := buildLinuxX11TypeArgs(barcode)
	if err := exec.Command(typeArgs[0], typeArgs[1:]...).Run(); err != nil {
		return err
	}
	retArgs := buildLinuxX11ReturnArgs()
	return exec.Command(retArgs[0], retArgs[1:]...).Run()
}

func runWayland(barcode string) error {
	typeArgs := buildLinuxWaylandTypeArgs(barcode)
	if err := exec.Command(typeArgs[0], typeArgs[1:]...).Run(); err != nil {
		return err
	}
	retArgs := buildLinuxWaylandReturnArgs()
	return exec.Command(retArgs[0], retArgs[1:]...).Run()
}

func typeLinux(barcode string) error {
	display := os.Getenv("DISPLAY")
	waylandDisplay := os.Getenv("WAYLAND_DISPLAY")

	xdotoolPath, _ := exec.LookPath("xdotool")
	ydotoolPath, _ := exec.LookPath("ydotool")

	// X11: DISPLAY set and xdotool in PATH
	if display != "" && xdotoolPath != "" {
		return runX11(barcode)
	}

	// Wayland: WAYLAND_DISPLAY set and ydotool in PATH
	if waylandDisplay != "" && ydotoolPath != "" {
		return runWayland(barcode)
	}

	// XWayland fallback: DISPLAY set but xdotool not in PATH, attempt anyway
	if display != "" {
		return runX11(barcode)
	}

	return ErrNoTool
}
