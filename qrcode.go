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
	Version            int
	ErrorCorrection    ErrorCorrectionLevel
	Mode              EncodingMode
	Data              string
	Matrix            [][]bool
	Size              int
	Logo              image.Image
	LogoSize          int
}

type StyleConfig struct {
	ModuleSize       int
	QuietZone        int
	RoundedCorners   bool
	CircularDots     bool
	BackgroundColor  color.Color
	ForegroundColor  color.Color
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
		qr.ErrorCorrection = adjustErrorCorrectionForLogo(qr.ErrorCorrection, placement, len(qr.Data))
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

