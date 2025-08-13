package myqrcode

import (
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestChromeGappedCircleRatios(t *testing.T) {
	outputDir := "chrome_gap_tests"
	
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}
	
	// Use the same data as Chrome reference
	testData := "https://meet.google.com/"
	
	// Test different size ratios to find the Chrome match
	ratios := []float64{0.80, 0.85, 0.87, 0.90, 0.92, 0.95}
	
	for _, ratio := range ratios {
		t.Run(formatRatioName(ratio), func(t *testing.T) {
			// Create QR code
			qr, err := New(testData, High)
			if err != nil {
				t.Fatalf("Failed to create QR code: %v", err)
			}
			
			err = qr.Encode()
			if err != nil {
				t.Fatalf("Failed to encode QR code: %v", err)
			}
			
			// Test the improved gapped circle drawer
			config := StyleConfig{
				ModuleSize:       10,
				QuietZone:        4,
				BackgroundColor:  color.RGBA{255, 255, 255, 255},
				ForegroundColor:  color.RGBA{0, 0, 0, 255},
				ModuleDrawer:     NewGappedCircleModuleDrawer(ratio),
			}
			
			// Generate image
			img, err := qr.ToImage(config)
			if err != nil {
				t.Fatalf("Failed to generate image: %v", err)
			}
			
			// Save image
			filename := filepath.Join(outputDir, formatRatioName(ratio)+".png")
			file, err := os.Create(filename)
			if err != nil {
				t.Fatalf("Failed to create file: %v", err)
			}
			defer file.Close()
			
			err = png.Encode(file, img)
			if err != nil {
				t.Fatalf("Failed to encode PNG: %v", err)
			}
			
			t.Logf("Generated ratio %.2f test", ratio)
		})
	}
	
	// Copy reference image for comparison
	srcFile := "qrcode_meet.google.com.png"
	if _, err := os.Stat(srcFile); err == nil {
		dstFile := filepath.Join(outputDir, "00_chrome_reference.png")
		if err := copyFile(srcFile, dstFile); err != nil {
			t.Logf("Warning: Could not copy reference image: %v", err)
		} else {
			t.Logf("Copied Chrome reference for comparison")
		}
	}
	
	t.Logf("Chrome gapped circle ratio tests generated in %s/", outputDir)
	t.Logf("Compare with 00_chrome_reference.png to find best match")
}

func TestChromeGappedComparison(t *testing.T) {
	outputDir := "chrome_comparison_tests" 
	
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}
	
	testData := "https://meet.google.com/"
	
	// Test old vs new gapped circle implementation
	tests := []struct {
		name   string
		config StyleConfig
		desc   string
	}{
		{
			name: "old_approach_90",
			config: StyleConfig{
				ModuleSize:       10,
				QuietZone:        4,
				BackgroundColor:  color.RGBA{255, 255, 255, 255},
				ForegroundColor:  color.RGBA{0, 0, 0, 255},
				CircularDots:     true, // Use old approach
			},
			desc: "Old circular dots approach",
		},
		{
			name: "new_gapped_85",
			config: StyleConfig{
				ModuleSize:       10,
				QuietZone:        4,
				BackgroundColor:  color.RGBA{255, 255, 255, 255},
				ForegroundColor:  color.RGBA{0, 0, 0, 255},
				ModuleDrawer:     NewGappedCircleModuleDrawer(0.85),
			},
			desc: "New gapped circles (85% ratio)",
		},
		{
			name: "new_gapped_87",
			config: StyleConfig{
				ModuleSize:       10,
				QuietZone:        4,
				BackgroundColor:  color.RGBA{255, 255, 255, 255},
				ForegroundColor:  color.RGBA{0, 0, 0, 255},
				ModuleDrawer:     NewGappedCircleModuleDrawer(0.87),
			},
			desc: "New gapped circles (87% ratio)",
		},
		{
			name: "new_gapped_90",
			config: StyleConfig{
				ModuleSize:       10,
				QuietZone:        4,
				BackgroundColor:  color.RGBA{255, 255, 255, 255},
				ForegroundColor:  color.RGBA{0, 0, 0, 255},
				ModuleDrawer:     NewGappedCircleModuleDrawer(0.90),
			},
			desc: "New gapped circles (90% ratio)",
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
	
	t.Logf("Chrome comparison tests generated in %s/", outputDir)
}

// Helper functions

func formatRatioName(ratio float64) string {
	return "ratio_" + intToString(int(ratio*100))
}

func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	
	var result string
	for n > 0 {
		result = string(rune('0'+(n%10))) + result
		n /= 10
	}
	return result
}