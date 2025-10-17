import type { TranslationsType } from './en';

export const sk: TranslationsType = {
  // Navigation
  nav: {
    appName: 'Chefly',
    generateRecipe: 'Vytvoriť Recept',
    myRecipes: 'Moje Recepty',
    favorites: 'Obľúbené',
    shoppingList: 'Nákupný Zoznam',
    logout: 'Odhlásiť sa',
    hello: 'Ahoj',
  },

  // Auth
  auth: {
    login: 'Prihlásiť sa',
    register: 'Vytvoriť účet',
    email: 'Email',
    password: 'Heslo',
    confirmPassword: 'Potvrdiť heslo',
    username: 'Používateľské meno',
    signingIn: 'Prihlasovanie...',
    creatingAccount: 'Vytváranie účtu...',
    noAccount: 'Nemáte účet?',
    haveAccount: 'Už máte účet?',
    signUp: 'Zaregistrovať sa',
    signIn: 'Prihlásiť sa',
    signInToAccount: 'Prihláste sa do svojho účtu',
    createAccount: 'Vytvorte si svoj účet',
    emailPlaceholder: 'Zadajte váš email',
    passwordPlaceholder: 'Zadajte vaše heslo',
    usernamePlaceholder: 'Zvoľte používateľské meno',
    confirmPasswordPlaceholder: 'Potvrďte vaše heslo',
    passwordsMismatch: 'Heslá sa nezhodujú',
    passwordTooShort: 'Heslo musí mať aspoň 6 znakov',
    passwordRequirements: 'Heslo musí mať aspoň 8 znakov a obsahovať veľké písmeno, malé písmeno a číslo',
  },

  // Dashboard
  dashboard: {
    welcome: 'Vitajte v Chefly',
    subtitle: 'Generujte úžasné recepty pomocou AI a spravujte svoju osobnú kuchársku knihu',
    generateRecipeTitle: 'Vytvoriť Recept',
    generateRecipeDesc: 'Vytvorte vlastné recepty pomocou AI na základe vašich preferencií, surovín a typu kuchyne',
    myRecipesTitle: 'Moje Recepty',
    myRecipesDesc: 'Zobrazte všetky vaše uložené recepty, spravujte obľúbené a pristupujte k svojej osobnej kuchárskej knihe',
    favoritesTitle: 'Obľúbené',
    favoritesDesc: 'Rýchly prístup k vašim obľúbeným receptom pre jednoduché plánovanie jedál',
    howItWorks: '🍳 Ako to funguje',
    step1: 'Vyberte si typ mäsa, kuchyňu a diétne preferencie',
    step2: 'Kliknite na "Vytvoriť Recept" a nechajte Claude AI vytvoriť vlastný recept',
    step3: 'Zobrazte podrobné ingrediencie, postupné inštrukcie a tipy na varenie',
    step4: 'Uložte si obľúbené a vytvorte si osobnú kuchársku knihu',
  },

  // Generate Recipe
  generate: {
    title: 'Vytvoriť Recept',
    meatType: 'Typ Mäsa',
    cuisine: 'Kuchyňa',
    cookingTime: 'Čas Prípravy',
    difficulty: 'Náročnosť',
    sideIngredients: 'Prílohy',
    dietaryPreferences: 'Diétne Preferencie (Voliteľné)',
    generateButton: 'Vytvoriť Recept',
    generating: 'Vytváranie Receptu...',
    loadingMessage: 'Môže to trvať 10-30 sekúnd, kým Claude AI vytvorí váš recept...',
    successMessage: 'Recept úspešne vytvorený!',
    viewRecipe: 'Zobraziť Recept',
    selectMeat: 'Vyberte Typ Mäsa',
    selectCuisine: 'Vyberte Kuchyňu',
    selectTime: 'Vyberte Čas Prípravy',
    selectDifficulty: 'Vyberte Náročnosť',
  },

  // Dropdown option translations
  options: {
    // Meat types
    meats: {
      'Chicken': 'Kura',
      'Beef': 'Hovädzie mäso',
      'Pork': 'Bravčové mäso',
      'Fish': 'Ryba',
      'Seafood': 'Morské plody',
      'Lamb': 'Jahňacie mäso',
      'Turkey': 'Morka',
      'None (Vegetarian)': 'Žiadne (Vegetariánske)',
    },
    // Cuisines
    cuisines: {
      'Italian': 'Talianska',
      'Mexican': 'Mexická',
      'Chinese': 'Čínska',
      'Indian': 'Indická',
      'Japanese': 'Japonská',
      'Thai': 'Thajská',
      'Mediterranean': 'Stredomorská',
      'American': 'Americká',
      'French': 'Francúzska',
      'Greek': 'Grécka',
      'Korean': 'Kórejská',
      'Vietnamese': 'Vietnamská',
      'Spanish': 'Španielska',
      'Middle Eastern': 'Blízkovýchodná',
    },
    // Side ingredients
    ingredients: {
      'Vegetables': 'Zelenina',
      'Rice': 'Ryža',
      'Pasta': 'Cestoviny',
      'Potatoes': 'Zemiaky',
      'Grains': 'Obilniny',
      'Legumes': 'Strukoviny',
      'Noodles': 'Rezance',
      'Bread': 'Chlieb',
      'Quinoa': 'Quinoa',
      'Couscous': 'Kuskus',
    },
    // Cooking times
    cookingTimes: {
      'quick': 'Rýchly (menej ako 30 min)',
      'medium': 'Stredný (30-60 min)',
      'long': 'Dlhý (viac ako 60 min)',
    },
    // Difficulty levels
    difficulties: {
      'easy': 'Ľahký',
      'medium': 'Stredný',
      'hard': 'Ťažký',
    },
    // Dietary preferences
    dietary: {
      'Vegetarian': 'Vegetariánske',
      'Vegan': 'Vegánske',
      'Gluten-free': 'Bezlepkové',
      'Dairy-free': 'Bez mliečnych výrobkov',
      'Low-carb': 'Nízkosacharidové',
      'Keto': 'Keto',
    },
  },

  // Recipe Detail
  recipe: {
    backToRecipes: 'Späť na Recepty',
    minutes: 'minút',
    ingredients: 'Ingrediencie',
    instructions: 'Postup',
    deleteConfirm: 'Naozaj chcete odstrániť tento recept?',
    deleteFailed: 'Nepodarilo sa odstrániť recept',
    notFound: 'Recept sa nenašiel',
    shareRecipe: 'Zdieľať Recept',
    copyLink: 'Kopírovať Odkaz',
    linkCopied: 'Odkaz skopírovaný do schránky!',
    sharedRecipe: 'Zdieľaný Recept',
  },

  // My Recipes
  myRecipes: {
    title: 'Moje Recepty',
    noRecipes: 'Ešte ste nevytvorili žiadne recepty.',
    generateFirst: 'Vytvoriť Prvý Recept',
    min: 'min',
    search: 'Hľadať recepty...',
    noResults: 'Nenašli sa žiadne recepty podľa vašej požiadavky.',
    clearSearch: 'Vymazať vyhľadávanie',
    showFavorites: 'Zobraziť Len Obľúbené',
  },

  // Profile
  profile: {
    title: 'Môj Profil',
    stats: 'Štatistiky',
    totalRecipes: 'Celkový Počet Receptov',
    favoriteRecipes: 'Obľúbené Recepty',
    memberSince: 'Člen Od',
    accountInfo: 'Informácie o Účte',
    editProfile: 'Upraviť Profil',
    updateProfile: 'Aktualizovať Profil',
    profileUpdated: 'Profil úspešne aktualizovaný!',
    updateFailed: 'Nepodarilo sa aktualizovať profil',
  },

  // Shopping List
  shoppingList: {
    title: 'Nákupný Zoznam',
    noItems: 'Váš nákupný zoznam je prázdny.',
    addFromRecipe: 'Pridajte ingrediencie z receptov a začnite nakupovať!',
    clearChecked: 'Vymazať Označené',
    clearAll: 'Vymazať Všetko',
    confirmClearAll: 'Naozaj chcete vymazať všetky položky?',
    itemsGrouped: 'Položky zoskupené podľa receptu',
    addedToList: 'Ingrediencie pridané do nákupného zoznamu!',
    addToShoppingList: 'Pridať do Nákupného Zoznamu',
    from: 'z',
  },

  // Common
  common: {
    loading: 'Načítavanie...',
    error: 'Chyba',
    success: 'Úspech',
    delete: 'Odstrániť',
    save: 'Uložiť',
    cancel: 'Zrušiť',
    select: 'Vyberte',
  },
};
