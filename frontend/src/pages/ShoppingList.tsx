import React, { useState, useEffect } from 'react';
import { useLanguage } from '../contexts/LanguageContext';
import { usePageTitle } from '../hooks/usePageTitle';
import { apiClient } from '../api/client';
import { SkeletonShoppingList } from '../components/SkeletonShoppingList';
import { ShoppingBag, Trash2, CheckCircle2, Circle } from 'lucide-react';
import type { ShoppingListItem } from '../types';

export const ShoppingList: React.FC = () => {
  const { t } = useLanguage();
  usePageTitle('Shopping List');
  const [items, setItems] = useState<ShoppingListItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadShoppingList();
  }, []);

  const loadShoppingList = async () => {
    try {
      setLoading(true);
      const data = await apiClient.getShoppingList();
      setItems(data.items || []);
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to load shopping list');
    } finally {
      setLoading(false);
    }
  };

  const handleToggleItem = async (id: string) => {
    try {
      await apiClient.toggleShoppingItem(id);
      // Update local state
      setItems(items.map(item =>
        item.id === id ? { ...item, is_checked: !item.is_checked } : item
      ));
    } catch (err) {
      console.error('Failed to toggle item:', err);
    }
  };

  const handleDeleteItem = async (id: string) => {
    try {
      await apiClient.deleteShoppingItem(id);
      setItems(items.filter(item => item.id !== id));
    } catch (err) {
      console.error('Failed to delete item:', err);
    }
  };

  const handleClearChecked = async () => {
    try {
      await apiClient.clearCheckedItems();
      setItems(items.filter(item => !item.is_checked));
    } catch (err) {
      console.error('Failed to clear checked items:', err);
    }
  };

  const handleClearAll = async () => {
    if (window.confirm(t.shoppingList.confirmClearAll)) {
      try {
        await apiClient.clearAllShoppingItems();
        setItems([]);
      } catch (err) {
        console.error('Failed to clear all items:', err);
      }
    }
  };

  // Group items by recipe
  const groupedItems = items.reduce((groups, item) => {
    const key = item.recipe_title || 'Other';
    if (!groups[key]) {
      groups[key] = [];
    }
    groups[key].push(item);
    return groups;
  }, {} as Record<string, ShoppingListItem[]>);

  const checkedCount = items.filter(item => item.is_checked).length;

  if (loading) {
    return <SkeletonShoppingList />;
  }

  return (
    <div className="max-w-4xl mx-auto">
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center">
          <ShoppingBag className="h-8 w-8 text-primary-600 dark:text-primary-400 mr-3" />
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">{t.shoppingList.title}</h1>
        </div>
        {items.length > 0 && (
          <div className="flex gap-2">
            {checkedCount > 0 && (
              <button
                onClick={handleClearChecked}
                className="px-4 py-2 bg-orange-600 hover:bg-orange-700 text-white rounded-lg transition-colors text-sm"
              >
                {t.shoppingList.clearChecked} ({checkedCount})
              </button>
            )}
            <button
              onClick={handleClearAll}
              className="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors text-sm"
            >
              {t.shoppingList.clearAll}
            </button>
          </div>
        )}
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-6">
          {error}
        </div>
      )}

      {items.length === 0 ? (
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-12 text-center">
          <ShoppingBag className="h-16 w-16 text-gray-300 dark:text-gray-600 mx-auto mb-4" />
          <p className="text-gray-600 dark:text-gray-300 text-lg mb-2">{t.shoppingList.noItems}</p>
          <p className="text-gray-500 dark:text-gray-400">{t.shoppingList.addFromRecipe}</p>
          <a
            href="/recipes"
            className="inline-block mt-6 bg-primary-600 hover:bg-primary-700 text-white px-6 py-3 rounded-lg transition-colors"
          >
            {t.nav.myRecipes}
          </a>
        </div>
      ) : (
        <div className="space-y-6">
          {Object.entries(groupedItems).map(([recipeTitle, recipeItems]) => (
            <div key={recipeTitle} className="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden">
              <div className="bg-primary-50 dark:bg-primary-900/20 px-4 py-3 border-b border-primary-100 dark:border-primary-800">
                <h2 className="font-semibold text-primary-900 dark:text-primary-200">
                  {recipeTitle === 'Other' ? t.shoppingList.itemsGrouped : recipeTitle}
                </h2>
              </div>
              <div className="divide-y divide-gray-100 dark:divide-gray-700">
                {recipeItems.map((item) => (
                  <div
                    key={item.id}
                    className={`flex items-center p-4 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors ${
                      item.is_checked ? 'opacity-60' : ''
                    }`}
                  >
                    <button
                      onClick={() => handleToggleItem(item.id)}
                      className="mr-4 text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300 transition-colors"
                    >
                      {item.is_checked ? (
                        <CheckCircle2 className="h-6 w-6" />
                      ) : (
                        <Circle className="h-6 w-6" />
                      )}
                    </button>
                    <div className="flex-1">
                      <p
                        className={`font-medium ${
                          item.is_checked ? 'line-through text-gray-500 dark:text-gray-600' : 'text-gray-900 dark:text-white'
                        }`}
                      >
                        {item.ingredient_name}
                      </p>
                      <p className="text-sm text-gray-500 dark:text-gray-400">
                        {item.quantity} {item.unit}
                        {item.recipe_title && recipeTitle === 'Other' && (
                          <span className="ml-2">
                            {t.shoppingList.from} {item.recipe_title}
                          </span>
                        )}
                      </p>
                    </div>
                    <button
                      onClick={() => handleDeleteItem(item.id)}
                      className="ml-4 text-red-600 hover:text-red-700 transition-colors"
                    >
                      <Trash2 className="h-5 w-5" />
                    </button>
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
