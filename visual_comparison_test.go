package myqrcode

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
)

// TestVisualComparison generates QR codes for visual comparison with Chrome's style
func TestVisualComparison(t *testing.T) {
	// Create test output directory
	os.Mkdir("visual_tests", 0755)

	// Test 1: Basic Chrome style without logo
	t.Run("ChromeStyleBasic", func(t *testing.T) {
		qr, err := New("https://meet.google.com/abc-defg-hij", High)
		if err != nil {
			t.Fatalf("Failed to create QR code: %v", err)
		}

		err = qr.Encode()
		if err != nil {
			t.Fatalf("Failed to encode QR code: %v", err)
		}

		config := StyleConfig{
			ModuleSize:      10,
			QuietZone:       40,
			RoundedCorners:  true,
			CircularDots:    true,
			BackgroundColor: color.RGBA{255, 255, 255, 255},
			ForegroundColor: color.RGBA{0, 0, 0, 255},
		}

		img, err := qr.ToImage(config)
		if err != nil {
			t.Fatalf("Failed to generate image: %v", err)
		}

		saveVisualTest(t, img, "chrome_style_basic.png")
	})

	// Test 2: Chrome style with dinosaur logo
	t.Run("ChromeStyleWithDinosaur", func(t *testing.T) {
		qr, err := New("https://meet.google.com/abc-defg-hij", High)
		if err != nil {
			t.Fatalf("Failed to create QR code: %v", err)
		}

		// Create better dinosaur logo
		dinosaur := createBetterDinosaurLogo()
		qr.SetLogo(dinosaur, 15)

		err = qr.Encode()
		if err != nil {
			t.Fatalf("Failed to encode QR code: %v", err)
		}

		config := StyleConfig{
			ModuleSize:      10,
			QuietZone:       40,
			RoundedCorners:  true,
			CircularDots:    true,
			BackgroundColor: color.RGBA{255, 255, 255, 255},
			ForegroundColor: color.RGBA{0, 0, 0, 255},
		}

		img, err := qr.ToImage(config)
		if err != nil {
			t.Fatalf("Failed to generate image: %v", err)
		}

		saveVisualTest(t, img, "chrome_style_with_dinosaur.png")
	})

	// Test 3: Different finder pattern styles
	t.Run("FinderPatternStyles", func(t *testing.T) {
		qr, err := New("https://meet.google.com/test", High)
		if err != nil {
			t.Fatalf("Failed to create QR code: %v", err)
		}

		err = qr.Encode()
		if err != nil {
			t.Fatalf("Failed to encode QR code: %v", err)
		}

		// Test different styles
		styles := []struct {
			name     string
			rounded  bool
			circular bool
		}{
			{"square_modules", false, false},
			{"rounded_corners", true, false},
			{"circular_dots", false, true},
			{"chrome_style", true, true},
		}

		for _, style := range styles {
			config := StyleConfig{
				ModuleSize:     8,
				QuietZone:      32,
				RoundedCorners: style.rounded,
				CircularDots:   style.circular,
			}

			img, err := qr.ToImage(config)
			if err != nil {
				t.Fatalf("Failed to generate %s image: %v", style.name, err)
			}

			saveVisualTest(t, img, "style_"+style.name+".png")
		}
	})

	// Test 4: Size variations
	t.Run("SizeVariations", func(t *testing.T) {
		qr, err := New("https://meet.google.com/size-test", High)
		if err != nil {
			t.Fatalf("Failed to create QR code: %v", err)
		}

		err = qr.Encode()
		if err != nil {
			t.Fatalf("Failed to encode QR code: %v", err)
		}

		sizes := []int{6, 8, 10, 12}
		for _, size := range sizes {
			config := StyleConfig{
				ModuleSize:     size,
				QuietZone:      size * 4,
				RoundedCorners: true,
				CircularDots:   true,
			}

			img, err := qr.ToImage(config)
			if err != nil {
				t.Fatalf("Failed to generate size %d image: %v", size, err)
			}

			saveVisualTest(t, img, fmt.Sprintf("size_%dpx.png", size))
		}
	})
}

// TestTargetReplication attempts to replicate the exact target image
func TestTargetReplication(t *testing.T) {
	os.Mkdir("visual_tests", 0755)
	os.Mkdir("visual_tests/target_replication", 0755)

	// Try to match the exact style from qrcode_meet.google.com.png
	qr, err := New("https://meet.google.com/abc-defg-hij", High)
	if err != nil {
		t.Fatalf("Failed to create QR code: %v", err)
	}

	// Add dinosaur logo
	dinosaur := createBetterDinosaurLogo()
	qr.SetLogo(dinosaur, 18) // Adjust size to match target

	err = qr.Encode()
	if err != nil {
		t.Fatalf("Failed to encode QR code: %v", err)
	}

	// Try to match exact visual parameters
	config := StyleConfig{
		ModuleSize:      8,  // Smaller modules like target
		QuietZone:       32, // Appropriate quiet zone
		RoundedCorners:  true,
		CircularDots:    true,
		BackgroundColor: color.RGBA{255, 255, 255, 255},
		ForegroundColor: color.RGBA{0, 0, 0, 255},
	}

	img, err := qr.ToImage(config)
	if err != nil {
		t.Fatalf("Failed to generate target replication: %v", err)
	}

	saveVisualTest(t, img, "target_replication/attempt_1.png")
	t.Log("Generated target replication - compare with qrcode_meet.google.com.png")
}

func saveVisualTest(t *testing.T, img image.Image, filename string) {
	file, err := os.Create("visual_tests/" + filename)
	if err != nil {
		t.Fatalf("Could not create visual test file %s: %v", filename, err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		t.Fatalf("Could not encode visual test image %s: %v", filename, err)
	}

	t.Logf("Generated visual test: visual_tests/%s", filename)
}

func createBetterDinosaurLogo() image.Image {
	// Create a more accurate Chrome dinosaur logo
	size := 48
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Chrome dinosaur pattern (simplified but more accurate)
	// This is a rough approximation - would need exact pixel data for perfect match
	dinosaurPattern := []string{
		"      ████████      ",
		"    ████████████    ",
		"   ██████  ██████   ",
		"  ████████████████  ",
		" ██████████████████ ",
		"████████████████████",
		"████████████████████",
		"██████████████████  ",
		"████████████████    ",
		"██████████████      ",
		"████████████        ",
		"██████████          ",
		"████████            ",
		"██████              ",
		"████                ",
		"██                  ",
	}

	// Convert pattern to image
	scale := 2
	for y, row := range dinosaurPattern {
		for x, char := range row {
			if char == '█' {
				for dy := 0; dy < scale; dy++ {
					for dx := 0; dx < scale; dx++ {
						px := x*scale + dx
						py := y*scale + dy
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