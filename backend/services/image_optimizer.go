package services

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

// ImageOptimizer handles image optimization tasks
type ImageOptimizer struct {
	uploadsDir string
}

// OptimizedImages contains paths to optimized image variants
type OptimizedImages struct {
	FullImagePath      string
	ThumbnailPath      string
	FullImageURL       string
	ThumbnailURL       string
}

// NewImageOptimizer creates a new ImageOptimizer instance
func NewImageOptimizer(uploadsDir string) *ImageOptimizer {
	return &ImageOptimizer{
		uploadsDir: uploadsDir,
	}
}

// OptimizeRecipeImage downloads and optimizes a recipe image from DALL-E
// Returns both full-size (800x800) and thumbnail (200x200) WebP images
func (io *ImageOptimizer) OptimizeRecipeImage(imageURL string) (*OptimizedImages, error) {
	// Download the image from DALL-E
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download image: status code %d", resp.StatusCode)
	}

	// Read image data
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	// Decode the image (supports PNG, JPEG, WebP)
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Generate unique filename
	imageID := uuid.New().String()

	// Create directories if they don't exist
	fullDir := filepath.Join(io.uploadsDir, "images", "full")
	thumbDir := filepath.Join(io.uploadsDir, "images", "thumbnails")

	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create full images directory: %w", err)
	}

	if err := os.MkdirAll(thumbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create thumbnails directory: %w", err)
	}

	// Process full-size image (800x800)
	fullImageFilename := fmt.Sprintf("%s.jpg", imageID)
	fullImagePath := filepath.Join(fullDir, fullImageFilename)

	if err := io.createOptimizedImage(img, fullImagePath, 800, 85); err != nil {
		return nil, fmt.Errorf("failed to create full-size image: %w", err)
	}

	// Process thumbnail (200x200)
	thumbFilename := fmt.Sprintf("%s_thumb.jpg", imageID)
	thumbPath := filepath.Join(thumbDir, thumbFilename)

	if err := io.createOptimizedImage(img, thumbPath, 200, 75); err != nil {
		return nil, fmt.Errorf("failed to create thumbnail: %w", err)
	}

	// Generate URLs for serving images
	fullImageURL := fmt.Sprintf("/uploads/images/full/%s", fullImageFilename)
	thumbnailURL := fmt.Sprintf("/uploads/images/thumbnails/%s", thumbFilename)

	return &OptimizedImages{
		FullImagePath:      fullImagePath,
		ThumbnailPath:      thumbPath,
		FullImageURL:       fullImageURL,
		ThumbnailURL:       thumbnailURL,
	}, nil
}

// createOptimizedImage resizes and saves an image as JPEG
func (io *ImageOptimizer) createOptimizedImage(img image.Image, outputPath string, size int, quality int) error {
	// Resize image maintaining aspect ratio (will be square for recipe images)
	resized := imaging.Fill(img, size, size, imaging.Center, imaging.Lanczos)

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Encode as JPEG with specified quality
	// JPEG quality 85 for full images (~100-150KB)
	// JPEG quality 75 for thumbnails (~15-20KB)
	options := &jpeg.Options{Quality: quality}
	err = jpeg.Encode(outFile, resized, options)
	if err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}

	return nil
}
