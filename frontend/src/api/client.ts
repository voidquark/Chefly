import axios, { type AxiosInstance, type AxiosError } from 'axios';
import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  Recipe,
  RecipeGenerationRequest,
  RecipesResponse,
  FilterOptions,
  User,
  ShoppingListResponse,
  AdminStats,
  UsersResponse,
} from '../types';

// Use empty string for relative URLs (will use current host via Vite proxy in dev, or same host in production)
const API_BASE_URL = import.meta.env.VITE_API_URL || '';

class APIClient {
  private client: AxiosInstance;
  private refreshPromise: Promise<string> | null = null;
  private retryClient: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: `${API_BASE_URL}/api`,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Create a separate axios instance for retries (without interceptors)
    this.retryClient = axios.create({
      baseURL: `${API_BASE_URL}/api`,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // Add request interceptor to include auth token
    this.client.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem('access_token');
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );

    // Add response interceptor for error handling and token refresh
    this.client.interceptors.response.use(
      (response) => response,
      async (error: AxiosError) => {
        const originalRequest = error.config as any;

        // If 401 and we haven't tried refreshing yet for this specific request
        if (error.response?.status === 401 && !originalRequest._isRetry) {
          // Mark this request as being retried to prevent infinite loops
          originalRequest._isRetry = true;

          // Get or create the refresh promise
          if (!this.refreshPromise) {
            const refreshToken = localStorage.getItem('refresh_token');

            if (!refreshToken) {
              // No refresh token, redirect to login
              this.clearAuthAndRedirect();
              return Promise.reject(error);
            }

            // Create a promise that all waiting requests can share
            this.refreshPromise = (async () => {
              try {
                // Try to refresh the token using a separate axios instance (no interceptors)
                const response = await axios.post(`${API_BASE_URL}/api/auth/refresh`, {
                  refresh_token: refreshToken,
                });

                const { access_token, refresh_token: newRefreshToken } = response.data;

                // Store new tokens
                localStorage.setItem('access_token', access_token);
                localStorage.setItem('refresh_token', newRefreshToken);

                return access_token;
              } catch (refreshError) {
                // Refresh failed, clear auth and redirect
                this.clearAuthAndRedirect();
                throw refreshError;
              } finally {
                // Reset refresh state after a small delay to allow all pending requests to grab the promise
                setTimeout(() => {
                  this.refreshPromise = null;
                }, 100);
              }
            })();
          }

          try {
            // Wait for the token refresh to complete
            const newToken = await this.refreshPromise;

            // Update the original request with new token
            originalRequest.headers.Authorization = `Bearer ${newToken}`;

            // Retry the original request using the retry client (bypasses interceptors completely)
            return this.retryClient.request(originalRequest);
          } catch (refreshError) {
            // Refresh failed
            return Promise.reject(refreshError);
          }
        }

        return Promise.reject(error);
      }
    );
  }

  private clearAuthAndRedirect() {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user');
    window.location.href = '/login';
  }

  // Authentication
  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await this.client.post<AuthResponse>('/auth/register', data);
    return response.data;
  }

  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await this.client.post<AuthResponse>('/auth/login', data);
    return response.data;
  }

  async logout(): Promise<void> {
    const refreshToken = localStorage.getItem('refresh_token');
    if (refreshToken) {
      try {
        await this.client.post('/auth/logout', { refresh_token: refreshToken });
      } catch (error) {
        // Ignore errors on logout - user is logging out anyway
        console.error('Logout error:', error);
      }
    }
  }

  // User Profile
  async getProfile(): Promise<User> {
    const response = await this.client.get<User>('/user/profile');
    return response.data;
  }

  async updateProfile(data: Partial<User>): Promise<User> {
    const response = await this.client.put<User>('/user/profile', data);
    return response.data;
  }

  async getUserStats(): Promise<any> {
    const response = await this.client.get('/user/stats');
    return response.data;
  }

  // Recipe Generation
  async generateRecipe(data: RecipeGenerationRequest): Promise<Recipe> {
    const response = await this.client.post<Recipe>('/recipes/generate', data);
    return response.data;
  }

  // Recipe Management
  async getRecipes(): Promise<RecipesResponse> {
    const response = await this.client.get<RecipesResponse>('/recipes');
    return response.data;
  }

  async getRecipe(id: string): Promise<Recipe> {
    const response = await this.client.get<Recipe>(`/recipes/${id}`);
    return response.data;
  }

  async deleteRecipe(id: string): Promise<void> {
    await this.client.delete(`/recipes/${id}`);
  }

  async toggleFavorite(id: string): Promise<void> {
    await this.client.post(`/recipes/${id}/favorite`);
  }

  // Filter Options
  async getCountries(): Promise<string[]> {
    const response = await this.client.get<{ countries: string[] }>('/filters/countries');
    return response.data.countries;
  }

  async getMeatTypes(): Promise<string[]> {
    const response = await this.client.get<{ meat_types: string[] }>('/filters/meats');
    return response.data.meat_types;
  }

  async getIngredients(): Promise<string[]> {
    const response = await this.client.get<{ ingredients: string[] }>('/filters/ingredients');
    return response.data.ingredients;
  }

  async getFilterOptions(): Promise<FilterOptions> {
    const [countries, meat_types, ingredients] = await Promise.all([
      this.getCountries(),
      this.getMeatTypes(),
      this.getIngredients(),
    ]);
    return { countries, meat_types, ingredients };
  }

  // Shopping List
  async getShoppingList(): Promise<ShoppingListResponse> {
    const response = await this.client.get<ShoppingListResponse>('/shopping-list');
    return response.data;
  }

  async addRecipeToShoppingList(recipeId: string): Promise<void> {
    await this.client.post('/shopping-list/add-recipe', { recipe_id: recipeId });
  }

  async toggleShoppingItem(id: string): Promise<void> {
    await this.client.post(`/shopping-list/${id}/toggle`);
  }

  async deleteShoppingItem(id: string): Promise<void> {
    await this.client.delete(`/shopping-list/${id}`);
  }

  async clearCheckedItems(): Promise<void> {
    await this.client.delete('/shopping-list/clear/checked');
  }

  async clearAllShoppingItems(): Promise<void> {
    await this.client.delete('/shopping-list/clear/all');
  }

  // Admin endpoints
  async getAdminStats(): Promise<AdminStats> {
    const response = await this.client.get<AdminStats>('/admin/stats');
    return response.data;
  }

  async getAllUsers(): Promise<UsersResponse> {
    const response = await this.client.get<UsersResponse>('/admin/users');
    return response.data;
  }

  async deleteUser(userId: string): Promise<void> {
    await this.client.delete(`/admin/users/${userId}`);
  }

  async updateUserRecipeLimit(userId: string, recipeLimit: number | null): Promise<void> {
    await this.client.put(`/admin/users/${userId}/recipe-limit`, { recipe_limit: recipeLimit });
  }
}

export const apiClient = new APIClient();
