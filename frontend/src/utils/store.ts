import { create } from 'zustand';
import { User } from '@/services/authService';
import { ChatMessage, ChatTitle } from '@/services/chatService';

// Helper function to check if code is running in browser
const isBrowser = typeof window !== 'undefined';

// User store
interface UserState {
  user: User | null;
  isAuthenticated: boolean;
  setUser: (user: User | null) => void;
  setAuthenticated: (isAuthenticated: boolean) => void;
  logout: () => void;
}

export const useUserStore = create<UserState>((set) => ({
  user: null,
  isAuthenticated: false,
  setUser: (user) => set({ user }),
  setAuthenticated: (isAuthenticated) => set({ isAuthenticated }),
  logout: () => set({ user: null, isAuthenticated: false }),
}));

// Chat store
interface ChatState {
  chatTitles: ChatTitle[];
  currentChatId: number | null;
  messages: ChatMessage[];
  modelType: 'Google' | 'OpenAI' | 'Anthropic';
  model: string;
  setChatTitles: (chatTitles: ChatTitle[]) => void;
  setCurrentChatId: (id: number | null) => void;
  setMessages: (messages: ChatMessage[]) => void;
  addMessage: (message: ChatMessage) => void;
  setModelType: (modelType: 'Google' | 'OpenAI' | 'Anthropic') => void;
  setModel: (model: string) => void;
}

export const useChatStore = create<ChatState>((set) => ({
  chatTitles: [],
  currentChatId: null,
  messages: [],
  modelType: 'OpenAI',
  model: 'gpt-4o',
  setChatTitles: (chatTitles) => set({ chatTitles }),
  setCurrentChatId: (currentChatId) => set({ currentChatId }),
  setMessages: (messages) => set({ messages }),
  addMessage: (message) => set((state) => ({ messages: [...state.messages, message] })),
  setModelType: (modelType) => set({ modelType }),
  setModel: (model) => set({ model }),
}));

// API key store
interface ApiKeyState {
  apiKey: string | null;
  setApiKey: (apiKey: string | null) => void;
}

export const useApiKeyStore = create<ApiKeyState>((set) => ({
  apiKey: isBrowser ? localStorage.getItem('api_key') : null,
  setApiKey: (apiKey) => {
    if (isBrowser) {
      if (apiKey) {
        localStorage.setItem('api_key', apiKey);
      } else {
        localStorage.removeItem('api_key');
      }
    }
    set({ apiKey });
  },
}));

// Initialize stores from localStorage on an app load
export const initializeStores = () => {
  // Only run in browser environment
  if (!isBrowser) return;

  // Initialize user store
  const userStr = localStorage.getItem('user');
  const sessionToken = localStorage.getItem('session_token');

  if (userStr && sessionToken) {
    const user = JSON.parse(userStr);
    useUserStore.getState().setUser(user);
    useUserStore.getState().setAuthenticated(true);
  }

  // Initialize API key store
  const apiKey = localStorage.getItem('api_key');
  if (apiKey) {
    useApiKeyStore.getState().setApiKey(apiKey);
  }
};
