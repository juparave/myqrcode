package myqrcode

type VersionInfo struct {
	Version           int
	Size              int
	TotalCodewords    int
	ECCodewordsPerBlock []int
	DataCodewordsPerBlock []int
	NumBlocks         []int
}

var versionTable = []VersionInfo{
	{1, 21, 26, []int{7, 10, 13, 17}, []int{19, 16, 13, 9}, []int{1, 1, 1, 1}},
	{2, 25, 44, []int{10, 16, 22, 28}, []int{34, 28, 22, 16}, []int{1, 1, 1, 1}},
	{3, 29, 70, []int{15, 26, 36, 44}, []int{55, 44, 34, 26}, []int{1, 1, 2, 2}},
	{4, 33, 100, []int{20, 36, 52, 64}, []int{80, 64, 48, 36}, []int{1, 2, 2, 4}},
	{5, 37, 134, []int{26, 48, 72, 88}, []int{108, 86, 62, 46}, []int{1, 2, 4, 4}},
	{6, 41, 172, []int{36, 64, 96, 112}, []int{136, 108, 76, 60}, []int{2, 4, 4, 4}},
	{7, 45, 196, []int{40, 72, 108, 130}, []int{156, 124, 88, 66}, []int{2, 4, 6, 5}},
	{8, 49, 242, []int{48, 88, 132, 156}, []int{194, 154, 110, 86}, []int{2, 4, 6, 6}},
	{9, 53, 292, []int{60, 110, 160, 192}, []int{232, 182, 132, 100}, []int{2, 5, 8, 8}},
	{10, 57, 346, []int{72, 130, 192, 224}, []int{274, 216, 154, 122}, []int{4, 5, 8, 8}},
}

func determineVersion(data string, mode EncodingMode, level ErrorCorrectionLevel) int {
	dataLength := len(data)
	
	for _, info := range versionTable {
		capacity := getDataCapacity(info.Version, mode, level)
		if dataLength <= capacity {
			return info.Version
		}
	}
	
	return 10
}

func getDataCapacity(version int, mode EncodingMode, level ErrorCorrectionLevel) int {
	if version < 1 || version > len(versionTable) {
		return 0
	}
	
	info := versionTable[version-1]
	totalDataCodewords := info.TotalCodewords - info.ECCodewordsPerBlock[level]
	
	modeBits := 4
	countBits := getCharCountBits(mode, version)
	headerBits := modeBits + countBits
	
	availableBits := (totalDataCodewords * 8) - headerBits
	
	switch mode {
	case Numeric:
		return (availableBits / 10) * 3
	case Alphanumeric:
		return (availableBits / 11) * 2
	case Byte:
		return availableBits / 8
	}
	
	return 0
}

func getVersionInfo(version int) VersionInfo {
	if version < 1 || version > len(versionTable) {
		return versionTable[0]
	}
	return versionTable[version-1]
}