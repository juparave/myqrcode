package myqrcode

import (
	"fmt"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestStyleComparison(t *testing.T) {
	testData := "https://meet.google.com/abc-defg-hij"
	outputDir := "style_comparison_tests"

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// Test configurations
	tests := []struct {
		name   string
		config StyleConfig
		desc   string
	}{
		{
			name: "01_basic_squares",
			config: StyleConfig{
				ModuleSize:      10,
				QuietZone:       4,
				BackgroundColor: color.RGBA{255, 255, 255, 255},
				ForegroundColor: color.RGBA{0, 0, 0, 255},
				ModuleDrawer:    NewSquareModuleDrawer(),
			},
			desc: "Basic square modules (default)",
		},
		{
			name: "02_circles",
			config: StyleConfig{
				ModuleSize:      10,
				QuietZone:       4,
				BackgroundColor: color.RGBA{255, 255, 255, 255},
				ForegroundColor: color.RGBA{0, 0, 0, 255},
				ModuleDrawer:    NewCircleModuleDrawer(),
			},
			desc: "Circular dots (Chrome style)",
		},
		{
			name: "03_rounded_corners",
			config: StyleConfig{
				ModuleSize:      10,
				QuietZone:       4,
				BackgroundColor: color.RGBA{255, 255, 255, 255},
				ForegroundColor: color.RGBA{0, 0, 0, 255},
				ModuleDrawer:    NewRoundedModuleDrawer(1.0),
			},
			desc: "Context-aware rounded corners (Chrome finder patterns)",
		},
		{
			name: "04_gapped_squares_80",
			config: StyleConfig{
				ModuleSize:      10,
				QuietZone:       4,
				BackgroundColor: color.RGBA{255, 255, 255, 255},
				ForegroundColor: color.RGBA{0, 0, 0, 255},
				ModuleDrawer:    NewGappedSquareModuleDrawer(0.8),
			},
			desc: "Gapped squares (80% size)",
		},
		{
			name: "05_gapped_circles_90",
			config: StyleConfig{
				ModuleSize:      10,
				QuietZone:       4,
				BackgroundColor: color.RGBA{255, 255, 255, 255},
				ForegroundColor: color.RGBA{0, 0, 0, 255},
				ModuleDrawer:    NewGappedCircleModuleDrawer(0.9),
			},
			desc: "Gapped circles (90% size)",
		},
		{
			name: "06_gapped_circles_70",
			config: StyleConfig{
				ModuleSize:      10,
				QuietZone:       4,
				BackgroundColor: color.RGBA{255, 255, 255, 255},
				ForegroundColor: color.RGBA{0, 0, 0, 255},
				ModuleDrawer:    NewGappedCircleModuleDrawer(0.7),
			},
			desc: "Gapped circles (70% size)",
		},
		{
			name: "07_rounded_medium",
			config: StyleConfig{
				ModuleSize:      10,
				QuietZone:       4,
				BackgroundColor: color.RGBA{255, 255, 255, 255},
				ForegroundColor: color.RGBA{0, 0, 0, 255},
				ModuleDrawer:    NewRoundedModuleDrawer(0.5),
			},
			desc: "Rounded corners (50% radius)",
		},
		{
			name:   "08_chrome_preset",
			config: ChromeStyleConfig(),
			desc:   "Chrome preset configuration",
		},
		{
			name:   "09_chrome_finder_preset",
			config: ChromeFinderPatternStyleConfig(),
			desc:   "Chrome with rounded finder patterns",
		},
		{
			name: "10_colored_blue",
			config: StyleConfig{
				ModuleSize:      10,
				QuietZone:       4,
				BackgroundColor: color.RGBA{240, 248, 255, 255}, // Alice blue
				ForegroundColor: color.RGBA{0, 100, 200, 255},   // Blue
				ModuleDrawer:    NewCircleModuleDrawer(),
			},
			desc: "Colored circles (blue theme)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create QR code
			qr, err := New(testData, High)
			if err != nil {
				t.Fatalf("Failed to create QR code: %v", err)
			}

			err = qr.Encode()
			if err != nil {
				t.Fatalf("Failed to encode QR code: %v", err)
			}

			// Generate image
			img, err := qr.ToImage(test.config)
			if err != nil {
				t.Fatalf("Failed to generate image: %v", err)
			}

			// Save image
			filename := filepath.Join(outputDir, test.name+".png")
			file, err := os.Create(filename)
			if err != nil {
				t.Fatalf("Failed to create file: %v", err)
			}
			defer file.Close()

			err = png.Encode(file, img)
			if err != nil {
				t.Fatalf("Failed to encode PNG: %v", err)
			}

			t.Logf("Generated %s: %s", test.name, test.desc)
		})
	}

	t.Logf("All style comparison images generated in %s/", outputDir)
}

