import React from 'react';
import { ChefHat } from 'lucide-react';

export const LoadingSplash: React.FC = () => {
  return (
    <div className="fixed inset-0 bg-gradient-to-br from-primary-50 to-primary-100 dark:from-gray-900 dark:to-gray-800 flex items-center justify-center z-50">
      <div className="text-center animate-fadeIn">
        {/* Logo */}
        <div className="mb-6 animate-bounce">
          <ChefHat className="h-20 w-20 text-primary-600 dark:text-primary-400 mx-auto" />
        </div>

        {/* App Name */}
        <h1 className="text-4xl font-bold text-gray-900 dark:text-white mb-4">
          Chefly
        </h1>

        {/* Loading Spinner */}
        <div className="flex justify-center mb-4">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 dark:border-primary-400"></div>
        </div>

        {/* Loading Text */}
        <p className="text-gray-600 dark:text-gray-300 text-sm">
          Loading your recipes...
        </p>
      </div>
    </div>
  );
};
