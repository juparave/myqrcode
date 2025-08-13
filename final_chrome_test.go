package myqrcode

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestFinalChromeReplication(t *testing.T) {
	outputDir := "final_chrome_tests"
	
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}
	
	// Use exact same data as Chrome reference
	testData := "https://meet.google.com/"
	
	// Test the best implementations against Chrome reference
	tests := []struct {
		name   string
		config StyleConfig
		desc   string
	}{
		{
			name:   "chrome_original_style",
			config: ChromeStyleConfig(),
			desc:   "Original Chrome preset (circles)",
		},
		{
			name:   "chrome_finder_style", 
			config: ChromeFinderPatternStyleConfig(),
			desc:   "Chrome with rounded finder patterns",
		},
		{
			name:   "chrome_gapped_87",
			config: ChromeGappedStyleConfig(),
			desc:   "New Chrome gapped style (87% ratio)",
		},
		{
			name:   "chrome_gapped_85",
			config: ChromeGappedStyleConfigWithRatio(0.85),
			desc:   "Chrome gapped style (85% ratio)",
		},
		{
			name:   "chrome_gapped_90",
			config: ChromeGappedStyleConfigWithRatio(0.90),
			desc:   "Chrome gapped style (90% ratio)",
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
	
	// Copy reference for direct comparison
	srcFile := "qrcode_meet.google.com.png"
	if _, err := os.Stat(srcFile); err == nil {
		dstFile := filepath.Join(outputDir, "00_chrome_reference.png")
		if err := copyFile(srcFile, dstFile); err != nil {
			t.Logf("Warning: Could not copy reference image: %v", err)
		} else {
			t.Logf("Copied Chrome reference for direct comparison")
		}
	}
	
	t.Logf("Final Chrome replication tests generated in %s/", outputDir)
	t.Logf("Compare all variants with 00_chrome_reference.png")
}

func TestUpdatedExampleWithLogo(t *testing.T) {
	outputDir := "updated_logo_tests"
	
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}
	
	// Test the updated example with the new gapped circle implementation
	testData := "https://meet.google.com/abc-defg-hij"
	
	// Load logo if available
	var logo image.Image
	logoFile, err := os.Open("res/goqrlogp.jpg")
	if err == nil {
		defer logoFile.Close()
		img, err := jpeg.Decode(logoFile)
		if err == nil {
			logo = img
		}
	}
	
	tests := []struct {
		name     string
		logoSize int
		config   StyleConfig
		desc     string
	}{
		{
			name:     "updated_chrome_gapped_logo",
			logoSize: 20,
			config:   ChromeGappedStyleConfig(),
			desc:     "New Chrome gapped style with logo",
		},
		{
			name:     "updated_chrome_original_logo",
			logoSize: 20,
			config:   ChromeStyleConfig(),
			desc:     "Original Chrome style with logo",
		},
	}
	
	for _, test := range tests {
		if logo == nil && test.logoSize > 0 {
			t.Logf("Skipping %s: logo not available", test.name)
			continue
		}
		
		t.Run(test.name, func(t *testing.T) {
			// Create QR code
			qr, err := New(testData, High)
			if err != nil {
				t.Fatalf("Failed to create QR code: %v", err)
			}
			
			// Set logo if available
			if logo != nil && test.logoSize > 0 {
				qr.SetLogo(logo, test.logoSize)
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
	
	t.Logf("Updated logo tests generated in %s/", outputDir)
}