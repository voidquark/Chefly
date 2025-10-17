package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"chefly/models"
	"chefly/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RecipeHandler handles recipe operations
type RecipeHandler struct {
	db                    *sql.DB
	claudeService         *services.ClaudeService
	openaiService         *services.OpenAIService
	recipeGenerationLimit string
}

// NewRecipeHandler creates a new recipe handler
func NewRecipeHandler(db *sql.DB, claudeAPIKey, claudeModel, openaiAPIKey, openaiModel, recipeGenerationLimit string) *RecipeHandler {
	return &RecipeHandler{
		db:                    db,
		claudeService:         services.NewClaudeService(claudeAPIKey, claudeModel),
		openaiService:         services.NewOpenAIService(openaiAPIKey, openaiModel),
		recipeGenerationLimit: recipeGenerationLimit,
	}
}

// GenerateRecipe generates a new recipe using Claude AI
func (h *RecipeHandler) GenerateRecipe(c *gin.Context) {
	userID := c.GetString("user_id")

	// Get audit logger from context
	auditLogger, _ := c.Get("audit_logger")
	logger, _ := auditLogger.(*services.AuditLogger)
	requestID := c.GetString("request_id")

	// Check recipe generation limit
	if !h.canGenerateRecipe(userID, logger, requestID, c.ClientIP()) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Recipe generation limit reached"})
		return
	}

	var req models.RecipeGenerationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Log recipe generation start
	if logger != nil {
		logger.Info("recipe.generate_start", "Recipe generation initiated", &models.AuditContext{
			RequestID: requestID,
			UserID:    userID,
			IPAddress: c.ClientIP(),
			Metadata: map[string]interface{}{
				"meat_type":    req.MeatType,
				"cuisine_type": req.CuisineType,
				"difficulty":   req.Difficulty,
			},
		})
	}

	// Generate recipe using Claude
	recipe, err := h.claudeService.GenerateRecipe(req)
	if err != nil {
		// Determine specific error message based on error type
		var errorMessage string
		var statusCode int

		if errors.Is(err, services.ErrRateLimit) {
			errorMessage = "Rate limit reached. Please wait a moment before generating another recipe."
			statusCode = http.StatusTooManyRequests
		} else if errors.Is(err, services.ErrEmptyResponse) {
			errorMessage = "AI service returned an empty response. Please try again."
			statusCode = http.StatusInternalServerError
		} else if errors.Is(err, services.ErrInvalidJSON) {
			errorMessage = "Failed to parse AI response. The service may be experiencing issues. Please try again."
			statusCode = http.StatusInternalServerError
		} else if errors.Is(err, services.ErrParsingFailed) {
			errorMessage = "Failed to process AI response. Please try again with different settings."
			statusCode = http.StatusInternalServerError
		} else if errors.Is(err, services.ErrAPIConnection) {
			errorMessage = "AI service is temporarily unavailable. Please check your connection and try again."
			statusCode = http.StatusServiceUnavailable
		} else {
			errorMessage = "Failed to generate recipe. Please try again."
			statusCode = http.StatusInternalServerError
		}

		// Log recipe generation failure
		if logger != nil {
			logger.Error("recipe.generate_failure", "Recipe generation failed", err, &models.AuditContext{
				RequestID: requestID,
				UserID:    userID,
				IPAddress: c.ClientIP(),
				Metadata: map[string]interface{}{
					"error_type": errorMessage,
					"meat_type":  req.MeatType,
				},
			})
		}

		c.JSON(statusCode, gin.H{
			"error": errorMessage,
		})
		return
	}

	// Generate realistic food image using OpenAI DALL-E 3
	imageDataURL, err := h.openaiService.GenerateFoodImage(recipe.Title, recipe.CuisineType, recipe.Description)
	if err != nil {
		// Log error but don't fail - use fallback
		println("Warning: Failed to generate image:", err.Error())
		imageDataURL = ""
	}
	recipe.ImagePath = imageDataURL

	// Save recipe to database
	recipeID := uuid.New().String()
	recipe.ID = recipeID
	recipe.UserID = userID

	// Encode ingredients and steps as JSON
	ingredientsJSON, err := json.Marshal(recipe.Ingredients)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode ingredients"})
		return
	}

	stepsJSON, err := json.Marshal(recipe.Steps)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode steps"})
		return
	}

	dietaryTagsJSON, err := json.Marshal(recipe.DietaryTags)
	if err != nil {
		dietaryTagsJSON = []byte("[]")
	}

	// Insert into database
	_, err = h.db.Exec(`
		INSERT INTO recipes (
			id, user_id, title, description, ingredients, steps,
			cooking_time, difficulty, cuisine_type, meat_type,
			dietary_tags, is_favorite, image_path
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, ?)
	`, recipeID, userID, recipe.Title, recipe.Description,
		string(ingredientsJSON), string(stepsJSON),
		recipe.CookingTime, recipe.Difficulty, recipe.CuisineType,
		recipe.MeatType, string(dietaryTagsJSON), recipe.ImagePath)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save recipe",
		})
		return
	}

	// Log successful recipe generation
	if logger != nil {
		logger.Info("recipe.generate_success", "Recipe generated successfully", &models.AuditContext{
			RequestID: requestID,
			UserID:    userID,
			IPAddress: c.ClientIP(),
			Metadata: map[string]interface{}{
				"recipe_id":    recipeID,
				"recipe_title": recipe.Title,
				"meat_type":    recipe.MeatType,
				"cuisine_type": recipe.CuisineType,
				"difficulty":   recipe.Difficulty,
			},
		})
	}

	// Return the generated recipe
	c.JSON(http.StatusCreated, recipe)
}

