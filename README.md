# MyQRCode - Chrome-Style QR Code Generator

A Go library for generating QR codes in the exact visual style of Chrome's QR generator, featuring rounded finder patterns, circular dots, and smart logo embedding with proper error correction.

![Chrome Style QR Code](visual_tests/target_replication/attempt_1.png)

## Features

- 🎨 **Chrome-Style Rendering** - Rounded finder patterns and circular dots matching Google's design
- 🦕 **Smart Logo Embedding** - Embeds logos while maintaining QR code readability
- 🔧 **Built from Scratch** - Custom implementation addressing limitations of existing libraries
- 📱 **Fully Readable** - Generates valid QR codes that scan properly on all devices
- ⚡ **High Performance** - Efficient Reed-Solomon error correction using `rsc.io/qr/gf256`
- 🎯 **Multiple Formats** - Supports Numeric, Alphanumeric, and Byte encoding modes

## Quick Start

### Installation

```bash
go get github.com/juparave/myqrcode
```

### Basic Usage

```go
package main

import (
    "image/png"
    "os"
    "github.com/juparave/myqrcode"
)

func main() {
    // Create a QR code
    qr, err := myqrcode.New("https://meet.google.com/abc-defg-hij", myqrcode.High)
    if err != nil {
        panic(err)
    }

    // Encode the data
    err = qr.Encode()
    if err != nil {
        panic(err)
    }

    // Configure Chrome-style rendering
    config := myqrcode.StyleConfig{
        ModuleSize:     8,
        QuietZone:      32,
        RoundedCorners: true,
        CircularDots:   true,
    }

    // Generate image
    img, err := qr.ToImage(config)
    if err != nil {
        panic(err)
    }

    // Save to file
    file, _ := os.Create("qrcode.png")
    defer file.Close()
    png.Encode(file, img)
}
```

### With Logo

```go
// Load your logo image
logoFile, _ := os.Open("logo.png")
logo, _, _ := image.Decode(logoFile)

// Create QR code with logo
qr, _ := myqrcode.New("https://example.com", myqrcode.High)
qr.SetLogo(logo, 20) // 20% of QR code size

qr.Encode()
img, _ := qr.ToImage(myqrcode.StyleConfig{
    CircularDots: true,
    RoundedCorners: true,
})
```

## API Reference

### Creating QR Codes

```go
// Create new QR code
qr, err := myqrcode.New(data string, level ErrorCorrectionLevel) (*QRCode, error)

// Error correction levels
myqrcode.Low      // ~7% correction
myqrcode.Medium   // ~15% correction  
myqrcode.Quartile // ~25% correction
myqrcode.High     // ~30% correction
```

### Adding Logos

```go
// Set logo (automatically adjusts error correction as needed)
qr.SetLogo(logo image.Image, sizePercent int)
```

### Styling Options

```go
type StyleConfig struct {
    ModuleSize      int         // Size of each QR module in pixels
    QuietZone       int         // Border size around QR code
    RoundedCorners  bool        // Rounded finder patterns (Chrome style)
    CircularDots    bool        // Circular data modules
    BackgroundColor color.Color // Background color
    ForegroundColor color.Color // QR code color
}
```

## Examples

### Different Styles

```go
// Classic square style
config := myqrcode.StyleConfig{
    ModuleSize: 10,
    QuietZone:  40,
}

// Chrome style (recommended)
config := myqrcode.StyleConfig{
    ModuleSize:     10,
    QuietZone:      40,
    RoundedCorners: true,
    CircularDots:   true,
}

// Custom colors
config := myqrcode.StyleConfig{
    ModuleSize:      8,
    QuietZone:       32,
    CircularDots:    true,
    BackgroundColor: color.RGBA{240, 240, 240, 255},
    ForegroundColor: color.RGBA{50, 50, 200, 255},
}
```

### Data Types

The library automatically detects the optimal encoding mode:

```go
// Numeric (most efficient for numbers)
qr, _ := myqrcode.New("1234567890", myqrcode.Medium)

// Alphanumeric (for text with limited character set)
qr, _ := myqrcode.New("HELLO WORLD 123", myqrcode.Medium)

// Byte mode (for any text, URLs, etc.)
qr, _ := myqrcode.New("https://example.com/path?param=value", myqrcode.High)
```

## Testing

The library includes comprehensive tests for validation:

```bash
# Test basic functionality
go test -v -run TestBasicQRGeneration

# Test QR code readability  
go test -v -run TestMinimalQRCode
go test -v -run TestCompareWithReference

# Generate visual comparisons
go test -v -run TestVisualComparison
go test -v -run TestTargetReplication

# Debug QR generation step-by-step
go test -v -run TestDebugQRGeneration
```

Test outputs are saved in:
- `readability_tests/` - QR codes optimized for scanning
- `visual_tests/` - Style comparison outputs  
- `debug_tests/` - Debug visualization outputs

## Why This Library?

Existing Go QR libraries have limitations when embedding logos:

- **skip2/go-qrcode** - Cannot embed logos while maintaining readability
- **yeqown/go-qrcode** - Logo placement often corrupts critical QR areas

**MyQRCode solves this by:**

1. **Smart Logo Placement** - Avoids finder patterns, timing patterns, and format information areas
2. **Automatic Error Correction** - Increases error correction level based on logo size
3. **Chrome-Style Rendering** - Matches Google's modern QR code aesthetic
4. **Built-in Validation** - Comprehensive testing ensures readability

## Technical Details

### Architecture

- **Reed-Solomon Error Correction** using `rsc.io/qr/gf256` with proper GF(256) arithmetic
- **QR Standard Compliance** - Implements ISO/IEC 18004 specification
- **Optimized Mask Selection** - Tests all 8 mask patterns for best readability
- **Logo-Aware Generation** - Reserves logo area during data placement phase

### Supported Features

- **Versions**: 1-10 (21×21 to 57×57 modules)
- **Error Correction**: All levels (L, M, Q, H)
- **Encoding Modes**: Numeric, Alphanumeric, Byte
- **Logo Sizes**: Up to 30% of QR code area (with High error correction)

## Examples Directory

Check out the `example/` directory for complete working examples:

- `main.go` - Basic Chrome-style QR code
- `main_with_logo.go` - QR code with embedded logo

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass: `go test ./...`
5. Submit a pull request

## License

MIT License - see LICENSE file for details.

## Development Status

- ✅ **Core QR Generation** - Fully functional and tested
- ✅ **Logo Embedding** - Smart placement with error correction
- ✅ **Chrome-Style Rendering** - Rounded corners and circular dots
- 🔄 **Visual Polish** - Fine-tuning to exactly match Chrome's appearance

---

Built with ❤️ to create beautiful, scannable QR codes with logos.