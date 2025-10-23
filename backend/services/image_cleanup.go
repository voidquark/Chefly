package services

import (
	"os"
	"path/filepath"
	"strings"

	"chefly/models"
)

// ImageCleanupService handles cleanup of recipe image files
type ImageCleanupService struct {
	uploadsDir  string
	auditLogger *AuditLogger
}

// NewImageCleanupService creates a new image cleanup service
func NewImageCleanupService(uploadsDir string, auditLogger *AuditLogger) *ImageCleanupService {
	return &ImageCleanupService{
		uploadsDir:  uploadsDir,
		auditLogger: auditLogger,
	}
}

// DeleteRecipeImages deletes the full image and thumbnail for a recipe
// Returns error only for critical failures; missing files are logged as warnings
func (s *ImageCleanupService) DeleteRecipeImages(imagePath, thumbnailPath, recipeID, userID, requestID string) error {
	deletedFiles := []string{}
	failedFiles := []string{}

	// Delete full image
	if imagePath != "" && !strings.HasPrefix(imagePath, "data:") {
		fullPath := s.getAbsolutePath(imagePath)
		if err := os.Remove(fullPath); err != nil {
			if !os.IsNotExist(err) {
				failedFiles = append(failedFiles, imagePath)
				if s.auditLogger != nil {
					s.auditLogger.Warn("recipe.image_cleanup_failed", "Failed to delete recipe image file", &models.AuditContext{
						RequestID: requestID,
						UserID:    userID,
						Metadata: map[string]interface{}{
							"recipe_id":  recipeID,
							"image_path": imagePath,
							"error":      err.Error(),
						},
					})
				}
			}
			// File doesn't exist - not an error, just skip
		} else {
			deletedFiles = append(deletedFiles, imagePath)
		}
	}

	// Delete thumbnail
	if thumbnailPath != "" && !strings.HasPrefix(thumbnailPath, "data:") {
		thumbPath := s.getAbsolutePath(thumbnailPath)
		if err := os.Remove(thumbPath); err != nil {
			if !os.IsNotExist(err) {
				failedFiles = append(failedFiles, thumbnailPath)
				if s.auditLogger != nil {
					s.auditLogger.Warn("recipe.image_cleanup_failed", "Failed to delete recipe thumbnail file", &models.AuditContext{
						RequestID: requestID,
						UserID:    userID,
						Metadata: map[string]interface{}{
							"recipe_id":      recipeID,
							"thumbnail_path": thumbnailPath,
							"error":          err.Error(),
						},
					})
				}
			}
			// File doesn't exist - not an error, just skip
		} else {
			deletedFiles = append(deletedFiles, thumbnailPath)
		}
	}

	// Log successful cleanup if any files were deleted
	if len(deletedFiles) > 0 && s.auditLogger != nil {
		s.auditLogger.Info("recipe.image_cleanup", "Recipe images deleted from disk", &models.AuditContext{
			RequestID: requestID,
			UserID:    userID,
			Metadata: map[string]interface{}{
				"recipe_id":      recipeID,
				"deleted_files":  deletedFiles,
				"deleted_count":  len(deletedFiles),
				"failed_files":   failedFiles,
				"cleanup_success": len(failedFiles) == 0,
			},
		})
	}

	return nil
}

// DeleteUserImages deletes all image files for a user's recipes
// Takes a list of recipe image paths and deletes them in bulk
func (s *ImageCleanupService) DeleteUserImages(recipes []struct{ ImagePath, ThumbnailPath, RecipeID string }, userID, requestID string) error {
	totalDeleted := 0
	totalFailed := 0
	deletedFiles := []string{}

	for _, recipe := range recipes {
		// Delete full image
		if recipe.ImagePath != "" && !strings.HasPrefix(recipe.ImagePath, "data:") {
			fullPath := s.getAbsolutePath(recipe.ImagePath)
			if err := os.Remove(fullPath); err != nil {
				if !os.IsNotExist(err) {
					totalFailed++
				}
			} else {
				totalDeleted++
				deletedFiles = append(deletedFiles, recipe.ImagePath)
			}
		}

		// Delete thumbnail
		if recipe.ThumbnailPath != "" && !strings.HasPrefix(recipe.ThumbnailPath, "data:") {
			thumbPath := s.getAbsolutePath(recipe.ThumbnailPath)
			if err := os.Remove(thumbPath); err != nil {
				if !os.IsNotExist(err) {
					totalFailed++
				}
			} else {
				totalDeleted++
				deletedFiles = append(deletedFiles, recipe.ThumbnailPath)
			}
		}
	}

	// Log cleanup results
	if s.auditLogger != nil {
		if totalFailed > 0 {
			s.auditLogger.Warn("admin.user_images_cleanup", "User images cleanup completed with failures", &models.AuditContext{
				RequestID: requestID,
				UserID:    userID,
				Metadata: map[string]interface{}{
					"deleted_user_id":   userID,
					"recipe_count":      len(recipes),
					"files_deleted":     totalDeleted,
					"files_failed":      totalFailed,
					"cleanup_success":   totalFailed == 0,
				},
			})
		} else if totalDeleted > 0 {
			s.auditLogger.Info("admin.user_images_cleanup", "All user images deleted from disk", &models.AuditContext{
				RequestID: requestID,
				UserID:    userID,
				Metadata: map[string]interface{}{
					"deleted_user_id": userID,
					"recipe_count":    len(recipes),
					"files_deleted":   totalDeleted,
					"cleanup_success": true,
				},
			})
		}
	}

	return nil
}

// getAbsolutePath converts a URL path to absolute filesystem path
// Handles paths like "/uploads/images/full/abc.jpg" -> "./uploads/images/full/abc.jpg"
func (s *ImageCleanupService) getAbsolutePath(urlPath string) string {
	// Remove leading slash if present
	cleanPath := strings.TrimPrefix(urlPath, "/")
	// Join with uploads directory
	return filepath.Join(s.uploadsDir, cleanPath)
}
