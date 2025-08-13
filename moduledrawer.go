package myqrcode

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

const AntialiasingFactor = 4

// ModuleDrawer interface defines how QR code modules are rendered
type ModuleDrawer interface {
	// Initialize sets up the drawer with the image and style configuration
	Initialize(img *image.RGBA, config StyleConfig)

	// DrawModule renders a single module at the given box coordinates
	// box is [x1, y1, x2, y2] coordinates
	// isActive indicates if this module should be drawn (true) or left as background (false)
	// neighbors provides context about surrounding modules (nil if not needed)
	DrawModule(box [4]int, isActive bool, neighbors *ActiveWithNeighbors)

	// NeedsNeighbors returns true if this drawer requires neighbor context
	NeedsNeighbors() bool
}

// ActiveWithNeighbors provides context about the 8 surrounding modules
type ActiveWithNeighbors struct {
	NW bool // Northwest
	N  bool // North
	NE bool // Northeast
	W  bool // West
	Me bool // Center (this module)
	E  bool // East
	SW bool // Southwest
	S  bool // South
	SE bool // Southeast
}

// BaseModuleDrawer provides common functionality for all module drawers
type BaseModuleDrawer struct {
	img    *image.RGBA
	config StyleConfig
}

func (b *BaseModuleDrawer) Initialize(img *image.RGBA, config StyleConfig) {
	b.img = img
	b.config = config
}

func (b *BaseModuleDrawer) NeedsNeighbors() bool {
	return false
}

// createAntialiasingImage creates a larger image for anti-aliasing
func createAntialiasingImage(size int, bgColor color.Color) *image.RGBA {
	bigSize := size * AntialiasingFactor
	img := image.NewRGBA(image.Rect(0, 0, bigSize, bigSize))

	// Fill with background color
	c := color.RGBAModel.Convert(bgColor).(color.RGBA)
	for y := 0; y < bigSize; y++ {
		for x := 0; x < bigSize; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

// resizeImage downscales an image with anti-aliasing
func resizeImage(src *image.RGBA, targetSize int) *image.RGBA {
	srcBounds := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, targetSize, targetSize))

	// Simple bilinear downsampling
	scale := float64(srcBounds.Dx()) / float64(targetSize)

	for y := 0; y < targetSize; y++ {
		for x := 0; x < targetSize; x++ {
			srcX := float64(x) * scale
			srcY := float64(y) * scale

			// Get the four surrounding pixels
			x1, y1 := int(srcX), int(srcY)
			x2, y2 := x1+1, y1+1

			// Clamp to bounds
			if x2 >= srcBounds.Dx() {
				x2 = srcBounds.Dx() - 1
			}
			if y2 >= srcBounds.Dy() {
				y2 = srcBounds.Dy() - 1
			}

			// Get colors
			c1 := src.RGBAAt(x1, y1)
			c2 := src.RGBAAt(x2, y1)
			c3 := src.RGBAAt(x1, y2)
			c4 := src.RGBAAt(x2, y2)

			// Interpolation weights
			fx := srcX - float64(x1)
			fy := srcY - float64(y1)

			// Bilinear interpolation
			r := uint8(float64(c1.R)*(1-fx)*(1-fy) + float64(c2.R)*fx*(1-fy) +
				float64(c3.R)*(1-fx)*fy + float64(c4.R)*fx*fy)
			g := uint8(float64(c1.G)*(1-fx)*(1-fy) + float64(c2.G)*fx*(1-fy) +
				float64(c3.G)*(1-fx)*fy + float64(c4.G)*fx*fy)
			b := uint8(float64(c1.B)*(1-fx)*(1-fy) + float64(c2.B)*fx*(1-fy) +
				float64(c3.B)*(1-fx)*fy + float64(c4.B)*fx*fy)
			a := uint8(float64(c1.A)*(1-fx)*(1-fy) + float64(c2.A)*fx*(1-fy) +
				float64(c3.A)*(1-fx)*fy + float64(c4.A)*fx*fy)

			dst.Set(x, y, color.RGBA{r, g, b, a})
		}
	}

	return dst
}

