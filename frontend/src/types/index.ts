// User types
export interface User {
  id: string;
  email: string;
  username: string;
  is_admin: boolean;
  created_at: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  user: User;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  username: string;
}

// Recipe types
export interface Ingredient {
  name: string;
  quantity: string;
  unit: string;
}

export interface CookingStep {
  step_number: number;
  instruction: string;
  timing?: string;
  temperature?: string;
}

export interface Recipe {
  id: string;
  user_id: string;
  title: string;
  description: string;
  ingredients: Ingredient[];
  steps: CookingStep[];
  cooking_time: number;
  difficulty: string;
  cuisine_type: string;
  meat_type: string;
  dietary_tags: string[];
  is_favorite: boolean;
  image_path: string;
  created_at: string;
}

export interface RecipeGenerationRequest {
  meat_type: string;
  side_ingredients: string[];
  cuisine_type: string;
  dietary_preferences: string[];
  cooking_time: string; // "quick", "medium", "long"
  difficulty: string; // "easy", "medium", "hard"
  language?: string; // "en" or "sk"
}

export interface RecipeSummary {
  id: string;
  title: string;
  description: string;
  cuisine_type: string;
  difficulty: string;
  cooking_time: number;
  is_favorite: boolean;
  image_path: string;
  created_at: string;
}

export interface RecipesResponse {
  recipes: RecipeSummary[];
}

// Filter options
export interface FilterOptions {
  countries?: string[];
  meat_types?: string[];
  ingredients?: string[];
}

// Shopping List types
export interface ShoppingListItem {
  id: string;
  user_id: string;
  recipe_id?: string;
  recipe_title?: string;
  ingredient_name: string;
  quantity: string;
  unit: string;
  is_checked: boolean;
  created_at: string;
}

export interface ShoppingListResponse {
  items: ShoppingListItem[];
}

// Admin types
export interface UserWithStats {
  id: string;
  email: string;
  username: string;
  is_admin: boolean;
  created_at: string;
  recipe_count: number;
  shopping_items: number;
  last_recipe_date?: string;
  recipe_limit?: number | null; // null = use global, -1 = unlimited, 0 = blocked, >0 = custom
}

export interface AdminStats {
  total_users: number;
  total_recipes: number;
  total_shopping_items: number;
  average_recipes_per_user: number;
  most_active_user?: string;
  most_active_user_count: number;
  recent_registrations: number;
  first_user_date: string;
}

export interface UsersResponse {
  users: UserWithStats[];
}
