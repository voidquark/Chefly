import React from 'react';

export const SkeletonRecipeDetail: React.FC = () => {
  return (
    <div className="max-w-4xl mx-auto animate-pulse">
      {/* Back button skeleton */}
      <div className="h-6 bg-gray-300 dark:bg-gray-700 rounded w-32 mb-6"></div>

      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-lg overflow-hidden">
        {/* Image skeleton */}
        <div className="w-full h-64 bg-gray-300 dark:bg-gray-700"></div>

        <div className="p-8">
          {/* Title skeleton */}
          <div className="h-8 bg-gray-300 dark:bg-gray-700 rounded mb-4 w-2/3"></div>

          {/* Description skeleton */}
          <div className="space-y-2 mb-6">
            <div className="h-4 bg-gray-300 dark:bg-gray-700 rounded w-full"></div>
            <div className="h-4 bg-gray-300 dark:bg-gray-700 rounded w-5/6"></div>
          </div>

          {/* Meta info skeleton */}
          <div className="flex gap-4 mb-6">
            <div className="h-6 bg-gray-300 dark:bg-gray-700 rounded w-24"></div>
            <div className="h-6 bg-gray-300 dark:bg-gray-700 rounded w-24"></div>
            <div className="h-6 bg-gray-300 dark:bg-gray-700 rounded w-20"></div>
            <div className="h-6 bg-gray-300 dark:bg-gray-700 rounded w-20"></div>
          </div>

          {/* Ingredients section */}
          <div className="mb-8">
            <div className="h-6 bg-gray-300 dark:bg-gray-700 rounded w-32 mb-4"></div>
            <div className="space-y-2">
              {[1, 2, 3, 4, 5].map((i) => (
                <div key={i} className="h-4 bg-gray-300 dark:bg-gray-700 rounded w-full"></div>
              ))}
            </div>
          </div>

          {/* Steps section */}
          <div>
            <div className="h-6 bg-gray-300 dark:bg-gray-700 rounded w-32 mb-4"></div>
            <div className="space-y-4">
              {[1, 2, 3].map((i) => (
                <div key={i} className="flex gap-4">
                  <div className="flex-shrink-0 w-8 h-8 bg-gray-300 dark:bg-gray-700 rounded-full"></div>
                  <div className="flex-1 space-y-2">
                    <div className="h-4 bg-gray-300 dark:bg-gray-700 rounded w-full"></div>
                    <div className="h-4 bg-gray-300 dark:bg-gray-700 rounded w-4/5"></div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
