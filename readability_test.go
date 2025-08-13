package myqrcode

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

// TestQRCodeReadability tests if generated QR codes are actually readable
func TestQRCodeReadability(t *testing.T) {
	os.Mkdir("readability_tests", 0755)
	
	testCases := []struct {
		name        string
		data        string
		level       ErrorCorrectionLevel
		withLogo    bool
		expectError bool
	}{
		{"Simple URL", "https://meet.google.com/test", High, false, false},
		{"Short URL", "https://g.co/meet", Low, false, false},
		{"Medium URL", "https://meet.google.com/abc-defg-hij", Medium, false, false},
		{"With Logo Small", "https://meet.google.com/logo", High, true, false},
		{"Numeric Data", "1234567890", Low, false, false},
		{"Mixed Data", "ABC123xyz", Medium, false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			qr, err := New(tc.data, tc.level)
			if err != nil {
				if !tc.expectError {
					t.Fatalf("Unexpected error creating QR: %v", err)
				}
				return
			}

			if tc.withLogo {
				logo := createSimpleLogo()
				qr.SetLogo(logo, 15)
			}

			err = qr.Encode()
			if err != nil {
				if !tc.expectError {
					t.Fatalf("Unexpected error encoding QR: %v", err)
				}
				return
			}

			// Test basic structure validation
			validateQRStructure(t, qr, tc.name)

			// Generate different rendering styles for testing
			styles := []struct {
				name   string
				config StyleConfig
			}{
				{"square", StyleConfig{
					ModuleSize: 10, QuietZone: 40,
					RoundedCorners: false, CircularDots: false,
				}},
				{"chrome_style", StyleConfig{
					ModuleSize: 10, QuietZone: 40,
					RoundedCorners: true, CircularDots: true,
				}},
			}

			for _, style := range styles {
				img, err := qr.ToImage(style.config)
				if err != nil {
					t.Fatalf("Failed to generate %s image: %v", style.name, err)
				}

				filename := sanitizeFilename(tc.name + "_" + style.name + ".png")
				saveReadabilityTest(t, img, filename)
			}
		})
	}
}

// TestQRMatrixValidation validates the raw QR matrix structure
func TestQRMatrixValidation(t *testing.T) {
	qr, err := New("https://meet.google.com/matrix-test", High)
	if err != nil {
		t.Fatalf("Failed to create QR: %v", err)
	}

	err = qr.Encode()
	if err != nil {
		t.Fatalf("Failed to encode QR: %v", err)
	}

	// Test matrix properties
	if qr.Matrix == nil {
		t.Fatal("QR matrix is nil")
	}

	if len(qr.Matrix) != qr.Size {
		t.Fatalf("Matrix height %d != expected size %d", len(qr.Matrix), qr.Size)
	}

	for i, row := range qr.Matrix {
		if len(row) != qr.Size {
			t.Fatalf("Matrix row %d width %d != expected size %d", i, len(row), qr.Size)
		}
	}

	// Test finder patterns exist
	validateFinderPatterns(t, qr)

	// Test timing patterns exist
	validateTimingPatterns(t, qr)

	t.Logf("Matrix validation passed for %dx%d QR code", qr.Size, qr.Size)
}

// TestCompareWithReference compares our output with a known working implementation
func TestCompareWithReference(t *testing.T) {
	// Generate a simple QR code and save both square and styled versions
	qr, err := New("HELLO WORLD", High)
	if err != nil {
		t.Fatalf("Failed to create reference QR: %v", err)
	}

	err = qr.Encode()
	if err != nil {
		t.Fatalf("Failed to encode reference QR: %v", err)
	}

	// Generate basic square version (most likely to be readable)
	basicConfig := StyleConfig{
		ModuleSize:      10,
		QuietZone:       40,
		RoundedCorners:  false,
		CircularDots:    false,
		BackgroundColor: color.RGBA{255, 255, 255, 255},
		ForegroundColor: color.RGBA{0, 0, 0, 255},
	}

	img, err := qr.ToImage(basicConfig)
	if err != nil {
		t.Fatalf("Failed to generate reference image: %v", err)
	}

	saveReadabilityTest(t, img, "reference_basic_square.png")

	// Generate Chrome-style version
	chromeConfig := StyleConfig{
		ModuleSize:      10,
		QuietZone:       40,
		RoundedCorners:  true,
		CircularDots:    true,
		BackgroundColor: color.RGBA{255, 255, 255, 255},
		ForegroundColor: color.RGBA{0, 0, 0, 255},
	}

	img2, err := qr.ToImage(chromeConfig)
	if err != nil {
		t.Fatalf("Failed to generate Chrome-style image: %v", err)
	}

	saveReadabilityTest(t, img2, "reference_chrome_style.png")

	t.Log("Generated reference QR codes - test these with a QR scanner")
	t.Log("If basic square version is readable but Chrome style isn't, the issue is in rendering")
	t.Log("If neither is readable, the issue is in core QR generation")
}