// resizeImageHighQuality provides higher quality resizing similar to Python's Lanczos
func resizeImageHighQuality(src *image.RGBA, targetSize int) *image.RGBA {
	srcBounds := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, targetSize, targetSize))
	
	// Use a higher quality interpolation approach
	scaleX := float64(srcBounds.Dx()) / float64(targetSize)
	scaleY := float64(srcBounds.Dy()) / float64(targetSize)
	
	for y := 0; y < targetSize; y++ {
		for x := 0; x < targetSize; x++ {
			// Use more sophisticated sampling for better quality
			srcX := float64(x) * scaleX
			srcY := float64(y) * scaleY
			
			// Sample multiple points for better anti-aliasing (mimics Lanczos approach)
			var totalR, totalG, totalB, totalA float64
			var samples int
			
			for dy := -0.5; dy <= 0.5; dy += 0.25 {
				for dx := -0.5; dx <= 0.5; dx += 0.25 {
					sampleX := int(srcX + dx)
					sampleY := int(srcY + dy)
					
					if sampleX >= 0 && sampleX < srcBounds.Dx() && sampleY >= 0 && sampleY < srcBounds.Dy() {
						c := src.RGBAAt(sampleX, sampleY)
						totalR += float64(c.R)
						totalG += float64(c.G)
						totalB += float64(c.B)
						totalA += float64(c.A)
						samples++
					}
				}
			}
			
			if samples > 0 {
				dst.Set(x, y, color.RGBA{
					uint8(totalR / float64(samples)),
					uint8(totalG / float64(samples)),
					uint8(totalB / float64(samples)),
					uint8(totalA / float64(samples)),
				})
			}
		}
	}
	
	return dst
}

// SquareModuleDrawer draws basic square modules
type SquareModuleDrawer struct {
	BaseModuleDrawer
}

func NewSquareModuleDrawer() *SquareModuleDrawer {
	return &SquareModuleDrawer{}
}

func (s *SquareModuleDrawer) DrawModule(box [4]int, isActive bool, neighbors *ActiveWithNeighbors) {
	if !isActive {
		return
	}

	rect := image.Rect(box[0], box[1], box[2], box[3])
	draw.Draw(s.img, rect, &image.Uniform{s.config.ForegroundColor}, image.Point{}, draw.Src)
}

// CircleModuleDrawer draws circular modules with anti-aliasing
type CircleModuleDrawer struct {
	BaseModuleDrawer
	circle *image.RGBA
}

func NewCircleModuleDrawer() *CircleModuleDrawer {
	return &CircleModuleDrawer{}
}

func (c *CircleModuleDrawer) Initialize(img *image.RGBA, config StyleConfig) {
	c.BaseModuleDrawer.Initialize(img, config)
	c.createCircle()
}

func (c *CircleModuleDrawer) createCircle() {
	size := c.config.ModuleSize
	bigImg := createAntialiasingImage(size, c.config.BackgroundColor)
	bigSize := size * AntialiasingFactor
	center := float64(bigSize) / 2
	radius := center

	fgColor := color.RGBAModel.Convert(c.config.ForegroundColor).(color.RGBA)

	// Draw anti-aliased circle
	for y := 0; y < bigSize; y++ {
		for x := 0; x < bigSize; x++ {
			dx := float64(x) - center + 0.5
			dy := float64(y) - center + 0.5
			distance := math.Sqrt(dx*dx + dy*dy)

			if distance <= radius {
				bigImg.Set(x, y, fgColor)
			}
		}
	}

	c.circle = resizeImage(bigImg, size)
}

func (c *CircleModuleDrawer) DrawModule(box [4]int, isActive bool, neighbors *ActiveWithNeighbors) {
	if !isActive {
		return
	}

	// Paste the pre-rendered circle
	dst := image.Rect(box[0], box[1], box[2], box[3])
	src := c.circle.Bounds()
	draw.Draw(c.img, dst, c.circle, src.Min, draw.Src)
}

// GappedSquareModuleDrawer draws squares with configurable gaps
type GappedSquareModuleDrawer struct {
	BaseModuleDrawer
	SizeRatio float64
}

func NewGappedSquareModuleDrawer(sizeRatio float64) *GappedSquareModuleDrawer {
	if sizeRatio <= 0 || sizeRatio > 1 {
		sizeRatio = 0.8
	}
	return &GappedSquareModuleDrawer{SizeRatio: sizeRatio}
}

