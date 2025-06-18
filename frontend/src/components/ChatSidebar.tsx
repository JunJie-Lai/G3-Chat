import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useChatStore } from '../utils/store';
import { chatService } from '../services/chatService';

export default function ChatSidebar() {
  const router = useRouter();
  const { chatTitles, setChatTitles, setCurrentChatId, currentChatId } = useChatStore();

  useEffect(() => {
    // Fetch chat titles when component mounts
    const fetchChatTitles = async () => {
      try {
        const titles = await chatService.getChatTitles();
        setChatTitles(titles);
      } catch (error) {
        console.error('Error fetching chat titles:', error);
      }
    };

    fetchChatTitles();
  }, [setChatTitles]);

  const handleChatSelect = async (id: number) => {
    setCurrentChatId(id);
    
    try {
      const chatHistory = await chatService.getChatHistory(id);
      useChatStore.getState().setMessages(chatHistory.chatHistory);
    } catch (error) {
      console.error('Error fetching chat history:', error);
    }
  };

  const handleNewChat = () => {
    setCurrentChatId(null);
    useChatStore.getState().setMessages([]);
  };

  return (
    <div className="w-64 h-full bg-gray-100 dark:bg-gray-900 border-r border-gray-200 dark:border-gray-700 flex flex-col">
      <div className="p-4 border-b border-gray-200 dark:border-gray-700">
        <button
          onClick={handleNewChat}
          className="w-full px-4 py-2 text-sm font-medium text-white bg-blue-500 rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 dark:bg-blue-600 dark:hover:bg-blue-700"
        >
          New Chat
        </button>
      </div>
      
      <div className="flex-1 overflow-y-auto p-2">
        <h3 className="px-2 py-1 text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">
          Recent Chats
        </h3>
        
        {chatTitles.length === 0 ? (
          <div className="text-center text-gray-500 dark:text-gray-400 text-sm p-4">
            No chats yet
          </div>
        ) : (
          <ul className="space-y-1 mt-1">
            {chatTitles.map((chat) => (
              <li key={chat.id}>
                <button
                  onClick={() => handleChatSelect(chat.id)}
                  className={`w-full text-left px-3 py-2 text-sm rounded-md ${
                    currentChatId === chat.id
                      ? 'bg-blue-100 text-blue-900 dark:bg-blue-900 dark:text-blue-100'
                      : 'text-gray-700 hover:bg-gray-200 dark:text-gray-200 dark:hover:bg-gray-800'
                  }`}
                >
                  <div className="truncate">{chat.title}</div>
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}