import { useEffect } from 'react';

export const usePageTitle = (title: string) => {
  useEffect(() => {
    const baseTitle = 'Chefly';
    document.title = title ? `${title} - ${baseTitle}` : baseTitle;

    // Cleanup: restore base title when component unmounts
    return () => {
      document.title = baseTitle;
    };
  }, [title]);
};
