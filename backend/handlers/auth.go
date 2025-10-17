package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"chefly/models"
	"chefly/services"
	"chefly/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler handles authentication operations
type AuthHandler struct {
	db                  *sql.DB
	jwtSecret           string
	registrationEnabled bool
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(db *sql.DB, jwtSecret string, registrationEnabled bool) *AuthHandler {
	return &AuthHandler{
		db:                  db,
		jwtSecret:           jwtSecret,
		registrationEnabled: registrationEnabled,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	// Check if registration is enabled
	if !h.registrationEnabled {
		c.JSON(http.StatusForbidden, gin.H{"error": "Registration is currently disabled"})
		return
	}

	var req models.UserRegistration
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Validate email
	if err := utils.ValidateEmail(req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate username
	if err := utils.ValidateUsername(req.Username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate password
	if err := utils.ValidatePassword(req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Sanitize username
	req.Username = utils.SanitizeHTML(req.Username)

	// Get audit logger from context
	auditLogger, _ := c.Get("audit_logger")
	logger, _ := auditLogger.(*services.AuditLogger)
	requestID := c.GetString("request_id")

	// Check if user already exists
	var exists bool
	err := h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", req.Email).Scan(&exists)
	if err != nil {
		if logger != nil {
			logger.Error("auth.register.failed", "Registration failed: database error", err, &models.AuditContext{
				RequestID:      requestID,
				EmailAttempted: req.Email,
				IPAddress:      c.ClientIP(),
				UserAgent:      c.Request.UserAgent(),
				Metadata:       map[string]interface{}{"failure_reason": "database_error"},
			})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists {
		if logger != nil {
			logger.Warn("auth.register.failed", "Registration failed: email already exists", &models.AuditContext{
				RequestID:      requestID,
				EmailAttempted: req.Email,
				IPAddress:      c.ClientIP(),
				UserAgent:      c.Request.UserAgent(),
				Metadata:       map[string]interface{}{"failure_reason": "email_exists"},
			})
		}
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Check if this is the first user (will be admin)
	var userCount int
	err = h.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	isAdmin := userCount == 0 // First user becomes admin

	// Create user
	userID := uuid.New().String()
	_, err = h.db.Exec(
		"INSERT INTO users (id, email, password_hash, username, is_admin) VALUES (?, ?, ?, ?, ?)",
		userID, req.Email, passwordHash, req.Username, isAdmin,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Fetch created user
	var user models.User
	var isAdminInt int
	err = h.db.QueryRow(
		"SELECT id, email, username, is_admin, created_at FROM users WHERE id = ?",
		userID,
	).Scan(&user.ID, &user.Email, &user.Username, &isAdminInt, &user.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}
	user.IsAdmin = isAdminInt == 1

	// Generate access token (15 minutes)
	accessToken, err := utils.GenerateJWT(userID, req.Email, user.Username, user.IsAdmin, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Generate refresh token (7 days)
	refreshToken, err := h.createRefreshToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Log successful registration
	if logger != nil {
		logger.Info("auth.register.success", "User registered successfully", &models.AuditContext{
			RequestID: requestID,
			UserID:    userID,
			Username:  user.Username,
			Email:     user.Email,
			IsAdmin:   user.IsAdmin,
			IPAddress: c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Metadata:  map[string]interface{}{"is_first_user": isAdmin},
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":          user.ToResponse(),
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.UserLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Get audit logger from context
	auditLogger, _ := c.Get("audit_logger")
	logger, _ := auditLogger.(*services.AuditLogger)
	requestID := c.GetString("request_id")

	// Fetch user
	var user models.User
	var isAdminInt int
	err := h.db.QueryRow(
		"SELECT id, email, password_hash, username, is_admin, created_at FROM users WHERE email = ?",
		req.Email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Username, &isAdminInt, &user.CreatedAt)
	if err == sql.ErrNoRows {
		// Log failed login attempt for non-registered user
		if logger != nil {
			logger.Warn("auth.login.failed", "Login failed: user not found", &models.AuditContext{
				RequestID:      requestID,
				EmailAttempted: req.Email,
				IPAddress:      c.ClientIP(),
				UserAgent:      c.Request.UserAgent(),
				Metadata:       map[string]interface{}{"failure_reason": "user_not_found"},
			})
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	if err != nil {
		if logger != nil {
			logger.Error("auth.login.failed", "Login failed: database error", err, &models.AuditContext{
				RequestID:      requestID,
				EmailAttempted: req.Email,
				IPAddress:      c.ClientIP(),
				UserAgent:      c.Request.UserAgent(),
				Metadata:       map[string]interface{}{"failure_reason": "database_error"},
			})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	user.IsAdmin = isAdminInt == 1

	// Check password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		// Log failed login attempt for registered user with wrong password
		if logger != nil {
			logger.Warn("auth.login.failed", "Login failed: invalid password", &models.AuditContext{
				RequestID:      requestID,
				EmailAttempted: req.Email,
				IPAddress:      c.ClientIP(),
				UserAgent:      c.Request.UserAgent(),
				Metadata:       map[string]interface{}{"failure_reason": "invalid_password"},
			})
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate access token (15 minutes)
	accessToken, err := utils.GenerateJWT(user.ID, user.Email, user.Username, user.IsAdmin, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Generate refresh token (7 days)
	refreshToken, err := h.createRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Log successful login
	if logger != nil {
		logger.Info("auth.login.success", "User logged in successfully", &models.AuditContext{
			RequestID: requestID,
			UserID:    user.ID,
			Username:  user.Username,
			Email:     user.Email,
			IsAdmin:   user.IsAdmin,
			IPAddress: c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"user":          user.ToResponse(),
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// GetProfile gets user profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	var user models.User
	var isAdminInt int
	err := h.db.QueryRow(
		"SELECT id, email, username, is_admin, created_at FROM users WHERE id = ?",
		userID,
	).Scan(&user.ID, &user.Email, &user.Username, &isAdminInt, &user.CreatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user.IsAdmin = isAdminInt == 1

	c.JSON(http.StatusOK, user.ToResponse())
}

// UpdateProfile updates user profile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		Username string `json:"username" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Validate username
	if err := utils.ValidateUsername(req.Username); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Sanitize username
	req.Username = utils.SanitizeHTML(req.Username)

	_, err := h.db.Exec(
		"UPDATE users SET username = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		req.Username, userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// GetStats gets user statistics
func (h *AuthHandler) GetStats(c *gin.Context) {
	userID := c.GetString("user_id")

	var totalRecipes int
	var favoriteRecipes sql.NullInt64
	err := h.db.QueryRow(
		"SELECT COUNT(*), COALESCE(SUM(CASE WHEN is_favorite = 1 THEN 1 ELSE 0 END), 0) FROM recipes WHERE user_id = ?",
		userID,
	).Scan(&totalRecipes, &favoriteRecipes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stats"})
		return
	}

	favCount := int64(0)
	if favoriteRecipes.Valid {
		favCount = favoriteRecipes.Int64
	}

	c.JSON(http.StatusOK, gin.H{
		"total_recipes":    totalRecipes,
		"favorite_recipes": favCount,
	})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Get audit logger from context
	auditLogger, _ := c.Get("audit_logger")
	logger, _ := auditLogger.(*services.AuditLogger)
	requestID := c.GetString("request_id")

	// Validate refresh token
	var tokenRecord models.RefreshToken
	var revokedInt int
	err := h.db.QueryRow(`
		SELECT id, user_id, token, expires_at, created_at, revoked
		FROM refresh_tokens
		WHERE token = ?
	`, req.RefreshToken).Scan(
		&tokenRecord.ID,
		&tokenRecord.UserID,
		&tokenRecord.Token,
		&tokenRecord.ExpiresAt,
		&tokenRecord.CreatedAt,
		&revokedInt,
	)
	tokenRecord.Revoked = revokedInt == 1

	if err == sql.ErrNoRows {
		if logger != nil {
			logger.Warn("auth.refresh.invalid_token", "Refresh token not found", &models.AuditContext{
				RequestID: requestID,
				IPAddress: c.ClientIP(),
			})
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check if token is revoked
	if tokenRecord.Revoked {
		if logger != nil {
			logger.Warn("auth.refresh.revoked_token", "Refresh token revoked", &models.AuditContext{
				RequestID: requestID,
				UserID:    tokenRecord.UserID,
				IPAddress: c.ClientIP(),
			})
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token revoked"})
		return
	}

	// Check if token is expired
	if time.Now().After(tokenRecord.ExpiresAt) {
		if logger != nil {
			logger.Warn("auth.refresh.expired_token", "Refresh token expired", &models.AuditContext{
				RequestID: requestID,
				UserID:    tokenRecord.UserID,
				IPAddress: c.ClientIP(),
			})
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
		return
	}

	// Revoke old refresh token (rotating tokens)
	_, err = h.db.Exec("UPDATE refresh_tokens SET revoked = 1 WHERE id = ?", tokenRecord.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke old token"})
		return
	}

	// Get user details
	var user models.User
	var isAdminInt int
	err = h.db.QueryRow(`
		SELECT id, email, username, is_admin
		FROM users
		WHERE id = ?
	`, tokenRecord.UserID).Scan(&user.ID, &user.Email, &user.Username, &isAdminInt)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	user.IsAdmin = isAdminInt == 1

	// Generate new access token
	accessToken, err := utils.GenerateJWT(user.ID, user.Email, user.Username, user.IsAdmin, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Generate new refresh token
	newRefreshToken, err := h.createRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Log successful token refresh
	if logger != nil {
		logger.Info("auth.refresh.success", "Token refreshed successfully", &models.AuditContext{
			RequestID: requestID,
			UserID:    user.ID,
			Username:  user.Username,
			Email:     user.Email,
			IPAddress: c.ClientIP(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	})
}

// Logout handles user logout by revoking refresh token
func (h *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Get audit logger from context
	auditLogger, _ := c.Get("audit_logger")
	logger, _ := auditLogger.(*services.AuditLogger)
	requestID := c.GetString("request_id")
	userID := c.GetString("user_id")

	// Revoke refresh token
	result, err := h.db.Exec("UPDATE refresh_tokens SET revoked = 1 WHERE token = ?", req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke token"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		// Token not found, but that's OK - user is logging out anyway
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
		return
	}

	// Log successful logout
	if logger != nil {
		logger.Info("auth.logout.success", "User logged out successfully", &models.AuditContext{
			RequestID: requestID,
			UserID:    userID,
			IPAddress: c.ClientIP(),
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// createRefreshToken creates a new refresh token for a user (7 days expiry)
func (h *AuthHandler) createRefreshToken(userID string) (string, error) {
	token, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", err
	}

	tokenID := uuid.New().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days

	_, err = h.db.Exec(`
		INSERT INTO refresh_tokens (id, user_id, token, expires_at)
		VALUES (?, ?, ?, ?)
	`, tokenID, userID, token, expiresAt)

	if err != nil {
		return "", err
	}

	return token, nil
}
