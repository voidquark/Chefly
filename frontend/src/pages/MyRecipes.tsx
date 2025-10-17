import React, { useState, useEffect } from 'react';
import { useLanguage } from '../contexts/LanguageContext';
import { usePageTitle } from '../hooks/usePageTitle';
import { apiClient } from '../api/client';
import { RecipeCard } from '../components/RecipeCard';
import { SkeletonCard } from '../components/SkeletonCard';
import { Search, X, Heart } from 'lucide-react';
import type { RecipeSummary } from '../types';

export const MyRecipes: React.FC = () => {
  const { t } = useLanguage();
  usePageTitle('My Recipes');
  const [recipes, setRecipes] = useState<RecipeSummary[]>([]);
  const [filteredRecipes, setFilteredRecipes] = useState<RecipeSummary[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const [showFavoritesOnly, setShowFavoritesOnly] = useState(false);

  useEffect(() => {
    loadRecipes();
  }, []);

  const loadRecipes = async () => {
    try {
      setLoading(true);
      const data = await apiClient.getRecipes();
      setRecipes(data.recipes || []);
      setFilteredRecipes(data.recipes || []);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load recipes');
    } finally {
      setLoading(false);
    }
  };

  // Filter recipes based on search query and favorites toggle
  useEffect(() => {
    let filtered = recipes;

    // Apply favorites filter first
    if (showFavoritesOnly) {
      filtered = filtered.filter(recipe => recipe.is_favorite);
    }

    // Then apply search filter
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter((recipe) => {
        return (
          recipe.title.toLowerCase().includes(query) ||
          recipe.description?.toLowerCase().includes(query) ||
          recipe.cuisine_type?.toLowerCase().includes(query) ||
          recipe.difficulty?.toLowerCase().includes(query)
        );
      });
    }

    setFilteredRecipes(filtered);
  }, [searchQuery, recipes, showFavoritesOnly]);

  const handleToggleFavorite = async (id: string) => {
    try {
      await apiClient.toggleFavorite(id);
      // Reload recipes to update favorite status
      await loadRecipes();
    } catch (err) {
      console.error('Failed to toggle favorite:', err);
    }
  };

  if (loading) {
    return (
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-6">{t.myRecipes.title}</h1>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {[1, 2, 3, 4, 5, 6].map((i) => (
            <SkeletonCard key={i} />
          ))}
        </div>
      </div>
    );
  }

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-6">{t.myRecipes.title}</h1>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-6">
          {error}
        </div>
      )}

      {recipes.length === 0 ? (
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-8 text-center">
          <p className="text-gray-600 dark:text-gray-300 mb-4">{t.myRecipes.noRecipes}</p>
          <a
            href="/generate"
            className="inline-block bg-primary-600 hover:bg-primary-700 text-white px-6 py-2 rounded transition-colors"
          >
            {t.myRecipes.generateFirst}
          </a>
        </div>
      ) : (
        <>
          {/* Search Bar and Favorites Filter */}
          <div className="mb-6">
            <div className="flex flex-col sm:flex-row gap-3">
              <div className="relative flex-1 max-w-md">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <Search className="h-5 w-5 text-gray-400" />
                </div>
                <input
                  type="text"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  placeholder={t.myRecipes.search}
                  className="block w-full pl-10 pr-10 py-2 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-gray-200 rounded-lg focus:ring-primary-500 focus:border-primary-500"
                />
                {searchQuery && (
                  <button
                    onClick={() => setSearchQuery('')}
                    className="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-600"
                  >
                    <X className="h-5 w-5" />
                  </button>
                )}
              </div>
              <button
                onClick={() => setShowFavoritesOnly(!showFavoritesOnly)}
                className={`flex items-center justify-center px-4 py-2 rounded-lg font-medium transition-colors ${
                  showFavoritesOnly
                    ? 'bg-red-600 text-white hover:bg-red-700'
                    : 'bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-600'
                }`}
              >
                <Heart className={`h-5 w-5 mr-2 ${showFavoritesOnly ? 'fill-current' : ''}`} />
                {t.myRecipes.showFavorites}
              </button>
            </div>
            {(searchQuery || showFavoritesOnly) && (
              <p className="text-sm text-gray-600 dark:text-gray-400 mt-2">
                {filteredRecipes.length} {filteredRecipes.length === 1 ? 'recipe' : 'recipes'} found
                {showFavoritesOnly && ' (favorites only)'}
              </p>
            )}
          </div>

          {/* Recipe Grid */}
          {filteredRecipes.length === 0 ? (
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-8 text-center">
              <p className="text-gray-600 dark:text-gray-300 mb-4">{t.myRecipes.noResults}</p>
              <button
                onClick={() => setSearchQuery('')}
                className="inline-block bg-primary-600 hover:bg-primary-700 text-white px-6 py-2 rounded transition-colors"
              >
                {t.myRecipes.clearSearch}
              </button>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {filteredRecipes.map((recipe) => (
                <RecipeCard
                  key={recipe.id}
                  recipe={recipe}
                  onToggleFavorite={handleToggleFavorite}
                />
              ))}
            </div>
          )}
        </>
      )}
    </div>
  );
};
