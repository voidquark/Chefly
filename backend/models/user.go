package models

import "time"

// User represents a user in the system
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never expose password hash in JSON
	Username     string    `json:"username"`
	IsAdmin      bool      `json:"is_admin"`
	RecipeLimit  *int      `json:"recipe_limit,omitempty"` // NULL = use global, -1 = unlimited, 0 = blocked, >0 = custom
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserRegistration represents registration request
type UserRegistration struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Username string `json:"username" binding:"required,min=3"`
}

// UserLogin represents login request
type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse represents user data in responses (without sensitive info)
type UserResponse struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	IsAdmin     bool      `json:"is_admin"`
	RecipeLimit *int      `json:"recipe_limit,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		Username:    u.Username,
		IsAdmin:     u.IsAdmin,
		RecipeLimit: u.RecipeLimit,
		CreatedAt:   u.CreatedAt,
	}
}
