package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"embed"

	// "github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	// "golang.org/x/image/font"
	// "golang.org/x/image/math/fixed"
)

//go:embed assets/Roboto-Regular.ttf
var fontData embed.FS

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) UploadAndWatermark(files []string, watermarkText string, outputDir string) ([]string, error) {
	var processedImages []string

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return nil, err
	}

	for _, file := range files {
		img, format, err := loadImage(file)
		if err != nil {
			return nil, fmt.Errorf("error loading image %s: %v", file, err)
		}

		// Add watermark text
		watermarkedImg := addTextWatermark(img, watermarkText)

		// Determine output file path
		outputPath := filepath.Join(outputDir, "watermarked_"+filepath.Base(file))

		// Save the watermarked image
		err = saveImage(outputPath, watermarkedImg, format)
		if err != nil {
			return nil, fmt.Errorf("error saving image %s: %v", outputPath, err)
		}

		processedImages = append(processedImages, outputPath)
	}

	return processedImages, nil
}

// loadImage loads an image from a file
func loadImage(filePath string) (image.Image, string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	img, format, err := image.Decode(f)
	return img, format, err
}

// addTextWatermark adds a simple watermark to the image
func addTextWatermark(img image.Image, text string) image.Image {
	// fmt.Println("[ADD WATERMARK]")
	// // Create a semi-transparent watermark image
	// watermark := imaging.New(200, 50, image.Transparent)
	// watermark = imaging.Fill(watermark, 200, 50, imaging.Center, imaging.Lanczos)

	// // Draw the watermark onto the image
	// offset := image.Pt(0, 0) // Adjust position
	// bounds := img.Bounds()
	// fmt.Println("[IMG BOUNDS]", bounds)
	// result := image.NewRGBA(bounds)
	// draw.Draw(result, bounds, img, image.Point{}, draw.Src)
	// draw.Draw(result, watermark.Bounds().Add(offset), watermark, image.Point{}, draw.Over)

	// Load font
	fontBytes, err := fontData.ReadFile("assets/Roboto-Regular.ttf") // Ensure this font file exists
	if err != nil {
		fmt.Println("Error loading font:", err)
		return img
	}

	f, err := truetype.Parse(fontBytes)
	if err != nil {
		fmt.Println("Error parsing font:", err)
		return img
	}

	// return result
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)

	// Draw original image onto new RGBA image
	draw.Draw(rgba, bounds, img, image.Point{}, draw.Src)

	// Set up font drawer
	c := freetype.NewContext()
	c.SetDPI(300)
	c.SetFont(f)
	c.SetFontSize(48) // Adjust watermark text size
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(image.NewUniform(color.RGBA{255, 255, 255, 80})) // White semi-transparent

	// Define text placement
	pt := freetype.Pt(bounds.Dx()/10, bounds.Dy()/2) // Adjust positioning
	_, err = c.DrawString(text, pt)
	if err != nil {
		fmt.Println("Error drawing string:", err)
	}

	return rgba
}

// saveImage saves the watermarked image to a file
func saveImage(outputPath string, img image.Image, format string) error {
	fmt.Println("[SAVE IMAGE]")
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if strings.ToLower(format) == "jpeg" || strings.ToLower(format) == "jpg" {
		return jpeg.Encode(outFile, img, &jpeg.Options{Quality: 90})
	} else if strings.ToLower(format) == "png" {
		return png.Encode(outFile, img)
	}
	return fmt.Errorf("unsupported format: %s", format)
}

func (a *App) GetDefaultOutputPath() (string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    return filepath.Join(homeDir, "WatermarkAppOutput"), nil
}

func (a *App) SelectFile() ([]string, error) {
	files, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return nil, err
	}
	return files, nil 
}
