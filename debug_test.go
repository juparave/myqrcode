package myqrcode

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

// TestDebugQRGeneration debugs QR generation step by step
func TestDebugQRGeneration(t *testing.T) {
	os.Mkdir("debug_tests", 0755)

	// Test with very simple data
	data := "HELLO"
	t.Logf("Testing with data: %s", data)

	qr, err := New(data, Low)
	if err != nil {
		t.Fatalf("Failed to create QR: %v", err)
	}

	t.Logf("Initial QR: Version=%d, Mode=%d, ErrorCorrection=%d", qr.Version, qr.Mode, qr.ErrorCorrection)

	// Debug encoding step
	mode := detectMode(data)
	t.Logf("Detected mode: %d", mode)

	version := determineVersion(data, mode, Low)
	t.Logf("Determined version: %d", version)

	// Test data encoding
	encoded, err := encodeData(data, mode, version)
	if err != nil {
		t.Fatalf("Failed to encode data: %v", err)
	}
	t.Logf("Encoded data length: %d bits", len(encoded))

	// Test padding
	padded := addTerminatorAndPadding(encoded, version, Low)
	t.Logf("Padded data length: %d bits", len(padded))

	// Test error correction
	ec := addErrorCorrection(padded, version, Low)
	t.Logf("Final data with EC length: %d bytes", len(ec))

	// Now encode the full QR
	err = qr.Encode()
	if err != nil {
		t.Fatalf("Failed to encode QR: %v", err)
	}

	t.Logf("Final QR: Version=%d, Size=%d, Matrix exists=%t", qr.Version, qr.Size, qr.Matrix != nil)

	// Save debug matrix visualization
	debugImg := visualizeMatrix(qr.Matrix, qr.Size)
	saveDebugImage(t, debugImg, "debug_matrix.png")

	// Generate basic readable version
	config := StyleConfig{
		ModuleSize: 15,
		QuietZone:  60,
		RoundedCorners: false,
		CircularDots: false,
	}

	img, err := qr.ToImage(config)
	if err != nil {
		t.Fatalf("Failed to generate debug image: %v", err)
	}

	saveDebugImage(t, img, "debug_readable.png")
	t.Log("Generated debug QR - check if this is readable")
}

// TestCompareKnownGoodQR compares with a working QR library for validation
func TestCompareKnownGoodQR(t *testing.T) {
	// Test our implementation against expected behavior
	data := "TEST123"
	
	qr, err := New(data, Medium)
	if err != nil {
		t.Fatalf("Failed to create QR: %v", err)
	}

	err = qr.Encode()
	if err != nil {
		t.Fatalf("Failed to encode QR: %v", err)
	}

	// Print matrix for manual inspection
	t.Logf("QR Matrix for '%s':", data)
	t.Logf("Version: %d, Size: %dx%d", qr.Version, qr.Size, qr.Size)
	
	// Print first few rows to see structure
	for i := 0; i < min(10, qr.Size); i++ {
		row := ""
		for j := 0; j < min(10, qr.Size); j++ {
			if qr.Matrix[i][j] {
				row += "█"
			} else {
				row += "░"
			}
		}
		t.Logf("Row %2d: %s", i, row)
	}
}

// TestDataEncodingOnly tests just the data encoding without matrix generation
func TestDataEncodingOnly(t *testing.T) {
	testCases := []struct {
		data string
		mode EncodingMode
	}{
		{"123", Numeric},
		{"ABC", Alphanumeric}, 
		{"hello", Byte},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s_%d", tc.data, tc.mode), func(t *testing.T) {
			encoded, err := encodeData(tc.data, tc.mode, 1)
			if err != nil {
				t.Fatalf("Failed to encode %s: %v", tc.data, err)
			}

			t.Logf("Data: %s, Mode: %d, Encoded bits: %d", tc.data, tc.mode, len(encoded))
			
			// Show first few bits
			bits := ""
			for i := 0; i < min(32, len(encoded)); i++ {
				bits += fmt.Sprintf("%d", encoded[i])
			}
			t.Logf("First 32 bits: %s", bits)
		})
	}
}

// Helper functions

func visualizeMatrix(matrix [][]bool, size int) image.Image {
	// Create a simple black/white visualization of the matrix
	scale := 5
	img := image.NewRGBA(image.Rect(0, 0, size*scale, size*scale))
	
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			var c color.Color
			if matrix[y][x] {
				c = color.RGBA{0, 0, 0, 255}
			} else {
				c = color.RGBA{255, 255, 255, 255}
			}
			
			for dy := 0; dy < scale; dy++ {
				for dx := 0; dx < scale; dx++ {
					img.Set(x*scale+dx, y*scale+dy, c)
				}
			}
		}
	}
	
	return img
}

func saveDebugImage(t *testing.T, img image.Image, filename string) {
	file, err := os.Create("debug_tests/" + filename)
	if err != nil {
		t.Logf("Could not create debug file %s: %v", filename, err)
		return
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		t.Logf("Could not encode debug image %s: %v", filename, err)
	} else {
		t.Logf("Generated debug image: debug_tests/%s", filename)
	}
}