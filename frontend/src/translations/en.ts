export const en = {
  // Navigation
  nav: {
    appName: 'Chefly',
    generateRecipe: 'Generate Recipe',
    myRecipes: 'My Recipes',
    favorites: 'Favorites',
    shoppingList: 'Shopping List',
    logout: 'Logout',
    hello: 'Hello',
  },

  // Auth
  auth: {
    login: 'Sign In',
    register: 'Create Account',
    email: 'Email',
    password: 'Password',
    confirmPassword: 'Confirm Password',
    username: 'Username',
    signingIn: 'Signing in...',
    creatingAccount: 'Creating account...',
    noAccount: "Don't have an account?",
    haveAccount: 'Already have an account?',
    signUp: 'Sign up',
    signIn: 'Sign in',
    signInToAccount: 'Sign in to your account',
    createAccount: 'Create your account',
    emailPlaceholder: 'Enter your email',
    passwordPlaceholder: 'Enter your password',
    usernamePlaceholder: 'Choose a username',
    confirmPasswordPlaceholder: 'Confirm your password',
    passwordsMismatch: 'Passwords do not match',
    passwordTooShort: 'Password must be at least 6 characters long',
    passwordRequirements: 'Password must be at least 8 characters and contain uppercase, lowercase, and a number',
  },

  // Dashboard
  dashboard: {
    welcome: 'Welcome to Chefly',
    subtitle: 'Generate amazing recipes with AI and manage your personal cookbook',
    generateRecipeTitle: 'Generate Recipe',
    generateRecipeDesc: 'Create custom recipes with AI based on your preferences, ingredients, and cuisine type',
    myRecipesTitle: 'My Recipes',
    myRecipesDesc: 'View all your saved recipes, manage favorites, and access your personal cookbook',
    favoritesTitle: 'Favorites',
    favoritesDesc: 'Quick access to your favorite recipes for easy meal planning',
    howItWorks: '🍳 How it works',
    step1: 'Choose your meat type, cuisine, and dietary preferences',
    step2: 'Click "Generate Recipe" and let Claude AI create a custom recipe',
    step3: 'View detailed ingredients, step-by-step instructions, and cooking tips',
    step4: 'Save your favorites and build your personal cookbook',
  },

  // Generate Recipe
  generate: {
    title: 'Generate Recipe',
    meatType: 'Meat Type',
    cuisine: 'Cuisine',
    cookingTime: 'Cooking Time',
    difficulty: 'Difficulty',
    sideIngredients: 'Side Ingredients',
    dietaryPreferences: 'Dietary Preferences (Optional)',
    generateButton: 'Generate Recipe',
    generating: 'Generating Recipe...',
    loadingMessage: 'This may take 10-30 seconds while Claude AI creates your recipe...',
    successMessage: 'Recipe generated successfully!',
    viewRecipe: 'View Recipe',
    selectMeat: 'Select Meat Type',
    selectCuisine: 'Select Cuisine',
    selectTime: 'Select Cooking Time',
    selectDifficulty: 'Select Difficulty',
  },

  // Dropdown option translations
  options: {
    // Meat types
    meats: {
      'Chicken': 'Chicken',
      'Beef': 'Beef',
      'Pork': 'Pork',
      'Fish': 'Fish',
      'Seafood': 'Seafood',
      'Lamb': 'Lamb',
      'Turkey': 'Turkey',
      'None (Vegetarian)': 'None (Vegetarian)',
    },
    // Cuisines
    cuisines: {
      'Italian': 'Italian',
      'Mexican': 'Mexican',
      'Chinese': 'Chinese',
      'Indian': 'Indian',
      'Japanese': 'Japanese',
      'Thai': 'Thai',
      'Mediterranean': 'Mediterranean',
      'American': 'American',
      'French': 'French',
      'Greek': 'Greek',
      'Korean': 'Korean',
      'Vietnamese': 'Vietnamese',
      'Spanish': 'Spanish',
      'Middle Eastern': 'Middle Eastern',
    },
    // Side ingredients
    ingredients: {
      'Vegetables': 'Vegetables',
      'Rice': 'Rice',
      'Pasta': 'Pasta',
      'Potatoes': 'Potatoes',
      'Grains': 'Grains',
      'Legumes': 'Legumes',
      'Noodles': 'Noodles',
      'Bread': 'Bread',
      'Quinoa': 'Quinoa',
      'Couscous': 'Couscous',
    },
    // Cooking times
    cookingTimes: {
      'quick': 'Quick (under 30 min)',
      'medium': 'Medium (30-60 min)',
      'long': 'Long (over 60 min)',
    },
    // Difficulty levels
    difficulties: {
      'easy': 'Easy',
      'medium': 'Medium',
      'hard': 'Hard',
    },
    // Dietary preferences
    dietary: {
      'Vegetarian': 'Vegetarian',
      'Vegan': 'Vegan',
      'Gluten-free': 'Gluten-free',
      'Dairy-free': 'Dairy-free',
      'Low-carb': 'Low-carb',
      'Keto': 'Keto',
    },
  },

  // Recipe Detail
  recipe: {
    backToRecipes: 'Back to Recipes',
    minutes: 'minutes',
    ingredients: 'Ingredients',
    instructions: 'Instructions',
    deleteConfirm: 'Are you sure you want to delete this recipe?',
    deleteFailed: 'Failed to delete recipe',
    notFound: 'Recipe not found',
    shareRecipe: 'Share Recipe',
    copyLink: 'Copy Link',
    linkCopied: 'Link copied to clipboard!',
    sharedRecipe: 'Shared Recipe',
  },

  // My Recipes
  myRecipes: {
    title: 'My Recipes',
    noRecipes: "You haven't generated any recipes yet.",
    generateFirst: 'Generate Your First Recipe',
    min: 'min',
    search: 'Search recipes...',
    noResults: 'No recipes found matching your search.',
    clearSearch: 'Clear search',
    showFavorites: 'Show Favorites Only',
  },

  // Profile
  profile: {
    title: 'My Profile',
    stats: 'Statistics',
    totalRecipes: 'Total Recipes',
    favoriteRecipes: 'Favorite Recipes',
    memberSince: 'Member Since',
    accountInfo: 'Account Information',
    editProfile: 'Edit Profile',
    updateProfile: 'Update Profile',
    profileUpdated: 'Profile updated successfully!',
    updateFailed: 'Failed to update profile',
  },

  // Shopping List
  shoppingList: {
    title: 'Shopping List',
    noItems: 'Your shopping list is empty.',
    addFromRecipe: 'Add ingredients from your recipes to start shopping!',
    clearChecked: 'Clear Checked Items',
    clearAll: 'Clear All',
    confirmClearAll: 'Are you sure you want to clear all items?',
    itemsGrouped: 'Items grouped by recipe',
    addedToList: 'Ingredients added to shopping list!',
    addToShoppingList: 'Add to Shopping List',
    from: 'from',
  },

  // Common
  common: {
    loading: 'Loading...',
    error: 'Error',
    success: 'Success',
    delete: 'Delete',
    save: 'Save',
    cancel: 'Cancel',
    select: 'Select',
  },
};

export type TranslationsType = typeof en;
