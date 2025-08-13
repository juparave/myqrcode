package myqrcode

func placeData(matrix *Matrix, data []byte) {
	bits := bytesToBits(data)
	bitIndex := 0

	size := matrix.Size
	direction := -1 // Start going up

	for col := size - 1; col >= 0; col -= 2 {
		if col == 6 {
			col-- // Skip timing column
		}

		for row := 0; row < size; row++ {
			var actualRow int
			if direction == -1 {
				actualRow = size - 1 - row
			} else {
				actualRow = row
			}

			// Place data in the two columns
			for c := 0; c < 2; c++ {
				x := col - c
				y := actualRow

				if !matrix.IsReserved(x, y) {
					var bit bool
					if bitIndex < len(bits) {
						bit = bits[bitIndex] == 1
						bitIndex++
					} else {
						bit = false
					}
					matrix.Set(x, y, bit)
				}
			}
		}

		direction *= -1 // Change direction
	}
}

func applyMask(matrix *Matrix, maskPattern int) *Matrix {
	size := matrix.Size
	masked := NewMatrix(size)

	// Copy reserved areas
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			masked.Reserve[y][x] = matrix.Reserve[y][x]
			masked.Modules[y][x] = matrix.Modules[y][x]
		}
	}

	// Apply mask to non-reserved areas
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if !masked.IsReserved(x, y) {
				if shouldMask(x, y, maskPattern) {
					masked.Set(x, y, !matrix.Get(x, y))
				}
			}
		}
	}

	return masked
}

func shouldMask(x, y, pattern int) bool {
	switch pattern {
	case 0:
		return (x+y)%2 == 0
	case 1:
		return y%2 == 0
	case 2:
		return x%3 == 0
	case 3:
		return (x+y)%3 == 0
	case 4:
		return (y/2+x/3)%2 == 0
	case 5:
		return (x*y)%2+(x*y)%3 == 0
	case 6:
		return ((x*y)%2+(x*y)%3)%2 == 0
	case 7:
		return ((x+y)%2+(x*y)%3)%2 == 0
	default:
		return false
	}
}

func evaluateMask(matrix *Matrix) int {
	size := matrix.Size
	penalty := 0

	// Rule 1: Adjacent modules in row/column in same color
	for y := 0; y < size; y++ {
		count := 1
		for x := 1; x < size; x++ {
			if matrix.Get(x, y) == matrix.Get(x-1, y) {
				count++
			} else {
				if count >= 5 {
					penalty += count - 2
				}
				count = 1
			}
		}
		if count >= 5 {
			penalty += count - 2
		}
	}

	for x := 0; x < size; x++ {
		count := 1
		for y := 1; y < size; y++ {
			if matrix.Get(x, y) == matrix.Get(x, y-1) {
				count++
			} else {
				if count >= 5 {
					penalty += count - 2
				}
				count = 1
			}
		}
		if count >= 5 {
			penalty += count - 2
		}
	}

	// Rule 2: Block of modules in same color (2x2)
	for y := 0; y < size-1; y++ {
		for x := 0; x < size-1; x++ {
			color := matrix.Get(x, y)
			if matrix.Get(x+1, y) == color &&
				matrix.Get(x, y+1) == color &&
				matrix.Get(x+1, y+1) == color {
				penalty += 3
			}
		}
	}

	// Rule 3: Finder-like patterns
	finderPattern := []bool{true, false, true, true, true, false, true}
	whitePattern := []bool{false, false, false, false, true, false, true, true, true, false, true}

	for y := 0; y < size; y++ {
		for x := 0; x <= size-11; x++ {
			match := true
			for i := 0; i < 11; i++ {
				expected := false
				if i < 4 {
					expected = false
				} else if i < 11 {
					expected = finderPattern[i-4]
				}
				if matrix.Get(x+i, y) != expected {
					match = false
					break
				}
			}
			if match {
				penalty += 40
			}

			match = true
			for i := 0; i < 11; i++ {
				if matrix.Get(x+i, y) != whitePattern[i] {
					match = false
					break
				}
			}
			if match {
				penalty += 40
			}
		}
	}

	// Rule 4: Proportion of dark modules
	darkCount := 0
	totalCount := size * size

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if matrix.Get(x, y) {
				darkCount++
			}
		}
	}

	percentage := (darkCount * 100) / totalCount
	deviation := abs(percentage - 50)
	penalty += (deviation / 5) * 10

	return penalty
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func selectBestMask(matrix *Matrix, level ErrorCorrectionLevel) (*Matrix, int) {
	bestMask := 0
	bestPenalty := int(^uint(0) >> 1) // Max int
	var bestMatrix *Matrix

	for mask := 0; mask < 8; mask++ {
		// Apply mask
		maskedMatrix := applyMask(matrix, mask)

		// Add format info with this mask
		maskedMatrix.AddFormatInfo(level, mask)

		// Evaluate penalty
		penalty := evaluateMask(maskedMatrix)

		if penalty < bestPenalty {
			bestPenalty = penalty
			bestMask = mask
			bestMatrix = maskedMatrix
		}
	}

	return bestMatrix, bestMask
}
