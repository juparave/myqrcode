package myqrcode

// GetModuleNeighbors analyzes the 8 surrounding modules and returns neighbor context
func GetModuleNeighbors(matrix [][]bool, row, col int) *ActiveWithNeighbors {
	rows := len(matrix)
	cols := len(matrix[0])

	// Helper function to safely check if a module is active
	isActive := func(r, c int) bool {
		if r < 0 || r >= rows || c < 0 || c >= cols {
			return false // Out of bounds = inactive
		}
		return matrix[r][c]
	}

	return &ActiveWithNeighbors{
		NW: isActive(row-1, col-1), // Northwest
		N:  isActive(row-1, col),   // North
		NE: isActive(row-1, col+1), // Northeast
		W:  isActive(row, col-1),   // West
		Me: isActive(row, col),     // Center (this module)
		E:  isActive(row, col+1),   // East
		SW: isActive(row+1, col-1), // Southwest
		S:  isActive(row+1, col),   // South
		SE: isActive(row+1, col+1), // Southeast
	}
}

// IsFinderPattern checks if the current position is part of a finder pattern
func IsFinderPattern(row, col, size int) bool {
	// Top-left finder pattern (0,0 to 6,6)
	if row >= 0 && row < 7 && col >= 0 && col < 7 {
		return true
	}

	// Top-right finder pattern
	if row >= 0 && row < 7 && col >= size-7 && col < size {
		return true
	}

	// Bottom-left finder pattern
	if row >= size-7 && row < size && col >= 0 && col < 7 {
		return true
	}

	return false
}

// IsTimingPattern checks if the current position is part of a timing pattern
func IsTimingPattern(row, col int) bool {
	// Row 6 (horizontal timing pattern)
	if row == 6 && col >= 8 {
		return true
	}

	// Column 6 (vertical timing pattern)
	if col == 6 && row >= 8 {
		return true
	}

	return false
}

// GetFinderPatternNeighbors creates special neighbor context for finder patterns
// This ensures finder patterns get properly rounded corners only on outer edges
func GetFinderPatternNeighbors(matrix [][]bool, row, col int) *ActiveWithNeighbors {
	size := len(matrix)

	// Determine which finder pattern we're in
	var patternRow, patternCol int
	var patternSize = 7

	if row >= 0 && row < 7 && col >= 0 && col < 7 {
		// Top-left finder pattern
		patternRow, patternCol = row, col
	} else if row >= 0 && row < 7 && col >= size-7 && col < size {
		// Top-right finder pattern
		patternRow, patternCol = row, col-(size-7)
	} else if row >= size-7 && row < size && col >= 0 && col < 7 {
		// Bottom-left finder pattern
		patternRow, patternCol = row-(size-7), col
	} else {
		// Not in a finder pattern, use regular neighbor detection
		return GetModuleNeighbors(matrix, row, col)
	}

	// Helper to check if position is within finder pattern and active
	isFinderActive := func(r, c int) bool {
		if r < 0 || r >= patternSize || c < 0 || c >= patternSize {
			return false
		}

		// Finder pattern structure:
		// - Outer ring (border)
		// - Inner white ring
		// - Center 3x3 black square
		if r == 0 || r == 6 || c == 0 || c == 6 {
			return true // Outer border
		}
		if r >= 2 && r <= 4 && c >= 2 && c <= 4 {
			return true // Center square
		}
		return false // Inner white ring
	}

	return &ActiveWithNeighbors{
		NW: isFinderActive(patternRow-1, patternCol-1),
		N:  isFinderActive(patternRow-1, patternCol),
		NE: isFinderActive(patternRow-1, patternCol+1),
		W:  isFinderActive(patternRow, patternCol-1),
		Me: isFinderActive(patternRow, patternCol),
		E:  isFinderActive(patternRow, patternCol+1),
		SW: isFinderActive(patternRow+1, patternCol-1),
		S:  isFinderActive(patternRow+1, patternCol),
		SE: isFinderActive(patternRow+1, patternCol+1),
	}
}