// GetRecipes gets user's recipes
func (h *RecipeHandler) GetRecipes(c *gin.Context) {
	userID := c.GetString("user_id")

	rows, err := h.db.Query(`
		SELECT id, title, description, cuisine_type, difficulty, cooking_time, is_favorite, image_path, created_at
		FROM recipes
		WHERE user_id = ?
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipes"})
		return
	}
	defer rows.Close()

	recipes := []gin.H{}
	for rows.Next() {
		var id, title, description, cuisineType, difficulty, imagePath string
		var cookingTime int
		var isFavorite bool
		var createdAt string

		err := rows.Scan(&id, &title, &description, &cuisineType, &difficulty, &cookingTime, &isFavorite, &imagePath, &createdAt)
		if err != nil {
			continue
		}

		recipes = append(recipes, gin.H{
			"id":           id,
			"title":        title,
			"description":  description,
			"cuisine_type": cuisineType,
			"difficulty":   difficulty,
			"cooking_time": cookingTime,
			"is_favorite":  isFavorite,
			"image_path":   imagePath,
			"created_at":   createdAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"recipes": recipes})
}

// GetRecipe gets a single recipe with parsed ingredients and steps
func (h *RecipeHandler) GetRecipe(c *gin.Context) {
	recipeID := c.Param("id")
	userID := c.GetString("user_id")

	var id, title, description, ingredientsJSON, stepsJSON, cuisineType, meatType, difficulty, dietaryTagsJSON, imagePath string
	var cookingTime int
	var isFavorite bool
	var createdAt string

	err := h.db.QueryRow(`
		SELECT id, title, description, ingredients, steps, cuisine_type, meat_type, difficulty, dietary_tags, cooking_time, is_favorite, image_path, created_at
		FROM recipes
		WHERE id = ? AND user_id = ?
	`, recipeID, userID).Scan(&id, &title, &description, &ingredientsJSON, &stepsJSON, &cuisineType, &meatType, &difficulty, &dietaryTagsJSON, &cookingTime, &isFavorite, &imagePath, &createdAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipe"})
		return
	}

	// Parse JSON fields
	var ingredients []models.Ingredient
	var steps []models.CookingStep
	var dietaryTags []string

	json.Unmarshal([]byte(ingredientsJSON), &ingredients)
	json.Unmarshal([]byte(stepsJSON), &steps)
	json.Unmarshal([]byte(dietaryTagsJSON), &dietaryTags)

	recipe := models.RecipeDetail{
		ID:          id,
		UserID:      userID,
		Title:       title,
		Description: description,
		Ingredients: ingredients,
		Steps:       steps,
		CookingTime: cookingTime,
		Difficulty:  difficulty,
		CuisineType: cuisineType,
		MeatType:    meatType,
		DietaryTags: dietaryTags,
		IsFavorite:  isFavorite,
		ImagePath:   imagePath,
	}

	c.JSON(http.StatusOK, recipe)
}

// DeleteRecipe deletes a recipe
func (h *RecipeHandler) DeleteRecipe(c *gin.Context) {
	recipeID := c.Param("id")
	userID := c.GetString("user_id")

	// Get audit logger from context
	auditLogger, _ := c.Get("audit_logger")
	logger, _ := auditLogger.(*services.AuditLogger)
	requestID := c.GetString("request_id")

	// Get recipe title before deleting for audit log
	var recipeTitle string
	h.db.QueryRow("SELECT title FROM recipes WHERE id = ? AND user_id = ?", recipeID, userID).Scan(&recipeTitle)

	result, err := h.db.Exec("DELETE FROM recipes WHERE id = ? AND user_id = ?", recipeID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete recipe"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}

	// Log recipe deletion
	if logger != nil {
		logger.Info("recipe.delete", "Recipe deleted", &models.AuditContext{
			RequestID: requestID,
			UserID:    userID,
			IPAddress: c.ClientIP(),
			Metadata: map[string]interface{}{
				"recipe_id":    recipeID,
				"recipe_title": recipeTitle,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe deleted successfully"})
}

// ToggleFavorite toggles favorite status
func (h *RecipeHandler) ToggleFavorite(c *gin.Context) {
	recipeID := c.Param("id")
	userID := c.GetString("user_id")

	// Get audit logger from context
	auditLogger, _ := c.Get("audit_logger")
	logger, _ := auditLogger.(*services.AuditLogger)
	requestID := c.GetString("request_id")

	// Get recipe title and current favorite status for audit log
	var recipeTitle string
	var isFavorite bool
	h.db.QueryRow("SELECT title, is_favorite FROM recipes WHERE id = ? AND user_id = ?", recipeID, userID).Scan(&recipeTitle, &isFavorite)

	_, err := h.db.Exec(`
		UPDATE recipes
		SET is_favorite = NOT is_favorite
		WHERE id = ? AND user_id = ?
	`, recipeID, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update favorite status"})
		return
	}

	// Log favorite toggle
	if logger != nil {
		logger.Info("recipe.favorite_toggle", "Recipe favorite status toggled", &models.AuditContext{
			RequestID: requestID,
			UserID:    userID,
			IPAddress: c.ClientIP(),
			Metadata: map[string]interface{}{
				"recipe_id":     recipeID,
				"recipe_title":  recipeTitle,
				"new_is_favorite": !isFavorite,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Favorite status updated"})
}

// GetCountries returns available countries
func (h *RecipeHandler) GetCountries(c *gin.Context) {
	countries := []string{
		"Italian", "Mexican", "Chinese", "Indian", "Japanese",
		"Thai", "Mediterranean", "American", "French", "Greek",
		"Korean", "Vietnamese", "Spanish", "Middle Eastern",
	}
	c.JSON(http.StatusOK, gin.H{"countries": countries})
}

// GetMeatTypes returns available meat types
func (h *RecipeHandler) GetMeatTypes(c *gin.Context) {
	meatTypes := []string{
		"Chicken", "Beef", "Pork", "Fish", "Seafood",
		"Lamb", "Turkey", "None (Vegetarian)",
	}
	c.JSON(http.StatusOK, gin.H{"meat_types": meatTypes})
}

// GetIngredients returns available side ingredients
func (h *RecipeHandler) GetIngredients(c *gin.Context) {
	ingredients := []string{
		"Vegetables", "Rice", "Pasta", "Potatoes", "Grains",
		"Legumes", "Noodles", "Bread", "Quinoa", "Couscous",
	}
	c.JSON(http.StatusOK, gin.H{"ingredients": ingredients})
}

// GetPublicRecipe gets a recipe by ID without authentication (for sharing)
func (h *RecipeHandler) GetPublicRecipe(c *gin.Context) {
	recipeID := c.Param("id")

	var id, title, description, ingredientsJSON, stepsJSON, cuisineType, meatType, difficulty, dietaryTagsJSON, imagePath string
	var cookingTime int
	var createdAt string

	err := h.db.QueryRow(`
		SELECT id, title, description, ingredients, steps, cuisine_type, meat_type, difficulty, dietary_tags, cooking_time, image_path, created_at
		FROM recipes
		WHERE id = ?
	`, recipeID).Scan(&id, &title, &description, &ingredientsJSON, &stepsJSON, &cuisineType, &meatType, &difficulty, &dietaryTagsJSON, &cookingTime, &imagePath, &createdAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipe"})
		return
	}

	// Parse JSON fields
	var ingredients []models.Ingredient
	var steps []models.CookingStep
	var dietaryTags []string

	json.Unmarshal([]byte(ingredientsJSON), &ingredients)
	json.Unmarshal([]byte(stepsJSON), &steps)
	json.Unmarshal([]byte(dietaryTagsJSON), &dietaryTags)

	recipe := gin.H{
		"id":           id,
		"title":        title,
		"description":  description,
		"ingredients":  ingredients,
		"steps":        steps,
		"cooking_time": cookingTime,
		"difficulty":   difficulty,
		"cuisine_type": cuisineType,
		"meat_type":    meatType,
		"dietary_tags": dietaryTags,
		"image_path":   imagePath,
		"created_at":   createdAt,
	}

	c.JSON(http.StatusOK, recipe)
}

// canGenerateRecipe checks if user can generate a recipe based on limits
func (h *RecipeHandler) canGenerateRecipe(userID string, logger *services.AuditLogger, requestID, ipAddress string) bool {
	// Get user's recipe_limit setting from database
	var recipeLimit sql.NullInt64
	err := h.db.QueryRow("SELECT recipe_limit FROM users WHERE id = ?", userID).Scan(&recipeLimit)
	if err != nil {
		// If error, allow generation (fail open, but log error)
		if logger != nil {
			logger.Error("recipe.limit_check_failed", "Failed to check recipe limit", err, &models.AuditContext{
				RequestID: requestID,
				UserID:    userID,
				IPAddress: ipAddress,
			})
		}
		return true
	}

	// Determine effective limit
	var effectiveLimit int
	var hasLimit bool

	if recipeLimit.Valid {
		// User has a personal limit set
		effectiveLimit = int(recipeLimit.Int64)
		hasLimit = true

		// -1 means unlimited for this user
		if effectiveLimit == -1 {
			return true
		}
	} else {
		// Use global limit from config
		globalLimit := h.recipeGenerationLimit

		// Check if global limit is "unlimited" or empty
		if globalLimit == "unlimited" || globalLimit == "" {
			return true
		}

		// Try to parse as number
		parsedLimit, err := strconv.Atoi(globalLimit)
		if err != nil {
			// Invalid config, default to unlimited
			return true
		}

		effectiveLimit = parsedLimit
		hasLimit = true
	}

	// If limit is 0, block generation
	if hasLimit && effectiveLimit == 0 {
		if logger != nil {
			logger.Warn("recipe.limit_blocked", "Recipe generation blocked by limit", &models.AuditContext{
				RequestID: requestID,
				UserID:    userID,
				IPAddress: ipAddress,
				Metadata: map[string]interface{}{
					"effective_limit": effectiveLimit,
					"has_personal_limit": recipeLimit.Valid,
				},
			})
		}
		return false
	}

	// Count user's existing recipes
	var recipeCount int
	err = h.db.QueryRow("SELECT COUNT(*) FROM recipes WHERE user_id = ?", userID).Scan(&recipeCount)
	if err != nil {
		// If error, allow generation (fail open)
		return true
	}

	// Check if user has reached their limit
	if hasLimit && recipeCount >= effectiveLimit {
		if logger != nil {
			logger.Warn("recipe.limit_reached", "Recipe generation limit reached", &models.AuditContext{
				RequestID: requestID,
				UserID:    userID,
				IPAddress: ipAddress,
				Metadata: map[string]interface{}{
					"current_count":      recipeCount,
					"effective_limit":    effectiveLimit,
					"has_personal_limit": recipeLimit.Valid,
				},
			})
		}
		return false
	}

	return true
}
