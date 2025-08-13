package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"

	"github.com/juparave/myqrcode"
)

func main() {
	// Load the logo image from res/goqrlogp.jpg
	logo, err := loadLogo()
	if err != nil {
		log.Printf("Could not load logo from res/goqrlogp.jpg, using simple placeholder: %v", err)
		logo = createSimpleLogo()
	}

	// Create a QR code with Chrome-style appearance and logo
	qr, err := myqrcode.New("https://meet.google.com/abc-defg-hij", myqrcode.High)
	if err != nil {
		log.Fatal(err)
	}

	// Set the logo (20% of QR code size)
	qr.SetLogo(logo, 20)

	// Encode the QR code
	err = qr.Encode()
	if err != nil {
		log.Fatal(err)
	}

	// Configure Chrome-style rendering with circular dots using new module drawer system
	config := myqrcode.StyleConfig{
		ModuleSize:       10,
		QuietZone:        4,
		RoundedCorners:   true,
		CircularDots:     true,
		BackgroundColor:  myqrcode.DefaultStyleConfig().BackgroundColor,
		ForegroundColor:  myqrcode.DefaultStyleConfig().ForegroundColor,
		ModuleDrawer:     myqrcode.NewCircleModuleDrawer(), // Use the new circular module drawer
	}

	// Generate the image
	img, err := qr.ToImage(config)
	if err != nil {
		log.Fatal(err)
	}

	// Save to file
	file, err := os.Create("chrome_style_qr_with_logo.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Chrome-style QR code with logo generated: chrome_style_qr_with_logo.png")
}

func createSimpleLogo() image.Image {
	// Create a simple 32x32 logo (placeholder)
	size := 32
	img := image.NewRGBA(image.Rect(0, 0, size, size))

	// Fill with black to simulate a simple logo
	for y := 8; y < size-8; y++ {
		for x := 8; x < size-8; x++ {
			img.Set(x, y, image.Black)
		}
	}

	return img
}

func loadLogo() (image.Image, error) {
	file, err := os.Open("res/goqrlogp.jpg")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}
	
	return img, nil
}
