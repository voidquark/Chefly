import React from 'react';
import { useNavigate } from 'react-router-dom';
import { useLanguage } from '../contexts/LanguageContext';
import { usePageTitle } from '../hooks/usePageTitle';
import { PlusCircle, BookOpen, Heart } from 'lucide-react';

export const Dashboard: React.FC = () => {
  const navigate = useNavigate();
  const { t } = useLanguage();
  usePageTitle('Dashboard');

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-2">{t.dashboard.welcome}</h1>
      <p className="text-gray-600 dark:text-gray-300 mb-8">
        {t.dashboard.subtitle}
      </p>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <button
          onClick={() => navigate('/generate')}
          className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md hover:shadow-xl hover:scale-105 transition-all duration-300 text-left group animate-fadeIn"
        >
          <div className="flex items-center mb-4">
            <div className="p-3 bg-primary-100 dark:bg-primary-900 rounded-lg group-hover:bg-primary-200 dark:group-hover:bg-primary-800 transition-colors">
              <PlusCircle className="h-8 w-8 text-primary-600 dark:text-primary-400" />
            </div>
          </div>
          <h3 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">{t.dashboard.generateRecipeTitle}</h3>
          <p className="text-gray-600 dark:text-gray-300 text-sm">
            {t.dashboard.generateRecipeDesc}
          </p>
        </button>

        <button
          onClick={() => navigate('/recipes')}
          className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md hover:shadow-xl hover:scale-105 transition-all duration-300 text-left group animate-fadeIn"
          style={{ animationDelay: '0.1s' }}
        >
          <div className="flex items-center mb-4">
            <div className="p-3 bg-green-100 dark:bg-green-900 rounded-lg group-hover:bg-green-200 dark:group-hover:bg-green-800 transition-colors">
              <BookOpen className="h-8 w-8 text-green-600 dark:text-green-400" />
            </div>
          </div>
          <h3 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">{t.dashboard.myRecipesTitle}</h3>
          <p className="text-gray-600 dark:text-gray-300 text-sm">
            {t.dashboard.myRecipesDesc}
          </p>
        </button>

        <button
          onClick={() => navigate('/recipes')}
          className="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-md hover:shadow-xl hover:scale-105 transition-all duration-300 text-left group animate-fadeIn"
          style={{ animationDelay: '0.2s' }}
        >
          <div className="flex items-center mb-4">
            <div className="p-3 bg-red-100 dark:bg-red-900 rounded-lg group-hover:bg-red-200 dark:group-hover:bg-red-800 transition-colors">
              <Heart className="h-8 w-8 text-red-600 dark:text-red-400" />
            </div>
          </div>
          <h3 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">{t.dashboard.favoritesTitle}</h3>
          <p className="text-gray-600 dark:text-gray-300 text-sm">
            {t.dashboard.favoritesDesc}
          </p>
        </button>
      </div>

      <div className="mt-12 bg-primary-50 dark:bg-primary-900/20 border border-primary-200 dark:border-primary-800 rounded-lg p-6 animate-slideUp">
        <h2 className="text-xl font-semibold text-primary-900 dark:text-primary-200 mb-2">{t.dashboard.howItWorks}</h2>
        <ol className="list-decimal list-inside space-y-2 text-primary-800 dark:text-primary-300">
          <li>{t.dashboard.step1}</li>
          <li>{t.dashboard.step2}</li>
          <li>{t.dashboard.step3}</li>
          <li>{t.dashboard.step4}</li>
        </ol>
      </div>
    </div>
  );
};
