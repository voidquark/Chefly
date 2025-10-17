package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"chefly/models"
	"chefly/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ShoppingListHandler handles shopping list operations
type ShoppingListHandler struct {
	db *sql.DB
}

// NewShoppingListHandler creates a new shopping list handler
func NewShoppingListHandler(db *sql.DB) *ShoppingListHandler {
	return &ShoppingListHandler{db: db}
}

// GetShoppingList gets all shopping list items for the user
func (h *ShoppingListHandler) GetShoppingList(c *gin.Context) {
	userID := c.GetString("user_id")

	rows, err := h.db.Query(`
		SELECT id, user_id, COALESCE(recipe_id, ''), COALESCE(recipe_title, ''),
		       ingredient_name, quantity, unit, is_checked, created_at
		FROM shopping_list_items
		WHERE user_id = ?
		ORDER BY is_checked ASC, created_at DESC
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shopping list"})
		return
	}
	defer rows.Close()

	var items []models.ShoppingListItem
	for rows.Next() {
		var item models.ShoppingListItem
		var isChecked int
		err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.RecipeID,
			&item.RecipeTitle,
			&item.IngredientName,
			&item.Quantity,
			&item.Unit,
			&isChecked,
			&item.CreatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan shopping list item"})
			return
		}
		item.IsChecked = isChecked == 1
		items = append(items, item)
	}

	if items == nil {
		items = []models.ShoppingListItem{}
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

// AddRecipeToShoppingList adds all ingredients from a recipe to the shopping list
func (h *ShoppingListHandler) AddRecipeToShoppingList(c *gin.Context) {
	userID := c.GetString("user_id")

	var req models.AddToShoppingListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Fetch recipe
	var recipe models.Recipe
	var ingredientsJSON, stepsJSON, dietaryTagsJSON string
	err := h.db.QueryRow(`
		SELECT id, user_id, title, description, ingredients, steps,
		       cooking_time, difficulty, cuisine_type, meat_type, dietary_tags,
		       is_favorite, image_path, created_at
		FROM recipes
		WHERE id = ? AND user_id = ?
	`, req.RecipeID, userID).Scan(
		&recipe.ID,
		&recipe.UserID,
		&recipe.Title,
		&recipe.Description,
		&ingredientsJSON,
		&stepsJSON,
		&recipe.CookingTime,
		&recipe.Difficulty,
		&recipe.CuisineType,
		&recipe.MeatType,
		&dietaryTagsJSON,
		&recipe.IsFavorite,
		&recipe.ImagePath,
		&recipe.CreatedAt,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipe"})
		return
	}

	// Parse ingredients
	var ingredients []models.Ingredient
	if err := json.Unmarshal([]byte(ingredientsJSON), &ingredients); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse ingredients"})
		return
	}

	// Get audit logger from context
	auditLogger, _ := c.Get("audit_logger")
	logger, _ := auditLogger.(*services.AuditLogger)
	requestID := c.GetString("request_id")

	// Add each ingredient to shopping list
	addedCount := 0
	for _, ingredient := range ingredients {
		itemID := uuid.New().String()
		_, err := h.db.Exec(`
			INSERT INTO shopping_list_items (id, user_id, recipe_id, recipe_title, ingredient_name, quantity, unit)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`, itemID, userID, recipe.ID, recipe.Title, ingredient.Name, ingredient.Quantity, ingredient.Unit)
		if err != nil {
			continue // Skip if fails (e.g., duplicate)
		}
		addedCount++
	}

	// Log adding recipe to shopping list
	if logger != nil {
		logger.Info("shopping.add_recipe", "Recipe added to shopping list", &models.AuditContext{
			RequestID: requestID,
			UserID:    userID,
			IPAddress: c.ClientIP(),
			Metadata: map[string]interface{}{
				"recipe_id":        recipe.ID,
				"recipe_title":     recipe.Title,
				"items_added":      addedCount,
				"total_ingredients": len(ingredients),
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Ingredients added to shopping list",
		"added":   addedCount,
	})
}

// ToggleItemChecked toggles the checked status of a shopping list item
func (h *ShoppingListHandler) ToggleItemChecked(c *gin.Context) {
	userID := c.GetString("user_id")
	itemID := c.Param("id")

	// Toggle the is_checked field
	result, err := h.db.Exec(`
		UPDATE shopping_list_items
		SET is_checked = CASE WHEN is_checked = 0 THEN 1 ELSE 0 END
		WHERE id = ? AND user_id = ?
	`, itemID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle item"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item toggled successfully"})
}

// DeleteItem deletes a shopping list item
func (h *ShoppingListHandler) DeleteItem(c *gin.Context) {
	userID := c.GetString("user_id")
	itemID := c.Param("id")

	result, err := h.db.Exec(`
		DELETE FROM shopping_list_items
		WHERE id = ? AND user_id = ?
	`, itemID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}

// ClearCheckedItems deletes all checked items from the shopping list
func (h *ShoppingListHandler) ClearCheckedItems(c *gin.Context) {
	userID := c.GetString("user_id")

	// Get audit logger from context
	auditLogger, _ := c.Get("audit_logger")
	logger, _ := auditLogger.(*services.AuditLogger)
	requestID := c.GetString("request_id")

	result, err := h.db.Exec(`
		DELETE FROM shopping_list_items
		WHERE user_id = ? AND is_checked = 1
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear checked items"})
		return
	}

	rowsAffected, _ := result.RowsAffected()

	// Log clearing checked items
	if logger != nil {
		logger.Info("shopping.clear_checked", "Cleared checked shopping items", &models.AuditContext{
			RequestID: requestID,
			UserID:    userID,
			IPAddress: c.ClientIP(),
			Metadata: map[string]interface{}{
				"items_deleted": rowsAffected,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Checked items cleared",
		"deleted": rowsAffected,
	})
}

// ClearAllItems deletes all items from the shopping list
func (h *ShoppingListHandler) ClearAllItems(c *gin.Context) {
	userID := c.GetString("user_id")

	// Get audit logger from context
	auditLogger, _ := c.Get("audit_logger")
	logger, _ := auditLogger.(*services.AuditLogger)
	requestID := c.GetString("request_id")

	result, err := h.db.Exec(`
		DELETE FROM shopping_list_items
		WHERE user_id = ?
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear shopping list"})
		return
	}

	rowsAffected, _ := result.RowsAffected()

	// Log clearing all items
	if logger != nil {
		logger.Info("shopping.clear_all", "Cleared all shopping items", &models.AuditContext{
			RequestID: requestID,
			UserID:    userID,
			IPAddress: c.ClientIP(),
			Metadata: map[string]interface{}{
				"items_deleted": rowsAffected,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Shopping list cleared",
		"deleted": rowsAffected,
	})
}
