package myqrcode

import (
	"math"
)

type LogoPlacement struct {
	X, Y   int
	Width  int
	Height int
}

func calculateLogoPlacement(matrix *Matrix, logoSize int) LogoPlacement {
	size := matrix.Size
	
	// Calculate logo dimensions based on percentage of QR size
	logoWidth := (size * logoSize) / 100
	logoHeight := logoWidth
	
	// Ensure logo size is odd for better centering
	if logoWidth%2 == 0 {
		logoWidth++
	}
	if logoHeight%2 == 0 {
		logoHeight++
	}
	
	// Center the logo
	x := (size - logoWidth) / 2
	y := (size - logoHeight) / 2
	
	return LogoPlacement{
		X:      x,
		Y:      y,
		Width:  logoWidth,
		Height: logoHeight,
	}
}

func isLogoArea(x, y int, placement LogoPlacement) bool {
	return x >= placement.X && x < placement.X+placement.Width &&
		y >= placement.Y && y < placement.Y+placement.Height
}

func isCriticalArea(x, y, size int) bool {
	// Finder patterns (corners)
	if (x < 9 && y < 9) ||                    // Top-left
		(x >= size-8 && y < 9) ||             // Top-right  
		(x < 9 && y >= size-8) {              // Bottom-left
		return true
	}
	
	// Timing patterns
	if x == 6 || y == 6 {
		return true
	}
	
	// Dark module
	if x == 8 && y == 4*((size-17)/4)+9 {
		return true
	}
	
	// Format information areas
	if (x == 8 && (y < 9 || y >= size-8)) ||
		(y == 8 && (x < 9 || x >= size-7)) {
		return true
	}
	
	return false
}

func optimizeLogoPlacement(matrix *Matrix, logoSize int) LogoPlacement {
	placement := calculateLogoPlacement(matrix, logoSize)
	size := matrix.Size
	
	// Adjust placement to avoid critical areas if possible
	maxOffset := 3
	bestPlacement := placement
	minCriticalOverlap := countCriticalOverlap(placement, size)
	
	for offsetX := -maxOffset; offsetX <= maxOffset; offsetX++ {
		for offsetY := -maxOffset; offsetY <= maxOffset; offsetY++ {
			newPlacement := LogoPlacement{
				X:      placement.X + offsetX,
				Y:      placement.Y + offsetY,
				Width:  placement.Width,
				Height: placement.Height,
			}
			
			// Check if placement is within bounds
			if newPlacement.X >= 0 && newPlacement.Y >= 0 &&
				newPlacement.X+newPlacement.Width <= size &&
				newPlacement.Y+newPlacement.Height <= size {
				
				criticalOverlap := countCriticalOverlap(newPlacement, size)
				if criticalOverlap < minCriticalOverlap {
					minCriticalOverlap = criticalOverlap
					bestPlacement = newPlacement
				}
			}
		}
	}
	
	return bestPlacement
}

func countCriticalOverlap(placement LogoPlacement, size int) int {
	count := 0
	for y := placement.Y; y < placement.Y+placement.Height; y++ {
		for x := placement.X; x < placement.X+placement.Width; x++ {
			if isCriticalArea(x, y, size) {
				count++
			}
		}
	}
	return count
}

func calculateLogoErrorCorrection(placement LogoPlacement, totalDataBits int) float64 {
	logoArea := placement.Width * placement.Height
	totalArea := int(math.Pow(float64(placement.Width+placement.Height), 2)) // Approximate QR area
	
	// Estimate percentage of data that will be obscured
	obscuredPercentage := float64(logoArea) / float64(totalArea)
	
	// Add safety margin
	requiredCorrection := obscuredPercentage * 1.5
	
	// Ensure we don't exceed maximum error correction capability
	if requiredCorrection > 0.30 { // QR High level is ~30%
		requiredCorrection = 0.30
	}
	
	return requiredCorrection
}

func adjustErrorCorrectionForLogo(level ErrorCorrectionLevel, placement LogoPlacement, dataSize int) ErrorCorrectionLevel {
	requiredCorrection := calculateLogoErrorCorrection(placement, dataSize*8)
	
	// Map correction percentages to levels
	// Low: ~7%, Medium: ~15%, Quartile: ~25%, High: ~30%
	if requiredCorrection > 0.25 {
		return High
	} else if requiredCorrection > 0.15 {
		return Quartile  
	} else if requiredCorrection > 0.07 {
		return Medium
	}
	
	return level // Keep original if no adjustment needed
}

func reserveLogoArea(matrix *Matrix, placement LogoPlacement) {
	for y := placement.Y; y < placement.Y+placement.Height; y++ {
		for x := placement.X; x < placement.X+placement.Width; x++ {
			if x >= 0 && x < matrix.Size && y >= 0 && y < matrix.Size {
				matrix.SetReserved(x, y)
				// Set to false (white) to create space for logo
				matrix.Set(x, y, false)
			}
		}
	}
}