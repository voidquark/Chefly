package models

import "time"

// Recipe represents a recipe in the database
type Recipe struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Ingredients string    `json:"ingredients"` // JSON encoded array
	Steps       string    `json:"steps"`       // JSON encoded array
	CookingTime int       `json:"cooking_time"` // minutes
	Difficulty  string    `json:"difficulty"`
	CuisineType string    `json:"cuisine_type"`
	MeatType    string    `json:"meat_type"`
	DietaryTags string    `json:"dietary_tags"` // JSON encoded array
	IsFavorite  bool      `json:"is_favorite"`
	ImagePath   string    `json:"image_path"`
	CreatedAt   time.Time `json:"created_at"`
}

// RecipeGenerationRequest represents a recipe generation request
type RecipeGenerationRequest struct {
	MeatType            string   `json:"meat_type"`
	SideIngredients     []string `json:"side_ingredients"`
	CuisineType         string   `json:"cuisine_type"`
	DietaryPreferences  []string `json:"dietary_preferences"`
	CookingTime         string   `json:"cooking_time"`  // "quick", "medium", "long"
	Difficulty          string   `json:"difficulty"`    // "easy", "medium", "hard"
	Language            string   `json:"language"`      // "en" or "sk"
}

// RecipeFilter represents filters for searching recipes
type RecipeFilter struct {
	ID                 string `json:"id"`
	RecipeID           string `json:"recipe_id"`
	MeatCategory       string `json:"meat_category"`
	SideIngredients    string `json:"side_ingredients"` // JSON encoded
	Country            string `json:"country"`
	DietaryPreferences string `json:"dietary_preferences"` // JSON encoded
}

// Ingredient represents a single ingredient with quantity
type Ingredient struct {
	Name     string `json:"name"`
	Quantity string `json:"quantity"`
	Unit     string `json:"unit"`
}

// CookingStep represents a single cooking step
type CookingStep struct {
	StepNumber  int    `json:"step_number"`
	Instruction string `json:"instruction"`
	Timing      string `json:"timing,omitempty"`      // e.g., "5 minutes"
	Temperature string `json:"temperature,omitempty"` // e.g., "180°C"
}

// RecipeDetail represents the full recipe with parsed ingredients and steps
type RecipeDetail struct {
	ID          string        `json:"id"`
	UserID      string        `json:"user_id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Ingredients []Ingredient  `json:"ingredients"`
	Steps       []CookingStep `json:"steps"`
	CookingTime int           `json:"cooking_time"`
	Difficulty  string        `json:"difficulty"`
	CuisineType string        `json:"cuisine_type"`
	MeatType    string        `json:"meat_type"`
	DietaryTags []string      `json:"dietary_tags"`
	IsFavorite  bool          `json:"is_favorite"`
	ImagePath   string        `json:"image_path"`
	CreatedAt   time.Time     `json:"created_at"`
}