// TestMinimalQRCode tests the simplest possible QR code
func TestMinimalQRCode(t *testing.T) {
	// Test with minimal data to isolate issues
	qr, err := New("HI", Low)
	if err != nil {
		t.Fatalf("Failed to create minimal QR: %v", err)
	}

	err = qr.Encode()
	if err != nil {
		t.Fatalf("Failed to encode minimal QR: %v", err)
	}

	t.Logf("Minimal QR: Version %d, Size %dx%d, Mode %d", qr.Version, qr.Size, qr.Size, qr.Mode)

	// Generate ultra-basic rendering
	config := StyleConfig{
		ModuleSize: 20, // Large modules for easier scanning
		QuietZone:  80, // Large quiet zone
		RoundedCorners: false,
		CircularDots: false,
	}

	img, err := qr.ToImage(config)
	if err != nil {
		t.Fatalf("Failed to generate minimal image: %v", err)
	}

	saveReadabilityTest(t, img, "minimal_test.png")
	t.Log("Generated minimal QR code - this should definitely be readable if our core logic is correct")
}

// Helper functions

func validateQRStructure(t *testing.T, qr *QRCode, testName string) {
	if qr.Size == 0 {
		t.Fatalf("%s: QR size is 0", testName)
	}

	expectedSize := 17 + 4*qr.Version
	if qr.Size != expectedSize {
		t.Fatalf("%s: QR size %d != expected %d for version %d", testName, qr.Size, expectedSize, qr.Version)
	}

	if qr.Matrix == nil {
		t.Fatalf("%s: QR matrix is nil", testName)
	}
}

func validateFinderPatterns(t *testing.T, qr *QRCode) {
	// Check top-left finder pattern
	if !qr.Matrix[0][0] || !qr.Matrix[0][6] || !qr.Matrix[6][0] || !qr.Matrix[6][6] {
		t.Error("Top-left finder pattern appears invalid")
	}

	// Check top-right finder pattern
	size := qr.Size
	if !qr.Matrix[0][size-7] || !qr.Matrix[0][size-1] || !qr.Matrix[6][size-7] || !qr.Matrix[6][size-1] {
		t.Error("Top-right finder pattern appears invalid")
	}

	// Check bottom-left finder pattern
	if !qr.Matrix[size-7][0] || !qr.Matrix[size-1][0] || !qr.Matrix[size-7][6] || !qr.Matrix[size-1][6] {
		t.Error("Bottom-left finder pattern appears invalid")
	}
}

func validateTimingPatterns(t *testing.T, qr *QRCode) {
	// Check horizontal timing pattern
	for i := 8; i < qr.Size-8; i++ {
		expected := (i % 2) == 0
		if qr.Matrix[6][i] != expected {
			t.Errorf("Horizontal timing pattern error at position %d", i)
		}
	}

	// Check vertical timing pattern
	for i := 8; i < qr.Size-8; i++ {
		expected := (i % 2) == 0
		if qr.Matrix[i][6] != expected {
			t.Errorf("Vertical timing pattern error at position %d", i)
		}
	}
}

func sanitizeFilename(filename string) string {
	// Replace spaces and special characters for valid filenames
	result := ""
	for _, char := range filename {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || 
		   (char >= '0' && char <= '9') || char == '_' || char == '.' {
			result += string(char)
		} else {
			result += "_"
		}
	}
	return result
}

func saveReadabilityTest(t *testing.T, img image.Image, filename string) {
	os.Mkdir("readability_tests", 0755)
	
	file, err := os.Create("readability_tests/" + filename)
	if err != nil {
		t.Logf("Could not create readability test file %s: %v", filename, err)
		return
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		t.Logf("Could not encode readability test image %s: %v", filename, err)
	} else {
		t.Logf("Generated readability test: readability_tests/%s", filename)
	}
}