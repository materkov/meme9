'use client';

import { createContext, useContext, useState, ReactNode } from 'react';
import Snackbar from '@/components/Snackbar';

interface SnackbarContextType {
  showSnackbar: (message: string) => void;
}

const SnackbarContext = createContext<SnackbarContextType | undefined>(undefined);

export function SnackbarProvider({ children }: { children: ReactNode }) {
  const [message, setMessage] = useState<string>('');
  const [isVisible, setIsVisible] = useState(false);

  const showSnackbar = (newMessage: string) => {
    setMessage(newMessage);
    setIsVisible(true);
  };

  const hideSnackbar = () => {
    setIsVisible(false);
  };

  return (
    <SnackbarContext.Provider value={{ showSnackbar }}>
      {children}
      <Snackbar message={message} isVisible={isVisible} onClose={hideSnackbar} />
    </SnackbarContext.Provider>
  );
}

export function useSnackbar() {
  const context = useContext(SnackbarContext);
  if (context === undefined) {
    throw new Error('useSnackbar must be used within a SnackbarProvider');
  }
  return context;
}

