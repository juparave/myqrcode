package myqrcode

// BlockInfo contains information about a group of blocks for a specific version and error correction level.
type BlockInfo struct {
	NumBlocks      int
	DataCodewords  int
	TotalCodewords int
}

// VersionInfo contains all information for a specific QR code version.
type VersionInfo struct {
	Version     int
	Size        int
	ECBlockInfo [][]BlockInfo // [level][group]
}

// versionTable holds the information for QR code versions 1-40.
var versionTable = []VersionInfo{
	{1, 21, [][]BlockInfo{{{1, 19, 26}}, {{1, 16, 26}}, {{1, 13, 26}}, {{1, 9, 26}}}},
	{2, 25, [][]BlockInfo{{{1, 34, 44}}, {{1, 28, 44}}, {{1, 22, 44}}, {{1, 16, 44}}}},
	{3, 29, [][]BlockInfo{{{1, 55, 70}}, {{1, 44, 70}}, {{2, 17, 35}}, {{2, 13, 35}}}},
	{4, 33, [][]BlockInfo{{{1, 80, 100}}, {{2, 32, 50}}, {{2, 24, 50}}, {{4, 9, 25}}}},
	{5, 37, [][]BlockInfo{{{1, 108, 134}}, {{2, 43, 67}}, {{2, 15, 33}, {2, 16, 34}}, {{2, 11, 33}, {2, 12, 34}}}},
	{6, 41, [][]BlockInfo{{{2, 68, 86}}, {{4, 27, 43}}, {{4, 19, 43}}, {{4, 15, 43}}}},
	{7, 45, [][]BlockInfo{{{2, 78, 98}}, {{4, 31, 49}}, {{2, 14, 32}, {4, 15, 33}}, {{4, 13, 39}, {1, 14, 40}}}},
	{8, 49, [][]BlockInfo{{{2, 97, 121}}, {{2, 38, 60}, {2, 39, 61}}, {{4, 18, 40}, {2, 19, 41}}, {{4, 14, 40}, {2, 15, 41}}}},
	{9, 53, [][]BlockInfo{{{2, 116, 146}}, {{3, 36, 58}, {2, 37, 59}}, {{4, 16, 36}, {4, 17, 37}}, {{4, 12, 36}, {4, 13, 37}}}},
	{10, 57, [][]BlockInfo{{{2, 68, 86}, {2, 69, 87}}, {{4, 43, 69}, {1, 44, 70}}, {{6, 19, 43}, {2, 20, 44}}, {{6, 15, 43}, {2, 16, 44}}}},
	{11, 61, [][]BlockInfo{{{4, 81, 101}}, {{1, 50, 80}, {4, 51, 81}}, {{4, 22, 50}, {4, 23, 51}}, {{3, 12, 36}, {8, 13, 37}}}},
	{12, 65, [][]BlockInfo{{{2, 92, 116}, {2, 93, 117}}, {{6, 36, 58}, {2, 37, 59}}, {{4, 20, 46}, {6, 21, 47}}, {{7, 14, 42}, {4, 15, 43}}}},
	{13, 69, [][]BlockInfo{{{4, 107, 133}}, {{8, 37, 59}, {1, 38, 60}}, {{8, 20, 44}, {4, 21, 45}}, {{12, 11, 33}, {4, 12, 34}}}},
	{14, 73, [][]BlockInfo{{{3, 115, 145}, {1, 116, 146}}, {{4, 40, 64}, {5, 41, 65}}, {{11, 16, 36}, {5, 17, 37}}, {{11, 12, 36}, {5, 13, 37}}}},
	{15, 77, [][]BlockInfo{{{5, 87, 109}, {1, 88, 110}}, {{5, 41, 65}, {5, 42, 66}}, {{5, 24, 54}, {7, 25, 55}}, {{11, 12, 36}, {7, 13, 37}}}},
	{16, 81, [][]BlockInfo{{{5, 98, 122}, {1, 99, 123}}, {{7, 45, 73}, {3, 46, 74}}, {{15, 19, 43}, {2, 20, 44}}, {{3, 15, 45}, {13, 16, 46}}}},
	{17, 85, [][]BlockInfo{{{1, 107, 135}, {5, 108, 136}}, {{10, 46, 74}, {1, 47, 75}}, {{1, 22, 50}, {15, 23, 51}}, {{2, 14, 42}, {17, 15, 43}}}},
	{18, 89, [][]BlockInfo{{{5, 120, 150}, {1, 121, 151}}, {{9, 43, 69}, {4, 44, 70}}, {{17, 22, 50}, {1, 23, 51}}, {{2, 14, 42}, {19, 15, 43}}}},
	{19, 93, [][]BlockInfo{{{3, 113, 141}, {4, 114, 142}}, {{3, 44, 70}, {11, 45, 71}}, {{17, 21, 47}, {4, 22, 48}}, {{9, 13, 39}, {16, 14, 40}}}},
	{20, 97, [][]BlockInfo{{{3, 107, 135}, {5, 108, 136}}, {{3, 41, 67}, {13, 42, 68}}, {{15, 24, 54}, {5, 25, 55}}, {{15, 15, 43}, {10, 16, 44}}}},
	{21, 101, [][]BlockInfo{{{4, 116, 144}, {4, 117, 145}}, {{17, 42, 68}}, {{17, 22, 50}, {6, 23, 51}}, {{19, 16, 46}, {6, 17, 47}}}},
	{22, 105, [][]BlockInfo{{{2, 111, 139}, {7, 112, 140}}, {{17, 46, 74}}, {{7, 24, 54}, {16, 25, 55}}, {{34, 13, 37}}}},
	{23, 109, [][]BlockInfo{{{4, 121, 151}, {5, 122, 152}}, {{4, 47, 75}, {14, 48, 76}}, {{11, 24, 54}, {14, 25, 55}}, {{16, 15, 45}, {14, 16, 46}}}},
	{24, 113, [][]BlockInfo{{{6, 117, 147}, {4, 118, 148}}, {{6, 45, 73}, {14, 46, 74}}, {{11, 24, 54}, {16, 25, 55}}, {{30, 16, 46}, {2, 17, 47}}}},
	{25, 117, [][]BlockInfo{{{8, 106, 132}, {4, 107, 133}}, {{8, 47, 75}, {13, 48, 76}}, {{7, 24, 54}, {22, 25, 55}}, {{22, 15, 45}, {13, 16, 46}}}},
	{26, 121, [][]BlockInfo{{{10, 114, 142}, {2, 115, 143}}, {{19, 46, 74}, {4, 47, 75}}, {{28, 22, 50}, {6, 23, 51}}, {{33, 16, 46}, {4, 17, 47}}}},
	{27, 125, [][]BlockInfo{{{8, 122, 152}, {4, 123, 153}}, {{22, 45, 73}, {3, 46, 74}}, {{8, 23, 53}, {26, 24, 54}}, {{12, 15, 45}, {28, 16, 46}}}},
	{28, 129, [][]BlockInfo{{{3, 117, 147}, {10, 118, 148}}, {{3, 45, 73}, {23, 46, 74}}, {{4, 24, 54}, {31, 25, 55}}, {{11, 15, 45}, {31, 16, 46}}}},
	{29, 133, [][]BlockInfo{{{7, 116, 146}, {7, 117, 147}}, {{21, 45, 73}, {7, 46, 74}}, {{1, 23, 53}, {37, 24, 54}}, {{19, 15, 45}, {26, 16, 46}}}},
	{30, 137, [][]BlockInfo{{{5, 115, 145}, {10, 116, 146}}, {{19, 47, 75}, {10, 48, 76}}, {{15, 24, 54}, {25, 25, 55}}, {{23, 15, 45}, {25, 16, 46}}}},
	{31, 141, [][]BlockInfo{{{13, 115, 145}, {3, 116, 146}}, {{2, 46, 74}, {29, 47, 75}}, {{42, 24, 54}, {1, 25, 55}}, {{23, 15, 45}, {28, 16, 46}}}},
	{32, 145, [][]BlockInfo{{{17, 115, 145}}, {{10, 46, 74}, {23, 47, 75}}, {{10, 24, 54}, {35, 25, 55}}, {{19, 15, 45}, {35, 16, 46}}}},
	{33, 149, [][]BlockInfo{{{17, 115, 145}, {1, 116, 146}}, {{14, 46, 74}, {21, 47, 75}}, {{29, 24, 54}, {19, 25, 55}}, {{11, 15, 45}, {46, 16, 46}}}},
	{34, 153, [][]BlockInfo{{{13, 115, 145}, {6, 116, 146}}, {{14, 46, 74}, {23, 47, 75}}, {{44, 24, 54}, {7, 25, 55}}, {{59, 16, 46}, {1, 17, 47}}}},
	{35, 157, [][]BlockInfo{{{12, 121, 151}, {7, 122, 152}}, {{12, 47, 75}, {26, 48, 76}}, {{39, 24, 54}, {14, 25, 55}}, {{22, 15, 45}, {41, 16, 46}}}},
	{36, 161, [][]BlockInfo{{{6, 121, 151}, {14, 122, 152}}, {{6, 47, 75}, {34, 48, 76}}, {{46, 24, 54}, {10, 25, 55}}, {{2, 15, 45}, {64, 16, 46}}}},
	{37, 165, [][]BlockInfo{{{17, 122, 152}, {4, 123, 153}}, {{29, 46, 74}, {14, 47, 75}}, {{49, 24, 54}, {10, 25, 55}}, {{24, 15, 45}, {46, 16, 46}}}},
	{38, 169, [][]BlockInfo{{{4, 122, 152}, {18, 123, 153}}, {{13, 46, 74}, {32, 47, 75}}, {{48, 24, 54}, {14, 25, 55}}, {{42, 15, 45}, {32, 16, 46}}}},
	{39, 173, [][]BlockInfo{{{20, 117, 147}, {4, 118, 148}}, {{40, 47, 75}, {7, 48, 76}}, {{43, 24, 54}, {22, 25, 55}}, {{10, 15, 45}, {67, 16, 46}}}},
	{40, 177, [][]BlockInfo{{{19, 118, 148}, {6, 119, 149}}, {{18, 47, 75}, {31, 48, 76}}, {{34, 24, 54}, {34, 25, 55}}, {{20, 15, 45}, {61, 16, 46}}}},
}

