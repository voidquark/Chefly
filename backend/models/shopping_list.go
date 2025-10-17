package models

// ShoppingListItem represents a shopping list item
type ShoppingListItem struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	RecipeID       string `json:"recipe_id,omitempty"`
	RecipeTitle    string `json:"recipe_title,omitempty"`
	IngredientName string `json:"ingredient_name"`
	Quantity       string `json:"quantity"`
	Unit           string `json:"unit"`
	IsChecked      bool   `json:"is_checked"`
	CreatedAt      string `json:"created_at"`
}

// AddToShoppingListRequest represents a request to add ingredients to shopping list
type AddToShoppingListRequest struct {
	RecipeID string `json:"recipe_id" binding:"required"`
}
