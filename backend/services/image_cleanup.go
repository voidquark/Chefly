package services

import (
	"os"
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
// Always logs cleanup attempts for audit trail
func (s *ImageCleanupService) DeleteRecipeImages(imagePath, thumbnailPath, recipeID, userID, requestID string) error {
	deletedFiles := []string{}
	failedFiles := []string{}
	var errors []string

	// Delete full image
	if imagePath != "" && !strings.HasPrefix(imagePath, "data:") {
		fullPath := s.getAbsolutePath(imagePath)
		if err := os.Remove(fullPath); err != nil {
			failedFiles = append(failedFiles, imagePath)
			errors = append(errors, err.Error())

			// Always log failure (even if file doesn't exist - could be wrong path)
			if s.auditLogger != nil {
				s.auditLogger.Warn("recipe.image_cleanup_failed", "Failed to delete recipe image file", &models.AuditContext{
					RequestID: requestID,
					UserID:    userID,
					Metadata: map[string]interface{}{
						"recipe_id":       recipeID,
						"image_url":       imagePath,
						"filesystem_path": fullPath,
						"error":           err.Error(),
						"is_not_exist":    os.IsNotExist(err),
					},
				})
			}
		} else {
			deletedFiles = append(deletedFiles, imagePath)
		}
	}

	// Delete thumbnail
	if thumbnailPath != "" && !strings.HasPrefix(thumbnailPath, "data:") {
		thumbPath := s.getAbsolutePath(thumbnailPath)
		if err := os.Remove(thumbPath); err != nil {
			failedFiles = append(failedFiles, thumbnailPath)
			errors = append(errors, err.Error())

			// Always log failure
			if s.auditLogger != nil {
				s.auditLogger.Warn("recipe.image_cleanup_failed", "Failed to delete recipe thumbnail file", &models.AuditContext{
					RequestID: requestID,
					UserID:    userID,
					Metadata: map[string]interface{}{
						"recipe_id":       recipeID,
						"thumbnail_url":   thumbnailPath,
						"filesystem_path": thumbPath,
						"error":           err.Error(),
						"is_not_exist":    os.IsNotExist(err),
					},
				})
			}
		} else {
			deletedFiles = append(deletedFiles, thumbnailPath)
		}
	}

	// Always log cleanup attempt (success, partial, or total failure)
	if s.auditLogger != nil {
		totalAttempted := 0
		if imagePath != "" && !strings.HasPrefix(imagePath, "data:") {
			totalAttempted++
		}
		if thumbnailPath != "" && !strings.HasPrefix(thumbnailPath, "data:") {
			totalAttempted++
		}

		if len(failedFiles) > 0 && len(deletedFiles) == 0 {
			// Total failure - all files failed to delete
			s.auditLogger.Error("recipe.image_cleanup_total_failure", "All recipe images failed to delete", nil, &models.AuditContext{
				RequestID: requestID,
				UserID:    userID,
				Metadata: map[string]interface{}{
					"recipe_id":       recipeID,
					"attempted_count": totalAttempted,
					"failed_count":    len(failedFiles),
					"failed_files":    failedFiles,
					"errors":          errors,
				},
			})
		} else if len(deletedFiles) > 0 && len(failedFiles) > 0 {
			// Partial success
			s.auditLogger.Warn("recipe.image_cleanup_partial", "Some recipe images failed to delete", &models.AuditContext{
				RequestID: requestID,
				UserID:    userID,
				Metadata: map[string]interface{}{
					"recipe_id":       recipeID,
					"attempted_count": totalAttempted,
					"deleted_count":   len(deletedFiles),
					"deleted_files":   deletedFiles,
					"failed_count":    len(failedFiles),
					"failed_files":    failedFiles,
					"cleanup_success": false,
				},
			})
		} else if len(deletedFiles) > 0 {
			// Complete success
			s.auditLogger.Info("recipe.image_cleanup", "Recipe images deleted from disk", &models.AuditContext{
				RequestID: requestID,
				UserID:    userID,
				Metadata: map[string]interface{}{
					"recipe_id":       recipeID,
					"deleted_count":   len(deletedFiles),
					"deleted_files":   deletedFiles,
					"cleanup_success": true,
				},
			})
		}
		// Note: If totalAttempted == 0 (no images to delete), we don't log anything
	}

	return nil
}

// DeleteUserImages deletes all image files for a user's recipes
// Takes a list of recipe image paths and deletes them in bulk
func (s *ImageCleanupService) DeleteUserImages(recipes []struct{ ImagePath, ThumbnailPath, RecipeID string }, userID, requestID string) error {
	totalAttempted := 0
	totalDeleted := 0
	totalFailed := 0
	deletedFiles := []string{}
	failedFiles := []string{}

	for _, recipe := range recipes {
		// Delete full image
		if recipe.ImagePath != "" && !strings.HasPrefix(recipe.ImagePath, "data:") {
			totalAttempted++
			fullPath := s.getAbsolutePath(recipe.ImagePath)
			if err := os.Remove(fullPath); err != nil {
				totalFailed++
				failedFiles = append(failedFiles, recipe.ImagePath)
			} else {
				totalDeleted++
				deletedFiles = append(deletedFiles, recipe.ImagePath)
			}
		}

		// Delete thumbnail
		if recipe.ThumbnailPath != "" && !strings.HasPrefix(recipe.ThumbnailPath, "data:") {
			totalAttempted++
			thumbPath := s.getAbsolutePath(recipe.ThumbnailPath)
			if err := os.Remove(thumbPath); err != nil {
				totalFailed++
				failedFiles = append(failedFiles, recipe.ThumbnailPath)
			} else {
				totalDeleted++
				deletedFiles = append(deletedFiles, recipe.ThumbnailPath)
			}
		}
	}

	// Always log cleanup results for user deletion
	if s.auditLogger != nil && totalAttempted > 0 {
		if totalFailed > 0 && totalDeleted == 0 {
			// Total failure
			s.auditLogger.Error("admin.user_images_cleanup_failed", "All user images failed to delete", nil, &models.AuditContext{
				RequestID: requestID,
				UserID:    userID,
				Metadata: map[string]interface{}{
					"deleted_user_id": userID,
					"recipe_count":    len(recipes),
					"attempted_count": totalAttempted,
					"failed_count":    totalFailed,
					"failed_files":    failedFiles,
					"cleanup_success": false,
				},
			})
		} else if totalFailed > 0 {
			// Partial success
			s.auditLogger.Warn("admin.user_images_cleanup_partial", "Some user images failed to delete", &models.AuditContext{
				RequestID: requestID,
				UserID:    userID,
				Metadata: map[string]interface{}{
					"deleted_user_id": userID,
					"recipe_count":    len(recipes),
					"attempted_count": totalAttempted,
					"deleted_count":   totalDeleted,
					"failed_count":    totalFailed,
					"failed_files":    failedFiles,
					"cleanup_success": false,
				},
			})
		} else if totalDeleted > 0 {
			// Complete success
			s.auditLogger.Info("admin.user_images_cleanup", "All user images deleted from disk", &models.AuditContext{
				RequestID: requestID,
				UserID:    userID,
				Metadata: map[string]interface{}{
					"deleted_user_id": userID,
					"recipe_count":    len(recipes),
					"attempted_count": totalAttempted,
					"deleted_count":   totalDeleted,
					"cleanup_success": true,
				},
			})
		}
	}

	return nil
}

// getAbsolutePath converts a URL path to absolute filesystem path
// URL format: /uploads/images/full/abc.jpg → ./uploads/images/full/abc.jpg
func (s *ImageCleanupService) getAbsolutePath(urlPath string) string {
	// Simply prepend "." to convert URL path to relative filesystem path
	// This works because the URL already includes the full path structure
	return "." + urlPath
}
