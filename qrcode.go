package myqrcode

import (
	"errors"
	"image"
	"image/color"
)

type ErrorCorrectionLevel int

const (
	Low ErrorCorrectionLevel = iota
	Medium
	Quartile
	High
)

type EncodingMode int

const (
	Numeric EncodingMode = iota
	Alphanumeric
	Byte
)

type QRCode struct {
	Version         int
	ErrorCorrection ErrorCorrectionLevel
	Mode            EncodingMode
	Data            string
	Matrix          [][]bool
	Size            int
	Logo            image.Image
	LogoSize        int
}

type StyleConfig struct {
	ModuleSize      int
	QuietZone       int
	RoundedCorners  bool
	CircularDots    bool
	BackgroundColor color.Color
	ForegroundColor color.Color
	ModuleDrawer    ModuleDrawer // New: pluggable module drawing system
}

func New(data string, level ErrorCorrectionLevel) (*QRCode, error) {
	if data == "" {
		return nil, errors.New("data cannot be empty")
	}

	qr := &QRCode{
		Data:            data,
		ErrorCorrection: level,
		LogoSize:        0,
	}

	return qr, nil
}

func (qr *QRCode) SetLogo(logo image.Image, size int) {
	qr.Logo = logo
	qr.LogoSize = size
}

func (qr *QRCode) Encode() error {
	// Detect encoding mode if not set
	if qr.Mode == 0 {
		qr.Mode = detectMode(qr.Data)
	}

	// Determine version based on data length
	if qr.Version == 0 {
		qr.Version = determineVersion(qr.Data, qr.Mode, qr.ErrorCorrection)
	}

	// Adjust error correction level if logo is present
	if qr.Logo != nil && qr.LogoSize > 0 {
		versionInfo := getVersionInfo(qr.Version)
		qr.Size = versionInfo.Size

		placement := optimizeLogoPlacement(&Matrix{Size: qr.Size}, qr.LogoSize)
		qr.ErrorCorrection = adjustErrorCorrectionForLogo(qr.ErrorCorrection, placement, qr.Version)
	}

	// Get version info
	versionInfo := getVersionInfo(qr.Version)
	qr.Size = versionInfo.Size

	// Encode data
	encodedData, err := encodeData(qr.Data, qr.Mode, qr.Version)
	if err != nil {
		return err
	}

	// Add terminator and padding
	encodedData = addTerminatorAndPadding(encodedData, qr.Version, qr.ErrorCorrection)

	// Add error correction
	finalData := addErrorCorrection(encodedData, qr.Version, qr.ErrorCorrection)

	// Create matrix and add patterns
	matrix := NewMatrix(qr.Size)
	matrix.AddFinderPatterns()
	matrix.AddTimingPatterns()
	matrix.AddDarkModule()

	// Reserve logo area if present
	if qr.Logo != nil && qr.LogoSize > 0 {
		placement := optimizeLogoPlacement(matrix, qr.LogoSize)
		reserveLogoArea(matrix, placement)
	}

	// Place data
	placeData(matrix, finalData)

	// Select best mask
	finalMatrix, _ := selectBestMask(matrix, qr.ErrorCorrection)

	// Convert to bool matrix
	qr.Matrix = make([][]bool, qr.Size)
	for i := range qr.Matrix {
		qr.Matrix[i] = make([]bool, qr.Size)
		for j := range qr.Matrix[i] {
			qr.Matrix[i][j] = finalMatrix.Get(j, i)
		}
	}

	return nil
}

// DefaultStyleConfig returns a basic style configuration
func DefaultStyleConfig() StyleConfig {
	return StyleConfig{
		ModuleSize:      8,
		QuietZone:       4,
		RoundedCorners:  false,
		CircularDots:    false,
		BackgroundColor: color.RGBA{255, 255, 255, 255}, // White
		ForegroundColor: color.RGBA{0, 0, 0, 255},       // Black
		ModuleDrawer:    NewSquareModuleDrawer(),
	}
}

