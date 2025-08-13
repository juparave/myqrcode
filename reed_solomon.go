package myqrcode

import (
	"rsc.io/qr/gf256"
)

func addErrorCorrection(data []int, version int, level ErrorCorrectionLevel) []byte {
	info := getVersionInfo(version)
	blockInfo := info.ECBlockInfo[level]

	dataBits := bitsToBytes(data)

	var dataBlocks [][]byte
	var ecBlocks [][]byte
	offset := 0

	totalDataCodewords := 0
	for _, group := range blockInfo {
		totalDataCodewords += group.NumBlocks * group.DataCodewords
	}

	if len(dataBits) > totalDataCodewords {
		dataBits = dataBits[:totalDataCodewords]
	}

	for len(dataBits) < totalDataCodewords {
		dataBits = append(dataBits, 0)
	}

	for _, group := range blockInfo {
		for i := 0; i < group.NumBlocks; i++ {
			dataCodewords := group.DataCodewords
			ecCodewords := group.TotalCodewords - dataCodewords

			blockData := make([]byte, dataCodewords)
			copy(blockData, dataBits[offset:offset+dataCodewords])
			offset += dataCodewords

			dataBlocks = append(dataBlocks, blockData)
			ecBlocks = append(ecBlocks, generateErrorCorrection(blockData, ecCodewords))
		}
	}

	var allCodewords []byte
	maxDataLen := 0
	for _, block := range dataBlocks {
		if len(block) > maxDataLen {
			maxDataLen = len(block)
		}
	}

	for i := 0; i < maxDataLen; i++ {
		for _, block := range dataBlocks {
			if i < len(block) {
				allCodewords = append(allCodewords, block[i])
			}
		}
	}

	maxEcLen := 0
	for _, block := range ecBlocks {
		if len(block) > maxEcLen {
			maxEcLen = len(block)
		}
	}

	for i := 0; i < maxEcLen; i++ {
		for _, block := range ecBlocks {
			if i < len(block) {
				allCodewords = append(allCodewords, block[i])
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
