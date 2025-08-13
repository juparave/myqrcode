package myqrcode

import (
	"errors"
	"image"
	"image/color"
	"image/draw"

	xdraw "golang.org/x/image/draw"
)

func (qr *QRCode) ToImage(config StyleConfig) (image.Image, error) {
	if qr.Matrix == nil {
		return nil, errors.New("QR code not encoded")
	}

	// Validate and set defaults
	moduleSize := config.ModuleSize
	if moduleSize <= 0 {
		moduleSize = 8
	}

	quietZone := config.QuietZone
	if quietZone <= 0 {
		quietZone = 4 * moduleSize
	}

	// Set default colors if not provided
	if config.BackgroundColor == nil {
		config.BackgroundColor = color.RGBA{255, 255, 255, 255}
	}
	if config.ForegroundColor == nil {
		config.ForegroundColor = color.RGBA{0, 0, 0, 255}
	}

	// Set default module drawer if not provided
	drawer := config.ModuleDrawer
	if drawer == nil {
		if config.CircularDots {
			drawer = NewCircleModuleDrawer()
		} else if config.RoundedCorners {
			drawer = NewRoundedModuleDrawer(1.0)
		} else {
			drawer = NewSquareModuleDrawer()
		}
	}

	// Create image
	imgSize := qr.Size*moduleSize + 2*quietZone
	img := image.NewRGBA(image.Rect(0, 0, imgSize, imgSize))

	// Fill background
	draw.Draw(img, img.Bounds(), &image.Uniform{config.BackgroundColor}, image.Point{}, draw.Src)

	// Initialize the module drawer
	drawer.Initialize(img, config)

	// Draw QR modules using the module drawer
	for y := 0; y < qr.Size; y++ {
		for x := 0; x < qr.Size; x++ {
			imgX := quietZone + x*moduleSize
			imgY := quietZone + y*moduleSize

			// Create box coordinates [x1, y1, x2, y2]
			box := [4]int{imgX, imgY, imgX + moduleSize, imgY + moduleSize}

			// Get neighbor context if the drawer needs it
			var neighbors *ActiveWithNeighbors
			if drawer.NeedsNeighbors() {
				if IsFinderPattern(y, x, qr.Size) {
					neighbors = GetFinderPatternNeighbors(qr.Matrix, y, x)
				} else {
					neighbors = GetModuleNeighbors(qr.Matrix, y, x)
				}
			}

			// Draw the module
			drawer.DrawModule(box, qr.Matrix[y][x], neighbors)
		}
	}

	// Draw logo if present
	if qr.Logo != nil && qr.LogoSize > 0 {
		placement := calculateLogoPlacement(&Matrix{Size: qr.Size}, qr.LogoSize)
		drawLogo(img, qr.Logo, placement, moduleSize, quietZone)
	}

	return img, nil
}

func isFinderPattern(x, y, size int) bool {
	// Top-left finder pattern
	if x < 7 && y < 7 {
		return true
	}
	// Top-right finder pattern
	if x >= size-7 && y < 7 {
		return true
	}
	// Bottom-left finder pattern
	if x < 7 && y >= size-7 {
		return true
	}
	return false
}

func drawCircle(img *image.RGBA, x, y, size int, color color.Color) {
	//-FIX--: Draw a simple square instead of a circle for better readability
	rect := image.Rect(x, y, x+size, y+size)
	draw.Draw(img, rect, &image.Uniform{color}, image.Point{}, draw.Src)
}

func drawRoundedModule(img *image.RGBA, x, y, size int, color color.Color) {
	cornerRadius := float64(size) / 4

	for py := y; py < y+size; py++ {
		for px := x; px < x+size; px++ {
			relX := float64(px - x)
			relY := float64(py - y)

			// Check if pixel is in rounded corner area
			inCorner := false

			// Top-left corner
			if relX < cornerRadius && relY < cornerRadius {
				dx := relX - cornerRadius
				dy := relY - cornerRadius
				if dx*dx+dy*dy > cornerRadius*cornerRadius {
					inCorner = true
				}
			}
			// Top-right corner
			if relX >= float64(size)-cornerRadius && relY < cornerRadius {
				dx := relX - (float64(size) - cornerRadius)
				dy := relY - cornerRadius
				if dx*dx+dy*dy > cornerRadius*cornerRadius {
					inCorner = true
				}
			}
			// Bottom-left corner
			if relX < cornerRadius && relY >= float64(size)-cornerRadius {
				dx := relX - cornerRadius
				dy := relY - (float64(size) - cornerRadius)
				if dx*dx+dy*dy > cornerRadius*cornerRadius {
					inCorner = true
				}
			}
			// Bottom-right corner
			if relX >= float64(size)-cornerRadius && relY >= float64(size)-cornerRadius {
				dx := relX - (float64(size) - cornerRadius)
				dy := relY - (float64(size) - cornerRadius)
				if dx*dx+dy*dy > cornerRadius*cornerRadius {
					inCorner = true
				}
			}

			if !inCorner {
				img.Set(px, py, color)
			}
		}
	}
}

func drawRoundedFinderPattern(img *image.RGBA, x, y, moduleSize int, color color.Color, qrX, qrY, qrSize int) {
	// Determine position within finder pattern
	var localX, localY int
	if qrX < 7 && qrY < 7 {
		// Top-left
		localX = qrX
		localY = qrY
	} else if qrX >= qrSize-7 && qrY < 7 {
		// Top-right
		localX = qrX - (qrSize - 7)
		localY = qrY
	} else if qrX < 7 && qrY >= qrSize-7 {
		// Bottom-left
		localX = qrX
		localY = qrY - (qrSize - 7)
	}

	// Draw with enhanced rounding for finder patterns
	if (localX == 0 || localX == 6) || (localY == 0 || localY == 6) {
		drawRoundedModule(img, x, y, moduleSize, color)
	} else {
		// Regular module for inner parts
		rect := image.Rect(x, y, x+moduleSize, y+moduleSize)
		draw.Draw(img, rect, &image.Uniform{color}, image.Point{}, draw.Src)
	}
}

func drawLogo(img *image.RGBA, logo image.Image, placement LogoPlacement, moduleSize, quietZone int) {
	logoWidth := placement.Width * moduleSize
	logoHeight := placement.Height * moduleSize

	logoX := quietZone + placement.X*moduleSize
	logoY := quietZone + placement.Y*moduleSize

	// Scale logo to fit
	logoRect := image.Rect(logoX, logoY, logoX+logoWidth, logoY+logoHeight)

	// Use bilinear scaling for better quality
	xdraw.BiLinear.Scale(img, logoRect, logo, logo.Bounds(), xdraw.Src, nil)
}