func (g *GappedSquareModuleDrawer) DrawModule(box [4]int, isActive bool, neighbors *ActiveWithNeighbors) {
	if !isActive {
		return
	}

	// Calculate smaller box based on size ratio
	width := box[2] - box[0]

	delta := int(float64(width) * (1.0 - g.SizeRatio) / 2.0)

	smallerBox := image.Rect(
		box[0]+delta,
		box[1]+delta,
		box[2]-delta,
		box[3]-delta,
	)

	draw.Draw(g.img, smallerBox, &image.Uniform{g.config.ForegroundColor}, image.Point{}, draw.Src)
}

// GappedCircleModuleDrawer draws circles with configurable gaps
type GappedCircleModuleDrawer struct {
	BaseModuleDrawer
	SizeRatio float64
	circle    *image.RGBA
}

func NewGappedCircleModuleDrawer(sizeRatio float64) *GappedCircleModuleDrawer {
	if sizeRatio <= 0 || sizeRatio > 1 {
		sizeRatio = 0.9
	}
	return &GappedCircleModuleDrawer{SizeRatio: sizeRatio}
}

func (g *GappedCircleModuleDrawer) Initialize(img *image.RGBA, config StyleConfig) {
	g.BaseModuleDrawer.Initialize(img, config)
	g.createGappedCircle()
}

func (g *GappedCircleModuleDrawer) createGappedCircle() {
	size := g.config.ModuleSize
	
	// Step 1: Create full-size anti-aliased circle (like Python's approach)
	bigImg := createAntialiasingImage(size, g.config.BackgroundColor)
	bigSize := size * AntialiasingFactor
	center := float64(bigSize) / 2
	radius := center

	fgColor := color.RGBAModel.Convert(g.config.ForegroundColor).(color.RGBA)

	// Draw full anti-aliased circle
	for y := 0; y < bigSize; y++ {
		for x := 0; x < bigSize; x++ {
			dx := float64(x) - center + 0.5
			dy := float64(y) - center + 0.5
			distance := math.Sqrt(dx*dx + dy*dy)

			if distance <= radius {
				bigImg.Set(x, y, fgColor)
			}
		}
	}

	// Step 2: Resize to full module size first (preserves anti-aliasing)
	fullSizeCircle := resizeImage(bigImg, size)
	
	// Step 3: Resize again to the gapped size (like Python's single resize)
	actualSize := int(float64(size) * g.SizeRatio)
	g.circle = resizeImageHighQuality(fullSizeCircle, actualSize)
}

func (g *GappedCircleModuleDrawer) DrawModule(box [4]int, isActive bool, neighbors *ActiveWithNeighbors) {
	if !isActive {
		return
	}

	// Use Python's approach: simple top-left positioning (like paste)
	actualSize := g.circle.Bounds().Dx()
	
	// Center the circle within the module box
	width := box[2] - box[0]
	offset := (width - actualSize) / 2

	dst := image.Rect(
		box[0]+offset,
		box[1]+offset,
		box[0]+offset+actualSize,
		box[1]+offset+actualSize,
	)

	draw.Draw(g.img, dst, g.circle, g.circle.Bounds().Min, draw.Src)
}

// RoundedModuleDrawer draws modules with context-aware rounded corners
// This is the key for Chrome-style finder patterns
type RoundedModuleDrawer struct {
	BaseModuleDrawer
	RadiusRatio float64
	cornerWidth int
	square      *image.RGBA
	nwRound     *image.RGBA
	neRound     *image.RGBA
	seRound     *image.RGBA
	swRound     *image.RGBA
}

func NewRoundedModuleDrawer(radiusRatio float64) *RoundedModuleDrawer {
	if radiusRatio <= 0 || radiusRatio > 1 {
		radiusRatio = 1.0
	}
	return &RoundedModuleDrawer{RadiusRatio: radiusRatio}
}

func (r *RoundedModuleDrawer) NeedsNeighbors() bool {
	return true
}

func (r *RoundedModuleDrawer) Initialize(img *image.RGBA, config StyleConfig) {
	r.BaseModuleDrawer.Initialize(img, config)
	r.cornerWidth = config.ModuleSize / 2
	r.setupCorners()
}

