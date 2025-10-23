import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Clock, ChefHat, Heart } from 'lucide-react';
import type { RecipeSummary } from '../types';

interface RecipeCardProps {
  recipe: RecipeSummary;
  onToggleFavorite?: (id: string) => void;
}

export const RecipeCard: React.FC<RecipeCardProps> = ({ recipe, onToggleFavorite }) => {
  const navigate = useNavigate();
  const [imageLoaded, setImageLoaded] = useState(false);

  // Use thumbnail if available, otherwise fallback to full image
  const thumbnailSrc = recipe.thumbnail_path || recipe.image_path;
  const fullImageSrc = recipe.image_path;

  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden hover:shadow-xl hover:scale-105 transition-all duration-300 cursor-pointer animate-fadeIn">
      <div onClick={() => navigate(`/recipe/${recipe.id}`)}>
        <div className="relative w-full h-48 bg-gray-200 dark:bg-gray-700">
          {/* Thumbnail - loads immediately, blurred */}
          {thumbnailSrc && (
            <img
              src={thumbnailSrc}
              alt=""
              className={`absolute inset-0 w-full h-full object-cover transition-opacity duration-300 ${
                imageLoaded ? 'opacity-0' : 'opacity-100 blur-sm'
              }`}
            />
          )}

          {/* Full image - loads lazily */}
          <img
            src={fullImageSrc}
            alt={recipe.title}
            loading="lazy"
            className={`absolute inset-0 w-full h-full object-cover transition-opacity duration-300 ${
              imageLoaded ? 'opacity-100' : 'opacity-0'
            }`}
            onLoad={() => setImageLoaded(true)}
            onError={(e) => {
              const target = e.target as HTMLImageElement;
              // Prevent infinite loop - only set fallback once
              if (!target.dataset.fallbackSet) {
                target.dataset.fallbackSet = 'true';
                // Use inline SVG with emoji instead of external placeholder
                target.src = `data:image/svg+xml,${encodeURIComponent(`
                  <svg xmlns="http://www.w3.org/2000/svg" width="400" height="300" viewBox="0 0 400 300">
                    <rect width="400" height="300" fill="#f3f4f6"/>
                    <text x="50%" y="50%" font-size="80" text-anchor="middle" dy=".3em">🍽️</text>
                    <text x="50%" y="70%" font-size="16" text-anchor="middle" fill="#6b7280">${recipe.title}</text>
                  </svg>
                `)}`;
              }
              setImageLoaded(true);
            }}
          />
        </div>
        <div className="p-4">
          <div className="flex justify-between items-start mb-2">
            <h3 className="text-xl font-bold text-gray-900 dark:text-white line-clamp-2">{recipe.title}</h3>
            {onToggleFavorite && (
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  onToggleFavorite(recipe.id);
                }}
                className="flex-shrink-0 ml-2"
              >
                <Heart
                  className={`h-6 w-6 ${
                    recipe.is_favorite ? 'fill-red-500 text-red-500' : 'text-gray-400'
                  } hover:text-red-500 transition-colors`}
                />
              </button>
            )}
          </div>

          <p className="text-gray-600 dark:text-gray-300 text-sm mb-3 line-clamp-2">{recipe.description}</p>

          <div className="flex items-center justify-between text-sm text-gray-500 dark:text-gray-400">
            <div className="flex items-center">
              <Clock className="h-4 w-4 mr-1" />
              <span>{recipe.cooking_time} min</span>
            </div>
            <div className="flex items-center">
              <ChefHat className="h-4 w-4 mr-1" />
              <span className="capitalize">{recipe.difficulty}</span>
            </div>
            <div className="px-2 py-1 bg-primary-100 dark:bg-primary-900 text-primary-700 dark:text-primary-200 rounded text-xs font-medium">
              {recipe.cuisine_type}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
