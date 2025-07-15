package renderer

import (
	"io"
	"net/http"
	"net/url"
	"os"

	aic "github.com/TheZoraiz/ascii-image-converter/aic_package"
)

// ImageRenderer handles conversion of images to ASCII art
type ImageRenderer struct {
	width  int
	height int
	colored bool
}

// NewImageRenderer creates a new image renderer with default settings
func NewImageRenderer() *ImageRenderer {
	return &ImageRenderer{
		width:  80,
		height: 40,
		colored: true,
	}
}

// SetDimensions sets the ASCII art dimensions
func (ir *ImageRenderer) SetDimensions(width, height int) {
	ir.width = width
	ir.height = height
}

// SetColored enables or disables colored ASCII output
func (ir *ImageRenderer) SetColored(colored bool) {
	ir.colored = colored
}

// RenderImageAsASCII fetches an image from the given src (handling relative URLs) and converts it to ASCII art
func (ir *ImageRenderer) RenderImageAsASCII(src, baseURL string) (string, error) {
	// Handle relative URLs
	u, err := url.Parse(src)
	if err != nil {
		return "", err
	}
	if !u.IsAbs() {
		base, err := url.Parse(baseURL)
		if err != nil {
			return "", err
		}
		src = base.ResolveReference(u).String()
	}

	resp, err := http.Get(src)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	tempFile, err := os.CreateTemp("", "brauser-img-*.tmp")
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write(imgData); err != nil {
		return "", err
	}
	if err := tempFile.Close(); err != nil {
		return "", err
	}

	flags := aic.DefaultFlags()
	flags.Dimensions = []int{ir.width, ir.height}
	flags.Colored = ir.colored

	asciiArt, err := aic.Convert(tempFile.Name(), flags)
	if err != nil {
		return "", err
	}

	return asciiArt, nil
}