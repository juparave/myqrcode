package myqrcode

import (
	"errors"
	"strconv"
	"strings"
)

var alphanumericChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:"

func detectMode(data string) EncodingMode {
	if isNumeric(data) {
		return Numeric
	}
	if isAlphanumeric(data) {
		return Alphanumeric
	}
	return Byte
}

func isNumeric(data string) bool {
	for _, c := range data {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func isAlphanumeric(data string) bool {
	for _, c := range data {
		if !strings.ContainsRune(alphanumericChars, c) {
			return false
		}
	}
	return true
}

func encodeData(data string, mode EncodingMode, version int) ([]int, error) {
	switch mode {
	case Numeric:
		return encodeNumeric(data, version)
	case Alphanumeric:
		return encodeAlphanumeric(data, version)
	case Byte:
		return encodeByte(data, version)
	}
	return nil, errors.New("unsupported encoding mode")
}

func encodeNumeric(data string, version int) ([]int, error) {
	var bits []int
	
	// Mode indicator (4 bits: 0001)
	bits = append(bits, 0, 0, 0, 1)
	
	// Character count indicator
	countBits := getCharCountBits(Numeric, version)
	count := len(data)
	for i := countBits - 1; i >= 0; i-- {
		bits = append(bits, (count>>i)&1)
	}
	
	// Data encoding
	for i := 0; i < len(data); i += 3 {
		group := data[i:min(i+3, len(data))]
		val, _ := strconv.Atoi(group)
		
		bitCount := 10
		if len(group) == 2 {
			bitCount = 7
		} else if len(group) == 1 {
			bitCount = 4
		}
		
		for j := bitCount - 1; j >= 0; j-- {
			bits = append(bits, (val>>j)&1)
		}
	}
	
	return bits, nil
}

func encodeAlphanumeric(data string, version int) ([]int, error) {
	var bits []int
	
	// Mode indicator (4 bits: 0010)
	bits = append(bits, 0, 0, 1, 0)
	
	// Character count indicator
	countBits := getCharCountBits(Alphanumeric, version)
	count := len(data)
	for i := countBits - 1; i >= 0; i-- {
		bits = append(bits, (count>>i)&1)
	}
	
	// Data encoding
	for i := 0; i < len(data); i += 2 {
		if i+1 < len(data) {
			val1 := strings.IndexRune(alphanumericChars, rune(data[i]))
			val2 := strings.IndexRune(alphanumericChars, rune(data[i+1]))
			val := val1*45 + val2
			
			for j := 10; j >= 0; j-- {
				bits = append(bits, (val>>j)&1)
			}
		} else {
			val := strings.IndexRune(alphanumericChars, rune(data[i]))
			for j := 5; j >= 0; j-- {
				bits = append(bits, (val>>j)&1)
			}
		}
	}
	
	return bits, nil
}

func encodeByte(data string, version int) ([]int, error) {
	var bits []int
	
	// Mode indicator (4 bits: 0100)
	bits = append(bits, 0, 1, 0, 0)
	
	// Character count indicator
	countBits := getCharCountBits(Byte, version)
	count := len(data)
	for i := countBits - 1; i >= 0; i-- {
		bits = append(bits, (count>>i)&1)
	}
	
	// Data encoding
	for _, char := range []byte(data) {
		for i := 7; i >= 0; i-- {
			bits = append(bits, int((char>>i)&1))
		}
	}
	
	return bits, nil
}

func getCharCountBits(mode EncodingMode, version int) int {
	switch mode {
	case Numeric:
		if version <= 9 {
			return 10
		} else if version <= 26 {
			return 12
		}
		return 14
	case Alphanumeric:
		if version <= 9 {
			return 9
		} else if version <= 26 {
			return 11
		}
		return 13
	case Byte:
		if version <= 9 {
			return 8
		}
		return 16
	}
	return 0
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func addTerminatorAndPadding(bits []int, version int, level ErrorCorrectionLevel) []int {
	info := getVersionInfo(version)
	maxBits := info.DataCodewordsPerBlock[level] * info.NumBlocks[level] * 8
	
	// Add terminator (up to 4 zeros)
	terminatorBits := min(4, maxBits-len(bits))
	for i := 0; i < terminatorBits; i++ {
		bits = append(bits, 0)
	}
	
	// Pad to byte boundary
	for len(bits)%8 != 0 {
		bits = append(bits, 0)
	}
	
	// Add padding bytes
	padBytes := []int{
		1, 1, 1, 0, 1, 1, 0, 0, // 0xEC
		0, 0, 0, 1, 0, 0, 0, 1, // 0x11
	}
	
	padIndex := 0
	for len(bits) < maxBits {
		bits = append(bits, padBytes[padIndex])
		padIndex = (padIndex + 1) % 16
	}
	
	// Truncate if too long
	if len(bits) > maxBits {
		bits = bits[:maxBits]
	}
	
	return bits
}