func TestMakeAPIComparison(t *testing.T) {
	outputDir := "make_api_tests"

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	testData := "https://example.com/test"

	// Test the Make API with different options
	tests := []struct {
		name    string
		options []func(*StyleConfig)
		desc    string
	}{
		{
			name:    "make_default",
			options: nil,
			desc:    "Default Make() configuration",
		},
		{
			name:    "make_circles",
			options: []func(*StyleConfig){WithCircles()},
			desc:    "Make() with circles",
		},
		{
			name:    "make_gapped_circles",
			options: []func(*StyleConfig){WithGappedCircles(0.8)},
			desc:    "Make() with gapped circles",
		},
		{
			name:    "make_rounded",
			options: []func(*StyleConfig){WithRoundedCorners(1.0)},
			desc:    "Make() with rounded corners",
		},
		{
			name:    "make_gapped_squares",
			options: []func(*StyleConfig){WithGappedSquares(0.7)},
			desc:    "Make() with gapped squares",
		},
		{
			name: "make_styled_blue",
			options: []func(*StyleConfig){
				WithCircles(),
				WithColors(color.RGBA{0, 50, 150, 255}, color.RGBA{240, 248, 255, 255}),
				WithModuleSize(12),
			},
			desc: "Make() with multiple style options",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Use the Make API
			img, err := Make(testData, test.options...)
			if err != nil {
				t.Fatalf("Failed to make QR code: %v", err)
			}

			// Save image
			filename := filepath.Join(outputDir, test.name+".png")
			file, err := os.Create(filename)
			if err != nil {
				t.Fatalf("Failed to create file: %v", err)
			}
			defer file.Close()

			err = png.Encode(file, img)
			if err != nil {
				t.Fatalf("Failed to encode PNG: %v", err)
			}

			t.Logf("Generated %s: %s", test.name, test.desc)
		})
	}

	t.Logf("All Make API test images generated in %s/", outputDir)
}

func TestChromeReplicationComparison(t *testing.T) {
	outputDir := "chrome_replication_tests"

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// Use the same data as the reference Chrome QR code
	testData := "https://meet.google.com/"

	// Test different approaches to replicate Chrome's style
	tests := []struct {
		name   string
		config StyleConfig
		desc   string
	}{
		{
			name:   "chrome_attempt_1_circles",
			config: ChromeStyleConfig(),
			desc:   "Chrome preset with circles",
		},
		{
			name:   "chrome_attempt_2_rounded",
			config: ChromeFinderPatternStyleConfig(),
			desc:   "Chrome preset with rounded finder patterns",
		},
		{
			name: "chrome_attempt_3_custom",
			config: StyleConfig{
				ModuleSize:      12,
				QuietZone:       6,
				BackgroundColor: color.RGBA{255, 255, 255, 255},
				ForegroundColor: color.RGBA{0, 0, 0, 255},
				ModuleDrawer:    NewRoundedModuleDrawer(0.8),
			},
			desc: "Custom Chrome-like configuration",
		},
		{
			name: "chrome_attempt_4_mixed",
			config: StyleConfig{
				ModuleSize:      10,
				QuietZone:       4,
				BackgroundColor: color.RGBA{255, 255, 255, 255},
				ForegroundColor: color.RGBA{0, 0, 0, 255},
				ModuleDrawer:    NewGappedCircleModuleDrawer(0.85),
			},
			desc: "Gapped circles matching Chrome proportions",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create QR code
			qr, err := New(testData, High)
			if err != nil {
				t.Fatalf("Failed to create QR code: %v", err)
			}

			err = qr.Encode()
			if err != nil {
				t.Fatalf("Failed to encode QR code: %v", err)
			}

			// Generate image
			img, err := qr.ToImage(test.config)
			if err != nil {
				t.Fatalf("Failed to generate image: %v", err)
			}

			// Save image
			filename := filepath.Join(outputDir, test.name+".png")
			file, err := os.Create(filename)
			if err != nil {
				t.Fatalf("Failed to create file: %v", err)
			}
			defer file.Close()

			err = png.Encode(file, img)
			if err != nil {
				t.Fatalf("Failed to encode PNG: %v", err)
			}

			t.Logf("Generated %s: %s", test.name, test.desc)
		})
	}

	// Copy reference image for comparison
	srcFile := "qrcode_meet.google.com.png"
	if _, err := os.Stat(srcFile); err == nil {
		dstFile := filepath.Join(outputDir, "00_reference_chrome.png")
		if err := copyFile(srcFile, dstFile); err != nil {
			t.Logf("Warning: Could not copy reference image: %v", err)
		} else {
			t.Logf("Copied reference Chrome QR code for comparison")
		}
	}

	t.Logf("Chrome replication test images generated in %s/", outputDir)
	t.Logf("Compare generated images with 00_reference_chrome.png")
}

// Helper function to copy files
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Decode and re-encode to ensure it's a valid PNG
	img, err := png.Decode(sourceFile)
	if err != nil {
		return err
	}

	return png.Encode(destFile, img)
}

func TestSizeAndScaleComparison(t *testing.T) {
	outputDir := "size_scale_tests"

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	testData := "Test"

	// Test different module sizes
	sizes := []int{5, 8, 10, 15, 20}

	for _, size := range sizes {
		t.Run(fmt.Sprintf("size_%d", size), func(t *testing.T) {
			config := StyleConfig{
				ModuleSize:      size,
				QuietZone:       4,
				BackgroundColor: color.RGBA{255, 255, 255, 255},
				ForegroundColor: color.RGBA{0, 0, 0, 255},
				ModuleDrawer:    NewCircleModuleDrawer(),
			}

			qr, err := New(testData, Medium)
			if err != nil {
				t.Fatalf("Failed to create QR code: %v", err)
			}

			err = qr.Encode()
			if err != nil {
				t.Fatalf("Failed to encode QR code: %v", err)
			}

			img, err := qr.ToImage(config)
			if err != nil {
				t.Fatalf("Failed to generate image: %v", err)
			}

			filename := filepath.Join(outputDir, fmt.Sprintf("module_size_%d.png", size))
			file, err := os.Create(filename)
			if err != nil {
				t.Fatalf("Failed to create file: %v", err)
			}
			defer file.Close()

			err = png.Encode(file, img)
			if err != nil {
				t.Fatalf("Failed to encode PNG: %v", err)
			}

			t.Logf("Generated module size %d test", size)
		})
	}

	t.Logf("Size and scale test images generated in %s/", outputDir)
}
