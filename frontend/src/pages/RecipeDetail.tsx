import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast';
import { useLanguage } from '../contexts/LanguageContext';
import { usePageTitle } from '../hooks/usePageTitle';
import { apiClient } from '../api/client';
import { SkeletonRecipeDetail } from '../components/SkeletonRecipeDetail';
import { ConfirmDialog } from '../components/ConfirmDialog';
import { Clock, ChefHat, Heart, ArrowLeft, Trash2, ShoppingBag, Share2 } from 'lucide-react';
import type { Recipe } from '../types';

export const RecipeDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { t } = useLanguage();
  const [recipe, setRecipe] = useState<Recipe | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [addingToList, setAddingToList] = useState(false);
  const [showShareModal, setShowShareModal] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [imageLoaded, setImageLoaded] = useState(false);
  usePageTitle(recipe?.title || 'Recipe Detail');

  useEffect(() => {
    loadRecipe();
  }, [id]);

  const loadRecipe = async () => {
    if (!id) return;

    try {
      setLoading(true);
      const data = await apiClient.getRecipe(id);
      setRecipe(data);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load recipe');
    } finally {
      setLoading(false);
    }
  };

  const handleToggleFavorite = async () => {
    if (!id || !recipe) return;

    try {
      await apiClient.toggleFavorite(id);
      setRecipe({ ...recipe, is_favorite: !recipe.is_favorite });
    } catch (err) {
      console.error('Failed to toggle favorite:', err);
    }
  };

  const handleDeleteClick = () => {
    setShowDeleteConfirm(true);
  };

  const handleDeleteConfirm = async () => {
    if (!id) return;

    try {
      await apiClient.deleteRecipe(id);
      toast.success('Recipe deleted successfully');
      setShowDeleteConfirm(false);
      navigate('/recipes');
    } catch (err) {
      toast.error(t.recipe.deleteFailed);
      setShowDeleteConfirm(false);
    }
  };

  const handleAddToShoppingList = async () => {
    if (!id) return;

    try {
      setAddingToList(true);
      await apiClient.addRecipeToShoppingList(id);
      toast.success(t.shoppingList.addedToList);
    } catch (err) {
      console.error('Failed to add to shopping list:', err);
      toast.error('Failed to add to shopping list');
    } finally {
      setAddingToList(false);
    }
  };

  const handleShareRecipe = () => {
    if (!id) return;
    setShowShareModal(true);
  };

  const copyToClipboard = async () => {
    if (!id) return;

    const shareUrl = `${window.location.origin}/shared/${id}`;

    try {
      // Try modern Clipboard API first
      if (navigator.clipboard && navigator.clipboard.writeText) {
        await navigator.clipboard.writeText(shareUrl);
        toast.success(t.recipe.linkCopied);
        return;
      }

      // Fallback for iOS and older browsers
      const textArea = document.createElement('textarea');
      textArea.value = shareUrl;
      textArea.style.position = 'fixed';
      textArea.style.left = '-999999px';
      textArea.style.top = '-999999px';
      document.body.appendChild(textArea);
      textArea.focus();
      textArea.select();

      try {
        document.execCommand('copy');
        textArea.remove();
        toast.success(t.recipe.linkCopied);
      } catch (err) {
        textArea.remove();
        toast.error('Failed to copy link');
      }
    } catch (err) {
      console.error('Failed to copy:', err);
      toast.error('Failed to copy link');
    }
  };

  if (loading) {
    return <SkeletonRecipeDetail />;
  }

  if (error || !recipe) {
    return (
      <div className="max-w-4xl mx-auto">
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
          {error || t.recipe.notFound}
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto">
      <button
        onClick={() => navigate('/recipes')}
        className="flex items-center text-gray-600 dark:text-gray-300 hover:text-primary-600 dark:hover:text-primary-400 mb-6 transition-colors"
      >
        <ArrowLeft className="h-5 w-5 mr-2" />
        {t.recipe.backToRecipes}
      </button>

      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-hidden">
        <div className="relative w-full h-64 bg-gray-200 dark:bg-gray-700">
          {/* Thumbnail - loads immediately, blurred */}
          {recipe.thumbnail_path && (
            <img
              src={recipe.thumbnail_path}
              alt=""
              className={`absolute inset-0 w-full h-full object-cover transition-opacity duration-300 ${
                imageLoaded ? 'opacity-0' : 'opacity-100 blur-sm'
              }`}
            />
          )}

          {/* Full image - loads lazily */}
          <img
            src={recipe.image_path}
            alt={recipe.title}
            loading="lazy"
            className={`absolute inset-0 w-full h-full object-cover transition-opacity duration-300 ${
              imageLoaded ? 'opacity-100' : 'opacity-0'
            }`}
            onLoad={() => setImageLoaded(true)}
            onError={(e) => {
              (e.target as HTMLImageElement).src = 'https://via.placeholder.com/800x400/FF6B6B/FFFFFF?text=' + encodeURIComponent(recipe.title);
              setImageLoaded(true);
            }}
          />
        </div>

        <div className="p-8">
          <div className="flex justify-between items-start mb-4">
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white">{recipe.title}</h1>
            <div className="flex space-x-2">
              <button
                onClick={handleShareRecipe}
                className="p-2 rounded-full hover:bg-blue-50 dark:hover:bg-blue-900 transition-colors"
                title={t.recipe.shareRecipe}
              >
                <Share2 className="h-6 w-6 text-blue-500 dark:text-blue-400" />
              </button>
              <button
                onClick={handleToggleFavorite}
                className="p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
              >
                <Heart
                  className={`h-6 w-6 ${
                    recipe.is_favorite ? 'fill-red-500 text-red-500' : 'text-gray-400'
                  }`}
                />
              </button>
              <button
                onClick={handleDeleteClick}
                className="p-2 rounded-full hover:bg-red-50 dark:hover:bg-red-900 transition-colors"
              >
                <Trash2 className="h-6 w-6 text-red-500 dark:text-red-400" />
              </button>
            </div>
          </div>

          <p className="text-gray-600 dark:text-gray-300 mb-6">{recipe.description}</p>

          <div className="flex flex-wrap gap-4 mb-6">
            <div className="flex items-center text-gray-700 dark:text-gray-300">
              <Clock className="h-5 w-5 mr-2 text-primary-600 dark:text-primary-400" />
              <span>{recipe.cooking_time} {t.recipe.minutes}</span>
            </div>
            <div className="flex items-center text-gray-700 dark:text-gray-300">
              <ChefHat className="h-5 w-5 mr-2 text-primary-600 dark:text-primary-400" />
              <span className="capitalize">{recipe.difficulty}</span>
            </div>
            <div className="px-3 py-1 bg-primary-100 dark:bg-primary-900 text-primary-700 dark:text-primary-200 rounded-full text-sm font-medium">
              {recipe.cuisine_type}
            </div>
            <div className="px-3 py-1 bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-200 rounded-full text-sm font-medium">
              {recipe.meat_type}
            </div>
          </div>

          {/* Ingredients */}
          <div className="mb-8">
            <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-3 mb-4">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white">{t.recipe.ingredients}</h2>
              <button
                onClick={handleAddToShoppingList}
                disabled={addingToList}
                className="flex items-center justify-center px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed whitespace-nowrap"
              >
                <ShoppingBag className="h-5 w-5 mr-2" />
                <span className="text-sm sm:text-base">{addingToList ? t.common.loading : t.shoppingList.addToShoppingList}</span>
              </button>
            </div>
            <ul className="space-y-2">
              {recipe.ingredients.map((ingredient, index) => (
                <li key={index} className="flex items-start">
                  <span className="text-primary-600 dark:text-primary-400 mr-2">•</span>
                  <span className="text-gray-700 dark:text-gray-300">
                    <span className="font-semibold">{ingredient.quantity} {ingredient.unit}</span> {ingredient.name}
                  </span>
                </li>
              ))}
            </ul>
          </div>

          {/* Steps */}
          <div>
            <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-4">{t.recipe.instructions}</h2>
            <div className="space-y-4">
              {recipe.steps.map((step) => (
                <div key={step.step_number} className="flex">
                  <div className="flex-shrink-0 w-8 h-8 bg-primary-600 dark:bg-primary-700 text-white rounded-full flex items-center justify-center font-bold mr-4">
                    {step.step_number}
                  </div>
                  <div className="flex-grow">
                    <p className="text-gray-700 dark:text-gray-300 mb-1">{step.instruction}</p>
                    {(step.timing || step.temperature) && (
                      <div className="flex gap-4 text-sm text-gray-500 dark:text-gray-400">
                        {step.timing && <span>⏱️ {step.timing}</span>}
                        {step.temperature && <span>🌡️ {step.temperature}</span>}
                      </div>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>

          {recipe.dietary_tags && recipe.dietary_tags.length > 0 && (
            <div className="mt-6 flex flex-wrap gap-2">
              {recipe.dietary_tags.map((tag) => (
                <span key={tag} className="px-3 py-1 bg-green-100 dark:bg-green-900 text-green-700 dark:text-green-200 rounded-full text-sm">
                  {tag}
                </span>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Share Modal */}
      {showShareModal && (
        <div
          className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
          onClick={() => setShowShareModal(false)}
        >
          <div
            className="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-md w-full p-6"
            onClick={(e) => e.stopPropagation()}
          >
            <h3 className="text-xl font-bold text-gray-900 dark:text-white mb-4">{t.recipe.shareRecipe}</h3>
            <p className="text-gray-600 dark:text-gray-300 mb-4">Copy this link to share the recipe:</p>

            <div className="bg-gray-50 dark:bg-gray-700 border border-gray-200 dark:border-gray-600 rounded-lg p-3 mb-4 break-all text-sm text-gray-900 dark:text-gray-100">
              {`${window.location.origin}/shared/${id}`}
            </div>

            <div className="flex gap-3">
              <button
                onClick={copyToClipboard}
                className="flex-1 bg-primary-600 hover:bg-primary-700 text-white px-4 py-2 rounded-lg transition-colors font-medium"
              >
                {t.recipe.copyLink}
              </button>
              <button
                onClick={() => setShowShareModal(false)}
                className="px-4 py-2 bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 text-gray-800 dark:text-gray-200 rounded-lg transition-colors font-medium"
              >
                {t.common.cancel}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Delete Confirmation Dialog */}
      <ConfirmDialog
        isOpen={showDeleteConfirm}
        title="Delete Recipe?"
        message="Are you sure you want to delete this recipe? This action cannot be undone."
        confirmLabel="Delete"
        cancelLabel="Cancel"
        variant="danger"
        onConfirm={handleDeleteConfirm}
        onCancel={() => setShowDeleteConfirm(false)}
      />
    </div>
  );
};
