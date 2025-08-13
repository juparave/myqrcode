package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/juparave/myqrcode"
)

func main() {
	testData := "https://meet.google.com/abc-defg-hij"

	log.Println("Generating QR codes with different styles...")

	// Example 1: Simple Make API with default style
	log.Println("1. Creating basic QR code...")
	img1, err := myqrcode.Make(testData)
	if err != nil {
		log.Fatal(err)
	}
	saveImage(img1, "01_basic_qr.png")

	// Example 2: Chrome-style circles
	log.Println("2. Creating Chrome-style with circles...")
	img2, err := myqrcode.Make(testData, myqrcode.WithCircles())
	if err != nil {
		log.Fatal(err)
	}
	saveImage(img2, "02_chrome_circles.png")

	// Example 3: Rounded corners (Chrome finder patterns)
	log.Println("3. Creating with rounded corners...")
	img3, err := myqrcode.Make(testData, myqrcode.WithRoundedCorners(1.0))
	if err != nil {
		log.Fatal(err)
	}
	saveImage(img3, "03_rounded_corners.png")

	// Example 4: Gapped circles
	log.Println("4. Creating with gapped circles...")
	img4, err := myqrcode.Make(testData, myqrcode.WithGappedCircles(0.8))
	if err != nil {
		log.Fatal(err)
	}
	saveImage(img4, "04_gapped_circles.png")

	// Example 5: Styled with colors
	log.Println("5. Creating styled with colors...")
	img5, err := myqrcode.Make(testData,
		myqrcode.WithCircles(),
		myqrcode.WithColors(
			color.RGBA{0, 100, 200, 255},   // Blue foreground
			color.RGBA{240, 248, 255, 255}, // Light blue background
		),
		myqrcode.WithModuleSize(12),
	)
	if err != nil {
		log.Fatal(err)
	}
	saveImage(img5, "05_styled_blue.png")

	// Example 6: Chrome preset
	log.Println("6. Creating with Chrome preset...")
	qr6, err := myqrcode.New(testData, myqrcode.High)
	if err != nil {
		log.Fatal(err)
	}
	err = qr6.Encode()
	if err != nil {
		log.Fatal(err)
	}
	img6, err := qr6.ToImage(myqrcode.ChromeStyleConfig())
	if err != nil {
		log.Fatal(err)
	}
	saveImage(img6, "06_chrome_preset.png")

	// Example 7: Chrome finder pattern preset
	log.Println("7. Creating with Chrome finder pattern preset...")
	qr7, err := myqrcode.New(testData, myqrcode.High)
	if err != nil {
		log.Fatal(err)
	}
	err = qr7.Encode()
	if err != nil {
		log.Fatal(err)
	}
	img7, err := qr7.ToImage(myqrcode.ChromeFinderPatternStyleConfig())
	if err != nil {
		log.Fatal(err)
	}
	saveImage(img7, "07_chrome_finder_preset.png")

	// Example 8: Custom module drawer
	log.Println("8. Creating with custom module drawer...")
	qr8, err := myqrcode.New(testData, myqrcode.High)
	if err != nil {
		log.Fatal(err)
	}
	err = qr8.Encode()
	if err != nil {
		log.Fatal(err)
	}

	config := myqrcode.StyleConfig{
		ModuleSize:      15,
		QuietZone:       6,
		BackgroundColor: color.RGBA{255, 255, 255, 255},
		ForegroundColor: color.RGBA{0, 0, 0, 255},
		ModuleDrawer:    myqrcode.NewGappedSquareModuleDrawer(0.7),
	}

	img8, err := qr8.ToImage(config)
	if err != nil {
		log.Fatal(err)
	}
	saveImage(img8, "08_custom_gapped_squares.png")

	log.Println("All QR code examples generated successfully!")
	log.Println("Generated files:")
	log.Println("  01_basic_qr.png - Basic square modules")
	log.Println("  02_chrome_circles.png - Chrome-style circular dots")
	log.Println("  03_rounded_corners.png - Context-aware rounded corners")
	log.Println("  04_gapped_circles.png - Circular dots with spacing")
	log.Println("  05_styled_blue.png - Custom colors and styling")
	log.Println("  06_chrome_preset.png - Chrome preset configuration")
	log.Println("  07_chrome_finder_preset.png - Chrome with rounded finder patterns")
	log.Println("  08_custom_gapped_squares.png - Custom gapped squares")
}

func saveImage(img image.Image, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Error creating %s: %v", filename, err)
		return
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		log.Printf("Error encoding %s: %v", filename, err)
		return
	}

	log.Printf("  Generated: %s", filename)
}
