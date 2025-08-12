package myqrcode

import (
	"rsc.io/qr/gf256"
)

func addErrorCorrection(data []int, version int, level ErrorCorrectionLevel) []byte {
	info := getVersionInfo(version)
	
	ecCodewords := info.ECCodewordsPerBlock[level]
	dataCodewords := info.DataCodewordsPerBlock[level]
	numBlocks := info.NumBlocks[level]
	
	dataBits := bitsToBytes(data)
	
	if len(dataBits) > dataCodewords*numBlocks {
		dataBits = dataBits[:dataCodewords*numBlocks]
	}
	
	for len(dataBits) < dataCodewords*numBlocks {
		dataBits = append(dataBits, 0)
	}
	
	var allCodewords []byte
	
	dataBlocks := make([][]byte, numBlocks)
	ecBlocks := make([][]byte, numBlocks)
	
	for i := 0; i < numBlocks; i++ {
		start := i * dataCodewords
		end := start + dataCodewords
		if end > len(dataBits) {
			end = len(dataBits)
		}
		
		blockData := make([]byte, dataCodewords)
		copy(blockData, dataBits[start:end])
		dataBlocks[i] = blockData
		
		ecBlocks[i] = generateErrorCorrection(blockData, ecCodewords)
	}
	
	for i := 0; i < dataCodewords; i++ {
		for j := 0; j < numBlocks; j++ {
			if i < len(dataBlocks[j]) {
				allCodewords = append(allCodewords, dataBlocks[j][i])
			}
		}
	}
	
	for i := 0; i < ecCodewords; i++ {
		for j := 0; j < numBlocks; j++ {
			if i < len(ecBlocks[j]) {
				allCodewords = append(allCodewords, ecBlocks[j][i])
			}
		}
	}
	
	return allCodewords
}

func generateErrorCorrection(data []byte, ecCodewords int) []byte {
	// Create QR code Galois Field (polynomial 0x11d, generator 2)
	field := gf256.NewField(0x11d, 2)
	
	// Create Reed-Solomon encoder
	encoder := gf256.NewRSEncoder(field, ecCodewords)
	
	// Generate error correction bytes
	result := make([]byte, ecCodewords)
	encoder.ECC(data, result)
	
	return result
}

func bitsToBytes(bits []int) []byte {
	for len(bits)%8 != 0 {
		bits = append(bits, 0)
	}
	
	bytes := make([]byte, len(bits)/8)
	for i := 0; i < len(bytes); i++ {
		for j := 0; j < 8; j++ {
			if bits[i*8+j] == 1 {
				bytes[i] |= 1 << (7 - j)
			}
		}
	}
	
	return bytes
}

func bytesToBits(data []byte) []int {
	var bits []int
	for _, b := range data {
		for i := 7; i >= 0; i-- {
			bits = append(bits, int((b>>i)&1))
		}
	}
	return bits
}