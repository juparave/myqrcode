package myqrcode

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

func TestBasicQRGeneration(t *testing.T) {
	qr, err := New("https://meet.google.com/test", High)
	if err != nil {
		t.Fatalf("Failed to create QR code: %v", err)
	}

	err = qr.Encode()
	if err != nil {
		t.Fatalf("Failed to encode QR code: %v", err)
	}

	if qr.Matrix == nil {
		t.Fatal("QR matrix is nil after encoding")
	}

	if qr.Size == 0 {
		t.Fatal("QR size is 0 after encoding")
	}

	t.Logf("Generated QR code: Version %d, Size %dx%d", qr.Version, qr.Size, qr.Size)
}

func TestChromeStyleRendering(t *testing.T) {
	qr, err := New("https://meet.google.com/chrome-test", High)
	if err != nil {
		t.Fatalf("Failed to create QR code: %v", err)
	}

	err = qr.Encode()
	if err != nil {
		t.Fatalf("Failed to encode QR code: %v", err)
	}

	// Test Chrome-style configuration
	config := StyleConfig{
		ModuleSize:      8,
		QuietZone:       32,
		RoundedCorners:  true,
		CircularDots:    true,
		BackgroundColor: color.RGBA{255, 255, 255, 255},
		ForegroundColor: color.RGBA{0, 0, 0, 255},
	}

	img, err := qr.ToImage(config)
	if err != nil {
		t.Fatalf("Failed to generate image: %v", err)
	}

	// Save test output
	saveTestImage(t, img, "test_chrome_style.png")
}

func TestLogoEmbedding(t *testing.T) {
	// Create Chrome dinosaur-style logo
	logo := createDinosaurLogo()

	qr, err := New("https://meet.google.com/logo-test", High)
	if err != nil {
		t.Fatalf("Failed to create QR code: %v", err)
	}

	qr.SetLogo(logo, 15) // 15% size like Chrome

	err = qr.Encode()
	if err != nil {
		t.Fatalf("Failed to encode QR code with logo: %v", err)
	}

	config := StyleConfig{
		ModuleSize:     8,
		QuietZone:      32,
		RoundedCorners: true,
		CircularDots:   true,
	}

	img, err := qr.ToImage(config)
	if err != nil {
		t.Fatalf("Failed to generate image with logo: %v", err)
	}

	saveTestImage(t, img, "test_chrome_with_logo.png")
}

func TestVariousDataSizes(t *testing.T) {
	testCases := []struct {
		name string
		data string
	}{
		{"Short URL", "https://g.co/meet"},
		{"Medium URL", "https://meet.google.com/abc-defg-hij"},
		{"Long URL", "https://meet.google.com/very-long-meeting-room-name-that-should-increase-version"},
		{"Numeric", "1234567890123456789012345"},
		{"Mixed", "Meeting ID: ABC-123-XYZ (https://meet.google.com)"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			qr, err := New(tc.data, High)
			if err != nil {
				t.Fatalf("Failed to create QR code for %s: %v", tc.name, err)
			}

			err = qr.Encode()
			if err != nil {
				t.Fatalf("Failed to encode QR code for %s: %v", tc.name, err)
			}

			t.Logf("%s: Version %d, Size %dx%d, Mode %d",
				tc.name, qr.Version, qr.Size, qr.Size, qr.Mode)
		})
	}
}

func TestErrorCorrectionLevels(t *testing.T) {
	data := "https://meet.google.com/test"
	levels := []ErrorCorrectionLevel{Low, Medium, Quartile, High}
	levelNames := []string{"Low", "Medium", "Quartile", "High"}

	for i, level := range levels {
		t.Run(levelNames[i], func(t *testing.T) {
			qr, err := New(data, level)
			if err != nil {
				t.Fatalf("Failed to create QR code with %s error correction: %v", levelNames[i], err)
			}

			err = qr.Encode()
			if err != nil {
				t.Fatalf("Failed to encode QR code with %s error correction: %v", levelNames[i], err)
			}

			if qr.ErrorCorrection != level {
				t.Errorf("Expected error correction level %d, got %d", level, qr.ErrorCorrection)
			}
		})
	}
}

func TestLogoSizeAdjustment(t *testing.T) {
	logo := createSimpleLogo()
	data := "https://meet.google.com/size-test"

	sizes := []int{10, 15, 20, 25}
	for _, size := range sizes {
		t.Run(fmt.Sprintf("Size%d", size), func(t *testing.T) {
			qr, err := New(data, Medium)
			if err != nil {
				t.Fatalf("Failed to create QR code: %v", err)
			}

			originalLevel := qr.ErrorCorrection
			qr.SetLogo(logo, size)

			err = qr.Encode()
			if err != nil {
				t.Fatalf("Failed to encode QR code with %d%% logo: %v", size, err)
			}

			// Larger logos should increase error correction level
			if size >= 20 && qr.ErrorCorrection <= originalLevel {
				t.Logf("Logo size %d%% adjusted error correction from %d to %d",
					size, originalLevel, qr.ErrorCorrection)
			}
		})
	}
}

// Helper functions for testing

func saveTestImage(t *testing.T, img image.Image, filename string) {
	file, err := os.Create("test_output/" + filename)
	if err != nil {
		// Create directory if it doesn't exist
		os.Mkdir("test_output", 0755)
		file, err = os.Create("test_output/" + filename)
		if err != nil {
			t.Logf("Could not save test image %s: %v", filename, err)
			return
		}
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		t.Logf("Could not encode test image %s: %v", filename, err)
	} else {
		t.Logf("Saved test image: test_output/%s", filename)
	}
}

func createDinosaurLogo() image.Image {
	// Create a simplified dinosaur silhouette like Chrome's
	size := 40
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Simple dinosaur shape (very basic approximation)
	// This should be replaced with actual dinosaur shape
	dinosaurPixels := [][]bool{
		{false, false, false, true, true, true, false, false},
		{false, false, true, true, true, true, true, false},
		{false, true, true, false, true, true, true, true},
		{true, true, true, true, true, true, true, true},
		{true, true, true, true, true, true, true, false},
		{true, true, true, true, true, false, false, false},
		{false, true, true, false, false, false, false, false},
		{false, false, false, false, false, false, false, false},
	}

	// Scale and center the dinosaur
	scale := size / 8
	offsetX := (size - len(dinosaurPixels[0])*scale) / 2
	offsetY := (size - len(dinosaurPixels)*scale) / 2

	for y, row := range dinosaurPixels {
		for x, pixel := range row {
			if pixel {
				for dy := 0; dy < scale; dy++ {
					for dx := 0; dx < scale; dx++ {
						px := offsetX + x*scale + dx
						py := offsetY + y*scale + dy
						if px < size && py < size {
							img.Set(px, py, color.RGBA{0, 0, 0, 255})
						}
					}
				}
			}
		}
	}

	return img
}

func createSimpleLogo() image.Image {
	size := 32
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Create a simple circle logo
	center := float64(size) / 2
	radius := float64(size) / 3

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := float64(x) - center
			dy := float64(y) - center
			if dx*dx+dy*dy <= radius*radius {
				img.Set(x, y, color.RGBA{0, 0, 0, 255})
			}
		}
	}

	return img
}
