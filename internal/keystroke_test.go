package barcomic

import (
	"strings"
	"testing"
)

func TestValidateInput_ValidDigits(t *testing.T) {
	cases := []string{
		"0",
		"12",
		"759606096015",
		"12345678901234567",
		"01234567890123456789", // 20 digits
	}
	for _, s := range cases {
		t.Run(s, func(t *testing.T) {
			if err := validateInput(s); err != nil {
				t.Fatalf("expected nil error for %q, got %v", s, err)
			}
		})
	}
}

func TestValidateInput_Empty(t *testing.T) {
	if err := validateInput(""); err != ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput for empty string, got %v", err)
	}
}

func TestValidateInput_NonDigit(t *testing.T) {
	cases := []string{"abc", "123abc", "abc123", "1a2"}
	for _, s := range cases {
		t.Run(s, func(t *testing.T) {
			if err := validateInput(s); err != ErrInvalidInput {
				t.Fatalf("expected ErrInvalidInput for %q, got %v", s, err)
			}
		})
	}
}

func TestValidateInput_Symbols(t *testing.T) {
	cases := []string{"+^%~()", "123+456", "^789"}
	for _, s := range cases {
		t.Run(s, func(t *testing.T) {
			if err := validateInput(s); err != ErrInvalidInput {
				t.Fatalf("expected ErrInvalidInput for %q, got %v", s, err)
			}
		})
	}
}

func TestValidateInput_NonASCIIDigits(t *testing.T) {
	// Arabic-Indic digits U+0660..U+0669 look like digits but are not ASCII
	cases := []string{"٠١٢٣", "١٢٣٤٥٦٧٨٩٠"}
	for _, s := range cases {
		t.Run(s, func(t *testing.T) {
			if err := validateInput(s); err != ErrInvalidInput {
				t.Fatalf("expected ErrInvalidInput for non-ASCII digits %q, got %v", s, err)
			}
		})
	}
}

func TestValidateInput_Whitespace(t *testing.T) {
	cases := []string{" ", "123 456", "789\n", "\t123"}
	for _, s := range cases {
		t.Run(s, func(t *testing.T) {
			if err := validateInput(s); err != ErrInvalidInput {
				t.Fatalf("expected ErrInvalidInput for %q, got %v", s, err)
			}
		})
	}
}

func TestBuildWindowsArgs(t *testing.T) {
	args := buildWindowsArgs("759606096015")
	script := strings.Join(args, " ")
	if !strings.Contains(script, "759606096015") {
		t.Fatal("expected barcode in Windows args")
	}
	if !strings.Contains(script, "{ENTER}") {
		t.Fatal("expected {ENTER} in Windows args")
	}
}

func TestBuildWindowsArgs_NoProfile(t *testing.T) {
	args := buildWindowsArgs("123456789012")
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "-NoProfile") {
		t.Fatal("expected -NoProfile in Windows args")
	}
	if !strings.Contains(joined, "-NonInteractive") {
		t.Fatal("expected -NonInteractive in Windows args")
	}
}

func TestBuildDarwinArgs(t *testing.T) {
	args := buildDarwinArgs("759606096015")
	script := strings.Join(args, " ")
	if !strings.Contains(script, "759606096015") {
		t.Fatal("expected barcode in Darwin args")
	}
	if !strings.Contains(script, "key code 36") {
		t.Fatal("expected key code 36 (Return) in Darwin args")
	}
}

func TestBuildLinuxX11TypeArgs(t *testing.T) {
	args := buildLinuxX11TypeArgs("759606096015")
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "--clearmodifiers") {
		t.Fatal("expected --clearmodifiers in X11 type args")
	}
	if !strings.Contains(joined, "-- 759606096015") {
		t.Fatal("expected -- separator before barcode in X11 type args")
	}
}

func TestBuildLinuxX11ReturnArgs(t *testing.T) {
	args := buildLinuxX11ReturnArgs()
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "Return") {
		t.Fatal("expected Return in X11 return args")
	}
	if args[0] != "xdotool" {
		t.Fatalf("expected xdotool as command, got %q", args[0])
	}
}

func TestBuildLinuxWaylandTypeArgs(t *testing.T) {
	args := buildLinuxWaylandTypeArgs("759606096015")
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "-- 759606096015") {
		t.Fatal("expected -- separator before barcode in Wayland type args")
	}
	if args[0] != "ydotool" {
		t.Fatalf("expected ydotool as command, got %q", args[0])
	}
}

func TestBuildLinuxWaylandReturnArgs(t *testing.T) {
	args := buildLinuxWaylandReturnArgs()
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "28:1") || !strings.Contains(joined, "28:0") {
		t.Fatal("expected evdev keydown 28:1 and keyup 28:0 in Wayland return args")
	}
}
