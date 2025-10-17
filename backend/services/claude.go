package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"chefly/models"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// Custom error types for better error handling
var (
	ErrAPIConnection  = errors.New("AI service connection error")
	ErrRateLimit      = errors.New("AI service rate limit reached")
	ErrEmptyResponse  = errors.New("AI service returned empty response")
	ErrParsingFailed  = errors.New("failed to parse AI response")
	ErrInvalidJSON    = errors.New("AI response did not contain valid JSON")
)

// ClaudeService handles Claude API interactions
type ClaudeService struct {
	client *anthropic.Client
	model  string
}

// NewClaudeService creates a new Claude service
func NewClaudeService(apiKey, model string) *ClaudeService {
	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)
	return &ClaudeService{
		client: &client,
		model:  model,
	}
}

// GenerateRecipe generates a recipe using Claude AI
func (s *ClaudeService) GenerateRecipe(req models.RecipeGenerationRequest) (*models.RecipeDetail, error) {
	// Build the prompt based on filters
	prompt := s.buildRecipePrompt(req)

	// Call Claude API using configured model
	message, err := s.client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.Model(s.model),
		MaxTokens: 4096,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})

	if err != nil {
		// Check if it's a rate limit error
		errStr := err.Error()
		if strings.Contains(errStr, "429") || strings.Contains(strings.ToLower(errStr), "rate limit") {
			return nil, fmt.Errorf("%w: please wait a moment before generating another recipe", ErrRateLimit)
		}
		// Check for invalid model error
		if strings.Contains(strings.ToLower(errStr), "model") || strings.Contains(errStr, "404") {
			return nil, fmt.Errorf("%w: model '%s' not found or not accessible. API error: %v", ErrAPIConnection, s.model, err)
		}
		// Generic API connection error
		return nil, fmt.Errorf("%w: %v", ErrAPIConnection, err)
	}

	// Extract text from response
	var responseText string
	for _, block := range message.Content {
		if block.Type == "text" {
			responseText += block.Text
		}
	}

	if responseText == "" {
		return nil, ErrEmptyResponse
	}

	// Parse the JSON response
	recipe, err := s.parseRecipeResponse(responseText, req)
	if err != nil {
		// Check if it's a JSON parsing error
		if errors.Is(err, ErrInvalidJSON) {
			return nil, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
		}
		return nil, fmt.Errorf("%w: %v", ErrParsingFailed, err)
	}

	return recipe, nil
}

