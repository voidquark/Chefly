package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	openai "github.com/sashabaranov/go-openai"
)

// OpenAIService handles OpenAI API interactions for image generation
type OpenAIService struct {
	client *openai.Client
	model  string
}

// NewOpenAIService creates a new OpenAI service
func NewOpenAIService(apiKey, model string) *OpenAIService {
	return &OpenAIService{
		client: openai.NewClient(apiKey),
		model:  model,
	}
}

// GenerateFoodImage generates a realistic food image using DALL-E 3
func (s *OpenAIService) GenerateFoodImage(recipeTitle, cuisineType, description string) (string, error) {
	// Build a detailed prompt for realistic food photography
	prompt := s.buildImagePrompt(recipeTitle, cuisineType, description)

	// Call OpenAI API using configured model
	req := openai.ImageRequest{
		Prompt:         prompt,
		Model:          s.model,
		N:              1,
		Size:           openai.CreateImageSize1024x1024,
		Quality:        openai.CreateImageQualityStandard,
		Style:          openai.CreateImageStyleNatural,
		ResponseFormat: openai.CreateImageResponseFormatURL,
	}

	ctx := context.Background()
	resp, err := s.client.CreateImage(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to generate image: %w", err)
	}

	if len(resp.Data) == 0 {
		return "", fmt.Errorf("no image generated")
	}

	// Get the image URL
	imageURL := resp.Data[0].URL

	// Download and convert to base64 data URL for storage
	dataURL, err := s.downloadAndConvertToDataURL(imageURL)
	if err != nil {
		// If download fails, return the URL as fallback
		return imageURL, nil
	}

	return dataURL, nil
}

// buildImagePrompt creates a detailed prompt for food photography
func (s *OpenAIService) buildImagePrompt(recipeTitle, cuisineType, description string) string {
	prompt := fmt.Sprintf(`Professional food photography of %s.
%s cuisine style dish.
High-quality, realistic, appetizing presentation on a clean white plate.
Natural lighting, restaurant quality, top-down view.
Sharp focus on the food with shallow depth of field.
Garnished beautifully with fresh ingredients.
Professional chef presentation, magazine quality photo.
%s
No text, no watermarks, no people, just the delicious food.`,
		recipeTitle,
		cuisineType,
		description)

	return prompt
}

// downloadAndConvertToDataURL downloads an image and converts it to a base64 data URL
func (s *OpenAIService) downloadAndConvertToDataURL(imageURL string) (string, error) {
	// Download the image
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	// Read image data
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image: %w", err)
	}

	// Convert to base64 data URL
	base64Data := base64.StdEncoding.EncodeToString(imageData)

	// Determine content type from response
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/png"
	}

	dataURL := fmt.Sprintf("data:%s;base64,%s", contentType, base64Data)

	return dataURL, nil
}
