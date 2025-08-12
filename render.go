package myqrcode

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"math"

	xdraw "golang.org/x/image/draw"
)

func (qr *QRCode) ToImage(config StyleConfig) (image.Image, error) {
	if qr.Matrix == nil {
		return nil, errors.New("QR code not encoded")
	}
	
	moduleSize := config.ModuleSize
	if moduleSize <= 0 {
		moduleSize = 8
	}
	
	quietZone := config.QuietZone
	if quietZone <= 0 {
		quietZone = 4 * moduleSize
	}
	
	imgSize := qr.Size*moduleSize + 2*quietZone
	img := image.NewRGBA(image.Rect(0, 0, imgSize, imgSize))
	
	// Fill background
	bg := config.BackgroundColor
	if bg == nil {
		bg = color.RGBA{255, 255, 255, 255}
	}
	draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)
	
	// Draw QR modules
	fg := config.ForegroundColor
	if fg == nil {
		fg = color.RGBA{0, 0, 0, 255}
	}
	
	for y := 0; y < qr.Size; y++ {
		for x := 0; x < qr.Size; x++ {
			if qr.Matrix[y][x] {
				imgX := quietZone + x*moduleSize
				imgY := quietZone + y*moduleSize
				
				if config.CircularDots {
					drawCircle(img, imgX, imgY, moduleSize, fg)
				} else if config.RoundedCorners && isFinderPattern(x, y, qr.Size) {
					drawRoundedFinderPattern(img, imgX, imgY, moduleSize, fg, x, y, qr.Size)
				} else if config.RoundedCorners {
					drawRoundedModule(img, imgX, imgY, moduleSize, fg)
				} else {
					rect := image.Rect(imgX, imgY, imgX+moduleSize, imgY+moduleSize)
					draw.Draw(img, rect, &image.Uniform{fg}, image.Point{}, draw.Src)
				}
			}
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
	radius := float64(size) / 2
	centerX := float64(x) + radius
	centerY := float64(y) + radius
	
	for py := y; py < y+size; py++ {
		for px := x; px < x+size; px++ {
			dx := float64(px) + 0.5 - centerX
			dy := float64(py) + 0.5 - centerY
			distance := math.Sqrt(dx*dx + dy*dy)
			
			if distance <= radius {
				img.Set(px, py, color)
			}
		}
	}
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
				if dx*dx + dy*dy > cornerRadius*cornerRadius {
					inCorner = true
				}
			}
			// Top-right corner
			if relX >= float64(size)-cornerRadius && relY < cornerRadius {
				dx := relX - (float64(size) - cornerRadius)
				dy := relY - cornerRadius
				if dx*dx + dy*dy > cornerRadius*cornerRadius {
					inCorner = true
				}
			}
			// Bottom-left corner
			if relX < cornerRadius && relY >= float64(size)-cornerRadius {
				dx := relX - cornerRadius
				dy := relY - (float64(size) - cornerRadius)
				if dx*dx + dy*dy > cornerRadius*cornerRadius {
					inCorner = true
				}
			}
			// Bottom-right corner
			if relX >= float64(size)-cornerRadius && relY >= float64(size)-cornerRadius {
				dx := relX - (float64(size) - cornerRadius)
				dy := relY - (float64(size) - cornerRadius)
				if dx*dx + dy*dy > cornerRadius*cornerRadius {
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