// ChromeStyleConfig returns a configuration that mimics Chrome's QR code style
func ChromeStyleConfig() StyleConfig {
	return StyleConfig{
		ModuleSize:      10,
		QuietZone:       4,
		RoundedCorners:  true,
		CircularDots:    true,
		BackgroundColor: color.RGBA{255, 255, 255, 255}, // White
		ForegroundColor: color.RGBA{0, 0, 0, 255},       // Black
		ModuleDrawer:    NewCircleModuleDrawer(),
	}
}

// ChromeFinderPatternStyleConfig returns Chrome style with rounded finder patterns
func ChromeFinderPatternStyleConfig() StyleConfig {
	return StyleConfig{
		ModuleSize:      10,
		QuietZone:       4,
		RoundedCorners:  true,
		CircularDots:    true,
		BackgroundColor: color.RGBA{255, 255, 255, 255}, // White
		ForegroundColor: color.RGBA{0, 0, 0, 255},       // Black
		ModuleDrawer:    NewRoundedModuleDrawer(1.0),    // Full rounding
	}
}

// ChromeGappedStyleConfig returns the most accurate Chrome QR style with gapped circles
func ChromeGappedStyleConfig() StyleConfig {
	return StyleConfig{
		ModuleSize:      10,
		QuietZone:       4,
		RoundedCorners:  false,
		CircularDots:    false,
		BackgroundColor: color.RGBA{255, 255, 255, 255}, // White
		ForegroundColor: color.RGBA{0, 0, 0, 255},       // Black
		ModuleDrawer:    NewGappedCircleModuleDrawer(0.87), // Chrome-like ratio
	}
}

// ChromeGappedStyleConfigWithRatio returns Chrome gapped style with custom ratio
func ChromeGappedStyleConfigWithRatio(ratio float64) StyleConfig {
	return StyleConfig{
		ModuleSize:      10,
		QuietZone:       4,
		RoundedCorners:  false,
		CircularDots:    false,
		BackgroundColor: color.RGBA{255, 255, 255, 255}, // White
		ForegroundColor: color.RGBA{0, 0, 0, 255},       // Black
		ModuleDrawer:    NewGappedCircleModuleDrawer(ratio),
	}
}

// Make is a convenience function for one-line QR code generation
func Make(data string, options ...func(*StyleConfig)) (image.Image, error) {
	qr, err := New(data, High)
	if err != nil {
		return nil, err
	}

	err = qr.Encode()
	if err != nil {
		return nil, err
	}

	config := DefaultStyleConfig()
	for _, opt := range options {
		opt(&config)
	}

	return qr.ToImage(config)
}

// WithModuleDrawer sets a custom module drawer
func WithModuleDrawer(drawer ModuleDrawer) func(*StyleConfig) {
	return func(config *StyleConfig) {
		config.ModuleDrawer = drawer
	}
}

// WithCircles configures circular dots
func WithCircles() func(*StyleConfig) {
	return func(config *StyleConfig) {
		config.CircularDots = true
		config.ModuleDrawer = NewCircleModuleDrawer()
	}
}

// WithGappedCircles configures circular dots with gaps
func WithGappedCircles(sizeRatio float64) func(*StyleConfig) {
	return func(config *StyleConfig) {
		config.CircularDots = true
		config.ModuleDrawer = NewGappedCircleModuleDrawer(sizeRatio)
	}
}

// WithRoundedCorners configures context-aware rounded corners
func WithRoundedCorners(radiusRatio float64) func(*StyleConfig) {
	return func(config *StyleConfig) {
		config.RoundedCorners = true
		config.ModuleDrawer = NewRoundedModuleDrawer(radiusRatio)
	}
}

// WithGappedSquares configures square modules with gaps
func WithGappedSquares(sizeRatio float64) func(*StyleConfig) {
	return func(config *StyleConfig) {
		config.ModuleDrawer = NewGappedSquareModuleDrawer(sizeRatio)
	}
}

// WithColors sets foreground and background colors
func WithColors(fg, bg color.Color) func(*StyleConfig) {
	return func(config *StyleConfig) {
		config.ForegroundColor = fg
		config.BackgroundColor = bg
	}
}

// WithModuleSize sets the module size
func WithModuleSize(size int) func(*StyleConfig) {
	return func(config *StyleConfig) {
		config.ModuleSize = size
	}
}
