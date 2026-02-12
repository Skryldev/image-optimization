package imageopt

import (
	"log"
	"os"
	"github.com/h2non/bimg"
)
const (
	maxWidthLarge  = 2000
	maxWidthMedium = 1600
	maxWidthSmall  = 1200
)

func calculateTargetWidth(originalWidth int) int {
	if originalWidth <= maxWidthSmall {
		return originalWidth
	}

	type breakpoint struct {
		minWidth int
		target   int
	}

	breakpoints := []breakpoint{
		{minWidth: 3600, target: maxWidthLarge},  // 6K–8K
		{minWidth: 2800, target: 2000},           // 5K
		{minWidth: 2200, target: maxWidthMedium}, // 4K
		{minWidth: 1600, target: maxWidthSmall},  // large web
	}

	for _, bp := range breakpoints {
		if originalWidth >= bp.minWidth {
			return bp.target
		}
	}

	return originalWidth
}

func calculateQuality(fileSize int, pixels int) int {
	const (
		baseQuality = 76
		minQuality  = 45
	)

	megapixels := float64(pixels) / 1_000_000
	if megapixels < 1 {
		return baseQuality
	}

	bytesPerMP := float64(fileSize) / megapixels

	q := baseQuality

	switch {
	case bytesPerMP > 8_000_000: // very dense / noisy
		q -= 28
	case bytesPerMP > 5_000_000:
		q -= 22
	case bytesPerMP > 3_500_000:
		q -= 15
	case bytesPerMP > 2_500_000:
		q -= 8
	}

	// Slight penalty for ultra-high resolution images
	switch {
	case megapixels > 24: // ~8K
		q -= 10
	case megapixels > 16: // ~6K
		q -= 6
	case megapixels > 12: // ~5K
		q -= 4
	}

	return clamp(q, minQuality, baseQuality)
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func OptimizeImageForWeb(inputFile string, outputFile string) error {
	input, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	img := bimg.NewImage(input)

	size, err := img.Size()
	if err != nil {
		return err
	}

	width := size.Width
	height := size.Height
	pixels := width * height
	fileSize := len(input)

	targetWidth := calculateTargetWidth(width)
	quality := calculateQuality(fileSize, pixels)
	resized := targetWidth < width

	options := bimg.Options{
		Width:         targetWidth,
		Quality:      quality,
		Type:          bimg.WEBP,
		StripMetadata: true,
		Interlace:    true,
		Lossless:     false,
	}

	if resized {
		options.Sharpen = bimg.Sharpen{
			Radius: 1,
			X1:     1,
			Y2:     1,
			M1:     1,
			M2:     2,
		}
	}

	output, err := img.Process(options)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputFile, output, 0644); err != nil {
		return err
	}

	log.Printf(
		"Done | %dx%d → %dpx | %dKB → %dKB | Q=%d",
		width,
		height,
		targetWidth,
		fileSize/1024,
		len(output)/1024,
		quality,
	)

	return nil
}