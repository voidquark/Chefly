import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Clock, ChefHat, Heart } from 'lucide-react';
import type { RecipeSummary } from '../types';

interface RecipeCardProps {
  recipe: RecipeSummary;
  onToggleFavorite?: (id: string) => void;
}

export const RecipeCard: React.FC<RecipeCardProps> = ({ recipe, onToggleFavorite }) => {
  const navigate = useNavigate();

  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden hover:shadow-xl hover:scale-105 transition-all duration-300 cursor-pointer animate-fadeIn">
      <div onClick={() => navigate(`/recipe/${recipe.id}`)}>
        <img
          src={recipe.image_path}
          alt={recipe.title}
          loading="lazy"
          className="w-full h-48 object-cover"
          onError={(e) => {
            (e.target as HTMLImageElement).src = 'https://via.placeholder.com/400x300/FF6B6B/FFFFFF?text=' + encodeURIComponent(recipe.title);
          }}
        />
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
