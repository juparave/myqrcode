# MyQRCode - Chrome-Style QR Code Generator

## Project Overview
A Go library for generating QR codes in the exact visual style of Chrome's QR generator (like qrcode_meet.google.com.png), featuring:
- **Rounded square finder patterns** (not just rounded corners)
- **Circular dots for all data modules** 
- **Smart logo embedding** (Chrome dinosaur style)
- **Professional, clean appearance** matching Google's design

Built from scratch to address existing libraries' inability to create readable QR codes with embedded logos while maintaining the distinctive Chrome aesthetic.

## Architecture

### Core Components

#### 1. Main QR Code Structure (`qrcode.go`)
- **QRCode struct**: Main structure containing version, error correction, mode, data, matrix, and logo info
- **ErrorCorrectionLevel**: Low, Medium, Quartile, High (maps to ~7%, ~15%, ~25%, ~30%)
- **EncodingMode**: Numeric, Alphanumeric, Byte
- **StyleConfig**: Rendering configuration (module size, quiet zone, rounded corners, circular dots, colors)

#### 2. Data Encoding (`encoding.go`)
- **Mode Detection**: Automatically detects optimal encoding mode based on input data
- **Numeric Encoding**: Groups of 3 digits → 10 bits, 2 digits → 7 bits, 1 digit → 4 bits
- **Alphanumeric Encoding**: Pairs of chars → 11 bits, single char → 6 bits
- **Byte Encoding**: Each byte → 8 bits
- **Character Count Indicators**: Variable bit length based on version and mode
- **Terminator & Padding**: Adds 0000 terminator and 0xEC/0x11 padding bytes

#### 3. Version Management (`version.go`)
- **Version Table**: Supports versions 1-10 (21x21 to 57x57 modules)
- **Capacity Calculation**: Determines minimum version needed for data + error correction
- **Block Structure**: Data/error correction codewords per block for each version/level

#### 4. Matrix Construction (`matrix.go`)
- **Finder Patterns**: 7x7 patterns in three corners with 1-module quiet border
- **Timing Patterns**: Alternating black/white modules on row 6 and column 6
- **Dark Module**: Single black module at (8, 4*version+9)
- **Format Information**: 15-bit error correction level + mask pattern info
- **Reserved Areas**: Marks areas that cannot contain data

#### 5. Reed-Solomon Error Correction (`reed_solomon.go`)
- **Dependencies**: Uses `rsc.io/qr/gf256` for Galois Field GF(256) arithmetic
- **Field Configuration**: Polynomial 0x11d (285), generator 2 (per ISO 18004)
- **Block Processing**: Splits data into blocks, generates EC codes, interleaves final output
- **Integration**: `generateErrorCorrection()` uses `gf256.NewRSEncoder(field, ecCodewords)`

#### 6. Data Placement (`placement.go`)
- **Zigzag Pattern**: Places data right-to-left, alternating up/down direction
- **Column Skipping**: Skips timing column (column 6)
- **Mask Patterns**: 8 different mask patterns (0-7) with specific formulas
- **Mask Evaluation**: 4-rule penalty system for selecting optimal mask
- **Best Mask Selection**: Tests all 8 masks, selects lowest penalty score

#### 7. Logo Integration (`logo.go`)
- **Smart Placement**: Centers logo while avoiding critical QR areas (finder patterns, timing, format info)
- **Size Calculation**: Logo size as percentage of QR code size
- **Error Correction Adjustment**: Increases EC level based on logo size (up to High/30%)
- **Critical Area Avoidance**: Identifies and avoids finder patterns, timing patterns, format areas
- **Area Reservation**: Reserves logo area during data placement phase

#### 8. Chrome-Style Rendering (`render.go`)
- **Rounded Corners**: Special rendering for finder patterns with rounded corners
- **Circular Dots**: Option to render data modules as circles instead of squares
- **Logo Embedding**: Scales and overlays logo image using bilinear interpolation
- **Style Configuration**: Customizable colors, module size, quiet zone
- **Image Output**: Generates RGBA images compatible with standard Go image packages

## Key Dependencies

### External Libraries
- `rsc.io/qr/gf256`: Galois Field arithmetic and Reed-Solomon encoding
- `golang.org/x/image/draw`: Advanced image scaling and manipulation

### Standard Library
- `image`, `image/color`, `image/draw`: Basic image operations
- `math`: Mathematical calculations for circular rendering
- `errors`: Error handling

## Usage Patterns

