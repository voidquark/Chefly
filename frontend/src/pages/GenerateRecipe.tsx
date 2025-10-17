import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useLanguage } from '../contexts/LanguageContext';
import { usePageTitle } from '../hooks/usePageTitle';
import { apiClient } from '../api/client';
import { Select } from '../components/Select';
import { MultiSelect } from '../components/MultiSelect';
import { Loader2 } from 'lucide-react';
import type { Recipe } from '../types';

export const GenerateRecipe: React.FC = () => {
  const navigate = useNavigate();
  const { t, language } = useLanguage();
  usePageTitle('Generate Recipe');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [generatedRecipe, setGeneratedRecipe] = useState<Recipe | null>(null);

  // Filter options (English values from backend)
  const [meatTypes, setMeatTypes] = useState<string[]>([]);
  const [cuisines, setCuisines] = useState<string[]>([]);
  const [ingredients, setIngredients] = useState<string[]>([]);

  // Form state (always store English values for API)
  const [meatType, setMeatType] = useState('');
  const [cuisineType, setCuisineType] = useState('');
  const [sideIngredients, setSideIngredients] = useState<string[]>([]);
  const [dietaryPreferences, setDietaryPreferences] = useState<string[]>([]);
  const [cookingTime, setCookingTime] = useState('');
  const [difficulty, setDifficulty] = useState('');

  const dietaryOptions = ['Vegetarian', 'Vegan', 'Gluten-free', 'Dairy-free', 'Low-carb', 'Keto'];
  const cookingTimeOptions = ['quick', 'medium', 'long'];
  const difficultyOptions = ['easy', 'medium', 'hard'];

  // Helper function to create {value, label} objects for translated options
  const createTranslatedOptions = (options: string[], translationMap: Record<string, string>) => {
    return options.map(opt => ({
      value: opt,
      label: translationMap[opt] || opt
    }));
  };

  // Create translated option arrays
  const translatedMeatTypes = createTranslatedOptions(meatTypes, t.options.meats);
  const translatedCuisines = createTranslatedOptions(cuisines, t.options.cuisines);
  const translatedIngredients = createTranslatedOptions(ingredients, t.options.ingredients);
  const translatedDietary = createTranslatedOptions(dietaryOptions, t.options.dietary);
  const translatedCookingTimes = createTranslatedOptions(cookingTimeOptions, t.options.cookingTimes);
  const translatedDifficulties = createTranslatedOptions(difficultyOptions, t.options.difficulties);

  // Load filter options on mount
  useEffect(() => {
    const loadOptions = async () => {
      try {
        const [meats, countries, sides] = await Promise.all([
          apiClient.getMeatTypes(),
          apiClient.getCountries(),
          apiClient.getIngredients(),
        ]);
        setMeatTypes(meats);
        setCuisines(countries);
        setIngredients(sides);
      } catch (err) {
        console.error('Failed to load filter options:', err);
      }
    };
    loadOptions();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    setGeneratedRecipe(null);

    try {
      const recipe = await apiClient.generateRecipe({
        meat_type: meatType,
        side_ingredients: sideIngredients,
        cuisine_type: cuisineType,
        dietary_preferences: dietaryPreferences,
        cooking_time: cookingTime,
        difficulty: difficulty,
        language: language,
      });
      setGeneratedRecipe(recipe);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to generate recipe. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const viewRecipe = () => {
    if (generatedRecipe) {
      navigate(`/recipe/${generatedRecipe.id}`);
    }
  };

  return (
    <div className="max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-6">{t.generate.title}</h1>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-6">
          {error}
        </div>
      )}

      {generatedRecipe && (
        <div className="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 px-4 py-3 rounded mb-6">
          <div className="flex items-center justify-between">
            <div>
              <p className="font-semibold text-green-900 dark:text-green-200">{t.generate.successMessage}</p>
              <p className="text-sm text-green-700 dark:text-green-300">{generatedRecipe.title}</p>
            </div>
            <button
              onClick={viewRecipe}
              className="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded transition-colors"
            >
              {t.generate.viewRecipe}
            </button>
          </div>
        </div>
      )}

      <form onSubmit={handleSubmit} className="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {/* Meat Type */}
          <Select
            id="meatType"
            label={t.generate.meatType}
            options={translatedMeatTypes}
            value={meatType}
            onChange={(e) => setMeatType(e.target.value)}
            placeholder={`${t.common.select} ${t.generate.meatType}`}
            required
          />

          {/* Cuisine Type */}
          <Select
            id="cuisineType"
            label={t.generate.cuisine}
            options={translatedCuisines}
            value={cuisineType}
            onChange={(e) => setCuisineType(e.target.value)}
            placeholder={`${t.common.select} ${t.generate.cuisine}`}
            required
          />

          {/* Cooking Time */}
          <Select
            id="cookingTime"
            label={t.generate.cookingTime}
            options={translatedCookingTimes}
            value={cookingTime}
            onChange={(e) => setCookingTime(e.target.value)}
            placeholder={`${t.common.select} ${t.generate.cookingTime}`}
            required
          />

          {/* Difficulty */}
          <Select
            id="difficulty"
            label={t.generate.difficulty}
            options={translatedDifficulties}
            value={difficulty}
            onChange={(e) => setDifficulty(e.target.value)}
            placeholder={`${t.common.select} ${t.generate.difficulty}`}
            required
          />
        </div>

        {/* Side Ingredients */}
        <MultiSelect
          label={t.generate.sideIngredients}
          options={translatedIngredients}
          selected={sideIngredients}
          onChange={setSideIngredients}
        />

        {/* Dietary Preferences */}
        <MultiSelect
          label={t.generate.dietaryPreferences}
          options={translatedDietary}
          selected={dietaryPreferences}
          onChange={setDietaryPreferences}
        />

        <button
          type="submit"
          disabled={loading}
          className="w-full bg-primary-600 hover:bg-primary-700 text-white font-bold py-3 px-4 rounded focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors flex items-center justify-center"
        >
          {loading ? (
            <>
              <Loader2 className="animate-spin h-5 w-5 mr-2" />
              {t.generate.generating}
            </>
          ) : (
            t.generate.generateButton
          )}
        </button>

        {loading && (
          <p className="text-center text-gray-600 dark:text-gray-400 mt-4 text-sm">
            {t.generate.loadingMessage}
          </p>
        )}
      </form>
    </div>
  );
};
