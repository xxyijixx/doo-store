// src/store/useTokenStore.ts
import { create } from 'zustand';

interface TokenState {
  token: string ;
  setToken: (newToken: string) => void;
  clearToken: () => void;
}

export const useTokenStore = create<TokenState>((set) => ({
  token: '',
  setToken: (newToken: string) => set({ token: newToken }),
  clearToken: () => set({ token: '' }),
}));