### Basic QR Code
```go
qr, _ := myqrcode.New("https://example.com", myqrcode.High)
qr.Encode()
img, _ := qr.ToImage(myqrcode.StyleConfig{
    ModuleSize:     8,
    RoundedCorners: true,
})
```

### QR Code with Logo
```go
qr, _ := myqrcode.New("https://example.com", myqrcode.High)
qr.SetLogo(logoImage, 20) // 20% of QR size
qr.Encode()
img, _ := qr.ToImage(myqrcode.StyleConfig{
    CircularDots: true,
})
```

## Technical Implementation Details

### Reed-Solomon Error Correction
- Uses QR-specific Galois Field configuration (polynomial 0x11d, generator 2)
- Handles block-based error correction as per ISO 18004 specification
- Automatically interleaves data and error correction codewords

### Logo-Aware Generation
- Calculates logo placement to minimize overlap with critical QR areas
- Reserves logo area before data placement to prevent data corruption
- Adjusts error correction level based on logo size for reliability
- Uses high error correction (Level H, ~30%) for larger logos

### Chrome-Style Aesthetics
- Rounded corners on finder patterns match Chrome's QR generator
- Circular dots option for modern appearance
- Configurable styling without compromising QR code functionality
- Maintains QR code readability while enhancing visual appeal

### Performance Considerations
- Supports versions 1-10 (adequate for most use cases)
- Efficient mask evaluation with penalty-based selection
- Minimal memory allocation through reuse of data structures
- Direct matrix manipulation for optimal performance

## Development Status

### ✅ Completed Features
- Basic QR code generation with Reed-Solomon error correction
- Circular dots for data modules
- Logo embedding with smart placement
- Basic rounded corners for finder patterns
- Multiple error correction levels
- Auto-adjusting error correction for logo size

### ✅ QR Code Functionality Status
- **Core QR Generation**: ✅ Working correctly - generates valid, scannable QR codes
- **Reed-Solomon Error Correction**: ✅ Implemented using proper GF(256) arithmetic
- **All QR Standards**: ✅ Finder patterns, timing patterns, format info, data placement all correct
- **Multiple Data Types**: ✅ Numeric, Alphanumeric, and Byte encoding working
- **Error Correction Levels**: ✅ Low, Medium, Quartile, High all functional

### 🔄 Visual Style Improvements Needed (Compared to Chrome's qrcode_meet.google.com.png)
1. **Finder Pattern Aesthetics**: Make corners more smoothly rounded like Chrome's style
2. **Dinosaur Logo**: Replace simple placeholder with accurate Chrome dinosaur silhouette
3. **Proportions & Spacing**: Fine-tune to exactly match Chrome's visual balance

### Testing Commands

#### Readability Testing (Priority 1)
- `go test -v -run TestMinimalQRCode` - Generate simplest readable QR code
- `go test -v -run TestCompareWithReference` - Basic vs Chrome style comparison
- `go test -v -run TestDebugQRGeneration` - Step-by-step QR generation debug
- `go test -v -run TestQRMatrixValidation` - Validate QR structure correctness

#### Visual Style Testing
- `go test -v -run TestTargetReplication` - Generate target comparison
- `go test -v -run TestVisualComparison` - Generate comprehensive visual tests
- `go run example/main.go` - Generate basic Chrome-style QR code
- `go run example/main_with_logo.go` - Generate QR code with embedded logo

#### Test Output Directories
- `readability_tests/` - QR codes optimized for scanning
- `visual_tests/` - Style comparison outputs
- `debug_tests/` - Debug visualization outputs

### Current Development Cycle
1. Run visual tests: `go test -v -run TestTargetReplication`
2. Compare output in `visual_tests/target_replication/attempt_1.png` with `qrcode_meet.google.com.png`
3. Identify specific visual differences
4. Update rendering code in `render.go`
5. Re-run tests and iterate

## File Structure
```
myqrcode/
├── qrcode.go           # Main QRCode struct and public API
├── encoding.go         # Data encoding (numeric, alphanumeric, byte)
├── version.go          # Version management and capacity calculation
├── matrix.go           # QR matrix construction and patterns
├── reed_solomon.go     # Error correction using rsc.io/qr/gf256
├── placement.go        # Data placement and mask selection
├── logo.go            # Logo integration and smart placement
├── render.go          # Chrome-style rendering and image generation
├── example/           # Usage examples and test files
└── CLAUDE.md          # This documentation file
```