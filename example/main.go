package main

import (
	"image/png"
	"log"
	"os"

	"github.com/juparave/myqrcode"
)

func main() {
	// Create a QR code with Chrome-style appearance
	qr, err := myqrcode.New("https://meet.google.com/abc-defg-hij", myqrcode.High)
	if err != nil {
		log.Fatal(err)
	}

	// Encode the QR code
	err = qr.Encode()
	if err != nil {
		log.Fatal(err)
	}

	// Configure Chrome-style rendering
	config := myqrcode.StyleConfig{
		ModuleSize:      8,
		QuietZone:       32,
		RoundedCorners:  true,
		CircularDots:    false,
	}

	// Generate the image
	img, err := qr.ToImage(config)
	if err != nil {
		log.Fatal(err)
	}

	// Save to file
	file, err := os.Create("chrome_style_qr.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Chrome-style QR code generated: chrome_style_qr.png")
}