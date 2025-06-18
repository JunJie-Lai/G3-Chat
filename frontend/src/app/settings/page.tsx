'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useUserStore, useApiKeyStore } from '../../utils/store';
import UserProfile from '../../components/UserProfile';
import { authService } from '../../services/authService';

export default function SettingsPage() {
  const router = useRouter();
  const { isAuthenticated, user } = useUserStore();
  const { apiKey, setApiKey } = useApiKeyStore();

  const [openAIKey, setOpenAIKey] = useState('');
  const [googleKey, setGoogleKey] = useState('');
  const [anthropicKey, setAnthropicKey] = useState('');
  const [savedMessage, setSavedMessage] = useState('');

  // Redirect to auth page if not authenticated
  useEffect(() => {
    // Use a timeout to ensure the stores are initialized before checking authentication
    const timer = setTimeout(() => {
      // Get the latest authentication state
      const isAuth = authService.isAuthenticated();
      if (!isAuth) {
        router.push('/auth');
      }
    }, 100);

    return () => clearTimeout(timer);
  }, [router]);

  // Load API keys from localStorage
  useEffect(() => {
    const storedOpenAIKey = localStorage.getItem('openai_api_key');
    const storedGoogleKey = localStorage.getItem('google_api_key');
    const storedAnthropicKey = localStorage.getItem('anthropic_api_key');

    if (storedOpenAIKey) setOpenAIKey(storedOpenAIKey);
    if (storedGoogleKey) setGoogleKey(storedGoogleKey);
    if (storedAnthropicKey) setAnthropicKey(storedAnthropicKey);
  }, []);

  const handleSaveKeys = () => {
    // Save API keys to localStorage
    if (openAIKey) {
      localStorage.setItem('openai_api_key', openAIKey);
    } else {
      localStorage.removeItem('openai_api_key');
    }

    if (googleKey) {
      localStorage.setItem('google_api_key', googleKey);
    } else {
      localStorage.removeItem('google_api_key');
    }

    if (anthropicKey) {
      localStorage.setItem('anthropic_api_key', anthropicKey);
    } else {
      localStorage.removeItem('anthropic_api_key');
    }

    // Set the active API key based on which one is available
    if (openAIKey) {
      setApiKey(openAIKey);
    } else if (googleKey) {
      setApiKey(googleKey);
    } else if (anthropicKey) {
      setApiKey(anthropicKey);
    } else {
      setApiKey(null);
    }

    setSavedMessage('API keys saved successfully!');
    setTimeout(() => setSavedMessage(''), 3000);
  };

  if (!isAuthenticated) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <p className="text-gray-500 dark:text-gray-400">Redirecting to login...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      {/* Header */}
      <header className="bg-white dark:bg-gray-800 shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
          <div className="flex items-center">
            <h1 className="text-xl font-semibold text-gray-800 dark:text-white mr-4">G3.Chat</h1>
            <button
              onClick={() => router.push('/chat')}
              className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-gray-200 dark:border-gray-600 dark:hover:bg-gray-600"
            >
              Back to Chat
            </button>
          </div>
          <UserProfile />
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
          <h2 className="text-lg font-medium text-gray-900 dark:text-white mb-6">API Keys</h2>

          {savedMessage && (
            <div className="mb-4 p-3 bg-green-100 text-green-700 rounded-md dark:bg-green-900 dark:text-green-100">
              {savedMessage}
            </div>
          )}

          <div className="space-y-6">
            <div>
              <label htmlFor="openai-key" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                OpenAI API Key
              </label>
              <input
                type="password"
                id="openai-key"
                value={openAIKey}
                onChange={(e) => setOpenAIKey(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                placeholder="sk-..."
              />
              <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                Used for OpenAI models like GPT-4o, GPT-4.1, etc.
              </p>
            </div>

            <div>
              <label htmlFor="google-key" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Google API Key
              </label>
              <input
                type="password"
                id="google-key"
                value={googleKey}
                onChange={(e) => setGoogleKey(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                placeholder="AIza..."
              />
              <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                Used for Google models like Gemini.
              </p>
            </div>

            <div>
              <label htmlFor="anthropic-key" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Anthropic API Key
              </label>
              <input
                type="password"
                id="anthropic-key"
                value={anthropicKey}
                onChange={(e) => setAnthropicKey(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                placeholder="sk-ant-..."
              />
              <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                Used for Anthropic models like Claude.
              </p>
            </div>

            <div className="pt-4">
              <button
                onClick={handleSaveKeys}
                className="px-4 py-2 text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-blue-600 dark:hover:bg-blue-700"
              >
                Save API Keys
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