func (r *RoundedModuleDrawer) setupCorners() {
	bgColor := r.config.BackgroundColor
	fgColor := r.config.ForegroundColor

	// Create square corner (no rounding)
	r.square = image.NewRGBA(image.Rect(0, 0, r.cornerWidth, r.cornerWidth))
	fgRGBA := color.RGBAModel.Convert(fgColor).(color.RGBA)
	for y := 0; y < r.cornerWidth; y++ {
		for x := 0; x < r.cornerWidth; x++ {
			r.square.Set(x, y, fgRGBA)
		}
	}

	// Create rounded corners with anti-aliasing
	fakeWidth := r.cornerWidth * AntialiasingFactor
	radius := r.RadiusRatio * float64(fakeWidth)

	// Create base image for northwest rounded corner
	base := createAntialiasingImage(r.cornerWidth, bgColor)
	bigSize := r.cornerWidth * AntialiasingFactor

	// Draw rounded corner: circle in top-left, rectangles extending right and down
	for y := 0; y < bigSize; y++ {
		for x := 0; x < bigSize; x++ {
			// Check if pixel is inside the rounded corner shape
			dx := float64(x)
			dy := float64(y)

			// Always fill if we're in the rectangle extensions
			inRightRect := dx >= radius && x < bigSize
			inBottomRect := dy >= radius && y < bigSize

			// Check if we're in the circular part
			inCircle := false
			if dx < radius && dy < radius {
				distFromCenter := math.Sqrt(dx*dx + dy*dy)
				inCircle = distFromCenter <= radius
			}

			if inCircle || inRightRect || inBottomRect {
				base.Set(x, y, fgRGBA)
			}
		}
	}

	// Resize to actual size
	r.nwRound = resizeImage(base, r.cornerWidth)

	// Create other corners by rotating/flipping
	r.neRound = r.flipHorizontal(r.nwRound)
	r.seRound = r.rotate180(r.nwRound)
	r.swRound = r.flipVertical(r.nwRound)
}

func (r *RoundedModuleDrawer) flipHorizontal(src *image.RGBA) *image.RGBA {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			srcX := bounds.Max.X - 1 - x
			dst.Set(x, y, src.At(srcX, y))
		}
	}
	return dst
}

func (r *RoundedModuleDrawer) flipVertical(src *image.RGBA) *image.RGBA {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			srcY := bounds.Max.Y - 1 - y
			dst.Set(x, y, src.At(x, srcY))
		}
	}
	return dst
}

func (r *RoundedModuleDrawer) rotate180(src *image.RGBA) *image.RGBA {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			srcX := bounds.Max.X - 1 - x
			srcY := bounds.Max.Y - 1 - y
			dst.Set(x, y, src.At(srcX, srcY))
		}
	}
	return dst
}

func (r *RoundedModuleDrawer) DrawModule(box [4]int, isActive bool, neighbors *ActiveWithNeighbors) {
	if !isActive || neighbors == nil {
		return
	}

	// Determine which corners should be rounded based on neighbors
	nwRounded := !neighbors.W && !neighbors.N
	neRounded := !neighbors.N && !neighbors.E
	seRounded := !neighbors.E && !neighbors.S
	swRounded := !neighbors.S && !neighbors.W

	// Select appropriate corner images
	nw := r.square
	if nwRounded {
		nw = r.nwRound
	}

	ne := r.square
	if neRounded {
		ne = r.neRound
	}

	se := r.square
	if seRounded {
		se = r.seRound
	}

	sw := r.square
	if swRounded {
		sw = r.swRound
	}

	// Draw the four corners
	// Northwest corner
	nwDst := image.Rect(box[0], box[1], box[0]+r.cornerWidth, box[1]+r.cornerWidth)
	draw.Draw(r.img, nwDst, nw, nw.Bounds().Min, draw.Src)

	// Northeast corner
	neDst := image.Rect(box[0]+r.cornerWidth, box[1], box[2], box[1]+r.cornerWidth)
	draw.Draw(r.img, neDst, ne, ne.Bounds().Min, draw.Src)

	// Southeast corner
	seDst := image.Rect(box[0]+r.cornerWidth, box[1]+r.cornerWidth, box[2], box[3])
	draw.Draw(r.img, seDst, se, se.Bounds().Min, draw.Src)

	// Southwest corner
	swDst := image.Rect(box[0], box[1]+r.cornerWidth, box[0]+r.cornerWidth, box[3])
	draw.Draw(r.img, swDst, sw, sw.Bounds().Min, draw.Src)
}