// buildRecipePrompt creates a detailed prompt for Claude
func (s *ClaudeService) buildRecipePrompt(req models.RecipeGenerationRequest) string {
	var builder strings.Builder

	builder.WriteString("You are a professional chef and recipe creator. Generate a detailed, high-quality recipe in JSON format.\n\n")

	// Language-specific instructions
	if req.Language == "sk" {
		builder.WriteString("IMPORTANT: Generate this recipe IN SLOVAK LANGUAGE (Slovenčina).\n")
		builder.WriteString("All text must be in Slovak:\n")
		builder.WriteString("- Recipe title in Slovak\n")
		builder.WriteString("- Description in Slovak\n")
		builder.WriteString("- Ingredient names in Slovak\n")
		builder.WriteString("- Instructions in Slovak\n")
		builder.WriteString("- Tips in Slovak\n\n")
	}

	builder.WriteString("Requirements:\n")

	// Meat type
	if req.MeatType != "" && req.MeatType != "None (Vegetarian)" {
		builder.WriteString(fmt.Sprintf("- Main protein: %s\n", req.MeatType))
	} else {
		builder.WriteString("- Vegetarian recipe (no meat)\n")
	}

	// Cuisine
	if req.CuisineType != "" {
		builder.WriteString(fmt.Sprintf("- Cuisine style: %s\n", req.CuisineType))
	}

	// Side ingredients
	if len(req.SideIngredients) > 0 {
		builder.WriteString(fmt.Sprintf("- Include these ingredients: %s\n", strings.Join(req.SideIngredients, ", ")))
	}

	// Dietary preferences
	if len(req.DietaryPreferences) > 0 {
		builder.WriteString(fmt.Sprintf("- Dietary requirements: %s\n", strings.Join(req.DietaryPreferences, ", ")))
	}

	// Cooking time
	cookingTimeMap := map[string]string{
		"quick":  "under 30 minutes",
		"medium": "30-60 minutes",
		"long":   "over 60 minutes",
	}
	if timeDesc, ok := cookingTimeMap[req.CookingTime]; ok {
		builder.WriteString(fmt.Sprintf("- Total cooking time: %s\n", timeDesc))
	}

	// Difficulty
	if req.Difficulty != "" {
		builder.WriteString(fmt.Sprintf("- Difficulty level: %s\n", req.Difficulty))
	}

	builder.WriteString("\nPlease provide a recipe with:\n")
	builder.WriteString("1. A creative and appetizing title\n")
	builder.WriteString("2. A brief description (2-3 sentences)\n")
	builder.WriteString("3. Precise ingredient list with measurements in METRIC/EUROPEAN units:\n")
	builder.WriteString("   - Use grams (g) or kilograms (kg) for solid ingredients\n")
	builder.WriteString("   - Use milliliters (ml) or liters (l) for liquids\n")
	builder.WriteString("   - Use pieces, cloves, pinches for items like garlic, spices\n")
	builder.WriteString("   - DO NOT use cups, teaspoons (tsp), tablespoons (tbsp), or ounces\n")
	builder.WriteString("   - Use Celsius (°C) for all temperatures\n")
	builder.WriteString("4. Detailed step-by-step cooking instructions\n")
	builder.WriteString("5. Each step should include timing and temperature where relevant\n")
	builder.WriteString("6. Professional cooking tips and techniques\n")
	builder.WriteString("7. Serving size (number of people)\n\n")

	builder.WriteString("IMPORTANT: Return ONLY valid JSON in this EXACT format:\n")
	builder.WriteString("- serving_size, cooking_time, prep_time, step_number must be NUMBERS (not strings)\n")
	builder.WriteString("- ingredient quantity must be a STRING (e.g., \"500\" not 500)\n")
	builder.WriteString("- timing and temperature in steps must be STRINGS\n\n")
	builder.WriteString(`{
  "title": "Recipe Name",
  "description": "Brief description",
  "serving_size": 4,
  "cooking_time": 45,
  "prep_time": 15,
  "difficulty": "medium",
  "ingredients": [
    {"name": "chicken breast", "quantity": "500", "unit": "g"},
    {"name": "pasta", "quantity": "400", "unit": "g"}
  ],
  "steps": [
    {"step_number": 1, "instruction": "detailed instruction", "timing": "5 minutes", "temperature": "180°C"},
    {"step_number": 2, "instruction": "next instruction", "timing": "10 minutes", "temperature": ""}
  ],
  "tips": ["tip 1", "tip 2"],
  "cuisine_type": "Italian",
  "meat_type": "Chicken"
}`)

	builder.WriteString("\n\nGenerate the recipe now:")

	return builder.String()
}

