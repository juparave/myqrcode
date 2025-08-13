package myqrcode

import (
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestQRCodeWithLogos(t *testing.T) {
	outputDir := "logo_tests"
	
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}
	
	testData := "https://meet.google.com/abc-defg-hij"
	
	// Load the logo image
	logoFile, err := os.Open("res/goqrlogp.jpg")
	if err != nil {
		t.Skipf("Logo file not found (res/goqrlogp.jpg): %v", err)
		return
	}
	defer logoFile.Close()
	
	logo, err := jpeg.Decode(logoFile)
	if err != nil {
		t.Fatalf("Failed to decode logo image: %v", err)
	}
	
	// Test different logo sizes and styles
	tests := []struct {
		name     string
		logoSize int
		config   StyleConfig
		desc     string
	}{
		{
			name:     "logo_squares_15",
			logoSize: 15,
			config: StyleConfig{
				ModuleSize:       10,
				QuietZone:        4,
				BackgroundColor:  DefaultStyleConfig().BackgroundColor,
				ForegroundColor:  DefaultStyleConfig().ForegroundColor,
				ModuleDrawer:     NewSquareModuleDrawer(),
			},
			desc: "Basic squares with 15% logo",
		},
		{
			name:     "logo_circles_20",
			logoSize: 20,
			config: StyleConfig{
				ModuleSize:       10,
				QuietZone:        4,
				BackgroundColor:  DefaultStyleConfig().BackgroundColor,
				ForegroundColor:  DefaultStyleConfig().ForegroundColor,
				ModuleDrawer:     NewCircleModuleDrawer(),
			},
			desc: "Chrome-style circles with 20% logo",
		},
		{
			name:     "logo_rounded_25",
			logoSize: 25,
			config: StyleConfig{
				ModuleSize:       10,
				QuietZone:        4,
				BackgroundColor:  DefaultStyleConfig().BackgroundColor,
				ForegroundColor:  DefaultStyleConfig().ForegroundColor,
				ModuleDrawer:     NewRoundedModuleDrawer(1.0),
			},
			desc: "Rounded corners with 25% logo",
		},
		{
			name:     "logo_gapped_circles_30",
			logoSize: 30,
			config: StyleConfig{
				ModuleSize:       12,
				QuietZone:        6,
				BackgroundColor:  DefaultStyleConfig().BackgroundColor,
				ForegroundColor:  DefaultStyleConfig().ForegroundColor,
				ModuleDrawer:     NewGappedCircleModuleDrawer(0.8),
			},
			desc: "Gapped circles with 30% logo",
		},
		{
			name:     "logo_chrome_style_20",
			logoSize: 20,
			config:   ChromeStyleConfig(),
			desc:     "Chrome preset with 20% logo",
		},
		{
			name:     "logo_chrome_finder_25",
			logoSize: 25,
			config:   ChromeFinderPatternStyleConfig(),
			desc:     "Chrome finder pattern style with 25% logo",
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create QR code with logo
			qr, err := New(testData, High) // Use High error correction for logos
			if err != nil {
				t.Fatalf("Failed to create QR code: %v", err)
			}
			
			// Set the logo
			qr.SetLogo(logo, test.logoSize)
			
			// Encode
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
	
	t.Logf("All logo test images generated in %s/", outputDir)
}

func TestMakeAPIWithLogo(t *testing.T) {
	outputDir := "make_logo_tests"
	
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}
	
	// Load the logo image
	logoFile, err := os.Open("res/goqrlogp.jpg")
	if err != nil {
		t.Skipf("Logo file not found (res/goqrlogp.jpg): %v", err)
		return
	}
	defer logoFile.Close()
	
	logo, err := jpeg.Decode(logoFile)
	if err != nil {
		t.Fatalf("Failed to decode logo image: %v", err)
	}
	
	testData := "https://github.com/myqrcode"
	
	// Test the convenience function for creating QR codes with logos
	tests := []struct {
		name     string
		logoSize int
		options  []func(*StyleConfig)
		desc     string
	}{
		{
			name:     "make_logo_basic",
			logoSize: 20,
			options:  nil,
			desc:     "Basic QR with logo using Make API",
		},
		{
			name:     "make_logo_circles",
			logoSize: 25,
			options:  []func(*StyleConfig){WithCircles()},
			desc:     "Circular dots with logo",
		},
		{
			name:     "make_logo_rounded",
			logoSize: 20,
			options:  []func(*StyleConfig){WithRoundedCorners(1.0)},
			desc:     "Rounded corners with logo",
		},
		{
			name:     "make_logo_gapped",
			logoSize: 25,
			options:  []func(*StyleConfig){WithGappedCircles(0.75)},
			desc:     "Gapped circles with logo",
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create QR code
			qr, err := New(testData, High)
			if err != nil {
				t.Fatalf("Failed to create QR code: %v", err)
			}
			
			// Set logo
			qr.SetLogo(logo, test.logoSize)
			
			// Encode
			err = qr.Encode()
			if err != nil {
				t.Fatalf("Failed to encode QR code: %v", err)
			}
			
			// Create config with options
			config := DefaultStyleConfig()
			for _, opt := range test.options {
				opt(&config)
			}
			
			// Generate image
			img, err := qr.ToImage(config)
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
	
	t.Logf("All Make API logo test images generated in %s/", outputDir)
}

func TestErrorCorrectionWithLogos(t *testing.T) {
	outputDir := "error_correction_tests"
	
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}
	
	// Load the logo image
	logoFile, err := os.Open("res/goqrlogp.jpg")
	if err != nil {
		t.Skipf("Logo file not found (res/goqrlogp.jpg): %v", err)
		return
	}
	defer logoFile.Close()
	
	logo, err := jpeg.Decode(logoFile)
	if err != nil {
		t.Fatalf("Failed to decode logo image: %v", err)
	}
	
	testData := "Test QR code error correction"
	
	// Test different error correction levels with logos
	tests := []struct {
		name       string
		level      ErrorCorrectionLevel
		logoSize   int
		desc       string
	}{
		{
			name:     "ec_low_15",
			level:    Low,
			logoSize: 15,
			desc:     "Low error correction with 15% logo",
		},
		{
			name:     "ec_medium_20",
			level:    Medium,
			logoSize: 20,
			desc:     "Medium error correction with 20% logo",
		},
		{
			name:     "ec_quartile_25",
			level:    Quartile,
			logoSize: 25,
			desc:     "Quartile error correction with 25% logo",
		},
		{
			name:     "ec_high_30",
			level:    High,
			logoSize: 30,
			desc:     "High error correction with 30% logo",
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create QR code
			qr, err := New(testData, test.level)
			if err != nil {
				t.Fatalf("Failed to create QR code: %v", err)
			}
			
			// Set logo
			qr.SetLogo(logo, test.logoSize)
			
			// Encode (this should auto-adjust error correction if needed)
			err = qr.Encode()
			if err != nil {
				t.Fatalf("Failed to encode QR code: %v", err)
			}
			
			// Use Chrome style for consistent appearance
			config := ChromeStyleConfig()
			
			// Generate image
			img, err := qr.ToImage(config)
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
			
			t.Logf("Generated %s: %s (final level: %v)", test.name, test.desc, qr.ErrorCorrection)
		})
	}
	
	t.Logf("Error correction test images generated in %s/", outputDir)
}