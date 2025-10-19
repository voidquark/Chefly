package services

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
// Returns both full-size (800x800) and thumbnail (200x200) JPEG images
// Supports both HTTP URLs and base64 data URLs
func (opt *ImageOptimizer) OptimizeRecipeImage(imageURL string) (*OptimizedImages, error) {
	var imageData []byte
	var err error

	// Check if it's a data URL (base64 encoded)
	if strings.HasPrefix(imageURL, "data:image/") {
		// Extract base64 data from data URL
		// Format: data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA...
		parts := strings.Split(imageURL, ",")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid data URL format")
		}

		imageData, err = base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64 image: %w", err)
		}
	} else {
		// Regular HTTP URL - download the image
		resp, err := http.Get(imageURL)
		if err != nil {
			return nil, fmt.Errorf("failed to download image: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to download image: status code %d", resp.StatusCode)
		}

		imageData, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read image data: %w", err)
		}
	}

	// Decode the image (supports PNG, JPEG, WebP)
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Generate unique filename
	imageID := uuid.New().String()

	// Create directories if they don't exist
	fullDir := filepath.Join(opt.uploadsDir, "images", "full")
	thumbDir := filepath.Join(opt.uploadsDir, "images", "thumbnails")

	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create full images directory: %w", err)
	}

	if err := os.MkdirAll(thumbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create thumbnails directory: %w", err)
	}

	// Process full-size image (800x800)
	fullImageFilename := fmt.Sprintf("%s.jpg", imageID)
	fullImagePath := filepath.Join(fullDir, fullImageFilename)

	if err := opt.createOptimizedImage(img, fullImagePath, 800, 85); err != nil {
		return nil, fmt.Errorf("failed to create full-size image: %w", err)
	}

	// Process thumbnail (200x200)
	thumbFilename := fmt.Sprintf("%s_thumb.jpg", imageID)
	thumbPath := filepath.Join(thumbDir, thumbFilename)

	if err := opt.createOptimizedImage(img, thumbPath, 200, 75); err != nil {
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
func (opt *ImageOptimizer) createOptimizedImage(img image.Image, outputPath string, size int, quality int) error {
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
