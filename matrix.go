package myqrcode

type Matrix struct {
	Size    int
	Modules [][]bool
	Reserve [][]bool
}

func NewMatrix(size int) *Matrix {
	modules := make([][]bool, size)
	reserve := make([][]bool, size)
	for i := range modules {
		modules[i] = make([]bool, size)
		reserve[i] = make([]bool, size)
	}
	return &Matrix{
		Size:    size,
		Modules: modules,
		Reserve: reserve,
	}
}

func (m *Matrix) Set(x, y int, value bool) {
	if x >= 0 && x < m.Size && y >= 0 && y < m.Size {
		m.Modules[y][x] = value
	}
}

func (m *Matrix) Get(x, y int) bool {
	if x >= 0 && x < m.Size && y >= 0 && y < m.Size {
		return m.Modules[y][x]
	}
	return false
}

func (m *Matrix) SetReserved(x, y int) {
	if x >= 0 && x < m.Size && y >= 0 && y < m.Size {
		m.Reserve[y][x] = true
	}
}

func (m *Matrix) IsReserved(x, y int) bool {
	if x >= 0 && x < m.Size && y >= 0 && y < m.Size {
		return m.Reserve[y][x]
	}
	return true
}

func (m *Matrix) AddFinderPatterns() {
	positions := [][2]int{{0, 0}, {m.Size - 7, 0}, {0, m.Size - 7}}

	for _, pos := range positions {
		m.addFinderPattern(pos[0], pos[1])
	}
}

func (m *Matrix) addFinderPattern(x, y int) {
	pattern := [][]bool{
		{true, true, true, true, true, true, true},
		{true, false, false, false, false, false, true},
		{true, false, true, true, true, false, true},
		{true, false, true, true, true, false, true},
		{true, false, true, true, true, false, true},
		{true, false, false, false, false, false, true},
		{true, true, true, true, true, true, true},
	}

	for dy := 0; dy < 7; dy++ {
		for dx := 0; dx < 7; dx++ {
			m.Set(x+dx, y+dy, pattern[dy][dx])
			m.SetReserved(x+dx, y+dy)
		}
	}

	for dy := -1; dy <= 7; dy++ {
		for dx := -1; dx <= 7; dx++ {
			if dx == -1 || dx == 7 || dy == -1 || dy == 7 {
				if x+dx >= 0 && x+dx < m.Size && y+dy >= 0 && y+dy < m.Size {
					m.Set(x+dx, y+dy, false)
					m.SetReserved(x+dx, y+dy)
				}
			}
		}
	}
}

func (m *Matrix) AddTimingPatterns() {
	for i := 8; i < m.Size-8; i++ {
		value := (i % 2) == 0
		m.Set(i, 6, value)
		m.Set(6, i, value)
		m.SetReserved(i, 6)
		m.SetReserved(6, i)
	}
}

func (m *Matrix) AddDarkModule() {
	version := (m.Size - 17) / 4
	x := 8
	y := 4*version + 9
	m.Set(x, y, true)
	m.SetReserved(x, y)
}

func (m *Matrix) AddFormatInfo(level ErrorCorrectionLevel, maskPattern int) {
	formatBits := getFormatBits(level, maskPattern)

	for i := 0; i < 15; i++ {
		bit := (formatBits >> i) & 1
		value := bit == 1

		if i < 6 {
			m.Set(8, i, value)
			m.SetReserved(8, i)
		} else if i < 8 {
			m.Set(8, i+1, value)
			m.SetReserved(8, i+1)
		} else if i == 8 {
			m.Set(7, 8, value)
			m.SetReserved(7, 8)
		} else {
			m.Set(14-i, 8, value)
			m.SetReserved(14-i, 8)
		}
	}

	for i := 0; i < 15; i++ {
		bit := (formatBits >> i) & 1
		value := bit == 1

		if i < 8 {
			m.Set(m.Size-1-i, 8, value)
			m.SetReserved(m.Size-1-i, 8)
		} else {
			m.Set(8, m.Size-15+i, value)
			m.SetReserved(8, m.Size-15+i)
		}
	}
}

func getFormatBits(level ErrorCorrectionLevel, maskPattern int) int {
	formatTable := map[string]int{
		"L0": 0x77C4, "L1": 0x72F3, "L2": 0x7DAA, "L3": 0x789D,
		"L4": 0x662F, "L5": 0x6318, "L6": 0x6C41, "L7": 0x6976,
		"M0": 0x5412, "M1": 0x5125, "M2": 0x5E7C, "M3": 0x5B4B,
		"M4": 0x45F9, "M5": 0x40CE, "M6": 0x4F97, "M7": 0x4AA0,
		"Q0": 0x355F, "Q1": 0x3068, "Q2": 0x3F31, "Q3": 0x3A06,
		"Q4": 0x24B4, "Q5": 0x2183, "Q6": 0x2EDA, "Q7": 0x2BED,
		"H0": 0x1689, "H1": 0x13BE, "H2": 0x1CE7, "H3": 0x19D0,
		"H4": 0x0762, "H5": 0x0255, "H6": 0x0D0C, "H7": 0x083B,
	}

	levelChars := []string{"L", "M", "Q", "H"}
	key := levelChars[level] + string(rune('0'+maskPattern))

	if val, ok := formatTable[key]; ok {
		return val
	}
	return 0
}
