package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"chefly/models"
	"chefly/services"

	"github.com/gin-gonic/gin"
)

// AdminHandler handles admin operations
type AdminHandler struct {
	db *sql.DB
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(db *sql.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

// UserWithStats represents a user with additional statistics
type UserWithStats struct {
	ID              string    `json:"id"`
	Email           string    `json:"email"`
	Username        string    `json:"username"`
	IsAdmin         bool      `json:"is_admin"`
	CreatedAt       time.Time `json:"created_at"`
	RecipeCount     int       `json:"recipe_count"`
	ShoppingItems   int       `json:"shopping_items"`
	LastRecipeDate  *string   `json:"last_recipe_date"`
	RecipeLimit     *int      `json:"recipe_limit,omitempty"` // null = use global, -1 = unlimited, 0 = blocked, >0 = custom
}

// AdminStats represents admin dashboard statistics
type AdminStats struct {
	TotalUsers          int       `json:"total_users"`
	TotalRecipes        int       `json:"total_recipes"`
	TotalShoppingItems  int       `json:"total_shopping_items"`
	AverageRecipesUser  float64   `json:"average_recipes_per_user"`
	MostActiveUser      *string   `json:"most_active_user"`
	MostActiveUserCount int       `json:"most_active_user_count"`
	RecentRegistrations int       `json:"recent_registrations"`
	FirstUserDate       time.Time `json:"first_user_date"`
}

// GetAllUsers returns all users with their statistics
func (h *AdminHandler) GetAllUsers(c *gin.Context) {
	// Get audit logger from context
	auditLogger, _ := c.Get("audit_logger")
	logger, _ := auditLogger.(*services.AuditLogger)
	requestID := c.GetString("request_id")
	adminID := c.GetString("user_id")

	// Query all users with recipe count, shopping items, last recipe date, and recipe limit
	query := `
		SELECT
			u.id,
			u.email,
			u.username,
			u.is_admin,
			u.created_at,
			u.recipe_limit,
			COALESCE(COUNT(DISTINCT r.id), 0) as recipe_count,
			COALESCE(COUNT(DISTINCT s.id), 0) as shopping_items,
			MAX(r.created_at) as last_recipe_date
		FROM users u
		LEFT JOIN recipes r ON u.id = r.user_id
		LEFT JOIN shopping_list_items s ON u.id = s.user_id
		GROUP BY u.id, u.email, u.username, u.is_admin, u.created_at, u.recipe_limit
		ORDER BY u.created_at ASC
	`

	rows, err := h.db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	users := []UserWithStats{}
	for rows.Next() {
		var user UserWithStats
		var isAdminInt int
		var lastRecipeDate sql.NullString
		var recipeLimit sql.NullInt64

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Username,
			&isAdminInt,
			&user.CreatedAt,
			&recipeLimit,
			&user.RecipeCount,
			&user.ShoppingItems,
			&lastRecipeDate,
		)
		if err != nil {
			continue
		}

		user.IsAdmin = isAdminInt == 1
		if lastRecipeDate.Valid {
			user.LastRecipeDate = &lastRecipeDate.String
		}
		if recipeLimit.Valid {
			limit := int(recipeLimit.Int64)
			user.RecipeLimit = &limit
		}

		users = append(users, user)
	}

	// Log admin viewing user list
	if logger != nil {
		logger.Info("admin.user_list", "Admin viewed user list", &models.AuditContext{
			RequestID: requestID,
			UserID:    adminID,
			IPAddress: c.ClientIP(),
			Metadata: map[string]interface{}{
				"total_users": len(users),
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// DeleteUser deletes a user and all their data (CASCADE)
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	adminID := c.GetString("user_id")

	// Get audit logger from context
	auditLogger, _ := c.Get("audit_logger")
	logger, _ := auditLogger.(*services.AuditLogger)
	requestID := c.GetString("request_id")

	// Prevent admin from deleting themselves
	if userID == adminID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete your own admin account"})
		return
	}

	// Get user info before deleting for audit log
	var username, email string
	var recipeCount, shoppingCount int
	h.db.QueryRow("SELECT username, email FROM users WHERE id = ?", userID).Scan(&username, &email)
	h.db.QueryRow("SELECT COUNT(*) FROM recipes WHERE user_id = ?", userID).Scan(&recipeCount)
	h.db.QueryRow("SELECT COUNT(*) FROM shopping_list_items WHERE user_id = ?", userID).Scan(&shoppingCount)

	// Check if user exists
	var exists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", userID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Delete user (CASCADE will handle recipes and shopping items)
	result, err := h.db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Log user deletion by admin
	if logger != nil {
		logger.Warn("admin.user_delete", "Admin deleted user", &models.AuditContext{
			RequestID: requestID,
			UserID:    adminID,
			IPAddress: c.ClientIP(),
			Metadata: map[string]interface{}{
				"deleted_user_id":       userID,
				"deleted_username":      username,
				"deleted_email":         email,
				"deleted_recipe_count":  recipeCount,
				"deleted_shopping_count": shoppingCount,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "User and all their data deleted successfully"})
}

// GetAdminStats returns admin dashboard statistics
func (h *AdminHandler) GetAdminStats(c *gin.Context) {
	// Get audit logger from context
	auditLogger, _ := c.Get("audit_logger")
	logger, _ := auditLogger.(*services.AuditLogger)
	requestID := c.GetString("request_id")
	adminID := c.GetString("user_id")

	var stats AdminStats

	// Total users
	err := h.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&stats.TotalUsers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stats"})
		return
	}

	// Total recipes
	err = h.db.QueryRow("SELECT COUNT(*) FROM recipes").Scan(&stats.TotalRecipes)
	if err != nil {
		stats.TotalRecipes = 0
	}

	// Total shopping items
	err = h.db.QueryRow("SELECT COUNT(*) FROM shopping_list_items").Scan(&stats.TotalShoppingItems)
	if err != nil {
		stats.TotalShoppingItems = 0
	}

	// Average recipes per user
	if stats.TotalUsers > 0 {
		stats.AverageRecipesUser = float64(stats.TotalRecipes) / float64(stats.TotalUsers)
	}

	// Most active user (user with most recipes)
	var mostActiveUser sql.NullString
	err = h.db.QueryRow(`
		SELECT u.username, COUNT(r.id) as recipe_count
		FROM users u
		LEFT JOIN recipes r ON u.id = r.user_id
		GROUP BY u.id, u.username
		ORDER BY recipe_count DESC
		LIMIT 1
	`).Scan(&mostActiveUser, &stats.MostActiveUserCount)
	if err == nil && mostActiveUser.Valid {
		stats.MostActiveUser = &mostActiveUser.String
	}

	// Recent registrations (last 7 days)
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	err = h.db.QueryRow(
		"SELECT COUNT(*) FROM users WHERE created_at >= ?",
		sevenDaysAgo.Format("2006-01-02 15:04:05"),
	).Scan(&stats.RecentRegistrations)
	if err != nil {
		stats.RecentRegistrations = 0
	}

	// First user registration date
	err = h.db.QueryRow("SELECT MIN(created_at) FROM users").Scan(&stats.FirstUserDate)
	if err != nil {
		stats.FirstUserDate = time.Now()
	}

	// Log admin viewing stats
	if logger != nil {
		logger.Info("admin.stats_view", "Admin viewed dashboard statistics", &models.AuditContext{
			RequestID: requestID,
			UserID:    adminID,
			IPAddress: c.ClientIP(),
		})
	}

	c.JSON(http.StatusOK, stats)
}

// UpdateUserRecipeLimit updates a user's recipe generation limit
func (h *AdminHandler) UpdateUserRecipeLimit(c *gin.Context) {
	userID := c.Param("id")
	adminID := c.GetString("user_id")

	// Get audit logger from context
	auditLogger, _ := c.Get("audit_logger")
	logger, _ := auditLogger.(*services.AuditLogger)
	requestID := c.GetString("request_id")

	var req struct {
		RecipeLimit *int `json:"recipe_limit"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Validate recipe_limit value if not null
	if req.RecipeLimit != nil {
		limit := *req.RecipeLimit
		// Must be: null (use global), -1 (unlimited), 0 (blocked), or positive number
		if limit < -1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe limit value. Must be null, -1 (unlimited), 0 (blocked), or positive number"})
			return
		}
	}

	// Check if user exists
	var exists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", userID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Get user info before update for audit log
	var username, email string
	var oldLimit sql.NullInt64
	h.db.QueryRow("SELECT username, email, recipe_limit FROM users WHERE id = ?", userID).Scan(&username, &email, &oldLimit)

	// Update user's recipe_limit
	var result sql.Result
	if req.RecipeLimit == nil {
		// Set to NULL (use global limit)
		result, err = h.db.Exec("UPDATE users SET recipe_limit = NULL WHERE id = ?", userID)
	} else {
		// Set to specific value
		result, err = h.db.Exec("UPDATE users SET recipe_limit = ? WHERE id = ?", *req.RecipeLimit, userID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update recipe limit"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Log recipe limit update
	if logger != nil {
		metadata := map[string]interface{}{
			"target_user_id": userID,
			"target_username": username,
			"target_email":    email,
		}

		if oldLimit.Valid {
			metadata["old_limit"] = oldLimit.Int64
		} else {
			metadata["old_limit"] = "global"
		}

		if req.RecipeLimit == nil {
			metadata["new_limit"] = "global"
		} else {
			metadata["new_limit"] = *req.RecipeLimit
		}

		logger.Info("admin.user_recipe_limit_update", "Admin updated user recipe limit", &models.AuditContext{
			RequestID: requestID,
			UserID:    adminID,
			IPAddress: c.ClientIP(),
			Metadata:  metadata,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe limit updated successfully",
		"user_id": userID,
		"recipe_limit": req.RecipeLimit,
	})
}
