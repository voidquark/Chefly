import React from 'react';

export const SkeletonShoppingList: React.FC = () => {
  return (
    <div className="max-w-4xl mx-auto animate-pulse">
      {/* Header skeleton */}
      <div className="flex items-center justify-between mb-6">
        <div className="h-8 bg-gray-300 dark:bg-gray-700 rounded w-48"></div>
        <div className="h-10 bg-gray-300 dark:bg-gray-700 rounded w-32"></div>
      </div>

      {/* Shopping list items skeleton */}
      <div className="space-y-6">
        {[1, 2].map((section) => (
          <div key={section} className="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden">
            {/* Section header */}
            <div className="bg-primary-50 dark:bg-primary-900/20 px-4 py-3 border-b border-primary-100 dark:border-primary-800">
              <div className="h-5 bg-gray-300 dark:bg-gray-700 rounded w-40"></div>
            </div>

            {/* Items */}
            <div className="divide-y divide-gray-100 dark:divide-gray-700">
              {[1, 2, 3, 4].map((item) => (
                <div key={item} className="flex items-center p-4 gap-4">
                  <div className="w-6 h-6 bg-gray-300 dark:bg-gray-700 rounded-full"></div>
                  <div className="flex-1 space-y-2">
                    <div className="h-4 bg-gray-300 dark:bg-gray-700 rounded w-3/4"></div>
                    <div className="h-3 bg-gray-300 dark:bg-gray-700 rounded w-1/2"></div>
                  </div>
                  <div className="w-6 h-6 bg-gray-300 dark:bg-gray-700 rounded"></div>
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
