import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { useLanguage } from '../contexts/LanguageContext';
import { usePageTitle } from '../hooks/usePageTitle';
import { Clock, ChefHat } from 'lucide-react';
import axios from 'axios';
import type { Recipe } from '../types';

export const SharedRecipe: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const { t } = useLanguage();
  const [recipe, setRecipe] = useState<Recipe | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  usePageTitle(recipe?.title ? `${recipe.title} (Shared)` : 'Shared Recipe');

  useEffect(() => {
    loadRecipe();
  }, [id]);

  const loadRecipe = async () => {
    if (!id) return;

    try {
      setLoading(true);
      const API_BASE_URL = import.meta.env.VITE_API_URL || '';
      const response = await axios.get(`${API_BASE_URL}/api/recipes/shared/${id}`);
      setRecipe(response.data);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load recipe');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-gray-50">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      </div>
    );
  }

  if (error || !recipe) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center p-4">
        <div className="max-w-md w-full bg-white rounded-lg shadow-lg p-8 text-center">
          <h1 className="text-2xl font-bold text-gray-900 mb-4">{t.recipe.notFound}</h1>
          <p className="text-gray-600 mb-6">{error || t.recipe.notFound}</p>
          <a
            href="/login"
            className="inline-block bg-primary-600 hover:bg-primary-700 text-white px-6 py-3 rounded-lg transition-colors"
          >
            {t.auth.login}
          </a>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 py-8">
      <div className="max-w-4xl mx-auto px-4">
        {/* Header Banner */}
        <div className="bg-gradient-to-r from-primary-600 to-primary-700 dark:from-primary-700 dark:to-primary-800 rounded-lg shadow-lg p-6 mb-6 text-white">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm opacity-90 mb-1">{t.recipe.sharedRecipe}</p>
              <h1 className="text-3xl font-bold">{t.nav.appName}</h1>
            </div>
            <a
              href="/register"
              className="bg-white text-primary-600 px-6 py-2 rounded-lg font-medium hover:bg-gray-100 transition-colors"
            >
              {t.auth.register}
            </a>
          </div>
        </div>

        {/* Recipe Content */}
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-hidden">
          <img
            src={recipe.image_path}
            alt={recipe.title}
            loading="lazy"
            className="w-full h-64 object-cover"
            onError={(e) => {
              (e.target as HTMLImageElement).src = 'https://via.placeholder.com/800x400/FF6B6B/FFFFFF?text=' + encodeURIComponent(recipe.title);
            }}
          />

          <div className="p-8">
            <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-4">{recipe.title}</h1>
            <p className="text-gray-600 dark:text-gray-300 mb-6">{recipe.description}</p>

            <div className="flex flex-wrap gap-4 mb-8">
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
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-4">{t.recipe.ingredients}</h2>
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
            <div className="mb-8">
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
              <div className="flex flex-wrap gap-2 mb-8">
                {recipe.dietary_tags.map((tag) => (
                  <span key={tag} className="px-3 py-1 bg-green-100 dark:bg-green-900 text-green-700 dark:text-green-200 rounded-full text-sm">
                    {tag}
                  </span>
                ))}
              </div>
            )}

            {/* Call to Action */}
            <div className="bg-gradient-to-r from-primary-50 to-blue-50 dark:from-primary-900/20 dark:to-blue-900/20 rounded-lg p-6 text-center">
              <h3 className="text-xl font-bold text-gray-900 dark:text-white mb-2">Love this recipe?</h3>
              <p className="text-gray-600 dark:text-gray-300 mb-4">Create your own AI-powered recipes with {t.nav.appName}!</p>
              <a
                href="/register"
                className="inline-block bg-primary-600 hover:bg-primary-700 text-white px-8 py-3 rounded-lg font-medium transition-colors"
              >
                {t.auth.register}
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