// parseRecipeResponse parses Claude's JSON response into a Recipe struct
func (s *ClaudeService) parseRecipeResponse(response string, req models.RecipeGenerationRequest) (*models.RecipeDetail, error) {
	// Try to extract JSON from the response (Claude sometimes adds text before/after)
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}")

	if jsonStart == -1 || jsonEnd == -1 {
		return nil, ErrInvalidJSON
	}

	jsonStr := response[jsonStart : jsonEnd+1]

	// Parse into a temporary structure
	// Using json.RawMessage to handle flexible quantity types
	var temp struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		ServingSize int    `json:"serving_size"`
		CookingTime int    `json:"cooking_time"`
		PrepTime    int    `json:"prep_time"`
		Difficulty  string `json:"difficulty"`
		Ingredients []struct {
			Name     string          `json:"name"`
			Quantity json.RawMessage `json:"quantity"` // Can be string or number
			Unit     string          `json:"unit"`
		} `json:"ingredients"`
		Steps []struct {
			StepNumber  int    `json:"step_number"`
			Instruction string `json:"instruction"`
			Timing      string `json:"timing"`
			Temperature string `json:"temperature"`
		} `json:"steps"`
		Tips        []string `json:"tips"`
		CuisineType string   `json:"cuisine_type"`
		MeatType    string   `json:"meat_type"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &temp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Convert to RecipeDetail
	recipe := &models.RecipeDetail{
		Title:       temp.Title,
		Description: temp.Description,
		CookingTime: temp.CookingTime + temp.PrepTime, // Total time
		Difficulty:  temp.Difficulty,
		CuisineType: temp.CuisineType,
		MeatType:    temp.MeatType,
		Ingredients: make([]models.Ingredient, len(temp.Ingredients)),
		Steps:       make([]models.CookingStep, len(temp.Steps)),
		DietaryTags: req.DietaryPreferences,
	}

	// Convert ingredients (handle both string and number quantities)
	for i, ing := range temp.Ingredients {
		// Parse quantity - could be string "500" or number 500
		var quantityStr string
		if len(ing.Quantity) > 0 {
			// Try to unmarshal as string first
			if err := json.Unmarshal(ing.Quantity, &quantityStr); err != nil {
				// If that fails, try as number and convert to string
				var quantityNum float64
				if err := json.Unmarshal(ing.Quantity, &quantityNum); err == nil {
					quantityStr = fmt.Sprintf("%.0f", quantityNum)
				}
			}
		}

		recipe.Ingredients[i] = models.Ingredient{
			Name:     ing.Name,
			Quantity: quantityStr,
			Unit:     ing.Unit,
		}
	}

	// Convert steps
	for i, step := range temp.Steps {
		recipe.Steps[i] = models.CookingStep{
			StepNumber:  step.StepNumber,
			Instruction: step.Instruction,
			Timing:      step.Timing,
			Temperature: step.Temperature,
		}
	}

	return recipe, nil
}

// generateRecipeImageDescription generates a professional image description for the recipe
// This creates a data URL with an SVG that represents the dish
func (s *ClaudeService) generateRecipeImageDescription(title, cuisine, meatType string) string {
	// Create a professional-looking SVG placeholder with recipe information
	// Using food-related colors and professional styling

	// Color schemes based on cuisine/meat type
	colors := s.getRecipeColors(cuisine, meatType)

	// Clean title for display
	cleanTitle := title
	if len(cleanTitle) > 30 {
		cleanTitle = cleanTitle[:27] + "..."
	}

	// Generate SVG
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="800" height="600" viewBox="0 0 800 600">
		<defs>
			<linearGradient id="grad1" x1="0%%" y1="0%%" x2="100%%" y2="100%%">
				<stop offset="0%%" style="stop-color:%s;stop-opacity:1" />
				<stop offset="100%%" style="stop-color:%s;stop-opacity:1" />
			</linearGradient>
		</defs>
		<rect width="800" height="600" fill="url(#grad1)"/>
		<circle cx="400" cy="250" r="120" fill="rgba(255,255,255,0.1)" />
		<circle cx="400" cy="250" r="100" fill="rgba(255,255,255,0.15)" />
		<circle cx="400" cy="250" r="80" fill="rgba(255,255,255,0.2)" />
		<text x="400" y="430" font-family="Arial, sans-serif" font-size="36" font-weight="bold" fill="white" text-anchor="middle">%s</text>
		<text x="400" y="480" font-family="Arial, sans-serif" font-size="24" fill="rgba(255,255,255,0.9)" text-anchor="middle">%s Cuisine</text>
		<text x="400" y="520" font-family="Arial, sans-serif" font-size="20" fill="rgba(255,255,255,0.8)" text-anchor="middle">🍽️ Professional Recipe</text>
	</svg>`, colors[0], colors[1], cleanTitle, cuisine)

	// Return as data URL
	return "data:image/svg+xml;base64," + s.base64Encode(svg)
}

// getRecipeColors returns color scheme based on cuisine and meat type
func (s *ClaudeService) getRecipeColors(cuisine, meatType string) [2]string {
	cuisineColors := map[string][2]string{
		"Italian":        {"#E74C3C", "#C0392B"}, // Red tones
		"Mexican":        {"#E67E22", "#D35400"}, // Orange tones
		"Chinese":        {"#E74C3C", "#922B21"}, // Deep red
		"Indian":         {"#F39C12", "#D68910"}, // Golden
		"Japanese":       {"#EC7063", "#CB4335"}, // Salmon red
		"Thai":           {"#28B463", "#1E8449"}, // Green
		"Mediterranean":  {"#3498DB", "#2874A6"}, // Blue
		"American":       {"#E67E22", "#CA6F1E"}, // Orange
		"French":         {"#8E44AD", "#6C3483"}, // Purple
		"Greek":          {"#3498DB", "#21618C"}, // Blue
		"Korean":         {"#E74C3C", "#A93226"}, // Red
		"Vietnamese":     {"#2ECC71", "#27AE60"}, // Green
		"Spanish":        {"#E67E22", "#BA4A00"}, // Orange-red
		"Middle Eastern": {"#F4D03F", "#D4AC0D"}, // Gold
	}

	// Return cuisine-specific colors or default
	if colors, ok := cuisineColors[cuisine]; ok {
		return colors
	}

	// Default professional food colors
	return [2]string{"#FF6B6B", "#C44569"}
}

// base64Encode encodes a string to base64
func (s *ClaudeService) base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