func getVersionInfo(version int) VersionInfo {
	if version < 1 || version > len(versionTable) {
		// Default to version 1 if out of range
		return versionTable[0]
	}
	return versionTable[version-1]
}

func determineVersion(data string, mode EncodingMode, level ErrorCorrectionLevel) int {
	dataLength := len(data)

	for _, info := range versionTable {
		capacity := getDataCapacity(info.Version, mode, level)
		if dataLength <= capacity {
			return info.Version
		}
	}

	return 40 // Return max version if data is too large
}

func getDataCapacity(version int, mode EncodingMode, level ErrorCorrectionLevel) int {
	if version < 1 || version > len(versionTable) {
		return 0
	}

	info := getVersionInfo(version)
	blockInfo := info.ECBlockInfo[level]

	totalDataCodewords := 0
	for _, group := range blockInfo {
		totalDataCodewords += group.NumBlocks * group.DataCodewords
	}

	modeBits := 4
	countBits := getCharCountBits(mode, version)
	headerBits := modeBits + countBits

	availableBits := (totalDataCodewords * 8) - headerBits

	switch mode {
	case Numeric:
		// 10 bits for every 3 digits
		return (availableBits / 10) * 3
	case Alphanumeric:
		// 11 bits for every 2 characters
		return (availableBits / 11) * 2
	case Byte:
		// 8 bits for every character
		return availableBits / 8
	}

	return 0
}

func getCharCountBits(mode EncodingMode, version int) int {
	if version <= 9 {
		switch mode {
		case Numeric:
			return 10
		case Alphanumeric:
			return 9
		case Byte:
			return 8
		}
	} else if version <= 26 {
		switch mode {
		case Numeric:
			return 12
		case Alphanumeric:
			return 11
		case Byte:
			return 16
		}
	} else { // version 27-40
		switch mode {
		case Numeric:
			return 14
		case Alphanumeric:
			return 13
		case Byte:
			return 16
		}
	}
	return 0
}
