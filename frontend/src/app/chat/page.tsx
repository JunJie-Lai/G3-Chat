'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useUserStore, useChatStore } from '../../utils/store';
import { chatService, ChatMessage } from '../../services/chatService';
import { authService } from '../../services/authService';
import ChatSidebar from '../../components/ChatSidebar';
import MessageDisplay from '../../components/MessageDisplay';
import ChatInput from '../../components/ChatInput';
import ModelSelector from '../../components/ModelSelector';
import UserProfile from '../../components/UserProfile';

export default function ChatPage() {
  const router = useRouter();
  const { isAuthenticated } = useUserStore();
  const { messages, addMessage, currentChatId, modelType, model } = useChatStore();
  const [isLoading, setIsLoading] = useState(false);

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

  const handleSendMessage = async (message: string) => {
    setIsLoading(true);

    try {
      // Send message to API
      const response = await chatService.sendMessage({
        id: currentChatId || -1, // -1 for anonymous chat
        model_type: modelType,
        model: model,
        prompt: message
      });

      // Add AI response to chat
      if (response && response.text) {
        addMessage({ role: 'ai', text: response.text });
      }
    } catch (error) {
      console.error('Error sending message:', error);
      // Add error message to chat
      addMessage({ 
        role: 'ai', 
        text: 'Sorry, there was an error processing your request. Please try again.' 
      });
    } finally {
      setIsLoading(false);
    }
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
    <div className="flex h-screen bg-white dark:bg-gray-900">
      {/* Sidebar */}
      <ChatSidebar />

      {/* Main content */}
      <div className="flex flex-col flex-1 h-full">
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-gray-200 dark:border-gray-700">
          <div className="flex items-center space-x-4">
            <h1 className="text-xl font-semibold text-gray-800 dark:text-white">G3.Chat</h1>
            <ModelSelector />
          </div>
          <UserProfile />
        </div>

        {/* Messages */}
        <MessageDisplay messages={messages} />

        {/* Input */}
        <ChatInput onSendMessage={handleSendMessage} />
      </div>
    </div>
  );
}
